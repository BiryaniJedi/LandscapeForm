/**
 * Auth Client for Landscaping Forms Backend auth routes
 * Handles all HTTP requests to the Go backend API auth routes
 */

import { User, LoginRequest, RegisterRequest, AuthResponse, AuthError, ErrorResponse } from './types'
import ApiClient from './common'
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';

class AuthClient extends ApiClient {
    constructor(baseUrl: string) {
        super(baseUrl)
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
