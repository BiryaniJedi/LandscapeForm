/*
 * API types define all the types for users, forms, materials, etc
 * Types match backend/internal/handlers/types.go for seamless integration
 */

// ============================================================================
// Auth & User Types
// ============================================================================

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


// ============================================================================
// Forms API Request Types (match backend/internal/handlers/types.go)
// ============================================================================

export interface CreateShrubFormRequest {
    first_name: string;
    last_name: string;
    home_phone: string;
    num_shrubs: number;
}

export interface CreatePesticideFormRequest {
    first_name: string;
    last_name: string;
    home_phone: string;
    pesticide_name: string;
}

export interface UpdateShrubFormRequest {
    first_name: string;
    last_name: string;
    home_phone: string;
    num_shrubs: number;
}

export interface UpdatePesticideFormRequest {
    first_name: string;
    last_name: string;
    home_phone: string;
    pesticide_name: string;
}

/**
 * Query parameters for listing forms with pagination, filtering, and sorting
 */
export interface ListFormsParams {
    limit?: number;
    offset?: number;
    form_type?: 'shrub' | 'pesticide';
    search_name?: string;
    sort_by?: 'first_name' | 'last_name' | 'created_at';
    order?: 'ASC' | 'DESC';
}

// ============================================================================
// Forms API Response Types (match backend/internal/handlers/types.go)
// ============================================================================

/**
 * FormResponse - matches backend FormResponse exactly
 * Contains all form fields with optional num_shrubs (shrub forms) or pesticide_name (pesticide forms)
 */
export interface FormResponse {
    id: string;
    created_by: string;
    created_at: string;
    updated_at: string;
    form_type: 'shrub' | 'pesticide';
    first_name: string;
    last_name: string;
    home_phone: string;
    // Shrub-specific field (null if pesticide form)
    num_shrubs?: number | null;
    // Pesticide-specific field (null if shrub form)
    pesticide_name?: string | null;
}

/**
 * ListFormsResponse - matches backend ListFormsResponse exactly
 */
export interface ListFormsResponse {
    forms: FormResponse[];
    count: number;
}

// ============================================================================
// Generic API Response Types (match backend/internal/handlers/types.go)
// ============================================================================

export interface ErrorResponse {
    error: string;
    message?: string;
}

export interface SuccessResponse {
    message: string;
}

// ============================================================================
// Forms API Error Classes
// ============================================================================

export class FormNotFoundError extends Error {
    constructor(message = 'Form not found or you do not have permission to access it') {
        super(message);
        this.name = 'FormNotFoundError';
    }
}

export class FormValidationError extends Error {
    constructor(message = 'Invalid form data') {
        super(message);
        this.name = 'FormValidationError';
    }
}

export class FormServerError extends Error {
    constructor(message = 'Server error while processing form') {
        super(message);
        this.name = 'FormServerError';
    }
}

export class AuthError extends Error {
    constructor(message = 'Unauthorized') {
        super(message);
        this.name = 'AuthError'
    }
}
