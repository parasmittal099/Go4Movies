import React from 'react'
import { render, screen, fireEvent, waitFor, act } from '@testing-library/react'
import { ZipCodeModal } from '../zip-code-modal'
import * as api from '@/lib/api'

jest.mock('next/navigation', () => ({
  useRouter: () => ({ push: jest.fn() }),
}))

jest.mock('@/lib/api', () => ({
  fetchLocations: jest.fn(),
}))

const locations = [
  { id: 1, zipcode: '10001', city: 'New York', state: 'NY', country: 'US' },
  { id: 2, zipcode: '90210', city: 'Beverly Hills', state: 'CA', country: 'US' },
]

const defaultProps = {
  onSelectLocation: jest.fn(),
  onClose: jest.fn(),
}

beforeEach(() => {
  jest.clearAllMocks()
  localStorage.clear()
  ;(api.fetchLocations as jest.Mock).mockResolvedValue(locations)
})

describe('ZipCodeModal', () => {
  it('renders modal heading', async () => {
    render(<ZipCodeModal {...defaultProps} />)
    expect(screen.getByText(/hi there/i)).toBeInTheDocument()
  })

  it('shows search input', async () => {
    render(<ZipCodeModal {...defaultProps} />)
    expect(screen.getByPlaceholderText(/search for your city or zip code/i)).toBeInTheDocument()
  })

  it('does NOT render close button when no zip code is stored', async () => {
    render(<ZipCodeModal {...defaultProps} />)
    await waitFor(() => expect(api.fetchLocations).toHaveBeenCalled())
    // The close button is only rendered when zipCode is truthy
    expect(screen.queryByRole('button', { name: /close/i })).not.toBeInTheDocument()
  })

  it('renders close button when a zip code is stored', async () => {
    localStorage.setItem('selectedZipCode', '10001')
    render(<ZipCodeModal {...defaultProps} />)
    await waitFor(() => expect(api.fetchLocations).toHaveBeenCalled())
    expect(screen.getByRole('button', { name: /close/i })).toBeInTheDocument()
  })

  it('calls onClose when close button is clicked', async () => {
    localStorage.setItem('selectedZipCode', '10001')
    render(<ZipCodeModal {...defaultProps} />)
    await waitFor(() => expect(api.fetchLocations).toHaveBeenCalled())
    fireEvent.click(screen.getByRole('button', { name: /close/i }))
    expect(defaultProps.onClose).toHaveBeenCalledTimes(1)
  })

  it('filters locations when user types in search', async () => {
    render(<ZipCodeModal {...defaultProps} />)
    await waitFor(() => expect(api.fetchLocations).toHaveBeenCalled())
    const input = screen.getByPlaceholderText(/search for your city or zip code/i)
    await act(async () => {
      fireEvent.change(input, { target: { value: 'New York' } })
    })
    await waitFor(() => expect(screen.getByText('New York')).toBeInTheDocument())
    expect(screen.queryByText('Beverly Hills')).not.toBeInTheDocument()
  })

  it('filters by zip code prefix', async () => {
    render(<ZipCodeModal {...defaultProps} />)
    await waitFor(() => expect(api.fetchLocations).toHaveBeenCalled())
    const input = screen.getByPlaceholderText(/search for your city or zip code/i)
    await act(async () => {
      fireEvent.change(input, { target: { value: '902' } })
    })
    await waitFor(() => expect(screen.getByText('Beverly Hills')).toBeInTheDocument())
    expect(screen.queryByText('New York')).not.toBeInTheDocument()
  })

  it('calls onSelectLocation when a result is clicked', async () => {
    render(<ZipCodeModal {...defaultProps} />)
    await waitFor(() => expect(api.fetchLocations).toHaveBeenCalled())
    const input = screen.getByPlaceholderText(/search for your city or zip code/i)
    await act(async () => {
      fireEvent.change(input, { target: { value: 'New' } })
    })
    await waitFor(() => fireEvent.click(screen.getByText('New York')))
    expect(defaultProps.onSelectLocation).toHaveBeenCalledWith(locations[0])
  })

  it('handles API fetch error gracefully', async () => {
    ;(api.fetchLocations as jest.Mock).mockRejectedValue(new Error('Network error'))
    // Should render without crashing
    render(<ZipCodeModal {...defaultProps} />)
    await waitFor(() => expect(api.fetchLocations).toHaveBeenCalled())
    // Modal is still visible
    expect(screen.getByText(/hi there/i)).toBeInTheDocument()
  })
})
