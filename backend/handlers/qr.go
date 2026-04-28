package handlers

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parasmittal099/backend-project/database"
	"github.com/parasmittal099/backend-project/models"
	"gorm.io/gorm"
)

// GET /api/v1/bookings/by-ticket?ticket_code=<code>
func GetBookingByTicketCode(c *gin.Context) {
	ticketCode := strings.TrimSpace(c.Query("ticket_code"))
	if ticketCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_code query parameter is required"})
		return
	}

	var qrTicket models.QRTicket
	if err := database.DB.Where("ticket_code = ? AND is_active = ?", ticketCode, true).First(&qrTicket).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to resolve ticket"})
		return
	}

	if qrTicket.ExpiresAt != nil && qrTicket.ExpiresAt.Before(time.Now()) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	var booking models.Booking
	if err := database.DB.
		Where("id = ?", qrTicket.BookingID).
		Preload("Showtime.Movie").
		Preload("Showtime.Screen.Theater").
		Preload("BookingSeats.Seat").
		First(&booking).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch booking"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"booking":     mapBookingHistoryItem(booking),
		"ticket_code": qrTicket.TicketCode,
		"qr_value":    "G4M:" + qrTicket.TicketCode,
	})
}
