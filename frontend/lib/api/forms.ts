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

/**
 * Client for interacting with Landscaping Forms backend API.
 *
 * This client wraps all `/api/forms/*` endpoints and provides
 * strongly-typed methods for creating, updating, and retrieving forms.
 *
 * @extends ApiClient
 */
export class FormsClient extends ApiClient {
    /**
     * Create a new shrub form.
     *
     * Sends a `POST` request to `/api/forms/shrub` with the provided form data.
     *
     * @param createShrubFormRequest - Payload used to create the shrub form
     * @returns A promise that resolves to the created for
     *
     * @throws {FormValidationError} If the request payload is invalid
     * @throws {AuthError} If the user is not authenticated
     * @throws {FormServerError} If the server encounters an unexpected error
     */
    async CreateShrubForm(createShrubFormRequest: CreateShrubFormRequest): Promise<FormResponse> {
        return this.request<FormResponse>('/forms/shrub', {
            method: 'POST',
            body: JSON.stringify(createShrubFormRequest),
            credentials: 'include',
        })
    }
    /**
     * Create a new pesticide form.
     *
     * Sends a `POST` request to `/api/forms/pesticide`.
     *
     * @param createPesticideFormRequest - Payload used to create the pesticide form
     * @returns A promise that resolves to the created form
     *
     * @throws {FormValidationError} If the request payload is invalid
     * @throws {AuthError} If the user is not authenticated
     * @throws {FormServerError} If the server encounters an unexpected error
     */
    async CreatePesticideForm(createPesticideFormRequest: CreatePesticideFormRequest): Promise<FormResponse> {
        return this.request<FormResponse>('/forms/pesticide', {
            method: 'POST',
            body: JSON.stringify(createPesticideFormRequest),
            credentials: 'include',
        })
    }

    /**
     * Update an existing shrub form.
     *
     * Sends a `PUT` request to `/api/forms/shrub/{formID}`.
     *
     * @param formID - Unique identifier of the shrub form
     * @param updateShrubFormRequest - Updated form values
     * @returns A promise that resolves to the updated form
     *
     * @throws {FormNotFoundError} If the form does not exist
     * @throws {FormValidationError} If the update payload is invalid
     * @throws {AuthError} If the user is not authenticated
     */
    async UpdateShrubForm(formID: string, updateShrubFormRequest: UpdateShrubFormRequest): Promise<FormResponse> {
        return this.request<FormResponse>(`/forms/shrub/${formID}`, {
            method: 'PUT',
            body: JSON.stringify(updateShrubFormRequest),
            credentials: 'include',
        })
    }

    /**
     * Update an existing pesticide form.
     *
     * Sends a `PUT` request to `/api/forms/pesticide/{formID}`.
     *
     * @param formID - Unique identifier of the pesticide form
     * @param updatePesticideFormRequest - Updated form values
     * @returns A promise that resolves to the updated form
     *
     * @throws {FormNotFoundError} If the form does not exist
     * @throws {FormValidationError} If the update payload is invalid
     * @throws {AuthError} If the user is not authenticated
     */
    async UpdatePesticideForm(formID: string, updatePesticideFormRequest: UpdatePesticideFormRequest): Promise<FormResponse> {
        return this.request<FormResponse>(`/forms/pesticide/${formID}`, {
            method: 'PUT',
            body: JSON.stringify(updatePesticideFormRequest),
            credentials: 'include',
        })
    }

    /**
     * Retrieve a single form by its ID.
     *
     * Sends a `GET` request to `/api/forms/{formID}` and returns a
     * generic {@link FormResponse}, regardless of form type.
     *
     * @param formID - Unique identifier of the form to retrieve
     * @returns A promise that resolves to the requested form
     *
     * @throws {FormNotFoundError} If the form does not exist
     * @throws {AuthError} If the user is not authenticated
     */
    async GetFormView(formID: string): Promise<FormResponse> {
        return this.request<FormResponse>(`/forms/${formID}`, {
            method: 'GET',
            credentials: 'include',
        })
    }
}

/**
 * Singleton instance of {@link FormsClient}.
 *
 * Use this instance for all form-related API interactions.
 */
export const formsClient = new FormsClient(); Object
