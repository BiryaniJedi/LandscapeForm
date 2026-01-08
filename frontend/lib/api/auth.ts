/**
 * API Client for Landscaping Forms Backend auth routes
 * Handles all HTTP requests to the Go backend API auth routes
 */

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';
console.log(`API_BASE_URL: ${API_BASE_URL}`)

export interface User {
    id: string;
    created_at: string;
    updated_at: string;
    pending: boolean;
    role: 'employee' | 'admin';
    first_name: string;
    last_name: string;
    date_of_birth: string;
    username: string;
}

export interface LoginRequest {
    username: string;
    password: string;
}

export interface RegisterRequest {
    first_name: string;
    last_name: string;
    date_of_birth: string;
    username: string;
    password: string;
}

export interface AuthResponse {
    token: string;
    user: User;
}

export class AuthError extends Error {
    constructor(message = 'Unauthorized') {
        super(message);
        this.name = 'AuthError'
    }
}

class ApiClient {
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
                const error = data as AuthError;
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
        let response = {} as AuthResponse
        try {
            response = await this.request<AuthResponse>('/auth/login', {
                method: 'POST',
                body: JSON.stringify(credentials),
            });
        } catch (error) {
            if (error instanceof AuthError) {
                throw error as AuthError
            } else if (error instanceof Error) {
                throw new Error(error.message)
            } else {
                throw new Error('Unexpected Error')
            }
        }

        return response;
    }

    /**
     * Register new user
     * Cookie is automatically stored by browser
     */
    async register(userData: RegisterRequest): Promise<AuthResponse> {
        let response = {} as AuthResponse
        try {
            response = await this.request<AuthResponse>('/auth/register', {
                method: 'POST',
                body: JSON.stringify(userData),
            });
        } catch (error) {
            if (error instanceof AuthError) {
                throw error as AuthError
            } else if (error instanceof Error) {
                throw new Error(error.message)
            } else {
                throw new Error('Unexpected Error')
            }
        }

        return response;
    }

    /**
     * Get current user from jwt stored in the http cookie
     */
    async me(): Promise<User> {
        let response = {} as User
        try {
            response = await this.request<User>('/auth/me', {
                method: 'GET',
                credentials: 'include',
            })
        } catch (error) {
            if (error instanceof AuthError) {
                throw error as AuthError
            } else if (error instanceof Error) {
                throw new Error(error.message)
            } else {
                throw new Error('Unexpected Error')
            }
        }

        return response;
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
export const api = new ApiClient(API_BASE_URL);
