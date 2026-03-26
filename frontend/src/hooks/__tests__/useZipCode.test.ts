import { renderHook, act } from '@testing-library/react'
import { useZipCode } from '../useZipCode'

const mockPush = jest.fn()

jest.mock('next/navigation', () => ({
  useRouter: () => ({ push: mockPush }),
}))

describe('useZipCode', () => {
  beforeEach(() => {
    localStorage.clear()
    mockPush.mockClear()
  })

  it('getZipCode returns null when nothing stored', () => {
    const { result } = renderHook(() => useZipCode())
    expect(result.current.getZipCode()).toBeNull()
  })

  it('setZipCode stores value in localStorage', () => {
    const { result } = renderHook(() => useZipCode())
    act(() => {
      result.current.setZipCode('10001')
    })
    expect(localStorage.getItem('selectedZipCode')).toBe('10001')
  })

  it('getZipCode returns stored zip code', () => {
    localStorage.setItem('selectedZipCode', '90210')
    const { result } = renderHook(() => useZipCode())
    expect(result.current.getZipCode()).toBe('90210')
  })

  it('clearZipCode removes stored zip code', () => {
    localStorage.setItem('selectedZipCode', '10001')
    const { result } = renderHook(() => useZipCode())
    act(() => {
      result.current.clearZipCode()
    })
    expect(localStorage.getItem('selectedZipCode')).toBeNull()
  })

  it('requireZipCode returns zip code if set', () => {
    localStorage.setItem('selectedZipCode', '55401')
    const { result } = renderHook(() => useZipCode())
    let zip: string | null = null
    act(() => {
      zip = result.current.requireZipCode()
    })
    expect(zip).toBe('55401')
    expect(mockPush).not.toHaveBeenCalled()
  })

  it('requireZipCode redirects to "/" when no zip code set', () => {
    const { result } = renderHook(() => useZipCode())
    let zip: string | null = null
    act(() => {
      zip = result.current.requireZipCode()
    })
    expect(zip).toBeNull()
    expect(mockPush).toHaveBeenCalledWith('/')
  })
})
