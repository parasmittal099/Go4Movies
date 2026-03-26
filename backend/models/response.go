package models

type ShowtimeEntry struct {
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

type TheaterGroup struct {
	TheaterID uint            `json:"theater_id"`
	Name      string          `json:"name"`
	Address   *string         `json:"address,omitempty"`
	Showtimes []ShowtimeEntry `json:"showtimes"`
}

type SeatResponse struct {
	ID        uint    `json:"id"`
	RowLabel  string  `json:"row_label"`
	ColNumber int     `json:"col_number"`
	SeatType  string  `json:"seat_type"`
	Price     float64 `json:"price"`
	Status    string  `json:"status"`
}

type SeatSummary struct {
	Total     int `json:"total"`
	Available int `json:"available"`
	Reserved  int `json:"reserved"`
	Booked    int `json:"booked"`
}
