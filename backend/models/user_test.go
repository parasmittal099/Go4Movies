package models

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupModelsDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	db.AutoMigrate(&User{}, &Location{}, &Movie{}, &Theater{}, &Screen{}, &Seat{}, &Showtime{}, &Booking{}, &BookingSeat{}, &Payment{}, &QRTicket{})
	sqlDB, _ := db.DB()
	sqlDB.Exec("PRAGMA foreign_keys = ON")
	return db
}

func TestUser_Create(t *testing.T) {
	db := setupModelsDB(t)
	user := User{Email: "u@t.com", Username: "user1", Password: "pass", FullName: "User One"}
	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}
	if user.ID == 0 {
		t.Error("expected non-zero ID")
	}
}

func TestUser_UniqueEmail(t *testing.T) {
	db := setupModelsDB(t)
	db.Create(&User{Email: "dup@t.com", Username: "u1", Password: "p", FullName: "N"})
	err := db.Create(&User{Email: "dup@t.com", Username: "u2", Password: "p", FullName: "N"}).Error
	if err == nil {
		t.Error("expected unique constraint violation on email")
	}
}

func TestUser_UniqueUsername(t *testing.T) {
	db := setupModelsDB(t)
	db.Create(&User{Email: "a@t.com", Username: "same", Password: "p", FullName: "N"})
	err := db.Create(&User{Email: "b@t.com", Username: "same", Password: "p", FullName: "N"}).Error
	if err == nil {
		t.Error("expected unique constraint violation on username")
	}
}

func TestUser_Read(t *testing.T) {
	db := setupModelsDB(t)
	db.Create(&User{Email: "read@t.com", Username: "reader", Password: "p", FullName: "Reader"})

	var user User
	if err := db.Where("email = ?", "read@t.com").First(&user).Error; err != nil {
		t.Fatalf("expected to find user: %v", err)
	}
	if user.Username != "reader" {
		t.Errorf("expected username 'reader', got %q", user.Username)
	}
}

func TestUser_Update(t *testing.T) {
	db := setupModelsDB(t)
	user := User{Email: "upd@t.com", Username: "upd", Password: "p", FullName: "Old"}
	db.Create(&user)

	db.Model(&user).Update("full_name", "New Name")

	var fetched User
	db.First(&fetched, user.ID)
	if fetched.FullName != "New Name" {
		t.Errorf("expected 'New Name', got %q", fetched.FullName)
	}
}

func TestUser_Delete(t *testing.T) {
	db := setupModelsDB(t)
	user := User{Email: "del@t.com", Username: "del", Password: "p", FullName: "Del"}
	db.Create(&user)

	db.Delete(&user)

	var count int64
	db.Model(&User{}).Where("id = ?", user.ID).Count(&count)
	if count != 0 {
		t.Error("expected user to be deleted")
	}
}
