import React from 'react'
import { render, screen, fireEvent } from '@testing-library/react'
import { MovieGrid } from '../movie-grid'
import type { Movie } from '@/lib/types'

const movies: Movie[] = [
  { id: 1, title: 'Dune', description: 'Sci-fi epic', poster_url: '/dune.jpg' },
  { id: 2, title: 'Inception', description: 'Dream heist', poster_url: '/inception.jpg' },
  { id: 3, title: 'Interstellar', description: 'Space odyssey', poster_url: '/interstellar.jpg' },
]

describe('MovieGrid', () => {
  it('renders all movie cards', () => {
    render(<MovieGrid movies={movies} onMovieClick={jest.fn()} />)
    expect(screen.getByText('Dune')).toBeInTheDocument()
    expect(screen.getByText('Inception')).toBeInTheDocument()
    expect(screen.getByText('Interstellar')).toBeInTheDocument()
  })

  it('renders the correct number of movie cards', () => {
    render(<MovieGrid movies={movies} onMovieClick={jest.fn()} />)
    expect(screen.getAllByRole('img')).toHaveLength(3)
  })

  it('renders empty grid when no movies', () => {
    const { container } = render(<MovieGrid movies={[]} onMovieClick={jest.fn()} />)
    expect(container.querySelectorAll('img')).toHaveLength(0)
  })

  it('calls onMovieClick with the correct movie id', () => {
    const handleClick = jest.fn()
    render(<MovieGrid movies={movies} onMovieClick={handleClick} />)
    fireEvent.click(screen.getByText('Inception'))
    expect(handleClick).toHaveBeenCalledWith(2)
  })
})
