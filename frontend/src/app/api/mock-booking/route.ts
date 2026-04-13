import { NextResponse } from "next/server"

type MockBookingPayload = {
  userId: number | null
  showtimeId: number
  seatIds: number[]
  seatCodes: string[]
  paymentMethod: "CARD"
  cardHolderName: string
  cardLast4: string
  saveCard: boolean
}

export async function POST(request: Request) {
  const body = (await request.json().catch(() => null)) as MockBookingPayload | null
  if (!body) {
    return NextResponse.json({ error: "Invalid JSON payload" }, { status: 400 })
  }

  if (!body.showtimeId || !Array.isArray(body.seatIds) || body.seatIds.length === 0) {
    return NextResponse.json({ error: "showtimeId and seatIds are required" }, { status: 400 })
  }

  if (body.paymentMethod !== "CARD" || !body.cardHolderName || !body.cardLast4) {
    return NextResponse.json({ error: "Valid card details are required" }, { status: 400 })
  }

  const bookingRef = `MOCK-${Date.now().toString().slice(-8)}`

  return NextResponse.json(
    {
      message: "Mock booking payload accepted",
      reservation: {
        bookingRef,
        status: "PENDING",
        paymentStatus: "INITIATED",
        showtimeId: body.showtimeId,
        seatIds: body.seatIds,
        seatCodes: body.seatCodes,
        userId: body.userId,
      },
    },
    { status: 201 }
  )
}
