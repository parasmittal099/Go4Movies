import React from 'react'
import { render, screen, fireEvent, waitFor } from '@testing-library/react'
import { SignUpForm } from '../signup-form'
import { AuthProvider } from '@/context/auth-context'
import * as api from '@/lib/api'

jest.mock('next/navigation', () => ({
  useRouter: () => ({ push: jest.fn() }),
}))

jest.mock('@/lib/api', () => ({
  registerUser: jest.fn(),
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

function renderSignUpForm() {
  return render(
    <AuthProvider>
      <SignUpForm />
    </AuthProvider>
  )
}

// Use placeholder text to target specific inputs to avoid ambiguity with labels
const fullNamePh = /john doe/i
const usernamePh = /john_doe/i
const emailPh = /you@example.com/i
const passwordPh = /create a password/i
const confirmPasswordPh = /confirm your password/i

function fillValidForm() {
  fireEvent.change(screen.getByPlaceholderText(fullNamePh), { target: { value: 'Bob Smith' } })
  fireEvent.change(screen.getByPlaceholderText(usernamePh), { target: { value: 'bobsmith' } })
  fireEvent.change(screen.getByPlaceholderText(emailPh), { target: { value: 'a@b.com' } })
  fireEvent.change(screen.getByPlaceholderText(passwordPh), { target: { value: 'secret123' } })
  fireEvent.change(screen.getByPlaceholderText(confirmPasswordPh), { target: { value: 'secret123' } })
  fireEvent.click(screen.getByRole('checkbox'))
}

describe('SignUpForm', () => {
  beforeEach(() => {
    jest.clearAllMocks()
    localStorage.clear()
  })

  it('renders all form fields', () => {
    renderSignUpForm()
    expect(screen.getByPlaceholderText(fullNamePh)).toBeInTheDocument()
    expect(screen.getByPlaceholderText(usernamePh)).toBeInTheDocument()
    expect(screen.getByPlaceholderText(emailPh)).toBeInTheDocument()
    expect(screen.getByPlaceholderText(passwordPh)).toBeInTheDocument()
    expect(screen.getByPlaceholderText(confirmPasswordPh)).toBeInTheDocument()
  })

  it('shows error when full name is missing', async () => {
    renderSignUpForm()
    const form = screen.getByRole('button', { name: /create account/i }).closest('form')!
    fireEvent.submit(form)
    await waitFor(() =>
      expect(screen.getByRole('alert')).toHaveTextContent('Full name is required')
    )
  })

  it('shows error when username is too short', async () => {
    renderSignUpForm()
    fireEvent.change(screen.getByPlaceholderText(fullNamePh), { target: { value: 'Bob Smith' } })
    fireEvent.change(screen.getByPlaceholderText(usernamePh), { target: { value: 'ab' } })
    const form = screen.getByPlaceholderText(fullNamePh).closest('form')!
    fireEvent.submit(form)
    await waitFor(() =>
      expect(screen.getByRole('alert')).toHaveTextContent('Username must be at least 3 characters')
    )
  })

  it('shows error when password is too short', async () => {
    renderSignUpForm()
    fireEvent.change(screen.getByPlaceholderText(fullNamePh), { target: { value: 'Bob Smith' } })
    fireEvent.change(screen.getByPlaceholderText(usernamePh), { target: { value: 'bobsmith' } })
    fireEvent.change(screen.getByPlaceholderText(passwordPh), { target: { value: '123' } })
    const form = screen.getByPlaceholderText(fullNamePh).closest('form')!
    fireEvent.submit(form)
    await waitFor(() =>
      expect(screen.getByRole('alert')).toHaveTextContent('Password must be at least 6 characters')
    )
  })

  it('shows error when passwords do not match', async () => {
    renderSignUpForm()
    fireEvent.change(screen.getByPlaceholderText(fullNamePh), { target: { value: 'Bob Smith' } })
    fireEvent.change(screen.getByPlaceholderText(usernamePh), { target: { value: 'bobsmith' } })
    fireEvent.change(screen.getByPlaceholderText(passwordPh), { target: { value: 'secret123' } })
    fireEvent.change(screen.getByPlaceholderText(confirmPasswordPh), { target: { value: 'different' } })
    const form = screen.getByPlaceholderText(fullNamePh).closest('form')!
    fireEvent.submit(form)
    await waitFor(() =>
      expect(screen.getByRole('alert')).toHaveTextContent('Passwords do not match')
    )
  })

  it('calls registerUser with correct payload on valid submit', async () => {
    ;(api.registerUser as jest.Mock).mockResolvedValue({ message: 'ok', user: mockUser })
    renderSignUpForm()
    fillValidForm()
    fireEvent.click(screen.getByRole('button', { name: /create account/i }))
    await waitFor(() =>
      expect(api.registerUser).toHaveBeenCalledWith({
        email: 'a@b.com',
        username: 'bobsmith',
        password: 'secret123',
        full_name: 'Bob Smith',
      })
    )
  })

  it('shows server error message on registration failure', async () => {
    ;(api.registerUser as jest.Mock).mockRejectedValue(new Error('Email already exists'))
    renderSignUpForm()
    fillValidForm()
    fireEvent.click(screen.getByRole('button', { name: /create account/i }))
    await waitFor(() =>
      expect(screen.getByRole('alert')).toHaveTextContent('Email already exists')
    )
  })

  it('toggles password visibility', () => {
    renderSignUpForm()
    const passwordInput = screen.getByPlaceholderText(passwordPh)
    expect(passwordInput).toHaveAttribute('type', 'password')
    // There are 2 toggle buttons (password + confirm password), get the first one
    fireEvent.click(screen.getAllByRole('button', { name: /show password/i })[0])
    expect(passwordInput).toHaveAttribute('type', 'text')
  })
})
