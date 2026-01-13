import {
    CreateShrubFormRequest,
    CreatePesticideFormRequest,
    UpdateShrubFormRequest,
    UpdatePesticideFormRequest,
    ListFormsParams,
    ShrubForm,
    PesticideForm,
    FormViewResponse,
    CreateFormResponse,
    ListFormsResponse,
    SuccessResponse,
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
     * @returns A promise that resolves to the created form's ID
     *
     * @throws {FormValidationError} If the request payload is invalid
     * @throws {AuthError} If the user is not authenticated
     * @throws {FormServerError} If the server encounters an unexpected error
     */
    async createShrubForm(createShrubFormRequest: CreateShrubFormRequest): Promise<CreateFormResponse> {
        return await this.request<CreateFormResponse>('/forms/shrub', {
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
     * @returns A promise that resolves to the created form's ID
     *
     * @throws {FormValidationError} If the request payload is invalid
     * @throws {AuthError} If the user is not authenticated
     * @throws {FormServerError} If the server encounters an unexpected error
     */
    async createPesticideForm(createPesticideFormRequest: CreatePesticideFormRequest): Promise<CreateFormResponse> {
        return await this.request<CreateFormResponse>('/forms/pesticide', {
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
    async updateShrubForm(formID: string, updateShrubFormRequest: UpdateShrubFormRequest): Promise<ShrubForm> {
        return await this.request<ShrubForm>(`/forms/shrub/${formID}`, {
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
    async updatePesticideForm(formID: string, updatePesticideFormRequest: UpdatePesticideFormRequest): Promise<PesticideForm> {
        return await this.request<PesticideForm>(`/forms/pesticide/${formID}`, {
            method: 'PUT',
            body: JSON.stringify(updatePesticideFormRequest),
            credentials: 'include',
        })
    }

    /**
     * Retrieve a single form by its ID.
     *
     * Sends a `GET` request to `/api/forms/{formID}` and returns a
     * generic {@link FormViewResponse}, regardless of form type.
     *
     * @param formID - Unique identifier of the form to retrieve
     * @returns A promise that resolves to the requested form
     *
     * @throws {FormNotFoundError} If the form does not exist
     * @throws {AuthError} If the user is not authenticated
     */
    async getFormView(formID: string): Promise<FormViewResponse> {
        return await this.request<FormViewResponse>(`/forms/${formID}`, {
            method: 'GET',
            credentials: 'include',
        })
    }

    /**
     * Retrieve a single shrub form by its ID.
     *
     * Sends a `GET` request to `/api/forms/shrub/{formID}` and returns a
     * generic {@link FormViewResponse}, regardless of form type.
     *
     * @param formID - Unique identifier of the form to retrieve
     * @returns A promise that resolves to the requested form
     *
     * @throws {FormNotFoundError} If the form does not exist
     * @throws {AuthError} If the user is not authenticated
     */
    async getShrubForm(formID: string): Promise<ShrubForm> {
        return await this.request<ShrubForm>(`/forms/shrub/${formID}`, {
            method: 'GET',
            credentials: 'include',
        })
    }

    /**
     * Retrieve a single pesticide form by its ID.
     *
     * Sends a `GET` request to `/api/forms/shrub/{formID}` and returns a
     * generic {@link FormViewResponse}, regardless of form type.
     *
     * @param formID - Unique identifier of the form to retrieve
     * @returns A promise that resolves to the requested form
     *
     * @throws {FormNotFoundError} If the form does not exist
     * @throws {AuthError} If the user is not authenticated
     */
    async getPesticideForm(formID: string): Promise<PesticideForm> {
        return await this.request<PesticideForm>(`/forms/pesticide/${formID}`, {
            method: 'GET',
            credentials: 'include',
        })
    }

    /**
     * List forms for the authenticated user.
     *
     * Sends a `GET` request to `/api/forms` with optional query parameters
     * for pagination, filtering, and sorting.
     *
     * @param params - Optional query parameters (limit, offset, form_type, search_name, sort_by, order)
     * @returns A promise that resolves to a list of forms and total count
     *
     * @throws {AuthError} If the user is not authenticated
     * @throws {FormServerError} If the server encounters an unexpected error
     */
    async listFormsByUserId(params?: ListFormsParams): Promise<ListFormsResponse> {
        const queryParams = new URLSearchParams()

        if (params?.limit !== undefined) queryParams.append('limit', params.limit.toString())
        if (params?.offset !== undefined) queryParams.append('offset', params.offset.toString())
        if (params?.form_type) queryParams.append('type', params.form_type)
        if (params?.search_name) queryParams.append('search', params.search_name)
        if (params?.sort_by) queryParams.append('sort_by', params.sort_by)
        if (params?.order) queryParams.append('order', params.order)

        const queryString = queryParams.toString()
        const url = queryString ? `/forms?${queryString}` : '/forms'

        return await this.request<ListFormsResponse>(url, {
            method: 'GET',
            credentials: 'include',
        })
    }

    /**
     * List all forms from all users (admin only).
     *
     * Sends a `GET` request to `/api/admin/forms` with optional query parameters
     * for pagination, filtering, and sorting.
     *
     * @param params - Optional query parameters (limit, offset, form_type, search_name, sort_by, order)
     * @returns A promise that resolves to a list of all forms and total count
     *
     * @throws {AuthError} If the user is not authenticated or not an admin
     * @throws {FormServerError} If the server encounters an unexpected error
     */
    async listAllForms(params?: ListFormsParams): Promise<ListFormsResponse> {
        const queryParams = new URLSearchParams()

        if (params?.limit !== undefined) queryParams.append('limit', params.limit.toString())
        if (params?.offset !== undefined) queryParams.append('offset', params.offset.toString())
        if (params?.form_type) queryParams.append('type', params.form_type)
        if (params?.search_name) queryParams.append('search', params.search_name)
        if (params?.sort_by) queryParams.append('sort_by', params.sort_by)
        if (params?.order) queryParams.append('order', params.order)

        const queryString = queryParams.toString()
        const url = queryString ? `/admin/forms?${queryString}` : '/admin/forms'

        return await this.request<ListFormsResponse>(url, {
            method: 'GET',
            credentials: 'include',
        })
    }

    /**
     * Delete a form by its ID.
     *
     * Sends a `DELETE` request to `/api/forms/{formID}`.
     *
     * @param formID - Unique identifier of the form to delete
     * @returns A promise that resolves to a success message
     *
     * @throws {FormNotFoundError} If the form does not exist
     * @throws {AuthError} If the user is not authenticated
     * @throws {FormServerError} If the server encounters an unexpected error
     */
    async deleteForm(formID: string): Promise<SuccessResponse> {
        return await this.request<SuccessResponse>(`/forms/${formID}`, {
            method: 'DELETE',
            credentials: 'include',
        })
    }
}

/**
 * Singleton instance of {@link FormsClient}.
 *
 * Use this instance for all form-related API interactions.
 */
export const formsClient = new FormsClient();
