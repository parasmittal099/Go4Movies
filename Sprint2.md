# Sprint 2 Report — Go4Movies

Project URL: https://github.com/parasmittal099/Go4Movies
**Sprint Duration:** Feb 19, 2026 – Mar 26, 2026
**Project Board:** https://github.com/users/parasmittal099/projects/2

## Team Members

| Member            | Role     |
| ----------------- | -------- |
| Harsh Soni        | Frontend |
| Paresh Devlekar   | Frontend |
| Omkar Salkade     | Backend  |
| Paras Mittal      | Backend  |

---

## User Stories

### US-06: View Showtimes Grouped by Theater
> *As a user, I want to see available showtimes for a movie grouped by theater and date so I can pick a convenient time and location.*

### US-07: Select Seats for a Showtime
> *As a user, I want to see a visual seat map for my chosen showtime, select seats, and see the running total so I can proceed to checkout.*

### US-08: Seat Map Reflects Real Availability
> *As a user, I want sold seats to appear as unavailable on the seat map so I only select open seats.*

### US-09: Frontend Unit & Integration Tests
> *As a developer, I want comprehensive frontend tests so the UI is reliable and regressions are caught early.*

### US-10: Backend Unit Test Coverage
> *As a developer, I want 1:1 unit tests for every backend package so I can maintain code quality and catch regressions.*

### US-11: Developer Experience — Auto-Reseed & Deterministic Data
> *As a developer, I want the database to reset on every container restart with deterministic data so the team always works with a predictable dataset.*

---

## Issues Planned

| #  | Issue / Task                                                  | Assignee         | PR(s)      | Status |
| -- | ------------------------------------------------------------- | ---------------- | ---------- | ------ |
| 1  | Seat selection component & test page                          | Harsh Soni       | #16        | Done   |
| 2  | Showtimes API for a movie                                     | Paras Mittal     | #17        | Done   |
| 3  | Wire showtimes API to movie detail page                       | Paresh Devlekar  | #18        | Done   |
| 4  | Expand seed data — more showtimes & dates                     | Omkar Salkade    | #19        | Done   |
| 5  | Seat availability API (`GET /api/v1/seats`)                   | Paras Mittal     | #20        | Done   |
| 6  | Dynamic seat map integration with seats API                   | Paresh Devlekar  | #21        | Done   |
| 7  | Fix SQLite WAL mode & deterministic seed dates                | Paras Mittal     | #22        | Done   |
| 8  | Dockerfile auto-reseed on container restart                   | Omkar Salkade    | #23        | Done   |
| 9  | Backend unit tests (Sprint 1 + Sprint 2 code)                 | Paras Mittal     | #24, #25   | Done   |
| 10 | Three-block seat layout with aisles & section labels          | Harsh Soni       | #26        | Done   |
| 11 | Frontend unit & integration tests (Jest + RTL)                | Harsh Soni       | #27        | Done   |

---

## Detail of Work Completed

### 1. US-06 — View Showtimes Grouped by Theater

**Backend (PR #17, #19):**
- `GET /api/v1/movies/:id/showtimes?zipcode=` — returns showtimes grouped by theater, with screen name, format (IMAX / 2D / 4DX), and price multiplier
- Expanded seed data from ~25 showtimes to 115 showtimes across 5 days with multiple time slots per theater
- Fixed base date (`2026-03-24`) for deterministic data across all developer machines

**Frontend (PR #18):**
- Movie detail page (`movies/[id]/page.tsx`) dynamically fetches showtimes via `fetchMovieShowtimes()`
- Date tab selector with sorted show dates
- Showtimes rendered by theater group with timing slots showing format, screen type, and start time
- "Book Tickets" button on each showtime triggers seat selection flow

### 2. US-07 — Select Seats for a Showtime

**Backend (PR #20):**
- `GET /api/v1/seats?showtime_id=&seat_type=&status=` — returns all seats with row, column, type (Premium / VIP), computed price (`base_price × price_multiplier`), and live status (AVAILABLE / RESERVED / BOOKED)
- Optional query filters: `seat_type` and `status`
- Response includes layout metadata (`total_rows`, `total_cols`) and availability summary

**Frontend (PR #16, #21, #26):**
- Full `SeatSelection` component (`components/booking/seat-selection.tsx`) with:
  - Three-block layout with aisles, matching the real theater screen
  - Section labels: "PREMIUM · FRONT", "EXECUTIVE · BEST VIEW", "VIP · RECLINER"
  - Color-coded seat buttons (gray / yellow / purple) with hover price tooltips
  - Click to select / deselect with max 20-seat limit and toast warning
  - Sticky bottom bar showing selected seats, ticket count, and running total
  - "Proceed to Payment" button (step 2 in the booking flow stepper)
  - Loading, error, and empty states for the async seat fetch
- TypeScript types: `SeatAPI`, `SeatsResponse`, `ShowtimeEntry`, `TheaterGroup`, `ShowtimesResponse` in `lib/types.ts`
- API client functions: `fetchMovieShowtimes()` and `fetchShowtimeSeats()` in `lib/api.ts`

### 3. US-08 — Seat Map Reflects Real Availability

**Backend (PR #20):**
- Seats API joins `booking_seats` + `bookings` tables to determine live status per showtime
- CONFIRMED bookings → BOOKED, PENDING bookings → RESERVED
- Sample bookings seeded (`SEED-0001` through `SEED-0008`) marking ~25 seats as BOOKED across 8 showtimes

**Frontend (PR #21, #26):**
- `mapApiSeat()` maps backend `BOOKED` / `RESERVED` to the `sold` status
- Sold seats render with an ✕ icon, `cursor-not-allowed`, and reduced opacity
- Clicking a sold seat is a no-op

### 4. US-11 — Developer Experience & DevOps

**PR #22, #23:**
- Disabled SQLite WAL mode to prevent data corruption on Docker Desktop
- Dockerfile CMD: `rm -f /app/data/*.db ... && go run main.go` — auto-wipes and reseeds on every container restart
- Fixed base date for all seed showtimes so data is deterministic across machines
- `.gitignore` updated to exclude `*.db-shm` and `*.db-wal` files

---

## Frontend Tests (Jest + React Testing Library)

**13 test files** (PR #27). No Cypress tests — all tests use Jest with React Testing Library.

### `components/booking/__tests__/seat-selection.test.tsx` (13 tests)
| # | Test Case |
|---|-----------|
| 1 | shows loading state initially |
| 2 | renders seat grid after loading |
| 3 | renders booked seat as disabled |
| 4 | selects an available seat on click |
| 5 | updates total price when seats are selected |
| 6 | calls onBack when back arrow button is clicked |
| 7 | calls onChangeMovie when Change Movie is clicked |
| 8 | calls onProceed with selected seats when Proceed button is clicked |
| 9 | Proceed button is disabled when no seats are selected |
| 10 | shows error message on API failure |
| 11 | shows error for invalid showtimeId (0) |
| 12 | shows no seats message when seats array is empty |
| 13 | renders movie title in header |

### `components/movies/__tests__/movie-detail.test.tsx` (11 tests)
| # | Test Case |
|---|-----------|
| 1 | renders the movie title |
| 2 | renders movie description |
| 3 | renders date buttons for provided dates |
| 4 | renders theater name |
| 5 | renders showtime buttons |
| 6 | shows loading skeletons when showtimesLoading is true |
| 7 | shows no showtimes message when cinemas is empty |
| 8 | selecting a timing shows the sticky booking bar |
| 9 | calls onBookTickets when Select Seats is clicked |
| 10 | changes selected date when another date is clicked |
| 11 | renders null showtimes correctly (loading state initial) |

### `components/movies/__tests__/movie-grid.test.tsx` (4 tests)
| # | Test Case |
|---|-----------|
| 1 | renders all movie cards |
| 2 | renders the correct number of movie cards |
| 3 | renders empty grid when no movies |
| 4 | calls onMovieClick with the correct movie id |

### `components/movies/__tests__/movie-card.test.tsx` (4 tests)
| # | Test Case |
|---|-----------|
| 1 | renders the movie title |
| 2 | renders the movie description |
| 3 | renders the poster image with correct src and alt |
| 4 | calls onClick with the movie id when clicked |

### `components/auth/__tests__/login-form.test.tsx` (8 tests)
| # | Test Case |
|---|-----------|
| 1 | renders email and password inputs |
| 2 | renders Sign In button |
| 3 | shows error when email and password are empty |
| 4 | shows error when only email is provided |
| 5 | calls loginUser API on valid input |
| 6 | shows server error message on login failure |
| 7 | toggles password visibility |
| 8 | disables submit button while submitting |

### `components/auth/__tests__/signup-form.test.tsx` (8 tests)
| # | Test Case |
|---|-----------|
| 1 | renders all form fields |
| 2 | shows error when full name is missing |
| 3 | shows error when username is too short |
| 4 | shows error when password is too short |
| 5 | shows error when passwords do not match |
| 6 | calls registerUser with correct payload on valid submit |
| 7 | shows server error message on registration failure |
| 8 | toggles password visibility |

### `context/__tests__/auth-context.test.tsx` (9 tests)
| # | Test Case |
|---|-----------|
| 1 | becomes ready after mount |
| 2 | starts with no user |
| 3 | login sets user in state and localStorage |
| 4 | logout clears user from state and localStorage |
| 5 | rehydrates user from localStorage on mount |
| 6 | setLocation persists location to localStorage |
| 7 | clearLocation removes location from state and storage |
| 8 | migrates legacy selectedZipCode key to go4movies_location |
| 9 | throws when used outside AuthProvider |

### `context/__tests__/movies-context.test.tsx` (4 tests)
| # | Test Case |
|---|-----------|
| 1 | starts with empty movies array |
| 2 | setMovies updates the movies list |
| 3 | clears movies list when setMovies is called with [] |
| 4 | throws when used outside MoviesProvider |

### `lib/__tests__/api.test.ts` (16 tests)
| # | Test Case | API Function |
|---|-----------|-------------|
| 1 | returns locations on success | `fetchLocations` |
| 2 | throws on error response | `fetchLocations` |
| 3 | returns movies for a given zip code | `fetchMoviesByZipCode` |
| 4 | throws on error | `fetchMoviesByZipCode` |
| 5 | returns a single movie by ID | `fetchMovieById` |
| 6 | throws Movie not found on 404 | `fetchMovieById` |
| 7 | returns auth response on success | `registerUser` |
| 8 | throws with server error message | `registerUser` |
| 9 | falls back to default message when no error field | `registerUser` |
| 10 | returns auth response on success | `loginUser` |
| 11 | throws with server error message on failure | `loginUser` |
| 12 | falls back to default login failure message | `loginUser` |
| 13 | returns showtime response | `fetchMovieShowtimes` |
| 14 | throws on failure | `fetchMovieShowtimes` |
| 15 | returns seats response | `fetchShowtimeSeats` |
| 16 | throws on failure | `fetchShowtimeSeats` |

### `components/layout/__tests__/header.test.tsx` (9 tests)
| # | Test Case |
|---|-----------|
| 1 | renders the Go4Movies logo link |
| 2 | shows Sign In link when not logged in |
| 3 | does not show Sign Out button when not logged in |
| 4 | shows user initials when logged in |
| 5 | shows full name when logged in |
| 6 | clicking Sign Out calls logout and redirects to / |
| 7 | opens search input when search icon button is clicked |
| 8 | closes search on Escape key |
| 9 | clears query when close (x) button inside search is clicked |

### `components/ui/__tests__/search-bar.test.tsx` (5 tests)
| # | Test Case |
|---|-----------|
| 1 | renders input and button |
| 2 | updates input value on change |
| 3 | calls alert with search query on form submit when query is non-empty |
| 4 | does not call alert when query is empty |
| 5 | does not call alert when query is only whitespace |

### `components/ui/__tests__/zip-code-modal.test.tsx` (9 tests)
| # | Test Case |
|---|-----------|
| 1 | renders modal heading |
| 2 | shows search input |
| 3 | does NOT render close button when no zip code is stored |
| 4 | renders close button when a zip code is stored |
| 5 | calls onClose when close button is clicked |
| 6 | filters locations when user types in search |
| 7 | filters by zip code prefix |
| 8 | calls onSelectLocation when a result is clicked |
| 9 | handles API fetch error gracefully |

### `hooks/__tests__/useZipCode.test.ts` (6 tests)
| # | Test Case |
|---|-----------|
| 1 | getZipCode returns null when nothing stored |
| 2 | setZipCode stores value in localStorage |
| 3 | getZipCode returns stored zip code |
| 4 | clearZipCode removes stored zip code |
| 5 | requireZipCode returns zip code if set |
| 6 | requireZipCode redirects to "/" when no zip code set |

**Frontend Test Total: 106 test cases across 13 files**

---

## Backend Unit Tests (Go `testing` package)

**15 test files, 88 tests** across 6 packages (PR #24, #25). Run with `go test ./... -v`.

### `config/config_test.go` (4 tests)
| # | Test Case | What It Verifies |
|---|-----------|-----------------|
| 1 | TestLoadConfig_Defaults | Default values when no env vars set |
| 2 | TestLoadConfig_EnvOverrides | Env vars override defaults |
| 3 | TestGenENV_Fallback | Fallback returned for missing key |
| 4 | TestGenENV_EnvSet | Env value returned when present |

### `database/database_test.go` (4 tests)
| # | Test Case | What It Verifies |
|---|-----------|-----------------|
| 1 | TestConnect_CreatesDirectory | DB directory auto-created |
| 2 | TestConnect_SetsDB | Package-level DB is non-nil and pingable |
| 3 | TestMigrate_CreatesAllTables | All 10 tables exist after migration |
| 4 | TestMigrate_Idempotent | Running Migrate twice doesn't panic |

### `database/seed_test.go` (9 tests)
| # | Test Case | What It Verifies |
|---|-----------|-----------------|
| 1 | TestSeed_PopulatesLocations | 7 locations seeded |
| 2 | TestSeed_PopulatesMovies | 10 movies seeded |
| 3 | TestSeed_PopulatesTheaters | 7 theaters seeded |
| 4 | TestSeed_PopulatesScreens | 7 screens seeded |
| 5 | TestSeed_PopulatesSeats | 51 seats per screen (357 total) |
| 6 | TestSeed_PopulatesShowtimes | 115 showtimes seeded |
| 7 | TestSeed_CreatesBookings | Sample bookings + booking_seats exist |
| 8 | TestSeed_CreatesUser | Seed user `tester@go4movies.com` exists |
| 9 | TestSeed_Idempotent | Second Seed() call doesn't duplicate data |

### `handlers/auth_test.go` (8 tests)
| # | Test Case | What It Verifies |
|---|-----------|-----------------|
| 1 | TestRegister_Success | 201 Created, user in DB |
| 2 | TestRegister_MissingFields | 400 on incomplete payload |
| 3 | TestRegister_DuplicateEmail | 409 Conflict on duplicate email |
| 4 | TestRegister_EmailNormalized | Email lowercased before storage |
| 5 | TestLogin_Success | 200 OK with valid credentials |
| 6 | TestLogin_WrongPassword | 401 on wrong password |
| 7 | TestLogin_NonexistentUser | 401 for unknown email |
| 8 | TestLogin_MissingFields | 400 on empty body |

### `handlers/location_test.go` (3 tests)
| # | Test Case | What It Verifies |
|---|-----------|-----------------|
| 1 | TestGetLocations_ReturnsAll | Returns all seeded locations |
| 2 | TestGetLocations_Empty | Returns empty list on empty DB |
| 3 | TestGetLocations_OrderedByCityZip | Results sorted by city, zipcode |

### `handlers/movie_test.go` (8 tests)
| # | Test Case | What It Verifies |
|---|-----------|-----------------|
| 1 | TestListMovies_NoFilter | Only active movies returned |
| 2 | TestListMovies_WithZipcode | Filters by zip → city → theater chain |
| 3 | TestListMovies_UnknownZipcode | Empty list for unknown zip |
| 4 | TestListMovies_SearchByTitle | `?q=` matches title |
| 5 | TestListMovies_SearchByGenre | `?q=` matches genre |
| 6 | TestListMovies_SearchNoMatch | Empty list on no match |
| 7 | TestGetMovie_Found | 200 with correct movie |
| 8 | TestGetMovie_NotFound | 404 for invalid ID |

### `handlers/showtime_test.go` (12 tests)
| # | Test Case | What It Verifies |
|---|-----------|-----------------|
| 1 | TestGetMovieShowtimes_Success | 200 with theaters and dates |
| 2 | TestGetMovieShowtimes_MultipleShowtimesPerTheater | Multiple entries grouped correctly |
| 3 | TestGetMovieShowtimes_MovieNotFound | 404 for invalid movie |
| 4 | TestGetMovieShowtimes_MissingZipcode | 400 when zipcode missing |
| 5 | TestGetMovieShowtimes_UnknownZipcode | Empty theaters for unknown zip |
| 6 | TestGetMovieShowtimes_NoShowtimesInCity | Empty when no showtimes in city |
| 7 | TestGetMovieShowtimes_ShowtimeFieldValues | Correct field mapping in response |
| 8 | TestGetShowtimeSeats_Success | 200 with all seats and summary |
| 9 | TestGetShowtimeSeats_MissingShowtimeID | 400 when showtime_id missing |
| 10 | TestGetShowtimeSeats_ShowtimeNotFound | 404 for invalid showtime |
| 11 | TestGetShowtimeSeats_FilterBySeatType | `?seat_type=VIP` filters correctly |
| 12 | TestGetShowtimeSeats_BookedSeatStatus | Booked seats have status BOOKED |
| 13 | TestGetShowtimeSeats_FilterByStatus | `?status=BOOKED` returns only booked |
| 14 | TestGetShowtimeSeats_PriceCalculation | Price = base_price × price_multiplier |
| 15 | TestGetShowtimeSeats_ShowtimeMetadata | Response contains showtime/layout info |

### `middleware/cors_test.go` (2 tests)
| # | Test Case | What It Verifies |
|---|-----------|-----------------|
| 1 | TestCORS_SetsHeaders | All CORS headers present |
| 2 | TestCORS_OptionsReturns204 | OPTIONS preflight returns 204 |

### `models/user_test.go` (5 tests)
| # | Test Case | What It Verifies |
|---|-----------|-----------------|
| 1 | TestUser_Create | Insert + auto-generated ID |
| 2 | TestUser_UniqueEmail | Unique constraint on email |
| 3 | TestUser_UniqueUsername | Unique constraint on username |
| 4 | TestUser_Read | Query by email |
| 5 | TestUser_Update | Update full_name |
| 6 | TestUser_Delete | Soft delete |

### `models/location_test.go` (5 tests)
| # | Test Case | What It Verifies |
|---|-----------|-----------------|
| 1 | TestLocation_Create | Insert + ID |
| 2 | TestLocation_UniqueZipcode | Unique constraint on zipcode |
| 3 | TestLocation_Read | Query by zipcode |
| 4 | TestLocation_Update | Update city |
| 5 | TestLocation_Delete | Delete record |

### `models/movie_test.go` (6 tests)
| # | Test Case | What It Verifies |
|---|-----------|-----------------|
| 1 | TestMovie_Create | Insert + ID |
| 2 | TestMovie_OptionalFields | Nullable fields persisted |
| 3 | TestMovie_Read | Query by title |
| 4 | TestMovie_Update | Update title |
| 5 | TestMovie_Delete | Delete record |
| 6 | TestMovie_FilterActive | `is_active` filter works |

### `models/theater_test.go` (7 tests)
| # | Test Case | What It Verifies |
|---|-----------|-----------------|
| 1 | TestTheater_Create | Insert with FK to location |
| 2 | TestTheater_WithAddress | Optional address field |
| 3 | TestScreen_Create | Insert with FK to theater |
| 4 | TestScreen_BelongsToTheater | Preload Theater association |
| 5 | TestSeat_Create | Insert with FK to screen |
| 6 | TestSeat_UniqueSeatPerScreen | Composite unique (screen, row, col) |
| 7 | TestTheater_HasManyScreens | Preload Screens association |

### `models/showtime_test.go` (4 tests)
| # | Test Case | What It Verifies |
|---|-----------|-----------------|
| 1 | TestShowtime_Create | Insert with movie + screen FKs |
| 2 | TestShowtime_BelongsToMovie | Preload Movie association |
| 3 | TestShowtime_BelongsToScreen | Preload Screen association |
| 4 | TestShowtime_DefaultPriceMultiplier | Default 1.0 multiplier |

### `models/booking_test.go` (6 tests)
| # | Test Case | What It Verifies |
|---|-----------|-----------------|
| 1 | TestBooking_Create | Insert with user + showtime FKs |
| 2 | TestBooking_UniqueRef | Unique constraint on booking_ref |
| 3 | TestBookingSeat_Create | Insert with booking + seat FKs |
| 4 | TestBooking_HasManySeats | Preload BookingSeats association |
| 5 | TestBooking_DefaultStatus | Default PENDING / UNPAID |
| 6 | TestPayment_Create | Insert payment with FK to booking |

### `routes/routes_test.go` (2 tests)
| # | Test Case | What It Verifies |
|---|-----------|-----------------|
| 1 | TestRegisterRoutes_AllEndpoints | All 7 routes registered |
| 2 | TestRegisterRoutes_CountRoutes | Route count ≥ 7 |

**Backend Test Total: 88 tests across 15 files, all passing**

---

## Backend API Documentation

Base URL: `http://localhost:8080`

### 1. `GET /api/v1/locations`

Returns all supported locations.

**Response** `200 OK`:
```json
{
  "locations": [
    { "id": 1, "zipcode": "33101", "city": "Miami", "state": "FL" },
    { "id": 2, "zipcode": "33139", "city": "Miami", "state": "FL" }
  ]
}
```

---

### 2. `GET /api/v1/movies`

Returns active movies, optionally filtered by location and/or search query.

| Query Param | Required | Description |
|-------------|----------|-------------|
| `zipcode`   | No       | Filter movies to those with showtimes near this zip code |
| `q`         | No       | Search by title or genre (LIKE match) |

**Response** `200 OK`:
```json
[
  {
    "id": 1,
    "title": "The Karate Kid",
    "description": "...",
    "genre": "Action,Drama,Family",
    "language": "English",
    "duration_min": 140,
    "rating": "PG",
    "poster_url": "https://image.tmdb.org/t/p/original/...",
    "cast": "Jackie Chan, Jaden Smith, ...",
    "trailer_url": "https://youtube.com/watch?v=...",
    "release_date": "2010-06-11T00:00:00Z",
    "is_active": true,
    "created_at": "2026-03-26T00:00:00Z"
  }
]
```

---

### 3. `GET /api/v1/movies/:id`

Returns a single movie by ID.

**Response** `200 OK`: Single movie object (same shape as above).
**Response** `404 Not Found`: `{ "error": "Movie not found" }`

---

### 4. `GET /api/v1/movies/:id/showtimes`

Returns showtimes for a movie grouped by theater.

| Query Param | Required | Description |
|-------------|----------|-------------|
| `zipcode`   | **Yes**  | Zip code to resolve city for theater lookup |

**Response** `200 OK`:
```json
{
  "movie_id": 1,
  "title": "The Karate Kid",
  "dates": ["2026-03-24", "2026-03-25"],
  "theaters": [
    {
      "theater_id": 1,
      "name": "AMC Aventura",
      "address": "19501 Biscayne Blvd",
      "showtimes": [
        {
          "id": 1,
          "show_date": "2026-03-24",
          "start_time": "10:00",
          "end_time": "12:20",
          "language": "English",
          "format": "IMAX",
          "price_multiplier": 1.5,
          "screen_name": "Screen 1",
          "screen_type": "IMAX"
        }
      ]
    }
  ]
}
```

**Response** `400 Bad Request`: `{ "error": "zipcode query parameter is required" }`
**Response** `404 Not Found`: `{ "error": "Movie not found" }`

---

### 5. `GET /api/v1/seats`

Returns seat availability for a specific showtime.

| Query Param   | Required | Description |
|---------------|----------|-------------|
| `showtime_id` | **Yes**  | ID of the showtime |
| `seat_type`   | No       | Filter by seat type (`Premium`, `VIP`) |
| `status`      | No       | Filter by status (`AVAILABLE`, `RESERVED`, `BOOKED`) |

**Response** `200 OK`:
```json
{
  "showtime": {
    "id": 1,
    "movie_title": "The Karate Kid",
    "screen_name": "Screen 1",
    "screen_type": "IMAX",
    "theater_name": "AMC Aventura",
    "show_date": "2026-03-24",
    "start_time": "10:00",
    "format": "IMAX",
    "language": "English"
  },
  "layout": {
    "total_rows": 5,
    "total_cols": 11
  },
  "seats": [
    {
      "id": 1,
      "row_label": "A",
      "col_number": 1,
      "seat_type": "Premium",
      "price": 27.00,
      "status": "AVAILABLE"
    },
    {
      "id": 3,
      "row_label": "A",
      "col_number": 3,
      "seat_type": "Premium",
      "price": 27.00,
      "status": "BOOKED"
    }
  ],
  "summary": {
    "total": 51,
    "available": 45,
    "reserved": 0,
    "booked": 6
  }
}
```

**Response** `400 Bad Request`: `{ "error": "showtime_id query parameter is required" }`
**Response** `404 Not Found`: `{ "error": "Showtime not found" }`

**Price calculation:** `price = seat.base_price × showtime.price_multiplier` (rounded to 2 decimal places)

**Seat layout per screen:**

| Rows | Count | Type    | Base Price |
|------|-------|---------|------------|
| A–D  | 11 each | Premium | $18.00   |
| E    | 7     | VIP     | $25.00     |
| **Total** | **51 seats per screen** | | |

---

### 6. `POST /api/v1/auth/register`

Creates a new user account.

**Request body:**
```json
{
  "email": "alice@example.com",
  "username": "alice",
  "password": "secret123",
  "full_name": "Alice Smith"
}
```

| Field     | Validation |
|-----------|------------|
| `email`    | Required, valid email format |
| `username` | Required, min 3 characters |
| `password` | Required, min 6 characters |
| `full_name`| Required |

**Response** `201 Created`:
```json
{
  "message": "User registered successfully",
  "user": { "id": 1, "email": "alice@example.com", "username": "alice", "full_name": "Alice Smith" }
}
```

**Response** `400 Bad Request`: `{ "error": "validation error details" }`
**Response** `409 Conflict`: `{ "error": "Email or username already taken" }`

---

### 7. `POST /api/v1/auth/login`

Authenticates a user.

**Request body:**
```json
{
  "email": "alice@example.com",
  "password": "secret123"
}
```

**Response** `200 OK`:
```json
{
  "message": "Login successful",
  "user": { "id": 1, "email": "alice@example.com", "username": "alice", "full_name": "Alice Smith" }
}
```

**Response** `400 Bad Request`: `{ "error": "validation error details" }`
**Response** `401 Unauthorized`: `{ "error": "Invalid email or password" }`

---

## What Didn't Get Completed & Why

### 1. JWT Token Generation & Session Management
| Status | Not implemented |
| ------ | --- |
| **Why** | Sprint 2 prioritized the booking flow (showtimes → seat selection) and test coverage. The auth middleware file (`middleware/auth.go`) remains empty. |
| **Impact** | No token-based session persistence or protected routes. Carry-over for Sprint 3. |

### 2. Password Hashing (bcrypt)
| Status | Not implemented |
| ------ | --- |
| **Why** | Carried over from Sprint 1. Sprint 2 focused on the booking UI and API integration rather than security hardening. |
| **Impact** | Passwords still stored in plain text. Will be resolved in Sprint 3. |

### 3. Payment Flow (Checkout & Confirmation)
| Status | Not started |
| ------ | --- |
| **Why** | The "Proceed to Payment" button is wired in the frontend but currently shows an `alert()` placeholder. The backend `Payment` model exists and is auto-migrated, but no payment API endpoints have been built. |
| **Impact** | Users can select seats but cannot complete a booking. This is the primary Sprint 3 deliverable. |

### 4. Booking Confirmation & History
| Status | Not started |
| ------ | --- |
| **Why** | Depends on the payment flow. The `Booking` and `BookingSeat` models are fully defined and used by the seed/seats API, but no user-facing booking creation or history endpoints exist. |
| **Impact** | No order confirmation page or "My Bookings" section. Planned for Sprint 3. |

---

## Summary

Sprint 2 delivered the core booking flow — from viewing showtimes grouped by theater and date, to selecting seats on a fully interactive, color-coded seat map that reflects real-time availability. The frontend and backend are now fully integrated end-to-end for the showtime → seat selection journey. The seed dataset was expanded to 115 showtimes across 5 days with sample bookings, and the Dockerfile now auto-reseeds for a consistent developer experience. Comprehensive test suites were added on both sides: **88 backend tests** (Go) and **106 frontend tests** across 13 files (Jest + RTL). All 7 REST API endpoints are documented with request/response examples. The outstanding items — JWT auth, bcrypt, payment processing, and booking confirmation — are well-scoped carry-overs for Sprint 3.
