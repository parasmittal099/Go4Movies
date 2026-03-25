package database

import (
	"testing"

	"github.com/parasmittal099/backend-project/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupSeedTestDB(t *testing.T) {
	t.Helper()
	var err error
	DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open in-memory db: %v", err)
	}
	sqlDB, _ := DB.DB()
	sqlDB.Exec("PRAGMA foreign_keys = ON")

	err = DB.AutoMigrate(
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
		t.Fatalf("failed to migrate: %v", err)
	}
}

func TestSeed_PopulatesLocations(t *testing.T) {
	setupSeedTestDB(t)
	Seed()

	var count int64
	DB.Model(&models.Location{}).Count(&count)
	if count == 0 {
		t.Error("expected locations to be seeded")
	}
}

func TestSeed_PopulatesMovies(t *testing.T) {
	setupSeedTestDB(t)
	Seed()

	var count int64
	DB.Model(&models.Movie{}).Count(&count)
	if count != 10 {
		t.Errorf("expected 10 movies, got %d", count)
	}
}

func TestSeed_PopulatesTheaters(t *testing.T) {
	setupSeedTestDB(t)
	Seed()

	var count int64
	DB.Model(&models.Theater{}).Count(&count)
	if count == 0 {
		t.Error("expected theaters to be seeded")
	}
}

func TestSeed_PopulatesScreens(t *testing.T) {
	setupSeedTestDB(t)
	Seed()

	var count int64
	DB.Model(&models.Screen{}).Count(&count)
	if count == 0 {
		t.Error("expected screens to be seeded")
	}
}

func TestSeed_PopulatesSeats(t *testing.T) {
	setupSeedTestDB(t)
	Seed()

	var screenCount int64
	DB.Model(&models.Screen{}).Count(&screenCount)

	var seatCount int64
	DB.Model(&models.Seat{}).Count(&seatCount)

	expectedPerScreen := int64(51)
	if seatCount != screenCount*expectedPerScreen {
		t.Errorf("expected %d seats (%d screens * %d), got %d",
			screenCount*expectedPerScreen, screenCount, expectedPerScreen, seatCount)
	}
}

func TestSeed_PopulatesShowtimes(t *testing.T) {
	setupSeedTestDB(t)
	Seed()

	var count int64
	DB.Model(&models.Showtime{}).Count(&count)
	if count == 0 {
		t.Error("expected showtimes to be seeded")
	}
}

func TestSeed_CreatesBookings(t *testing.T) {
	setupSeedTestDB(t)
	Seed()

	var count int64
	DB.Model(&models.Booking{}).Count(&count)
	if count == 0 {
		t.Error("expected sample bookings to be seeded")
	}

	var bsCount int64
	DB.Model(&models.BookingSeat{}).Count(&bsCount)
	if bsCount == 0 {
		t.Error("expected booking_seats to be seeded")
	}
}

func TestSeed_CreatesUser(t *testing.T) {
	setupSeedTestDB(t)
	Seed()

	var user models.User
	err := DB.Where("email = ?", "tester@go4movies.com").First(&user).Error
	if err != nil {
		t.Error("expected seed user to exist")
	}
}

func TestSeed_Idempotent(t *testing.T) {
	setupSeedTestDB(t)
	Seed()

	var before int64
	DB.Model(&models.Movie{}).Count(&before)

	Seed()

	var after int64
	DB.Model(&models.Movie{}).Count(&after)
	if after != before {
		t.Errorf("seed should be idempotent: movies before=%d after=%d", before, after)
	}
}
