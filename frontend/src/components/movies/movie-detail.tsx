"use client"

import { useState, useEffect, useMemo } from "react"
import { Movie, ShowtimesResponse } from "@/lib/types"

interface MovieDetailProps {
  movie: Movie
  showtimes: ShowtimesResponse | null
  showtimesLoading: boolean
  onBookTickets: (showtimeId: number) => void
}

type TimingStatus = "available" | "filling_fast" | "sold_out"

interface Timing {
  time: string
  status: TimingStatus
  showtimeId: number
  priceMultiplier: number
}

interface Screen {
  type: string
  timings: Timing[]
}

interface Cinema {
  id: string
  name: string
  address: string
  screens: Screen[]
}

interface SelectedTiming {
  showtimeId: number
  cinemaId: string
  cinemaName: string
  screenType: string
  time: string
  price: number
}

const BASE_TICKET_PRICE = 15



export function MovieDetail({ movie, showtimes, showtimesLoading, onBookTickets }: MovieDetailProps) {
  const [selectedDate, setSelectedDate] = useState<string>("")
  const [selectedTiming, setSelectedTiming] = useState<SelectedTiming | null>(null)

  const sortedDates = useMemo(() => {
    return [...(showtimes?.dates ?? [])].sort((a, b) => a.localeCompare(b))
  }, [showtimes?.dates])

  // Set initial selected date once showtimes arrive
  useEffect(() => {
    if (sortedDates.length && !selectedDate) {
      setSelectedDate(sortedDates[0])
    }
  }, [sortedDates, selectedDate])

  // Map API theaters → Cinema objects for the selected date
  const cinemas = useMemo<Cinema[]>(() => {
    if (!showtimes?.theaters || !selectedDate) return []
    return showtimes.theaters
      .map((theater) => {
        const filtered = theater.showtimes.filter((st) => st.show_date === selectedDate)
        const screenMap = new Map<string, Timing[]>()
        for (const st of filtered) {
          const key = st.format
          if (!screenMap.has(key)) screenMap.set(key, [])
          screenMap.get(key)!.push({
            time: st.start_time,
            status: "available",
            showtimeId: st.id,
            priceMultiplier: st.price_multiplier,
          })
        }
        const screens: Screen[] = Array.from(screenMap.entries()).map(
          ([type, timings]) => ({ type, timings })
        )
        return {
          id: theater.theater_id.toString(),
          name: theater.name,
          address: theater.address || "",
          screens,
        } as Cinema
      })
      .filter((c) => c.screens.length > 0)
  }, [showtimes, selectedDate])

  const handleTimingSelect = (cinema: Cinema, screen: Screen, timing: Timing) => {
    if (timing.status === "sold_out") return
    setSelectedTiming({
      showtimeId: timing.showtimeId,
      cinemaId: cinema.id,
      cinemaName: cinema.name,
      screenType: screen.type,
      time: timing.time,
      price: Math.round(BASE_TICKET_PRICE * timing.priceMultiplier * 100) / 100,
    })
  }

  const handleBookTickets = () => {
    if (selectedTiming) {
      onBookTickets(selectedTiming.showtimeId)
    }
  }

  // Helper to parse API dates (YYYY-MM-DD)
  const parseDate = (dateStr: string) => {
    const parts = dateStr.split("-")
    const d = new Date(parseInt(parts[0]), parseInt(parts[1]) - 1, parseInt(parts[2]))
    return {
      dayStr: d.getDate(),
      monthStr: d.toLocaleString("en-US", { month: "short" }),
    }
  }

  return (
    <div className="flex-grow w-full max-w-[1400px] mx-auto px-6 py-8 md:px-10 md:py-12 lg:flex lg:gap-16 relative">
      {/* Left Column: Movie Context (Sticky) */}
      <aside className="w-full lg:w-[360px] flex-shrink-0 mb-10 lg:mb-0">
        <div className="lg:sticky lg:top-28 flex flex-col gap-6">
          {/* Poster */}
          <div className="relative group w-full aspect-[2/3] rounded-xl overflow-hidden shadow-2xl shadow-primary/10">
            <div
              className="absolute inset-0 bg-cover bg-center transition-transform duration-700 group-hover:scale-105"
              style={{ backgroundImage: `url(${movie.poster_url || "/dune-poster.jpg"})` }}
            />
            <div className="absolute inset-0 bg-gradient-to-t from-[#211111] via-transparent to-transparent opacity-60"></div>
            <div className="absolute top-4 left-4">
              <span className="inline-flex items-center justify-center px-2.5 py-1 rounded-lg bg-white/10 backdrop-blur-md text-white text-xs font-bold ring-1 ring-white/20">
                IMAX
              </span>
            </div>
          </div>

          {/* Info */}
          <div className="flex flex-col gap-4">
            <h1 className="text-3xl md:text-4xl font-bold tracking-tight text-white leading-tight">
              {movie.title}
            </h1>
            <div className="flex flex-wrap items-center gap-3 text-sm text-neutral-400">
              <span className="px-2 py-0.5 rounded-lg border border-neutral-700 text-neutral-300 text-xs font-medium">
                PG-13
              </span>
              <span>•</span>
              <span>2h 46m</span>
              <span>•</span>
              <span>Sci-Fi, Adventure</span>
            </div>

            <div className="flex items-center gap-2">
              <span className="material-symbols-outlined text-yellow-500 text-xl font-bold">star</span>
              <span className="text-white font-bold text-lg">8.9</span>
              <span className="text-neutral-500 text-sm">/ 10 (45k votes)</span>
            </div>

            <p className="text-neutral-300 text-sm leading-relaxed max-h-40 overflow-y-auto pr-2">
              {movie.description || "Paul Atreides unites with Chani and the Fremen while on a warpath of revenge against the conspirators who destroyed his family."}
            </p>

            {/* Cast - using downloaded images */}
            <div className="pt-2">
              <h4 className="text-xs font-bold text-neutral-500 uppercase tracking-wider mb-3">Cast</h4>
              <div className="flex -space-x-3 overflow-hidden p-1 pl-0">
                <img
                  src="/cast1.jpg"
                  alt="Timothée Chalamet"
                  className="inline-block w-10 h-10 rounded-full ring-2 ring-[#211111] object-cover"
                />
                <img
                  src="/cast2.jpg"
                  alt="Zendaya"
                  className="inline-block w-10 h-10 rounded-full ring-2 ring-[#211111] object-cover"
                />
                <img
                  src="/cast3.jpg"
                  alt="Rebecca Ferguson"
                  className="inline-block w-10 h-10 rounded-full ring-2 ring-[#211111] object-cover"
                />
                <div className="inline-flex w-10 h-10 rounded-full ring-2 ring-[#211111] bg-surface-dark items-center justify-center text-xs font-medium text-white ring-offset-2 ring-offset-[#211111]">
                  +8
                </div>
              </div>
            </div>
          </div>
        </div>
      </aside>

      {/* Right Column: Selection Interface */}
      <section className="flex-1 min-w-0 flex flex-col gap-8 pb-32">
        {/* Date Selector */}
        <div className="w-full">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-bold text-white">Select Date</h3>
            <div className="flex gap-2">
              <button className="w-8 h-8 rounded-full bg-surface-dark hover:bg-neutral-800 flex items-center justify-center text-neutral-400 hover:text-white transition-colors">
                <span className="material-symbols-outlined text-sm">chevron_left</span>
              </button>
              <button className="w-8 h-8 rounded-full bg-surface-dark hover:bg-neutral-800 flex items-center justify-center text-neutral-400 hover:text-white transition-colors">
                <span className="material-symbols-outlined text-sm">chevron_right</span>
              </button>
            </div>
          </div>
          <div className="flex gap-3 overflow-x-auto [&::-webkit-scrollbar]:hidden [-ms-overflow-style:none] [scrollbar-width:none] pb-2">
            {sortedDates.map((date) => {
              const { dayStr, monthStr } = parseDate(date)
              const isActive = selectedDate === date
              const today = new Date().toISOString().split("T")[0]
              const dayLabel = date === today ? "TODAY" : new Date(date + "T00:00:00").toLocaleString("en-US", { weekday: "short" }).toUpperCase()

              return (
                <button
                  key={date}
                  onClick={() => { setSelectedDate(date); setSelectedTiming(null) }}
                  className={`flex flex-col items-center justify-center min-w-[72px] h-[84px] rounded-xl shrink-0 transition-all ${
                    isActive
                      ? "bg-primary text-white shadow-lg shadow-primary/20 transform hover:scale-105"
                      : "bg-surface-dark hover:bg-neutral-800 border border-transparent hover:border-neutral-700 text-neutral-400 hover:text-white group"
                  }`}
                >
                  <span
                    className={`text-xs font-medium uppercase tracking-wide ${
                      isActive ? "opacity-80" : "group-hover:text-primary"
                    }`}
                  >
                    {dayLabel}
                  </span>
                  <span className={`text-2xl font-bold mt-1 ${isActive ? "" : "text-white"}`}>{dayStr}</span>
                  <span className={`text-xs font-medium ${isActive ? "opacity-80" : ""}`}>{monthStr}</span>
                </button>
              )
            })}
          </div>
        </div>

        {/* Filters */}
        <div className="flex flex-wrap items-center gap-3 border-b border-neutral-800 pb-6">
          <button className="flex items-center gap-2 px-4 py-2 bg-neutral-800 rounded-lg text-sm font-medium text-white hover:bg-neutral-700 transition-colors">
            <span>All Formats</span>
            <span className="material-symbols-outlined text-sm">expand_more</span>
          </button>
          <button className="flex items-center gap-2 px-4 py-2 bg-surface-dark border border-neutral-800 rounded-lg text-sm font-medium text-neutral-300 hover:text-white hover:border-neutral-600 transition-colors">
            <span>Language</span>
            <span className="material-symbols-outlined text-sm">expand_more</span>
          </button>
          <button className="flex items-center gap-2 px-4 py-2 bg-surface-dark border border-neutral-800 rounded-lg text-sm font-medium text-neutral-300 hover:text-white hover:border-neutral-600 transition-colors">
            <span>Price Range</span>
            <span className="material-symbols-outlined text-sm">expand_more</span>
          </button>
          <div className="ml-auto text-sm text-neutral-400">
            Showing <span className="text-white font-bold">{cinemas.length}</span> theaters nearby
          </div>
        </div>

        {/* Theaters List */}
        <div className="flex flex-col gap-6">
          {showtimesLoading ? (
            // Skeleton Loader
            [1, 2].map((i) => (
              <div key={i} className="bg-surface-dark rounded-xl p-6 border border-neutral-800 animate-pulse">
                <div className="h-6 w-1/3 bg-neutral-800 rounded-lg mb-2"></div>
                <div className="h-4 w-1/2 bg-neutral-800 rounded-lg mb-6"></div>
                <div className="h-10 w-full bg-neutral-800 rounded-lg"></div>
              </div>
            ))
          ) : cinemas.length === 0 ? (
            <div className="text-neutral-400 text-sm py-8 text-center">No showtimes available for this date.</div>
          ) : (
            cinemas.map((cinema) => (
              <div
                key={cinema.id}
                className="bg-surface-dark rounded-xl p-6 border border-transparent hover:border-neutral-700 transition-all"
              >
                <div className="flex flex-col sm:flex-row sm:items-start justify-between gap-4 mb-6">
                  <div className="flex items-start gap-4">
                    <div className="hidden sm:flex items-center justify-center w-10 h-10 rounded-full bg-neutral-800 text-neutral-500 hover:text-primary cursor-pointer transition-colors">
                      <span className="material-symbols-outlined">theater_comedy</span>
                    </div>
                    <div>
                      <h3 className="text-lg font-bold text-white">{cinema.name}</h3>
                      {cinema.address && (
                        <p className="text-sm text-neutral-500 mt-1">{cinema.address}</p>
                      )}
                    </div>
                  </div>
                </div>

                <div className="space-y-4">
                  {cinema.screens.map((screen, sIdx) => (
                    <div key={screen.type}>
                      {sIdx > 0 && <div className="w-full h-px bg-neutral-800/50 my-4"></div>}
                      <div className="flex flex-col sm:flex-row gap-3 sm:items-start pt-2">
                        <span className="text-xs font-bold text-neutral-400 uppercase w-20 shrink-0 mt-3">
                          {screen.type}
                        </span>
                        <div className="flex flex-wrap gap-3">
                          {screen.timings.map((timing) => {
                            const isSelected =
                              selectedTiming?.time === timing.time &&
                              selectedTiming?.cinemaId === cinema.id &&
                              selectedTiming?.screenType === screen.type

                            if (timing.status === "sold_out") {
                              return (
                                <div key={timing.time} className="flex flex-col gap-1">
                                  <button
                                    disabled
                                    className="px-5 py-2.5 rounded-lg border border-neutral-700 bg-neutral-800 opacity-50 cursor-not-allowed text-sm font-semibold text-white"
                                  >
                                    {timing.time}
                                  </button>
                                  <span className="block text-[10px] text-red-400 font-medium text-center">
                                    Sold Out
                                  </span>
                                </div>
                              )
                            }

                            if (isSelected) {
                              return (
                                <div key={timing.time} className="flex flex-col gap-1">
                                  <button className="group relative px-5 py-2.5 rounded-lg bg-primary border border-primary hover:bg-red-600 transition-all shadow-lg shadow-primary/20">
                                    <span className="text-sm font-semibold text-white">{timing.time}</span>
                                    <div className="absolute -top-2 -right-2 flex w-5 h-5 bg-white rounded-full items-center justify-center text-primary shadow-sm">
                                      <span className="material-symbols-outlined text-[12px] font-bold">check</span>
                                    </div>
                                  </button>
                                  {timing.status === "filling_fast" && (
                                    <span className="block text-[10px] text-green-400 font-medium text-center">
                                      Filling Fast
                                    </span>
                                  )}
                                </div>
                              )
                            }

                            // Available or Filling Fast (but not selected)
                            return (
                              <div key={timing.time} className="flex flex-col gap-1">
                                <button
                                  onClick={() => handleTimingSelect(cinema, screen, timing)}
                                  className="group relative px-5 py-2.5 rounded-lg border border-neutral-700 bg-surface-dark hover:bg-white hover:border-white transition-all"
                                >
                                  <span className="text-sm font-semibold text-white group-hover:text-background transition-colors">
                                    {timing.time}
                                  </span>
                                </button>
                                {timing.status === "filling_fast" && (
                                  <span className="block text-[10px] text-green-400 font-medium text-center">
                                    Filling Fast
                                  </span>
                                )}
                              </div>
                            )
                          })}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            ))
          )}
        </div>
      </section>

      {/* Sticky Bottom Bar */}
      {selectedTiming && (
        <div className="fixed bottom-0 left-0 right-0 bg-neutral-900 border-t border-neutral-800 px-6 py-4 shadow-[0_-4px_20px_rgba(0,0,0,0.6)] z-40 animate-[fade-in-up_0.3s_ease-out]">
          <div className="max-w-[1400px] mx-auto flex flex-col sm:flex-row items-center justify-between gap-4">
            <div className="flex items-center gap-6">
              <div className="hidden sm:block">
                <p className="text-xs text-neutral-400 uppercase tracking-wide mb-1">Selected Show</p>
                <p className="text-white font-bold text-lg">
                  {selectedTiming.screenType} • {selectedTiming.time}
                </p>
              </div>
              <div className="h-10 w-px bg-neutral-800 hidden sm:block"></div>
              <div>
                <p className="text-xs text-neutral-400 uppercase tracking-wide mb-1">Location</p>
                <p className="text-white font-medium">{selectedTiming.cinemaName}</p>
              </div>
              {/* Group booking avatars layout (dummy) */}
              <div className="h-10 w-px bg-neutral-800 hidden lg:block"></div>
              <div className="hidden lg:flex items-center gap-2">
                <p className="text-xs text-neutral-400 uppercase tracking-wide mr-2">Going with</p>
                <div className="flex -space-x-3">
                  <img
                    src="/cast1.jpg"
                    alt="Friend 1"
                    className="w-8 h-8 rounded-full ring-2 ring-neutral-900 object-cover"
                  />
                  <img
                    src="/cast2.jpg"
                    alt="Friend 2"
                    className="w-8 h-8 rounded-full ring-2 ring-neutral-900 object-cover"
                  />
                  <div className="w-8 h-8 rounded-full ring-2 ring-neutral-900 bg-neutral-800 flex items-center justify-center text-[10px] font-bold text-white">
                    +2
                  </div>
                </div>
              </div>
            </div>

            <div className="flex w-full sm:w-auto gap-4 items-center">
              <div className="flex-1 sm:flex-none flex flex-col justify-center sm:text-right px-2">
                <span className="text-xs text-neutral-400">Total Price</span>
                <span className="text-xl font-bold text-white uppercase">${selectedTiming.price.toFixed(2)}</span>
              </div>
              <button
                onClick={handleBookTickets}
                className="flex-1 sm:flex-none bg-primary hover:bg-red-600 text-white font-bold py-3 px-8 rounded-lg shadow-lg shadow-primary/30 transition-all transform hover:scale-[1.02] flex items-center justify-center gap-2"
              >
                <span>Select Seats</span>
                <span className="material-symbols-outlined text-[20px]">arrow_forward</span>
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}
