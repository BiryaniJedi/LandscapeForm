'use client';

import { useState, FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import { formsClient } from '@/lib/api/forms';

export default function CreateLawnFormPage() {
    const router = useRouter();
    const [firstName, setFirstName] = useState('');
    const [lastName, setLastName] = useState('');
    const [streetNumber, setStreetNumber] = useState('');
    const [streetName, setStreetName] = useState('');
    const [town, setTown] = useState('');
    const [zipCode, setZipCode] = useState('');
    const [homePhone, setHomePhone] = useState('');
    const [otherPhone, setOtherPhone] = useState('');
    const [callBefore, setCallBefore] = useState(false);
    const [isHoliday, setIsHoliday] = useState(false);
    const [lawnAreaSqFt, setLawnAreaSqFt] = useState('');
    const [fertOnly, setFertOnly] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);

    const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setError(null);
        setIsSubmitting(true);

        try {
            const response = await formsClient.createLawnForm({
                first_name: firstName,
                last_name: lastName,
                street_number: streetNumber,
                street_name: streetName,
                town: town,
                zip_code: zipCode,
                home_phone: homePhone,
                other_phone: otherPhone,
                call_before: callBefore,
                is_holiday: isHoliday,
                lawn_area_sq_ft: parseInt(lawnAreaSqFt) || 0,
                fert_only: fertOnly,
            });

            router.push(`/forms/lawn/${response.id}`);
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
                        Create Lawn Form
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

                        <div className="grid grid-cols-2 gap-4">
                            <div>
                                <label htmlFor="streetNumber" className="block text-sm font-medium text-zinc-900 dark:text-zinc-50 mb-2">
                                    Street Number
                                </label>
                                <input
                                    id="streetNumber"
                                    type="text"
                                    value={streetNumber}
                                    onChange={(e) => setStreetNumber(e.target.value)}
                                    required
                                    className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50"
                                />
                            </div>
                            <div>
                                <label htmlFor="streetName" className="block text-sm font-medium text-zinc-900 dark:text-zinc-50 mb-2">
                                    Street Name
                                </label>
                                <input
                                    id="streetName"
                                    type="text"
                                    value={streetName}
                                    onChange={(e) => setStreetName(e.target.value)}
                                    required
                                    className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50"
                                />
                            </div>
                        </div>

                        <div className="grid grid-cols-2 gap-4">
                            <div>
                                <label htmlFor="town" className="block text-sm font-medium text-zinc-900 dark:text-zinc-50 mb-2">
                                    Town
                                </label>
                                <input
                                    id="town"
                                    type="text"
                                    value={town}
                                    onChange={(e) => setTown(e.target.value)}
                                    required
                                    className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50"
                                />
                            </div>
                            <div>
                                <label htmlFor="zipCode" className="block text-sm font-medium text-zinc-900 dark:text-zinc-50 mb-2">
                                    Zip Code
                                </label>
                                <input
                                    id="zipCode"
                                    type="text"
                                    value={zipCode}
                                    onChange={(e) => setZipCode(e.target.value)}
                                    required
                                    className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50"
                                />
                            </div>
                        </div>

                        <div className="grid grid-cols-2 gap-4">
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
                                <label htmlFor="otherPhone" className="block text-sm font-medium text-zinc-900 dark:text-zinc-50 mb-2">
                                    Other Phone
                                </label>
                                <input
                                    id="otherPhone"
                                    type="text"
                                    value={otherPhone}
                                    onChange={(e) => setOtherPhone(e.target.value)}
                                    className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50"
                                />
                            </div>
                        </div>

                        <div>
                            <label htmlFor="lawnAreaSqFt" className="block text-sm font-medium text-zinc-900 dark:text-zinc-50 mb-2">
                                Lawn Area (sq ft)
                            </label>
                            <input
                                id="lawnAreaSqFt"
                                type="number"
                                value={lawnAreaSqFt}
                                onChange={(e) => setLawnAreaSqFt(e.target.value)}
                                required
                                min="0"
                                className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50"
                            />
                        </div>

                        <div className="space-y-2">
                            <label className="flex items-center">
                                <input
                                    type="checkbox"
                                    checked={callBefore}
                                    onChange={(e) => setCallBefore(e.target.checked)}
                                    className="mr-2"
                                />
                                <span className="text-sm text-zinc-900 dark:text-zinc-50">Call Before Visit</span>
                            </label>

                            <label className="flex items-center">
                                <input
                                    type="checkbox"
                                    checked={isHoliday}
                                    onChange={(e) => setIsHoliday(e.target.checked)}
                                    className="mr-2"
                                />
                                <span className="text-sm text-zinc-900 dark:text-zinc-50">Holiday Property</span>
                            </label>

                            <label className="flex items-center">
                                <input
                                    type="checkbox"
                                    checked={fertOnly}
                                    onChange={(e) => setFertOnly(e.target.checked)}
                                    className="mr-2"
                                />
                                <span className="text-sm text-zinc-900 dark:text-zinc-50">Fertilizer Only</span>
                            </label>
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
