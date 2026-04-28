package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/parasmittal099/backend-project/database"
	"github.com/parasmittal099/backend-project/models"
)

type bookingSeatHistory struct {
	SeatID    uint    `json:"seat_id"`
	RowLabel  string  `json:"row_label"`
	ColNumber int     `json:"col_number"`
	SeatType  string  `json:"seat_type"`
	SeatPrice float64 `json:"seat_price"`
}

type bookingHistoryItem struct {
	ID             uint                 `json:"id"`
	BookingRef     string               `json:"booking_ref"`
	Status         string               `json:"status"`
	TotalAmount    float64              `json:"total_amount"`
	ConvenienceFee float64              `json:"convenience_fee"`
	TaxAmount      float64              `json:"tax_amount"`
	PaymentStatus  string               `json:"payment_status"`
	BookedAt       string               `json:"booked_at"`
	MovieTitle     string               `json:"movie_title"`
	MoviePoster    string               `json:"movie_poster"`
	TheaterName    string               `json:"theater_name"`
	ScreenName     string               `json:"screen_name"`
	ScreenType     string               `json:"screen_type"`
	ShowDate       string               `json:"show_date"`
	StartTime      string               `json:"start_time"`
	Format         string               `json:"format"`
	Language       string               `json:"language"`
	Seats          []bookingSeatHistory `json:"seats"`
}

// GET /api/v1/bookings — accepts JWT (user_id from token) or ?user_id= query param
func GetUserBookings(c *gin.Context) {
	var userID uint
	if uid, exists := c.Get("user_id"); exists {
		userID = uid.(uint)
	} else if param := c.Query("user_id"); param != "" {
		id, err := strconv.Atoi(param)
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "user_id must be a positive integer"})
			return
		}
		userID = uint(id)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id query parameter or Authorization header is required"})
		return
	}

	var bookings []models.Booking
	if err := database.DB.
		Where("user_id = ?", userID).
		Preload("Showtime.Movie").
		Preload("Showtime.Screen.Theater").
		Preload("BookingSeats.Seat").
		Order("booked_at DESC").
		Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bookings"})
		return
	}

	respBookings := make([]bookingHistoryItem, 0, len(bookings))
	for _, b := range bookings {
		item := bookingHistoryItem{
			ID:             b.ID,
			BookingRef:     b.BookingRef,
			Status:         b.Status,
			TotalAmount:    b.TotalAmount,
			ConvenienceFee: b.ConvenienceFee,
			TaxAmount:      b.TaxAmount,
			PaymentStatus:  b.PaymentStatus,
			BookedAt:       b.BookedAt.UTC().Format("2006-01-02T15:04:05Z"),
			ShowDate:       b.Showtime.ShowDate,
			StartTime:      b.Showtime.StartTime,
			Format:         b.Showtime.Format,
			Language:       b.Showtime.Language,
		}

		if b.Showtime.Movie.Title != "" {
			item.MovieTitle = b.Showtime.Movie.Title
		}
		if b.Showtime.Movie.PosterURL != nil {
			item.MoviePoster = *b.Showtime.Movie.PosterURL
		}
		item.ScreenName = b.Showtime.Screen.Name
		item.ScreenType = b.Showtime.Screen.ScreenType
		item.TheaterName = b.Showtime.Screen.Theater.Name

		seats := make([]bookingSeatHistory, 0, len(b.BookingSeats))
		for _, bs := range b.BookingSeats {
			seats = append(seats, bookingSeatHistory{
				SeatID:    bs.SeatID,
				RowLabel:  bs.Seat.RowLabel,
				ColNumber: bs.Seat.ColNumber,
				SeatType:  bs.Seat.SeatType,
				SeatPrice: bs.SeatPrice,
			})
		}
		item.Seats = seats
		respBookings = append(respBookings, item)
	}

	c.JSON(http.StatusOK, gin.H{"bookings": respBookings})
}

