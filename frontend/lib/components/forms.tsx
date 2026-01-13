'use client';


/**
 * Forms Context Provider
 * Manages CRUD operations on forms across the application
 */

import { createContext, useContext, useState, useEffect, useCallback, ReactNode } from 'react';
import {
    CreateShrubFormRequest,
    CreatePesticideFormRequest,
    UpdateShrubFormRequest,
    UpdatePesticideFormRequest,
    FormViewResponse,
    ListFormsResponse

} from '../api/types';
import { authClient } from '../api/auth'
import { formsClient } from '../api/forms'

interface FormsContextType {
    formview: FormViewResponse;
    formviewList: ListFormsResponse;
    isLoading: boolean;
    login: (username: string, password: string) => Promise<void>;
    register: (userData: {
        first_name: string;
        last_name: string;
        date_of_birth: string;
        username: string;
        password: string;
    }) => Promise<void>;
    logout: () => void;
    refreshUser: () => Promise<void>;
}

const FormsContext = createContext<FormsContextType | undefined>(undefined);

export function FormsProvider({ children }: { children: ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    const refreshUser = useCallback(async () => {
        setIsLoading(true);
        try {
            const response = await authClient.me();
            setUser(response);
        } catch (error) {
            throw error;
        } finally {
            setIsLoading(false);
        }
    }, []);

    // fetch user from db on mount
    useEffect(() => {
        const init = async () => {
            try {
                await refreshUser();
            } catch (error) {
                console.error('Failed to fetch user:', error);
            }
        };

        init();
    }, [refreshUser]);

    const login = async (username: string, password: string) => {
        setIsLoading(true);
        try {
            const response = await authClient.login({ username, password });
            setUser(response.user);
        } catch (error) {
            throw error;
        } finally {
            setIsLoading(false);
        }
    };

    const register = async (userData: {
        first_name: string;
        last_name: string;
        date_of_birth: string;
        username: string;
        password: string;
    }) => {
        setIsLoading(true);
        try {
            const response = await authClient.register(userData);
            setUser(response.user);
        } catch (error) {
            throw error;
        } finally {
            setIsLoading(false);
        }
    };

    const logout = async () => {
        await authClient.logout();
        setUser(null);
    };

    const value: FormsContextType = {
        user,
        isAuthenticated: !!user,
        isLoading,
        login,
        register,
        logout,
        refreshUser,
    };

    return <FormsContext.Provider value={value}>{children}</AuthContext.Provider>;
}

/**
 * Hook to use auth context
 */
export function useForms() {
    const context = useContext(FormsContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
}

