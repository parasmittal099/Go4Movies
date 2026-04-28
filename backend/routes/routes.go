package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/parasmittal099/backend-project/config"
	"github.com/parasmittal099/backend-project/handlers"
	"github.com/parasmittal099/backend-project/middleware"
)

func RegisterRoutes(r *gin.Engine, cfg *config.Config) {
	r.Use(middleware.CORS())

	v1 := r.Group("/api/v1")
	{
		v1.GET("/locations", handlers.GetLocations)
		v1.GET("/movies", handlers.ListMovies)
		v1.GET("/movies/:id", handlers.GetMovie)
		v1.GET("/movies/:id/showtimes", handlers.GetMovieShowtimes)
		v1.GET("/seats", handlers.GetShowtimeSeats)

		v1.POST("/checkout/preview", handlers.PreviewCheckout)
		v1.POST("/checkout/confirm", middleware.OptionalJWTAuth(cfg.JWTSecret), handlers.ConfirmCheckout)
		v1.GET("/bookings", middleware.OptionalJWTAuth(cfg.JWTSecret), handlers.GetUserBookings)

		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}
	}
}
