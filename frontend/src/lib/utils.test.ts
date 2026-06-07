import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { formatDate, timeAgo, formatPoemFormat } from './utils';

describe('formatDate', () => {
  it('formats an ISO date string to a readable US date', () => {
    var result = formatDate('2024-03-15T12:00:00Z');
    expect(result).toBe('March 15, 2024');
  });

  it('handles dates at the start of the year', () => {
    var result = formatDate('2023-01-01T00:00:00Z');
    expect(result).toBe('January 1, 2023');
  });

  it('handles dates at the end of the year', () => {
    var result = formatDate('2022-12-31T23:59:59Z');
    expect(result).toBe('December 31, 2022');
  });
});

describe('timeAgo', () => {
  var realDateNow: () => number;

  beforeEach(() => {
    realDateNow = Date.now;
  });

  afterEach(() => {
    Date.now = realDateNow;
  });

  it('returns "just now" for dates less than 60 seconds ago', () => {
    var now = new Date('2024-06-01T12:00:00Z').getTime();
    Date.now = vi.fn(() => now + 30_000);
    expect(timeAgo('2024-06-01T12:00:00Z')).toBe('just now');
  });

  it('returns minutes for dates between 1 and 59 minutes ago', () => {
    var now = new Date('2024-06-01T12:00:00Z').getTime();
    Date.now = vi.fn(() => now + 5 * 60 * 1000);
    expect(timeAgo('2024-06-01T12:00:00Z')).toBe('5m ago');
  });

  it('returns "1m ago" at exactly 60 seconds', () => {
    var now = new Date('2024-06-01T12:00:00Z').getTime();
    Date.now = vi.fn(() => now + 60_000);
    expect(timeAgo('2024-06-01T12:00:00Z')).toBe('1m ago');
  });

  it('returns hours for dates between 1 and 23 hours ago', () => {
    var now = new Date('2024-06-01T12:00:00Z').getTime();
    Date.now = vi.fn(() => now + 3 * 3600 * 1000);
    expect(timeAgo('2024-06-01T12:00:00Z')).toBe('3h ago');
  });

  it('returns days for dates between 1 and 6 days ago', () => {
    var now = new Date('2024-06-01T12:00:00Z').getTime();
    Date.now = vi.fn(() => now + 4 * 86400 * 1000);
    expect(timeAgo('2024-06-01T12:00:00Z')).toBe('4d ago');
  });

  it('falls back to formatDate for dates 7 or more days ago', () => {
    var now = new Date('2024-06-01T12:00:00Z').getTime();
    Date.now = vi.fn(() => now + 8 * 86400 * 1000);
    var result = timeAgo('2024-06-01T12:00:00Z');
    expect(result).toBe('June 1, 2024');
  });
});

describe('formatPoemFormat', () => {
  it('capitalises a single-word format', () => {
    expect(formatPoemFormat('haiku')).toBe('Haiku');
  });

  it('capitalises and joins underscore-separated words', () => {
    expect(formatPoemFormat('free_verse')).toBe('Free Verse');
  });

  it('handles multi-segment formats like iambic_pentameter', () => {
    expect(formatPoemFormat('iambic_pentameter')).toBe('Iambic Pentameter');
  });

  it('handles rhyming_couplets', () => {
    expect(formatPoemFormat('rhyming_couplets')).toBe('Rhyming Couplets');
  });

  it('handles already-capitalised input gracefully', () => {
    expect(formatPoemFormat('Custom')).toBe('Custom');
  });
});
