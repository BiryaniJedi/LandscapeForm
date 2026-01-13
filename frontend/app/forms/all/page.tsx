'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { formsClient } from '@/lib/api/forms';
import { ListFormsResponse, FormViewResponse } from '@/lib/api/types';

export default function ListFormsAllUsersPage() {
    const router = useRouter();

    const [formviewList, setFormviewList] = useState<ListFormsResponse | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    useEffect(() => {
        const fetchForms = async () => {
            try {
                const data = await formsClient.listFormsAllUsers();
                setFormviewList(data);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Failed to load form');
            } finally {
                setIsLoading(false);
            }
        };

        fetchForms();
    }, []);

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
                            Error Loading Forms for All users
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

    if (formviewList == null) {
        return (<div>Not available! FormviewList is null</div>);
    }
    console.log("Reached: ", formviewList)
    return (
        <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950">
            <header className="bg-white dark:bg-zinc-900 shadow">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
                    <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">
                        Your Forms
                    </h1>
                </div>
            </header>

            <main className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                {formviewList && (
                    <div className="space-y-4">
                        {formviewList.forms.map((formview: FormViewResponse) => (
                            <div key={formview.id} className="bg-white dark:bg-zinc-900 rounded-lg shadow p-6">
                                <div>
                                    <label className="block text-sm font-medium text-zinc-500 dark:text-zinc-400 mb-2">
                                        First Name
                                    </label>
                                    <p className="text-zinc-900 dark:text-zinc-50">{formview.first_name}</p>
                                </div>

                                <div>
                                    <label className="block text-sm font-medium text-zinc-500 dark:text-zinc-400 mb-2">
                                        Last Name
                                    </label>
                                    <p className="text-zinc-900 dark:text-zinc-50">{formview.last_name}</p>
                                </div>

                                <div>
                                    <label className="block text-sm font-medium text-zinc-500 dark:text-zinc-400 mb-2">
                                        Home Phone
                                    </label>
                                    <p className="text-zinc-900 dark:text-zinc-50">{formview.home_phone}</p>
                                </div>

                                <div>
                                    <label className="block text-sm font-medium text-zinc-500 dark:text-zinc-400 mb-2">
                                        Form Type
                                    </label>
                                    <p className="text-zinc-900 dark:text-zinc-50">{formview.form_type}</p>
                                </div>

                                <div>
                                    <label className="block text-sm font-medium text-zinc-500 dark:text-zinc-400 mb-2">
                                        {formview.form_type == 'shrub' ? 'Number of Shrubs' : 'Name of Pesticide'}
                                    </label>
                                    <p className="text-zinc-900 dark:text-zinc-50">{formview.form_type == 'shrub' ? formview.num_shrubs : formview.pesticide_name}</p>
                                </div>

                                <div>
                                    <label className="block text-sm font-medium text-zinc-500 dark:text-zinc-400 mb-2">
                                        Created By
                                    </label>
                                    <p className="text-zinc-900 dark:text-zinc-50">{formview.created_by}</p>
                                </div>
                            </div>
                        ))}
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
            </main>
        </div>

    );
}
