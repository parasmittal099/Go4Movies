import React from 'react'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import PaymentPage from '../page'
import * as api from '@/lib/api'

// ── Navigation mocks ──────────────────────────────────────────────────────────
const mockRouterBack = jest.fn()
const mockRouterPush = jest.fn()

jest.mock('next/navigation', () => ({
  useRouter: () => ({ back: mockRouterBack, push: mockRouterPush }),
  useSearchParams: () => ({
    get: (key: string) => {
      const params: Record<string, string> = {
        showtimeId: '1',
        seats: 'A1,A2',
        seatIds: '10,11',
      }
      return params[key] ?? null
    },
  }),
}))

// ── Auth mock ─────────────────────────────────────────────────────────────────
jest.mock('@/context/auth-context', () => ({
  useAuth: () => ({
    user: { id: 42, email: 'test@example.com', username: 'tester' },
  }),
}))

// ── API mocks ─────────────────────────────────────────────────────────────────
jest.mock('@/lib/api', () => ({
  previewCheckout: jest.fn().mockResolvedValue({
    line_items: [],
    totals: {
      subtotal: 20,
      convenience_fee: 2,
      tax_amount: 1.76,
      discount_amount: 0,
      discount_code: '',
      total_due: 23.76,
    },
  }),
  confirmCheckout: jest.fn(),
  fetchShowtimeSeats: jest.fn().mockResolvedValue({
    showtime: {
      movie_title: 'Test Movie',
      format: 'IMAX',
      language: 'English',
      theater_name: 'CineMax',
      screen_name: 'Screen 1',
      show_date: '2026-04-28',
      start_time: '18:00',
    },
    seats: [
      { id: 10, price: 10 },
      { id: 11, price: 10 },
    ],
  }),
}))

// ── Helpers ───────────────────────────────────────────────────────────────────
function renderPage() {
  return render(<PaymentPage />)
}

function getCardNumberInput() {
  return screen.getByPlaceholderText('0000-0000-0000-0000')
}

function getExpiryInput() {
  return screen.getByPlaceholderText('MM/YY')
}

// ─────────────────────────────────────────────────────────────────────────────
describe('PaymentPage — card number formatting', () => {
  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('renders the card number input with correct placeholder', () => {
    renderPage()
    expect(getCardNumberInput()).toBeInTheDocument()
  })

  it('enforces maxLength of 19 on card number input', () => {
    renderPage()
    expect(getCardNumberInput()).toHaveAttribute('maxLength', '19')
  })

  it('inserts a hyphen after every 4 digits while typing', () => {
    renderPage()
    const input = getCardNumberInput()

    fireEvent.change(input, { target: { value: '12345678' } })
    expect(input).toHaveValue('1234-5678')
  })

  it('formats a full 16-digit card number as XXXX-XXXX-XXXX-XXXX', () => {
    renderPage()
    const input = getCardNumberInput()

    fireEvent.change(input, { target: { value: '1234567890123456' } })
    expect(input).toHaveValue('1234-5678-9012-3456')
  })

  it('strips non-digit characters entered by the user', () => {
    renderPage()
    const input = getCardNumberInput()

    fireEvent.change(input, { target: { value: 'abcd1234efgh5678' } })
    expect(input).toHaveValue('1234-5678')
  })

  it('does not add a trailing hyphen after the last 4-digit group', () => {
    renderPage()
    const input = getCardNumberInput()

    fireEvent.change(input, { target: { value: '1234567890123456' } })
    expect(input).toHaveValue('1234-5678-9012-3456')
    expect((input as HTMLInputElement).value.endsWith('-')).toBe(false)
  })

  it('truncates input beyond 16 digits', () => {
    renderPage()
    const input = getCardNumberInput()

    // 20 digits — only first 16 should survive
    fireEvent.change(input, { target: { value: '12345678901234567890' } })
    expect(input).toHaveValue('1234-5678-9012-3456')
  })

  it('shows only digits (no hyphens) for fewer than 4 digits', () => {
    renderPage()
    const input = getCardNumberInput()

    fireEvent.change(input, { target: { value: '123' } })
    expect(input).toHaveValue('123')
  })

  it('adds first hyphen exactly after the 4th digit', () => {
    renderPage()
    const input = getCardNumberInput()

    fireEvent.change(input, { target: { value: '1234' } })
    expect(input).toHaveValue('1234')

    fireEvent.change(input, { target: { value: '12345' } })
    expect(input).toHaveValue('1234-5')
  })
})

// ─────────────────────────────────────────────────────────────────────────────
describe('PaymentPage — expiry date formatting', () => {
  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('renders the expiry input with correct placeholder', () => {
    renderPage()
    expect(getExpiryInput()).toBeInTheDocument()
  })

  it('enforces maxLength of 5 on expiry input', () => {
    renderPage()
    expect(getExpiryInput()).toHaveAttribute('maxLength', '5')
  })

  it('does not add a slash for 2 or fewer digits', () => {
    renderPage()
    const input = getExpiryInput()

    fireEvent.change(input, { target: { value: '12' } })
    expect(input).toHaveValue('12')
  })

  it('inserts a forward slash after the 2nd digit (MM → MM/)', () => {
    renderPage()
    const input = getExpiryInput()

    fireEvent.change(input, { target: { value: '123' } })
    expect(input).toHaveValue('12/3')
  })

  it('formats a full expiry as MM/YY', () => {
    renderPage()
    const input = getExpiryInput()

    fireEvent.change(input, { target: { value: '1226' } })
    expect(input).toHaveValue('12/26')
  })

  it('strips non-digit characters from expiry input', () => {
    renderPage()
    const input = getExpiryInput()

    fireEvent.change(input, { target: { value: '12/26' } })
    // The slash is re-inserted from the digit logic, result should still be 12/26
    expect(input).toHaveValue('12/26')
  })

  it('truncates expiry beyond 4 digits', () => {
    renderPage()
    const input = getExpiryInput()

    fireEvent.change(input, { target: { value: '123456' } })
    expect(input).toHaveValue('12/34')
  })
})

// ─────────────────────────────────────────────────────────────────────────────
describe('PaymentPage — form validation with formatted inputs', () => {
  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('shows validation error when card number (after stripping hyphens) is too short', async () => {
    renderPage()

    // Wait for the preview to finish loading (button is disabled while previewLoading=true)
    await waitFor(() =>
      expect(screen.getByRole('button', { name: /pay.*securely/i })).not.toBeDisabled()
    )

    fireEvent.change(screen.getByPlaceholderText('John Doe'), { target: { value: 'John Doe' } })
    // Enter only 8 digits → "1234-5678" (8 raw digits, less than 12)
    fireEvent.change(getCardNumberInput(), { target: { value: '12345678' } })
    fireEvent.change(getExpiryInput(), { target: { value: '1226' } })
    fireEvent.change(screen.getByPlaceholderText('***'), { target: { value: '123' } })

    fireEvent.click(screen.getByRole('button', { name: /pay.*securely/i }))

    await waitFor(() =>
      expect(screen.getByText(/please enter valid card details/i)).toBeInTheDocument()
    )
  })

  it('accepts a fully formatted 16-digit card number and proceeds to submit', async () => {
    ;(api.confirmCheckout as jest.Mock).mockResolvedValue({ booking_ref: 'MOCK-001' })

    renderPage()

    await waitFor(() =>
      expect(screen.getByRole('button', { name: /pay.*securely/i })).not.toBeDisabled()
    )

    fireEvent.change(screen.getByPlaceholderText('John Doe'), { target: { value: 'John Doe' } })
    fireEvent.change(getCardNumberInput(), { target: { value: '1234567890123456' } })
    fireEvent.change(getExpiryInput(), { target: { value: '1226' } })
    fireEvent.change(screen.getByPlaceholderText('***'), { target: { value: '123' } })

    fireEvent.click(screen.getByRole('button', { name: /pay.*securely/i }))

    await waitFor(() =>
      expect(api.confirmCheckout).toHaveBeenCalled()
    )
  })

  it('shows error message when confirmCheckout API fails', async () => {
    ;(api.confirmCheckout as jest.Mock).mockRejectedValue(new Error('Payment declined'))

    renderPage()

    await waitFor(() =>
      expect(screen.getByRole('button', { name: /pay.*securely/i })).not.toBeDisabled()
    )

    fireEvent.change(screen.getByPlaceholderText('John Doe'), { target: { value: 'Jane Smith' } })
    fireEvent.change(getCardNumberInput(), { target: { value: '4111111111111111' } })
    fireEvent.change(getExpiryInput(), { target: { value: '0927' } })
    fireEvent.change(screen.getByPlaceholderText('***'), { target: { value: '321' } })

    fireEvent.click(screen.getByRole('button', { name: /pay.*securely/i }))

    await waitFor(() =>
      expect(screen.getByText(/payment declined/i)).toBeInTheDocument()
    )
  })
})

// ─────────────────────────────────────────────────────────────────────────────
describe('PaymentPage — order summary', () => {
  beforeEach(() => {
    jest.clearAllMocks()
  })

  it('renders the page heading', () => {
    renderPage()
    expect(screen.getByRole('heading', { name: /secure payment checkout/i })).toBeInTheDocument()
  })

  it('displays selected seat codes from query params', async () => {
    renderPage()
    await waitFor(() =>
      expect(screen.getByText('A1, A2')).toBeInTheDocument()
    )
  })

  it('navigates back when Change button is clicked', () => {
    renderPage()
    fireEvent.click(screen.getByRole('button', { name: /change/i }))
    expect(mockRouterBack).toHaveBeenCalled()
  })

  it('shows countdown timer', () => {
    renderPage()
    // Timer starts at 10:00
    expect(screen.getByText(/10:00/)).toBeInTheDocument()
  })
})
