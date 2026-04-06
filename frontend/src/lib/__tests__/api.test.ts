/**
 * Tests for API client functions in src/lib/api.ts
 */


const BASE = 'http://localhost:8080'

beforeEach(() => {
  jest.resetAllMocks()
})

const mockFetch = (body: unknown, ok = true, status = 200) => {
  global.fetch = jest.fn().mockResolvedValue({
    ok,
    status,
    json: jest.fn().mockResolvedValue(body),
  } as unknown as Response)
}

describe('fetchLocations', () => {
  it('returns locations on success', async () => {
    const locations = [{ id: 1, zipcode: '10001', city: 'New York', state: 'NY', country: 'US' }]
    mockFetch({ locations })
    const result = await fetchLocations()
    expect(result).toEqual(locations)
    expect(global.fetch).toHaveBeenCalledWith(`${BASE}/api/v1/locations`)
  })

  it('throws on error response', async () => {
    mockFetch({}, false, 500)
    await expect(fetchLocations()).rejects.toThrow('Failed to fetch locations')
  })
})

describe('fetchMoviesByZipCode', () => {
  it('returns movies for a given zip code', async () => {
    const movies = [{ id: 1, title: 'Dune', description: '', poster_url: '' }]
    mockFetch(movies)
    const result = await fetchMoviesByZipCode('10001')
    expect(result).toEqual(movies)
    expect(global.fetch).toHaveBeenCalledWith(`${BASE}/api/v1/movies?zipcode=10001`)
  })

  it('throws on error', async () => {
    mockFetch({}, false, 404)
    await expect(fetchMoviesByZipCode('99999')).rejects.toThrow('Failed to fetch movies')
  })
})

describe('fetchMovieById', () => {
  it('returns a single movie by ID', async () => {
    const movie = { id: 5, title: 'Inception', description: '', poster_url: '' }
    mockFetch(movie)
    const result = await fetchMovieById('5')
    expect(result).toEqual(movie)
    expect(global.fetch).toHaveBeenCalledWith(`${BASE}/api/v1/movies/5`)
  })

  it('throws Movie not found on 404', async () => {
    mockFetch({}, false, 404)
    await expect(fetchMovieById('999')).rejects.toThrow('Movie not found')
  })
})

describe('registerUser', () => {
  const payload = { email: 'a@b.com', username: 'bob', password: 'secret', full_name: 'Bob' }

  it('returns auth response on success', async () => {
    const successResp = { message: 'ok', user: { id: 1, email: 'a@b.com', username: 'bob', full_name: 'Bob', phone: null, created_at: '', updated_at: '' } }
    mockFetch(successResp)
    const result = await registerUser(payload)
    expect(result).toEqual(successResp)
    expect(global.fetch).toHaveBeenCalledWith(
      `${BASE}/api/v1/auth/register`,
      expect.objectContaining({ method: 'POST', body: JSON.stringify(payload) })
    )
  })

  it('throws with server error message', async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: false,
      status: 400,
      json: jest.fn().mockResolvedValue({ error: 'Email already exists' }),
    } as unknown as Response)
    await expect(registerUser(payload)).rejects.toThrow('Email already exists')
  })

  it('falls back to default message when no error field', async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: false,
      status: 500,
      json: jest.fn().mockResolvedValue(null),
    } as unknown as Response)
    await expect(registerUser(payload)).rejects.toThrow('Failed to register user')
  })
})

describe('loginUser', () => {
  const payload = { email: 'a@b.com', password: 'secret' }

  it('returns auth response on success', async () => {
    const successResp = { message: 'ok', user: { id: 1, email: 'a@b.com', username: 'bob', full_name: 'Bob', phone: null, created_at: '', updated_at: '' } }
    mockFetch(successResp)
    const result = await loginUser(payload)
    expect(result).toEqual(successResp)
  })

  it('throws with server error message on failure', async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: false,
      status: 401,
      json: jest.fn().mockResolvedValue({ error: 'Invalid credentials' }),
    } as unknown as Response)
    await expect(loginUser(payload)).rejects.toThrow('Invalid credentials')
  })

  it('falls back to default login failure message', async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: false,
      status: 500,
      json: jest.fn().mockRejectedValue(new Error('parse fail')), // simulate JSON parse failure
    } as unknown as Response)
    await expect(loginUser(payload)).rejects.toThrow('Failed to login')
  })
})

describe('fetchMovieShowtimes', () => {
  it('returns showtime response', async () => {
    const showtimesResp = { movie_id: 1, title: 'Dune', dates: ['2026-01-01'], theaters: [] }
    mockFetch(showtimesResp)
    const result = await fetchMovieShowtimes(1, '10001')
    expect(result).toEqual(showtimesResp)
    expect(global.fetch).toHaveBeenCalledWith(`${BASE}/api/v1/movies/1/showtimes?zipcode=10001`)
  })

  it('throws on failure', async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: false,
      status: 500,
      json: jest.fn().mockResolvedValue({}),
    } as unknown as Response)
    await expect(fetchMovieShowtimes(1, '10001')).rejects.toThrow('Failed to fetch showtimes')
  })
})

describe('fetchShowtimeSeats', () => {
  it('returns seats response', async () => {
    const seatsResp = { showtime: {}, layout: {}, summary: {}, seats: [] }
    mockFetch(seatsResp)
    const result = await fetchShowtimeSeats(42)
    expect(result).toEqual(seatsResp)
    expect(global.fetch).toHaveBeenCalledWith(`${BASE}/api/v1/seats?showtime_id=42`)
  })

  it('throws on failure', async () => {
    global.fetch = jest.fn().mockResolvedValue({
      ok: false,
      status: 500,
      json: jest.fn().mockResolvedValue({}),
    } as unknown as Response)
    await expect(fetchShowtimeSeats(42)).rejects.toThrow('Failed to fetch seats')
  })
})
