import React from 'react'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { LoginForm } from '../login-form'
import { AuthProvider } from '@/context/auth-context'
import * as api from '@/lib/api'

jest.mock('next/navigation', () => ({
  useRouter: () => ({ push: jest.fn() }),
}))

jest.mock('@/lib/api', () => ({
  loginUser: jest.fn(),
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

function renderLoginForm() {
  return render(
    <AuthProvider>
      <LoginForm />
    </AuthProvider>
  )
}

// Use placeholder text to target specific inputs (avoids ambiguity with toggle button aria-label)
const emailPlaceholder = /you@example.com/i
const passwordPlaceholder = /enter your password/i

describe('LoginForm', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    localStorage.clear()
  })

  it('renders email and password inputs', () => {
    renderLoginForm()
    expect(screen.getByPlaceholderText(emailPlaceholder)).toBeInTheDocument()
    expect(screen.getByPlaceholderText(passwordPlaceholder)).toBeInTheDocument()
  })

  it('renders Sign In button', () => {
    renderLoginForm()
    expect(screen.getByRole('button', { name: /sign in/i })).toBeInTheDocument()
  })

  it('shows error when email and password are empty', async () => {
    renderLoginForm()
    const form = screen.getByRole('button', { name: /sign in/i }).closest('form')!
    fireEvent.submit(form)
    await waitFor(() =>
      expect(screen.getByRole('alert')).toHaveTextContent('Email and password are required')
    )
  })

  it('shows error when only email is provided', async () => {
    renderLoginForm()
    fireEvent.change(screen.getByPlaceholderText(emailPlaceholder), { target: { value: 'a@b.com' } })
    const form = screen.getByPlaceholderText(emailPlaceholder).closest('form')!
    fireEvent.submit(form)
    await waitFor(() =>
      expect(screen.getByRole('alert')).toHaveTextContent('Email and password are required')
    )
  })

  it('calls loginUser API on valid input', async () => {
    ;(api.loginUser as jest.Mock).mockResolvedValue({ message: 'ok', user: mockUser })
    renderLoginForm()
    fireEvent.change(screen.getByPlaceholderText(emailPlaceholder), { target: { value: 'a@b.com' } })
    fireEvent.change(screen.getByPlaceholderText(passwordPlaceholder), { target: { value: 'secret' } })
    fireEvent.click(screen.getByRole('button', { name: /sign in/i }))
    await waitFor(() => expect(api.loginUser).toHaveBeenCalledWith({ email: 'a@b.com', password: 'secret' }))
  })

  it('shows server error message on login failure', async () => {
    ;(api.loginUser as jest.Mock).mockRejectedValue(new Error('Invalid credentials'))
    renderLoginForm()
    fireEvent.change(screen.getByPlaceholderText(emailPlaceholder), { target: { value: 'a@b.com' } })
    fireEvent.change(screen.getByPlaceholderText(passwordPlaceholder), { target: { value: 'wrongpass' } })
    fireEvent.click(screen.getByRole('button', { name: /sign in/i }))
    await waitFor(() =>
      expect(screen.getByRole('alert')).toHaveTextContent('Invalid credentials')
    )
  })

  it('toggles password visibility', () => {
    renderLoginForm()
    const passwordInput = screen.getByPlaceholderText(passwordPlaceholder)
    expect(passwordInput).toHaveAttribute('type', 'password')
    const toggleBtn = screen.getByRole('button', { name: /show password/i })
    fireEvent.click(toggleBtn)
    expect(passwordInput).toHaveAttribute('type', 'text')
    fireEvent.click(screen.getByRole('button', { name: /hide password/i }))
    expect(passwordInput).toHaveAttribute('type', 'password')
  })

  it('disables submit button while submitting', async () => {
    let resolveLogin: (val: unknown) => void
    ;(api.loginUser as jest.Mock).mockReturnValue(
      new Promise((res) => { resolveLogin = res })
    )
    renderLoginForm()
    fireEvent.change(screen.getByPlaceholderText(emailPlaceholder), { target: { value: 'a@b.com' } })
    fireEvent.change(screen.getByPlaceholderText(passwordPlaceholder), { target: { value: 'secret' } })
    fireEvent.click(screen.getByRole('button', { name: /sign in/i }))
    await waitFor(() => expect(screen.getByRole('button', { name: /signing in/i })).toBeDisabled())
    resolveLogin!({ message: 'ok', user: mockUser })
  })
})
