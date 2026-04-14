# Sprint 3 Report — Go4Movies

Project URL: https://github.com/parasmittal099/Go4Movies
**Sprint Duration:** Mar 27, 2026 – Apr 14, 2026
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

### US-12: Checkout Preview & Price Breakdown
> *As a user, I want to see a detailed price breakdown (subtotal, convenience fee, tax, discount, total) before confirming my booking so I can review costs.*

### US-13: Confirm Booking & Mock Payment
> *As a user, I want to confirm my seat selection and complete a mock payment so I receive a booking reference.*

### US-14: Concurrent Booking Safety
> *As a user, I want to be sure that if two people try to book the same seat at the same time, only one succeeds, so I never get a double-booked ticket.*

### US-15: Seat Hold & Automatic Expiry
> *As a user, I want my selected seats held for 10 minutes during checkout so no one else can take them, and released automatically if I don't complete payment.*

### US-16: Payment Checkout Page (Frontend)
> *As a user, I want a payment page with card entry, real-time pricing, a countdown timer, and discount code support so I can complete my booking end-to-end.*

### US-17: Concurrency & Checkout Test Coverage
> *As a developer, I want unit tests covering checkout pricing, double-booking prevention, concurrent access, and seat hold expiry so the booking flow is reliable.*

---

## Issues Planned

| #  | Issue / Task                                                        | Assignee         | PR(s)    | Status |
| -- | ------------------------------------------------------------------- | ---------------- | -------- | ------ |
| 1  | Checkout preview API (`POST /checkout/preview`)                     | Paras Mittal     | #28      | Done   |
| 2  | Checkout confirm API (`POST /checkout/confirm`)                     | Omkar Salkade    | #28      | Done   |
| 3  | Transaction-wrapped booking with re-validation                      | Paras Mittal     | #29      | Done   |
| 4  | Unique index on `(showtime_id, seat_id)` in BookingSeat             | Omkar Salkade    | #29      | Done   |
| 5  | Booking `ExpiresAt` field & 10-minute hold window                   | Paras Mittal     | #29      | Done   |
| 6  | Background booking expiry worker (`workers/expiry.go`)              | Omkar Salkade    | #29      | Done   |
| 7  | Seats API excludes expired PENDING bookings                         | Paras Mittal     | #29      | Done   |
| 8  | Payment checkout page (frontend)                                    | Harsh Soni       | #30      | Done   |
| 9  | Checkout API client functions & TypeScript types                    | Paresh Devlekar  | #30      | Done   |
| 10 | Seat selection → payment page navigation                            | Paresh Devlekar  | #30      | Done   |
| 11 | Checkout & concurrency unit tests (backend)                         | Paras Mittal     | #31      | Done   |
| 12 | Test infrastructure: shared-cache SQLite for concurrent tests       | Omkar Salkade    | #31      | Done   |

---

## Detail of Work Completed

### 1. US-12 & US-13 — Checkout Preview & Confirm APIs

**Backend (`handlers/checkout.go`, `routes/routes.go`):**

- `POST /api/v1/checkout/preview` — accepts `user_id`, `showtime_id`, `seat_ids[]`, optional `discount_code`; validates the user, showtime, and seat ownership; checks seat availability; returns a full price quote with per-seat line items, subtotal, $2.00 convenience fee, 8% tax, discount, and total due
- `POST /api/v1/checkout/confirm` — same request shape; creates a `Booking` (PENDING → CONFIRMED), `BookingSeat` rows, and a `Payment` record inside a single DB transaction; returns `booking_ref` and `payment_id`
- Mock discount code `MOCK100` waives the entire total for testing
- Helper functions: `validateRequest()` (stateless checks), `buildQuote()` (pure price math), `checkConflicts()` (availability query), `newBookingRef()` (crypto-random `G4M-*` reference)

### 2. US-14 — Concurrent Booking Safety

**Three layers of protection in `handlers/checkout.go` and `models/booking.go`:**

1. **Transaction wrapping** — `ConfirmCheckout` wraps all writes (booking, booking_seats, payment, status update) inside `database.DB.Transaction()`. Any failure rolls back everything.
2. **Re-validation inside the transaction** — `checkConflicts()` is called inside the write transaction, serialised by SQLite's database-level write lock, eliminating the time-of-check-to-time-of-use (TOCTOU) race condition.
3. **Database-level unique index** — `BookingSeat` now has `uniqueIndex:idx_showtime_seat` on `(showtime_id, seat_id)`. Even if application logic misses a conflict, the database rejects the duplicate insert.

**Error handling:** Both `errConflict` (from re-validation) and `UNIQUE constraint failed` (from the index) return `409 Conflict` to the client.

### 3. US-15 — Seat Hold & Automatic Expiry

**Model change (`models/booking.go`):**
- New `ExpiresAt *time.Time` field on `Booking` with a database index for efficient sweeps

**Background worker (`workers/expiry.go`):**
- `StartBookingExpiry()` launches a goroutine with a 1-minute ticker
- `expireStaleBookings()` updates all PENDING bookings where `expires_at <= now` to status `EXPIRED` / payment_status `EXPIRED`
- Started in `main.go` after database setup, before the HTTP server

**Seats API update (`handlers/showtime.go`):**
- The booked-seats query now adds: `bookings.status = 'CONFIRMED' OR bookings.expires_at IS NULL OR bookings.expires_at > now`
- Expired PENDING bookings no longer block seats on the seat map

**Checkout flow (`handlers/checkout.go`):**
- `ConfirmCheckout` sets `ExpiresAt` to `now + 10 minutes` on the initial PENDING booking
- On successful payment, `ExpiresAt` is cleared (`nil`) as the booking moves to CONFIRMED

### 4. US-16 — Payment Checkout Page (Frontend)

**New page (`app/payment/page.tsx`, 383 lines):**
- Card payment form: Name on Card, Card Number, Expiry, CVV with VISA/MC/AMEX badge
- Real-time checkout preview: calls `previewCheckout()` on mount, displays subtotal, convenience fee, tax, discount, and total due
- Discount/promo code input with "Apply" button — re-fetches preview with the applied code
- 10-minute countdown timer (matching the backend hold window) with `MM:SS` display
- Booking summary sidebar: movie title, format, language, theater, screen, date, time, selected seats with "Change" link
- "Pay $X.XX Securely" button calls `confirmCheckout()`, disables during submission
- Success state: displays booking reference in green
- Error states: validation errors (missing login, invalid card), seat conflict (409), and generic API failures
- SSL encryption notice for trust signaling

**API client additions (`lib/api.ts`):**
- `previewCheckout(payload)` → `POST /api/v1/checkout/preview`
- `confirmCheckout(payload)` → `POST /api/v1/checkout/confirm`
- `CheckoutPayload` interface with `user_id`, `showtime_id`, `seat_ids[]`, optional `discount_code`

**Type definitions (`lib/types.ts`):**
- `CheckoutLineItem`, `CheckoutTotals`, `CheckoutQuote`, `CheckoutConfirmResponse`

**Navigation (`app/movies/[id]/page.tsx`):**
- "Proceed to Payment" in seat selection now routes to `/payment?showtimeId=&seats=&seatIds=` with real selected seat data

### 5. US-17 — Concurrency & Checkout Test Coverage

**Test infrastructure (`testutil/testutil.go`):**
- Upgraded from `:memory:` to `file:testdb_N?mode=memory&cache=shared` with atomic counter for unique naming — enables concurrent goroutines to operate on the same in-memory database while maintaining test isolation
- Added `PRAGMA busy_timeout = 5000` for SQLite write-lock tolerance

---

## Frontend Tests (Jest + React Testing Library)

**13 test files, 106 test cases** — unchanged from Sprint 2. All existing tests continue to pass.

No new frontend test files were added in Sprint 3. The payment page (`app/payment/page.tsx`) is a candidate for test coverage in Sprint 4.

---

## Backend Unit Tests (Go `testing` package)

**17 test files, 104 tests** across 7 packages. Run with `go test ./... -v`.

### New in Sprint 3

#### `handlers/checkout_test.go` (15 tests)

| # | Test Case | What It Verifies |
|---|-----------|-----------------|
| 1 | TestPreviewCheckout_Success_NoDiscount | 200 with correct subtotal ($100), fee ($2), tax ($8.16), total ($110.16) |
| 2 | TestPreviewCheckout_Mock100_ZeroTotal | MOCK100 discount zeroes out total_due |
| 3 | TestPreviewCheckout_TwoSeats | Line items array contains both seats |
| 4 | TestPreviewCheckout_UserNotFound | 404 for nonexistent user_id |
| 5 | TestPreviewCheckout_ShowtimeNotFound | 404 for nonexistent showtime_id |
| 6 | TestPreviewCheckout_InactiveShowtime | 400 when showtime is_active = false |
| 7 | TestPreviewCheckout_InvalidSeatWrongScreen | 400 when seat belongs to a different screen |
| 8 | TestPreviewCheckout_DuplicateSeatID | 400 on duplicate seat_id in request |
| 9 | TestPreviewCheckout_SeatAlreadyBooked | 409 when seat has a CONFIRMED booking |
| 10 | TestPreviewCheckout_EmptyBody | 400 on missing required fields |
| 11 | TestConfirmCheckout_CreatesBookingAndPayment | 201 with booking (CONFIRMED/PAID), 1 booking_seat, 1 payment |
| 12 | TestConfirmCheckout_ThenPreviewConflict | Confirm then preview same seat → 409 |
| 13 | TestConfirmCheckout_ConcurrentSameSeat_OnlyOneWins | 3 goroutines race for same seat; at most 1 booking_seat row exists |
| 14 | TestConfirmCheckout_SequentialDoubleBook_Rejected | First confirm 201, second confirm 409, exactly 1 row in DB |
| 15 | TestPreviewCheckout_ExpiredPendingSeatsAreAvailable | Expired PENDING booking does not block preview |

#### `workers/expiry_test.go` (1 test)

| # | Test Case | What It Verifies |
|---|-----------|-----------------|
| 1 | TestExpireStaleBookings | Expired PENDING → EXPIRED; valid PENDING stays PENDING; CONFIRMED untouched |

### Sprint 2 Tests (unchanged, still passing)

| File | Tests |
|------|-------|
| `config/config_test.go` | 4 |
| `database/database_test.go` | 4 |
| `database/seed_test.go` | 9 |
| `handlers/auth_test.go` | 8 |
| `handlers/location_test.go` | 3 |
| `handlers/movie_test.go` | 8 |
| `handlers/showtime_test.go` | 15 |
| `middleware/cors_test.go` | 2 |
| `models/user_test.go` | 6 |
| `models/location_test.go` | 5 |
| `models/movie_test.go` | 6 |
| `models/theater_test.go` | 7 |
| `models/showtime_test.go` | 4 |
| `models/booking_test.go` | 6 |
| `routes/routes_test.go` | 2 |

**Backend Test Total: 104 tests across 17 files, all passing**

---

## Backend API Documentation

Base URL: `http://localhost:8080`

### Existing Endpoints (Sprint 1 & 2)

| # | Method | Endpoint | Description |
|---|--------|----------|-------------|
| 1 | GET | `/api/v1/locations` | All supported locations |
| 2 | GET | `/api/v1/movies` | Active movies (optional `?zipcode=`, `?q=`) |
| 3 | GET | `/api/v1/movies/:id` | Single movie by ID |
| 4 | GET | `/api/v1/movies/:id/showtimes` | Showtimes grouped by theater (`?zipcode=` required) |
| 5 | GET | `/api/v1/seats` | Seat availability (`?showtime_id=` required) |
| 6 | POST | `/api/v1/auth/register` | Create user account |
| 7 | POST | `/api/v1/auth/login` | Authenticate user |

### New Endpoints (Sprint 3)

#### 8. `POST /api/v1/checkout/preview`

Returns a price breakdown without creating a booking.

**Request body:**
```json
{
  "user_id": 1,
  "showtime_id": 1,
  "seat_ids": [1, 2],
  "discount_code": "MOCK100"
}
```

| Field | Required | Description |
|-------|----------|-------------|
| `user_id` | Yes | Authenticated user ID |
| `showtime_id` | Yes | Target showtime |
| `seat_ids` | Yes | Array of seat IDs (min 1, no duplicates) |
| `discount_code` | No | Promo code (`MOCK100` waives full total) |

**Response** `200 OK`:
```json
{
  "showtime_id": 1,
  "user_id": 1,
  "line_items": [
    {
      "seat_id": 1,
      "row_label": "A",
      "col_number": 1,
      "seat_type": "Premium",
      "unit_price": 27.00
    }
  ],
  "totals": {
    "subtotal": 27.00,
    "convenience_fee": 2.00,
    "tax_amount": 2.32,
    "discount_code": "MOCK100",
    "discount_amount": 31.32,
    "total_due": 0
  }
}
```

**Error responses:**
- `400` — missing/invalid fields, duplicate seat IDs, inactive showtime, seat not on correct screen
- `404` — user or showtime not found
- `409` — one or more seats already booked/reserved

---

#### 9. `POST /api/v1/checkout/confirm`

Creates a booking, assigns seats, records payment, and confirms — all atomically.

**Request body:** Same as preview.

**Response** `201 Created`:
```json
{
  "message": "Booking confirmed",
  "booking_id": 1,
  "booking_ref": "G4M-a1b2c3d4e5f6g7h8",
  "quote": { ... },
  "payment_id": 1
}
```

**Error responses:**
- `400` — validation errors (same as preview)
- `404` — user or showtime not found
- `409` — seat conflict (detected by re-validation or unique index)
- `500` — internal error (transaction rolled back, no partial data)

**Price calculation:**
- `unit_price = seat.base_price × showtime.price_multiplier`
- `subtotal = sum(unit_prices)`
- `convenience_fee = $2.00` (flat per order)
- `tax = (subtotal + convenience_fee) × 8%`
- `total_due = subtotal + convenience_fee + tax - discount`

---

## Summary

Sprint 3 delivered the complete booking and payment flow end-to-end. Users can now select seats, preview pricing with a detailed breakdown, apply discount codes, and confirm bookings through a polished payment page with a 10-minute countdown timer. The backend checkout is protected against double-booking with three layers of defense: transaction wrapping, in-transaction re-validation, and a database-level unique constraint. A background expiry worker automatically releases abandoned seat holds. The backend test suite grew from 88 to 104 tests, with 15 new checkout tests covering pricing math, validation, sequential double-booking rejection, concurrent race conditions, and seat hold expiry. All 104 backend tests and 106 frontend tests pass.
