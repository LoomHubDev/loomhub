import { describe, it, expect, vi, afterEach } from 'vitest'
import { timeAgo, fileSize } from '@/utils/format'

describe('timeAgo', () => {
  afterEach(() => {
    vi.useRealTimers()
  })

  it('returns "just now" for less than 60 seconds ago', () => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date('2025-06-15T12:00:30Z'))
    expect(timeAgo('2025-06-15T12:00:00Z')).toBe('just now')
  })

  it('returns minutes ago', () => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date('2025-06-15T12:05:00Z'))
    expect(timeAgo('2025-06-15T12:00:00Z')).toBe('5m ago')
  })

  it('returns hours ago', () => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date('2025-06-15T15:00:00Z'))
    expect(timeAgo('2025-06-15T12:00:00Z')).toBe('3h ago')
  })

  it('returns days ago', () => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date('2025-06-20T12:00:00Z'))
    expect(timeAgo('2025-06-15T12:00:00Z')).toBe('5d ago')
  })

  it('returns months ago', () => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date('2025-09-15T12:00:00Z'))
    expect(timeAgo('2025-06-15T12:00:00Z')).toBe('3mo ago')
  })

  it('returns years ago', () => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date('2027-06-15T12:00:00Z'))
    expect(timeAgo('2025-06-15T12:00:00Z')).toBe('2y ago')
  })

  it('returns "just now" at exactly 0 seconds', () => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date('2025-06-15T12:00:00Z'))
    expect(timeAgo('2025-06-15T12:00:00Z')).toBe('just now')
  })

  it('returns "1m ago" at exactly 60 seconds', () => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date('2025-06-15T12:01:00Z'))
    expect(timeAgo('2025-06-15T12:00:00Z')).toBe('1m ago')
  })

  it('returns "1h ago" at exactly 60 minutes', () => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date('2025-06-15T13:00:00Z'))
    expect(timeAgo('2025-06-15T12:00:00Z')).toBe('1h ago')
  })

  it('returns "1d ago" at exactly 24 hours', () => {
    vi.useFakeTimers()
    vi.setSystemTime(new Date('2025-06-16T12:00:00Z'))
    expect(timeAgo('2025-06-15T12:00:00Z')).toBe('1d ago')
  })
})

describe('fileSize', () => {
  it('returns "0 B" for 0 bytes', () => {
    expect(fileSize(0)).toBe('0 B')
  })

  it('returns bytes without decimals', () => {
    expect(fileSize(500)).toBe('500 B')
  })

  it('returns 1 B for 1 byte', () => {
    expect(fileSize(1)).toBe('1 B')
  })

  it('converts to KB with one decimal', () => {
    expect(fileSize(1024)).toBe('1.0 KB')
  })

  it('converts to KB with fractional value', () => {
    expect(fileSize(1536)).toBe('1.5 KB')
  })

  it('converts to MB', () => {
    expect(fileSize(1048576)).toBe('1.0 MB')
  })

  it('converts to MB with fraction', () => {
    expect(fileSize(2621440)).toBe('2.5 MB')
  })

  it('converts to GB', () => {
    expect(fileSize(1073741824)).toBe('1.0 GB')
  })

  it('handles large MB values', () => {
    expect(fileSize(524288000)).toBe('500.0 MB')
  })
})
