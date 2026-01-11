'use client';

/**
 * Auth Context Provider
 * Manages authentication state across the application
 */

import { createContext, useContext, useState, useEffect, useCallback, ReactNode } from 'react';
import { User } from '../api/types';
import { authClient } from '../api/auth'

interface AuthContextType {
    user: User | null;
    isAuthenticated: boolean;
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

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
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

    const value: AuthContextType = {
        user,
        isAuthenticated: !!user,
        isLoading,
        login,
        register,
        logout,
        refreshUser,
    };

    return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

/**
 * Hook to use auth context
 */
export function useAuth() {
    const context = useContext(AuthContext);
    if (context === undefined) {
        throw new Error('useAuth must be used within an AuthProvider');
    }
    return context;
}
