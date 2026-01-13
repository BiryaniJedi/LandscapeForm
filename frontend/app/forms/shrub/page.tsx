'use client';

import { useState, FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import { formsClient } from '@/lib/api/forms';

export default function CreateShrubFormPage() {
    const router = useRouter();
    const [firstName, setFirstName] = useState('');
    const [lastName, setLastName] = useState('');
    const [homePhone, setHomePhone] = useState('');
    const [numShrubs, setNumShrubs] = useState('');
    const [error, setError] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);

    const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setError(null);
        setIsSubmitting(true);

        try {
            const response = await formsClient.createShrubForm({
                first_name: firstName,
                last_name: lastName,
                home_phone: homePhone,
                num_shrubs: parseInt(numShrubs, 10),
            });

            // Redirect to the created form's detail page
            router.push(`/forms/shrub/${response.id}`);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to create form');
        } finally {
            setIsSubmitting(false);
        }
    };

    return (
        <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950">
            <header className="bg-white dark:bg-zinc-900 shadow">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
                    <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">
                        Create Shrub Form
                    </h1>
                </div>
            </header>

            <main className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                <div className="bg-white dark:bg-zinc-900 rounded-lg shadow p-6">
                    <form onSubmit={handleSubmit} className="space-y-4">
                        {error && (
                            <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
                                <p className="text-red-800 dark:text-red-200">{error}</p>
                            </div>
                        )}

                        <div>
                            <label htmlFor="firstName" className="block text-sm font-medium text-zinc-900 dark:text-zinc-50 mb-2">
                                First Name
                            </label>
                            <input
                                id="firstName"
                                type="text"
                                value={firstName}
                                onChange={(e) => setFirstName(e.target.value)}
                                required
                                className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50"
                            />
                        </div>

                        <div>
                            <label htmlFor="lastName" className="block text-sm font-medium text-zinc-900 dark:text-zinc-50 mb-2">
                                Last Name
                            </label>
                            <input
                                id="lastName"
                                type="text"
                                value={lastName}
                                onChange={(e) => setLastName(e.target.value)}
                                required
                                className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50"
                            />
                        </div>

                        <div>
                            <label htmlFor="homePhone" className="block text-sm font-medium text-zinc-900 dark:text-zinc-50 mb-2">
                                Home Phone
                            </label>
                            <input
                                id="homePhone"
                                type="text"
                                value={homePhone}
                                onChange={(e) => setHomePhone(e.target.value)}
                                required
                                className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50"
                            />
                        </div>

                        <div>
                            <label htmlFor="numShrubs" className="block text-sm font-medium text-zinc-900 dark:text-zinc-50 mb-2">
                                Number of Shrubs
                            </label>
                            <input
                                id="numShrubs"
                                type="number"
                                value={numShrubs}
                                onChange={(e) => setNumShrubs(e.target.value)}
                                required
                                min="1"
                                className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50"
                            />
                        </div>

                        <div className="flex gap-4">
                            <button
                                type="submit"
                                disabled={isSubmitting}
                                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                            >
                                {isSubmitting ? 'Creating...' : 'Create Form'}
                            </button>
                            <button
                                type="button"
                                onClick={() => router.push('/dashboard')}
                                className="px-4 py-2 bg-zinc-200 dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50 rounded-lg hover:bg-zinc-300 dark:hover:bg-zinc-700 transition-colors"
                            >
                                Cancel
                            </button>
                        </div>
                    </form>
                </div>
            </main>
        </div>
    );
}
