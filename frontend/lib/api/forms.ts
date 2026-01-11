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

export class FormsClient extends ApiClient {
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

    async UpdateShrubForm(formID: string, updateShrubFormRequest: UpdateShrubFormRequest): Promise<FormResponse> {
        return this.request<FormResponse>(`/forms/shrub/${formID}`, {
            method: 'PUT',
            body: JSON.stringify(updateShrubFormRequest),
            credentials: 'include',
        })
    }

    async UpdatePesticideForm(formID: string, updatePesticideFormRequest: UpdatePesticideFormRequest): Promise<FormResponse> {
        return this.request<FormResponse>(`/forms/pesticide/${formID}`, {
            method: 'PUT',
            body: JSON.stringify(updatePesticideFormRequest),
            credentials: 'include',
        })
    }

    async GetFormView(formID: string): Promise<FormResponse> {
        return this.request<FormResponse>(`/forms/${formID}`, {
            method: 'GET',
            credentials: 'include',
        })
    }
}

// Export singleton instance
export const formsClient = new FormsClient();
