import React from 'react'
import { render, screen, fireEvent } from '@testing-library/react'
import { SearchBar } from '../search-bar'

// SearchBar uses useRouter — mock it
jest.mock('next/navigation', () => ({
  useRouter: () => ({ push: jest.fn() }),
}))

describe('SearchBar', () => {
  beforeEach(() => {
    // jsdom does not implement window.alert – spy and suppress
    jest.spyOn(window, 'alert').mockImplementation(() => {})
  })

  afterEach(() => {
    jest.restoreAllMocks()
  })

  it('renders input and button', () => {
    render(<SearchBar />)
    expect(screen.getByPlaceholderText(/search movies/i)).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /search/i })).toBeInTheDocument()
  })

  it('updates input value on change', () => {
    render(<SearchBar />)
    const input = screen.getByPlaceholderText(/search movies/i)
    fireEvent.change(input, { target: { value: 'Dune' } })
    expect(input).toHaveValue('Dune')
  })

  it('calls alert with search query on form submit when query is non-empty', () => {
    render(<SearchBar />)
    const input = screen.getByPlaceholderText(/search movies/i)
    fireEvent.change(input, { target: { value: 'Inception' } })
    fireEvent.submit(screen.getByRole('button', { name: /search/i }).closest('form')!)
    expect(window.alert).toHaveBeenCalledWith('Searching for: Inception')
  })

  it('does not call alert when query is empty', () => {
    render(<SearchBar />)
    fireEvent.submit(screen.getByRole('button', { name: /search/i }).closest('form')!)
    expect(window.alert).not.toHaveBeenCalled()
  })

  it('does not call alert when query is only whitespace', () => {
    render(<SearchBar />)
    const input = screen.getByPlaceholderText(/search movies/i)
    fireEvent.change(input, { target: { value: '   ' } })
    fireEvent.submit(screen.getByRole('button', { name: /search/i }).closest('form')!)
    expect(window.alert).not.toHaveBeenCalled()
  })
})
