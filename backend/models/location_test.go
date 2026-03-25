package models

import "testing"

func TestLocation_Create(t *testing.T) {
	db := setupModelsDB(t)
	loc := Location{Zipcode: "33101", City: "Miami", State: "FL"}
	if err := db.Create(&loc).Error; err != nil {
		t.Fatalf("failed to create location: %v", err)
	}
	if loc.ID == 0 {
		t.Error("expected non-zero ID")
	}
}

func TestLocation_UniqueZipcode(t *testing.T) {
	db := setupModelsDB(t)
	db.Create(&Location{Zipcode: "33101", City: "Miami", State: "FL"})
	err := db.Create(&Location{Zipcode: "33101", City: "Other", State: "FL"}).Error
	if err == nil {
		t.Error("expected unique constraint violation on zipcode")
	}
}

func TestLocation_Read(t *testing.T) {
	db := setupModelsDB(t)
	db.Create(&Location{Zipcode: "10001", City: "New York", State: "NY"})

	var loc Location
	if err := db.Where("zipcode = ?", "10001").First(&loc).Error; err != nil {
		t.Fatalf("expected to find location: %v", err)
	}
	if loc.City != "New York" {
		t.Errorf("expected 'New York', got %q", loc.City)
	}
}

func TestLocation_Update(t *testing.T) {
	db := setupModelsDB(t)
	loc := Location{Zipcode: "90210", City: "Beverly Hills", State: "CA"}
	db.Create(&loc)

	db.Model(&loc).Update("city", "Updated City")

	var fetched Location
	db.First(&fetched, loc.ID)
	if fetched.City != "Updated City" {
		t.Errorf("expected 'Updated City', got %q", fetched.City)
	}
}

func TestLocation_Delete(t *testing.T) {
	db := setupModelsDB(t)
	loc := Location{Zipcode: "00000", City: "Delete Me", State: "XX"}
	db.Create(&loc)

	db.Delete(&loc)

	var count int64
	db.Model(&Location{}).Where("id = ?", loc.ID).Count(&count)
	if count != 0 {
		t.Error("expected location to be deleted")
	}
}
