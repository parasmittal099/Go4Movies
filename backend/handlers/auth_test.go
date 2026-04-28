package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/parasmittal099/backend-project/database"
	"github.com/parasmittal099/backend-project/models"
	"github.com/parasmittal099/backend-project/testutil"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupAuthRouter() *gin.Engine {
	JWTSecret = "test-jwt-secret"
	r := gin.New()
	r.POST("/register", Register)
	r.POST("/login", Login)
	return r
}

func TestRegister_Success(t *testing.T) {
	testutil.SetupTestDB(t)
	r := setupAuthRouter()

	body, _ := json.Marshal(map[string]string{
		"email":     "alice@test.com",
		"username":  "alice",
		"password":  "Secret@123",
		"full_name": "Alice Smith",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["message"] != "User registered successfully" {
		t.Errorf("unexpected message: %v", resp["message"])
	}
	if resp["token"] == nil || resp["token"] == "" {
		t.Errorf("expected a JWT token in register response")
	}

	var user models.User
	if err := database.DB.Where("email = ?", "alice@test.com").First(&user).Error; err != nil {
		t.Fatalf("failed to fetch created user: %v", err)
	}
	if user.Password == "secret123" {
		t.Fatalf("password must be hashed, found plain-text")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("Secret@123")); err != nil {
		t.Fatalf("stored password hash does not match original password: %v", err)
	}
}

func TestRegister_MissingFields(t *testing.T) {
	testutil.SetupTestDB(t)
	r := setupAuthRouter()

	body, _ := json.Marshal(map[string]string{"email": "a@b.com"})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	testutil.SetupTestDB(t)
	database.DB.Create(&models.User{
		Email: "dup@test.com", Username: "dup", Password: "pass", FullName: "Dup",
	})
	r := setupAuthRouter()

	body, _ := json.Marshal(map[string]string{
		"email":     "dup@test.com",
		"username":  "dup2",
		"password":  "Secret@123",
		"full_name": "Dup Two",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestRegister_EmailNormalized(t *testing.T) {
	testutil.SetupTestDB(t)
	r := setupAuthRouter()

	body, _ := json.Marshal(map[string]string{
		"email":     "UPPER@TEST.COM",
		"username":  "upper",
		"password":  "Secret@123",
		"full_name": "Upper Case",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}

	var user models.User
	database.DB.Where("email = ?", "upper@test.com").First(&user)
	if user.Email != "upper@test.com" {
		t.Errorf("email should be lowercased, got %q", user.Email)
	}
}

func TestRegister_UsernameValidation(t *testing.T) {
	testutil.SetupTestDB(t)
	r := setupAuthRouter()

	body, _ := json.Marshal(map[string]string{
		"email":     "user1@test.com",
		"username":  "bad user!",
		"password":  "Secret@123",
		"full_name": "User One",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRegister_WeakPasswordRejected(t *testing.T) {
	testutil.SetupTestDB(t)
	r := setupAuthRouter()

	body, _ := json.Marshal(map[string]string{
		"email":     "user2@test.com",
		"username":  "user2",
		"password":  "secret123", // missing uppercase and special
		"full_name": "User Two",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestRegister_TrimsAndNormalizesInput(t *testing.T) {
	testutil.SetupTestDB(t)
	r := setupAuthRouter()

	body, _ := json.Marshal(map[string]string{
		"email":     "MixedCase@Test.com",
		"username":  "  mixed.user  ",
		"password":  "Secret@123",
		"full_name": "  Mixed User  ",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}

	var user models.User
	if err := database.DB.Where("username = ?", "mixed.user").First(&user).Error; err != nil {
		t.Fatalf("failed to fetch normalized user: %v", err)
	}
	if user.Email != "mixedcase@test.com" {
		t.Fatalf("expected normalized email, got %q", user.Email)
	}
	if user.FullName != "Mixed User" {
		t.Fatalf("expected trimmed full name, got %q", user.FullName)
	}
}

func TestLogin_Success(t *testing.T) {
	testutil.SetupTestDB(t)
	hash, err := bcrypt.GenerateFromPassword([]byte("mypass"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash generation failed: %v", err)
	}
	database.DB.Create(&models.User{
		Email: "bob@test.com", Username: "bob", Password: string(hash), FullName: "Bob",
	})
	r := setupAuthRouter()

	body, _ := json.Marshal(map[string]string{
		"email":    "bob@test.com",
		"password": "mypass",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["message"] != "Login successful" {
		t.Errorf("unexpected message: %v", resp["message"])
	}
	if resp["token"] == nil || resp["token"] == "" {
		t.Errorf("expected a JWT token in login response")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	testutil.SetupTestDB(t)
	hash, err := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash generation failed: %v", err)
	}
	database.DB.Create(&models.User{
		Email: "bob2@test.com", Username: "bob2", Password: string(hash), FullName: "Bob",
	})
	r := setupAuthRouter()

	body, _ := json.Marshal(map[string]string{
		"email":    "bob2@test.com",
		"password": "wrong",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLogin_NonexistentUser(t *testing.T) {
	testutil.SetupTestDB(t)
	r := setupAuthRouter()

	body, _ := json.Marshal(map[string]string{
		"email":    "nobody@test.com",
		"password": "pass",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestLogin_MissingFields(t *testing.T) {
	testutil.SetupTestDB(t)
	r := setupAuthRouter()

	body, _ := json.Marshal(map[string]string{})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestLogin_MalformedStoredHashReturnsUnauthorized(t *testing.T) {
	testutil.SetupTestDB(t)
	database.DB.Create(&models.User{
		Email: "malformed@test.com", Username: "malformed", Password: "not-a-bcrypt-hash", FullName: "Malformed",
	})
	r := setupAuthRouter()

	body, _ := json.Marshal(map[string]string{
		"email":    "malformed@test.com",
		"password": "AnyPass@1",
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d: %s", w.Code, w.Body.String())
	}
}
