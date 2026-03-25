"use client";

import React, { useMemo, useState, useEffect } from 'react';
import { fetchShowtimeSeats } from '@/lib/api';
import { SeatAPI } from '@/lib/types';

export interface Seat {
  id: string; // e.g. "B5"
  dbId: number;
  row: string;
  number: number;
  type: 'standard' | 'premium' | 'vip';
  status: 'available' | 'sold';
  price: number; 
}

export interface SeatSelectionProps {
  showtimeId: number;
  movie: { title: string; time: string; format: string; hall: string };
  onProceed: (seats: Seat[]) => void;
  onBack: () => void;
  onChangeMovie: () => void;
}

const seatTypeMap = (seatType: string): 'standard' | 'premium' | 'vip' => {
  const normalized = seatType.toLowerCase();
  if (normalized.includes('vip')) return 'vip';
  if (normalized.includes('premium')) return 'premium';
  return 'standard';
};

const mapApiSeat = (seat: SeatAPI): Seat => ({
  id: `${seat.row_label}${seat.col_number}`,
  dbId: seat.id,
  row: seat.row_label,
  number: seat.col_number,
  type: seatTypeMap(seat.seat_type),
  status: seat.status === 'AVAILABLE' ? 'available' : 'sold',
  price: seat.price,
});

const sectionConfig: Record<Seat['type'], { label: string; color: string }> = {
  standard: { label: 'PREMIUM · FRONT',         color: 'text-gray-400' },
  premium:  { label: 'EXECUTIVE · BEST VIEW',   color: 'text-primary' },
  vip:      { label: 'VIP · RECLINER',           color: 'text-purple-400' },
};

export function SeatSelection({ showtimeId, movie, onProceed, onBack, onChangeMovie }: SeatSelectionProps) {
  const [seats, setSeats] = useState<Seat[]>([]);
  const [totalCols, setTotalCols] = useState<number>(0);
  const [loadingSeats, setLoadingSeats] = useState(true);
  const [seatLoadError, setSeatLoadError] = useState<string | null>(null);
  const [selectedSeatIds, setSelectedSeatIds] = useState<string[]>([]);
  const [toastMessage, setToastMessage] = useState<string | null>(null);

  useEffect(() => {
    let active = true;

    const loadSeats = async () => {
      if (!showtimeId) {
        setSeatLoadError('Invalid showtime selected. Please go back and pick a showtime.');
        setLoadingSeats(false);
        return;
      }

      setLoadingSeats(true);
      setSeatLoadError(null);

      try {
        const data = await fetchShowtimeSeats(showtimeId);
        if (!active) return;
        setSeats(data.seats.map((seat) => mapApiSeat(seat)));
        setTotalCols(data.layout?.total_cols ?? 0);
      } catch (error) {
        if (!active) return;
        const message = error instanceof Error ? error.message : 'Failed to load seats';
        setSeatLoadError(message);
      } finally {
        if (active) setLoadingSeats(false);
      }
    };

    loadSeats();

    return () => {
      active = false;
    };
  }, [showtimeId]);

  const selectedSeats = seats.filter(s => selectedSeatIds.includes(s.id));
  const totalPrice = selectedSeats.reduce((sum, seat) => sum + seat.price, 0);

  const rows = useMemo(() => {
    const grouped = new Map<string, Seat[]>();
    for (const seat of seats) {
      if (!grouped.has(seat.row)) grouped.set(seat.row, []);
      grouped.get(seat.row)!.push(seat);
    }
    return Array.from(grouped.entries())
      .map(([rowLabel, rowSeats]) => ({
        rowLabel,
        seats: rowSeats.sort((a, b) => a.number - b.number),
        type: (rowSeats[0]?.type ?? 'standard') as Seat['type'],
      }))
      .sort((a, b) => a.rowLabel.localeCompare(b.rowLabel));
  }, [seats]);

  // Group consecutive rows of the same type into sections
  const sections = useMemo(() => {
    const result: { type: Seat['type']; rows: typeof rows }[] = [];
    for (const row of rows) {
      const last = result[result.length - 1];
      if (last && last.type === row.type) {
        last.rows.push(row);
      } else {
        result.push({ type: row.type, rows: [row] });
      }
    }
    return result;
  }, [rows]);

  const handleSeatClick = (seat: Seat) => {
    if (seat.status === 'sold') return;

    setSelectedSeatIds(prev => {
      // Deselect
      if (prev.includes(seat.id)) {
        setToastMessage(null);
        return prev.filter(id => id !== seat.id);
      }
      
      // Select
      if (prev.length >= 20) {
        setToastMessage("Maximum 20 seats can be selected at once.");
        // Clear toast after 3s (using a timeout here is safe enough for UI)
        setTimeout(() => setToastMessage(null), 3000);
        return prev;
      }
      
      return [...prev, seat.id];
    });
  };

  // UI mapping helpers
  const renderSeatButton = (seat: Seat, sizeClass: string, colorClasses: string) => {
    const isSelected = selectedSeatIds.includes(seat.id);
    const isSold = seat.status === 'sold';
    
    if (isSold) {
      return (
        <button key={seat.id} className={`${sizeClass} rounded-md bg-seat-booked border-none cursor-not-allowed opacity-50 flex items-center justify-center`} aria-label={`${seat.id} sold`} disabled>
          <span className="material-symbols-outlined text-white/20 text-[14px]">close</span>
        </button>
      );
    }
    
    if (isSelected) {
      return (
        <button
          key={seat.id}
          onClick={() => handleSeatClick(seat)}
          className={`${sizeClass} rounded-md bg-seat-selected border-none shadow-[0_0_15px_rgba(34,197,94,0.4)] transition-all flex items-center justify-center text-white font-bold text-xs transform scale-105`}
        >
          {seat.id}
        </button>
      );
    }

    // Available
    return (
      <button
        key={seat.id}
        onClick={() => handleSeatClick(seat)}
        aria-label={`${seat.id} available`}
        className={`${sizeClass} rounded-md border ${colorClasses} hover:shadow-[0_0_10px_currentColor] transition-all bg-transparent group relative`}
      >
        <span className="hidden group-hover:flex absolute -top-8 left-1/2 -translate-x-1/2 bg-gray-800 text-white text-[10px] px-2 py-1 rounded whitespace-nowrap z-10">
          ${seat.price.toFixed(2)}
        </span>
      </button>
    );
  };

  const styleMap: Record<Seat['type'], { sizeClass: string; colorClasses: string }> = {
    standard: {
      sizeClass: 'w-8 h-8 md:w-9 md:h-9',
      colorClasses: 'border-gray-500 hover:border-primary',
    },
    premium: {
      sizeClass: 'w-8 h-8 md:w-9 md:h-9',
      colorClasses: 'border-yellow-500/60 bg-yellow-500/5 hover:border-yellow-400',
    },
    vip: {
      sizeClass: 'w-14 h-12',
      colorClasses: 'border-purple-500/60 bg-purple-500/5 hover:border-purple-400',
    },
  };

  const renderRow = (rowLabel: string, rowSeats: Seat[]) => {
    const n = totalCols > 0 ? totalCols : rowSeats.length;
    // Each of the first two blocks gets ceil(n/3) seats; the last gets the remainder.
    // For n=11: split1=4, split2=8  →  4 | 4 | 3  (cols 1-4 | 5-8 | 9-11)
    const blockSize = Math.ceil(n / 3);
    const split1 = blockSize;
    const split2 = blockSize * 2;

    const leftSeats   = rowSeats.filter(s => s.number <= split1);
    const centerSeats = rowSeats.filter(s => s.number > split1 && s.number <= split2);
    const rightSeats  = rowSeats.filter(s => s.number > split2);

    const renderBlock = (blockSeats: Seat[]) => (
      <div className="flex gap-1.5 md:gap-2">
        {blockSeats.map(seat => renderSeatButton(seat, styleMap[seat.type].sizeClass, styleMap[seat.type].colorClasses))}
      </div>
    );

    return (
      <div key={rowLabel} className="flex items-center gap-0">
        {/* Row label - left */}
        <span className="w-7 shrink-0 text-xs text-gray-500 font-mono text-right pr-2">{rowLabel}</span>
        {/* Left block */}
        {renderBlock(leftSeats)}
        {/* Aisle 1 */}
        <div className="w-6 md:w-10 shrink-0" />
        {/* Center block */}
        {renderBlock(centerSeats)}
        {/* Aisle 2 */}
        <div className="w-6 md:w-10 shrink-0" />
        {/* Right block */}
        {renderBlock(rightSeats)}
        {/* Row label - right */}
        <span className="w-7 shrink-0 text-xs text-gray-500 font-mono text-left pl-2">{rowLabel}</span>
      </div>
    );
  };

  return (
    <div className="bg-[#211111] font-display text-white min-h-screen flex flex-col antialiased selection:bg-primary selection:text-white pb-32">
      {/* Toast Warning */}
      {toastMessage && (
        <div className="fixed top-20 left-1/2 -translate-x-1/2 z-50 bg-red-500 text-white px-6 py-3 rounded-full shadow-lg font-bold text-sm animate-fade-in-up">
          {toastMessage}
        </div>
      )}

      {/* Top Navigation */}
      <header className="sticky top-0 z-40 w-full border-b border-[#382929] bg-[#181111]/95 backdrop-blur">
        <div className="px-6 md:px-10 py-4 flex items-center justify-between max-w-[1400px] mx-auto">
          <div className="flex items-center gap-6">
            <button onClick={onBack} className="flex items-center justify-center p-2 rounded-lg hover:bg-[#382929] transition-colors">
              <span className="material-symbols-outlined text-gray-300">arrow_back</span>
            </button>
            <div className="flex flex-col">
              <h1 className="text-xl font-bold leading-tight tracking-tight">{movie.title}</h1>
              <div className="flex items-center gap-2 text-sm text-[#b89d9f]">
                <span className="material-symbols-outlined text-[16px]">schedule</span>
                <span>{movie.time}</span>
                <span>•</span>
                <span>{movie.format}</span>
                <span>•</span>
                <span>{movie.hall}</span>
              </div>
            </div>
          </div>
          
          {/* Steps */}
          <div className="hidden md:flex items-center gap-1">
            <div className="flex items-center gap-2 px-3 py-1 opacity-50">
              <div className="w-6 h-6 rounded-full bg-[#382929] flex items-center justify-center text-xs font-bold">1</div>
              <span className="text-sm font-medium">Showtime</span>
            </div>
            <div className="w-8 h-[1px] bg-[#382929]"></div>
            <div className="flex items-center gap-2 px-3 py-1">
              <div className="w-6 h-6 rounded-full bg-primary text-white flex items-center justify-center text-xs font-bold">2</div>
              <span className="text-sm font-bold text-primary">Seats</span>
            </div>
            <div className="w-8 h-[1px] bg-[#382929]"></div>
            <div className="flex items-center gap-2 px-3 py-1 opacity-50">
              <div className="w-6 h-6 rounded-full bg-[#382929] flex items-center justify-center text-xs font-bold">3</div>
              <span className="text-sm font-medium">Payment</span>
            </div>
          </div>

          <div className="flex items-center gap-4">
            <button onClick={onChangeMovie} className="hidden md:flex min-w-[84px] cursor-pointer items-center justify-center overflow-hidden rounded-lg h-9 px-4 bg-[#382929] text-white text-sm font-bold leading-normal border border-transparent hover:border-gray-600 transition-all">
              Change Movie
            </button>
          </div>
        </div>
      </header>

      {/* Main Content Area */}
      <main className="flex-grow flex flex-col items-center justify-center relative w-full pt-10 px-4">
        
        {/* Legend */}
        <div className="flex flex-wrap justify-center gap-6 mb-12">
          <div className="flex items-center gap-3">
            <div className="w-5 h-5 rounded border border-gray-500 bg-transparent"></div>
            <span className="text-sm text-gray-300">Available</span>
          </div>
          <div className="flex items-center gap-3">
            <div className="w-5 h-5 rounded bg-seat-selected shadow-[0_0_10px_rgba(34,197,94,0.3)]"></div>
            <span className="text-sm text-gray-300">Selected</span>
          </div>
          <div className="flex items-center gap-3">
            <div className="w-5 h-5 rounded bg-seat-booked border border-transparent opacity-60"></div>
            <span className="text-sm text-gray-300">Sold</span>
          </div>
          <div className="flex items-center gap-3 ml-4 pl-4 border-l border-[#382929]">
            <span className="w-5 h-5 rounded border border-yellow-500/50 bg-yellow-500/10"></span>
            <span className="text-sm text-gray-300">Premium ($18)</span>
          </div>
          <div className="flex items-center gap-3">
            <span className="w-5 h-5 rounded border border-purple-500/50 bg-purple-500/10"></span>
            <span className="text-sm text-gray-300">VIP ($25)</span>
          </div>
        </div>

        {/* Cinema Hall */}
        <div className="relative w-full max-w-5xl mx-auto flex flex-col items-center" style={{ perspective: '1000px' }}>
          
          {/* Screen Curve */}
          <div className="w-3/4 h-16 mb-16 relative flex justify-center group">
            <div 
              className="absolute inset-0 bg-gradient-to-b from-white/10 to-transparent rounded-[50%] blur-sm"
              style={{ boxShadow: '0 15px 40px rgba(255, 255, 255, 0.1)', transform: 'rotateX(-30deg) scale(0.9)', transformOrigin: 'top' }}
            ></div>
            <div className="absolute top-4 bg-[#211111] px-4 py-1 rounded-full border border-white/10 shadow-lg z-10">
              <span className="text-xs uppercase tracking-[0.2em] text-gray-500 font-semibold">Screen this way</span>
            </div>
            <div className="absolute -top-10 left-1/2 -translate-x-1/2 w-1/2 h-20 bg-primary/20 blur-[60px] rounded-full pointer-events-none"></div>
          </div>

          <div className="flex flex-col gap-4 w-full overflow-x-auto pb-10 px-4" style={{ msOverflowStyle: 'none', scrollbarWidth: 'none' }}>
            {loadingSeats ? (
              <div className="text-gray-400 text-sm text-center py-12">Loading seats...</div>
            ) : seatLoadError ? (
              <div className="text-red-400 text-sm text-center py-12">{seatLoadError}</div>
            ) : rows.length === 0 ? (
              <div className="text-gray-400 text-sm text-center py-12">No seats available for this showtime.</div>
            ) : (
              <div className="flex flex-col gap-8 items-center min-w-[600px]">
                {sections.map((section, sIdx) => {
                  const cfg = sectionConfig[section.type];
                  return (
                    <div key={sIdx} className="flex flex-col items-center gap-3 w-full">
                      {/* Section label */}
                      <span className={`text-[11px] font-bold uppercase tracking-[0.18em] ${cfg.color}`}>
                        {cfg.label}
                      </span>
                      {/* Rows in this section */}
                      <div className="flex flex-col gap-2 items-center w-full">
                        {section.rows.map(row => renderRow(row.rowLabel, row.seats))}
                      </div>
                    </div>
                  );
                })}
              </div>
            )}
          </div>
        </div>
      </main>

      {/* Bottom Bar */}
      <div className="fixed bottom-0 left-0 w-full p-4 md:p-6 z-50 pointer-events-none">
        <div className="max-w-[1200px] mx-auto bg-[#2d1b1b]/95 backdrop-blur-md shadow-2xl rounded-xl p-4 md:px-8 border border-[#382929] flex flex-col md:flex-row items-center justify-between gap-4 pointer-events-auto ring-1 ring-white/10">
          <div className="flex items-center gap-6 w-full md:w-auto justify-between md:justify-start">
            <div className="flex flex-col">
              <span className="text-xs font-semibold uppercase text-gray-400 tracking-wider">Seats</span>
              <div className="flex items-center gap-2 mt-1 min-h-[28px]">
                {selectedSeatIds.length > 0 ? (
                  <>
                    <span className="text-xl font-bold text-white truncate max-w-[200px]">
                      {selectedSeatIds.join(', ')}
                    </span>
                    <span className="text-sm bg-[#382929] px-2 py-0.5 rounded text-gray-300 font-medium whitespace-nowrap">
                      {selectedSeatIds.length} ticket{selectedSeatIds.length !== 1 && 's'}
                    </span>
                  </>
                ) : (
                  <span className="text-gray-500 font-medium">None selected</span>
                )}
              </div>
            </div>
            <div className="h-10 w-[1px] bg-[#382929] hidden md:block"></div>
            <div className="flex flex-col text-right md:text-left">
              <span className="text-xs font-semibold uppercase text-gray-400 tracking-wider">Total Price</span>
              <div className="flex items-end gap-1 mt-1 justify-end md:justify-start min-h-[28px]">
                <span className="text-2xl font-bold text-primary">${totalPrice.toFixed(2)}</span>
              </div>
            </div>
          </div>
          
          <div className="w-full md:w-auto flex gap-3">
            <button 
              onClick={() => onProceed(selectedSeats)}
              disabled={selectedSeatIds.length === 0}
              className={`w-full md:w-auto md:min-w-[200px] h-12 rounded-lg font-bold text-lg transition-all flex items-center justify-center gap-2
                ${selectedSeatIds.length > 0 
                  ? 'bg-primary hover:bg-red-600 text-white shadow-[0_4px_14px_0_rgba(234,42,51,0.39)] transform hover:scale-[1.02] active:scale-[0.98]' 
                  : 'bg-[#382929] text-gray-500 cursor-not-allowed'
                }`}
            >
              <span>Proceed to Payment</span>
              <span className="material-symbols-outlined text-[20px]">arrow_forward</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
