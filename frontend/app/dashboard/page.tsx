'use client';

import { useState, FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import { useAuth } from '@/lib/components/auth';

export default function DashboardPage() {
    const router = useRouter();
    const { user, isAuthenticated, isLoading, logout, } = useAuth();
    const [error, setError] = useState<string | null>(null);

    if (isLoading) {
        return (
            <div className="flex min-h-screen items-center justify-center bg-zinc-50 dark:bg-zinc-950">
                <div className="text-center">
                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
                    <p className="mt-4 text-zinc-600 dark:text-zinc-400">Loading...</p>
                </div>
            </div>
        );
    }

    if (!isAuthenticated) {
        console.log("Unauthenticated")
        return null; // Will redirect
    } else {
        console.log("AUTHENTICATED!!")
    }

    const handleLogout = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setError(null);

        try {
            logout();
            // Redirect to login on success
            router.push('/login');
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Logout Failed');
        } finally {
        }
    };

    return (
        <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950">
            {/* Header */}
            <header className="bg-white dark:bg-zinc-900 shadow">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex items-center justify-between">
                    <div>
                        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">
                            Dashboard
                        </h1>
                        {user && (
                            <p className="text-sm text-zinc-600 dark:text-zinc-400 mt-1">
                                Welcome, {user.first_name} {user.last_name}
                                {user.role === 'admin' && (
                                    <span className="ml-2 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800 dark:bg-purple-900 dark:text-purple-200">
                                        Admin
                                    </span>
                                )}
                                {user.pending && (
                                    <span className="ml-2 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800 dark:bg-yellow-900 dark:text-yellow-200">
                                        Pending Approval
                                    </span>
                                )}
                            </p>
                        )}
                    </div>
                    <form onSubmit={handleLogout} className="space-y-6">
                        <button
                            type="submit"
                            className="px-4 py-2 bg-zinc-200 dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50 rounded-lg hover:bg-zinc-300 dark:hover:bg-zinc-700 transition-colors"
                        >
                            Logout
                        </button>
                    </form>
                </div>
            </header>

            {/* Main Content */}
            <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                {user?.pending ? (
                    // Pending Approval Message
                    <div className="bg-yellow-50 dark:bg-yellow-900/20 border border-yellow-200 dark:border-yellow-800 rounded-lg p-6">
                        <h2 className="text-lg font-semibold text-yellow-800 dark:text-yellow-200 mb-2">
                            Account Pending Approval
                        </h2>
                        <p className="text-yellow-700 dark:text-yellow-300">
                            Your account is pending admin approval. You'll be able to access all features once an administrator approves your registration.
                        </p>
                    </div>
                ) : (
                    // Dashboard Content
                    <div className="space-y-6">
                        <div className="bg-white dark:bg-zinc-900 rounded-lg shadow p-6">
                            <h2 className="text-xl font-semibold text-zinc-900 dark:text-zinc-50 mb-4">
                                Welcome to Landscaping Forms
                            </h2>
                            <p className="text-zinc-600 dark:text-zinc-400 mb-4">
                                This is your dashboard. Forms management features coming soon!
                            </p>

                            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-6">
                                <div className="border border-zinc-200 dark:border-zinc-800 rounded-lg p-4">
                                    <h3 className="font-medium text-zinc-900 dark:text-zinc-50 mb-2">
                                        Shrub Forms
                                    </h3>
                                    <p className="text-sm text-zinc-600 dark:text-zinc-400">
                                        Manage shrub landscaping forms
                                    </p>
                                    <p className="text-2xl font-bold text-blue-600 dark:text-blue-400 mt-2">
                                        0
                                    </p>
                                </div>

                                <div className="border border-zinc-200 dark:border-zinc-800 rounded-lg p-4">
                                    <h3 className="font-medium text-zinc-900 dark:text-zinc-50 mb-2">
                                        Pesticide Forms
                                    </h3>
                                    <p className="text-sm text-zinc-600 dark:text-zinc-400">
                                        Manage pesticide application forms
                                    </p>
                                    <p className="text-2xl font-bold text-green-600 dark:text-green-400 mt-2">
                                        0
                                    </p>
                                </div>

                                <div className="border border-zinc-200 dark:border-zinc-800 rounded-lg p-4">
                                    <h3 className="font-medium text-zinc-900 dark:text-zinc-50 mb-2">
                                        Total Forms
                                    </h3>
                                    <p className="text-sm text-zinc-600 dark:text-zinc-400">
                                        All forms combined
                                    </p>
                                    <p className="text-2xl font-bold text-purple-600 dark:text-purple-400 mt-2">
                                        0
                                    </p>
                                </div>
                            </div>
                        </div>

                        {/* User Info Card */}
                        {user && (
                            <div className="bg-white dark:bg-zinc-900 rounded-lg shadow p-6">
                                <h2 className="text-xl font-semibold text-zinc-900 dark:text-zinc-50 mb-4">
                                    Your Profile
                                </h2>
                                <dl className="grid grid-cols-1 gap-4 sm:grid-cols-2">
                                    <div>
                                        <dt className="text-sm font-medium text-zinc-500 dark:text-zinc-400">
                                            Username
                                        </dt>
                                        <dd className="mt-1 text-sm text-zinc-900 dark:text-zinc-50">
                                            {user.username}
                                        </dd>
                                    </div>
                                    <div>
                                        <dt className="text-sm font-medium text-zinc-500 dark:text-zinc-400">
                                            Role
                                        </dt>
                                        <dd className="mt-1 text-sm text-zinc-900 dark:text-zinc-50 capitalize">
                                            {user.role}
                                        </dd>
                                    </div>
                                    <div>
                                        <dt className="text-sm font-medium text-zinc-500 dark:text-zinc-400">
                                            Account Status
                                        </dt>
                                        <dd className="mt-1 text-sm text-zinc-900 dark:text-zinc-50">
                                            {user.pending ? 'Pending Approval' : 'Active'}
                                        </dd>
                                    </div>
                                    <div>
                                        <dt className="text-sm font-medium text-zinc-500 dark:text-zinc-400">
                                            Member Since
                                        </dt>
                                        <dd className="mt-1 text-sm text-zinc-900 dark:text-zinc-50">
                                            {new Date(user.created_at).toLocaleDateString()}
                                        </dd>
                                    </div>
                                </dl>
                            </div>
                        )}
                    </div>
                )}
            </main>
        </div>
    );
}
