"use client"

import { useEffect, useState } from "react"
import { useRouter } from "next/navigation"
import Link from "next/link"
import { useAuth } from "@/context/auth-context"
import { fetchUserBookings } from "@/lib/api"
import type { BookingDetail } from "@/lib/types"

type TabKey = "upcoming" | "past"

/** Compare a show_date string (e.g. "2025-06-15") to today to determine upcoming vs past. */
function isUpcoming(showDate: string): boolean {
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const date = new Date(showDate + "T00:00:00")
  return date >= today
}

function formatDate(showDate: string): string {
  try {
    const d = new Date(showDate + "T00:00:00")
    return d.toLocaleDateString("en-US", {
      weekday: "long",
      month: "short",
      day: "numeric",
    })
  } catch {
    return showDate
  }
}

export default function BookingsPage() {
  const router = useRouter()
  const { user, isReady } = useAuth()

  const [bookings, setBookings] = useState<BookingDetail[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [tab, setTab] = useState<TabKey>("upcoming")

  useEffect(() => {
    if (!isReady) return
    if (!user?.id) {
      setLoading(false)
      return
    }
    fetchUserBookings(user.id)
      .then(setBookings)
      .catch((err) =>
        setError(err instanceof Error ? err.message : "Failed to fetch bookings")
      )
      .finally(() => setLoading(false))
  }, [user?.id, isReady])

  const upcoming = bookings.filter(
    (b) => b.status === "CONFIRMED" && isUpcoming(b.show_date)
  )
  const past = bookings.filter(
    (b) => b.status !== "CONFIRMED" || !isUpcoming(b.show_date)
  )
  const activeList = tab === "upcoming" ? upcoming : past

  if (!isReady || loading) {
    return (
      <main className="min-h-screen flex items-center justify-center">
        <div className="flex flex-col items-center gap-4 animate-pulse">
          <span className="material-symbols-outlined text-5xl text-primary">
            confirmation_number
          </span>
          <p className="text-neutral-400 text-lg">Loading your bookings…</p>
        </div>
      </main>
    )
  }

  if (!user) {
    return (
      <main className="min-h-screen flex items-center justify-center px-4">
        <div className="text-center space-y-5">
          <span className="material-symbols-outlined text-6xl text-neutral-600">
            lock
          </span>
          <h2 className="text-2xl font-bold text-white">Sign in to view your bookings</h2>
          <Link
            href="/login"
            className="inline-flex items-center gap-2 px-6 py-3 bg-primary text-white font-bold rounded-lg hover:brightness-110 transition"
          >
            Sign In
          </Link>
        </div>
      </main>
    )
  }

  return (
    <main className="min-h-screen pb-24">
      <div className="max-w-7xl mx-auto px-4 md:px-8 pt-8 md:pt-12">
        {/* ─── Header ─── */}
        <header className="mb-10">
          <h1 className="text-4xl md:text-5xl font-black tracking-tight text-white mb-2">
            My Bookings
          </h1>
          <p className="text-neutral-500 text-lg">
            Manage your upcoming cinematic experiences and past memories.
          </p>
        </header>

        {/* ─── Tabs ─── */}
        <div className="flex gap-8 mb-10 border-b border-neutral-800">
          <button
            onClick={() => setTab("upcoming")}
            className={`pb-4 font-bold text-lg transition-colors ${
              tab === "upcoming"
                ? "text-primary border-b-2 border-primary"
                : "text-neutral-500 hover:text-neutral-300"
            }`}
          >
            Upcoming ({upcoming.length})
          </button>
          <button
            onClick={() => setTab("past")}
            className={`pb-4 font-bold text-lg transition-colors ${
              tab === "past"
                ? "text-primary border-b-2 border-primary"
                : "text-neutral-500 hover:text-neutral-300"
            }`}
          >
            Past History ({past.length})
          </button>
        </div>

        {/* ─── Error State ─── */}
        {error && (
          <div className="bg-red-500/10 border border-red-500/20 rounded-xl p-4 mb-8 flex items-center gap-3">
            <span className="material-symbols-outlined text-red-400">
              error_outline
            </span>
            <p className="text-red-300 text-sm">{error}</p>
          </div>
        )}

        {/* ─── Empty State ─── */}
        {activeList.length === 0 && !error && (
          <div className="flex flex-col items-center justify-center py-24 text-center">
            <span className="material-symbols-outlined text-6xl text-neutral-700 mb-4">
              {tab === "upcoming" ? "event_available" : "history"}
            </span>
            <h3 className="text-xl font-bold text-neutral-400 mb-2">
              {tab === "upcoming"
                ? "No upcoming bookings"
                : "No past bookings yet"}
            </h3>
            <p className="text-neutral-600 mb-6 max-w-md">
              {tab === "upcoming"
                ? "Browse our latest movies and book your next cinematic adventure."
                : "Your viewing history will appear here after you attend your first show."}
            </p>
            <Link
              href="/movies"
              className="inline-flex items-center gap-2 px-6 py-3 bg-primary text-white font-bold rounded-lg hover:brightness-110 transition"
            >
              <span className="material-symbols-outlined text-sm">
                movie
              </span>
              Browse Movies
            </Link>
          </div>
        )}

        {/* ─── Bookings Grid ─── */}
        {activeList.length > 0 && (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {activeList.map((b) => (
              <BookingCard
                key={b.id}
                booking={b}
                isPast={tab === "past"}
                onViewTicket={() =>
                  router.push(`/booking-confirmation?ref=${b.booking_ref}`)
                }
              />
            ))}
          </div>
        )}
      </div>
    </main>
  )
}

/* ─────────────────────────────────────────────────────────────────────────────
 *  Booking Card Component
 * ─────────────────────────────────────────────────────────────────────────── */

function BookingCard({
  booking,
  isPast,
  onViewTicket,
}: {
  booking: BookingDetail
  isPast: boolean
  onViewTicket: () => void
}) {
  const seatLabels = booking.seats
    .map((s) => `${s.row_label}${s.col_number}`)
    .join(", ")

  return (
    <div
      className={`group relative flex flex-col md:flex-row bg-neutral-900/60 backdrop-blur-sm rounded-xl border border-neutral-800 overflow-hidden transition-all duration-300 shadow-xl shadow-black/30 ${
        isPast
          ? "opacity-70 hover:opacity-100"
          : "hover:scale-[1.02] hover:border-neutral-700"
      }`}
    >
      {/* ─── Poster ─── */}
      <div className="relative w-full md:w-44 aspect-[2/3] shrink-0">
        {booking.movie_poster ? (
          <img
            alt={booking.movie_title}
            className={`w-full h-full object-cover ${
              isPast ? "grayscale group-hover:grayscale-0 transition-all duration-500" : ""
            }`}
            src={booking.movie_poster}
          />
        ) : (
          <div className="w-full h-full bg-neutral-800 flex items-center justify-center">
            <span className="material-symbols-outlined text-5xl text-neutral-600">
              movie
            </span>
          </div>
        )}
        {!isPast && (
          <div className="absolute top-3 left-3 bg-primary text-white text-[10px] font-bold px-2 py-1 rounded-sm uppercase tracking-widest">
            Digital Ticket
          </div>
        )}
        {isPast && (
          <div className="absolute top-3 left-3 bg-neutral-700 text-neutral-300 text-[10px] font-bold px-2 py-1 rounded-sm uppercase tracking-widest">
            Watched
          </div>
        )}
      </div>

      {/* ─── Details ─── */}
      <div className="p-5 md:p-6 flex flex-col justify-between flex-grow">
        <div>
          <div className="flex justify-between items-start mb-2 gap-3">
            <h2 className="text-xl md:text-2xl font-bold text-white group-hover:text-primary transition-colors leading-tight">
              {booking.movie_title}
            </h2>
            <span className="bg-white/10 px-2 py-1 rounded text-[10px] font-bold uppercase text-cyan-300 shrink-0">
              {booking.format}
            </span>
          </div>

          <div className="space-y-2.5 mt-4">
            <div className="flex items-center gap-3 text-neutral-400">
              <span className="material-symbols-outlined text-primary text-lg">
                calendar_today
              </span>
              <span className="font-semibold text-white text-sm">
                {formatDate(booking.show_date)} • {booking.start_time}
              </span>
            </div>
            <div className="flex items-center gap-3 text-neutral-400">
              <span className="material-symbols-outlined text-primary text-lg">
                location_on
              </span>
              <span className="text-sm">
                {booking.theater_name}, {booking.screen_name}
              </span>
            </div>
            <div className="flex items-center gap-3 text-neutral-400">
              <span className="material-symbols-outlined text-primary text-lg">
                event_seat
              </span>
              <span className="text-sm">{seatLabels}</span>
            </div>
          </div>
        </div>

        <div className="mt-6 flex items-center justify-between">
          <div className="flex flex-col">
            <span className="text-[10px] font-bold uppercase tracking-widest text-neutral-600">
              Booking ID
            </span>
            <span className="font-mono text-xs text-neutral-400">
              {booking.booking_ref}
            </span>
          </div>
          <button
            onClick={onViewTicket}
            className={`px-5 py-2.5 rounded-lg font-bold text-sm flex items-center gap-2 transition-all ${
              isPast
                ? "bg-neutral-800 text-neutral-300 hover:bg-neutral-700"
                : "bg-primary text-white hover:scale-105 active:scale-95 shadow-lg shadow-primary/25"
            }`}
          >
            <span className="material-symbols-outlined text-sm">
              confirmation_number
            </span>
            {isPast ? "View Details" : "View Ticket"}
          </button>
        </div>
      </div>
    </div>
  )
}
