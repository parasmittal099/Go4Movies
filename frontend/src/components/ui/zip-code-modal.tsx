"use client"

import { useState, useEffect, useMemo } from "react"
import { fetchLocations } from "@/lib/api"
import { Location } from "@/lib/types"
import { useZipCode } from "@/hooks/useZipCode"

interface ZipCodeModalProps {
  onSelectLocation: (location: Location) => void
  onClose: () => void
}

export function ZipCodeModal({ onSelectLocation, onClose }: ZipCodeModalProps) {
  const [input, setInput] = useState("")
  const [locations, setLocations] = useState<Location[]>([])
  const [loading, setLoading] = useState(true)
  const zipCode = useZipCode().getZipCode()

  useEffect(() => {
    fetchLocations()
      .then((data) => {
        setLocations(data)
        setLoading(false)
      })
      .catch((error) => {
        console.error("Error fetching locations:", error)
        setLoading(false)
      })
  }, [])

  // Unique cities (de-duped by city name, keep first occurrence)
  const uniqueCities = useMemo(() => {
    const cityMap = new Map<string, Location>()
    locations.forEach((loc) => {
      if (!cityMap.has(loc.city)) {
        cityMap.set(loc.city, loc)
      }
    })
    return Array.from(cityMap.values())
  }, [locations])

  // First 7 cities as "popular"
  const popularCities = uniqueCities.slice(0, 7)

  // Filter locations by search input (zip or city name)
  const filteredLocations = useMemo(() => {
    if (!input.trim()) return []
    const query = input.toLowerCase()
    return locations.filter(
      (loc) =>
        loc.zipcode.startsWith(input) ||
        loc.city.toLowerCase().includes(query) ||
        loc.state.toLowerCase().includes(query)
    )
  }, [input, locations])

  const handleSelect = (loc: Location) => {
    onSelectLocation(loc)
  }

  const getZipCodeByLocation = async (lat: number, lon: number) => {
  try {
    const response = await fetch(
      `https://nominatim.openstreetmap.org/reverse?format=jsonv2&lat=${lat}&lon=${lon}`
    );
    const data = await response.json();
    
    // The zip code is usually in the address object
    const zipCode = data.address.postcode;
    console.log("zipcode haha:", zipCode);
    return zipCode;
  } catch (error) {
    console.error("Error fetching zip code:", error);
  }
};

  const handleDetectLocation = async () => {
    if ("geolocation" in navigator) {
      navigator.geolocation.getCurrentPosition(async (position) => {
        const detectedZip = await getZipCodeByLocation(position.coords.latitude, position.coords.longitude)
        // Try to find a matching location from our list
        const matched = locations.find((l) => l.zipcode === detectedZip)
        if (matched) {
          onSelectLocation(matched)
        } else {
          // Fallback: build a minimal Location object
          onSelectLocation({ id: 0, zipcode: detectedZip, city: "", state: "", country: "" })
        }
      })
    } else {
      alert("Geolocation is not available")
    }
  }

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm">
      <div className="relative w-full max-w-[640px] bg-surface-dark border border-neutral-800 rounded-lg shadow-2xl flex flex-col overflow-hidden animate-fade-in-up">
        {/* Close button */}
        {zipCode && (
          <button
            onClick={onClose}
            className="absolute top-4 right-4 p-2 text-neutral-400 hover:text-white transition-colors rounded-full hover:bg-neutral-800 z-10"
          >
            <span className="material-symbols-outlined text-xl leading-none">
            close
          </span>
        </button>
        )}
        {/* Header */}
        <div className="pt-10 pb-2 px-8 text-center">
          <div className="inline-flex items-center justify-center w-12 h-12 rounded-full bg-primary/20 text-primary mb-4 ring-1 ring-primary/30">
            <span className="material-symbols-outlined text-2xl">
              location_on
            </span>
          </div>
          <h2 className="text-white text-2xl font-bold leading-tight mb-2">
            Hi there!
          </h2>
          <p className="text-neutral-400 text-base">
            Let&apos;s find movies playing near you.
          </p>
        </div>

        {/* Search & Detect */}
        <div className="px-8 py-6 space-y-4">
          <div className="group relative flex flex-col w-full">
            <div className="relative flex items-center">
              <div className="absolute left-4 text-neutral-400 group-focus-within:text-primary transition-colors pointer-events-none flex items-center">
                <span className="material-symbols-outlined">search</span>
              </div>
              <input
                value={input}
                onChange={(e) => setInput(e.target.value)}
                className="w-full h-12 pl-12 pr-4 rounded-lg bg-neutral-800 border border-neutral-700 text-white placeholder:text-neutral-500 focus:ring-2 focus:ring-primary focus:border-transparent transition-all text-base outline-none"
                placeholder="Search for your city or zip code"
                type="text"
                autoFocus
              />
            </div>

            {/* Search results dropdown */}
            {filteredLocations.length > 0 && input.trim() && (
              <ul className="mt-2 bg-neutral-800 border border-neutral-700 rounded-lg max-h-48 overflow-y-auto divide-y divide-neutral-700/50 shadow-xl">
                {filteredLocations.map((loc) => (
                  <li
                    key={loc.id}
                    onClick={() => handleSelect(loc)}
                    className="px-4 py-3 hover:bg-surface-dark-hover cursor-pointer text-neutral-300 hover:text-white transition-colors flex justify-between items-center"
                  >
                    <div className="flex flex-col">
                      <span className="font-medium">{loc.city}</span>
                      <span className="text-[10px] text-neutral-500 uppercase tracking-wider">
                        {loc.state}
                      </span>
                    </div>
                    <span className="text-neutral-500 text-sm font-mono">
                      {loc.zipcode}
                    </span>
                  </li>
                ))}
              </ul>
            )}
          </div>

          <button
            onClick={handleDetectLocation}
            className="w-full flex items-center justify-center gap-2 h-12 rounded-lg bg-neutral-800/50 border border-primary/40 text-primary hover:bg-primary hover:text-white transition-all duration-300 font-semibold text-sm group"
          >
            <span className="material-symbols-outlined group-hover:animate-pulse text-lg">
              my_location
            </span>
            Detect My Location
          </button>
        </div>

        {/* Popular Cities */}
        {/* <div className="px-8 pb-10">
          <div className="relative flex items-center py-4">
            <div className="flex-grow border-t border-neutral-800"></div>
            <span className="flex-shrink-0 mx-4 text-[10px] font-bold uppercase tracking-[0.2em] text-neutral-500">
              Popular Cities
            </span>
            <div className="flex-grow border-t border-neutral-800"></div>
          </div>

          <div className="grid grid-cols-3 sm:grid-cols-4 gap-3 mt-2">
            {!loading ? (
              popularCities.map((loc) => (
                <button
                  key={loc.id}
                  onClick={() => handleSelect(loc.zipcode)}
                  className="flex flex-col items-center justify-center gap-2 p-3 rounded-lg hover:bg-surface-dark-hover text-neutral-400 hover:text-white transition-all group cursor-pointer border border-transparent hover:border-neutral-700"
                >
                  <div className="bg-neutral-800 rounded-full w-14 h-14 flex items-center justify-center group-hover:bg-neutral-700 group-hover:shadow-[0_0_20px_rgba(234,42,51,0.15)] transition-all ring-1 ring-neutral-700 group-hover:ring-primary/40 overflow-hidden relative">
                    <span className="material-symbols-outlined text-2xl opacity-60 group-hover:opacity-100 transition-opacity">
                      location_city
                    </span>
                  </div>
                  <span className="text-xs font-semibold">{loc.city}</span>
                </button>
              ))
            ) : (
              Array.from({ length: 4 }).map((_, i) => (
                <div
                  key={i}
                  className="flex flex-col items-center gap-2 p-3 animate-pulse"
                >
                  <div className="bg-neutral-800 rounded-full w-14 h-14"></div>
                  <div className="h-2 w-12 bg-neutral-800 rounded"></div>
                </div>
              ))
            )}

            {!loading && (
              <button className="flex flex-col items-center justify-center gap-2 p-3 rounded-lg hover:bg-surface-dark-hover text-primary transition-all group cursor-pointer border border-transparent hover:border-neutral-700">
                <div className="bg-neutral-800 rounded-full w-14 h-14 flex items-center justify-center group-hover:bg-primary/20 group-hover:shadow-md transition-all ring-1 ring-neutral-700 group-hover:ring-primary/60">
                  <span className="material-symbols-outlined text-2xl">
                    apps
                  </span>
                </div>
                <span className="text-xs font-semibold">View All</span>
              </button>
            )}
          </div>
        </div> */}

        {/* Footer */}
        {/* <div className="bg-neutral-900/50 py-4 px-8 text-center border-t border-neutral-800">
          <button
            onClick={onClose}
            className="text-[11px] font-bold uppercase tracking-widest text-neutral-500 hover:text-primary transition-colors"
          >
            Skip for now
          </button>
        </div> */}
      </div>
    </div>
  )
}
