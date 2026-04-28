package models

import "time"

type Booking struct {
	ID             uint          `gorm:"primaryKey" json:"id"`
	UserID         uint          `gorm:"index;not null" json:"user_id"`
	User           User          `gorm:"foreignKey:UserID" json:"-"`
	ShowtimeID     uint          `gorm:"index;not null" json:"showtime_id"`
	Showtime       Showtime      `gorm:"foreignKey:ShowtimeID" json:"showtime,omitempty"`
	BookingRef     string        `gorm:"uniqueIndex;not null" json:"booking_ref"`
	Status         string        `gorm:"not null;default:PENDING" json:"status"`
	TotalAmount    float64       `gorm:"not null" json:"total_amount"`
	ConvenienceFee float64       `gorm:"not null;default:0" json:"convenience_fee"`
	TaxAmount      float64       `gorm:"not null;default:0" json:"tax_amount"`
	PaymentStatus  string        `gorm:"not null;default:UNPAID" json:"payment_status"`
	BookedAt       time.Time     `gorm:"not null" json:"booked_at"`
	ExpiresAt      *time.Time    `gorm:"index" json:"expires_at,omitempty"`
	CancelledAt    *time.Time    `json:"cancelled_at,omitempty"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
	BookingSeats   []BookingSeat `gorm:"foreignKey:BookingID" json:"booking_seats,omitempty"`
	QRTicket       *QRTicket     `gorm:"foreignKey:BookingID" json:"qr_ticket,omitempty"`
}

type BookingSeat struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	BookingID  uint    `gorm:"uniqueIndex:idx_bs_unique;not null" json:"booking_id"`
	Booking    Booking `gorm:"foreignKey:BookingID;constraint:OnDelete:CASCADE" json:"-"`
	SeatID     uint    `gorm:"uniqueIndex:idx_bs_unique;uniqueIndex:idx_showtime_seat;not null" json:"seat_id"`
	Seat       Seat    `gorm:"foreignKey:SeatID" json:"seat,omitempty"`
	ShowtimeID uint    `gorm:"uniqueIndex:idx_showtime_seat;not null" json:"showtime_id"`
	SeatPrice  float64 `gorm:"not null" json:"seat_price"`
}

type Payment struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	BookingID       uint       `gorm:"index;not null" json:"booking_id"`
	Booking         Booking    `gorm:"foreignKey:BookingID" json:"-"`
	Amount          float64    `gorm:"not null" json:"amount"`
	PaymentMethod   string     `gorm:"not null" json:"payment_method"`
	TransactionID   *string    `gorm:"uniqueIndex" json:"transaction_id,omitempty"`
	Status          string     `gorm:"not null;default:INITIATED" json:"status"`
	InitiatedAt     time.Time  `gorm:"not null" json:"initiated_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	GatewayResponse *string    `json:"gateway_response,omitempty"`
}
