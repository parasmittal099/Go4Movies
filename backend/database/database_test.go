package database

import (
	"os"
	"path/filepath"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestConnect_CreatesDirectory(t *testing.T) {
	tmp := t.TempDir()
	dbPath := filepath.Join(tmp, "sub", "test.db")

	Connect(dbPath)

	if DB == nil {
		t.Fatal("expected DB to be non-nil after Connect")
	}

	if _, err := os.Stat(filepath.Dir(dbPath)); os.IsNotExist(err) {
		t.Error("expected directory to be created")
	}
}

func TestConnect_SetsDB(t *testing.T) {
	tmp := t.TempDir()
	dbPath := filepath.Join(tmp, "test.db")

	Connect(dbPath)

	if DB == nil {
		t.Fatal("DB should not be nil")
	}

	sqlDB, err := DB.DB()
	if err != nil {
		t.Fatalf("failed to get sql.DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
		t.Fatalf("failed to ping DB: %v", err)
	}
}

func TestMigrate_CreatesAllTables(t *testing.T) {
	tmp := t.TempDir()
	dbPath := filepath.Join(tmp, "migrate_test.db")

	Connect(dbPath)
	Migrate()

	expectedTables := []string{
		"users", "locations", "movies", "theaters", "screens",
		"seats", "showtimes", "bookings", "booking_seats", "payments",
	}
	for _, table := range expectedTables {
		if !DB.Migrator().HasTable(table) {
			t.Errorf("expected table %q to exist after migration", table)
		}
	}
}

func TestMigrate_Idempotent(t *testing.T) {
	tmp := t.TempDir()
	dbPath := filepath.Join(tmp, "idem_test.db")

	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}

	Migrate()
	Migrate() // should not panic or error
}
