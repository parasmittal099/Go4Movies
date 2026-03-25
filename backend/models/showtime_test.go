package models

import "testing"

func TestShowtime_Create(t *testing.T) {
	db := setupModelsDB(t)
	loc := Location{Zipcode: "33201", City: "Miami", State: "FL"}
	db.Create(&loc)
	theater := Theater{Name: "ST Theater", LocationID: loc.ID, TotalScreens: 1}
	db.Create(&theater)
	screen := Screen{TheaterID: theater.ID, Name: "S1", TotalRows: 5, TotalCols: 11, ScreenType: "Standard"}
	db.Create(&screen)
	movie := Movie{Title: "ST Movie", Language: "English", DurationMin: 120, IsActive: true}
	db.Create(&movie)

	st := Showtime{
		MovieID: movie.ID, ScreenID: screen.ID,
		ShowDate: "2026-02-17", StartTime: "10:00", EndTime: "12:00",
		Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true,
	}
	if err := db.Create(&st).Error; err != nil {
		t.Fatalf("failed to create showtime: %v", err)
	}
	if st.ID == 0 {
		t.Error("expected non-zero ID")
	}
}

func TestShowtime_BelongsToMovie(t *testing.T) {
	db := setupModelsDB(t)
	loc := Location{Zipcode: "33202", City: "Miami", State: "FL"}
	db.Create(&loc)
	theater := Theater{Name: "T", LocationID: loc.ID, TotalScreens: 1}
	db.Create(&theater)
	screen := Screen{TheaterID: theater.ID, Name: "S", TotalRows: 5, TotalCols: 11, ScreenType: "Standard"}
	db.Create(&screen)
	movie := Movie{Title: "Linked Movie", Language: "English", DurationMin: 90, IsActive: true}
	db.Create(&movie)
	st := Showtime{MovieID: movie.ID, ScreenID: screen.ID, ShowDate: "2026-02-17", StartTime: "14:00", EndTime: "15:30", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true}
	db.Create(&st)

	var fetched Showtime
	db.Preload("Movie").First(&fetched, st.ID)
	if fetched.Movie.Title != "Linked Movie" {
		t.Errorf("expected 'Linked Movie', got %q", fetched.Movie.Title)
	}
}

func TestShowtime_BelongsToScreen(t *testing.T) {
	db := setupModelsDB(t)
	loc := Location{Zipcode: "33203", City: "Miami", State: "FL"}
	db.Create(&loc)
	theater := Theater{Name: "T", LocationID: loc.ID, TotalScreens: 1}
	db.Create(&theater)
	screen := Screen{TheaterID: theater.ID, Name: "IMAX Room", TotalRows: 5, TotalCols: 11, ScreenType: "IMAX"}
	db.Create(&screen)
	movie := Movie{Title: "M", Language: "English", DurationMin: 90, IsActive: true}
	db.Create(&movie)
	st := Showtime{MovieID: movie.ID, ScreenID: screen.ID, ShowDate: "2026-02-17", StartTime: "18:00", EndTime: "19:30", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true}
	db.Create(&st)

	var fetched Showtime
	db.Preload("Screen").First(&fetched, st.ID)
	if fetched.Screen.Name != "IMAX Room" {
		t.Errorf("expected 'IMAX Room', got %q", fetched.Screen.Name)
	}
}

func TestShowtime_DefaultPriceMultiplier(t *testing.T) {
	db := setupModelsDB(t)
	loc := Location{Zipcode: "33204", City: "Miami", State: "FL"}
	db.Create(&loc)
	theater := Theater{Name: "T", LocationID: loc.ID, TotalScreens: 1}
	db.Create(&theater)
	screen := Screen{TheaterID: theater.ID, Name: "S", TotalRows: 5, TotalCols: 11, ScreenType: "Standard"}
	db.Create(&screen)
	movie := Movie{Title: "M", Language: "English", DurationMin: 90, IsActive: true}
	db.Create(&movie)

	st := Showtime{MovieID: movie.ID, ScreenID: screen.ID, ShowDate: "2026-02-17", StartTime: "10:00", EndTime: "11:30", Language: "English", Format: "2D", IsActive: true}
	db.Create(&st)

	var fetched Showtime
	db.First(&fetched, st.ID)
	if fetched.PriceMultiplier != 1.0 {
		t.Errorf("expected default PriceMultiplier 1.0, got %v", fetched.PriceMultiplier)
	}
}
