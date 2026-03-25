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

func setupLocationRouter() *gin.Engine {
	r := gin.New()
	r.GET("/locations", GetLocations)
	return r
}

func TestGetLocations_ReturnsAll(t *testing.T) {
	testutil.SetupTestDB(t)
	database.DB.Create(&[]models.Location{
		{Zipcode: "33101", City: "Miami", State: "FL"},
		{Zipcode: "10001", City: "New York", State: "NY"},
	})

	r := setupLocationRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/locations", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string][]models.Location
	json.Unmarshal(w.Body.Bytes(), &resp)

	locs := resp["locations"]
	if len(locs) != 2 {
		t.Errorf("expected 2 locations, got %d", len(locs))
	}
}

func TestGetLocations_Empty(t *testing.T) {
	testutil.SetupTestDB(t)

	r := setupLocationRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/locations", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var resp map[string][]models.Location
	json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp["locations"]) != 0 {
		t.Errorf("expected empty list, got %d", len(resp["locations"]))
	}
}

func TestGetLocations_OrderedByCityZip(t *testing.T) {
	testutil.SetupTestDB(t)
	database.DB.Create(&[]models.Location{
		{Zipcode: "90210", City: "Zebra City", State: "CA"},
		{Zipcode: "10001", City: "Alpha City", State: "NY"},
	})

	r := setupLocationRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/locations", nil)
	r.ServeHTTP(w, req)

	var resp map[string][]models.Location
	json.Unmarshal(w.Body.Bytes(), &resp)

	locs := resp["locations"]
	if len(locs) >= 2 && locs[0].City != "Alpha City" {
		t.Errorf("expected Alpha City first, got %q", locs[0].City)
	}
}
