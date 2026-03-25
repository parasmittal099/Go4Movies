package models

import "testing"

func TestTheater_Create(t *testing.T) {
	db := setupModelsDB(t)
	loc := Location{Zipcode: "33101", City: "Miami", State: "FL"}
	db.Create(&loc)

	theater := Theater{Name: "AMC Test", LocationID: loc.ID, TotalScreens: 2}
	if err := db.Create(&theater).Error; err != nil {
		t.Fatalf("failed to create theater: %v", err)
	}
	if theater.ID == 0 {
		t.Error("expected non-zero ID")
	}
}

func TestTheater_WithAddress(t *testing.T) {
	db := setupModelsDB(t)
	loc := Location{Zipcode: "33102", City: "Miami", State: "FL"}
	db.Create(&loc)

	addr := "123 Main St"
	theater := Theater{Name: "Regal", LocationID: loc.ID, TotalScreens: 1, Address: &addr}
	db.Create(&theater)

	var fetched Theater
	db.First(&fetched, theater.ID)
	if fetched.Address == nil || *fetched.Address != "123 Main St" {
		t.Error("address not persisted")
	}
}

func TestScreen_Create(t *testing.T) {
	db := setupModelsDB(t)
	loc := Location{Zipcode: "33103", City: "Miami", State: "FL"}
	db.Create(&loc)
	theater := Theater{Name: "T", LocationID: loc.ID, TotalScreens: 1}
	db.Create(&theater)

	screen := Screen{TheaterID: theater.ID, Name: "IMAX 1", TotalRows: 5, TotalCols: 11, ScreenType: "IMAX"}
	if err := db.Create(&screen).Error; err != nil {
		t.Fatalf("failed to create screen: %v", err)
	}
	if screen.ID == 0 {
		t.Error("expected non-zero ID")
	}
}

func TestScreen_BelongsToTheater(t *testing.T) {
	db := setupModelsDB(t)
	loc := Location{Zipcode: "33104", City: "Miami", State: "FL"}
	db.Create(&loc)
	theater := Theater{Name: "T2", LocationID: loc.ID, TotalScreens: 1}
	db.Create(&theater)
	screen := Screen{TheaterID: theater.ID, Name: "S1", TotalRows: 5, TotalCols: 11, ScreenType: "Standard"}
	db.Create(&screen)

	var fetched Screen
	db.Preload("Theater").First(&fetched, screen.ID)
	if fetched.Theater.Name != "T2" {
		t.Errorf("expected theater name 'T2', got %q", fetched.Theater.Name)
	}
}

func TestSeat_Create(t *testing.T) {
	db := setupModelsDB(t)
	loc := Location{Zipcode: "33105", City: "Miami", State: "FL"}
	db.Create(&loc)
	theater := Theater{Name: "T3", LocationID: loc.ID, TotalScreens: 1}
	db.Create(&theater)
	screen := Screen{TheaterID: theater.ID, Name: "S2", TotalRows: 5, TotalCols: 11, ScreenType: "Standard"}
	db.Create(&screen)

	seat := Seat{ScreenID: screen.ID, RowLabel: "A", ColNumber: 1, SeatType: "Premium", BasePrice: 18.0}
	if err := db.Create(&seat).Error; err != nil {
		t.Fatalf("failed to create seat: %v", err)
	}
	if seat.ID == 0 {
		t.Error("expected non-zero ID")
	}
}

func TestSeat_UniqueSeatPerScreen(t *testing.T) {
	db := setupModelsDB(t)
	loc := Location{Zipcode: "33106", City: "Miami", State: "FL"}
	db.Create(&loc)
	theater := Theater{Name: "T4", LocationID: loc.ID, TotalScreens: 1}
	db.Create(&theater)
	screen := Screen{TheaterID: theater.ID, Name: "S3", TotalRows: 5, TotalCols: 11, ScreenType: "Standard"}
	db.Create(&screen)

	db.Create(&Seat{ScreenID: screen.ID, RowLabel: "A", ColNumber: 1, SeatType: "Premium", BasePrice: 18.0})
	err := db.Create(&Seat{ScreenID: screen.ID, RowLabel: "A", ColNumber: 1, SeatType: "Premium", BasePrice: 18.0}).Error
	if err == nil {
		t.Error("expected unique constraint violation on duplicate seat")
	}
}

func TestTheater_HasManyScreens(t *testing.T) {
	db := setupModelsDB(t)
	loc := Location{Zipcode: "33107", City: "Miami", State: "FL"}
	db.Create(&loc)
	theater := Theater{Name: "Multi", LocationID: loc.ID, TotalScreens: 2}
	db.Create(&theater)
	db.Create(&Screen{TheaterID: theater.ID, Name: "S1", TotalRows: 5, TotalCols: 11, ScreenType: "Standard"})
	db.Create(&Screen{TheaterID: theater.ID, Name: "S2", TotalRows: 5, TotalCols: 11, ScreenType: "IMAX"})

	var fetched Theater
	db.Preload("Screens").First(&fetched, theater.ID)
	if len(fetched.Screens) != 2 {
		t.Errorf("expected 2 screens, got %d", len(fetched.Screens))
	}
}
