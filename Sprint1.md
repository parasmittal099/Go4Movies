# Sprint 1 Report — Go4Movies

**Sprint Duration:** Feb 4, 2026 – Feb 19, 2026
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

### US-01: Browse Movies by Location
> *As a user, I want to browse movies playing near my location so I can decide what to watch.*

### US-02: Search Movies
> *As a user, I want to search movies by title or genre so I can quickly find what I'm looking for.*

### US-03: View Movie Details
> *As a user, I want to see full movie details (poster, cast, rating, duration) so I can make an informed choice.*

### US-04: Real Movie Data with Posters
> *As a user, I want to see real movie posters and accurate information so I can recognize the movies.*

### US-05: User Authentication (Register, Login, Signup Page)
> *As a user, I want to register and log in so I can access personalized features.*

---

## Issues Planned

| #  | Issue / Task                                          | Assignee         | PR   | Status  |
| -- | ----------------------------------------------------- | ---------------- | ---- | ------- |
| 1  | Initialize Next.js frontend with Docker               | Harsh Soni       | #2   | Done    |
| 2  | Initialize Go backend with Docker                     | Paras Mittal     | #3   | Done    |
| 3  | Login page UI                                         | Harsh Soni       | #4   | Done    |
| 4  | Backend scaffold — models, seed data, APIs            | Paras Mittal     | #6   | Done    |
| 5  | Auth backend — Register & Login APIs                  | Omkar Salkade    | #7   | Done    |
| 6  | Frontend — movie listing, detail, zip modal, API      | Paresh Devlekar  | #8   | Done    |
| 7  | Home page UI — components & styling                   | Harsh Soni       | #9   | Done    |
| 8  | Seed real movie data with TMDB poster URLs            | Paras Mittal     | #10  | Done    |
| 9  | Signup page + auth API integration                    | Paresh Devlekar  | #11  | Done    |
| 10 | Global auth & location context with persistence       | Harsh Soni       | #12  | Done    |
| 11 | Global movie context and search functionality         | Harsh Soni       | #13  | Done    |
| 12 | Sample user seed + email case insensitivity           | Omkar Salkade    | #14  | Done    |

---

## Successfully Completed

### 1. Project Infrastructure & DevOps
- Next.js 16 frontend with Tailwind CSS v4 and Docker (PR #2)
- Go/Gin REST API with SQLite via GORM and Docker (PR #3)
- Docker Compose orchestrating both services with development volumes
- CORS middleware for frontend–backend communication

### 2. US-01 — Browse Movies by Location
**Backend (PR #6):**
- `GET /api/v1/locations` — returns all supported zip codes / cities
- `GET /api/v1/movies?zipcode=` — movies filtered by city (resolves zip → city → theaters → showtimes → movies)
- Seed data: 7 Florida locations, 7 theaters, 7 screens, 10 movies, 25+ showtimes

**Frontend (PR #8, #9, #12):**
- Zip code selection modal with browser geolocation detection (OpenStreetMap Nominatim)
- Zip code persistence via localStorage and global location context
- Movie listing page with responsive grid layout
- Change location button on movies page

### 3. US-02 — Search Movies
**Backend (PR #6):** `GET /api/v1/movies?q=` with LIKE search on title and genre
**Frontend (PR #13):** Global movie context with integrated search bar, real-time filtering

### 4. US-03 — View Movie Details
**Backend (PR #6):** `GET /api/v1/movies/:id` returns full movie details
**Frontend (PR #8, #9):** Movie detail page with poster, description, cast, genre, rating, duration, and "Book Tickets" button

### 5. US-04 — Real Movie Data with Posters
**PR #10:** Replaced placeholder data with 10 real movies sourced from the jsonfakery movies API. All movies have working TMDB poster URLs, complete metadata (genre, cast, MPAA rating, duration, release date, trailer URL), and showtimes distributed across all theaters.

### 6. US-05 — User Authentication
**Backend (PR #7, #14):**
- `POST /api/v1/auth/register` — creates user with email, username, password, full name
- `POST /api/v1/auth/login` — validates email + password against DB
- Sample user seeded for testing; email lookup is case-insensitive

**Frontend (PR #4, #11, #12):**
- Login page with styled email/password form
- Signup page with full name, email, password, and confirm password fields
- Auth API integration — login and register forms call backend endpoints
- Global auth context with persistent state management

---

## What Didn't Get Completed & Why

### 1. JWT Token Generation & Session Management
| Status | Not implemented |
| ------ | --- |
| **Why** | Sprint 1 prioritized getting the auth endpoints and frontend integration functional. The auth middleware file (`middleware/auth.go`) exists but is empty. Login and register APIs return success responses but do not issue JWT tokens. |
| **Impact** | No token-based session persistence. Users can register and log in, but there are no protected routes. Planned for Sprint 2. |

### 2. Password Hashing (bcrypt)
| Status | Not implemented |
| ------ | --- |
| **Why** | Deprioritized to deliver the core auth flow and movie browsing features within the sprint timeline. Passwords are currently stored in plain text. |
| **Impact** | Security gap. Will be resolved in Sprint 2 using bcrypt. |

### 3. Token Refresh Endpoint (`POST /auth/refresh`)
| Status | Not implemented |
| ------ | --- |
| **Why** | Depends on JWT implementation which was deferred. The endpoint was planned in the original sprint backlog but could not be built without the token infrastructure in place. |
| **Impact** | No token refresh capability. Will be built alongside JWT in Sprint 2. |

### 4. Booking Flow (Seat Selection, Showtimes, Payment)
| Status | Not started — planned for Sprint 2 |
| ------ | --- |
| **Why** | Sprint 1 scope was intentionally limited to infrastructure, movie browsing, and authentication foundations. The booking system is the primary Sprint 2 deliverable. |
| **Impact** | Database models for Booking, BookingSeat, Seat, and Payment are already defined and auto-migrated. No API endpoints or frontend pages exist yet. |

---

## Summary

Sprint 1 delivered a working full-stack movie browsing platform with location-based filtering, search, real TMDB movie posters, and user authentication (register + login + signup). All 12 planned issues were completed and merged across 13 pull requests. The outstanding items — JWT session management, bcrypt password hashing, refresh tokens, and the booking flow — are well-scoped carry-overs for Sprint 2, where the focus will shift to seat selection and ticket booking.
