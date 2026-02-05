'use client';

import { useEffect, useState } from 'react';
import { useRouter, useParams } from 'next/navigation';

/**
 * View Form Page
 *
 * Displays a single form's details. Currently redirects to print preview.
 */
export default function ViewFormPage() {
    const router = useRouter();
    const params = useParams();
    const formId = params.id as string;

    return (
        <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950">
            <header className="bg-white dark:bg-zinc-900 shadow">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
                    <div className="flex justify-between items-center">
                        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">
                            View Form
                        </h1>
                        <button
                            onClick={() => router.push('/dashboard')}
                            className="px-4 py-2 bg-zinc-200 dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50 rounded-lg hover:bg-zinc-300 dark:hover:bg-zinc-700 transition-colors"
                        >
                            Back to Dashboard
                        </button>
                    </div>
                </div>
            </header>

            <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                <div className="bg-white dark:bg-zinc-900 rounded-lg shadow p-6">
                    <p className="text-zinc-600 dark:text-zinc-400">
                        Form ID: {formId}
                    </p>
                    <p className="mt-4 text-zinc-600 dark:text-zinc-400">
                        This page is under construction.
                    </p>
                </div>
            </main>
        </div>
    );
}
