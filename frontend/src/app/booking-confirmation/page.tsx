"use client"

import { useEffect, useState } from "react"
import { useSearchParams, useRouter } from "next/navigation"
import Link from "next/link"
import QRCode from "react-qr-code"
import { useAuth } from "@/context/auth-context"
import { fetchUserBookings } from "@/lib/api"
import type { BookingDetail } from "@/lib/types"

export default function BookingConfirmationPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { user } = useAuth()

  const bookingRef = searchParams.get("ref") ?? ""
  const [booking, setBooking] = useState<BookingDetail | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!user?.id) {
      setLoading(false)
      return
    }
    fetchUserBookings(user.id)
      .then((bookings) => {
        const found = bookings.find((b) => b.booking_ref === bookingRef)
        if (found) {
          setBooking(found)
        } else {
          setError("Booking not found")
        }
      })
      .catch((err) => {
        setError(err instanceof Error ? err.message : "Failed to load booking")
      })
      .finally(() => setLoading(false))
  }, [user?.id, bookingRef])

  /** Format seat labels, e.g. "G12, G13" */
  const seatLabels = booking
    ? booking.seats.map((s) => `${s.row_label}${s.col_number}`).join(", ")
    : ""

  if (loading) {
    return (
      <main className="min-h-screen flex items-center justify-center">
        <div className="flex flex-col items-center gap-4 animate-pulse">
          <span className="material-symbols-outlined text-5xl text-primary">
            confirmation_number
          </span>
          <p className="text-neutral-400 text-lg">Loading your booking…</p>
        </div>
      </main>
    )
  }

  if (error || !booking) {
    return (
      <main className="min-h-screen flex items-center justify-center px-4">
        <div className="text-center space-y-4">
          <span className="material-symbols-outlined text-5xl text-red-400">
            error_outline
          </span>
          <p className="text-neutral-300 text-lg">{error ?? "Booking not found"}</p>
          <Link
            href="/movies"
            className="inline-flex items-center gap-2 px-6 py-3 bg-primary text-white font-bold rounded-lg hover:brightness-110 transition"
          >
            Browse Movies
          </Link>
        </div>
      </main>
    )
  }

  return (
    <main className="min-h-screen relative overflow-hidden pb-24">
      {/* ─── Background cinematic blur ─── */}
      {booking.movie_poster && (
        <div className="absolute inset-0 z-0 opacity-15 pointer-events-none">
          <div className="absolute inset-0 bg-gradient-to-b from-transparent via-background to-background z-10" />
          <img
            src={booking.movie_poster}
            alt=""
            className="w-full h-full object-cover grayscale brightness-50 blur-sm"
          />
        </div>
      )}

      <div className="screen-content relative z-10 max-w-7xl mx-auto px-4 md:px-8 pt-12 md:pt-16">
        {/* ─── Confirmation Badge ─── */}
        <div className="flex flex-col items-center text-center mb-10 animate-fade-in-up">
          <div className="inline-flex items-center gap-2 px-5 py-2.5 bg-green-500/10 border border-green-500/20 rounded-full mb-6">
            <span
              className="material-symbols-outlined text-green-400 text-sm"
              style={{ fontVariationSettings: "'FILL' 1" }}
            >
              check_circle
            </span>
            <span className="text-xs font-bold text-green-400 uppercase tracking-widest">
              Booking Confirmed
            </span>
          </div>
          <h1 className="text-4xl md:text-5xl font-black tracking-tight text-white mb-2">
            Enjoy the Show
          </h1>
          <p className="text-neutral-400 max-w-lg">
            Your seats are reserved. Check the details below and get ready for an amazing cinematic experience.
          </p>
        </div>

        {/* ─── Content Grid ─── */}
        <div className="max-w-5xl mx-auto grid grid-cols-1 lg:grid-cols-12 gap-8 items-start">
          {/* ─── Main Ticket Card ─── */}
          <div className="lg:col-span-8 flex flex-col gap-6">
            <div className="bg-neutral-900/80 backdrop-blur-xl rounded-2xl overflow-hidden border border-neutral-800 shadow-2xl shadow-black/50">
              <div className="flex flex-col md:flex-row">
                {/* Poster */}
                <div className="w-full md:w-64 h-80 md:h-auto overflow-hidden shrink-0">
                  {booking.movie_poster ? (
                    <img
                      src={booking.movie_poster}
                      alt={booking.movie_title}
                      className="w-full h-full object-cover"
                    />
                  ) : (
                    <div className="w-full h-full bg-neutral-800 flex items-center justify-center">
                      <span className="material-symbols-outlined text-6xl text-neutral-600">
                        movie
                      </span>
                    </div>
                  )}
                </div>

                {/* Details */}
                <div className="flex-1 p-6 md:p-8 flex flex-col justify-between">
                  <div>
                    <h2 className="text-2xl md:text-3xl font-bold text-white mb-2">
                      {booking.movie_title}
                    </h2>
                    <div className="flex flex-wrap items-center gap-2 mb-6">
                      <span className="px-2.5 py-1 bg-neutral-800 border border-neutral-700 rounded text-[10px] font-bold text-neutral-300 uppercase">
                        {booking.format}
                      </span>
                      <span className="px-2.5 py-1 bg-neutral-800 border border-neutral-700 rounded text-[10px] font-bold text-neutral-300 uppercase">
                        {booking.language}
                      </span>
                      <span className="px-2.5 py-1 bg-neutral-800 border border-neutral-700 rounded text-[10px] font-bold text-neutral-300 uppercase">
                        {booking.screen_type}
                      </span>
                    </div>

                    <div className="grid grid-cols-2 gap-y-5 gap-x-4">
                      <div>
                        <p className="text-[10px] font-bold uppercase tracking-widest text-neutral-500 mb-1">
                          Date
                        </p>
                        <p className="font-semibold text-white">{booking.show_date}</p>
                      </div>
                      <div>
                        <p className="text-[10px] font-bold uppercase tracking-widest text-neutral-500 mb-1">
                          Time
                        </p>
                        <p className="font-semibold text-white">{booking.start_time}</p>
                      </div>
                      <div>
                        <p className="text-[10px] font-bold uppercase tracking-widest text-neutral-500 mb-1">
                          Cinema
                        </p>
                        <p className="font-semibold text-white">{booking.theater_name}</p>
                        <p className="text-xs text-neutral-500">{booking.screen_name}</p>
                      </div>
                      <div>
                        <p className="text-[10px] font-bold uppercase tracking-widest text-neutral-500 mb-1">
                          Seats
                        </p>
                        <p className="font-semibold text-white">{seatLabels}</p>
                        <p className="text-xs text-neutral-500">
                          {booking.seats[0]?.seat_type ?? "Standard"}
                        </p>
                      </div>
                    </div>
                  </div>

                  {/* Booking ID & Total */}
                  <div className="mt-8 pt-6 border-t border-neutral-800/60 flex flex-wrap items-end justify-between gap-4">
                    <div>
                      <p className="text-[10px] font-bold uppercase tracking-widest text-neutral-500 mb-1">
                        Booking ID
                      </p>
                      <p className="font-mono text-white tracking-widest font-bold">
                        {booking.booking_ref}
                      </p>
                    </div>
                    <div className="text-right">
                      <p className="text-[10px] font-bold uppercase tracking-widest text-neutral-500 mb-1">
                        Total Paid
                      </p>
                      <p className="text-2xl md:text-3xl font-black text-primary">
                        ${booking.total_amount.toFixed(2)}
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* ─── Action Buttons ─── */}
            <div className="flex flex-col sm:flex-row gap-4">
              <button
                onClick={() => window.print()}
                className="flex-1 flex items-center justify-center gap-2 bg-primary text-white font-bold py-4 rounded-xl hover:scale-[1.02] active:scale-95 transition-all shadow-lg shadow-primary/25"
              >
                <span className="material-symbols-outlined">download</span>
                Download Ticket
              </button>
              <Link
                href="/movies"
                className="flex-1 flex items-center justify-center gap-2 bg-white/5 backdrop-blur-md border border-white/10 text-white font-bold py-4 rounded-xl hover:bg-white/10 active:scale-95 transition-all"
              >
                <span className="material-symbols-outlined">home</span>
                Back to Home
              </Link>
            </div>
          </div>

          {/* ─── QR Code / Side Panel ─── */}
          <div className="lg:col-span-4">
            <div className="bg-neutral-900/80 backdrop-blur-xl rounded-2xl p-8 border border-neutral-800 flex flex-col items-center text-center shadow-xl">
              <p className="text-[10px] font-bold uppercase tracking-widest text-neutral-500 mb-6">
                Scan at Entry
              </p>

              {/* QR Code placeholder */}
              <div className="bg-white p-5 rounded-2xl mb-8 shadow-inner">
                <div className="w-44 h-44 bg-white flex items-center justify-center relative">
                  <QRCode
                    value={booking.qr_value || booking.ticket_code || booking.booking_ref}
                    size={176}
                    style={{ height: "auto", maxWidth: "100%", width: "100%" }}
                  />
                  <div className="absolute inset-0 border-2 border-primary/15 rounded-lg" />
                </div>
              </div>

              {/* Info notices */}
              <div className="space-y-3 w-full">
                <div className="flex items-center gap-3 p-4 bg-neutral-800/50 rounded-lg text-left">
                  <span className="material-symbols-outlined text-primary shrink-0">
                    info
                  </span>
                  <p className="text-sm text-neutral-400 leading-tight">
                    Please arrive at least 15 minutes before the showtime.
                  </p>
                </div>
                <div className="flex items-center gap-3 p-4 bg-neutral-800/50 rounded-lg text-left">
                  <span className="material-symbols-outlined text-primary shrink-0">
                    no_photography
                  </span>
                  <p className="text-sm text-neutral-400 leading-tight">
                    Photography and recording are strictly prohibited inside.
                  </p>
                </div>
              </div>

              {/* Share button */}
              <div className="mt-8 pt-6 border-t border-neutral-800/60 w-full">
                <button className="text-primary font-bold flex items-center justify-center gap-2 w-full hover:underline decoration-2 underline-offset-4 transition-all">
                  <span className="material-symbols-outlined">share</span>
                  Share Ticket
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <section className="ticket-print">
        <div className="ticket-print-card">
          <div className="ticket-print-qr-wrap">
            <QRCode
              value={booking.qr_value || booking.ticket_code || booking.booking_ref}
              size={220}
              style={{ height: "auto", maxWidth: "100%", width: "100%" }}
            />
          </div>

          <div className="ticket-print-details">
            <h1 className="ticket-print-movie">{booking.movie_title}</h1>
            <p><strong>Date:</strong> {booking.show_date}</p>
            <p><strong>Time:</strong> {booking.start_time}</p>
            <p><strong>Seats:</strong> {seatLabels}</p>
          </div>

          <p className="ticket-print-booking-id">Booking ID: {booking.booking_ref}</p>
        </div>
      </section>

      <style jsx global>{`
        .ticket-print {
          display: none;
        }

        @media print {
          @page {
            size: auto;
            margin: 12mm;
          }

          .screen-content {
            display: none !important;
          }

          .ticket-print {
            display: block !important;
            width: 100%;
          }

          .ticket-print-card {
            margin: 0 auto;
            width: 100%;
            max-width: 420px;
            background: #ffffff;
            color: #111111;
            border: 1px solid #e5e7eb;
            border-radius: 12px;
            padding: 20px;
            box-shadow: none;
            break-inside: avoid;
          }

          .ticket-print-qr-wrap {
            margin: 0 auto 20px;
            width: 240px;
            background: #ffffff;
            padding: 10px;
            border: 1px solid #e5e7eb;
            border-radius: 8px;
          }

          .ticket-print-details {
            margin-bottom: 18px;
            line-height: 1.6;
            font-size: 14px;
          }

          .ticket-print-movie {
            margin: 0 0 10px;
            font-size: 20px;
            font-weight: 700;
            line-height: 1.3;
          }

          .ticket-print-booking-id {
            margin: 0;
            font-size: 11px;
            color: #4b5563;
            text-align: left;
          }
        }
      `}</style>
    </main>
  )
}
