import { formatDate, truncateString, generateId, debounce } from '../../lib/utils.js';

describe('Utils Module', () => {
  describe('formatDate', () => {
    test('formats valid date string correctly', () => {
      // Mock Date.prototype.toLocaleDateString to return a fixed string
      const originalToLocaleDateString = Date.prototype.toLocaleDateString;
      Date.prototype.toLocaleDateString = jest.fn().mockReturnValue('Mar 19, 2025, 10:00 AM');
      
      const result = formatDate('2025-03-19T10:00:00Z');
      expect(result).toBe('Mar 19, 2025, 10:00 AM');
      
      // Restore original method
      Date.prototype.toLocaleDateString = originalToLocaleDateString;
    });

    test('returns "Unknown date" for null or undefined input', () => {
      expect(formatDate(null)).toBe('Unknown date');
      expect(formatDate(undefined)).toBe('Unknown date');
      expect(formatDate('')).toBe('Unknown date');
    });

    test('returns "Invalid date" for invalid date string', () => {
      expect(formatDate('not-a-date')).toBe('Invalid date');
    });
  });

  describe('truncateString', () => {
    test('does not truncate strings shorter than maxLength', () => {
      expect(truncateString('Hello', 10)).toBe('Hello');
    });

    test('truncates strings longer than maxLength', () => {
      expect(truncateString('Hello World', 5)).toBe('Hello...');
    });

    test('handles null or undefined input', () => {
      expect(truncateString(null)).toBeNull();
      expect(truncateString(undefined)).toBeUndefined();
    });

    test('uses default maxLength if not provided', () => {
      const longString = 'a'.repeat(150);
      const result = truncateString(longString);
      
      expect(result.length).toBe(103); // 100 chars + 3 for ellipsis
      expect(result.endsWith('...')).toBe(true);
    });
  });

  describe('generateId', () => {
    test('generates a string', () => {
      const id = generateId();
      expect(typeof id).toBe('string');
    });

    test('generates unique IDs', () => {
      const ids = new Set();
      for (let i = 0; i < 100; i++) {
        ids.add(generateId());
      }
      
      // If all IDs are unique, the Set size should be 100
      expect(ids.size).toBe(100);
    });
  });

  describe('debounce', () => {
    jest.useFakeTimers();
    
    test('delays function execution', () => {
      const mockFn = jest.fn();
      const debouncedFn = debounce(mockFn, 500);
      
      debouncedFn();
      
      // Function should not be called immediately
      expect(mockFn).not.toHaveBeenCalled();
      
      // Fast-forward time
      jest.advanceTimersByTime(500);
      
      // Function should be called after the delay
      expect(mockFn).toHaveBeenCalledTimes(1);
    });

    test('only executes once for multiple rapid calls', () => {
      const mockFn = jest.fn();
      const debouncedFn = debounce(mockFn, 500);
      
      debouncedFn();
      debouncedFn();
      debouncedFn();
      
      // Function should not be called yet
      expect(mockFn).not.toHaveBeenCalled();
      
      // Fast-forward time
      jest.advanceTimersByTime(500);
      
      // Function should be called only once
      expect(mockFn).toHaveBeenCalledTimes(1);
    });

    test('resets the timer on subsequent calls', () => {
      const mockFn = jest.fn();
      const debouncedFn = debounce(mockFn, 500);
      
      debouncedFn();
      
      // Fast-forward time partially
      jest.advanceTimersByTime(300);
      
      // Function should not be called yet
      expect(mockFn).not.toHaveBeenCalled();
      
      // Call again, which should reset the timer
      debouncedFn();
      
      // Fast-forward time partially again
      jest.advanceTimersByTime(300);
      
      // Function should still not be called
      expect(mockFn).not.toHaveBeenCalled();
      
      // Fast-forward the remaining time
      jest.advanceTimersByTime(200);
      
      // Function should now be called
      expect(mockFn).toHaveBeenCalledTimes(1);
    });

    test('uses default wait time if not provided', () => {
      const mockFn = jest.fn();
      const debouncedFn = debounce(mockFn); // Default is 300ms
      
      debouncedFn();
      
      // Fast-forward time partially
      jest.advanceTimersByTime(200);
      
      // Function should not be called yet
      expect(mockFn).not.toHaveBeenCalled();
      
      // Fast-forward the remaining time
      jest.advanceTimersByTime(100);
      
      // Function should now be called
      expect(mockFn).toHaveBeenCalledTimes(1);
    });
  });
});
