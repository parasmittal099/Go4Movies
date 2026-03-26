import React from 'react'
import { render, screen, fireEvent } from '@testing-library/react'
import { MovieCard } from '../movie-card'
import type { Movie } from '@/lib/types'

const movie: Movie = {
  id: 1,
  title: 'Dune: Part Two',
  description: 'Epic sci-fi epic set on Arrakis.',
  poster_url: '/dune.jpg',
}

describe('MovieCard', () => {
  it('renders the movie title', () => {
    render(<MovieCard movie={movie} onClick={jest.fn()} />)
    expect(screen.getByText('Dune: Part Two')).toBeInTheDocument()
  })

  it('renders the movie description', () => {
    render(<MovieCard movie={movie} onClick={jest.fn()} />)
    expect(screen.getByText(movie.description)).toBeInTheDocument()
  })

  it('renders the poster image with correct src and alt', () => {
    render(<MovieCard movie={movie} onClick={jest.fn()} />)
    const img = screen.getByRole('img')
    expect(img).toHaveAttribute('src', '/dune.jpg')
    expect(img).toHaveAttribute('alt', 'Dune: Part Two')
  })

  it('calls onClick with the movie id when clicked', () => {
    const handleClick = jest.fn()
    render(<MovieCard movie={movie} onClick={handleClick} />)
    fireEvent.click(screen.getByText('Dune: Part Two'))
    expect(handleClick).toHaveBeenCalledWith(1)
  })
})
