// Shared type definitions for the application

export interface Movie {
  id: number
  title: string
  description: string
  poster_url: string
}

export interface Location {
  id: number
  zipcode: string
  city: string
  state: string
  country: string
}

export interface User {
  id: number
  email: string
  username: string
  full_name: string
  phone: string | null
  created_at: string
  updated_at: string
}

export interface RegisterPayload {
  email: string
  username: string
  password: string
  full_name: string
}

export interface LoginPayload {
  email: string
  password: string
}

export interface AuthSuccessResponse {
  message: string
  user: User
}

export interface ShowtimeEntry {
  id: number
  show_date: string
  start_time: string
  end_time: string
  language: string
  format: string
  price_multiplier: number
  screen_name: string
  screen_type: string
}

export interface TheaterGroup {
  theater_id: number
  name: string
  address?: string
  showtimes: ShowtimeEntry[]
}

export interface ShowtimesResponse {
  movie_id: number
  title: string
  dates: string[]
  theaters: TheaterGroup[]
}

