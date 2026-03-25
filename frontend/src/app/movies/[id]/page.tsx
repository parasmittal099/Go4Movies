"use client"

import { useState, useEffect } from "react"
import { useParams, useRouter } from "next/navigation"
import { Movie, ShowtimesResponse } from "@/lib/types"
import { fetchMovieById, fetchMovieShowtimes } from "@/lib/api"
import { MovieDetail } from "@/components/movies/movie-detail"
import { SeatSelection } from "@/components/booking/seat-selection"

export default function MoviePage() {
  const params = useParams()
  const router = useRouter()
  const [movie, setMovie] = useState<Movie | null>(null)
  const [loading, setLoading] = useState(true)
  const [showtimes, setShowtimes] = useState<ShowtimesResponse | null>(null)
  const [showtimesLoading, setShowtimesLoading] = useState(true)
  const [showSeats, setShowSeats] = useState(false)
  const [selectedShowtimeId, setSelectedShowtimeId] = useState<number | null>(null)

  useEffect(() => {
    const movieId = params.id as string
    const zipCode = localStorage.getItem("selectedZipCode") || ""

    fetchMovieById(movieId)
      .then((data) => {
        setMovie(data)
        setLoading(false)
      })
      .catch((error) => {
        console.error("Error fetching movie:", error)
        router.push("/movies")
      })

    if (zipCode) {
      fetchMovieShowtimes(movieId, zipCode)
        .then((data) => {
          setShowtimes(data)
          setShowtimesLoading(false)
        })
        .catch((error) => {
          console.error("Error fetching showtimes:", error)
          setShowtimesLoading(false)
        })
    } else {
      setShowtimesLoading(false)
    }
  }, [params.id, router])

  const handleBookTickets = (showtimeId: number) => {
    setSelectedShowtimeId(showtimeId)
    setShowSeats(true)
  }

  if (loading) {
    return <div className="container mx-auto p-4">Loading movie...</div>
  }

  if (!movie) {
    return <div className="container mx-auto p-4">Movie not found</div>
  }

  if (showSeats) {
    const selectedShowtime = showtimes?.theaters
      .flatMap((t) => t.showtimes)
      .find((st) => st.id === selectedShowtimeId)

    return (
      <SeatSelection
        movie={{
          title: movie.title,
          time: selectedShowtime ? `${selectedShowtime.show_date}, ${selectedShowtime.start_time}` : "",
          format: selectedShowtime?.format ?? "",
          hall: selectedShowtime?.screen_name ?? "",
        }}
        onProceed={(seats) => alert(`Proceeding to payment with seats: ${seats.map(s => s.id).join(', ')}`)}
        onBack={() => setShowSeats(false)}
        onChangeMovie={() => router.push("/movies")}
      />
    )
  }

  return (
    <MovieDetail
      movie={movie}
      showtimes={showtimes}
      showtimesLoading={showtimesLoading}
      onBookTickets={handleBookTickets}
    />
  )
}