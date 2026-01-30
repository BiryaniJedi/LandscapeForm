'use client';

import { useState, useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import dynamic from 'next/dynamic';
import { formsClient } from '@/lib/api/forms';
import { chemicalsClient } from '@/lib/api/chemicals';
import { FormViewResponse, ListChemicalsResponse } from '@/lib/api/types';

// Dynamically import PDF components (they don't work with SSR)
const PDFDownloadLink = dynamic(
    () => import('@react-pdf/renderer').then((mod) => mod.PDFDownloadLink),
    { ssr: false }
);

const PDFViewer = dynamic(
    () => import('@react-pdf/renderer').then((mod) => mod.PDFViewer),
    { ssr: false }
);

// Import the PDF document component
import FormPDFDocument from '@/lib/pdf/FormPDFDocument';

export default function PrintFormPage() {
    const params = useParams();
    const router = useRouter();
    const formId = params.id as string;

    const [form, setForm] = useState<FormViewResponse | null>(null);
    const [chemList, setChemList] = useState<ListChemicalsResponse | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [showPreview, setShowPreview] = useState(true);

    useEffect(() => {
        const fetchForm = async () => {
            try {
                setIsLoading(true);
                const data = await formsClient.getFormView(formId);
                setForm(data);
                fetchChemicals(data.form_type)
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Failed to load form');
            } finally {
                setIsLoading(false);
            }
        };
        const fetchChemicals = async (formType: 'lawn' | 'shrub') => {
            try {
                setIsLoading(true);
                const data = await chemicalsClient.listChemicalsByCategory(formType);
                setChemList(data);
            } catch (err) {
                setError(err instanceof Error ? err.message : 'Failed to fetch chemicals');
            } finally {
                setIsLoading(false);
            }
        }

        if (formId) {
            fetchForm();
        }
    }, [formId]);

    // You can add any data processing logic here before rendering the PDF
    // Example: format dates, calculate totals, etc.
    const processedData = form ? {
        ...form,
        // Add any computed fields or formatted data here
        formattedCreatedAt: new Date(form.created_at).toLocaleDateString(),
        formattedUpdatedAt: new Date(form.updated_at).toLocaleDateString(),
        fullName: `${form.first_name} ${form.last_name}`,
        fullAddress: `${form.street_number} ${form.street_name}, ${form.town} ${form.zip_code}`,
    } : null;

    if (isLoading) {
        return (
            <div className="flex min-h-screen items-center justify-center bg-zinc-50 dark:bg-zinc-950">
                <div className="text-center">
                    <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
                    <p className="mt-4 text-zinc-600 dark:text-zinc-400">Loading form data...</p>
                </div>
            </div>
        );
    }

    if (error || !form) {
        return (
            <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950 flex items-center justify-center">
                <div className="bg-white dark:bg-zinc-900 rounded-lg shadow p-8 max-w-md">
                    <h1 className="text-2xl font-bold text-red-600 dark:text-red-400 mb-4">Error</h1>
                    <p className="text-zinc-900 dark:text-zinc-50 mb-4">
                        {error || 'Form not found'}
                    </p>
                    <button
                        onClick={() => router.push('/forms')}
                        className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
                    >
                        Back to Forms
                    </button>
                </div>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950">
            <header className="bg-white dark:bg-zinc-900 shadow">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
                    <div className="flex justify-between items-center">
                        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">
                            Print Form - {form.first_name} {form.last_name}
                        </h1>
                        <div className="flex gap-2">
                            <button
                                onClick={() => setShowPreview(!showPreview)}
                                className="px-4 py-2 bg-zinc-200 dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50 rounded-lg hover:bg-zinc-300 dark:hover:bg-zinc-700 transition-colors"
                            >
                                {showPreview ? 'Hide Preview' : 'Show Preview'}
                            </button>
                            <button
                                onClick={() => router.push('/dashboard')}
                                className="px-4 py-2 bg-zinc-200 dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50 rounded-lg hover:bg-zinc-300 dark:hover:bg-zinc-700 transition-colors"
                            >
                                Back to Dashboard
                            </button>
                        </div>
                    </div>
                </div>
            </header>

            <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                {showPreview && processedData && (
                    <div className="bg-white dark:bg-zinc-900 rounded-lg shadow p-6">
                        <h2 className="text-xl font-semibold text-zinc-900 dark:text-zinc-50 mb-4">
                            PDF Preview
                        </h2>
                        <div style={{ height: '800px', width: '100%' }}>
                            <PDFViewer style={{ width: '100%', height: '100%' }}>
                                <FormPDFDocument form={processedData} chemicalList={chemList!} />
                            </PDFViewer>
                        </div>
                    </div>
                )}
            </main>
        </div>
    );
}
