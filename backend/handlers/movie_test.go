package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/parasmittal099/backend-project/database"
	"github.com/parasmittal099/backend-project/models"
	"github.com/parasmittal099/backend-project/testutil"
)

func strPtr(s string) *string { return &s }

func setupMovieRouter() *gin.Engine {
	r := gin.New()
	r.GET("/movies", ListMovies)
	r.GET("/movies/:id", GetMovie)
	return r
}

func seedMovieTestData(t *testing.T) {
	t.Helper()
	loc := models.Location{Zipcode: "33101", City: "Miami", State: "FL"}
	database.DB.Create(&loc)

	theater := models.Theater{Name: "Test Theater", LocationID: loc.ID, TotalScreens: 1}
	database.DB.Create(&theater)

	screen := models.Screen{TheaterID: theater.ID, Name: "Screen 1", TotalRows: 5, TotalCols: 11, ScreenType: "Standard"}
	database.DB.Create(&screen)

	movies := []models.Movie{
		{Title: "Active Movie", Language: "English", DurationMin: 120, IsActive: true, Genre: strPtr("Action")},
		{Title: "Inactive Movie", Language: "English", DurationMin: 90, IsActive: true},
	}
	database.DB.Create(&movies)
	database.DB.Model(&movies[1]).Update("is_active", false)

	database.DB.Create(&models.Showtime{
		MovieID: movies[0].ID, ScreenID: screen.ID,
		ShowDate: "2026-02-17", StartTime: "10:00", EndTime: "12:00",
		Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true,
	})
}

func TestListMovies_NoFilter(t *testing.T) {
	testutil.SetupTestDB(t)
	seedMovieTestData(t)
	r := setupMovieRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/movies", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var movies []models.Movie
	json.Unmarshal(w.Body.Bytes(), &movies)
	if len(movies) != 1 {
		t.Errorf("expected 1 active movie, got %d", len(movies))
	}
}

func TestListMovies_WithZipcode(t *testing.T) {
	testutil.SetupTestDB(t)
	seedMovieTestData(t)
	r := setupMovieRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/movies?zipcode=33101", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var movies []models.Movie
	json.Unmarshal(w.Body.Bytes(), &movies)
	if len(movies) != 1 {
		t.Errorf("expected 1 movie for zipcode 33101, got %d", len(movies))
	}
}

func TestListMovies_UnknownZipcode(t *testing.T) {
	testutil.SetupTestDB(t)
	seedMovieTestData(t)
	r := setupMovieRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/movies?zipcode=99999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var movies []models.Movie
	json.Unmarshal(w.Body.Bytes(), &movies)
	if len(movies) != 0 {
		t.Errorf("expected 0 movies for unknown zip, got %d", len(movies))
	}
}

func TestListMovies_SearchByTitle(t *testing.T) {
	testutil.SetupTestDB(t)
	seedMovieTestData(t)
	r := setupMovieRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/movies?q=active", nil)
	r.ServeHTTP(w, req)

	var movies []models.Movie
	json.Unmarshal(w.Body.Bytes(), &movies)
	if len(movies) != 1 {
		t.Errorf("expected 1 movie matching 'active', got %d", len(movies))
	}
}

func TestListMovies_SearchByGenre(t *testing.T) {
	testutil.SetupTestDB(t)
	seedMovieTestData(t)
	r := setupMovieRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/movies?q=Action", nil)
	r.ServeHTTP(w, req)

	var movies []models.Movie
	json.Unmarshal(w.Body.Bytes(), &movies)
	if len(movies) != 1 {
		t.Errorf("expected 1 movie matching genre 'Action', got %d", len(movies))
	}
}

func TestListMovies_SearchNoMatch(t *testing.T) {
	testutil.SetupTestDB(t)
	seedMovieTestData(t)
	r := setupMovieRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/movies?q=zzzzz", nil)
	r.ServeHTTP(w, req)

	var movies []models.Movie
	json.Unmarshal(w.Body.Bytes(), &movies)
	if len(movies) != 0 {
		t.Errorf("expected 0 movies, got %d", len(movies))
	}
}

func TestGetMovie_Found(t *testing.T) {
	testutil.SetupTestDB(t)
	m := models.Movie{Title: "Test", Language: "English", DurationMin: 100, IsActive: true}
	database.DB.Create(&m)

	r := setupMovieRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/movies/1", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var movie models.Movie
	json.Unmarshal(w.Body.Bytes(), &movie)
	if movie.Title != "Test" {
		t.Errorf("expected title 'Test', got %q", movie.Title)
	}
}

func TestGetMovie_NotFound(t *testing.T) {
	testutil.SetupTestDB(t)
	r := setupMovieRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/movies/999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}
