package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parasmittal099/backend-project/database"
	"github.com/parasmittal099/backend-project/models"
)

type showtimeEntry struct {
	ID              uint    `json:"id"`
	ShowDate        string  `json:"show_date"`
	StartTime       string  `json:"start_time"`
	EndTime         string  `json:"end_time"`
	Language        string  `json:"language"`
	Format          string  `json:"format"`
	PriceMultiplier float64 `json:"price_multiplier"`
	ScreenName      string  `json:"screen_name"`
	ScreenType      string  `json:"screen_type"`
}

type theaterGroup struct {
	TheaterID  uint            `json:"theater_id"`
	Name       string          `json:"name"`
	Address    *string         `json:"address,omitempty"`
	Showtimes  []showtimeEntry `json:"showtimes"`
}

// GET /api/v1/movies/:id/showtimes?zipcode=33101
func GetMovieShowtimes(c *gin.Context) {
	movieID := c.Param("id")

	var movie models.Movie
	if err := database.DB.First(&movie, movieID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	zipcode := c.Query("zipcode")
	if zipcode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "zipcode query parameter is required"})
		return
	}

	var location models.Location
	if err := database.DB.Where("zipcode = ?", zipcode).First(&location).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{
			"movie_id": movie.ID,
			"dates":    []string{},
			"theaters": []theaterGroup{},
		})
		return
	}

	type rawRow struct {
		ShowtimeID      uint
		ShowDate        string
		StartTime       string
		EndTime         string
		Language        string
		Format          string
		PriceMultiplier float64
		ScreenName      string
		ScreenType      string
		TheaterID       uint
		TheaterName     string
		TheaterAddress  *string
	}

	var rows []rawRow
	err := database.DB.Table("showtimes").
		Select(`showtimes.id        AS showtime_id,
				showtimes.show_date  AS show_date,
				showtimes.start_time AS start_time,
				showtimes.end_time   AS end_time,
				showtimes.language   AS language,
				showtimes.format     AS format,
				showtimes.price_multiplier AS price_multiplier,
				screens.name         AS screen_name,
				screens.screen_type  AS screen_type,
				theaters.id          AS theater_id,
				theaters.name        AS theater_name,
				theaters.address     AS theater_address`).
		Joins("JOIN screens   ON screens.id   = showtimes.screen_id").
		Joins("JOIN theaters  ON theaters.id  = screens.theater_id").
		Joins("JOIN locations ON locations.id = theaters.location_id").
		Where("showtimes.movie_id = ? AND showtimes.is_active = ? AND locations.city = ?",
			movie.ID, true, location.City).
		Order("showtimes.show_date, showtimes.start_time").
		Scan(&rows).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch showtimes"})
		return
	}

	dateSet := map[string]bool{}
	theaterMap := map[uint]*theaterGroup{}
	theaterOrder := []uint{}

	for _, r := range rows {
		dateSet[r.ShowDate] = true

		tg, exists := theaterMap[r.TheaterID]
		if !exists {
			tg = &theaterGroup{
				TheaterID: r.TheaterID,
				Name:      r.TheaterName,
				Address:   r.TheaterAddress,
			}
			theaterMap[r.TheaterID] = tg
			theaterOrder = append(theaterOrder, r.TheaterID)
		}

		tg.Showtimes = append(tg.Showtimes, showtimeEntry{
			ID:              r.ShowtimeID,
			ShowDate:        r.ShowDate,
			StartTime:       r.StartTime,
			EndTime:         r.EndTime,
			Language:        r.Language,
			Format:          r.Format,
			PriceMultiplier: r.PriceMultiplier,
			ScreenName:      r.ScreenName,
			ScreenType:      r.ScreenType,
		})
	}

	dates := make([]string, 0, len(dateSet))
	for d := range dateSet {
		dates = append(dates, d)
	}

	theaters := make([]theaterGroup, 0, len(theaterOrder))
	for _, tid := range theaterOrder {
		theaters = append(theaters, *theaterMap[tid])
	}

	c.JSON(http.StatusOK, gin.H{
		"movie_id": movie.ID,
		"title":    movie.Title,
		"dates":    dates,
		"theaters": theaters,
	})
}
