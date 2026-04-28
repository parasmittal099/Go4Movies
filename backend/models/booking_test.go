package models

import (
	"testing"
	"time"

	"gorm.io/gorm"
)

func seedBookingData(t *testing.T, db *gorm.DB) (User, Showtime, Seat) {
	t.Helper()

	loc := Location{Zipcode: "33301", City: "Miami", State: "FL"}
	db.Create(&loc)
	theater := Theater{Name: "BT", LocationID: loc.ID, TotalScreens: 1}
	db.Create(&theater)
	screen := Screen{TheaterID: theater.ID, Name: "S1", TotalRows: 5, TotalCols: 11, ScreenType: "Standard"}
	db.Create(&screen)
	seat := Seat{ScreenID: screen.ID, RowLabel: "A", ColNumber: 1, SeatType: "Premium", BasePrice: 18.0}
	db.Create(&seat)
	movie := Movie{Title: "BM", Language: "English", DurationMin: 90, IsActive: true}
	db.Create(&movie)
	st := Showtime{MovieID: movie.ID, ScreenID: screen.ID, ShowDate: "2026-02-17", StartTime: "10:00", EndTime: "11:30", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true}
	db.Create(&st)
	user := User{Email: "booker@t.com", Username: "booker", Password: "p", FullName: "Booker"}
	db.Create(&user)

	return user, st, seat
}

func TestBooking_Create(t *testing.T) {
	db := setupModelsDB(t)
	user, st, _ := seedBookingData(t, db)

	booking := Booking{
		UserID: user.ID, ShowtimeID: st.ID,
		BookingRef: "REF-001", Status: "CONFIRMED",
		TotalAmount: 18.0, PaymentStatus: "PAID", BookedAt: time.Now(),
	}
	if err := db.Create(&booking).Error; err != nil {
		t.Fatalf("failed to create booking: %v", err)
	}
	if booking.ID == 0 {
		t.Error("expected non-zero ID")
	}
}

func TestBooking_UniqueRef(t *testing.T) {
	db := setupModelsDB(t)
	user, st, _ := seedBookingData(t, db)

	db.Create(&Booking{UserID: user.ID, ShowtimeID: st.ID, BookingRef: "DUP-001", Status: "CONFIRMED", TotalAmount: 18.0, PaymentStatus: "PAID", BookedAt: time.Now()})
	err := db.Create(&Booking{UserID: user.ID, ShowtimeID: st.ID, BookingRef: "DUP-001", Status: "CONFIRMED", TotalAmount: 18.0, PaymentStatus: "PAID", BookedAt: time.Now()}).Error
	if err == nil {
		t.Error("expected unique constraint violation on booking_ref")
	}
}

func TestBookingSeat_Create(t *testing.T) {
	db := setupModelsDB(t)
	user, st, seat := seedBookingData(t, db)

	booking := Booking{UserID: user.ID, ShowtimeID: st.ID, BookingRef: "BS-001", Status: "CONFIRMED", TotalAmount: 18.0, PaymentStatus: "PAID", BookedAt: time.Now()}
	db.Create(&booking)

	bs := BookingSeat{BookingID: booking.ID, SeatID: seat.ID, ShowtimeID: st.ID, SeatPrice: 18.0}
	if err := db.Create(&bs).Error; err != nil {
		t.Fatalf("failed to create booking_seat: %v", err)
	}
	if bs.ID == 0 {
		t.Error("expected non-zero ID")
	}
}

func TestBooking_HasManySeats(t *testing.T) {
	db := setupModelsDB(t)
	user, st, seat := seedBookingData(t, db)

	seat2 := Seat{ScreenID: seat.ScreenID, RowLabel: "A", ColNumber: 2, SeatType: "Premium", BasePrice: 18.0}
	db.Create(&seat2)

	booking := Booking{UserID: user.ID, ShowtimeID: st.ID, BookingRef: "MULTI-001", Status: "CONFIRMED", TotalAmount: 36.0, PaymentStatus: "PAID", BookedAt: time.Now()}
	db.Create(&booking)
	db.Create(&BookingSeat{BookingID: booking.ID, SeatID: seat.ID, ShowtimeID: st.ID, SeatPrice: 18.0})
	db.Create(&BookingSeat{BookingID: booking.ID, SeatID: seat2.ID, ShowtimeID: st.ID, SeatPrice: 18.0})

	var fetched Booking
	db.Preload("BookingSeats").First(&fetched, booking.ID)
	if len(fetched.BookingSeats) != 2 {
		t.Errorf("expected 2 booking seats, got %d", len(fetched.BookingSeats))
	}
}

func TestBooking_DefaultStatus(t *testing.T) {
	db := setupModelsDB(t)
	user, st, _ := seedBookingData(t, db)

	booking := Booking{UserID: user.ID, ShowtimeID: st.ID, BookingRef: "DEF-001", TotalAmount: 18.0, BookedAt: time.Now()}
	db.Create(&booking)

	var fetched Booking
	db.First(&fetched, booking.ID)
	if fetched.Status != "PENDING" {
		t.Errorf("expected default status PENDING, got %q", fetched.Status)
	}
	if fetched.PaymentStatus != "UNPAID" {
		t.Errorf("expected default payment_status UNPAID, got %q", fetched.PaymentStatus)
	}
}

func TestPayment_Create(t *testing.T) {
	db := setupModelsDB(t)
	user, st, _ := seedBookingData(t, db)

	booking := Booking{UserID: user.ID, ShowtimeID: st.ID, BookingRef: "PAY-001", Status: "CONFIRMED", TotalAmount: 18.0, PaymentStatus: "PAID", BookedAt: time.Now()}
	db.Create(&booking)

	txnID := "TXN-123"
	payment := Payment{
		BookingID: booking.ID, Amount: 18.0,
		PaymentMethod: "CARD", TransactionID: &txnID,
		Status: "COMPLETED", InitiatedAt: time.Now(),
	}
	if err := db.Create(&payment).Error; err != nil {
		t.Fatalf("failed to create payment: %v", err)
	}
	if payment.ID == 0 {
		t.Error("expected non-zero ID")
	}
}

func TestQRTicket_UniqueTicketCode(t *testing.T) {
	db := setupModelsDB(t)
	user, st, _ := seedBookingData(t, db)

	booking1 := Booking{UserID: user.ID, ShowtimeID: st.ID, BookingRef: "QR-REF-1", Status: "CONFIRMED", TotalAmount: 18.0, PaymentStatus: "PAID", BookedAt: time.Now()}
	booking2 := Booking{UserID: user.ID, ShowtimeID: st.ID, BookingRef: "QR-REF-2", Status: "CONFIRMED", TotalAmount: 18.0, PaymentStatus: "PAID", BookedAt: time.Now()}
	db.Create(&booking1)
	db.Create(&booking2)

	db.Create(&QRTicket{BookingID: booking1.ID, TicketCode: "same-code"})
	err := db.Create(&QRTicket{BookingID: booking2.ID, TicketCode: "same-code"}).Error
	if err == nil {
		t.Error("expected unique constraint violation on ticket_code")
	}
}

func TestQRTicket_OnePerBooking(t *testing.T) {
	db := setupModelsDB(t)
	user, st, _ := seedBookingData(t, db)

	booking := Booking{UserID: user.ID, ShowtimeID: st.ID, BookingRef: "QR-ONE-BOOKING", Status: "CONFIRMED", TotalAmount: 18.0, PaymentStatus: "PAID", BookedAt: time.Now()}
	db.Create(&booking)

	db.Create(&QRTicket{BookingID: booking.ID, TicketCode: "code-1"})
	err := db.Create(&QRTicket{BookingID: booking.ID, TicketCode: "code-2"}).Error
	if err == nil {
		t.Error("expected unique constraint violation on booking_id")
	}
}
