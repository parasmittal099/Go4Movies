"use client"

import { useMemo, useState } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { useAuth } from "@/context/auth-context"

export default function PaymentPage() {
  const router = useRouter()
  const searchParams = useSearchParams()
  const { user } = useAuth()

  const showtimeId = searchParams.get("showtimeId") ?? ""
  const selectedSeats = useMemo(
    () => (searchParams.get("seats") ?? "").split(",").filter(Boolean),
    [searchParams]
  )
  const selectedSeatDbIds = useMemo(
    () =>
      (searchParams.get("seatIds") ?? "")
        .split(",")
        .filter(Boolean)
        .map((id) => Number(id))
        .filter((id) => Number.isFinite(id)),
    [searchParams]
  )

  const [cardName, setCardName] = useState("")
  const [cardNumber, setCardNumber] = useState("")
  const [expiry, setExpiry] = useState("")
  const [cvv, setCvv] = useState("")
  const [saveCard, setSaveCard] = useState(false)
  const [submitting, setSubmitting] = useState(false)
  const [submitError, setSubmitError] = useState<string | null>(null)
  const [successMessage, setSuccessMessage] = useState<string | null>(null)

  const handlePay = async () => {
    setSubmitError(null)
    setSuccessMessage(null)

    const normalizedCard = cardNumber.replace(/\s+/g, "")
    if (!showtimeId || selectedSeatDbIds.length === 0) {
      setSubmitError("Missing showtime or seat details. Please reselect seats and try again.")
      return
    }
    if (!cardName.trim() || normalizedCard.length < 12 || !expiry.trim() || cvv.trim().length < 3) {
      setSubmitError("Please enter valid card details before proceeding.")
      return
    }

    setSubmitting(true)
    try {
      const payload = {
        userId: user?.id ?? null,
        showtimeId: Number(showtimeId),
        seatIds: selectedSeatDbIds,
        seatCodes: selectedSeats,
        paymentMethod: "CARD",
        cardHolderName: cardName.trim(),
        cardLast4: normalizedCard.slice(-4),
        saveCard,
      }

      const response = await fetch("/api/mock-booking", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(payload),
      })

      const data = await response.json().catch(() => null)
      if (!response.ok) {
        const message = data?.error ?? "Failed to submit booking payload"
        throw new Error(message)
      }

      const ref = data?.reservation?.bookingRef ?? "N/A"
      setSuccessMessage(`Payment payload submitted successfully. Mock booking ref: ${ref}`)
    } catch (error) {
      setSubmitError(error instanceof Error ? error.message : "Something went wrong while submitting")
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <main className="min-h-screen bg-background px-4 py-8 text-white lg:px-8">
      <div className="mx-auto grid w-full max-w-7xl gap-8 lg:grid-cols-[1fr_420px]">
        <section className="space-y-6">
          <div className="space-y-1">
            <h1 className="text-3xl font-bold tracking-tight">Secure Payment Checkout</h1>
            <p className="text-neutral-500">Pay with card to complete your booking.</p>
          </div>

          <div className="flex items-center gap-2 border-b border-neutral-800 pb-3 text-sm font-semibold text-white">
            <span className="material-symbols-outlined align-middle text-base">credit_card</span>
            Card Payment
          </div>

          <div className="rounded-2xl border border-neutral-800 bg-surface-dark p-6 shadow-xl">
            <div className="mb-6 flex items-center justify-between">
              <h2 className="text-lg font-semibold">Enter Card Details</h2>
              <div className="flex gap-2 text-[10px] text-neutral-400">
                <span className="rounded bg-neutral-800 px-2 py-1">VISA</span>
                <span className="rounded bg-neutral-800 px-2 py-1">MC</span>
                <span className="rounded bg-neutral-800 px-2 py-1">AMEX</span>
              </div>
            </div>

            <form className="space-y-4">
              <div className="space-y-2">
                <label className="text-sm text-neutral-500" htmlFor="cardName">
                  Name on Card
                </label>
                <input
                  id="cardName"
                  type="text"
                  placeholder="John Doe"
                  value={cardName}
                  onChange={(e) => setCardName(e.target.value)}
                  className="w-full rounded-lg border border-neutral-800 bg-neutral-900 px-4 py-3 placeholder:text-neutral-500 focus:border-primary focus:outline-none"
                />
              </div>

              <div className="space-y-2">
                <label className="text-sm text-neutral-500" htmlFor="cardNumber">
                  Card Number
                </label>
                <input
                  id="cardNumber"
                  type="text"
                  placeholder="0000 0000 0000 0000"
                  value={cardNumber}
                  onChange={(e) => setCardNumber(e.target.value)}
                  className="w-full rounded-lg border border-neutral-800 bg-neutral-900 px-4 py-3 font-mono placeholder:text-neutral-500 focus:border-primary focus:outline-none"
                />
              </div>

              <div className="grid gap-4 sm:grid-cols-2">
                <div className="space-y-2">
                  <label className="text-sm text-neutral-500" htmlFor="expiry">
                    Expiry Date
                  </label>
                  <input
                    id="expiry"
                    type="text"
                    placeholder="MM / YY"
                    value={expiry}
                    onChange={(e) => setExpiry(e.target.value)}
                    className="w-full rounded-lg border border-neutral-800 bg-neutral-900 px-4 py-3 font-mono placeholder:text-neutral-500 focus:border-primary focus:outline-none"
                  />
                </div>
                <div className="space-y-2">
                  <label className="text-sm text-neutral-500" htmlFor="cvv">
                    CVV
                  </label>
                  <input
                    id="cvv"
                    type="password"
                    placeholder="***"
                    value={cvv}
                    onChange={(e) => setCvv(e.target.value)}
                    className="w-full rounded-lg border border-neutral-800 bg-neutral-900 px-4 py-3 font-mono placeholder:text-neutral-500 focus:border-primary focus:outline-none"
                  />
                </div>
              </div>

              <label className="flex items-center gap-2 text-sm text-neutral-500">
                <input
                  type="checkbox"
                  checked={saveCard}
                  onChange={(e) => setSaveCard(e.target.checked)}
                  className="h-4 w-4 rounded border-neutral-700 bg-neutral-900"
                />
                Save this card securely for future payments
              </label>

              <button
                type="button"
                onClick={handlePay}
                disabled={submitting}
                className="mt-2 w-full rounded-xl bg-primary px-6 py-3.5 font-bold text-white transition hover:brightness-110"
              >
                {submitting ? "Submitting..." : "Pay $36.40 Securely"}
              </button>

              {submitError && <p className="text-sm text-red-400">{submitError}</p>}
              {successMessage && <p className="text-sm text-green-400">{successMessage}</p>}

              <p className="text-center text-xs text-neutral-500">
                This build posts to a frontend mock endpoint only. Backend booking API is not called.
              </p>

              <p className="text-center text-xs text-neutral-500">
                Your transaction is secured with 256-bit SSL encryption.
              </p>
            </form>
          </div>
        </section>

        <aside className="space-y-4 lg:sticky lg:top-24 lg:self-start">
          <div className="flex items-center justify-between rounded-xl border border-primary/20 bg-primary/10 px-4 py-3 text-primary">
            <div className="flex items-center gap-2 text-sm font-semibold">
              <span className="material-symbols-outlined text-base">timer</span>
              Complete your payment in
            </div>
            <span className="font-mono text-lg font-bold">09:59</span>
          </div>

          <div className="overflow-hidden rounded-2xl border border-neutral-800 bg-surface-dark shadow-2xl">
            <div className="h-44 bg-[url('https://images.unsplash.com/photo-1536440136628-849c177e76a1?q=80&w=1600&auto=format&fit=crop')] bg-cover bg-center" />

            <div className="space-y-5 p-6">
              <div>
                <h3 className="text-2xl font-bold">Dune: Part Two</h3>
                <p className="text-sm text-neutral-500">PG-13 · 2h 46m</p>
              </div>

              <div className="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <p className="text-xs uppercase tracking-wide text-neutral-500">Theater</p>
                  <p className="font-medium">Grand Cinema</p>
                  <p className="text-xs text-neutral-500">Hall 4, IMAX 3D</p>
                </div>
                <div>
                  <p className="text-xs uppercase tracking-wide text-neutral-500">Date & Time</p>
                  <p className="font-medium">Sat, 24 Feb</p>
                  <p className="text-xs text-neutral-500">19:30 PM</p>
                </div>
              </div>

              <div>
                <p className="text-xs uppercase tracking-wide text-neutral-500">Selected Seats</p>
                <div className="mt-1 flex items-center justify-between text-sm">
                  <p className="font-medium">
                    {selectedSeats.length > 0 ? selectedSeats.join(", ") : "No seats selected"}
                  </p>
                  <button
                    type="button"
                    onClick={() => router.back()}
                    className="text-primary hover:underline"
                  >
                    Change
                  </button>
                </div>
              </div>

              <div>
                <p className="text-xs uppercase tracking-wide text-neutral-500">Showtime ID</p>
                <p className="mt-1 text-sm font-medium">{showtimeId || "N/A"}</p>
              </div>

              <div>
                <p className="text-xs uppercase tracking-wide text-neutral-500">Seat IDs (Backend)</p>
                <p className="mt-1 break-all text-xs text-neutral-500">
                  {selectedSeatDbIds.length > 0 ? selectedSeatDbIds.join(", ") : "N/A"}
                </p>
              </div>

              <div className="space-y-2 border-t border-neutral-800 pt-4 text-sm">
                <div className="flex justify-between text-neutral-500">
                  <span>Ticket Price (2x)</span>
                  <span>$32.00</span>
                </div>
                <div className="flex justify-between text-neutral-500">
                  <span>Convenience Fee</span>
                  <span>$2.50</span>
                </div>
                <div className="flex justify-between text-neutral-500">
                  <span>Tax</span>
                  <span>$1.90</span>
                </div>
                <div className="mt-3 flex justify-between text-lg font-bold">
                  <span>Total Amount</span>
                  <span className="text-primary">$36.40</span>
                </div>
              </div>

              <div className="flex gap-2">
                <input
                  type="text"
                  placeholder="Have a promo code?"
                  className="w-full rounded-lg border border-neutral-800 bg-neutral-900 px-3 py-2 text-sm placeholder:text-neutral-500 focus:border-primary focus:outline-none"
                />
                <button type="button" className="rounded-lg border border-neutral-700 px-3 py-2 text-xs font-semibold">
                  Apply
                </button>
              </div>
            </div>
          </div>

          <p className="text-center text-sm text-neutral-500">Need help with your booking?</p>
        </aside>
      </div>
    </main>
  )
}
