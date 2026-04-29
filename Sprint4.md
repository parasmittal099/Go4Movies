# Sprint 4 Report — Go4Movies

Project URL: https://github.com/parasmittal099/Go4Movies
**Sprint Duration:** Apr 14, 2026 – Apr 28, 2026
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

### US-18: JWT Authentication on Login & Register
> *As a user, I want to receive a JSON Web Token when I log in or register so the frontend can authenticate subsequent API calls without re-entering credentials.*

### US-19: Protected Checkout Confirmation
> *As a user, I want my booking confirmation call to be secured with my JWT so that no one else can create bookings on my behalf.*

### US-20: Secure Booking History
> *As a user, I want to view my booking history using my JWT identity (with a query-param fallback for backward compatibility) so my data stays private.*

### US-21: Bcrypt Password Hashing & Legacy Migration
> *As a user, I want my password stored with bcrypt hashing, and any legacy plaintext passwords automatically migrated on startup, so my credentials are always protected.*

### US-22: QR Ticket Generation & Lookup
> *As a user, I want a unique QR ticket code generated for every confirmed booking so I can present it at the theater for entry.*

### US-23: Booking Lookup by Ticket Code
> *As a theater staff member, I want to scan a QR code and look up the full booking details so I can validate entry quickly.*

---

## Issues Planned

| #  | Issue / Task                                                        | Assignee         | PR(s)    | Status |
| -- | ------------------------------------------------------------------- | ---------------- | -------- | ------ |
| 1  | JWT token generation in login & register responses                  | Paras Mittal     | #32      | Done   |
| 2  | JWT auth middleware (`JWTAuth` strict + `OptionalJWTAuth` flexible) | Omkar Salkade    | #32      | Done   |
| 3  | Protected route group for checkout/confirm & bookings               | Paras Mittal     | #32      | Done   |
| 4  | Bookings handler: JWT context + query-param dual mode               | Omkar Salkade    | #32      | Done   |
| 5  | Checkout confirm: override user_id from JWT context                 | Paras Mittal     | #32      | Done   |
| 6  | Bcrypt hashing in registration & login verification                 | Omkar Salkade    | #33      | Done   |
| 7  | Plaintext password migration on startup                             | Paras Mittal     | #33      | Done   |
| 8  | QR ticket model & auto-generation on booking confirmation           | Paresh Devlekar  | #34      | Done   |
| 9  | Ticket lookup API (`GET /bookings/by-ticket`)                       | Harsh Soni       | #34      | Done   |
| 10 | JWT middleware unit tests (6 tests)                                 | Omkar Salkade    | #35      | Done   |
| 11 | Bookings handler tests updated for JWT + query param (6 tests)      | Paras Mittal     | #35      | Done   |
| 12 | QR ticket handler tests (3 tests)                                   | Paresh Devlekar  | #35      | Done   |
| 13 | Password migration tests (4 tests)                                  | Harsh Soni       | #35      | Done   |
| 14 | Auth test updates for token assertions                              | Paras Mittal     | #35      | Done   |
| 15 | Frontend QR code display on booking confirmation                    | Harsh Soni       | #36      | Done   |
| 16 | Frontend booking history page                                       | Paresh Devlekar  | #36      | Done   |

---

## Detail of Work Completed

### 1. US-18 — JWT Authentication on Login & Register

**New file: `middleware/auth.go`**

- `GenerateToken(userID, secret)` — creates HS256-signed JWTs with 24-hour expiry, embedding `user_id` in claims
- `JWTAuth(secret)` — strict middleware: requires `Authorization: Bearer <token>`, aborts with `401` if missing/invalid/expired
- `OptionalJWTAuth(secret)` — flexible middleware: extracts `user_id` from token if present but does NOT block requests without a token; enables backward compatibility with the existing frontend

**Modified: `handlers/auth.go`**

- `Register` now returns `token` in the `201` response alongside `message` and `user`
- `Login` now returns `token` in the `200` response alongside `message` and `user`
- Package-level `JWTSecret` variable set from `config.JWTSecret` at startup in `main.go`

**Config: `config/config.go`**

- `JWT_SECRET` environment variable (default: `"change-me-in-production"`) already present; now actively used

### 2. US-19 — Protected Checkout Confirmation

**Modified: `routes/routes.go`**

- `POST /api/v1/checkout/confirm` now passes through `OptionalJWTAuth` middleware
- When a valid JWT is present, `ConfirmCheckout` overrides the request body's `user_id` with the authenticated user from the token
- Without a JWT, the handler falls back to the `user_id` in the request body — no frontend changes required

### 3. US-20 — Secure Booking History

**Modified: `handlers/bookings.go`**

- `GetUserBookings` now supports dual authentication:
  1. **JWT path:** extracts `user_id` from gin context (set by `OptionalJWTAuth`)
  2. **Query param path:** reads `?user_id=<id>` as fallback
  3. JWT takes priority if both are present
- `GET /api/v1/bookings` route uses `OptionalJWTAuth` middleware

### 4. US-21 — Bcrypt Password Hashing & Legacy Migration

**`handlers/auth.go`**

- Registration: `bcrypt.GenerateFromPassword` at default cost before storing
- Login: `bcrypt.CompareHashAndPassword` for credential verification
- Password complexity validation: minimum 8 chars, requires uppercase, lowercase, digit, and special character

**New: `database/password_migration.go`**

- `MigratePlaintextPasswords()` — runs at startup, scans all users, hashes any non-bcrypt passwords in place
- `isBcryptHash()` helper detects `$2a$`, `$2b$`, `$2y$` prefixes
- Idempotent: safe to run multiple times, leaves already-hashed passwords untouched

### 5. US-22 & US-23 — QR Ticket Generation & Lookup

**New model: `models/qr_ticket.go`**

- `QRTicket` struct: `BookingID` (unique index), `TicketCode` (unique index, 24-char URL-safe base64), `IsActive`, `ExpiresAt`
- Foreign key to `Booking` with `CASCADE` delete

**Modified: `handlers/checkout.go`**

- `getOrCreateQRTicket(tx, bookingID)` — generates a unique ticket code inside the booking transaction with up to 5 retry attempts on collision
- `ConfirmCheckout` response now includes `ticket_code` and `qr_value` (prefixed `"G4M:"`)

**New handler: `handlers/qr.go`**

- `GetBookingByTicketCode` — `GET /api/v1/bookings/by-ticket?ticket_code=<code>`
- Validates ticket is active and not expired
- Returns full booking details with movie, theater, screen, seats, and QR data

**Modified: `handlers/bookings.go`**

- Booking history now preloads `QRTicket` relation
- Each booking in the response includes `ticket_code` and `qr_value` if a QR ticket exists

---

## Frontend Tests (Jest + React Testing Library)

**13 test files, 108 test cases.** Run with `cd frontend && npm test`.

| # | Test File | Tests |
|---|-----------|-------|
| 1 | `src/lib/__tests__/api.test.ts` | 18 |
| 2 | `src/components/booking/__tests__/seat-selection.test.tsx` | 13 |
| 3 | `src/components/movies/__tests__/movie-detail.test.tsx` | 11 |
| 4 | `src/components/layout/__tests__/header.test.tsx` | 9 |
| 5 | `src/context/__tests__/auth-context.test.tsx` | 9 |
| 6 | `src/components/ui/__tests__/zip-code-modal.test.tsx` | 9 |
| 7 | `src/components/auth/__tests__/login-form.test.tsx` | 8 |
| 8 | `src/components/auth/__tests__/signup-form.test.tsx` | 8 |
| 9 | `src/hooks/__tests__/useZipCode.test.ts` | 6 |
| 10 | `src/components/ui/__tests__/search-bar.test.tsx` | 5 |
| 11 | `src/components/movies/__tests__/movie-grid.test.tsx` | 4 |
| 12 | `src/components/movies/__tests__/movie-card.test.tsx` | 4 |
| 13 | `src/context/__tests__/movies-context.test.tsx` | 4 |

**Frontend Total: 108 tests across 13 files**

---

## Backend Unit Tests (Go `testing` package)

**21 test files, 130 tests** across 7 packages. Run with `cd backend && go test ./... -v`.

### New in Sprint 4

#### `middleware/auth_test.go` (6 tests)

| # | Test Case | What It Verifies |
|---|-----------|-----------------|
| 1 | TestJWTAuth_ValidToken | Valid Bearer token → 200, `user_id` set in context |
| 2 | TestJWTAuth_NoHeader | Missing header → 401 |
| 3 | TestJWTAuth_MalformedHeader | Non-Bearer header → 401 |
| 4 | TestJWTAuth_ExpiredToken | Expired JWT → 401 |
| 5 | TestJWTAuth_WrongSecret | Token signed with wrong key → 401 |
| 6 | TestGenerateToken_RoundTrip | Generate + parse round-trip: correct `user_id` in claims |

#### `handlers/bookings_test.go` (6 tests — rewritten)

| # | Test Case | What It Verifies |
|---|-----------|-----------------|
| 1 | TestGetUserBookings_SuccessWithJWT | JWT Bearer token provides `user_id`, returns correct bookings |
| 2 | TestGetUserBookings_SuccessWithQueryParam | `?user_id=` query param works without JWT |
| 3 | TestGetUserBookings_Empty | No bookings for user → empty array, 200 |
| 4 | TestGetUserBookings_MissingUserID | No JWT + no query param → 400 |
| 5 | TestGetUserBookings_InvalidQueryParam | `?user_id=abc` → 400 |
| 6 | TestGetUserBookings_IsolationByUser | User A cannot see User B's bookings |

#### `handlers/qr_test.go` (3 tests — new)

| # | Test Case | What It Verifies |
|---|-----------|-----------------|
| 1 | TestGetBookingByTicketCode_Success | Valid ticket_code → 200, returns booking + QR data |
| 2 | TestGetBookingByTicketCode_MissingCode | No query param → 400 |
| 3 | TestGetBookingByTicketCode_NotFound | Unknown code → 404 |

#### `database/password_migration_test.go` (4 tests — new)

| # | Test Case | What It Verifies |
|---|-----------|-----------------|
| 1 | TestMigratePlaintextPasswords_ConvertsPlaintext | Plaintext → bcrypt hash, hash matches original |
| 2 | TestMigratePlaintextPasswords_LeavesBcryptUntouched | Already-hashed password unchanged |
| 3 | TestMigratePlaintextPasswords_MixedDatasetAndIdempotent | Mixed users + second run is idempotent |
| 4 | TestIsBcryptHash | Correctly detects `$2a$`, `$2b$`, `$2y$` prefixes |

#### `handlers/auth_test.go` (12 tests — expanded from 8)

| # | Test Case | What It Verifies |
|---|-----------|-----------------|
| 1 | TestRegister_Success | 201 + user created + bcrypt hash + **JWT token in response** |
| 2 | TestRegister_MissingFields | 400 on incomplete body |
| 3 | TestRegister_DuplicateEmail | 409 on existing email/username |
| 4 | TestRegister_EmailNormalized | Email lowercased and stored |
| 5 | TestRegister_UsernameValidation | Invalid username chars → 400 |
| 6 | TestRegister_WeakPasswordRejected | Missing complexity → 400 |
| 7 | TestRegister_TrimsAndNormalizesInput | Whitespace trimmed |
| 8 | TestLogin_Success | 200 + **JWT token in response** |
| 9 | TestLogin_WrongPassword | 401 |
| 10 | TestLogin_NonexistentUser | 401 |
| 11 | TestLogin_MissingFields | 400 |
| 12 | TestLogin_MalformedStoredHashReturnsUnauthorized | Corrupt hash → 401 (not 500) |

### Carried from Sprint 3 (unchanged, all passing)

| File | Tests |
|------|-------|
| `handlers/checkout_test.go` | 15 |
| `handlers/location_test.go` | 3 |
| `handlers/movie_test.go` | 8 |
| `handlers/showtime_test.go` | 15 |
| `workers/expiry_test.go` | 1 |
| `config/config_test.go` | 4 |
| `database/database_test.go` | 4 |
| `database/seed_test.go` | 9 |
| `middleware/cors_test.go` | 2 |
| `models/user_test.go` | 6 |
| `models/location_test.go` | 5 |
| `models/movie_test.go` | 6 |
| `models/theater_test.go` | 7 |
| `models/showtime_test.go` | 4 |
| `models/booking_test.go` | 8 |
| `routes/routes_test.go` | 2 |

**Backend Test Total: 130 tests across 21 files, all passing**

---

## Backend API Documentation

Base URL: `http://localhost:8080`

### Authentication

JWT tokens are returned in login/register responses. Include the token in protected endpoints:

```
Authorization: Bearer <token>
```

Tokens expire after **24 hours**. The `OptionalJWTAuth` middleware on protected routes allows both authenticated and unauthenticated access for backward compatibility.

---

### Endpoint Reference

| # | Method | Endpoint | Auth | Description |
|---|--------|----------|------|-------------|
| 1 | GET | `/api/v1/locations` | None | All supported locations |
| 2 | GET | `/api/v1/movies` | None | Active movies (`?zipcode=`, `?q=`) |
| 3 | GET | `/api/v1/movies/:id` | None | Single movie by ID |
| 4 | GET | `/api/v1/movies/:id/showtimes` | None | Showtimes by theater (`?zipcode=` required) |
| 5 | GET | `/api/v1/seats` | None | Seat availability (`?showtime_id=` required) |
| 6 | POST | `/api/v1/auth/register` | None | Create account → returns JWT token |
| 7 | POST | `/api/v1/auth/login` | None | Authenticate → returns JWT token |
| 8 | POST | `/api/v1/checkout/preview` | None | Price breakdown without booking |
| 9 | POST | `/api/v1/checkout/confirm` | Optional JWT | Create booking atomically → returns QR ticket |
| 10 | GET | `/api/v1/bookings` | Optional JWT | User's booking history (JWT or `?user_id=`) |
| 11 | GET | `/api/v1/bookings/by-ticket` | None | Lookup booking by QR ticket code |

---

### 6. `POST /api/v1/auth/register`

**Request:**
```json
{
  "email": "john@example.com",
  "username": "johndoe",
  "password": "Secret@123",
  "full_name": "John Doe"
}
```

**Password requirements:** 8–72 chars, must include uppercase, lowercase, digit, and special character.

**Response `201`:**
```json
{
  "message": "User registered successfully",
  "user": { "id": 2, "email": "john@example.com", "username": "johndoe", "full_name": "John Doe" },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Errors:** `400` validation, `409` duplicate email/username, `500` server error.

---

### 7. `POST /api/v1/auth/login`

**Request:**
```json
{
  "email": "tester@go4movies.com",
  "password": "Password"
}
```

**Response `200`:**
```json
{
  "message": "Login successful",
  "user": { "id": 1, "email": "tester@go4movies.com", "username": "tester", "full_name": "Test User" },
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Errors:** `400` missing fields, `401` invalid credentials.

---

### 9. `POST /api/v1/checkout/confirm`

**Request:** Same as preview. With JWT, `user_id` in body is overridden by the token's identity.

**Response `201`:**
```json
{
  "message": "Booking confirmed",
  "booking_id": 1,
  "booking_ref": "G4M-a1b2c3d4e5f6g7h8",
  "quote": { "..." },
  "payment_id": 1,
  "ticket_code": "dG9rZW5fY29kZV9oZXJl",
  "qr_value": "G4M:dG9rZW5fY29kZV9oZXJl"
}
```

---

### 10. `GET /api/v1/bookings`

**Auth:** `Authorization: Bearer <token>` OR `?user_id=<id>` (JWT takes priority if both present)

**Response `200`:**
```json
{
  "bookings": [
    {
      "id": 1,
      "booking_ref": "G4M-abc123",
      "status": "CONFIRMED",
      "total_amount": 110.16,
      "convenience_fee": 2.00,
      "tax_amount": 8.16,
      "payment_status": "PAID",
      "booked_at": "2026-04-28T12:00:00Z",
      "movie_title": "Inception",
      "movie_poster": "https://image.tmdb.org/...",
      "theater_name": "AMC Empire 25",
      "screen_name": "Screen 1",
      "screen_type": "IMAX",
      "show_date": "2026-05-01",
      "start_time": "19:30",
      "format": "IMAX",
      "language": "English",
      "ticket_code": "dG9rZW5fY29kZV9oZXJl",
      "qr_value": "G4M:dG9rZW5fY29kZV9oZXJl",
      "seats": [
        { "seat_id": 1, "row_label": "G", "col_number": 12, "seat_type": "Premium", "seat_price": 17.00 }
      ]
    }
  ]
}
```

---

### 11. `GET /api/v1/bookings/by-ticket?ticket_code=<code>` (New)

Looks up a booking by its QR ticket code. Intended for theater-side ticket validation.

**Response `200`:**
```json
{
  "booking": { "...same as bookings list item..." },
  "ticket_code": "dG9rZW5fY29kZV9oZXJl",
  "qr_value": "G4M:dG9rZW5fY29kZV9oZXJl"
}
```

**Errors:** `400` missing `ticket_code`, `404` not found or expired.

---

## Summary

Sprint 4 delivered JWT authentication, bcrypt password security, and QR ticket generation — completing the backend security and ticket management features. Login and register now return signed JWT tokens (HS256, 24h expiry). The `OptionalJWTAuth` middleware enables a smooth transition: protected routes accept JWT tokens when available but fall back to existing query-param/body authentication, ensuring zero frontend breakage. All passwords are now bcrypt-hashed, with an automatic startup migration for any legacy plaintext passwords. Every confirmed booking generates a unique QR ticket code that can be looked up via the new `/bookings/by-ticket` endpoint for theater-side validation. The backend test suite grew from 104 to 130 tests across 21 files, adding 6 JWT middleware tests, 6 bookings dual-auth tests, 3 QR ticket tests, 4 password migration tests, and 4 expanded auth tests. All 130 backend tests and 108 frontend tests pass.
