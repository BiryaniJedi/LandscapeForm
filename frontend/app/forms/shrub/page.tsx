'use client';

import { useState, useEffect, FormEvent } from 'react';
import { useRouter } from 'next/navigation';
import { formsClient } from '@/lib/api/forms';
import { chemicalsClient } from '@/lib/api/chemicals';
import { Chemical, PesticideApplication } from '@/lib/api/types';
import { siteCodesFirst, siteCodesSecond } from '@/lib/common/siteCodes';

/**
 * Create Shrub Form Page
 *
 * Form for creating new shrub pesticide application records.
 * Includes customer information and multiple pesticide applications.
 */
export default function CreateShrubFormPage() {
    const router = useRouter();
    const [chemicals, setChemicals] = useState<Chemical[]>([]);
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
    const [fleaOnly, setFleaOnly] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [isLoading, setIsLoading] = useState(false);

    const [applications, setApplications] = useState<PesticideApplication[]>([]);
    const [chemIndex, setChemIndex] = useState('');
    const [appTimestamp, setAppTimestamp] = useState('');
    const [rateAmt, setRateAmt] = useState('');
    const [rateSqFt, setRateSqFt] = useState('');
    const [rate, setRate] = useState('');
    const [amountApplied, setAmountApplied] = useState('');
    const [locationCode, setLocationCode] = useState('');

    const fetchShrubChemicals = async () => {
        setIsLoading(true)
        try {
            const data = await chemicalsClient.listChemicals('shrub');
            setChemicals(data.chemicals);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to load chemicals');
        } finally {
            setIsLoading(false)
        }
    };

    const handleAddApplication = () => {
        if (!chemIndex || !appTimestamp || !rateAmt || !rateSqFt || !amountApplied || !locationCode) {
            setError('Please fill in all application fields');
            return;
        }

        const chemIndexNum = parseInt(chemIndex);
        if (chemIndexNum < 0 || chemIndexNum >= chemicals.length) {
            setError(`Chemical index must be between 0 and ${chemicals.length - 1}`);
            return;
        }

        const rateAmtNum = parseInt(rateAmt);
        const rateSqFtNum = parseInt(rateSqFt);
        if (rateAmtNum < 0) {
            setError(`Rate Amount in units must be greater than 0`);
            return;
        }
        if (rateSqFtNum < 0) {
            setError(`Rate Amount of Square Feet must be greater than 0`);
            return;
        }

        // Validate location code format
        if (locationCode.length !== 2) {
            setError('Location code must be exactly 2 characters');
            return;
        }
        const firstChar = locationCode[0];
        const secondChar = locationCode[1].toUpperCase();

        if (!Object.keys(siteCodesFirst).includes(firstChar)) {
            setError(`First character of location code must be one of: ${Object.keys(siteCodesFirst).join(', ')}`);
            return;
        }
        if (!Object.keys(siteCodesSecond).includes(secondChar)) {
            setError(`Second character of location code must be one of: ${Object.keys(siteCodesSecond).join(', ')}`);
            return;
        }

        const newApplication: PesticideApplication = {
            chem_used: chemicals[chemIndexNum].id,
            app_timestamp: new Date(appTimestamp).toISOString(),
            rate: rate,
            amount_applied: parseFloat(amountApplied),
            location_code: firstChar + secondChar,
        };

        setApplications([...applications, newApplication]);
        setChemIndex('');
        setAppTimestamp('');
        setRateAmt('');
        setRateSqFt('');
        setRate('');
        setAmountApplied('');
        setLocationCode('');
        setError(null);
    };

    const handleRemoveApplication = (index: number) => {
        setApplications(applications.filter((_, i) => i !== index));
    };

    const handleSubmit = async (e: FormEvent<HTMLFormElement>) => {
        e.preventDefault();
        setError(null);
        setIsSubmitting(true);

        if (applications.length === 0) {
            setError('At least one pesticide application is required');
            setIsSubmitting(false);
            return;
        }

        try {
            const response = await formsClient.createShrubForm({
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
                flea_only: fleaOnly,
                applications: applications,
            });

            router.push(`/forms/shrub/${response.id}`);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to create form');
        } finally {
            setIsSubmitting(false);
        }
    };

    useEffect(() => {
        fetchShrubChemicals();
    }, []);

    useEffect(() => {
        if (rateAmt && rateSqFt) {
            setRate(`${rateAmt} units/${rateSqFt} sq ft`);
        } else {
            setRate('');
        }
    }, [rateAmt, rateSqFt]);

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

    return (
        <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950">
            <header className="bg-white dark:bg-zinc-900 shadow">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
                    <div className="flex justify-between items-center">
                        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">
                            Create Shrub Form
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

            <main className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-4">
                <div className="bg-white dark:bg-zinc-900 rounded-lg shadow p-6">
                    <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">
                        Chemical List
                    </h1>

                    <div className="bg-white dark:bg-zinc-900 rounded-lg shadow overflow-hidden">
                        <div className="overflow-x-auto">
                            <table className="min-w-full divide-y divide-zinc-200 dark:divide-zinc-800">
                                <thead className="bg-zinc-50 dark:bg-zinc-800">
                                    <tr>
                                        <th className="px-6 py-3 text-left text-xs font-medium text-zinc-500 dark:text-zinc-400 uppercase tracking-wider">
                                            Index
                                        </th>
                                        <th className="px-6 py-3 text-left text-xs font-medium text-zinc-500 dark:text-zinc-400 uppercase tracking-wider">
                                            Category
                                        </th>
                                        <th className="px-6 py-3 text-left text-xs font-medium text-zinc-500 dark:text-zinc-400 uppercase tracking-wider">
                                            Brand Name
                                        </th>
                                        <th className="px-6 py-3 text-left text-xs font-medium text-zinc-500 dark:text-zinc-400 uppercase tracking-wider">
                                            Chemical Name
                                        </th>
                                        <th className="px-6 py-3 text-left text-xs font-medium text-zinc-500 dark:text-zinc-400 uppercase tracking-wider">
                                            EPA Reg No
                                        </th>
                                        <th className="px-6 py-3 text-left text-xs font-medium text-zinc-500 dark:text-zinc-400 uppercase tracking-wider">
                                            Recipe
                                        </th>
                                        <th className="px-6 py-3 text-left text-xs font-medium text-zinc-500 dark:text-zinc-400 uppercase tracking-wider">
                                            Unit
                                        </th>
                                    </tr>
                                </thead>
                                <tbody className="bg-white dark:bg-zinc-900 divide-y divide-zinc-200 dark:divide-zinc-800">
                                    {chemicals.map((chem, index) => (
                                        <tr key={chem.id}>
                                            <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-zinc-900 dark:text-zinc-50">
                                                {index}
                                            </td>
                                            <td className="px-6 py-4 whitespace-nowrap">
                                                <span className={`px-2 inline-flex text-xs leading-5 font-semibold rounded-full ${chem.category === 'lawn'
                                                    ? 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                                                    : 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200'
                                                    }`}>
                                                    {chem.category}
                                                </span>
                                            </td>
                                            <td className="px-6 py-4 whitespace-nowrap text-sm text-zinc-900 dark:text-zinc-50">
                                                {chem.brand_name}
                                            </td>
                                            <td className="px-6 py-4 whitespace-nowrap text-sm text-zinc-900 dark:text-zinc-50">
                                                {chem.chemical_name}
                                            </td>
                                            <td className="px-6 py-4 whitespace-nowrap text-sm text-zinc-500 dark:text-zinc-400">
                                                {chem.epa_reg_no || 'N/A'}
                                            </td>
                                            <td className="px-6 py-4 whitespace-nowrap text-sm text-zinc-500 dark:text-zinc-400">
                                                {chem.recipe || 'N/A'}
                                            </td>
                                            <td className="px-6 py-4 whitespace-nowrap text-sm text-zinc-500 dark:text-zinc-400">
                                                {chem.unit || 'N/A'}
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>
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
                                <span className="text-sm text-zinc-900 dark:text-zinc-50">Jewish Holiday</span>
                            </label>

                            <label className="flex items-center">
                                <input
                                    type="checkbox"
                                    checked={fleaOnly}
                                    onChange={(e) => setFleaOnly(e.target.checked)}
                                    className="mr-2"
                                />
                                <span className="text-sm text-zinc-900 dark:text-zinc-50">Flea Only</span>
                            </label>
                        </div>

                        <div className="border-t border-zinc-200 dark:border-zinc-700 pt-4 mt-6">
                            <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50 mb-4">
                                Pesticide Applications
                            </h2>

                            {applications.length > 0 && (
                                <div className="mb-4 space-y-2">
                                    {applications.map((app, index) => {
                                        const chem = chemicals.find(c => c.id === app.chem_used);
                                        const chemIdx = chemicals.findIndex(c => c.id === app.chem_used);
                                        return (
                                            <div
                                                key={index}
                                                className="bg-zinc-50 dark:bg-zinc-800 p-3 rounded-lg flex justify-between items-start"
                                            >
                                                <div className="flex-1">
                                                    <p className="text-sm font-medium text-zinc-900 dark:text-zinc-50">
                                                        Chemical [{chemIdx}]: {chem?.brand_name || 'Unknown'}
                                                    </p>
                                                    <p className="text-xs text-zinc-600 dark:text-zinc-400">
                                                        Timestamp: {new Date(app.app_timestamp).toLocaleString()} |
                                                        Rate: {app.rate} |
                                                        Amount: {app.amount_applied} |
                                                        Location: {app.location_code}
                                                    </p>
                                                </div>
                                                <button
                                                    type="button"
                                                    onClick={() => handleRemoveApplication(index)}
                                                    className="ml-2 px-2 py-1 bg-red-600 text-white rounded hover:bg-red-700 transition-colors text-sm"
                                                >
                                                    Remove
                                                </button>
                                            </div>
                                        );
                                    })}
                                </div>
                            )}

                            <div className="bg-zinc-50 dark:bg-zinc-800 p-4 rounded-lg space-y-3">
                                <div className="grid grid-cols-2 gap-3">
                                    <div>
                                        <label htmlFor="chemIndex" className="block text-sm font-medium text-zinc-900 dark:text-zinc-50 mb-1">
                                            Chemical Index (from table above)
                                        </label>
                                        <input
                                            id="chemIndex"
                                            type="number"
                                            value={chemIndex}
                                            onChange={(e) => setChemIndex(e.target.value)}
                                            min="0"
                                            max={chemicals.length - 1}
                                            className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-900 text-zinc-900 dark:text-zinc-50"
                                            placeholder="0"
                                        />
                                    </div>
                                    <div>
                                        <label htmlFor="appTimestamp" className="block text-sm font-medium text-zinc-900 dark:text-zinc-50 mb-1">
                                            Application Date & Time
                                        </label>
                                        <input
                                            id="appTimestamp"
                                            type="datetime-local"
                                            value={appTimestamp}
                                            onChange={(e) => setAppTimestamp(e.target.value)}
                                            className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-900 text-zinc-900 dark:text-zinc-50"
                                        />
                                    </div>
                                </div>

                                <div className="grid grid-cols-3 gap-3">
                                    <div>
                                        <label htmlFor="rate" className="block text-sm font-medium text-zinc-900 dark:text-zinc-50 mb-1">
                                            Rate
                                        </label>
                                        <div className="flex">
                                            <input
                                                id="rateUnit"
                                                type="number"
                                                value={rateAmt}
                                                onChange={(e) => setRateAmt(e.target.value)}
                                                className="w-1/2 px-3 py-2 mr-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-900 text-zinc-900 dark:text-zinc-50"
                                                placeholder=""
                                            />
                                            <p> units <br /> per </p>
                                        </div>
                                        <div className="flex">
                                            <input
                                                id="rateSqFt"
                                                type="number"
                                                value={rateSqFt}
                                                onChange={(e) => setRateSqFt(e.target.value)}
                                                className="w-1/2 px-3 py-2 mr-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-900 text-zinc-900 dark:text-zinc-50"
                                                placeholder=""
                                            />
                                            <p> Sq Ft </p>
                                        </div>
                                    </div>
                                    <div>
                                        <label htmlFor="amountApplied" className="block text-sm font-medium text-zinc-900 dark:text-zinc-50 mb-1">
                                            Amount Applied
                                        </label>
                                        <input
                                            id="amountApplied"
                                            type="number"
                                            step="0.01"
                                            value={amountApplied}
                                            onChange={(e) => setAmountApplied(e.target.value)}
                                            className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-900 text-zinc-900 dark:text-zinc-50"
                                            placeholder="0.00"
                                        />
                                    </div>
                                    <div>
                                        <label htmlFor="locationCode" className="block text-sm font-medium text-zinc-900 dark:text-zinc-50 mb-1">
                                            Location Code
                                        </label>
                                        <input
                                            id="locationCode"
                                            type="text"
                                            value={locationCode}
                                            onChange={(e) => setLocationCode(e.target.value.slice(0, 2))}
                                            maxLength={2}
                                            className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-900 text-zinc-900 dark:text-zinc-50"
                                            placeholder="FL"
                                        />
                                    </div>
                                </div>

                                <button
                                    type="button"
                                    onClick={handleAddApplication}
                                    className="w-full px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors flex items-center justify-center"
                                >
                                    <span className="mr-2">+</span> Add Application
                                </button>
                            </div>
                        </div>

                        <div>
                            <h2 className="text-lg font-semibold text-zinc-900 dark:text-zinc-50 mb-4">
                                Location Code Reference (combine one from each table)
                            </h2>
                            <div className="grid grid-cols-2 gap-4">
                                {/* First Site Code Table */}
                                <div className="bg-white dark:bg-zinc-900 rounded-lg shadow overflow-hidden">
                                    <div className="overflow-x-auto">
                                        <table className="min-w-full divide-y divide-zinc-200 dark:divide-zinc-800">
                                            <thead className="bg-zinc-50 dark:bg-zinc-800">
                                                <tr>
                                                    <th className="px-4 py-3 text-left text-xs font-medium text-zinc-500 dark:text-zinc-400 uppercase tracking-wider">
                                                        Code
                                                    </th>
                                                    <th className="px-4 py-3 text-left text-xs font-medium text-zinc-500 dark:text-zinc-400 uppercase tracking-wider">
                                                        Location
                                                    </th>
                                                </tr>
                                            </thead>
                                            <tbody className="bg-white dark:bg-zinc-900 divide-y divide-zinc-200 dark:divide-zinc-800">
                                                {Object.entries(siteCodesFirst).map(([code, location]) => (
                                                    <tr key={code}>
                                                        <td className="px-4 py-3 whitespace-nowrap text-sm font-medium text-zinc-900 dark:text-zinc-50">
                                                            {code}
                                                        </td>
                                                        <td className="px-4 py-3 whitespace-nowrap text-sm text-zinc-900 dark:text-zinc-50">
                                                            {location}
                                                        </td>
                                                    </tr>
                                                ))}
                                            </tbody>
                                        </table>
                                    </div>
                                </div>

                                {/* Second Site Code Table */}
                                <div className="bg-white dark:bg-zinc-900 rounded-lg shadow overflow-hidden">
                                    <div className="overflow-x-auto">
                                        <table className="min-w-full divide-y divide-zinc-200 dark:divide-zinc-800">
                                            <thead className="bg-zinc-50 dark:bg-zinc-800">
                                                <tr>
                                                    <th className="px-4 py-3 text-left text-xs font-medium text-zinc-500 dark:text-zinc-400 uppercase tracking-wider">
                                                        Code
                                                    </th>
                                                    <th className="px-4 py-3 text-left text-xs font-medium text-zinc-500 dark:text-zinc-400 uppercase tracking-wider">
                                                        Location
                                                    </th>
                                                </tr>
                                            </thead>
                                            <tbody className="bg-white dark:bg-zinc-900 divide-y divide-zinc-200 dark:divide-zinc-800">
                                                {Object.entries(siteCodesSecond).map(([code, location]) => (
                                                    <tr key={code}>
                                                        <td className="px-4 py-3 whitespace-nowrap text-sm font-medium text-zinc-900 dark:text-zinc-50">
                                                            {code}
                                                        </td>
                                                        <td className="px-4 py-3 whitespace-nowrap text-sm text-zinc-900 dark:text-zinc-50">
                                                            {location}
                                                        </td>
                                                    </tr>
                                                ))}
                                            </tbody>
                                        </table>
                                    </div>
                                </div>
                            </div>
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
