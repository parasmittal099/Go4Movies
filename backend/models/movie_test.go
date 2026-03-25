package models

import "testing"

func sPtr(s string) *string { return &s }

func TestMovie_Create(t *testing.T) {
	db := setupModelsDB(t)
	m := Movie{Title: "Test Film", Language: "English", DurationMin: 120, IsActive: true}
	if err := db.Create(&m).Error; err != nil {
		t.Fatalf("failed to create movie: %v", err)
	}
	if m.ID == 0 {
		t.Error("expected non-zero ID")
	}
}

func TestMovie_OptionalFields(t *testing.T) {
	db := setupModelsDB(t)
	m := Movie{
		Title:       "Full Movie",
		Description: sPtr("A great movie"),
		Genre:       sPtr("Action,Drama"),
		Language:    "English",
		DurationMin: 150,
		Rating:      sPtr("PG-13"),
		PosterURL:   sPtr("https://example.com/poster.jpg"),
		Cast:        sPtr("Actor A, Actor B"),
		TrailerURL:  sPtr("https://youtube.com/watch?v=abc"),
		IsActive:    true,
	}
	if err := db.Create(&m).Error; err != nil {
		t.Fatalf("failed to create movie with optional fields: %v", err)
	}

	var fetched Movie
	db.First(&fetched, m.ID)
	if fetched.Description == nil || *fetched.Description != "A great movie" {
		t.Error("description not persisted")
	}
	if fetched.Genre == nil || *fetched.Genre != "Action,Drama" {
		t.Error("genre not persisted")
	}
}

func TestMovie_Read(t *testing.T) {
	db := setupModelsDB(t)
	db.Create(&Movie{Title: "Findable", Language: "English", DurationMin: 90, IsActive: true})

	var m Movie
	if err := db.Where("title = ?", "Findable").First(&m).Error; err != nil {
		t.Fatalf("expected to find movie: %v", err)
	}
}

func TestMovie_Update(t *testing.T) {
	db := setupModelsDB(t)
	m := Movie{Title: "Old Title", Language: "English", DurationMin: 90, IsActive: true}
	db.Create(&m)

	db.Model(&m).Update("title", "New Title")

	var fetched Movie
	db.First(&fetched, m.ID)
	if fetched.Title != "New Title" {
		t.Errorf("expected 'New Title', got %q", fetched.Title)
	}
}

func TestMovie_Delete(t *testing.T) {
	db := setupModelsDB(t)
	m := Movie{Title: "Delete Me", Language: "English", DurationMin: 90, IsActive: true}
	db.Create(&m)

	db.Delete(&m)

	var count int64
	db.Model(&Movie{}).Where("id = ?", m.ID).Count(&count)
	if count != 0 {
		t.Error("expected movie to be deleted")
	}
}

func TestMovie_FilterActive(t *testing.T) {
	db := setupModelsDB(t)
	db.Create(&Movie{Title: "Active", Language: "English", DurationMin: 90, IsActive: true})
	inactive := Movie{Title: "Inactive", Language: "English", DurationMin: 90, IsActive: true}
	db.Create(&inactive)
	db.Model(&inactive).Update("is_active", false)

	var movies []Movie
	db.Where("is_active = ?", true).Find(&movies)
	if len(movies) != 1 {
		t.Errorf("expected 1 active movie, got %d", len(movies))
	}
}
