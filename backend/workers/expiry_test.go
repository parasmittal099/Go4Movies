package workers

import (
	"testing"
	"time"

	"github.com/parasmittal099/backend-project/database"
	"github.com/parasmittal099/backend-project/models"
	"github.com/parasmittal099/backend-project/testutil"
)

func TestExpireStaleBookings(t *testing.T) {
	db := testutil.SetupTestDB(t)

	user := models.User{Email: "exp@test.com", Username: "expuser", Password: "pass", FullName: "Exp User"}
	db.Create(&user)

	loc := models.Location{Zipcode: "10001", City: "NYC", State: "NY"}
	db.Create(&loc)
	addr := "1 St"
	theater := models.Theater{Name: "T", LocationID: loc.ID, Address: &addr, TotalScreens: 1}
	db.Create(&theater)
	screen := models.Screen{TheaterID: theater.ID, Name: "S1", TotalRows: 1, TotalCols: 2, ScreenType: "Standard"}
	db.Create(&screen)
	movie := models.Movie{Title: "M", Language: "English", DurationMin: 90, IsActive: true}
	db.Create(&movie)
	st := models.Showtime{
		MovieID: movie.ID, ScreenID: screen.ID,
		ShowDate: "2026-03-01", StartTime: "10:00", EndTime: "12:00",
		Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true,
	}
	db.Create(&st)

	past := time.Now().Add(-15 * time.Minute)
	future := time.Now().Add(15 * time.Minute)

	expired := models.Booking{
		UserID: user.ID, ShowtimeID: st.ID, BookingRef: "EXP-001",
		Status: "PENDING", TotalAmount: 50, PaymentStatus: "UNPAID",
		BookedAt: past, ExpiresAt: &past,
	}
	db.Create(&expired)

	stillValid := models.Booking{
		UserID: user.ID, ShowtimeID: st.ID, BookingRef: "VALID-001",
		Status: "PENDING", TotalAmount: 50, PaymentStatus: "UNPAID",
		BookedAt: time.Now(), ExpiresAt: &future,
	}
	db.Create(&stillValid)

	confirmed := models.Booking{
		UserID: user.ID, ShowtimeID: st.ID, BookingRef: "CONF-001",
		Status: "CONFIRMED", TotalAmount: 50, PaymentStatus: "PAID",
		BookedAt: time.Now(),
	}
	db.Create(&confirmed)

	expireStaleBookings(database.DB)

	var b1, b2, b3 models.Booking
	db.First(&b1, expired.ID)
	db.First(&b2, stillValid.ID)
	db.First(&b3, confirmed.ID)

	if b1.Status != "EXPIRED" {
		t.Errorf("expired booking: want EXPIRED, got %s", b1.Status)
	}
	if b2.Status != "PENDING" {
		t.Errorf("valid pending booking: want PENDING, got %s", b2.Status)
	}
	if b3.Status != "CONFIRMED" {
		t.Errorf("confirmed booking: want CONFIRMED, got %s", b3.Status)
	}
}
