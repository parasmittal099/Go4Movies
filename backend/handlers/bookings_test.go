package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parasmittal099/backend-project/database"
	"github.com/parasmittal099/backend-project/middleware"
	"github.com/parasmittal099/backend-project/models"
	"github.com/parasmittal099/backend-project/testutil"
)

const bookingsTestSecret = "test-jwt-secret"

func setupBookingsRouter() *gin.Engine {
	r := gin.New()
	r.GET("/bookings", middleware.OptionalJWTAuth(bookingsTestSecret), GetUserBookings)
	return r
}

func strPtrBookings(s string) *string { return &s }

func bookingsAuthRequest(t *testing.T, method, url string, userID uint) *http.Request {
	t.Helper()
	req, _ := http.NewRequest(method, url, nil)
	tok, err := middleware.GenerateToken(userID, bookingsTestSecret)
	if err != nil {
		t.Fatalf("token: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	return req
}

func seedBookingHistoryData(t *testing.T) (models.User, models.User, []models.Booking) {
	t.Helper()

	user1 := models.User{Email: "u1@test.com", Username: "u1", Password: "hash", FullName: "User One"}
	user2 := models.User{Email: "u2@test.com", Username: "u2", Password: "hash", FullName: "User Two"}
	database.DB.Create(&user1)
	database.DB.Create(&user2)

	loc := models.Location{Zipcode: "33101", City: "Miami", State: "FL"}
	database.DB.Create(&loc)

	addr := "123 Main St"
	theater := models.Theater{Name: "AMC Miami", LocationID: loc.ID, Address: &addr, TotalScreens: 1}
	database.DB.Create(&theater)

	screen := models.Screen{TheaterID: theater.ID, Name: "Screen 1", TotalRows: 5, TotalCols: 10, ScreenType: "IMAX"}
	database.DB.Create(&screen)

	movie := models.Movie{Title: "Demo Movie", PosterURL: strPtrBookings("https://img.test/poster.jpg"), Language: "English", DurationMin: 120, IsActive: true}
	database.DB.Create(&movie)

	st := models.Showtime{
		MovieID: movie.ID, ScreenID: screen.ID,
		ShowDate: "2026-06-15", StartTime: "19:30", EndTime: "21:30",
		Language: "English", Format: "IMAX", PriceMultiplier: 1.0, IsActive: true,
	}
	database.DB.Create(&st)

	seat1 := models.Seat{ScreenID: screen.ID, RowLabel: "G", ColNumber: 12, SeatType: "Premium", BasePrice: 17}
	seat2 := models.Seat{ScreenID: screen.ID, RowLabel: "G", ColNumber: 13, SeatType: "Premium", BasePrice: 17}
	seat3 := models.Seat{ScreenID: screen.ID, RowLabel: "G", ColNumber: 14, SeatType: "Premium", BasePrice: 17}
	database.DB.Create(&seat1)
	database.DB.Create(&seat2)
	database.DB.Create(&seat3)

	oldBookedAt := time.Now().Add(-2 * time.Hour)
	newBookedAt := time.Now().Add(-1 * time.Hour)
	b1 := models.Booking{
		UserID: user1.ID, ShowtimeID: st.ID, BookingRef: "G4M-old", Status: "CONFIRMED",
		TotalAmount: 19, ConvenienceFee: 2, TaxAmount: 0, PaymentStatus: "PAID", BookedAt: oldBookedAt,
	}
	b2 := models.Booking{
		UserID: user1.ID, ShowtimeID: st.ID, BookingRef: "G4M-new", Status: "CONFIRMED",
		TotalAmount: 34, ConvenienceFee: 2, TaxAmount: 2.88, PaymentStatus: "PAID", BookedAt: newBookedAt,
	}
	bOther := models.Booking{
		UserID: user2.ID, ShowtimeID: st.ID, BookingRef: "G4M-other", Status: "CONFIRMED",
		TotalAmount: 34, ConvenienceFee: 2, TaxAmount: 2.88, PaymentStatus: "PAID", BookedAt: newBookedAt,
	}
	database.DB.Create(&b1)
	database.DB.Create(&b2)
	database.DB.Create(&bOther)

	database.DB.Create(&models.BookingSeat{BookingID: b1.ID, SeatID: seat3.ID, ShowtimeID: st.ID, SeatPrice: 17})
	database.DB.Create(&models.BookingSeat{BookingID: b2.ID, SeatID: seat1.ID, ShowtimeID: st.ID, SeatPrice: 17})
	database.DB.Create(&models.BookingSeat{BookingID: b2.ID, SeatID: seat2.ID, ShowtimeID: st.ID, SeatPrice: 17})
	database.DB.Create(&models.BookingSeat{BookingID: bOther.ID, SeatID: seat3.ID, ShowtimeID: st.ID, SeatPrice: 17})

	return user1, user2, []models.Booking{b1, b2, bOther}
}

func TestGetUserBookings_SuccessWithJWT(t *testing.T) {
	testutil.SetupTestDB(t)
	u1, _, _ := seedBookingHistoryData(t)
	r := setupBookingsRouter()

	w := httptest.NewRecorder()
	req := bookingsAuthRequest(t, "GET", "/bookings", u1.ID)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	bookings := resp["bookings"].([]any)
	if len(bookings) != 2 {
		t.Fatalf("expected 2 bookings for user1, got %d", len(bookings))
	}

	first := bookings[0].(map[string]any)
	second := bookings[1].(map[string]any)
	if first["booking_ref"] != "G4M-new" || second["booking_ref"] != "G4M-old" {
		t.Fatalf("expected bookings ordered by booked_at desc, got first=%v second=%v", first["booking_ref"], second["booking_ref"])
	}
	if first["movie_title"] != "Demo Movie" {
		t.Fatalf("expected movie_title mapped, got %v", first["movie_title"])
	}
	if first["movie_poster"] != "https://img.test/poster.jpg" {
		t.Fatalf("expected movie_poster mapped, got %v", first["movie_poster"])
	}
	if first["theater_name"] != "AMC Miami" || first["screen_name"] != "Screen 1" || first["screen_type"] != "IMAX" {
		t.Fatalf("expected theater/screen mapping, got theater=%v screen=%v type=%v", first["theater_name"], first["screen_name"], first["screen_type"])
	}
	seats := first["seats"].([]any)
	if len(seats) != 2 {
		t.Fatalf("expected 2 seats in latest booking, got %d", len(seats))
	}
}

func TestGetUserBookings_SuccessWithQueryParam(t *testing.T) {
	testutil.SetupTestDB(t)
	u1, _, _ := seedBookingHistoryData(t)
	r := setupBookingsRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", fmt.Sprintf("/bookings?user_id=%d", u1.ID), nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	bookings := resp["bookings"].([]any)
	if len(bookings) != 2 {
		t.Fatalf("expected 2 bookings for user1 via query param, got %d", len(bookings))
	}
}

func TestGetUserBookings_Empty(t *testing.T) {
	testutil.SetupTestDB(t)
	r := setupBookingsRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/bookings?user_id=9999", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	bookings := resp["bookings"].([]any)
	if len(bookings) != 0 {
		t.Fatalf("expected empty bookings array, got %d", len(bookings))
	}
}

func TestGetUserBookings_MissingUserID(t *testing.T) {
	testutil.SetupTestDB(t)
	r := setupBookingsRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/bookings", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetUserBookings_InvalidQueryParam(t *testing.T) {
	testutil.SetupTestDB(t)
	r := setupBookingsRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/bookings?user_id=abc", nil)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestGetUserBookings_IsolationByUser(t *testing.T) {
	testutil.SetupTestDB(t)
	_, u2, _ := seedBookingHistoryData(t)
	r := setupBookingsRouter()

	w := httptest.NewRecorder()
	req := bookingsAuthRequest(t, "GET", "/bookings", u2.ID)
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	bookings := resp["bookings"].([]any)
	if len(bookings) != 1 {
		t.Fatalf("expected only one booking for user2, got %d", len(bookings))
	}
	ref := bookings[0].(map[string]any)["booking_ref"]
	if ref != "G4M-other" {
		t.Fatalf("expected G4M-other booking for user2, got %v", ref)
	}
}
