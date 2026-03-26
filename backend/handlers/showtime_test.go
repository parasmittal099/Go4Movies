package handlers

import (
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

func setupShowtimeRouter() *gin.Engine {
	r := gin.New()
	r.GET("/movies/:id/showtimes", GetMovieShowtimes)
	r.GET("/seats", GetShowtimeSeats)
	return r
}

type showtimeTestData struct {
	location  models.Location
	theater1  models.Theater
	theater2  models.Theater
	screen1   models.Screen
	screen2   models.Screen
	movie     models.Movie
	showtime1 models.Showtime
	showtime2 models.Showtime
	showtime3 models.Showtime
}

func seedShowtimeTestData(t *testing.T) showtimeTestData {
	t.Helper()
	loc := models.Location{Zipcode: "33101", City: "Miami", State: "FL"}
	database.DB.Create(&loc)

	locOther := models.Location{Zipcode: "32801", City: "Orlando", State: "FL"}
	database.DB.Create(&locOther)

	addr1 := "123 Main St, Miami"
	theater1 := models.Theater{Name: "AMC Miami", LocationID: loc.ID, Address: &addr1, TotalScreens: 1}
	database.DB.Create(&theater1)

	addr2 := "456 Beach Rd, Miami"
	theater2 := models.Theater{Name: "Regal Miami", LocationID: loc.ID, Address: &addr2, TotalScreens: 1}
	database.DB.Create(&theater2)

	screen1 := models.Screen{TheaterID: theater1.ID, Name: "Screen 1", TotalRows: 5, TotalCols: 10, ScreenType: "IMAX"}
	database.DB.Create(&screen1)

	screen2 := models.Screen{TheaterID: theater2.ID, Name: "Screen A", TotalRows: 4, TotalCols: 8, ScreenType: "Standard"}
	database.DB.Create(&screen2)

	movie := models.Movie{Title: "Test Movie", Language: "English", DurationMin: 120, IsActive: true, Genre: strPtr("Action")}
	database.DB.Create(&movie)

	otherMovie := models.Movie{Title: "Other Movie", Language: "English", DurationMin: 90, IsActive: true}
	database.DB.Create(&otherMovie)

	today := time.Now().Format("2006-01-02")
	tomorrow := time.Now().Add(24 * time.Hour).Format("2006-01-02")

	st1 := models.Showtime{MovieID: movie.ID, ScreenID: screen1.ID, ShowDate: today, StartTime: "10:00", EndTime: "12:00", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true}
	database.DB.Create(&st1)

	st2 := models.Showtime{MovieID: movie.ID, ScreenID: screen1.ID, ShowDate: tomorrow, StartTime: "18:00", EndTime: "20:00", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true}
	database.DB.Create(&st2)

	st3 := models.Showtime{MovieID: movie.ID, ScreenID: screen2.ID, ShowDate: today, StartTime: "14:00", EndTime: "16:00", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true}
	database.DB.Create(&st3)

	return showtimeTestData{
		location: loc, theater1: theater1, theater2: theater2,
		screen1: screen1, screen2: screen2, movie: movie,
		showtime1: st1, showtime2: st2, showtime3: st3,
	}
}

// ---------- GetMovieShowtimes tests ----------

func TestGetMovieShowtimes_Success(t *testing.T) {
	testutil.SetupTestDB(t)
	td := seedShowtimeTestData(t)
	r := setupShowtimeRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/movies/1/showtimes?zipcode=33101", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	if uint(resp["movie_id"].(float64)) != td.movie.ID {
		t.Errorf("expected movie_id %d, got %v", td.movie.ID, resp["movie_id"])
	}
	if resp["title"] != "Test Movie" {
		t.Errorf("expected title 'Test Movie', got %v", resp["title"])
	}

	dates := resp["dates"].([]interface{})
	if len(dates) != 2 {
		t.Errorf("expected 2 dates, got %d", len(dates))
	}

	theaters := resp["theaters"].([]interface{})
	if len(theaters) != 2 {
		t.Errorf("expected 2 theaters, got %d", len(theaters))
	}
}

func TestGetMovieShowtimes_MultipleShowtimesPerTheater(t *testing.T) {
	testutil.SetupTestDB(t)
	seedShowtimeTestData(t)
	r := setupShowtimeRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/movies/1/showtimes?zipcode=33101", nil)
	r.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	theaters := resp["theaters"].([]interface{})
	for _, th := range theaters {
		tMap := th.(map[string]interface{})
		if tMap["name"] == "AMC Miami" {
			shows := tMap["showtimes"].([]interface{})
			if len(shows) != 2 {
				t.Errorf("expected 2 showtimes for AMC Miami, got %d", len(shows))
			}
		}
		if tMap["name"] == "Regal Miami" {
			shows := tMap["showtimes"].([]interface{})
			if len(shows) != 1 {
				t.Errorf("expected 1 showtime for Regal Miami, got %d", len(shows))
			}
		}
	}
}

func TestGetMovieShowtimes_MovieNotFound(t *testing.T) {
	testutil.SetupTestDB(t)
	r := setupShowtimeRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/movies/999/showtimes?zipcode=33101", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestGetMovieShowtimes_MissingZipcode(t *testing.T) {
	testutil.SetupTestDB(t)
	movie := models.Movie{Title: "X", Language: "English", DurationMin: 90, IsActive: true}
	database.DB.Create(&movie)
	r := setupShowtimeRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/movies/1/showtimes", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetMovieShowtimes_UnknownZipcode(t *testing.T) {
	testutil.SetupTestDB(t)
	movie := models.Movie{Title: "X", Language: "English", DurationMin: 90, IsActive: true}
	database.DB.Create(&movie)
	r := setupShowtimeRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/movies/1/showtimes?zipcode=99999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	dates := resp["dates"].([]interface{})
	theaters := resp["theaters"].([]interface{})
	if len(dates) != 0 {
		t.Errorf("expected 0 dates, got %d", len(dates))
	}
	if len(theaters) != 0 {
		t.Errorf("expected 0 theaters, got %d", len(theaters))
	}
}

func TestGetMovieShowtimes_NoShowtimesInCity(t *testing.T) {
	testutil.SetupTestDB(t)

	loc := models.Location{Zipcode: "32801", City: "Orlando", State: "FL"}
	database.DB.Create(&loc)
	movie := models.Movie{Title: "Lonely Movie", Language: "English", DurationMin: 100, IsActive: true}
	database.DB.Create(&movie)

	r := setupShowtimeRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/movies/1/showtimes?zipcode=32801", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	theaters := resp["theaters"].([]interface{})
	if len(theaters) != 0 {
		t.Errorf("expected 0 theaters, got %d", len(theaters))
	}
}

func TestGetMovieShowtimes_ShowtimeFieldValues(t *testing.T) {
	testutil.SetupTestDB(t)
	seedShowtimeTestData(t)
	r := setupShowtimeRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/movies/1/showtimes?zipcode=33101", nil)
	r.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	theaters := resp["theaters"].([]interface{})
	first := theaters[0].(map[string]interface{})
	shows := first["showtimes"].([]interface{})
	s := shows[0].(map[string]interface{})

	if s["format"] != "IMAX" {
		t.Errorf("expected format IMAX, got %v", s["format"])
	}
	if s["screen_type"] != "IMAX" {
		t.Errorf("expected screen_type IMAX, got %v", s["screen_type"])
	}
	if s["language"] != "English" {
		t.Errorf("expected language English, got %v", s["language"])
	}
	if s["start_time"] != "10:00" {
		t.Errorf("expected start_time 10:00, got %v", s["start_time"])
	}
}

// ---------- GetShowtimeSeats tests ----------

func seedSeatsForShowtime(t *testing.T) (models.Showtime, []models.Seat) {
	t.Helper()
	loc := models.Location{Zipcode: "33101", City: "Miami", State: "FL"}
	database.DB.Create(&loc)

	theater := models.Theater{Name: "Seat Theater", LocationID: loc.ID, TotalScreens: 1}
	database.DB.Create(&theater)

	screen := models.Screen{TheaterID: theater.ID, Name: "Screen 1", TotalRows: 2, TotalCols: 3, ScreenType: "Standard"}
	database.DB.Create(&screen)

	movie := models.Movie{Title: "Seat Movie", Language: "English", DurationMin: 100, IsActive: true}
	database.DB.Create(&movie)

	st := models.Showtime{
		MovieID: movie.ID, ScreenID: screen.ID,
		ShowDate: "2026-03-25", StartTime: "14:00", EndTime: "15:40",
		Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true,
	}
	database.DB.Create(&st)

	seats := []models.Seat{
		{ScreenID: screen.ID, RowLabel: "A", ColNumber: 1, SeatType: "Regular", BasePrice: 150},
		{ScreenID: screen.ID, RowLabel: "A", ColNumber: 2, SeatType: "Regular", BasePrice: 150},
		{ScreenID: screen.ID, RowLabel: "A", ColNumber: 3, SeatType: "Premium", BasePrice: 250},
		{ScreenID: screen.ID, RowLabel: "B", ColNumber: 1, SeatType: "Regular", BasePrice: 150},
		{ScreenID: screen.ID, RowLabel: "B", ColNumber: 2, SeatType: "Regular", BasePrice: 150},
		{ScreenID: screen.ID, RowLabel: "B", ColNumber: 3, SeatType: "Premium", BasePrice: 250},
	}
	database.DB.Create(&seats)

	return st, seats
}

func TestGetShowtimeSeats_Success(t *testing.T) {
	testutil.SetupTestDB(t)
	st, seats := seedSeatsForShowtime(t)
	r := setupShowtimeRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/seats?showtime_id=1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	seatList := resp["seats"].([]interface{})
	if len(seatList) != len(seats) {
		t.Errorf("expected %d seats, got %d", len(seats), len(seatList))
	}

	summary := resp["summary"].(map[string]interface{})
	if int(summary["total"].(float64)) != len(seats) {
		t.Errorf("expected total %d, got %v", len(seats), summary["total"])
	}
	if int(summary["available"].(float64)) != len(seats) {
		t.Errorf("expected all seats available, got %v available", summary["available"])
	}

	_ = st // used for seed
}

func TestGetShowtimeSeats_MissingShowtimeID(t *testing.T) {
	testutil.SetupTestDB(t)
	r := setupShowtimeRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/seats", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestGetShowtimeSeats_ShowtimeNotFound(t *testing.T) {
	testutil.SetupTestDB(t)
	r := setupShowtimeRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/seats?showtime_id=999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestGetShowtimeSeats_FilterBySeatType(t *testing.T) {
	testutil.SetupTestDB(t)
	seedSeatsForShowtime(t)
	r := setupShowtimeRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/seats?showtime_id=1&seat_type=Premium", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	seatList := resp["seats"].([]interface{})
	if len(seatList) != 2 {
		t.Errorf("expected 2 Premium seats, got %d", len(seatList))
	}
}

func TestGetShowtimeSeats_BookedSeatStatus(t *testing.T) {
	testutil.SetupTestDB(t)
	st, seats := seedSeatsForShowtime(t)

	user := models.User{Email: "test@test.com", Username: "tester", Password: "pass123", FullName: "Tester"}
	database.DB.Create(&user)

	booking := models.Booking{
		UserID: user.ID, ShowtimeID: st.ID, BookingRef: "TEST-001",
		Status: "CONFIRMED", TotalAmount: 150, PaymentStatus: "PAID",
		BookedAt: time.Now(),
	}
	database.DB.Create(&booking)

	bs := models.BookingSeat{BookingID: booking.ID, SeatID: seats[0].ID, ShowtimeID: st.ID, SeatPrice: 150}
	database.DB.Create(&bs)

	r := setupShowtimeRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/seats?showtime_id=1", nil)
	r.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	summary := resp["summary"].(map[string]interface{})
	if int(summary["booked"].(float64)) != 1 {
		t.Errorf("expected 1 booked seat, got %v", summary["booked"])
	}
	if int(summary["available"].(float64)) != 5 {
		t.Errorf("expected 5 available seats, got %v", summary["available"])
	}
}

func TestGetShowtimeSeats_FilterByStatus(t *testing.T) {
	testutil.SetupTestDB(t)
	st, seats := seedSeatsForShowtime(t)

	user := models.User{Email: "test@test.com", Username: "tester", Password: "pass123", FullName: "Tester"}
	database.DB.Create(&user)

	booking := models.Booking{
		UserID: user.ID, ShowtimeID: st.ID, BookingRef: "TEST-002",
		Status: "CONFIRMED", TotalAmount: 150, PaymentStatus: "PAID",
		BookedAt: time.Now(),
	}
	database.DB.Create(&booking)
	database.DB.Create(&models.BookingSeat{BookingID: booking.ID, SeatID: seats[0].ID, ShowtimeID: st.ID, SeatPrice: 150})

	r := setupShowtimeRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/seats?showtime_id=1&status=AVAILABLE", nil)
	r.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	seatList := resp["seats"].([]interface{})
	if len(seatList) != 5 {
		t.Errorf("expected 5 available seats in filtered list, got %d", len(seatList))
	}
}

func TestGetShowtimeSeats_PriceCalculation(t *testing.T) {
	testutil.SetupTestDB(t)

	loc := models.Location{Zipcode: "33101", City: "Miami", State: "FL"}
	database.DB.Create(&loc)
	theater := models.Theater{Name: "T", LocationID: loc.ID, TotalScreens: 1}
	database.DB.Create(&theater)
	screen := models.Screen{TheaterID: theater.ID, Name: "S1", TotalRows: 1, TotalCols: 1, ScreenType: "IMAX"}
	database.DB.Create(&screen)
	movie := models.Movie{Title: "M", Language: "English", DurationMin: 90, IsActive: true}
	database.DB.Create(&movie)

	st := models.Showtime{
		MovieID: movie.ID, ScreenID: screen.ID,
		ShowDate: "2026-03-25", StartTime: "10:00", EndTime: "11:30",
		Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true,
	}
	database.DB.Create(&st)

	seat := models.Seat{ScreenID: screen.ID, RowLabel: "A", ColNumber: 1, SeatType: "Regular", BasePrice: 200}
	database.DB.Create(&seat)

	r := setupShowtimeRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/seats?showtime_id=1", nil)
	r.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	seatList := resp["seats"].([]interface{})
	s := seatList[0].(map[string]interface{})

	expectedPrice := 300.0 // 200 * 1.5
	if s["price"].(float64) != expectedPrice {
		t.Errorf("expected price %.2f, got %.2f", expectedPrice, s["price"].(float64))
	}
}

func TestGetShowtimeSeats_ShowtimeMetadata(t *testing.T) {
	testutil.SetupTestDB(t)
	seedSeatsForShowtime(t)
	r := setupShowtimeRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/seats?showtime_id=1", nil)
	r.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	stInfo := resp["showtime"].(map[string]interface{})
	if stInfo["movie_title"] != "Seat Movie" {
		t.Errorf("expected movie_title 'Seat Movie', got %v", stInfo["movie_title"])
	}
	if stInfo["theater_name"] != "Seat Theater" {
		t.Errorf("expected theater_name 'Seat Theater', got %v", stInfo["theater_name"])
	}
	if stInfo["screen_type"] != "Standard" {
		t.Errorf("expected screen_type 'Standard', got %v", stInfo["screen_type"])
	}
	if stInfo["format"] != "2D" {
		t.Errorf("expected format '2D', got %v", stInfo["format"])
	}

	layout := resp["layout"].(map[string]interface{})
	if int(layout["total_rows"].(float64)) != 2 {
		t.Errorf("expected 2 rows, got %v", layout["total_rows"])
	}
	if int(layout["total_cols"].(float64)) != 3 {
		t.Errorf("expected 3 cols, got %v", layout["total_cols"])
	}
}
