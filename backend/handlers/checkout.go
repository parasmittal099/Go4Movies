package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parasmittal099/backend-project/database"
	"github.com/parasmittal099/backend-project/models"
	"gorm.io/gorm"
)

const (
	checkoutTaxRate        = 0.08
	checkoutConvenienceFee = 2.00
	discountCodeMock100    = "MOCK100"
	pendingHoldMinutes     = 10
)

type checkoutRequest struct {
	UserID       uint   `json:"user_id" binding:"required"`
	ShowtimeID   uint   `json:"showtime_id" binding:"required"`
	SeatIDs      []uint `json:"seat_ids" binding:"required,min=1"`
	DiscountCode string `json:"discount_code"`
}

type checkoutLineItem struct {
	SeatID    uint    `json:"seat_id"`
	RowLabel  string  `json:"row_label"`
	ColNumber int     `json:"col_number"`
	SeatType  string  `json:"seat_type"`
	UnitPrice float64 `json:"unit_price"`
}

type checkoutTotals struct {
	Subtotal       float64 `json:"subtotal"`
	ConvenienceFee float64 `json:"convenience_fee"`
	TaxAmount      float64 `json:"tax_amount"`
	DiscountCode   string  `json:"discount_code,omitempty"`
	DiscountAmount float64 `json:"discount_amount"`
	TotalDue       float64 `json:"total_due"`
}

type checkoutQuote struct {
	ShowtimeID uint               `json:"showtime_id"`
	UserID     uint               `json:"user_id"`
	LineItems  []checkoutLineItem `json:"line_items"`
	Totals     checkoutTotals     `json:"totals"`
}

func roundMoney(v float64) float64 {
	return math.Round(v*100) / 100
}

func seatUnitPrice(seat models.Seat, mult float64) float64 {
	return roundMoney(seat.BasePrice * mult)
}

var errConflict = errors.New("seat_conflict")

// validateRequest performs stateless + lightweight DB checks that
// don't need to be inside the serialised transaction.
func validateRequest(req checkoutRequest, c *gin.Context) (*models.User, *models.Showtime, []models.Seat, bool) {
	var user models.User
	if err := database.DB.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return nil, nil, nil, false
	}

	var showtime models.Showtime
	if err := database.DB.First(&showtime, req.ShowtimeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Showtime not found"})
		return nil, nil, nil, false
	}
	if !showtime.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Showtime is not active"})
		return nil, nil, nil, false
	}

	seen := make(map[uint]struct{})
	for _, id := range req.SeatIDs {
		if _, dup := seen[id]; dup {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate seat_id in request"})
			return nil, nil, nil, false
		}
		seen[id] = struct{}{}
	}

	var seats []models.Seat
	if err := database.DB.Where("id IN ? AND screen_id = ?", req.SeatIDs, showtime.ScreenID).Find(&seats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load seats"})
		return nil, nil, nil, false
	}
	if len(seats) != len(req.SeatIDs) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "One or more seat_ids are invalid for this showtime"})
		return nil, nil, nil, false
	}

	return &user, &showtime, seats, true
}

// buildQuote computes the price breakdown (pure calculation, no DB).
func buildQuote(user models.User, showtime models.Showtime, seats []models.Seat, req checkoutRequest) *checkoutQuote {
	seatByID := make(map[uint]models.Seat, len(seats))
	for _, s := range seats {
		seatByID[s.ID] = s
	}

	lineItems := make([]checkoutLineItem, 0, len(req.SeatIDs))
	var subtotal float64
	for _, id := range req.SeatIDs {
		seat := seatByID[id]
		price := seatUnitPrice(seat, showtime.PriceMultiplier)
		subtotal += price
		lineItems = append(lineItems, checkoutLineItem{
			SeatID:    seat.ID,
			RowLabel:  seat.RowLabel,
			ColNumber: seat.ColNumber,
			SeatType:  seat.SeatType,
			UnitPrice: price,
		})
	}
	subtotal = roundMoney(subtotal)

	taxable := subtotal + checkoutConvenienceFee
	tax := roundMoney(taxable * checkoutTaxRate)
	preDiscount := roundMoney(subtotal + checkoutConvenienceFee + tax)

	discountAmt := 0.0
	code := req.DiscountCode
	if code == discountCodeMock100 {
		discountAmt = preDiscount
	}

	totalDue := roundMoney(preDiscount - discountAmt)
	if totalDue < 0 {
		totalDue = 0
	}

	return &checkoutQuote{
		ShowtimeID: showtime.ID,
		UserID:     user.ID,
		LineItems:  lineItems,
		Totals: checkoutTotals{
			Subtotal:       subtotal,
			ConvenienceFee: checkoutConvenienceFee,
			TaxAmount:      tax,
			DiscountCode:   code,
			DiscountAmount: discountAmt,
			TotalDue:       totalDue,
		},
	}
}

// checkConflicts checks whether any of the requested seats are already
// held by a non-expired PENDING or CONFIRMED booking. Must be called
// inside a transaction for safety.
func checkConflicts(tx *gorm.DB, showtimeID uint, seatIDs []uint) error {
	var conflict int64
	err := tx.Table("booking_seats").
		Joins("JOIN bookings ON bookings.id = booking_seats.booking_id").
		Where("booking_seats.showtime_id = ? AND booking_seats.seat_id IN ?", showtimeID, seatIDs).
		Where("bookings.status IN ?", []string{"PENDING", "CONFIRMED"}).
		Where("bookings.status = 'CONFIRMED' OR bookings.expires_at IS NULL OR bookings.expires_at > ?", time.Now()).
		Count(&conflict).Error
	if err != nil {
		return err
	}
	if conflict > 0 {
		return errConflict
	}
	return nil
}

// POST /api/v1/checkout/preview
func PreviewCheckout(c *gin.Context) {
	var req checkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, showtime, seats, ok := validateRequest(req, c)
	if !ok {
		return
	}

	if err := checkConflicts(database.DB, showtime.ID, req.SeatIDs); err != nil {
		if errors.Is(err, errConflict) {
			c.JSON(http.StatusConflict, gin.H{"error": "One or more seats are no longer available"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check seat availability"})
		}
		return
	}

	quote := buildQuote(*user, *showtime, seats, req)
	c.JSON(http.StatusOK, quote)
}

func newBookingRef() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return "G4M-" + hex.EncodeToString(b)
}

// POST /api/v1/checkout/confirm (JWT-protected)
func ConfirmCheckout(c *gin.Context) {
	var req checkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if uid, exists := c.Get("user_id"); exists {
		req.UserID = uid.(uint)
	}

	user, showtime, seats, ok := validateRequest(req, c)
	if !ok {
		return
	}

	quote := buildQuote(*user, *showtime, seats, req)

	now := time.Now()
	expires := now.Add(pendingHoldMinutes * time.Minute)

	var booking models.Booking
	var payment models.Payment

	txErr := database.DB.Transaction(func(tx *gorm.DB) error {
		// Re-check availability inside the transaction (serialised by SQLite's write lock)
		if err := checkConflicts(tx, showtime.ID, req.SeatIDs); err != nil {
			return err
		}

		booking = models.Booking{
			UserID:         quote.UserID,
			ShowtimeID:     quote.ShowtimeID,
			BookingRef:     newBookingRef(),
			Status:         "PENDING",
			TotalAmount:    quote.Totals.TotalDue,
			ConvenienceFee: quote.Totals.ConvenienceFee,
			TaxAmount:      quote.Totals.TaxAmount,
			PaymentStatus:  "UNPAID",
			BookedAt:       now,
			ExpiresAt:      &expires,
		}
		if err := tx.Create(&booking).Error; err != nil {
			return err
		}

		bookingSeats := make([]models.BookingSeat, 0, len(quote.LineItems))
		for _, li := range quote.LineItems {
			bookingSeats = append(bookingSeats, models.BookingSeat{
				BookingID:  booking.ID,
				SeatID:     li.SeatID,
				ShowtimeID: quote.ShowtimeID,
				SeatPrice:  li.UnitPrice,
			})
		}
		if err := tx.Create(&bookingSeats).Error; err != nil {
			return err
		}

		payment = models.Payment{
			BookingID:     booking.ID,
			Amount:        quote.Totals.TotalDue,
			PaymentMethod: "MOCK",
			Status:        "COMPLETED",
			InitiatedAt:   now,
			CompletedAt:   &now,
		}
		if err := tx.Create(&payment).Error; err != nil {
			return err
		}

		if err := tx.Model(&booking).Updates(map[string]interface{}{
			"status":         "CONFIRMED",
			"payment_status": "PAID",
			"expires_at":     nil,
		}).Error; err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		if errors.Is(txErr, errConflict) {
			c.JSON(http.StatusConflict, gin.H{"error": "One or more seats are no longer available"})
			return
		}
		if strings.Contains(txErr.Error(), "UNIQUE constraint failed") {
			c.JSON(http.StatusConflict, gin.H{"error": "One or more seats are no longer available"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to complete booking"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     "Booking confirmed",
		"booking_id":  booking.ID,
		"booking_ref": booking.BookingRef,
		"quote":       quote,
		"payment_id":  payment.ID,
	})
}
