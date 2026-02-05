/**
 * Formats an ISO 8601 date string to a localized human-readable format.
 *
 * @param dateString - ISO 8601 date string (e.g., "2024-01-15T10:30:00Z")
 * @returns Localized date and time string based on user's browser locale
 *
 * @example
 * formatDate("2024-01-15T10:30:00Z") // "1/15/2024, 10:30:00 AM" (en-US)
 */
export const formatDate = (dateString: string) => {
    const dateObject = new Date(dateString);
    return dateObject.toLocaleString();
}
