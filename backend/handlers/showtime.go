package handlers

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/parasmittal099/backend-project/database"
	"github.com/parasmittal099/backend-project/models"
)

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
			"theaters": []models.TheaterGroup{},
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
	theaterMap := map[uint]*models.TheaterGroup{}
	theaterOrder := []uint{}

	for _, r := range rows {
		dateSet[r.ShowDate] = true

		tg, exists := theaterMap[r.TheaterID]
		if !exists {
			tg = &models.TheaterGroup{
				TheaterID: r.TheaterID,
				Name:      r.TheaterName,
				Address:   r.TheaterAddress,
			}
			theaterMap[r.TheaterID] = tg
			theaterOrder = append(theaterOrder, r.TheaterID)
		}

		tg.Showtimes = append(tg.Showtimes, models.ShowtimeEntry{
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

	theaters := make([]models.TheaterGroup, 0, len(theaterOrder))
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

// GET /api/v1/seats?showtime_id=&seat_type=&status=
func GetShowtimeSeats(c *gin.Context) {
	showtimeID := c.Query("showtime_id")
	if showtimeID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "showtime_id query parameter is required"})
		return
	}

	var showtime models.Showtime
	err := database.DB.
		Preload("Screen.Theater").
		Preload("Movie").
		First(&showtime, showtimeID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Showtime not found"})
		return
	}

	var seats []models.Seat
	seatQuery := database.DB.Where("screen_id = ?", showtime.ScreenID).Order("row_label, col_number")
	if seatType := c.Query("seat_type"); seatType != "" {
		seatQuery = seatQuery.Where("seat_type = ?", seatType)
	}
	if err := seatQuery.Find(&seats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch seats"})
		return
	}

	// Booked / reserved seats for this showtime
	type bsRow struct {
		SeatID        uint
		BookingStatus string
	}
	var bsRows []bsRow
	database.DB.Table("booking_seats").
		Select("booking_seats.seat_id, bookings.status AS booking_status").
		Joins("JOIN bookings ON bookings.id = booking_seats.booking_id").
		Where("booking_seats.showtime_id = ? AND bookings.status IN ?",
			showtime.ID, []string{"PENDING", "CONFIRMED"},
		).
		Scan(&bsRows)

	statusMap := make(map[uint]string, len(bsRows))
	for _, r := range bsRows {
		switch r.BookingStatus {
		case "CONFIRMED":
			statusMap[r.SeatID] = "BOOKED"
		case "PENDING":
			if statusMap[r.SeatID] != "BOOKED" {
				statusMap[r.SeatID] = "RESERVED"
			}
		}
	}

	filterStatus := c.Query("status")
	result := make([]models.SeatResponse, 0, len(seats))
	summary := models.SeatSummary{Total: len(seats)}

	for _, seat := range seats {
		status := "AVAILABLE"
		if s, ok := statusMap[seat.ID]; ok {
			status = s
		}

		switch status {
		case "AVAILABLE":
			summary.Available++
		case "RESERVED":
			summary.Reserved++
		case "BOOKED":
			summary.Booked++
		}

		if filterStatus != "" && status != filterStatus {
			continue
		}

		price := math.Round(seat.BasePrice*showtime.PriceMultiplier*100) / 100
		result = append(result, models.SeatResponse{
			ID:        seat.ID,
			RowLabel:  seat.RowLabel,
			ColNumber: seat.ColNumber,
			SeatType:  seat.SeatType,
			Price:     price,
			Status:    status,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"showtime": gin.H{
			"id":           showtime.ID,
			"movie_title":  showtime.Movie.Title,
			"screen_name":  showtime.Screen.Name,
			"screen_type":  showtime.Screen.ScreenType,
			"theater_name": showtime.Screen.Theater.Name,
			"show_date":    showtime.ShowDate,
			"start_time":   showtime.StartTime,
			"format":       showtime.Format,
			"language":     showtime.Language,
		},
		"layout": gin.H{
			"total_rows": showtime.Screen.TotalRows,
			"total_cols": showtime.Screen.TotalCols,
		},
		"seats":   result,
		"summary": summary,
	})
}
