package database

import (
	"fmt"
	"log"
	"time"

	"github.com/parasmittal099/backend-project/models"
)

func strPtr(s string) *string {
	return &s
}

func timePtr(t time.Time) *time.Time {
	return &t
}

// Seed inserts sample data into the database.
// It only runs if the locations table is empty (idempotent).
func Seed() {
	var count int64
	DB.Model(&models.Location{}).Count(&count)
	if count > 0 {
		log.Println("Database already seeded, skipping")
		return
	}

	// -------------------------------------------------------
	// 1. LOCATIONS (Florida zipcodes)
	// -------------------------------------------------------
	locations := []models.Location{
		{Zipcode: "33101", City: "Miami", State: "Florida"},
		{Zipcode: "33301", City: "Fort Lauderdale", State: "Florida"},
		{Zipcode: "32801", City: "Orlando", State: "Florida"},
		{Zipcode: "33601", City: "Tampa", State: "Florida"},
		{Zipcode: "32301", City: "Tallahassee", State: "Florida"},
		{Zipcode: "33401", City: "West Palm Beach", State: "Florida"},
		{Zipcode: "32601", City: "Gainesville", State: "Florida"},
	}
	DB.Create(&locations)
	log.Printf("Seeded %d locations", len(locations))

	// -------------------------------------------------------
	// 2. THEATERS
	// -------------------------------------------------------
	theaters := []models.Theater{
		{Name: "AMC Aventura", LocationID: locations[0].ID, Address: strPtr("19501 Biscayne Blvd, Miami"), TotalScreens: 3},
		{Name: "Regal South Beach", LocationID: locations[0].ID, Address: strPtr("1100 Lincoln Rd, Miami Beach"), TotalScreens: 2},
		{Name: "Cinemark Paradise", LocationID: locations[1].ID, Address: strPtr("5401 S University Dr, Fort Lauderdale"), TotalScreens: 2},
		{Name: "AMC Disney Springs", LocationID: locations[2].ID, Address: strPtr("1500 Buena Vista Dr, Orlando"), TotalScreens: 4},
		{Name: "Regal Pointe Orlando", LocationID: locations[2].ID, Address: strPtr("9101 International Dr, Orlando"), TotalScreens: 3},
		{Name: "AMC West Shore", LocationID: locations[3].ID, Address: strPtr("210 Westshore Plaza, Tampa"), TotalScreens: 2},
		{Name: "Regal Celebration Pointe", LocationID: locations[6].ID, Address: strPtr("3528 SW 45th St, Gainesville"), TotalScreens: 2},
	}
	DB.Create(&theaters)
	log.Printf("Seeded %d theaters", len(theaters))

	// -------------------------------------------------------
	// 3. SCREENS (1 per theater — all share the same 5-row
	//    layout that the frontend renders)
	// -------------------------------------------------------
	screens := []models.Screen{
		{TheaterID: theaters[0].ID, Name: "Screen 1", TotalRows: 5, TotalCols: 11, ScreenType: "IMAX"},
		{TheaterID: theaters[1].ID, Name: "Screen 1", TotalRows: 5, TotalCols: 11, ScreenType: "Standard"},
		{TheaterID: theaters[2].ID, Name: "Screen 1", TotalRows: 5, TotalCols: 11, ScreenType: "Standard"},
		{TheaterID: theaters[3].ID, Name: "Screen 1", TotalRows: 5, TotalCols: 11, ScreenType: "IMAX"},
		{TheaterID: theaters[4].ID, Name: "Screen 1", TotalRows: 5, TotalCols: 11, ScreenType: "4DX"},
		{TheaterID: theaters[5].ID, Name: "Screen 1", TotalRows: 5, TotalCols: 11, ScreenType: "Standard"},
		{TheaterID: theaters[6].ID, Name: "Screen 1", TotalRows: 5, TotalCols: 11, ScreenType: "Standard"},
	}
	DB.Create(&screens)
	log.Printf("Seeded %d screens", len(screens))

	// -------------------------------------------------------
	// 3b. SEATS (matching the frontend layout exactly)
	//     Rows A–B: Premium  (11 seats, $18)
	//     Rows C–D: Premium  (11 seats, $18) — "Executive" in UI
	//     Row  E:   VIP      ( 7 seats, $25) — recliners
	//     Total: 51 seats per screen
	// -------------------------------------------------------
	type seatDef struct {
		Row   string
		Count int
		Type  string
		Price float64
	}
	seatLayout := []seatDef{
		{"A", 11, "Premium", 18.0},
		{"B", 11, "Premium", 18.0},
		{"C", 11, "Premium", 18.0},
		{"D", 11, "Premium", 18.0},
		{"E", 7, "VIP", 25.0},
	}
	for _, screen := range screens {
		var seats []models.Seat
		for _, def := range seatLayout {
			for col := 1; col <= def.Count; col++ {
				seats = append(seats, models.Seat{
					ScreenID:  screen.ID,
					RowLabel:  def.Row,
					ColNumber: col,
					SeatType:  def.Type,
					BasePrice: def.Price,
				})
			}
		}
		DB.Create(&seats)
	}
	log.Printf("Seeded 51 seats per screen (%d screens)", len(screens))

	// -------------------------------------------------------
	// 4. MOVIES (poster URLs sourced from TMDB via jsonfakery)
	// -------------------------------------------------------
	movies := []models.Movie{
		{
			Title:       "The Karate Kid",
			Description: strPtr("Twelve-year-old Dre Parker could have been the most popular kid in Detroit, but his mother's latest career move has landed him in China. Dre immediately falls for his classmate Mei Ying but the cultural differences make such a friendship impossible. Even worse, Dre's feelings make him an enemy of the class bully, Cheng."),
			Genre:       strPtr("Action,Drama,Family"),
			Language:    "English",
			DurationMin: 140,
			Rating:      strPtr("PG"),
			PosterURL:   strPtr("https://image.tmdb.org/t/p/original/o4gwFcm6qRgnC9Gff4giy0POHdF.jpg"),
			Cast:        strPtr("Jackie Chan, Jaden Smith, Taraji P. Henson, Wenwen Han"),
			TrailerURL:  strPtr("https://youtube.com/watch?v=2zHSrABr0xY"),
			ReleaseDate: timePtr(time.Date(2010, 6, 10, 0, 0, 0, 0, time.UTC)),
			IsActive:    true,
		},
		{
			Title:       "Ghost in the Shell",
			Description: strPtr("In the near future, Major is the first of her kind: a human saved from a terrible crash, then cyber-enhanced to be a perfect soldier devoted to stopping the world's most dangerous criminals."),
			Genre:       strPtr("Sci-Fi,Action,Drama"),
			Language:    "English",
			DurationMin: 107,
			Rating:      strPtr("PG-13"),
			PosterURL:   strPtr("https://image.tmdb.org/t/p/original/myRzRzCxdfUWjkJWgpHHZ1oGkJd.jpg"),
			Cast:        strPtr("Scarlett Johansson, Pilou Asbæk, Takeshi Kitano, Juliette Binoche"),
			TrailerURL:  strPtr("https://youtube.com/watch?v=G4VmJcZR0Yg"),
			ReleaseDate: timePtr(time.Date(2017, 3, 29, 0, 0, 0, 0, time.UTC)),
			IsActive:    true,
		},
		{
			Title:       "Tomorrow Never Dies",
			Description: strPtr("A deranged media mogul is staging international incidents to pit the world's superpowers against each other. Now James Bond must take on this evil mastermind in an adrenaline-charged battle to end his reign of terror and prevent global pandemonium."),
			Genre:       strPtr("Action,Adventure,Thriller"),
			Language:    "English",
			DurationMin: 119,
			Rating:      strPtr("PG-13"),
			PosterURL:   strPtr("https://image.tmdb.org/t/p/original/gZm002w7q9yLOkltxT76TWGfdZX.jpg"),
			Cast:        strPtr("Pierce Brosnan, Michelle Yeoh, Jonathan Pryce, Judi Dench"),
			TrailerURL:  strPtr("https://youtube.com/watch?v=6gRqFSMqbOs"),
			ReleaseDate: timePtr(time.Date(1997, 12, 11, 0, 0, 0, 0, time.UTC)),
			IsActive:    true,
		},
		{
			Title:       "Black Death",
			Description: strPtr("As the plague decimates medieval Europe, rumours circulate of a village immune from the plague. There is talk of a necromancer who leads the village and is able to raise the dead. A fearsome knight joined by a cohort of soldiers and a young monk are charged by the church to investigate."),
			Genre:       strPtr("Drama,History,Thriller"),
			Language:    "English",
			DurationMin: 102,
			Rating:      strPtr("R"),
			PosterURL:   strPtr("https://image.tmdb.org/t/p/original/gXRERDpyT9s3m2yk6wNmrTWbZfG.jpg"),
			Cast:        strPtr("Sean Bean, Eddie Redmayne, Carice van Houten, David Warner"),
			TrailerURL:  strPtr("https://youtube.com/watch?v=THG_azR8hgM"),
			ReleaseDate: timePtr(time.Date(2010, 6, 7, 0, 0, 0, 0, time.UTC)),
			IsActive:    true,
		},
		{
			Title:       "Walk Hard: The Dewey Cox Story",
			Description: strPtr("Following a childhood tragedy, Dewey Cox follows a long and winding road to music stardom. Dewey perseveres through changing musical styles, an addiction to nearly every drug known and bouts of uncontrollable rage."),
			Genre:       strPtr("Comedy,Music"),
			Language:    "English",
			DurationMin: 96,
			Rating:      strPtr("R"),
			PosterURL:   strPtr("https://image.tmdb.org/t/p/original/Aa1IQ4Cuin3d7qIahvPheMmR4E5.jpg"),
			Cast:        strPtr("John C. Reilly, Jenna Fischer, Jack Black, Paul Rudd"),
			TrailerURL:  strPtr("https://youtube.com/watch?v=dSjEfPdkbJo"),
			ReleaseDate: timePtr(time.Date(2007, 12, 21, 0, 0, 0, 0, time.UTC)),
			IsActive:    true,
		},
		{
			Title:       "The Other Side of the Door",
			Description: strPtr("Grieving over the loss of her son, a mother struggles with her feelings for her daughter and her husband. She seeks out a ritual that allows her say goodbye to her dead child, opening the veil between the world of the dead and the living."),
			Genre:       strPtr("Horror,Thriller"),
			Language:    "English",
			DurationMin: 96,
			Rating:      strPtr("PG-13"),
			PosterURL:   strPtr("https://image.tmdb.org/t/p/original/z5HubVaGLp1EXQbEYee9.jpg"),
			Cast:        strPtr("Sarah Wayne Callies, Jeremy Sisto, Sofia Rosinsky, Javier Botet"),
			TrailerURL:  strPtr("https://youtube.com/watch?v=o_Kf3bGLdVA"),
			ReleaseDate: timePtr(time.Date(2016, 2, 25, 0, 0, 0, 0, time.UTC)),
			IsActive:    true,
		},
		{
			Title:       "Nick and Norah's Infinite Playlist",
			Description: strPtr("Nick cannot stop obsessing over his ex-girlfriend, Tris, until Tris' friend Norah suddenly shows interest in him at a club. Thus begins an odd night filled with ups and downs as the two keep running into Tris and her new boyfriend while searching for Norah's drunken friend."),
			Genre:       strPtr("Comedy,Drama,Romance"),
			Language:    "English",
			DurationMin: 90,
			Rating:      strPtr("PG-13"),
			PosterURL:   strPtr("https://image.tmdb.org/t/p/original/cndxWgZEHnSUtytrMjal97ye34E.jpg"),
			Cast:        strPtr("Michael Cera, Kat Dennings, Ari Graynor, Aaron Yoo"),
			TrailerURL:  strPtr("https://youtube.com/watch?v=AhXOGqWAB2g"),
			ReleaseDate: timePtr(time.Date(2008, 10, 3, 0, 0, 0, 0, time.UTC)),
			IsActive:    true,
		},
		{
			Title:       "On the Line",
			Description: strPtr("A provocative and edgy radio host must play a dangerous game of cat and mouse with a mysterious caller who's kidnapped his family and is threatening to blow up the whole station."),
			Genre:       strPtr("Thriller"),
			Language:    "English",
			DurationMin: 103,
			Rating:      strPtr("R"),
			PosterURL:   strPtr("https://image.tmdb.org/t/p/original/jVmWI8PqoVTHCnrLYAcyrclzeY0.jpg"),
			Cast:        strPtr("Mel Gibson, Kevin Dillon, William Moseley, Enrique Arce"),
			TrailerURL:  strPtr("https://youtube.com/watch?v=9i5zcvOkfEI"),
			ReleaseDate: timePtr(time.Date(2022, 10, 31, 0, 0, 0, 0, time.UTC)),
			IsActive:    true,
		},
		{
			Title:       "The Brave Little Toaster Goes to Mars",
			Description: strPtr("Whilst trying to protect their new Little Master the anthropomorphic appliances set off on an epic adventure to Mars and make many new friends along the way."),
			Genre:       strPtr("Animation,Family,Sci-Fi"),
			Language:    "English",
			DurationMin: 74,
			Rating:      strPtr("G"),
			PosterURL:   strPtr("https://image.tmdb.org/t/p/original/zO6SF0qLZfgSVRIfk8Y8l4hmN5e.jpg"),
			Cast:        strPtr("Deanna Oliver, Wayne Knight, Jim Cummings, Farrah Fawcett"),
			TrailerURL:  strPtr("https://youtube.com/watch?v=gYMljEaXMHI"),
			ReleaseDate: timePtr(time.Date(1998, 5, 19, 0, 0, 0, 0, time.UTC)),
			IsActive:    true,
		},
		{
			Title:       "The Humbling",
			Description: strPtr("Following a breakdown and suicide attempt, an aging actor becomes involved with a much younger woman but soon finds that it's difficult to keep pace with her."),
			Genre:       strPtr("Drama,Comedy"),
			Language:    "English",
			DurationMin: 112,
			Rating:      strPtr("R"),
			PosterURL:   strPtr("https://image.tmdb.org/t/p/original/8x8NUmtcuNqlIH1JDAoxUjoDX7P.jpg"),
			Cast:        strPtr("Al Pacino, Greta Gerwig, Kyra Sedgwick, Dianne Wiest"),
			TrailerURL:  strPtr("https://youtube.com/watch?v=DaVrEFGmrlY"),
			ReleaseDate: timePtr(time.Date(2015, 1, 23, 0, 0, 0, 0, time.UTC)),
			IsActive:    true,
		},
	}
	DB.Create(&movies)
	log.Printf("Seeded %d movies", len(movies))

	// -------------------------------------------------------
	// 5. SHOWTIMES
	//    Spread across theaters so different cities have
	//    different movie selections for testing.
	//    5 days of shows with multiple time slots per day.
	// -------------------------------------------------------
	day0 := time.Now().Format("2006-01-02")
	day1 := time.Now().Add(24 * time.Hour).Format("2006-01-02")
	day2 := time.Now().Add(48 * time.Hour).Format("2006-01-02")
	day3 := time.Now().Add(72 * time.Hour).Format("2006-01-02")
	day4 := time.Now().Add(96 * time.Hour).Format("2006-01-02")

	showtimes := []models.Showtime{
		// ===== Miami -- AMC Aventura (screen 0, IMAX) =====
		// The Karate Kid (140 min)
		{MovieID: movies[0].ID, ScreenID: screens[0].ID, ShowDate: day0, StartTime: "10:00", EndTime: "12:20", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[0].ID, ShowDate: day0, StartTime: "15:00", EndTime: "17:20", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[0].ID, ShowDate: day0, StartTime: "20:00", EndTime: "22:20", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[0].ID, ShowDate: day1, StartTime: "10:00", EndTime: "12:20", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[0].ID, ShowDate: day1, StartTime: "18:00", EndTime: "20:20", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[0].ID, ShowDate: day2, StartTime: "11:00", EndTime: "13:20", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[0].ID, ShowDate: day3, StartTime: "14:00", EndTime: "16:20", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[0].ID, ShowDate: day4, StartTime: "10:00", EndTime: "12:20", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[0].ID, ShowDate: day4, StartTime: "19:00", EndTime: "21:20", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true},
		// Tomorrow Never Dies (119 min)
		{MovieID: movies[2].ID, ScreenID: screens[0].ID, ShowDate: day0, StartTime: "13:00", EndTime: "14:59", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true},
		{MovieID: movies[2].ID, ScreenID: screens[0].ID, ShowDate: day1, StartTime: "13:00", EndTime: "14:59", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true},
		{MovieID: movies[2].ID, ScreenID: screens[0].ID, ShowDate: day2, StartTime: "16:00", EndTime: "17:59", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true},
		{MovieID: movies[2].ID, ScreenID: screens[0].ID, ShowDate: day3, StartTime: "10:00", EndTime: "11:59", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true},
		// Walk Hard (96 min)
		{MovieID: movies[4].ID, ScreenID: screens[0].ID, ShowDate: day1, StartTime: "15:30", EndTime: "17:06", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true},
		{MovieID: movies[4].ID, ScreenID: screens[0].ID, ShowDate: day2, StartTime: "19:00", EndTime: "20:36", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true},
		{MovieID: movies[4].ID, ScreenID: screens[0].ID, ShowDate: day4, StartTime: "14:00", EndTime: "15:36", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true},

		// ===== Miami -- Regal South Beach (screen 1, Standard) =====
		// Ghost in the Shell (107 min)
		{MovieID: movies[1].ID, ScreenID: screens[1].ID, ShowDate: day0, StartTime: "11:00", EndTime: "12:47", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[1].ID, ScreenID: screens[1].ID, ShowDate: day0, StartTime: "16:00", EndTime: "17:47", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[1].ID, ScreenID: screens[1].ID, ShowDate: day0, StartTime: "20:30", EndTime: "22:17", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[1].ID, ScreenID: screens[1].ID, ShowDate: day1, StartTime: "11:00", EndTime: "12:47", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[1].ID, ScreenID: screens[1].ID, ShowDate: day1, StartTime: "19:00", EndTime: "20:47", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[1].ID, ScreenID: screens[1].ID, ShowDate: day2, StartTime: "14:00", EndTime: "15:47", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[1].ID, ScreenID: screens[1].ID, ShowDate: day3, StartTime: "11:00", EndTime: "12:47", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[1].ID, ScreenID: screens[1].ID, ShowDate: day4, StartTime: "17:00", EndTime: "18:47", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		// The Other Side of the Door (96 min)
		{MovieID: movies[5].ID, ScreenID: screens[1].ID, ShowDate: day0, StartTime: "13:30", EndTime: "15:06", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[5].ID, ScreenID: screens[1].ID, ShowDate: day1, StartTime: "14:00", EndTime: "15:36", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[5].ID, ScreenID: screens[1].ID, ShowDate: day2, StartTime: "18:00", EndTime: "19:36", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[5].ID, ScreenID: screens[1].ID, ShowDate: day3, StartTime: "20:00", EndTime: "21:36", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		// The Humbling (112 min)
		{MovieID: movies[9].ID, ScreenID: screens[1].ID, ShowDate: day1, StartTime: "16:00", EndTime: "17:52", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[9].ID, ScreenID: screens[1].ID, ShowDate: day2, StartTime: "10:00", EndTime: "11:52", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[9].ID, ScreenID: screens[1].ID, ShowDate: day4, StartTime: "12:00", EndTime: "13:52", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},

		// ===== Fort Lauderdale -- Cinemark Paradise (screen 2, Standard) =====
		// The Karate Kid (140 min)
		{MovieID: movies[0].ID, ScreenID: screens[2].ID, ShowDate: day0, StartTime: "10:30", EndTime: "12:50", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[2].ID, ShowDate: day0, StartTime: "16:00", EndTime: "18:20", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[2].ID, ShowDate: day1, StartTime: "10:30", EndTime: "12:50", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[2].ID, ShowDate: day1, StartTime: "19:00", EndTime: "21:20", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[2].ID, ShowDate: day2, StartTime: "13:00", EndTime: "15:20", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[2].ID, ShowDate: day3, StartTime: "10:30", EndTime: "12:50", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[2].ID, ShowDate: day4, StartTime: "16:00", EndTime: "18:20", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		// Black Death (102 min)
		{MovieID: movies[3].ID, ScreenID: screens[2].ID, ShowDate: day0, StartTime: "14:00", EndTime: "15:42", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[3].ID, ScreenID: screens[2].ID, ShowDate: day0, StartTime: "19:30", EndTime: "21:12", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[3].ID, ScreenID: screens[2].ID, ShowDate: day2, StartTime: "17:00", EndTime: "18:42", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[3].ID, ScreenID: screens[2].ID, ShowDate: day3, StartTime: "14:00", EndTime: "15:42", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		// Nick and Norah's Infinite Playlist (90 min)
		{MovieID: movies[6].ID, ScreenID: screens[2].ID, ShowDate: day1, StartTime: "14:00", EndTime: "15:30", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[6].ID, ScreenID: screens[2].ID, ShowDate: day2, StartTime: "20:00", EndTime: "21:30", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[6].ID, ScreenID: screens[2].ID, ShowDate: day4, StartTime: "11:00", EndTime: "12:30", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[6].ID, ScreenID: screens[2].ID, ShowDate: day4, StartTime: "19:00", EndTime: "20:30", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},

		// ===== Orlando -- AMC Disney Springs (screen 3, IMAX) =====
		// Tomorrow Never Dies (119 min)
		{MovieID: movies[2].ID, ScreenID: screens[3].ID, ShowDate: day0, StartTime: "11:00", EndTime: "12:59", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true},
		{MovieID: movies[2].ID, ScreenID: screens[3].ID, ShowDate: day0, StartTime: "16:00", EndTime: "17:59", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true},
		{MovieID: movies[2].ID, ScreenID: screens[3].ID, ShowDate: day0, StartTime: "20:00", EndTime: "21:59", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true},
		{MovieID: movies[2].ID, ScreenID: screens[3].ID, ShowDate: day1, StartTime: "11:00", EndTime: "12:59", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true},
		{MovieID: movies[2].ID, ScreenID: screens[3].ID, ShowDate: day2, StartTime: "14:00", EndTime: "15:59", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true},
		{MovieID: movies[2].ID, ScreenID: screens[3].ID, ShowDate: day3, StartTime: "18:00", EndTime: "19:59", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true},
		{MovieID: movies[2].ID, ScreenID: screens[3].ID, ShowDate: day4, StartTime: "11:00", EndTime: "12:59", Language: "English", Format: "IMAX", PriceMultiplier: 1.5, IsActive: true},
		// Brave Little Toaster (74 min)
		{MovieID: movies[8].ID, ScreenID: screens[3].ID, ShowDate: day0, StartTime: "14:00", EndTime: "15:14", Language: "English", Format: "IMAX", PriceMultiplier: 1.2, IsActive: true},
		{MovieID: movies[8].ID, ScreenID: screens[3].ID, ShowDate: day1, StartTime: "14:00", EndTime: "15:14", Language: "English", Format: "IMAX", PriceMultiplier: 1.2, IsActive: true},
		{MovieID: movies[8].ID, ScreenID: screens[3].ID, ShowDate: day2, StartTime: "10:00", EndTime: "11:14", Language: "English", Format: "IMAX", PriceMultiplier: 1.2, IsActive: true},
		{MovieID: movies[8].ID, ScreenID: screens[3].ID, ShowDate: day3, StartTime: "10:00", EndTime: "11:14", Language: "English", Format: "IMAX", PriceMultiplier: 1.2, IsActive: true},
		// On the Line (103 min)
		{MovieID: movies[7].ID, ScreenID: screens[3].ID, ShowDate: day1, StartTime: "10:00", EndTime: "11:43", Language: "English", Format: "IMAX", PriceMultiplier: 1.2, IsActive: true},
		{MovieID: movies[7].ID, ScreenID: screens[3].ID, ShowDate: day1, StartTime: "17:00", EndTime: "18:43", Language: "English", Format: "IMAX", PriceMultiplier: 1.2, IsActive: true},
		{MovieID: movies[7].ID, ScreenID: screens[3].ID, ShowDate: day2, StartTime: "19:00", EndTime: "20:43", Language: "English", Format: "IMAX", PriceMultiplier: 1.2, IsActive: true},
		{MovieID: movies[7].ID, ScreenID: screens[3].ID, ShowDate: day4, StartTime: "15:00", EndTime: "16:43", Language: "English", Format: "IMAX", PriceMultiplier: 1.2, IsActive: true},

		// ===== Orlando -- Regal Pointe Orlando (screen 4, 4DX) =====
		// The Karate Kid (140 min)
		{MovieID: movies[0].ID, ScreenID: screens[4].ID, ShowDate: day0, StartTime: "12:00", EndTime: "14:20", Language: "English", Format: "4DX", PriceMultiplier: 1.8, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[4].ID, ShowDate: day0, StartTime: "18:00", EndTime: "20:20", Language: "English", Format: "4DX", PriceMultiplier: 1.8, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[4].ID, ShowDate: day1, StartTime: "10:00", EndTime: "12:20", Language: "English", Format: "4DX", PriceMultiplier: 1.8, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[4].ID, ShowDate: day1, StartTime: "15:00", EndTime: "17:20", Language: "English", Format: "4DX", PriceMultiplier: 1.8, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[4].ID, ShowDate: day2, StartTime: "12:00", EndTime: "14:20", Language: "English", Format: "4DX", PriceMultiplier: 1.8, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[4].ID, ShowDate: day3, StartTime: "18:00", EndTime: "20:20", Language: "English", Format: "4DX", PriceMultiplier: 1.8, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[4].ID, ShowDate: day4, StartTime: "10:00", EndTime: "12:20", Language: "English", Format: "4DX", PriceMultiplier: 1.8, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[4].ID, ShowDate: day4, StartTime: "19:00", EndTime: "21:20", Language: "English", Format: "4DX", PriceMultiplier: 1.8, IsActive: true},
		// Ghost in the Shell (107 min)
		{MovieID: movies[1].ID, ScreenID: screens[4].ID, ShowDate: day0, StartTime: "15:00", EndTime: "16:47", Language: "English", Format: "4DX", PriceMultiplier: 1.8, IsActive: true},
		{MovieID: movies[1].ID, ScreenID: screens[4].ID, ShowDate: day1, StartTime: "13:00", EndTime: "14:47", Language: "English", Format: "4DX", PriceMultiplier: 1.8, IsActive: true},
		{MovieID: movies[1].ID, ScreenID: screens[4].ID, ShowDate: day2, StartTime: "17:00", EndTime: "18:47", Language: "English", Format: "4DX", PriceMultiplier: 1.8, IsActive: true},
		{MovieID: movies[1].ID, ScreenID: screens[4].ID, ShowDate: day3, StartTime: "11:00", EndTime: "12:47", Language: "English", Format: "4DX", PriceMultiplier: 1.8, IsActive: true},
		// The Other Side of the Door (96 min)
		{MovieID: movies[5].ID, ScreenID: screens[4].ID, ShowDate: day1, StartTime: "19:00", EndTime: "20:36", Language: "English", Format: "4DX", PriceMultiplier: 1.8, IsActive: true},
		{MovieID: movies[5].ID, ScreenID: screens[4].ID, ShowDate: day2, StartTime: "10:00", EndTime: "11:36", Language: "English", Format: "4DX", PriceMultiplier: 1.8, IsActive: true},
		{MovieID: movies[5].ID, ScreenID: screens[4].ID, ShowDate: day3, StartTime: "15:00", EndTime: "16:36", Language: "English", Format: "4DX", PriceMultiplier: 1.8, IsActive: true},
		{MovieID: movies[5].ID, ScreenID: screens[4].ID, ShowDate: day4, StartTime: "14:00", EndTime: "15:36", Language: "English", Format: "4DX", PriceMultiplier: 1.8, IsActive: true},

		// ===== Tampa -- AMC West Shore (screen 5, Standard) =====
		// Walk Hard (96 min)
		{MovieID: movies[4].ID, ScreenID: screens[5].ID, ShowDate: day0, StartTime: "13:00", EndTime: "14:36", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[4].ID, ScreenID: screens[5].ID, ShowDate: day0, StartTime: "19:00", EndTime: "20:36", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[4].ID, ScreenID: screens[5].ID, ShowDate: day1, StartTime: "13:00", EndTime: "14:36", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[4].ID, ScreenID: screens[5].ID, ShowDate: day2, StartTime: "16:00", EndTime: "17:36", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[4].ID, ScreenID: screens[5].ID, ShowDate: day3, StartTime: "13:00", EndTime: "14:36", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		// Nick and Norah's Infinite Playlist (90 min)
		{MovieID: movies[6].ID, ScreenID: screens[5].ID, ShowDate: day0, StartTime: "17:00", EndTime: "18:30", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[6].ID, ScreenID: screens[5].ID, ShowDate: day1, StartTime: "10:00", EndTime: "11:30", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[6].ID, ScreenID: screens[5].ID, ShowDate: day1, StartTime: "17:00", EndTime: "18:30", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[6].ID, ScreenID: screens[5].ID, ShowDate: day2, StartTime: "20:00", EndTime: "21:30", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[6].ID, ScreenID: screens[5].ID, ShowDate: day4, StartTime: "11:00", EndTime: "12:30", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		// On the Line (103 min)
		{MovieID: movies[7].ID, ScreenID: screens[5].ID, ShowDate: day1, StartTime: "15:00", EndTime: "16:43", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[7].ID, ScreenID: screens[5].ID, ShowDate: day2, StartTime: "10:00", EndTime: "11:43", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[7].ID, ScreenID: screens[5].ID, ShowDate: day2, StartTime: "13:00", EndTime: "14:43", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[7].ID, ScreenID: screens[5].ID, ShowDate: day3, StartTime: "19:00", EndTime: "20:43", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		// The Humbling (112 min)
		{MovieID: movies[9].ID, ScreenID: screens[5].ID, ShowDate: day0, StartTime: "10:00", EndTime: "11:52", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[9].ID, ScreenID: screens[5].ID, ShowDate: day1, StartTime: "19:30", EndTime: "21:22", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[9].ID, ScreenID: screens[5].ID, ShowDate: day3, StartTime: "10:00", EndTime: "11:52", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[9].ID, ScreenID: screens[5].ID, ShowDate: day4, StartTime: "15:00", EndTime: "16:52", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[9].ID, ScreenID: screens[5].ID, ShowDate: day4, StartTime: "19:00", EndTime: "20:52", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},

		// ===== Gainesville -- Regal Celebration Pointe (screen 6, Standard) =====
		// Black Death (102 min)
		{MovieID: movies[3].ID, ScreenID: screens[6].ID, ShowDate: day0, StartTime: "11:00", EndTime: "12:42", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[3].ID, ScreenID: screens[6].ID, ShowDate: day0, StartTime: "18:00", EndTime: "19:42", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[3].ID, ScreenID: screens[6].ID, ShowDate: day1, StartTime: "11:00", EndTime: "12:42", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[3].ID, ScreenID: screens[6].ID, ShowDate: day2, StartTime: "15:00", EndTime: "16:42", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[3].ID, ScreenID: screens[6].ID, ShowDate: day3, StartTime: "18:00", EndTime: "19:42", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		// Brave Little Toaster (74 min)
		{MovieID: movies[8].ID, ScreenID: screens[6].ID, ShowDate: day0, StartTime: "14:00", EndTime: "15:14", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[8].ID, ScreenID: screens[6].ID, ShowDate: day1, StartTime: "14:00", EndTime: "15:14", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[8].ID, ScreenID: screens[6].ID, ShowDate: day2, StartTime: "10:00", EndTime: "11:14", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[8].ID, ScreenID: screens[6].ID, ShowDate: day3, StartTime: "14:00", EndTime: "15:14", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[8].ID, ScreenID: screens[6].ID, ShowDate: day4, StartTime: "10:00", EndTime: "11:14", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[8].ID, ScreenID: screens[6].ID, ShowDate: day4, StartTime: "16:00", EndTime: "17:14", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		// The Karate Kid (140 min)
		{MovieID: movies[0].ID, ScreenID: screens[6].ID, ShowDate: day1, StartTime: "10:00", EndTime: "12:20", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[6].ID, ShowDate: day1, StartTime: "16:00", EndTime: "18:20", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[6].ID, ShowDate: day2, StartTime: "12:00", EndTime: "14:20", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[6].ID, ShowDate: day3, StartTime: "10:00", EndTime: "12:20", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[0].ID, ScreenID: screens[6].ID, ShowDate: day4, StartTime: "13:00", EndTime: "15:20", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		// Ghost in the Shell (107 min)
		{MovieID: movies[1].ID, ScreenID: screens[6].ID, ShowDate: day2, StartTime: "18:00", EndTime: "19:47", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[1].ID, ScreenID: screens[6].ID, ShowDate: day3, StartTime: "20:00", EndTime: "21:47", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
		{MovieID: movies[1].ID, ScreenID: screens[6].ID, ShowDate: day4, StartTime: "19:00", EndTime: "20:47", Language: "English", Format: "2D", PriceMultiplier: 1.0, IsActive: true},
	}
	user := models.User{
		Username: "Tester 1",
		FullName: "Tester John",
		Email:    "tester@go4movies.com",
		Password: "Password",
		Phone:    strPtr("999999999"),
	}
	DB.Create(&user)
	DB.Create(&showtimes)
	log.Printf("Seeded %d showtimes", len(showtimes))

	// -------------------------------------------------------
	// 7. SAMPLE BOOKINGS (mark a few seats as sold so the
	//    seat map isn't empty on first load)
	// -------------------------------------------------------
	type soldSeat struct {
		Row string
		Col int
	}
	sampleBookings := []struct {
		ShowtimeIdx int
		Seats       []soldSeat
		Status      string
	}{
		{0, []soldSeat{{"A", 3}, {"A", 4}, {"A", 9}, {"A", 10}}, "CONFIRMED"},
		{0, []soldSeat{{"C", 4}, {"C", 11}}, "CONFIRMED"},
		{1, []soldSeat{{"B", 2}, {"B", 3}, {"B", 4}}, "CONFIRMED"},
		{1, []soldSeat{{"D", 7}, {"D", 8}}, "CONFIRMED"},
		{3, []soldSeat{{"A", 1}, {"A", 2}, {"C", 6}, {"C", 7}}, "CONFIRMED"},
		{4, []soldSeat{{"E", 3}, {"E", 4}, {"E", 5}}, "CONFIRMED"},
		{6, []soldSeat{{"A", 5}, {"A", 6}, {"B", 5}, {"B", 6}}, "CONFIRMED"},
		{9, []soldSeat{{"D", 1}, {"D", 2}, {"D", 3}}, "CONFIRMED"},
	}

	for i, sb := range sampleBookings {
		st := showtimes[sb.ShowtimeIdx]

		var seatRecords []models.Seat
		for _, s := range sb.Seats {
			var seat models.Seat
			if err := DB.Where("screen_id = ? AND row_label = ? AND col_number = ?",
				st.ScreenID, s.Row, s.Col).First(&seat).Error; err == nil {
				seatRecords = append(seatRecords, seat)
			}
		}
		if len(seatRecords) == 0 {
			continue
		}

		var total float64
		for _, seat := range seatRecords {
			total += seat.BasePrice * st.PriceMultiplier
		}

		booking := models.Booking{
			UserID:      user.ID,
			ShowtimeID:  st.ID,
			BookingRef:  fmt.Sprintf("SEED-%04d", i+1),
			Status:      sb.Status,
			TotalAmount: total,
			PaymentStatus: "PAID",
			BookedAt:    time.Now(),
		}
		DB.Create(&booking)

		for _, seat := range seatRecords {
			DB.Create(&models.BookingSeat{
				BookingID:  booking.ID,
				SeatID:     seat.ID,
				ShowtimeID: st.ID,
				SeatPrice:  seat.BasePrice * st.PriceMultiplier,
			})
		}
	}
	log.Println("Seeded sample bookings with sold seats")

	log.Println("Database seeding completed successfully")
}
