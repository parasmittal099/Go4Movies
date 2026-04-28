package database

import (
    "log"
    "os"
    "path/filepath"

    "gorm.io/driver/sqlite"
    "github.com/parasmittal099/backend-project/models"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(dbPath string) {
    // Ensure the directory exists
    dir := filepath.Dir(dbPath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        log.Fatalf("Failed to create DB directory: %v", err)
    }

    var err error
    DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Set SQLite PRAGMAs for concurrency and integrity
    sqlDB, _ := DB.DB()
    // Disabling WAL mode (using DELETE) because SQLite WAL over Docker Desktop 
    // virtualized bind mounts (Mac/Windows) causes silent data loss and corruption.
    sqlDB.Exec("PRAGMA journal_mode = DELETE")
    sqlDB.Exec("PRAGMA foreign_keys = ON")
    sqlDB.Exec("PRAGMA busy_timeout = 5000")

    log.Println("Database connected successfully")
}

func Migrate() {
    err := DB.AutoMigrate(
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
        log.Fatalf("Failed to auto-migrate: %v", err)
    }
    log.Println("Database migration completed")
}
