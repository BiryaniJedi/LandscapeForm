/**
 * API Client for Landscaping Forms Backend auth routes
 * Handles all HTTP requests to the Go backend API auth routes
 */
import {
    CreateShrubFormRequest,
    CreatePesticideFormRequest,
    UpdateShrubFormRequest,
    UpdatePesticideFormRequest,
    ListFormsParams,
    FormResponse,
    ListFormsResponse,
    ErrorResponse,
    FormNotFoundError,
    FormValidationError,
    FormServerError,
    AuthError
} from './types'

import ApiClient from './common'

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';
console.log(`API_BASE_URL: ${API_BASE_URL}`)

export class FormsClient extends ApiClient {
    constructor(baseUrl: string) {
        super(baseUrl)
    }

    async CreateShrubForm(createShrubFormRequest: CreateShrubFormRequest): Promise<FormResponse> {
        return this.request<FormResponse>('/forms/shrub', {
            method: 'POST',
            body: JSON.stringify(createShrubFormRequest),
            credentials: 'include',
        })
    }

    async CreatePesticideForm(createPesticideFormRequest: CreatePesticideFormRequest): Promise<FormResponse> {
        return this.request<FormResponse>('/forms/shrub', {
            method: 'POST',
            body: JSON.stringify(createPesticideFormRequest),
            credentials: 'include',
        })
    }

    async UpdateShrubForm(updateShrubFormRequest: UpdateShrubFormRequest): Promise<FormResponse> {
        return this.request<FormResponse>('/forms/shrub', {
            method: 'POST',
            body: JSON.stringify(updateShrubFormRequest),
            credentials: 'include',
        })
    }

    async CreatePesticideForm(createPesticideFormRequest: CreatePesticideFormRequest): Promise<FormResponse> {
        return this.request<FormResponse>('/forms/shrub', {
            method: 'POST',
            body: JSON.stringify(createPesticideFormRequest),
            credentials: 'include',
        })
    }
}

// Export singleton instance
export const formsClient = new FormsClient(API_BASE_URL);
