import { User, SuccessResponse } from './types'
import ApiClient from './common'

/**
 * Client for interacting with User Management API.
 *
 * This client wraps all `/api/users/*` endpoints and provides
 * strongly-typed methods for managing users (admin only).
 *
 * @extends ApiClient
 */
export class UsersClient extends ApiClient {
    /**
     * Approve a pending user (admin only).
     *
     * Sends a `POST` request to `/api/users/{userID}/approve`.
     *
     * @param userID - Unique identifier of the user to approve
     * @returns A promise that resolves to a success message
     *
     * @throws {AuthError} If the user is not authenticated or not an admin
     */
    async getUserById(userID: string): Promise<SuccessResponse> {
        return await this.request<SuccessResponse>(`/users/${userID}`, {
            method: 'GET',
            credentials: 'include',
        })
    }
    /**
     * List all users (admin only).
     *
     * Sends a `GET` request to `/api/users`.
     *
     * @returns A promise that resolves to a list of all users
     *
     * @throws {AuthError} If the user is not authenticated or not an admin
     */
    async listUsers(): Promise<{ users: User[]; count: number }> {
        return await this.request<{ users: User[]; count: number }>('/users', {
            method: 'GET',
            credentials: 'include',
        })
    }

    /**
     * Approve a pending user (admin only).
     *
     * Sends a `POST` request to `/api/users/{userID}/approve`.
     *
     * @param userID - Unique identifier of the user to approve
     * @returns A promise that resolves to a success message
     *
     * @throws {AuthError} If the user is not authenticated or not an admin
     */
    async approveUser(userID: string): Promise<SuccessResponse> {
        return await this.request<SuccessResponse>(`/users/${userID}/approve`, {
            method: 'POST',
            credentials: 'include',
        })
    }

    /**
     * Delete a user (admin only).
     *
     * Sends a `DELETE` request to `/api/users/{userID}`.
     *
     * @param userID - Unique identifier of the user to delete
     * @returns A promise that resolves to a success message
     *
     * @throws {AuthError} If the user is not authenticated or not an admin
     */
    async deleteUser(userID: string): Promise<SuccessResponse> {
        return await this.request<SuccessResponse>(`/users/${userID}`, {
            method: 'DELETE',
            credentials: 'include',
        })
    }
}

/**
 * Singleton instance of {@link UsersClient}.
 *
 * Use this instance for all user management API interactions.
 */
export const usersClient = new UsersClient();
