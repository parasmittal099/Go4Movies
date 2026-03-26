import React from 'react'
import { render, screen, act, waitFor } from '@testing-library/react'
import { AuthProvider, useAuth } from '../auth-context'
import type { User } from '@/lib/types'

const mockUser: User = {
  id: 1,
  email: 'test@test.com',
  username: 'testuser',
  full_name: 'Test User',
  phone: null,
  created_at: '2024-01-01',
  updated_at: '2024-01-01',
}

// Helper component that consumes the context
function TestConsumer() {
  const { user, location, isReady, login, logout, setLocation, clearLocation } = useAuth()
  return (
    <div>
      <div data-testid="is-ready">{String(isReady)}</div>
      <div data-testid="user-email">{user?.email ?? 'none'}</div>
      <div data-testid="location-city">{location?.city ?? 'none'}</div>
      <button data-testid="login-btn" onClick={() => login(mockUser)}>
        Login
      </button>
      <button data-testid="logout-btn" onClick={logout}>
        Logout
      </button>
      <button data-testid="set-location-btn" onClick={() => setLocation({ city: 'NYC', zipcode: '10001' })}>
        Set Location
      </button>
      <button data-testid="clear-location-btn" onClick={clearLocation}>
        Clear Location
      </button>
    </div>
  )
}

function renderWithProvider() {
  return render(
    <AuthProvider>
      <TestConsumer />
    </AuthProvider>
  )
}

describe('AuthProvider', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('becomes ready after mount', async () => {
    renderWithProvider()
    await waitFor(() => expect(screen.getByTestId('is-ready').textContent).toBe('true'))
  })

  it('starts with no user', async () => {
    renderWithProvider()
    await waitFor(() => expect(screen.getByTestId('is-ready').textContent).toBe('true'))
    expect(screen.getByTestId('user-email').textContent).toBe('none')
  })

  it('login sets user in state and localStorage', async () => {
    renderWithProvider()
    await waitFor(() => expect(screen.getByTestId('is-ready').textContent).toBe('true'))
    act(() => {
      screen.getByTestId('login-btn').click()
    })
    expect(screen.getByTestId('user-email').textContent).toBe('test@test.com')
    expect(JSON.parse(localStorage.getItem('go4movies_user')!).email).toBe('test@test.com')
  })

  it('logout clears user from state and localStorage', async () => {
    localStorage.setItem('go4movies_user', JSON.stringify(mockUser))
    renderWithProvider()
    await waitFor(() => expect(screen.getByTestId('user-email').textContent).toBe('test@test.com'))
    act(() => {
      screen.getByTestId('logout-btn').click()
    })
    expect(screen.getByTestId('user-email').textContent).toBe('none')
    expect(localStorage.getItem('go4movies_user')).toBeNull()
  })

  it('rehydrates user from localStorage on mount', async () => {
    localStorage.setItem('go4movies_user', JSON.stringify(mockUser))
    renderWithProvider()
    await waitFor(() => expect(screen.getByTestId('user-email').textContent).toBe('test@test.com'))
  })

  it('setLocation persists location to localStorage', async () => {
    renderWithProvider()
    await waitFor(() => expect(screen.getByTestId('is-ready').textContent).toBe('true'))
    act(() => {
      screen.getByTestId('set-location-btn').click()
    })
    expect(screen.getByTestId('location-city').textContent).toBe('NYC')
    expect(JSON.parse(localStorage.getItem('go4movies_location')!).zipcode).toBe('10001')
    expect(localStorage.getItem('selectedZipCode')).toBe('10001')
  })

  it('clearLocation removes location from state and storage', async () => {
    localStorage.setItem('go4movies_location', JSON.stringify({ city: 'NYC', zipcode: '10001' }))
    renderWithProvider()
    await waitFor(() => expect(screen.getByTestId('location-city').textContent).toBe('NYC'))
    act(() => {
      screen.getByTestId('clear-location-btn').click()
    })
    expect(screen.getByTestId('location-city').textContent).toBe('none')
    expect(localStorage.getItem('go4movies_location')).toBeNull()
  })

  it('migrates legacy selectedZipCode key to go4movies_location', async () => {
    localStorage.setItem('selectedZipCode', '90210')
    renderWithProvider()
    await waitFor(() => expect(screen.getByTestId('is-ready').textContent).toBe('true'))
    expect(screen.getByTestId('location-city').textContent).toBe('')
    expect(localStorage.getItem('go4movies_location')).toBeTruthy()
  })
})

describe('useAuth', () => {
  it('throws when used outside AuthProvider', () => {
    const spy = jest.spyOn(console, 'error').mockImplementation(() => {})
    expect(() => render(<TestConsumer />)).toThrow('useAuth must be used within an <AuthProvider>')
    spy.mockRestore()
  })
})
