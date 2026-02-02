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
// Pesticide Application Types
// ============================================================================

export interface PesticideApplication {
    chem_used: number;
    app_timestamp: string;
    rate: string;
    amount_applied: number;
    location_code: string;
}

export interface PesticideApplicationResponse extends PesticideApplication {
    id: number;
    form_id: string;
}

// ============================================================================
// Forms API Request Types (match backend/internal/handlers/types.go)
// ============================================================================

export interface CreateShrubFormRequest {
    first_name: string;
    last_name: string;
    street_number: string;
    street_name: string;
    town: string;
    zip_code: string;
    home_phone: string;
    other_phone: string;
    call_before: boolean;
    is_holiday: boolean;
    flea_only: boolean;
    applications?: PesticideApplication[];
}

export interface CreateLawnFormRequest {
    first_name: string;
    last_name: string;
    street_number: string;
    street_name: string;
    town: string;
    zip_code: string;
    home_phone: string;
    other_phone: string;
    call_before: boolean;
    is_holiday: boolean;
    lawn_area_sq_ft: number;
    fert_only: boolean;
    applications?: PesticideApplication[];
}

export interface UpdateShrubFormRequest {
    first_name: string;
    last_name: string;
    street_number: string;
    street_name: string;
    town: string;
    zip_code: string;
    home_phone: string;
    other_phone: string;
    call_before: boolean;
    is_holiday: boolean;
    flea_only: boolean;
}

export interface UpdateLawnFormRequest {
    first_name: string;
    last_name: string;
    street_number: string;
    street_name: string;
    town: string;
    zip_code: string;
    home_phone: string;
    other_phone: string;
    call_before: boolean;
    is_holiday: boolean;
    lawn_area_sq_ft: number;
    fert_only: boolean;
}

/**
 * Query parameters for listing forms with pagination, filtering, and sorting
 *
 *
 * limit?: number;
   offset?: number;
   form_type?: 'shrub' | 'pesticide';
   search_name?: string;
   sort_by?: 'first_name' | 'last_name' | 'created_at';
   order?: 'ASC' | 'DESC';
   chemical_ids?: number[];
 */
export interface ListFormsParams {
    limit?: number | null;
    offset?: number | null;
    form_type?: string | null;
    search_name?: string | null;
    sort_by?: string | null;
    order?: string | null;
    chemical_ids?: number[] | null;
    zip_code?: string | null;
    jewish_holiday?: string | null;
    date_low?: string | null;
    date_high?: string | null;
}

// ============================================================================
// Forms API Response Types (match backend/internal/handlers/types.go)
// ============================================================================
//
/**
 * ShrubForm - matches backend ShrubFormResponse exactly
 * Contains all form fields with all shrub form specific fields
 */

/**
 * CreateFormResponse - wrapper for hold the recently created form's ID 
 */
export interface CreateFormResponse {
    id: string;
}

export interface ShrubForm {
    id: string;
    created_by: string;
    created_at: string;
    updated_at: string;
    form_type: 'shrub' | 'lawn';
    first_name: string;
    last_name: string;
    street_number: string;
    street_name: string;
    town: string;
    zip_code: string;
    home_phone: string;
    other_phone: string;
    call_before: boolean;
    is_holiday: boolean;
    first_app_date: string;
    last_app_date: string;
    // Shrub-specific field
    flea_only: boolean;
    pest_apps: PesticideApplicationResponse[];
}

/**
 * LawnForm - matches backend LawnFormResponse exactly
 * Contains all form fields with all lawn form specific fields
 *
 * 
 *  id: string;
 *  created_by: string;
 *  created_at: string;
 *  updated_at: string;
 *  form_type: 'shrub' | 'lawn';
 *  first_name: string;
 *  last_name: string;
 *  street_number: string;
 *  street_name: string;
 *  town: string;
 *  zip_code: string;
 *  home_phone: string;
 *  other_phone: string;
 *  call_before: boolean;
 *  is_holiday: boolean;
 *  first_app_date: string;
 *  last_app_date: string;
 *  // Lawn-specific fields
 *  lawn_area_sq_ft: number;
 *  fert_only: boolean;
 *  pest_apps: PesticideApplicationResponse[];
 */
export interface LawnForm {
    id: string;
    created_by: string;
    created_at: string;
    updated_at: string;
    form_type: 'shrub' | 'lawn';
    first_name: string;
    last_name: string;
    street_number: string;
    street_name: string;
    town: string;
    zip_code: string;
    home_phone: string;
    other_phone: string;
    call_before: boolean;
    is_holiday: boolean;
    first_app_date: string;
    last_app_date: string;
    // Lawn-specific fields
    lawn_area_sq_ft: number;
    fert_only: boolean;
    pest_apps: PesticideApplicationResponse[];
}

/**
 * FormViewResponse - matches backend FormResponse exactly
 * Contains all form fields with optional flea_only (shrub forms) or lawn_area_sq_ft/fert_only (lawn forms)
 */
export interface FormViewResponse {
    id: string;
    created_by: string;
    created_at: string;
    updated_at: string;
    form_type: 'shrub' | 'lawn';
    first_name: string;
    last_name: string;
    street_number: string;
    street_name: string;
    town: string;
    zip_code: string;
    home_phone: string;
    other_phone: string;
    call_before: boolean;
    is_holiday: boolean;
    first_app_date: string;
    last_app_date: string;
    // Shrub-specific field (null if lawn form)
    flea_only?: boolean | null;
    // Lawn-specific fields (null if shrub form)
    lawn_area_sq_ft?: number | null;
    fert_only?: boolean | null;
    pest_apps: PesticideApplicationResponse[];
}

/**
 * ListFormsResponse - matches backend ListFormsResponse exactly
 */
export interface ListFormsResponse {
    forms: FormViewResponse[];
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
// Chemicals API Types
// ============================================================================

export interface Chemical {
    id: number;
    category: 'lawn' | 'shrub';
    brand_name: string;
    chemical_name: string;
    epa_reg_no: string;
    recipe: string;
    unit: string;
}

export interface CreateChemicalRequest {
    category: 'lawn' | 'shrub';
    brand_name: string;
    chemical_name: string;
    epa_reg_no: string;
    recipe: string;
    unit: string;
}

export interface ListChemicalsResponse {
    chemicals: Chemical[];
    count: number;
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
