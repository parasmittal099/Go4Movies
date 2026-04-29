"use client"

import { useEffect, useMemo, useState } from "react"
import { useRouter, useSearchParams } from "next/navigation"
import { useAuth } from "@/context/auth-context"
import { previewCheckout, confirmCheckout, fetchShowtimeSeats, type CheckoutPayload } from "@/lib/api"
import type { CheckoutQuote, SeatsResponse } from "@/lib/types"

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

  const [discountCode, setDiscountCode] = useState("")
  const [appliedDiscount, setAppliedDiscount] = useState("")
  const [preview, setPreview] = useState<CheckoutQuote | null>(null)
  const [previewLoading, setPreviewLoading] = useState(false)
  const [previewError, setPreviewError] = useState<string | null>(null)
  const [showtimeInfo, setShowtimeInfo] = useState<SeatsResponse["showtime"] | null>(null)
  const [seatSubtotal, setSeatSubtotal] = useState<number | null>(null)
  const [timerSeconds, setTimerSeconds] = useState(10 * 60)

  const showtimeIdNum = Number(showtimeId)

  useEffect(() => {
    const intervalId = window.setInterval(() => {
      setTimerSeconds((prev) => {
        if (prev <= 1) {
          window.clearInterval(intervalId)
          return 0
        }
        return prev - 1
      })
    }, 1000)

    return () => window.clearInterval(intervalId)
  }, [])

  // Fetch showtime info independently — does not require auth
  useEffect(() => {
    if (!showtimeIdNum) return
    fetchShowtimeSeats(showtimeIdNum)
      .then((seatsResp) => {
        setShowtimeInfo(seatsResp.showtime)

        const selectedSeatSet = new Set(selectedSeatDbIds)
        const subtotal = seatsResp.seats
          .filter((seat) => selectedSeatSet.has(seat.id))
          .reduce((sum, seat) => sum + seat.price, 0)
        setSeatSubtotal(Number(subtotal.toFixed(2)))
      })
      .catch(() => {/* showtime info is non-critical, silently skip */})
  }, [showtimeIdNum, selectedSeatDbIds])

  // Fetch pricing preview — requires user + seats
  useEffect(() => {
    if (!showtimeIdNum || selectedSeatDbIds.length === 0 || !user?.id) {
      return
    }
    setPreviewLoading(true)
    setPreviewError(null)

    const payload: CheckoutPayload = {
      user_id: user.id,
      showtime_id: showtimeIdNum,
      seat_ids: selectedSeatDbIds,
      ...(appliedDiscount ? { discount_code: appliedDiscount } : {}),
    }

    previewCheckout(payload)
      .then((quote) => setPreview(quote))
      .catch((err) => {
        setPreviewError(err instanceof Error ? err.message : "Failed to load checkout preview")
      })
      .finally(() => setPreviewLoading(false))
  }, [showtimeIdNum, appliedDiscount, user?.id])

  const handleApplyDiscount = () => {
    setAppliedDiscount(discountCode.trim())
  }
  const handlePay = async () => {
    setSubmitError(null)
    setSuccessMessage(null)

    const normalizedCard = cardNumber.replace(/[-\s]/g, "")
    if (!showtimeId || selectedSeatDbIds.length === 0) {
      setSubmitError("Missing showtime or seat details. Please reselect seats and try again.")
      return
    }
    if (!user?.id) {
      setSubmitError("You must be logged in to complete a booking.")
      return
    }
    if (!cardName.trim() || normalizedCard.length < 12 || !expiry.trim() || cvv.trim().length < 3) {
      setSubmitError("Please enter valid card details before proceeding.")
      return
    }

    setSubmitting(true)
    try {
      const data = await confirmCheckout({
        user_id: user.id,
        showtime_id: showtimeIdNum,
        seat_ids: selectedSeatDbIds,
        ...(appliedDiscount ? { discount_code: appliedDiscount } : {}),
      })
      // Redirect to the booking confirmation page
      router.push(`/booking-confirmation?ref=${data.booking_ref}`)
    } catch (error) {
      setSubmitError(error instanceof Error ? error.message : "Something went wrong while submitting")
    } finally {
      setSubmitting(false)
    }
  }

  const totalDue = preview?.totals.total_due
  const ticketSubtotal = preview?.totals.subtotal ?? seatSubtotal
  const convenienceFee = preview?.totals.convenience_fee ?? (ticketSubtotal !== null && ticketSubtotal !== undefined ? 2 : undefined)
  const taxAmount =
    preview?.totals.tax_amount ??
    (ticketSubtotal !== null && ticketSubtotal !== undefined && convenienceFee !== undefined
      ? Number(((ticketSubtotal + convenienceFee) * 0.08).toFixed(2))
      : undefined)
  const estimatedTotal =
    ticketSubtotal !== null && ticketSubtotal !== undefined && convenienceFee !== undefined && taxAmount !== undefined
      ? Number((ticketSubtotal + convenienceFee + taxAmount).toFixed(2))
      : undefined
  const formattedTicketSubtotal =
    ticketSubtotal !== null && ticketSubtotal !== undefined ? `$${ticketSubtotal.toFixed(2)}` : "—"
  const formattedTotal =
    totalDue !== undefined
      ? `$${totalDue.toFixed(2)}`
      : estimatedTotal !== undefined
        ? `$${estimatedTotal.toFixed(2)}`
        : "—"
  const timerMinutes = String(Math.floor(timerSeconds / 60)).padStart(2, "0")
  const timerRemainderSeconds = String(timerSeconds % 60).padStart(2, "0")

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
                  placeholder="0000-0000-0000-0000"
                  value={cardNumber}
                  maxLength={19}
                  onChange={(e) => {
                    const digits = e.target.value.replace(/\D/g, "").slice(0, 16)
                    const formatted = digits.match(/.{1,4}/g)?.join("-") ?? digits
                    setCardNumber(formatted)
                  }}
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
                    placeholder="MM/YY"
                    value={expiry}
                    maxLength={5}
                    onChange={(e) => {
                      const digits = e.target.value.replace(/\D/g, "").slice(0, 4)
                      const formatted = digits.length > 2 ? `${digits.slice(0, 2)}/${digits.slice(2)}` : digits
                      setExpiry(formatted)
                    }}
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
                disabled={submitting || previewLoading}
                className="mt-2 w-full rounded-xl bg-primary px-6 py-3.5 font-bold text-white transition hover:brightness-110 disabled:opacity-60"
              >
                {submitting ? "Processing..." : `Pay ${formattedTotal} Securely`}
              </button>

              {submitError && <p className="text-sm text-red-400">{submitError}</p>}
              {successMessage && <p className="text-sm text-green-400">{successMessage}</p>}
              {previewError && <p className="text-sm text-yellow-400">{previewError}</p>}

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
            <span className="font-mono text-lg font-bold">{timerMinutes}:{timerRemainderSeconds}</span>
          </div>

          <div className="overflow-hidden rounded-2xl border border-neutral-800 bg-surface-dark shadow-2xl">
            <div className="h-44 bg-[url('https://images.unsplash.com/photo-1536440136628-849c177e76a1?q=80&w=1600&auto=format&fit=crop')] bg-cover bg-center" />

            <div className="space-y-5 p-6">
              <div>
                <h3 className="text-2xl font-bold">{showtimeInfo?.movie_title ?? "—"}</h3>
                <p className="text-sm text-neutral-500">
                  {showtimeInfo ? `${showtimeInfo.format} · ${showtimeInfo.language}` : "Loading..."}
                </p>
              </div>

              <div className="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <p className="text-xs uppercase tracking-wide text-neutral-500">Theater</p>
                  <p className="font-medium">{showtimeInfo?.theater_name ?? "—"}</p>
                  <p className="text-xs text-neutral-500">{showtimeInfo?.screen_name ?? "—"}</p>
                </div>
                <div>
                  <p className="text-xs uppercase tracking-wide text-neutral-500">Date & Time</p>
                  <p className="font-medium">{showtimeInfo?.show_date ?? "—"}</p>
                  <p className="text-xs text-neutral-500">{showtimeInfo?.start_time ?? "—"}</p>
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

              <div className="space-y-2 border-t border-neutral-800 pt-4 text-sm">
                {previewLoading ? (
                  <p className="text-center text-xs text-neutral-500">Loading pricing...</p>
                ) : (
                  <>
                    <div className="flex justify-between text-neutral-500">
                      <span>Ticket Price ({preview?.line_items.length ?? selectedSeatDbIds.length}x)</span>
                      <span>{formattedTicketSubtotal}</span>
                    </div>
                    <div className="flex justify-between text-neutral-500">
                      <span>Convenience Fee</span>
                      <span>{convenienceFee !== undefined ? `$${convenienceFee.toFixed(2)}` : "—"}</span>
                    </div>
                    <div className="flex justify-between text-neutral-500">
                      <span>Tax</span>
                      <span>{taxAmount !== undefined ? `$${taxAmount.toFixed(2)}` : "—"}</span>
                    </div>
                    {preview && preview.totals.discount_amount > 0 && (
                      <div className="flex justify-between text-green-400">
                        <span>Discount ({preview.totals.discount_code})</span>
                        <span>-${preview.totals.discount_amount.toFixed(2)}</span>
                      </div>
                    )}
                    <div className="mt-3 flex justify-between text-lg font-bold">
                      <span>Total Amount</span>
                      <span className="text-primary">{formattedTotal}</span>
                    </div>
                  </>
                )}
              </div>

              <div className="flex gap-2">
                <input
                  type="text"
                  placeholder="Have a promo code?"
                  value={discountCode}
                  onChange={(e) => setDiscountCode(e.target.value)}
                  className="w-full rounded-lg border border-neutral-800 bg-neutral-900 px-3 py-2 text-sm placeholder:text-neutral-500 focus:border-primary focus:outline-none"
                />
                <button
                  type="button"
                  onClick={handleApplyDiscount}
                  className="rounded-lg border border-neutral-700 px-3 py-2 text-xs font-semibold"
                >
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
