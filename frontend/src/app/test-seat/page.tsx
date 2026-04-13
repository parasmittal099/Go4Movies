"use client";

import { useRouter } from "next/navigation";
import { SeatSelection } from "@/components/booking/seat-selection";

export default function TestSeatPage() {
  const router = useRouter();

  const movie = {
    title: "Dune: Part Two",
    time: "Today, 19:30",
    format: "IMAX 2D",
    hall: "Hall 4"
  };

  return (
    <SeatSelection
      showtimeId={1}
      movie={movie}
      onProceed={(seats) => {
        const ids = seats.map((s) => s.id).join(",");
        const dbSeatIds = seats.map((s) => s.dbId).join(",");
        router.push(
          `/payment?showtimeId=1&seats=${encodeURIComponent(ids)}&seatIds=${encodeURIComponent(dbSeatIds)}`
        );
      }}
      onBack={() => alert("Back button clicked")}
      onChangeMovie={() => alert("Change Movie clicked")}
    />
  );
}
