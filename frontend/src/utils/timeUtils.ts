/**
 * Utility functions for time formatting in IST (Indian Standard Time)
 * IST is UTC+5:30
 */

/**
 * Converts a date/time string to IST and formats it as time string
 * @param dateString - ISO date string from backend (assumed to be in UTC)
 * @returns Formatted time string in IST (HH:MM:SS)
 */
export const formatTimeIST = (dateString: string): string => {
  const date = new Date(dateString);
  
  // Use Intl.DateTimeFormat to format in IST timezone
  const formatter = new Intl.DateTimeFormat('en-IN', {
    timeZone: 'Asia/Kolkata',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
  });
  
  const parts = formatter.formatToParts(date);
  const hours = parts.find(p => p.type === 'hour')?.value.padStart(2, '0') || '00';
  const minutes = parts.find(p => p.type === 'minute')?.value.padStart(2, '0') || '00';
  const seconds = parts.find(p => p.type === 'second')?.value.padStart(2, '0') || '00';
  
  return `${hours}:${minutes}:${seconds}`;
};

/**
 * Converts a date/time string to IST and formats it as date string
 * @param dateString - ISO date string from backend
 * @returns Formatted date string in IST (DD/MM/YYYY)
 */
export const formatDateIST = (dateString: string): string => {
  const date = new Date(dateString);
  
  // Use Intl.DateTimeFormat to format in IST timezone
  const formatter = new Intl.DateTimeFormat('en-IN', {
    timeZone: 'Asia/Kolkata',
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  });
  
  return formatter.format(date);
};

/**
 * Gets current time in IST
 * @returns Current time string in IST format (HH:MM:SS)
 */
export const getCurrentTimeIST = (): string => {
  const now = new Date();
  
  // Use Intl.DateTimeFormat to format in IST timezone
  const formatter = new Intl.DateTimeFormat('en-IN', {
    timeZone: 'Asia/Kolkata',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: false,
  });
  
  const parts = formatter.formatToParts(now);
  const hours = parts.find(p => p.type === 'hour')?.value.padStart(2, '0') || '00';
  const minutes = parts.find(p => p.type === 'minute')?.value.padStart(2, '0') || '00';
  const seconds = parts.find(p => p.type === 'second')?.value.padStart(2, '0') || '00';
  
  return `${hours}:${minutes}:${seconds}`;
};

/**
 * Gets current date in IST
 * @returns Current date in IST timezone
 */
export const getCurrentDateIST = (): Date => {
  // Return current date - JavaScript Date objects are timezone-aware
  // When formatting, use the IST timezone formatters
  return new Date();
};

/**
 * Formats date with weekday in IST
 * @param dateString - ISO date string from backend
 * @returns Formatted date string with weekday
 */
export const formatDateWithWeekdayIST = (dateString: string): string => {
  const date = new Date(dateString);
  
  // Use Intl.DateTimeFormat to format in IST timezone
  const formatter = new Intl.DateTimeFormat('en-IN', {
    timeZone: 'Asia/Kolkata',
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  });
  
  return formatter.format(date);
};
