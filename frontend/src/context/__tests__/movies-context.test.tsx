import React from 'react'
import { render, screen, act } from '@testing-library/react'
import { MoviesProvider, useMovies } from '../movies-context'
import type { Movie } from '@/lib/types'

const sampleMovies: Movie[] = [
  { id: 1, title: 'Dune', description: 'Epic sci-fi', poster_url: '/dune.jpg' },
  { id: 2, title: 'Inception', description: 'Dream heist', poster_url: '/inception.jpg' },
]

function TestConsumer() {
  const { movies, setMovies } = useMovies()
  return (
    <div>
      <div data-testid="movie-count">{movies.length}</div>
      <ul>
        {movies.map((m) => (
          <li key={m.id} data-testid="movie-title">
            {m.title}
          </li>
        ))}
      </ul>
      <button data-testid="set-movies-btn" onClick={() => setMovies(sampleMovies)}>
        Set Movies
      </button>
      <button data-testid="clear-movies-btn" onClick={() => setMovies([])}>
        Clear Movies
      </button>
    </div>
  )
}

describe('MoviesProvider', () => {
  it('starts with empty movies array', () => {
    render(
      <MoviesProvider>
        <TestConsumer />
      </MoviesProvider>
    )
    expect(screen.getByTestId('movie-count').textContent).toBe('0')
  })

  it('setMovies updates the movies list', () => {
    render(
      <MoviesProvider>
        <TestConsumer />
      </MoviesProvider>
    )
    act(() => {
      screen.getByTestId('set-movies-btn').click()
    })
    expect(screen.getByTestId('movie-count').textContent).toBe('2')
    const titles = screen.getAllByTestId('movie-title').map((el) => el.textContent)
    expect(titles).toEqual(['Dune', 'Inception'])
  })

  it('clears movies list when setMovies is called with []', () => {
    render(
      <MoviesProvider>
        <TestConsumer />
      </MoviesProvider>
    )
    act(() => {
      screen.getByTestId('set-movies-btn').click()
    })
    expect(screen.getByTestId('movie-count').textContent).toBe('2')
    act(() => {
      screen.getByTestId('clear-movies-btn').click()
    })
    expect(screen.getByTestId('movie-count').textContent).toBe('0')
  })
})

describe('useMovies', () => {
  it('throws when used outside MoviesProvider', () => {
    const spy = jest.spyOn(console, 'error').mockImplementation(() => {})
    expect(() => render(<TestConsumer />)).toThrow(
      'useMovies must be used within a <MoviesProvider>'
    )
    spy.mockRestore()
  })
})
