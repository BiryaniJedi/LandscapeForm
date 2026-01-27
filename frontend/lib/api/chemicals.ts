import { Chemical, CreateChemicalRequest, ListChemicalsResponse, SuccessResponse } from './types'
import ApiClient from './common'

/**
 * Client for interacting with Chemicals Management API.
 *
 * This client wraps all `/api/chemicals/*` and `/api/admin/chemicals/*` endpoints
 * and provides strongly-typed methods for managing chemicals.
 *
 * @extends ApiClient
 */
export class ChemicalsClient extends ApiClient {
    /**
     * List all chemicals (admin only).
     *
     * Sends a `GET` request to `/api/admin/chemicals`.
     *
     * @param category - Optional category filter ('lawn' or 'shrub')
     * @returns A promise that resolves to a list of all chemicals
     *
     * @throws {AuthError} If the user is not authenticated or not an admin
     */
    async listChemicals(category?: 'lawn' | 'shrub'): Promise<ListChemicalsResponse> {
        const params = new URLSearchParams()
        if (category) params.append('category', category)

        const queryString = params.toString()
        const url = queryString ? `/chemicals?${queryString}` : '/chemicals'

        return await this.request<ListChemicalsResponse>(url, {
            method: 'GET',
            credentials: 'include',
        })
    }

    /**
     * List chemicals by category (public endpoint).
     *
     * Sends a `GET` request to `/api/chemicals/category/{category}`.
     *
     * @param category - Category filter ('lawn' or 'shrub')
     * @returns A promise that resolves to a list of chemicals in the category
     */
    async listChemicalsByCategory(category: 'lawn' | 'shrub'): Promise<ListChemicalsResponse> {
        return await this.request<ListChemicalsResponse>(`/chemicals/category/${category}`, {
            method: 'GET',
            credentials: 'include',
        })
    }

    /**
     * Create a new chemical (admin only).
     *
     * Sends a `POST` request to `/api/admin/chemicals`.
     *
     * @param chemical - The chemical data to create
     * @returns A promise that resolves to the created chemical's ID
     *
     * @throws {AuthError} If the user is not authenticated or not an admin
     */
    async createChemical(chemical: CreateChemicalRequest): Promise<{ id: string }> {
        return await this.request<{ id: string }>('/admin/chemicals', {
            method: 'POST',
            body: JSON.stringify(chemical),
            credentials: 'include',
        })
    }

    /**
     * Update an existing chemical (admin only).
     *
     * Sends a `PUT` request to `/api/admin/chemicals/{id}`.
     *
     * @param id - The chemical ID to update
     * @param chemical - The updated chemical data
     * @returns A promise that resolves to the updated chemical
     *
     * @throws {AuthError} If the user is not authenticated or not an admin
     */
    async updateChemical(id: number, chemical: CreateChemicalRequest): Promise<Chemical> {
        return await this.request<Chemical>(`/admin/chemicals/${id}`, {
            method: 'PUT',
            body: JSON.stringify(chemical),
            credentials: 'include',
        })
    }

    /**
     * Delete a chemical (admin only).
     *
     * Sends a `DELETE` request to `/api/admin/chemicals/{id}`.
     *
     * @param id - The chemical ID to delete
     * @returns A promise that resolves to a success message
     *
     * @throws {AuthError} If the user is not authenticated or not an admin
     */
    async deleteChemical(id: number): Promise<SuccessResponse> {
        return await this.request<SuccessResponse>(`/admin/chemicals/${id}`, {
            method: 'DELETE',
            credentials: 'include',
        })
    }
}

/**
 * Singleton instance of {@link ChemicalsClient}.
 *
 * Use this instance for all chemicals-related API interactions.
 */
export const chemicalsClient = new ChemicalsClient();
