import React from 'react'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { SeatSelection } from '../seat-selection'
import * as api from '@/lib/api'
import type { SeatsResponse } from '@/lib/types'

jest.mock('@/lib/api', () => ({
  fetchShowtimeSeats: jest.fn(),
}))

const mockMovie = { title: 'Dune', time: '10:00', format: 'IMAX', hall: 'Screen 1' }

const mockSeatsResponse: SeatsResponse = {
  showtime: {
    id: 1,
    movie_title: 'Dune',
    screen_name: 'Screen 1',
    screen_type: 'IMAX',
    theater_name: 'IMAX Cinema',
    show_date: '2026-04-01',
    start_time: '10:00',
    format: 'IMAX',
    language: 'English',
  },
  layout: { total_rows: 2, total_cols: 3 },
  summary: { total: 6, available: 5, reserved: 0, booked: 1 },
  seats: [
    { id: 1, row_label: 'A', col_number: 1, seat_type: 'standard', price: 15, status: 'AVAILABLE' },
    { id: 2, row_label: 'A', col_number: 2, seat_type: 'premium', price: 18, status: 'AVAILABLE' },
    { id: 3, row_label: 'A', col_number: 3, seat_type: 'vip', price: 25, status: 'AVAILABLE' },
    { id: 4, row_label: 'B', col_number: 1, seat_type: 'standard', price: 15, status: 'BOOKED' },
  ],
}

const defaultProps = {
  showtimeId: 42,
  movie: mockMovie,
  onProceed: jest.fn(),
  onBack: jest.fn(),
  onChangeMovie: jest.fn(),
}

beforeEach(() => {
  jest.clearAllMocks()
  ;(api.fetchShowtimeSeats as jest.Mock).mockResolvedValue(mockSeatsResponse)
})

describe('SeatSelection', () => {
  it('shows loading state initially', () => {
    ;(api.fetchShowtimeSeats as jest.Mock).mockReturnValue(new Promise(() => {}))
    render(<SeatSelection {...defaultProps} />)
    expect(screen.getByText(/loading seats/i)).toBeInTheDocument()
  })

  it('renders seat grid after loading', async () => {
    render(<SeatSelection {...defaultProps} />)
    await waitFor(() => expect(screen.getByLabelText(/A1 available/i)).toBeInTheDocument())
  })

  it('renders booked seat as disabled', async () => {
    render(<SeatSelection {...defaultProps} />)
    await waitFor(() => {
      const bookedSeat = screen.getByLabelText(/B1 sold/i)
      expect(bookedSeat).toBeDisabled()
    })
  })

  it('selects an available seat on click', async () => {
    render(<SeatSelection {...defaultProps} />)
    await waitFor(() => expect(screen.getByLabelText(/A1 available/i)).toBeInTheDocument())
    fireEvent.click(screen.getByLabelText(/A1 available/i))
    // After selection, the seat button shows the seat ID text internally
    // The bottom bar should show "A1" in the selected seats list
    await waitFor(() => {
      const bottomBar = document.querySelector('[class*="pointer-events-auto"]')
      expect(bottomBar?.textContent).toContain('A1')
    })
  })

  it('updates total price when seats are selected', async () => {
    render(<SeatSelection {...defaultProps} />)
    await waitFor(() => expect(screen.getByLabelText(/A1 available/i)).toBeInTheDocument())
    fireEvent.click(screen.getByLabelText(/A1 available/i))
    // Price of A1 is $15
    await waitFor(() => expect(screen.getByText('$15.00')).toBeInTheDocument())
  })

  it('calls onBack when back arrow button is clicked', async () => {
    render(<SeatSelection {...defaultProps} />)
    await waitFor(() => expect(api.fetchShowtimeSeats).toHaveBeenCalled())
    // The back button contains material icon "arrow_back"
    const backBtn = screen.getByText('arrow_back').closest('button')!
    fireEvent.click(backBtn)
    expect(defaultProps.onBack).toHaveBeenCalled()
  })

  it('calls onChangeMovie when Change Movie is clicked', async () => {
    render(<SeatSelection {...defaultProps} />)
    await waitFor(() => expect(api.fetchShowtimeSeats).toHaveBeenCalled())
    fireEvent.click(screen.getByText('Change Movie'))
    expect(defaultProps.onChangeMovie).toHaveBeenCalled()
  })

  it('calls onProceed with selected seats when Proceed button is clicked', async () => {
    render(<SeatSelection {...defaultProps} />)
    await waitFor(() => expect(screen.getByLabelText(/A1 available/i)).toBeInTheDocument())
    fireEvent.click(screen.getByLabelText(/A1 available/i))
    await waitFor(() => {
      const proceedBtn = screen.getByText(/proceed to payment/i).closest('button')!
      expect(proceedBtn).not.toBeDisabled()
      fireEvent.click(proceedBtn)
    })
    expect(defaultProps.onProceed).toHaveBeenCalledWith(
      expect.arrayContaining([expect.objectContaining({ id: 'A1' })])
    )
  })

  it('Proceed button is disabled when no seats are selected', async () => {
    render(<SeatSelection {...defaultProps} />)
    await waitFor(() => expect(api.fetchShowtimeSeats).toHaveBeenCalled())
    expect(screen.getByText(/proceed to payment/i).closest('button')).toBeDisabled()
  })

  it('shows error message on API failure', async () => {
    ;(api.fetchShowtimeSeats as jest.Mock).mockRejectedValue(new Error('Network error'))
    render(<SeatSelection {...defaultProps} />)
    await waitFor(() => expect(screen.getByText(/network error/i)).toBeInTheDocument())
  })

  it('shows error for invalid showtimeId (0)', async () => {
    ;(api.fetchShowtimeSeats as jest.Mock).mockClear()
    render(<SeatSelection {...defaultProps} showtimeId={0} />)
    await waitFor(() =>
      expect(screen.getByText(/invalid showtime selected/i)).toBeInTheDocument()
    )
    expect(api.fetchShowtimeSeats).not.toHaveBeenCalled()
  })

  it('shows no seats message when seats array is empty', async () => {
    ;(api.fetchShowtimeSeats as jest.Mock).mockResolvedValue({ ...mockSeatsResponse, seats: [] })
    render(<SeatSelection {...defaultProps} />)
    await waitFor(() => expect(screen.getByText(/no seats available/i)).toBeInTheDocument())
  })

  it('renders movie title in header', async () => {
    render(<SeatSelection {...defaultProps} />)
    await waitFor(() => expect(api.fetchShowtimeSeats).toHaveBeenCalled())
    expect(screen.getByRole('heading', { name: /dune/i })).toBeInTheDocument()
  })
})
