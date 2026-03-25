package testutil

import (
	"testing"

	"github.com/parasmittal099/backend-project/database"
	"github.com/parasmittal099/backend-project/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB opens an in-memory SQLite database, runs all
// migrations, and assigns it to database.DB so that handlers and
// other packages work transparently. Call it at the start of
// every test (or in TestMain) to get a fresh, isolated database.
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Location{},
		&models.Movie{},
		&models.Theater{},
		&models.Screen{},
		&models.Seat{},
		&models.Showtime{},
		&models.Booking{},
		&models.BookingSeat{},
		&models.Payment{},
	)
	if err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.Exec("PRAGMA foreign_keys = ON")

	database.DB = db
	return db
}
