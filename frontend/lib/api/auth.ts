/**
 * API Client for Landscaping Forms Backend auth routes
 * Handles all HTTP requests to the Go backend API auth routes
 */

import { User, LoginRequest, RegisterRequest, AuthResponse, AuthError, ErrorResponse } from './types'
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

class AuthClient {
    private baseUrl: string;

    constructor(baseUrl: string) {
        this.baseUrl = baseUrl;
    }

    /**
     * Generic fetch wrapper with error handling
     * Uses session cookies for authentication
     */
    private async request<T>(
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

    /**
     * Login user
     * Cookie is automatically stored by browser
     */
    async login(credentials: LoginRequest): Promise<AuthResponse> {
        return this.request<AuthResponse>('/auth/login', {
            method: 'POST',
            body: JSON.stringify(credentials),
        });
    }

    /**
     * Register new user
     * Cookie is automatically stored by browser
     */
    async register(userData: RegisterRequest): Promise<AuthResponse> {
        return this.request<AuthResponse>('/auth/register', {
            method: 'POST',
            body: JSON.stringify(userData),
        });
    }

    /**
     * Get current user from jwt stored in the http cookie
     */
    async me(): Promise<User> {
        return this.request<User>('/auth/me', {
            method: 'GET',
            credentials: 'include',
        });
    }

    /**
     * Logout user - calls backend to clear cookie
     */
    async logout(): Promise<void> {
        await this.request<{ message: string }>('/auth/logout', {
            method: 'POST',
        });
    }
}

// Export singleton instance
export const authClient = new AuthClient(API_BASE_URL);
