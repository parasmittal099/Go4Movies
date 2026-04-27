package database

import (
	"log"
	"strings"

	"github.com/parasmittal099/backend-project/models"
	"golang.org/x/crypto/bcrypt"
)

func isBcryptHash(s string) bool {
	return strings.HasPrefix(s, "$2a$") || strings.HasPrefix(s, "$2b$") || strings.HasPrefix(s, "$2y$")
}

// MigratePlaintextPasswords hashes any legacy plain-text passwords in users table.
func MigratePlaintextPasswords() {
	var users []models.User
	if err := DB.Find(&users).Error; err != nil {
		log.Printf("Password migration skipped: failed to load users: %v", err)
		return
	}

	migrated := 0
	for _, u := range users {
		if u.Password == "" || isBcryptHash(u.Password) {
			continue
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Password migration skipped for user %d: %v", u.ID, err)
			continue
		}

		if err := DB.Model(&models.User{}).Where("id = ?", u.ID).Update("password", string(hash)).Error; err != nil {
			log.Printf("Password migration update failed for user %d: %v", u.ID, err)
			continue
		}

		migrated++
	}

	log.Printf("Password migration completed: %d user(s) updated", migrated)
}

