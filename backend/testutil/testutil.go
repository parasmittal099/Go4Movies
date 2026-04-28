package testutil

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/parasmittal099/backend-project/database"
	"github.com/parasmittal099/backend-project/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var testDBCounter uint64

// SetupTestDB opens an in-memory SQLite database with shared-cache
// mode so that concurrent goroutines can operate on the same data.
// Each call gets a uniquely-named database for test isolation.
func SetupTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	id := atomic.AddUint64(&testDBCounter, 1)
	dsn := fmt.Sprintf("file:testdb_%d?mode=memory&cache=shared", id)

	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
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
		&models.QRTicket{},
	)
	if err != nil {
		t.Fatalf("failed to migrate test database: %v", err)
	}

	sqlDB, _ := db.DB()
	sqlDB.Exec("PRAGMA foreign_keys = ON")
	sqlDB.Exec("PRAGMA busy_timeout = 5000")

	database.DB = db
	return db
}
