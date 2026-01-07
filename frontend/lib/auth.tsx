'use client';

/**
 * Auth Context Provider
 * Manages authentication state across the application
 */

import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { api, User } from './api';

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
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    // Load user from localStorage on mount
    useEffect(() => {
        try {
            const storedUser = localStorage.getItem('user');
            if (storedUser) {
                setUser(JSON.parse(storedUser));
            }
        } catch (error) {
            console.error('Failed to load user from localStorage:', error);
        } finally {
            setIsLoading(false);
        }
    }, []);

    const login = async (username: string, password: string) => {
        try {
            const response = await api.login({ username, password });
            setUser(response.user);
            localStorage.setItem('user', JSON.stringify(response.user));
        } catch (error) {
            throw error;
        }
    };

    const register = async (userData: {
        first_name: string;
        last_name: string;
        date_of_birth: string;
        username: string;
        password: string;
    }) => {
        try {
            const response = await api.register(userData);
            setUser(response.user);
            localStorage.setItem('user', JSON.stringify(response.user));
        } catch (error) {
            throw error;
        }
    };

    const logout = async () => {
        await api.logout();
        setUser(null);
        localStorage.removeItem('user');
    };

    const value = {
        user,
        isAuthenticated: !!user,
        isLoading,
        login,
        register,
        logout,
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
