import React from 'react'
import { render, screen, fireEvent, waitFor, act } from '@testing-library/react'
import { Header } from '../header'
import { AuthProvider } from '@/context/auth-context'
import { MoviesProvider } from '@/context/movies-context'

const mockPush = jest.fn()

jest.mock('next/navigation', () => ({
  useRouter: () => ({ push: mockPush }),
}))

const mockUser = {
  id: 1,
  email: 'a@b.com',
  username: 'bob',
  full_name: 'Bob Smith',
  phone: null,
  created_at: '',
  updated_at: '',
}

function renderHeader() {
  return render(
    <AuthProvider>
      <MoviesProvider>
        <Header />
      </MoviesProvider>
    </AuthProvider>
  )
}

describe('Header (unauthenticated)', () => {
  beforeEach(() => {
    localStorage.clear()
    mockPush.mockClear()
  })

  it('renders the Go4Movies logo link', async () => {
    renderHeader()
    await waitFor(() => expect(screen.getByText('Go4Movies')).toBeInTheDocument())
  })

  it('shows Sign In link when not logged in', async () => {
    renderHeader()
    await waitFor(() => expect(screen.getByRole('link', { name: /sign in/i })).toBeInTheDocument())
  })

  it('does not show Sign Out button when not logged in', async () => {
    renderHeader()
    await waitFor(() => expect(screen.queryByRole('button', { name: /sign out/i })).not.toBeInTheDocument())
  })
})

describe('Header (authenticated)', () => {
  beforeEach(() => {
    localStorage.clear()
    localStorage.setItem('go4movies_user', JSON.stringify(mockUser))
    mockPush.mockClear()
  })

  it('shows user initials when logged in', async () => {
    renderHeader()
    await waitFor(() => expect(screen.getByText('BS')).toBeInTheDocument())
  })

  it('shows full name when logged in', async () => {
    renderHeader()
    await waitFor(() => expect(screen.getByText('Bob Smith')).toBeInTheDocument())
  })

  it('clicking Sign Out calls logout and redirects to /', async () => {
    renderHeader()
    await waitFor(() => expect(screen.getByRole('button', { name: /sign out/i })).toBeInTheDocument())
    fireEvent.click(screen.getByRole('button', { name: /sign out/i }))
    expect(mockPush).toHaveBeenCalledWith('/')
  })
})

describe('Header search', () => {
  beforeEach(() => {
    localStorage.clear()
    mockPush.mockClear()
  })

  it('opens search input when search icon button is clicked', async () => {
    renderHeader()
    await waitFor(() => expect(screen.getByText('Go4Movies')).toBeInTheDocument())
    // Find the search toggle button by its icon text
    const searchIconBtn = screen.getByText('search').closest('button')!
    fireEvent.click(searchIconBtn)
    await waitFor(() =>
      expect(screen.getByPlaceholderText(/search movies/i)).toBeInTheDocument()
    )
  })

  it('closes search on Escape key', async () => {
    renderHeader()
    await waitFor(() => expect(screen.getByText('Go4Movies')).toBeInTheDocument())
    const searchIconBtn = screen.getByText('search').closest('button')!
    fireEvent.click(searchIconBtn)
    await waitFor(() => expect(screen.getByPlaceholderText(/search movies/i)).toBeInTheDocument())
    act(() => {
      fireEvent.keyDown(document, { key: 'Escape' })
    })
    await waitFor(() => expect(screen.queryByPlaceholderText(/search movies/i)).not.toBeInTheDocument())
  })

  it('clears query when close (x) button inside search is clicked', async () => {
    renderHeader()
    await waitFor(() => expect(screen.getByText('Go4Movies')).toBeInTheDocument())
    const searchIconBtn = screen.getByText('search').closest('button')!
    fireEvent.click(searchIconBtn)
    await waitFor(() => expect(screen.getByPlaceholderText(/search movies/i)).toBeInTheDocument())
    const searchInput = screen.getByPlaceholderText(/search movies/i)
    fireEvent.change(searchInput, { target: { value: 'Dune' } })
    expect(searchInput).toHaveValue('Dune')
    // The clear (×) button appears when there is a query
    const closeBtn = screen.getByText('close').closest('button')!
    fireEvent.click(closeBtn)
    expect(searchInput).toHaveValue('')
  })
})
