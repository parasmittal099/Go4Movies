"use client"

import { useState, useEffect, useRef, useMemo } from "react"
import Link from "next/link"
import { useAuth } from "@/context/auth-context"
import { useMovies } from "@/context/movies-context"
import { useRouter } from "next/navigation"

export function Header() {
  const { user, isReady, logout } = useAuth()
  const { movies } = useMovies()
  const router = useRouter()

  const [searchOpen, setSearchOpen] = useState(false)
  const [searchQuery, setSearchQuery] = useState("")
  const [debouncedQuery, setDebouncedQuery] = useState("")
  const searchInputRef = useRef<HTMLInputElement>(null)
  const searchContainerRef = useRef<HTMLDivElement>(null)

  const handleLogout = () => {
    logout()
    router.push("/")
  }

  // Get initials from full_name
  const initials =
    user?.full_name
      ?.split(" ")
      .map((n) => n[0])
      .join("")
      .toUpperCase()
      .slice(0, 2) ?? ""

  /* ---- Debounce the search query (300ms) ---- */
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedQuery(searchQuery)
    }, 300)
    return () => clearTimeout(timer)
  }, [searchQuery])

  /* ---- Filter movies by debounced query ---- */
  const searchResults = useMemo(() => {
    if (!debouncedQuery.trim()) return []
    const q = debouncedQuery.toLowerCase()
    return movies.filter((m) => m.title.toLowerCase().includes(q))
  }, [debouncedQuery, movies])

  /* ---- Focus input when search opens ---- */
  useEffect(() => {
    if (searchOpen) {
      // Small delay so the DOM has rendered the input
      requestAnimationFrame(() => searchInputRef.current?.focus())
    } else {
      setSearchQuery("")
      setDebouncedQuery("")
    }
  }, [searchOpen])

  /* ---- Close search on click outside ---- */
  useEffect(() => {
    if (!searchOpen) return
    const handleClickOutside = (e: MouseEvent) => {
      if (
        searchContainerRef.current &&
        !searchContainerRef.current.contains(e.target as Node)
      ) {
        setSearchOpen(false)
      }
    }
    document.addEventListener("mousedown", handleClickOutside)
    return () => document.removeEventListener("mousedown", handleClickOutside)
  }, [searchOpen])

  /* ---- Close search on Escape ---- */
  useEffect(() => {
    if (!searchOpen) return
    const handleKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") setSearchOpen(false)
    }
    document.addEventListener("keydown", handleKey)
    return () => document.removeEventListener("keydown", handleKey)
  }, [searchOpen])

  const handleResultClick = (movieId: number) => {
    setSearchOpen(false)
    router.push(`/movies/${movieId}`)
  }

  return (
    <header className="flex items-center justify-between border-b border-neutral-800 px-6 md:px-10 py-4 bg-neutral-900">
      <Link href="/movies" className="flex items-center gap-3 text-white">
        <span className="material-symbols-outlined text-primary text-3xl">
          movie
        </span>
        <h1 className="text-xl font-bold tracking-tight">Go4Movies</h1>
      </Link>

      <div className="hidden md:flex flex-1 justify-end gap-8 items-center">
        <div className="flex gap-3 items-center">
          {/* ---- Search ---- */}
          <div ref={searchContainerRef} className="relative">
            {searchOpen ? (
              <div className="flex items-center gap-2">
                <div className="relative">
                  <input
                    ref={searchInputRef}
                    type="text"
                    value={searchQuery}
                    onChange={(e) => setSearchQuery(e.target.value)}
                    placeholder="Search movies..."
                    className="w-64 h-9 pl-9 pr-8 rounded-lg bg-neutral-800 border border-neutral-700 text-white text-sm placeholder:text-neutral-500 outline-none focus:ring-2 focus:ring-primary/50 focus:border-primary/50 transition-all"
                  />
                  <span className="material-symbols-outlined absolute left-2.5 top-1/2 -translate-y-1/2 text-neutral-400 text-lg pointer-events-none">
                    search
                  </span>
                  {searchQuery && (
                    <button
                      onClick={() => setSearchQuery("")}
                      className="absolute right-2 top-1/2 -translate-y-1/2 text-neutral-500 hover:text-neutral-300 transition-colors"
                    >
                      <span className="material-symbols-outlined text-base">
                        close
                      </span>
                    </button>
                  )}
                </div>
              </div>
            ) : (
              <button
                onClick={() => setSearchOpen(true)}
                className="p-2 rounded-lg bg-neutral-800 text-white hover:bg-neutral-700 transition-colors flex align-center"
              >
                <span className="material-symbols-outlined">search</span>
              </button>
            )}

            {/* ---- Results dropdown ---- */}
            {searchOpen && debouncedQuery.trim() && (
              <div className="absolute top-full right-0 mt-2 w-80 bg-neutral-900 border border-neutral-700 rounded-lg shadow-2xl overflow-hidden z-50 animate-fade-in-up">
                {searchResults.length > 0 ? (
                  <ul className="max-h-72 overflow-y-auto divide-y divide-neutral-800">
                    {searchResults.map((movie) => (
                      <li key={movie.id}>
                        <button
                          onClick={() => handleResultClick(movie.id)}
                          className="w-full flex items-center gap-3 px-4 py-3 hover:bg-neutral-800 transition-colors text-left group"
                        >
                          {movie.poster_url ? (
                            <img
                              src={movie.poster_url}
                              alt={movie.title}
                              className="w-10 h-14 rounded object-cover shrink-0 bg-neutral-800"
                            />
                          ) : (
                            <div className="w-10 h-14 rounded bg-neutral-800 flex items-center justify-center shrink-0">
                              <span className="material-symbols-outlined text-neutral-600 text-lg">
                                movie
                              </span>
                            </div>
                          )}
                          <div className="flex flex-col min-w-0">
                            <span className="text-sm font-medium text-neutral-200 group-hover:text-white transition-colors truncate">
                              {movie.title}
                            </span>
                            {movie.description && (
                              <span className="text-xs text-neutral-500 line-clamp-1 mt-0.5">
                                {movie.description}
                              </span>
                            )}
                          </div>
                          <span className="material-symbols-outlined text-neutral-600 group-hover:text-primary text-lg ml-auto shrink-0 transition-colors">
                            arrow_forward
                          </span>
                        </button>
                      </li>
                    ))}
                  </ul>
                ) : (
                  <div className="px-4 py-6 text-center">
                    <span className="material-symbols-outlined text-3xl text-neutral-600 mb-2 block">
                      search_off
                    </span>
                    <p className="text-sm text-neutral-400">
                      No movies found for &ldquo;{debouncedQuery}&rdquo;
                    </p>
                  </div>
                )}
              </div>
            )}
          </div>

          {/* ---- Auth section ---- */}
          {!isReady ? (
            <div className="w-20 h-9 rounded-lg bg-neutral-800 animate-pulse" />
          ) : user ? (
            <div className="flex items-center gap-3">
              <div className="flex items-center gap-2.5 px-3 py-1.5 rounded-lg hover:bg-neutral-800 transition-colors group">
                <div className="flex items-center justify-center w-8 h-8 rounded-full bg-primary text-white text-xs font-bold shrink-0">
                  {initials}
                </div>
                <span className="text-sm font-medium text-neutral-200 group-hover:text-white transition-colors max-w-[120px] truncate">
                  {user.full_name}
                </span>
              </div>
              <button
                onClick={handleLogout}
                className="px-4 py-2 rounded-lg bg-neutral-800 text-neutral-300 text-sm font-medium hover:bg-neutral-700 hover:text-white transition-all"
              >
                Sign Out
              </button>
            </div>
          ) : (
            <Link
              href="/login"
              className="px-5 py-2 rounded-lg bg-primary text-white text-sm font-bold hover:brightness-110 transition-all flex items-center"
            >
              Sign In
            </Link>
          )}
        </div>
      </div>
    </header>
  )
}
