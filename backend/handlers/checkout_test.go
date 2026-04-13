package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parasmittal099/backend-project/database"
	"github.com/parasmittal099/backend-project/models"
	"github.com/parasmittal099/backend-project/testutil"
)

func setupCheckoutRouter() *gin.Engine {
	r := gin.New()
	r.POST("/checkout/preview", PreviewCheckout)
	r.POST("/checkout/confirm", ConfirmCheckout)
	return r
}

type checkoutSeed struct {
	User     models.User
	Showtime models.Showtime
	Seat1    models.Seat
	Seat2    models.Seat
}

func seedCheckoutData(t *testing.T) checkoutSeed {
	t.Helper()

	user := models.User{
		Email: "checkout@test.com", Username: "checkoutuser",
		Password: "secret12", FullName: "Checkout User",
	}
	if err := database.DB.Create(&user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}

	loc := models.Location{Zipcode: "33101", City: "Miami", State: "FL"}
	database.DB.Create(&loc)

	addr := "1 Test St"
	theater := models.Theater{Name: "Test Plex", LocationID: loc.ID, Address: &addr, TotalScreens: 1}
	database.DB.Create(&theater)

	screen := models.Screen{TheaterID: theater.ID, Name: "S1", TotalRows: 5, TotalCols: 10, ScreenType: "Standard"}
	database.DB.Create(&screen)

	movie := models.Movie{Title: "Checkout Movie", Language: "English", DurationMin: 100, IsActive: true}
	database.DB.Create(&movie)

	today := time.Now().Format("2006-01-02")
	st := models.Showtime{
		MovieID: movie.ID, ScreenID: screen.ID,
		ShowDate: today, StartTime: "10:00", EndTime: "12:00",
		Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true,
	}
	database.DB.Create(&st)

	seat1 := models.Seat{ScreenID: screen.ID, RowLabel: "A", ColNumber: 1, SeatType: "Regular", BasePrice: 100}
	seat2 := models.Seat{ScreenID: screen.ID, RowLabel: "A", ColNumber: 2, SeatType: "Regular", BasePrice: 100}
	database.DB.Create(&seat1)
	database.DB.Create(&seat2)

	return checkoutSeed{User: user, Showtime: st, Seat1: seat1, Seat2: seat2}
}

func postJSONCheckout(t *testing.T, r *gin.Engine, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	b, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	return w
}

func TestPreviewCheckout_Success_NoDiscount(t *testing.T) {
	testutil.SetupTestDB(t)
	td := seedCheckoutData(t)
	r := setupCheckoutRouter()

	w := postJSONCheckout(t, r, "/checkout/preview", map[string]any{
		"user_id":     td.User.ID,
		"showtime_id": td.Showtime.ID,
		"seat_ids":    []uint{td.Seat1.ID},
	})

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)

	totals := resp["totals"].(map[string]any)
	// subtotal 100, fee 2, tax on 102 = 8.16, total 110.16
	if totals["subtotal"].(float64) != 100 {
		t.Errorf("subtotal want 100 got %v", totals["subtotal"])
	}
	if totals["convenience_fee"].(float64) != 2 {
		t.Errorf("fee want 2 got %v", totals["convenience_fee"])
	}
	if totals["tax_amount"].(float64) != 8.16 {
		t.Errorf("tax want 8.16 got %v", totals["tax_amount"])
	}
	if totals["total_due"].(float64) != 110.16 {
		t.Errorf("total_due want 110.16 got %v", totals["total_due"])
	}
	if totals["discount_amount"].(float64) != 0 {
		t.Errorf("discount want 0 got %v", totals["discount_amount"])
	}
}

func TestPreviewCheckout_Mock100_ZeroTotal(t *testing.T) {
	testutil.SetupTestDB(t)
	td := seedCheckoutData(t)
	r := setupCheckoutRouter()

	w := postJSONCheckout(t, r, "/checkout/preview", map[string]any{
		"user_id":        td.User.ID,
		"showtime_id":    td.Showtime.ID,
		"seat_ids":       []uint{td.Seat1.ID},
		"discount_code":  discountCodeMock100,
	})

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	totals := resp["totals"].(map[string]any)
	if totals["total_due"].(float64) != 0 {
		t.Errorf("total_due want 0 got %v", totals["total_due"])
	}
	pre := 100.0 + 2.0 + 8.16
	if totals["discount_amount"].(float64) != pre {
		t.Errorf("discount_amount want %v got %v", pre, totals["discount_amount"])
	}
}

func TestPreviewCheckout_TwoSeats(t *testing.T) {
	testutil.SetupTestDB(t)
	td := seedCheckoutData(t)
	r := setupCheckoutRouter()

	w := postJSONCheckout(t, r, "/checkout/preview", map[string]any{
		"user_id":     td.User.ID,
		"showtime_id": td.Showtime.ID,
		"seat_ids":    []uint{td.Seat1.ID, td.Seat2.ID},
	})
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	items := resp["line_items"].([]any)
	if len(items) != 2 {
		t.Errorf("want 2 line items, got %d", len(items))
	}
}

func TestPreviewCheckout_UserNotFound(t *testing.T) {
	testutil.SetupTestDB(t)
	td := seedCheckoutData(t)
	r := setupCheckoutRouter()

	w := postJSONCheckout(t, r, "/checkout/preview", map[string]any{
		"user_id":     99999,
		"showtime_id": td.Showtime.ID,
		"seat_ids":    []uint{td.Seat1.ID},
	})
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestPreviewCheckout_ShowtimeNotFound(t *testing.T) {
	testutil.SetupTestDB(t)
	td := seedCheckoutData(t)
	r := setupCheckoutRouter()

	w := postJSONCheckout(t, r, "/checkout/preview", map[string]any{
		"user_id":     td.User.ID,
		"showtime_id": 99999,
		"seat_ids":    []uint{td.Seat1.ID},
	})
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestPreviewCheckout_InactiveShowtime(t *testing.T) {
	testutil.SetupTestDB(t)
	td := seedCheckoutData(t)
	database.DB.Model(&td.Showtime).Update("is_active", false)
	r := setupCheckoutRouter()

	w := postJSONCheckout(t, r, "/checkout/preview", map[string]any{
		"user_id":     td.User.ID,
		"showtime_id": td.Showtime.ID,
		"seat_ids":    []uint{td.Seat1.ID},
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPreviewCheckout_InvalidSeatWrongScreen(t *testing.T) {
	testutil.SetupTestDB(t)
	td := seedCheckoutData(t)

	var screen models.Screen
	database.DB.First(&screen, td.Showtime.ScreenID)
	theaterID := screen.TheaterID
	os := models.Screen{TheaterID: theaterID, Name: "Screen 2", TotalRows: 2, TotalCols: 2, ScreenType: "Standard"}
	database.DB.Create(&os)
	badSeat := models.Seat{ScreenID: os.ID, RowLabel: "Z", ColNumber: 9, SeatType: "Regular", BasePrice: 50}
	database.DB.Create(&badSeat)

	r := setupCheckoutRouter()
	w := postJSONCheckout(t, r, "/checkout/preview", map[string]any{
		"user_id":     td.User.ID,
		"showtime_id": td.Showtime.ID,
		"seat_ids":    []uint{badSeat.ID},
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPreviewCheckout_DuplicateSeatID(t *testing.T) {
	testutil.SetupTestDB(t)
	td := seedCheckoutData(t)
	r := setupCheckoutRouter()

	w := postJSONCheckout(t, r, "/checkout/preview", map[string]any{
		"user_id":     td.User.ID,
		"showtime_id": td.Showtime.ID,
		"seat_ids":    []uint{td.Seat1.ID, td.Seat1.ID},
	})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPreviewCheckout_SeatAlreadyBooked(t *testing.T) {
	testutil.SetupTestDB(t)
	td := seedCheckoutData(t)

	booking := models.Booking{
		UserID: td.User.ID, ShowtimeID: td.Showtime.ID, BookingRef: "EXIST-1",
		Status: "CONFIRMED", TotalAmount: 100, ConvenienceFee: 0, TaxAmount: 0,
		PaymentStatus: "PAID", BookedAt: time.Now(),
	}
	database.DB.Create(&booking)
	database.DB.Create(&models.BookingSeat{
		BookingID: booking.ID, SeatID: td.Seat1.ID, ShowtimeID: td.Showtime.ID, SeatPrice: 100,
	})

	r := setupCheckoutRouter()
	w := postJSONCheckout(t, r, "/checkout/preview", map[string]any{
		"user_id":     td.User.ID,
		"showtime_id": td.Showtime.ID,
		"seat_ids":    []uint{td.Seat1.ID},
	})
	if w.Code != http.StatusConflict {
		t.Fatalf("expected 409, got %d: %s", w.Code, w.Body.String())
	}
}

func TestPreviewCheckout_EmptyBody(t *testing.T) {
	testutil.SetupTestDB(t)
	seedCheckoutData(t)
	r := setupCheckoutRouter()

	w := postJSONCheckout(t, r, "/checkout/preview", map[string]any{})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestConfirmCheckout_CreatesBookingAndPayment(t *testing.T) {
	testutil.SetupTestDB(t)
	td := seedCheckoutData(t)
	r := setupCheckoutRouter()

	w := postJSONCheckout(t, r, "/checkout/confirm", map[string]any{
		"user_id":        td.User.ID,
		"showtime_id":    td.Showtime.ID,
		"seat_ids":       []uint{td.Seat1.ID},
		"discount_code":  discountCodeMock100,
	})
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["booking_ref"] == nil || resp["booking_ref"] == "" {
		t.Error("missing booking_ref")
	}

	var count int64
	database.DB.Model(&models.Booking{}).Count(&count)
	if count != 1 {
		t.Errorf("want 1 booking, got %d", count)
	}
	database.DB.Model(&models.BookingSeat{}).Count(&count)
	if count != 1 {
		t.Errorf("want 1 booking_seat, got %d", count)
	}
	database.DB.Model(&models.Payment{}).Count(&count)
	if count != 1 {
		t.Errorf("want 1 payment, got %d", count)
	}

	var b models.Booking
	database.DB.First(&b)
	if b.Status != "CONFIRMED" || b.PaymentStatus != "PAID" {
		t.Errorf("booking status want CONFIRMED/PAID, got %s/%s", b.Status, b.PaymentStatus)
	}
}

func TestConfirmCheckout_ThenPreviewConflict(t *testing.T) {
	testutil.SetupTestDB(t)
	td := seedCheckoutData(t)
	r := setupCheckoutRouter()

	body := map[string]any{
		"user_id":     td.User.ID,
		"showtime_id": td.Showtime.ID,
		"seat_ids":    []uint{td.Seat2.ID},
	}
	w1 := postJSONCheckout(t, r, "/checkout/confirm", body)
	if w1.Code != http.StatusCreated {
		t.Fatalf("confirm: %d %s", w1.Code, w1.Body.String())
	}

	w2 := postJSONCheckout(t, r, "/checkout/preview", body)
	if w2.Code != http.StatusConflict {
		t.Fatalf("second preview want 409, got %d: %s", w2.Code, w2.Body.String())
	}
}
