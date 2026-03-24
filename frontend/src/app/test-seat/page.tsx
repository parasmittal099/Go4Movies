"use client";

import { SeatSelection } from "@/components/booking/seat-selection";

export default function TestSeatPage() {
  const movie = {
    title: "Dune: Part Two",
    time: "Today, 19:30",
    format: "IMAX 2D",
    hall: "Hall 4"
  };

  return (
    <SeatSelection
      movie={movie}
      onProceed={(seats) => alert(`Proceeding to payment with seats: ${seats.map(s => s.id).join(', ')}`)}
      onBack={() => alert("Back button clicked")}
      onChangeMovie={() => alert("Change Movie clicked")}
    />
  );
}
