'use client';

import { useEffect, useState } from 'react';
import { useRouter, useParams } from 'next/navigation';
import { formsClient } from '@/lib/api/forms';
import { ShrubForm } from '@/lib/api/types';

export default function ShrubFormDetailPage() {
    const router = useRouter();
    const params = useParams();
    const formId = params.id as string;

    const [form, setForm] = useState<ShrubForm | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const fetchForm = async () => {
            try {
                const data = await formsClient.getShrubForm(formId);
                setForm(data);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Failed to load form');
            } finally {
                setIsLoading(false);
            }
        };

        fetchForm();
    }, [formId]);

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

    if (error) {
        return (
            <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950">
                <header className="bg-white dark:bg-zinc-900 shadow">
                    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
                        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">
                            Error Loading Form
                        </h1>
                    </div>
                </header>
                <main className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                    <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
                        <p className="text-red-800 dark:text-red-200">{error}</p>
                    </div>
                    <button
                        onClick={() => router.push('/dashboard')}
                        className="mt-4 px-4 py-2 bg-zinc-200 dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50 rounded-lg hover:bg-zinc-300 dark:hover:bg-zinc-700 transition-colors"
                    >
                        Back to Dashboard
                    </button>
                </main>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950">
            <header className="bg-white dark:bg-zinc-900 shadow">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
                    <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">
                        Shrub Form Details
                    </h1>
                </div>
            </header>

            <main className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                <div className="bg-white dark:bg-zinc-900 rounded-lg shadow p-6">
                    {form && (
                        <div className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-zinc-500 dark:text-zinc-400 mb-2">
                                    Form ID
                                </label>
                                <p className="text-zinc-900 dark:text-zinc-50">{form.id}</p>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-zinc-500 dark:text-zinc-400 mb-2">
                                    First Name
                                </label>
                                <p className="text-zinc-900 dark:text-zinc-50">{form.first_name}</p>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-zinc-500 dark:text-zinc-400 mb-2">
                                    Last Name
                                </label>
                                <p className="text-zinc-900 dark:text-zinc-50">{form.last_name}</p>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-zinc-500 dark:text-zinc-400 mb-2">
                                    Home Phone
                                </label>
                                <p className="text-zinc-900 dark:text-zinc-50">{form.home_phone}</p>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-zinc-500 dark:text-zinc-400 mb-2">
                                    Number of Shrubs
                                </label>
                                <p className="text-zinc-900 dark:text-zinc-50">{form.num_shrubs}</p>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-zinc-500 dark:text-zinc-400 mb-2">
                                    Created At
                                </label>
                                <p className="text-zinc-900 dark:text-zinc-50">
                                    {new Date(form.created_at).toLocaleString()}
                                </p>
                            </div>

                            <div>
                                <label className="block text-sm font-medium text-zinc-500 dark:text-zinc-400 mb-2">
                                    Updated At
                                </label>
                                <p className="text-zinc-900 dark:text-zinc-50">
                                    {new Date(form.updated_at).toLocaleString()}
                                </p>
                            </div>
                        </div>
                    )}

                    <div className="mt-6">
                        <button
                            onClick={() => router.push('/dashboard')}
                            className="px-4 py-2 bg-zinc-200 dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50 rounded-lg hover:bg-zinc-300 dark:hover:bg-zinc-700 transition-colors"
                        >
                            Back to Dashboard
                        </button>
                    </div>
                </div>
            </main>
        </div>
    );
}
