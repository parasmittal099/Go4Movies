"use client"

import { useState, useEffect } from "react"
import { useParams, useRouter } from "next/navigation"
import { Movie } from "@/lib/types"
import { fetchMovieById } from "@/lib/api"
import { MovieDetail } from "@/components/movies/movie-detail"
import { SeatSelection } from "@/components/booking/seat-selection"

export default function MoviePage() {
  const params = useParams()
  const router = useRouter()
  const [movie, setMovie] = useState<Movie | null>(null)
  const [loading, setLoading] = useState(true)
  const [showSeats, setShowSeats] = useState(false)

  useEffect(() => {
    const movieId = params.id as string

    fetchMovieById(movieId)
      .then((data) => {
        setMovie(data)
        setLoading(false)
      })
      .catch((error) => {
        console.error("Error fetching movie:", error)
        router.push("/movies")
      })
  }, [params.id, router])

  const handleBookTickets = () => {
    setShowSeats(true)
  }

  if (loading) {
    return <div className="container mx-auto p-4">Loading movie...</div>
  }

  if (!movie) {
    return <div className="container mx-auto p-4">Movie not found</div>
  }

  if (showSeats) {
    return (
      <SeatSelection
        movie={{
          title: movie.title,
          time: "Today, 19:30",
          format: "IMAX 2D",
          hall: "Hall 4"
        }}
        onProceed={(seats) => alert(`Proceeding to payment with seats: ${seats.map(s => s.id).join(', ')}`)}
        onBack={() => setShowSeats(false)}
        onChangeMovie={() => router.push("/movies")}
      />
    )
  }

  return <MovieDetail movie={movie} onBookTickets={handleBookTickets} />
}