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

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api';
console.log(`API_BASE_URL: ${API_BASE_URL}`)

export class FormsClient {

}
