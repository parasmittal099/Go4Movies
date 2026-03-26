import React from 'react'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { MovieDetail } from '../movie-detail'
import type { Movie, ShowtimesResponse } from '@/lib/types'

const movie: Movie = { id: 1, title: 'Dune', description: 'Sci-fi epic', poster_url: '/dune.jpg' }

const showtimes: ShowtimesResponse = {
  movie_id: 1,
  title: 'Dune',
  dates: ['2026-04-10', '2026-04-11'],
  theaters: [
    {
      theater_id: 10,
      name: 'IMAX Cinema',
      address: '123 Main St',
      showtimes: [
        {
          id: 101,
          show_date: '2026-04-10',
          start_time: '10:00',
          end_time: '13:00',
          language: 'English',
          format: 'IMAX',
          price_multiplier: 1.5,
          screen_name: 'Screen 1',
          screen_type: 'IMAX',
        },
        {
          id: 102,
          show_date: '2026-04-10',
          start_time: '14:00',
          end_time: '17:00',
          language: 'English',
          format: 'IMAX',
          price_multiplier: 1.5,
          screen_name: 'Screen 1',
          screen_type: 'IMAX',
        },
      ],
    },
  ],
}

const defaultProps = {
  movie,
  showtimes,
  showtimesLoading: false,
  onBookTickets: jest.fn(),
}

describe('MovieDetail', () => {
  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('renders the movie title', () => {
    render(<MovieDetail {...defaultProps} />)
    expect(screen.getByRole('heading', { level: 1 })).toHaveTextContent('Dune')
  })

  it('renders movie description', () => {
    render(<MovieDetail {...defaultProps} />)
    expect(screen.getByText('Sci-fi epic')).toBeInTheDocument()
  })

  it('renders date buttons for provided dates', () => {
    render(<MovieDetail {...defaultProps} />)
    // Day numbers for April 10 and April 11 — use data from the date buttons
    expect(screen.getByText('10')).toBeInTheDocument()
    expect(screen.getByText('11')).toBeInTheDocument()
  })

  it('renders theater name', () => {
    render(<MovieDetail {...defaultProps} />)
    expect(screen.getByText('IMAX Cinema')).toBeInTheDocument()
  })

  it('renders showtime buttons', () => {
    render(<MovieDetail {...defaultProps} />)
    expect(screen.getByText('10:00')).toBeInTheDocument()
    expect(screen.getByText('14:00')).toBeInTheDocument()
  })

  it('shows loading skeletons when showtimesLoading is true', () => {
    render(<MovieDetail {...defaultProps} showtimesLoading={true} />)
    expect(screen.queryByText('IMAX Cinema')).not.toBeInTheDocument()
  })

  it('shows no showtimes message when cinemas is empty', () => {
    render(<MovieDetail {...defaultProps} showtimes={{ ...showtimes, theaters: [] }} />)
    expect(screen.getByText(/no showtimes available/i)).toBeInTheDocument()
  })

  it('selecting a timing shows the sticky booking bar', async () => {
    render(<MovieDetail {...defaultProps} />)
    expect(screen.queryByText(/select seats/i)).not.toBeInTheDocument()
    fireEvent.click(screen.getByText('10:00'))
    await waitFor(() => expect(screen.getByText(/select seats/i)).toBeInTheDocument())
  })

  it('calls onBookTickets when Select Seats is clicked', async () => {
    render(<MovieDetail {...defaultProps} />)
    fireEvent.click(screen.getByText('10:00'))
    await waitFor(() => fireEvent.click(screen.getByText(/select seats/i)))
    expect(defaultProps.onBookTickets).toHaveBeenCalledWith(101)
  })

  it('changes selected date when another date is clicked', () => {
    render(<MovieDetail {...defaultProps} />)
    // Click the second date (April 11 = day 11)
    fireEvent.click(screen.getByText('11'))
    // No showtimes on Apr 11, so theater should disappear
    expect(screen.queryByText('IMAX Cinema')).not.toBeInTheDocument()
  })

  it('renders null showtimes correctly (loading state initial)', () => {
    render(<MovieDetail {...defaultProps} showtimes={null} />)
    // No theater should be visible
    expect(screen.queryByText('IMAX Cinema')).not.toBeInTheDocument()
  })
})
