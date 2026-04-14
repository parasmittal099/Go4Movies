package workers

import (
	"log"
	"time"

	"gorm.io/gorm"
)

// StartBookingExpiry launches a background goroutine that periodically
// marks PENDING bookings whose hold window has elapsed as EXPIRED.
// interval controls how often the sweep runs (e.g. 1 minute).
func StartBookingExpiry(db *gorm.DB, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			expireStaleBookings(db)
		}
	}()
}

func expireStaleBookings(db *gorm.DB) {
	result := db.
		Table("bookings").
		Where("status = ? AND expires_at IS NOT NULL AND expires_at <= ?", "PENDING", time.Now()).
		Updates(map[string]interface{}{
			"status":         "EXPIRED",
			"payment_status": "EXPIRED",
		})
	if result.Error != nil {
		log.Printf("[expiry] error expiring bookings: %v", result.Error)
	} else if result.RowsAffected > 0 {
		log.Printf("[expiry] expired %d stale PENDING bookings", result.RowsAffected)
	}
}
