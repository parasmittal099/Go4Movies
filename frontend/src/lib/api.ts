// API client for backend communication
// Centralizes all fetch calls so endpoints and error handling live in one place.

import {
  AuthSuccessResponse,
  CheckoutConfirmResponse,
  CheckoutQuote,
  Location,
  LoginPayload,
  Movie,
  RegisterPayload,
  SeatsResponse,
  ShowtimesResponse,
} from "./types"

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080"

async function parseResponse<T>(response: Response, fallbackMessage: string): Promise<T> {
  const data = await response.json().catch(() => null)

  if (!response.ok) {
    const message =
      data && typeof data === "object" && "error" in data && typeof data.error === "string"
        ? data.error
        : fallbackMessage
    throw new Error(message)
  }

  return data as T
}

/**
 * Fetch all locations (zip codes) from the backend.
 */
export async function fetchLocations(): Promise<Location[]> {
  const response = await fetch(`${API_BASE_URL}/api/v1/locations`)
  if (!response.ok) {
    throw new Error("Failed to fetch locations")
  }
  const data = await response.json()
  return data.locations
}

/**
 * Fetch movies for a given zip code.
 */
export async function fetchMoviesByZipCode(zipCode: string): Promise<Movie[]> {
  const response = await fetch(
    `${API_BASE_URL}/api/v1/movies?zipcode=${zipCode}`
  )
  if (!response.ok) {
    throw new Error("Failed to fetch movies")
  }
  return response.json()
}

/**
 * Fetch a single movie by ID.
 */
export async function fetchMovieById(movieId: string): Promise<Movie> {
  const response = await fetch(`${API_BASE_URL}/api/v1/movies/${movieId}`)
  if (!response.ok) {
    throw new Error("Movie not found")
  }
  return response.json()
}

export async function registerUser(payload: RegisterPayload): Promise<AuthSuccessResponse> {
  const response = await fetch(`${API_BASE_URL}/api/v1/auth/register`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify(payload),
  })

  return parseResponse<AuthSuccessResponse>(response, "Failed to register user")
}

export async function loginUser(payload: LoginPayload): Promise<AuthSuccessResponse> {
  const response = await fetch(`${API_BASE_URL}/api/v1/auth/login`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    credentials: "include",
    body: JSON.stringify(payload),
  })

  return parseResponse<AuthSuccessResponse>(response, "Failed to login")
}

/**
 * Fetch showtimes for a movie+zipcode combination.
 */
export async function fetchMovieShowtimes(
  movieId: string | number,
  zipCode: string
): Promise<ShowtimesResponse> {
  const response = await fetch(
    `${API_BASE_URL}/api/v1/movies/${movieId}/showtimes?zipcode=${zipCode}`
  )
  return parseResponse<ShowtimesResponse>(response, "Failed to fetch showtimes")
}

export async function fetchShowtimeSeats(showtimeId: number): Promise<SeatsResponse> {
  const response = await fetch(`${API_BASE_URL}/api/v1/seats?showtime_id=${showtimeId}`)
  return parseResponse<SeatsResponse>(response, "Failed to fetch seats")
}

export interface CheckoutPayload {
  user_id: number
  showtime_id: number
  seat_ids: number[]
  discount_code?: string
}

export async function previewCheckout(payload: CheckoutPayload): Promise<CheckoutQuote> {
  const response = await fetch(`${API_BASE_URL}/api/v1/checkout/preview`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  })
  return parseResponse<CheckoutQuote>(response, "Failed to get checkout preview")
}

export async function confirmCheckout(payload: CheckoutPayload): Promise<CheckoutConfirmResponse> {
  const response = await fetch(`${API_BASE_URL}/api/v1/checkout/confirm`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  })
  return parseResponse<CheckoutConfirmResponse>(response, "Failed to confirm booking")
}

