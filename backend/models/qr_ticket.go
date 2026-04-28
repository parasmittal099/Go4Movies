package models

import "time"

type QRTicket struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	BookingID  uint       `gorm:"uniqueIndex;not null" json:"booking_id"`
	Booking    Booking    `gorm:"foreignKey:BookingID;constraint:OnDelete:CASCADE" json:"-"`
	TicketCode string     `gorm:"uniqueIndex;not null" json:"ticket_code"`
	IsActive   bool       `gorm:"not null;default:true" json:"is_active"`
	ExpiresAt  *time.Time `gorm:"index" json:"expires_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}
