/**
 * Api Client for Landscaping Forms Backend auth routes
 * Provides interface for all HTTP requests to the Go backend API routes 
 */

import { AuthError, ErrorResponse } from './types'
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

export default class ApiClient {
    private baseUrl: string;

    constructor(baseUrl: string) {
        this.baseUrl = baseUrl;
    }

    /**
     * Generic fetch wrapper with error handling
     * Uses session cookies for authentication
     */
    public async request<T>(
        endpoint: string,
        options: RequestInit = {}
    ): Promise<T> {
        const headers: HeadersInit = {
            'Content-Type': 'application/json',
            ...options.headers,
        };

        const url = `${this.baseUrl}${endpoint}`;

        try {
            const response = await fetch(url, {
                ...options,
                headers,
                credentials: 'include',
            });

            // Handle 401 Unauthorized 
            if (response.status === 401) {
                throw new AuthError();
            }

            // Parse response body
            const data = await response.json();

            // Handle non-OK responses
            if (!response.ok) {
                const error = data as ErrorResponse;
                throw new AuthError(error.message || `HTTP ${response.status}: ${response.statusText}`);
            }

            return data as T;
        } catch (error) {
            if (!(error instanceof Error)) {
                throw new Error('UNEXPECTED')
            }
            throw new AuthError(error.message);
        }
    }
}
