# Go4Movies
Project Name - Go4Movies 

Go4Movies is a web-based movie ticket booking platform designed to simulate a real-world, high-traffic e-commerce environment. The application allows users to browse movies, check showtimes, and book specific seats in an interactive theater map.

Core Problem Solved: The primary technical challenge this project solves is the "Race Condition" inherent in booking systems—where multiple users attempt to book the exact same seat at the same millisecond. GoMovies prevents double-booking errors through a robust concurrency control strategy.

Technical Architecture:

Frontend (React + Tailwind): A responsive Single Page Application (SPA) utilizing functional components and Hooks. It renders a dynamic visual seat map and manages complex local state, including reservation countdown timers and immediate visual feedback (e.g., blocking a seat the moment it is taken).

Backend (Go/Golang): A high-performance REST API built with Go (using Chi or Gin). Go would be beneficial for its efficient concurrency model (Goroutines) to handle multiple simultaneous booking requests without blocking.

Database (SQLite): A lightweight relational database. We will  allow concurrent reads and writes, simulating a production-grade database environment.

Concurrency Control: The system implements Optimistic Locking via SQL transactions.

Reservation: When a user clicks a seat, the backend atomically updates the status to RESERVED with a 5-minute expiry.

Conflict Resolution: If User B tries to book the seat while User A holds the lock, the transaction fails, and the API returns a 409 Conflict error.

Cleanup: A background Goroutine runs periodically to expire unpaid reservations, releasing seats back to the pool.

DevOps (Docker): The application is fully containerized. A single docker-compose.yml orchestrates the React frontend (served via Nginx or Vite preview) and the Go binary, ensuring the environment is identical for all developers.

# Frontend Testing

The frontend uses **Jest** and **React Testing Library** for unit testing, covering all major components, hooks, contexts, and API utilities.

**Coverage:** 106 tests across 13 test suites. All business-logic files (components, hooks, contexts, API client) achieve **80%+ coverage**.

**What's tested:**
- `lib/api.ts` — all fetch functions (success + error paths)
- `hooks/useZipCode.ts` — localStorage read/write/clear/redirect
- `context/auth-context.tsx` / `movies-context.tsx` — state, localStorage rehydration
- `components/auth/` — login + signup form validation, API calls, toggle states
- `components/movies/` — MovieCard, MovieGrid, MovieDetail (date selection, booking bar)
- `components/booking/SeatSelection` — seat loading, selection, pricing, error states
- `components/layout/Header` — auth states, search open/close/escape
- `components/ui/SearchBar`, `ZipCodeModal` — input, filtering, location selection

**Run tests:**
```bash
cd frontend
npm test                # Run all tests
npm run test:watch      # Watch mode
npm run test:coverage   # With coverage report
```

# Members 
Harsh Soni - FrontEnd \
Paresh Devlekar - FrontEnd \
Omkar Salkade - BackEnd \
Paras Mittal - BackEnd 
