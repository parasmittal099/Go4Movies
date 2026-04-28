package routes

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/parasmittal099/backend-project/config"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestRegisterRoutes_AllEndpoints(t *testing.T) {
	r := gin.New()
	cfg := &config.Config{
		DBPath:    ":memory:",
		JWTSecret: "test-secret",
		Port:      "8080",
	}

	RegisterRoutes(r, cfg)

	expected := map[string]string{
		"GET":  "/api/v1/locations",
		"GET2": "/api/v1/movies",
		"GET3": "/api/v1/movies/:id",
		"GET4": "/api/v1/movies/:id/showtimes",
		"GET5": "/api/v1/seats",
	}

	routes := r.Routes()
	routeMap := make(map[string]bool)
	for _, route := range routes {
		routeMap[route.Method+":"+route.Path] = true
	}

	checks := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/locations"},
		{"GET", "/api/v1/movies"},
		{"GET", "/api/v1/movies/:id"},
		{"GET", "/api/v1/movies/:id/showtimes"},
		{"GET", "/api/v1/seats"},
		{"GET", "/api/v1/bookings"},
		{"GET", "/api/v1/bookings/by-ticket"},
		{"POST", "/api/v1/auth/register"},
		{"POST", "/api/v1/auth/login"},
		{"POST", "/api/v1/checkout/preview"},
		{"POST", "/api/v1/checkout/confirm"},
	}

	_ = expected
	for _, c := range checks {
		key := c.method + ":" + c.path
		if !routeMap[key] {
			t.Errorf("expected route %s %s to be registered", c.method, c.path)
		}
	}
}

func TestRegisterRoutes_CountRoutes(t *testing.T) {
	r := gin.New()
	cfg := &config.Config{}

	RegisterRoutes(r, cfg)

	routes := r.Routes()
	if len(routes) < 11 {
		t.Errorf("expected at least 11 routes, got %d", len(routes))
	}
}
