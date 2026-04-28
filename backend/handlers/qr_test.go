package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/parasmittal099/backend-project/database"
	"github.com/parasmittal099/backend-project/models"
	"github.com/parasmittal099/backend-project/testutil"
)

func setupQRRouter() *gin.Engine {
	r := gin.New()
	r.GET("/bookings/by-ticket", GetBookingByTicketCode)
	return r
}

func TestGetBookingByTicketCode_Success(t *testing.T) {
	testutil.SetupTestDB(t)
	_, _, bookings := seedBookingHistoryData(t)
	target := bookings[1] // G4M-new

	qr := models.QRTicket{
		BookingID:  target.ID,
		TicketCode: "ticket-abc-123",
		IsActive:   true,
	}
	if err := database.DB.Create(&qr).Error; err != nil {
		t.Fatalf("create qr ticket: %v", err)
	}

	r := setupQRRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/bookings/by-ticket?ticket_code=ticket-abc-123", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp map[string]any
	_ = json.Unmarshal(w.Body.Bytes(), &resp)
	if resp["ticket_code"] != "ticket-abc-123" {
		t.Fatalf("expected ticket_code in response, got %v", resp["ticket_code"])
	}

	booking := resp["booking"].(map[string]any)
	if booking["booking_ref"] != "G4M-new" {
		t.Fatalf("expected booking_ref G4M-new, got %v", booking["booking_ref"])
	}
}

func TestGetBookingByTicketCode_MissingCode(t *testing.T) {
	testutil.SetupTestDB(t)
	r := setupQRRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/bookings/by-ticket", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetBookingByTicketCode_NotFound(t *testing.T) {
	testutil.SetupTestDB(t)
	r := setupQRRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/bookings/by-ticket?ticket_code=missing-code", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d: %s", w.Code, w.Body.String())
	}
}
