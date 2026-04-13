package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parasmittal099/backend-project/database"
	"github.com/parasmittal099/backend-project/models"
)

// Checkout defaults (no env for now; adjust later as needed).
const (
	checkoutTaxRate        = 0.08  // 8% on subtotal + convenience fee
	checkoutConvenienceFee = 2.00  // flat per order
	discountCodeMock100    = "MOCK100"
)

type checkoutRequest struct {
	UserID        uint     `json:"user_id" binding:"required"`
	ShowtimeID    uint     `json:"showtime_id" binding:"required"`
	SeatIDs       []uint   `json:"seat_ids" binding:"required,min=1"`
	DiscountCode  string   `json:"discount_code"`
}

type checkoutLineItem struct {
	SeatID    uint    `json:"seat_id"`
	RowLabel  string  `json:"row_label"`
	ColNumber int     `json:"col_number"`
	SeatType  string  `json:"seat_type"`
	UnitPrice float64 `json:"unit_price"`
}

type checkoutTotals struct {
	Subtotal        float64 `json:"subtotal"`
	ConvenienceFee  float64 `json:"convenience_fee"`
	TaxAmount       float64 `json:"tax_amount"`
	DiscountCode    string  `json:"discount_code,omitempty"`
	DiscountAmount  float64 `json:"discount_amount"`
	TotalDue        float64 `json:"total_due"`
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

// validateAndQuote loads data, checks seats belong to showtime and are free, returns quote or HTTP error via gin.
func validateAndQuote(req checkoutRequest, c *gin.Context) (*checkoutQuote, bool) {
	var user models.User
	if err := database.DB.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return nil, false
	}

	var showtime models.Showtime
	if err := database.DB.First(&showtime, req.ShowtimeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Showtime not found"})
		return nil, false
	}
	if !showtime.IsActive {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Showtime is not active"})
		return nil, false
	}

	seen := make(map[uint]struct{})
	for _, id := range req.SeatIDs {
		if _, dup := seen[id]; dup {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Duplicate seat_id in request"})
			return nil, false
		}
		seen[id] = struct{}{}
	}

	var seats []models.Seat
	if err := database.DB.Where("id IN ? AND screen_id = ?", req.SeatIDs, showtime.ScreenID).Find(&seats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load seats"})
		return nil, false
	}
	if len(seats) != len(req.SeatIDs) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "One or more seat_ids are invalid for this showtime"})
		return nil, false
	}

	var conflict int64
	err := database.DB.Table("booking_seats").
		Joins("JOIN bookings ON bookings.id = booking_seats.booking_id").
		Where("booking_seats.showtime_id = ? AND booking_seats.seat_id IN ? AND bookings.status IN ?",
			showtime.ID, req.SeatIDs, []string{"PENDING", "CONFIRMED"}).
		Count(&conflict).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check seat availability"})
		return nil, false
	}
	if conflict > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "One or more seats are no longer available"})
		return nil, false
	}

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

	quote := &checkoutQuote{
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
	return quote, true
}

// POST /api/v1/checkout/preview
func PreviewCheckout(c *gin.Context) {
	var req checkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	quote, ok := validateAndQuote(req, c)
	if !ok {
		return
	}
	c.JSON(http.StatusOK, quote)
}

func newBookingRef() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return "G4M-" + hex.EncodeToString(b)
}

// POST /api/v1/checkout/confirm
func ConfirmCheckout(c *gin.Context) {
	var req checkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	quote, ok := validateAndQuote(req, c)
	if !ok {
		return
	}

	now := time.Now()
	booking := models.Booking{
		UserID:         quote.UserID,
		ShowtimeID:     quote.ShowtimeID,
		BookingRef:     newBookingRef(),
		Status:         "PENDING",
		TotalAmount:    quote.Totals.TotalDue,
		ConvenienceFee: quote.Totals.ConvenienceFee,
		TaxAmount:      quote.Totals.TaxAmount,
		PaymentStatus:  "UNPAID",
		BookedAt:       now,
	}
	if err := database.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking"})
		return
	}

	for _, li := range quote.LineItems {
		bs := models.BookingSeat{
			BookingID:  booking.ID,
			SeatID:     li.SeatID,
			ShowtimeID: quote.ShowtimeID,
			SeatPrice:  li.UnitPrice,
		}
		if err := database.DB.Create(&bs).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save seat assignment", "booking_id": booking.ID})
			return
		}
	}

	// Mock checkout: always record payment as completed (real gateway would gate this).
	payment := models.Payment{
		BookingID:     booking.ID,
		Amount:        quote.Totals.TotalDue,
		PaymentMethod: "MOCK",
		Status:        "COMPLETED",
		InitiatedAt:   now,
		CompletedAt:   &now,
	}
	if err := database.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record payment", "booking_id": booking.ID})
		return
	}

	bookingUpdates := map[string]interface{}{
		"status":          "CONFIRMED",
		"payment_status":  "PAID",
		"total_amount":    quote.Totals.TotalDue,
		"convenience_fee": quote.Totals.ConvenienceFee,
		"tax_amount":      quote.Totals.TaxAmount,
	}
	if err := database.DB.Model(&booking).Updates(bookingUpdates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to finalize booking", "booking_id": booking.ID})
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
