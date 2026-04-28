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

export interface SeatAPI {
  id: number
  row_label: string
  col_number: number
  seat_type: string
  price: number
  status: "AVAILABLE" | "RESERVED" | "BOOKED"
}

export interface SeatsResponse {
  showtime: {
    id: number
    movie_title: string
    screen_name: string
    screen_type: string
    theater_name: string
    show_date: string
    start_time: string
    format: string
    language: string
  }
  layout: {
    total_rows: number
    total_cols: number
  }
  summary: {
    total: number
    available: number
    reserved: number
    booked: number
  }
  seats: SeatAPI[]
}

export interface CheckoutLineItem {
  seat_id: number
  row_label: string
  col_number: number
  seat_type: string
  unit_price: number
}

export interface CheckoutTotals {
  subtotal: number
  convenience_fee: number
  tax_amount: number
  discount_code?: string
  discount_amount: number
  total_due: number
}

export interface CheckoutQuote {
  showtime_id: number
  user_id: number
  line_items: CheckoutLineItem[]
  totals: CheckoutTotals
}

export interface CheckoutConfirmResponse {
  message: string
  booking_id: number
  booking_ref: string
  quote: CheckoutQuote
  payment_id: number
}

export interface BookingSeatDetail {
  seat_id: number
  row_label: string
  col_number: number
  seat_type: string
  seat_price: number
}

export interface BookingDetail {
  id: number
  booking_ref: string
  status: string
  total_amount: number
  convenience_fee: number
  tax_amount: number
  payment_status: string
  booked_at: string
  movie_title: string
  movie_poster: string
  theater_name: string
  screen_name: string
  screen_type: string
  show_date: string
  start_time: string
  format: string
  language: string
  seats: BookingSeatDetail[]
  ticket_code?: string
  qr_value?: string
}

