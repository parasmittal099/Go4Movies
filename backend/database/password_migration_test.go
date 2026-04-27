package database

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/parasmittal099/backend-project/models"
	"golang.org/x/crypto/bcrypt"
)

func TestMigratePlaintextPasswords_ConvertsPlaintext(t *testing.T) {
	tmp := t.TempDir()
	Connect(filepath.Join(tmp, "migrate_plaintext.db"))
	Migrate()

	user := models.User{
		Email: "plain@test.com", Username: "plain", Password: "plainpass", FullName: "Plain User",
	}
	if err := DB.Create(&user).Error; err != nil {
		t.Fatalf("create user failed: %v", err)
	}

	MigratePlaintextPasswords()

	var got models.User
	if err := DB.First(&got, user.ID).Error; err != nil {
		t.Fatalf("fetch user failed: %v", err)
	}
	if got.Password == "plainpass" {
		t.Fatalf("password was not migrated")
	}
	if !isBcryptHash(got.Password) {
		t.Fatalf("migrated password is not bcrypt hash: %q", got.Password)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(got.Password), []byte("plainpass")); err != nil {
		t.Fatalf("migrated hash does not match original password: %v", err)
	}
}

func TestMigratePlaintextPasswords_LeavesBcryptUntouched(t *testing.T) {
	tmp := t.TempDir()
	Connect(filepath.Join(tmp, "migrate_hashed.db"))
	Migrate()

	hash, err := bcrypt.GenerateFromPassword([]byte("alreadysecure"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash generation failed: %v", err)
	}
	user := models.User{
		Email: "hash@test.com", Username: "hash", Password: string(hash), FullName: "Hash User",
	}
	if err := DB.Create(&user).Error; err != nil {
		t.Fatalf("create user failed: %v", err)
	}

	before := user.Password
	MigratePlaintextPasswords()

	var got models.User
	if err := DB.First(&got, user.ID).Error; err != nil {
		t.Fatalf("fetch user failed: %v", err)
	}
	if got.Password != before {
		t.Fatalf("already-hashed password should remain unchanged")
	}
}

func TestMigratePlaintextPasswords_MixedDatasetAndIdempotent(t *testing.T) {
	tmp := t.TempDir()
	Connect(filepath.Join(tmp, "migrate_mixed.db"))
	Migrate()

	hashed, err := bcrypt.GenerateFromPassword([]byte("Strong@123"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash generation failed: %v", err)
	}

	users := []models.User{
		{Email: "plain1@test.com", Username: "plain1", Password: "Plain@123", FullName: "Plain One"},
		{Email: "hashed1@test.com", Username: "hashed1", Password: string(hashed), FullName: "Hashed One"},
		{Email: "empty1@test.com", Username: "empty1", Password: "", FullName: "Empty One"},
	}
	for _, u := range users {
		if err := DB.Create(&u).Error; err != nil {
			t.Fatalf("create user failed: %v", err)
		}
	}

	// First pass should migrate plain text and keep others safe.
	MigratePlaintextPasswords()

	var plainAfter, hashedAfter, emptyAfter models.User
	DB.Where("email = ?", "plain1@test.com").First(&plainAfter)
	DB.Where("email = ?", "hashed1@test.com").First(&hashedAfter)
	DB.Where("email = ?", "empty1@test.com").First(&emptyAfter)

	if !isBcryptHash(plainAfter.Password) {
		t.Fatalf("plain user should be migrated to bcrypt")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(plainAfter.Password), []byte("Plain@123")); err != nil {
		t.Fatalf("migrated plain user hash mismatch: %v", err)
	}
	if hashedAfter.Password != string(hashed) {
		t.Fatalf("already hashed password should remain unchanged")
	}
	if emptyAfter.Password != "" {
		t.Fatalf("empty password should remain unchanged")
	}

	// Second pass should be idempotent (no changes).
	before := plainAfter.Password
	MigratePlaintextPasswords()
	DB.Where("email = ?", "plain1@test.com").First(&plainAfter)
	if plainAfter.Password != before {
		t.Fatalf("migration must be idempotent on second run")
	}
}

func TestIsBcryptHash(t *testing.T) {
	if !isBcryptHash("$2b$10$abcdefghijklmnopqrstuv1234567890abcdefghi") {
		t.Fatalf("expected bcrypt prefix detection for $2b$")
	}
	if !isBcryptHash("$2a$10$abcdefghijklmnopqrstuv1234567890abcdefghi") {
		t.Fatalf("expected bcrypt prefix detection for $2a$")
	}
	if !isBcryptHash("$2y$10$abcdefghijklmnopqrstuv1234567890abcdefghi") {
		t.Fatalf("expected bcrypt prefix detection for $2y$")
	}

	if isBcryptHash("plainpass") {
		t.Fatalf("plain text should not be considered bcrypt")
	}
	if isBcryptHash(strings.Repeat("x", 60)) {
		t.Fatalf("random string should not be considered bcrypt")
	}
}

