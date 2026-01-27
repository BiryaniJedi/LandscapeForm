'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { formsClient } from '@/lib/api/forms';
import { chemicalsClient } from '@/lib/api/chemicals';
import { ListFormsResponse, FormViewResponse, AuthError, Chemical } from '@/lib/api/types';

export default function ListFormsPage() {
    const router = useRouter();

    const [formviewList, setFormviewList] = useState<ListFormsResponse | null>(null);
    const [error, setError] = useState<Error | AuthError | null>(null);
    const [isLoading, setIsLoading] = useState(true);
    const [deletingFormId, setDeletingFormId] = useState<string | null>(null);
    const [showDeleteConfirm, setShowDeleteConfirm] = useState(false);
    const [formToDelete, setFormToDelete] = useState<FormViewResponse | null>(null);

    //query string params
    const [limit, setLimit] = useState<number>(10);
    const [offset, setOffset] = useState<number>(0);

    const [formType, setFormType] = useState<string>('');
    const [formTypeInput, setFormTypeInput] = useState<string>('');

    const [searchName, setSearchName] = useState<string>('');
    const [searchNameInput, setSearchNameInput] = useState<string>('');

    const [sortBy, setSortBy] = useState<string>('created_at');
    const [sortByInput, setSortByInput] = useState<string>('created_at');

    const [order, setOrder] = useState<string>('DESC');
    const [orderInput, setOrderInput] = useState<string>('DESC');

    // Chemical filtering
    const [chemicals, setChemicals] = useState<Chemical[]>([]);
    const [selectedChemicalIds, setSelectedChemicalIds] = useState<number[]>([]);
    const [selectedChemicalDropdown, setSelectedChemicalDropdown] = useState<string>('');
    const [chemicalsFilter, setChemicalsFilter] = useState<number[]>([]);
    const [chemicalsFilterInput, setChemicalsFilterInput] = useState<number[]>([]);

    useEffect(() => {
        const fetchChemicals = async () => {
            try {
                const data = await chemicalsClient.listChemicals();
                setChemicals(data.chemicals);
            } catch (err) {
                console.error('Failed to load chemicals:', err);
            }
        };
        fetchChemicals();
    }, []);

    useEffect(() => {
        const fetchForms = async () => {
            try {
                const data = await formsClient.listFormsByUserId(
                    {
                        limit: limit,
                        offset: offset,
                        form_type: formType || null,
                        search_name: searchName || null,
                        sort_by: sortBy || null,
                        order: order || null,
                        chemical_ids: selectedChemicalIds.length > 0 ? selectedChemicalIds : null,
                    }
                );
                setFormviewList(data);
                setError(null);
            } catch (err) {
                let errMessage: Error;
                if (err instanceof AuthError) {
                    errMessage = err as Error
                    setError(new AuthError(errMessage.message))
                } else if (err instanceof Error) {
                    errMessage = err as Error
                    setError(new Error(errMessage.message))
                }
            } finally {
                setIsLoading(false);
            }
        };

        fetchForms();
    }, [limit, offset, formType, searchName, sortBy, order, selectedChemicalIds]);

    const handleApplyFilters = () => {
        setSearchName(searchNameInput);
        setFormType(formTypeInput);
        setSortBy(sortByInput);
        setOrder(orderInput);
        setOffset(0);
    };

    const handleResetFilters = () => {
        setSearchNameInput('');
        setFormTypeInput('');
        setSortByInput('created_at');
        setOrderInput('DESC');
        setSearchName('');
        setFormType('');
        setSortBy('created_at');
        setOrder('DESC');
        setSelectedChemicalIds([]);
        setSelectedChemicalDropdown('');
        setOffset(0);
    };

    const handleAddChemical = () => {
        if (selectedChemicalDropdown) {
            const chemId = parseInt(selectedChemicalDropdown);
            if (!selectedChemicalIds.includes(chemId)) {
                setSelectedChemicalIds([...selectedChemicalIds, chemId]);
            }
            setSelectedChemicalDropdown('');
        }
    };

    const handleRemoveChemical = (chemId: number) => {
        setSelectedChemicalIds(selectedChemicalIds.filter(id => id !== chemId));
    };

    const handleDeleteClick = (form: FormViewResponse) => {
        setFormToDelete(form);
        setShowDeleteConfirm(true);
    };

    const handleDeleteCancel = () => {
        setFormToDelete(null);
        setShowDeleteConfirm(false);
    };

    const handleDeleteConfirm = async () => {
        if (!formToDelete) return;

        setDeletingFormId(formToDelete.id);
        setShowDeleteConfirm(false);

        try {
            await formsClient.deleteForm(formToDelete.id);

            // Refresh the forms list
            const data = await formsClient.listFormsAllUsers({
                offset: offset,
                form_type: formType || null,
                search_name: searchName || null,
                sort_by: sortBy || null,
                order: order || null,
                chemical_ids: chemicalsFilter.length > 0 ? chemicalsFilter : null,
            });
            setFormviewList(data);
            setError(null);
        } catch (err) {
            if (err instanceof AuthError) {
                setError(new AuthError((err as Error).message));
            } else if (err instanceof Error) {
                setError(new Error((err as Error).message));
            }
        } finally {
            setDeletingFormId(null);
            setFormToDelete(null);
        }
    };

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

    if (error instanceof AuthError) {
        router.push('/login')
    }

    if (formviewList == null) {
        return (
            <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950">
                <header className="bg-white dark:bg-zinc-900 shadow">
                    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
                        <div className="flex justify-between items-center">
                            <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">
                                No forms yet!
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
            </div>
        );
    }
    return (
        <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950">
            <header className="bg-white dark:bg-zinc-900 shadow">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
                    <div className="flex justify-between items-center">
                        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">
                            Your Forms
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
                {formviewList && (
                    <div className="space-y-6">
                        {/* Filter and Sort Bar */}
                        <div className="bg-white dark:bg-zinc-900 rounded-lg shadow-md p-6">
                            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-4">
                                {/* Search Name */}
                                <div>
                                    <label htmlFor="searchNameInput" className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
                                        Search Name
                                    </label>
                                    <input
                                        id="searchNameInput"
                                        type="text"
                                        value={searchNameInput}
                                        onChange={(e) => setSearchNameInput(e.target.value)}
                                        placeholder="First or last name..."
                                        className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                    />
                                </div>

                                {/* Form Type Filter */}
                                <div>
                                    <label htmlFor="formTypeInput" className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
                                        Form Type
                                    </label>
                                    <select
                                        id="formTypeInput"
                                        value={formTypeInput}
                                        onChange={(e) => setFormTypeInput(e.target.value)}
                                        className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                    >
                                        <option value="">All Forms</option>
                                        <option value="shrub">Shrub</option>
                                        <option value="lawn">Lawn</option>
                                    </select>
                                </div>

                                {/* Sort By */}
                                <div>
                                    <label htmlFor="sortByInput" className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
                                        Sort By
                                    </label>
                                    <select
                                        id="sortByInput"
                                        value={sortByInput}
                                        onChange={(e) => setSortByInput(e.target.value)}
                                        className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                    >
                                        <option value="created_at">Date Created</option>
                                        <option value="first_name">First Name</option>
                                        <option value="last_name">Last Name</option>
                                    </select>
                                </div>

                                {/* Order Direction */}
                                <div>
                                    <label htmlFor="orderInput" className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
                                        Order
                                    </label>
                                    <select
                                        id="orderInput"
                                        value={orderInput}
                                        onChange={(e) => setOrderInput(e.target.value)}
                                        className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                    >
                                        <option value="DESC">Descending</option>
                                        <option value="ASC">Ascending</option>
                                    </select>
                                </div>
                            </div>

                            {/* Chemical Filter */}
                            <div className="mb-4">
                                <label htmlFor="chemicalFilter" className="block text-sm font-medium text-zinc-700 dark:text-zinc-300 mb-2">
                                    Filter by Chemical
                                </label>
                                <div className="flex gap-2">
                                    <select
                                        id="chemicalFilter"
                                        value={selectedChemicalDropdown}
                                        onChange={(e) => setSelectedChemicalDropdown(e.target.value)}
                                        className="flex-1 px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50 focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                                    >
                                        <option value="">Select a chemical...</option>
                                        {chemicals.map(chem => (
                                            <option key={chem.id} value={chem.id}>
                                                {chem.brand_name} - {chem.chemical_name} ({chem.category})
                                            </option>
                                        ))}
                                    </select>
                                    <button
                                        onClick={handleAddChemical}
                                        disabled={!selectedChemicalDropdown}
                                        className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                                    >
                                        Add
                                    </button>
                                </div>

                                {/* Selected Chemicals */}
                                {selectedChemicalIds.length > 0 && (
                                    <div className="mt-2 flex flex-wrap gap-2">
                                        {selectedChemicalIds.map(chemId => {
                                            const chem = chemicals.find(c => c.id === chemId);
                                            if (!chem) return null;
                                            return (
                                                <div key={chemId} className="inline-flex items-center bg-blue-100 dark:bg-blue-900 text-blue-800 dark:text-blue-200 px-3 py-1 rounded-full text-sm">
                                                    <span>{chem.brand_name} - {chem.chemical_name}</span>
                                                    <button
                                                        onClick={() => handleRemoveChemical(chemId)}
                                                        className="ml-2 text-blue-600 dark:text-blue-300 hover:text-blue-800 dark:hover:text-blue-100"
                                                    >
                                                        ×
                                                    </button>
                                                </div>
                                            );
                                        })}
                                    </div>
                                )}
                            </div>

                            {/* Action Buttons */}
                            <div className="flex flex-wrap gap-3">
                                <button
                                    onClick={handleApplyFilters}
                                    className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-medium"
                                >
                                    Apply Filters
                                </button>
                                <button
                                    onClick={handleResetFilters}
                                    className="px-6 py-2 bg-zinc-200 dark:bg-zinc-700 text-zinc-900 dark:text-zinc-50 rounded-lg hover:bg-zinc-300 dark:hover:bg-zinc-600 transition-colors font-medium"
                                >
                                    Reset
                                </button>
                            </div>
                        </div>

                        {/* Results Info */}
                        <div className="flex items-center justify-between text-sm text-zinc-600 dark:text-zinc-400">
                            <p>
                                Showing <span className="font-medium text-zinc-900 dark:text-zinc-50">{formviewList.forms.length}</span> of{' '}
                                <span className="font-medium text-zinc-900 dark:text-zinc-50">{formviewList.count}</span> forms
                            </p>
                        </div>

                        {/* Forms Grid */}
                        {formviewList.forms.length > 0 ? (
                            <div className="grid grid-cols-1 md:grid-cols-1 gap-4">
                                {formviewList.forms.map((formview: FormViewResponse) => (
                                    <div key={formview.id} className="bg-white dark:bg-zinc-900 rounded-lg shadow p-6">
                                        <div className="flex justify-between items-start">
                                            <div className="grid grid-cols-2 md:grid-cols-3 gap-4 flex-1">
                                                <div>
                                                    <label className="block text-sm font-medium text-zinc-500 dark:text-zinc-400 mb-1">
                                                        Name
                                                    </label>
                                                    <p className="text-zinc-900 dark:text-zinc-50">{formview.first_name} {formview.last_name}</p>
                                                </div>

                                                <div>
                                                    <label className="block text-sm font-medium text-zinc-500 dark:text-zinc-400 mb-1">
                                                        Address
                                                    </label>
                                                    <p className="text-zinc-900 dark:text-zinc-50">
                                                        {formview.street_number} {formview.street_name}
                                                    </p>
                                                </div>

                                                <div>
                                                    <label className="block text-sm font-medium text-zinc-500 dark:text-zinc-400 mb-1">
                                                        Location
                                                    </label>
                                                    <p className="text-zinc-900 dark:text-zinc-50">
                                                        {formview.town}, {formview.zip_code}
                                                    </p>
                                                </div>

                                                <div>
                                                    <label className="block text-sm font-medium text-zinc-500 dark:text-zinc-400 mb-1">
                                                        Form Type
                                                    </label>
                                                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${formview.form_type === 'shrub'
                                                        ? 'bg-blue-100 text-blue-800 dark:bg-blue-900 dark:text-blue-200'
                                                        : 'bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-200'
                                                        }`}>
                                                        {formview.form_type}
                                                    </span>
                                                </div>
                                            </div>
                                        </div>
                                        <div className="flex gap-2 pt-2 border-t border-zinc-200 dark:border-zinc-700">
                                            <button
                                                onClick={() => router.push(`/forms/${formview.form_type}/${formview.id}`)}
                                                className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-sm"
                                            >
                                                View Details
                                            </button>
                                            <button
                                                onClick={() => handleDeleteClick(formview)}
                                                disabled={deletingFormId === formview.id}
                                                className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed text-sm"
                                            >
                                                {deletingFormId === formview.id ? 'Deleting...' : 'Delete'}
                                            </button>
                                        </div>
                                    </div>
                                ))}
                            </div>
                        ) : (
                            <div className="bg-white dark:bg-zinc-900 rounded-lg shadow p-12 text-center">
                                <p className="text-zinc-600 dark:text-zinc-400 text-lg">No forms found matching your filters.</p>
                                <button
                                    onClick={handleResetFilters}
                                    className="mt-4 px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
                                >
                                    Clear Filters
                                </button>
                            </div>
                        )}
                    </div>
                )}
                {/* Delete Confirmation Modal */}
                {showDeleteConfirm && formToDelete && (
                    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
                        <div className="bg-white dark:bg-zinc-900 rounded-lg shadow-xl max-w-md w-full p-6">
                            <h2 className="text-xl font-bold text-zinc-900 dark:text-zinc-50 mb-4">
                                Confirm Deletion
                            </h2>
                            <p className="text-zinc-600 dark:text-zinc-400 mb-6">
                                Are you sure you want to delete the form for <strong>{formToDelete.first_name} {formToDelete.last_name}</strong>?
                                This action cannot be undone.
                            </p>
                            <div className="flex gap-3 justify-end">
                                <button
                                    onClick={handleDeleteCancel}
                                    className="px-4 py-2 bg-zinc-200 dark:bg-zinc-700 text-zinc-900 dark:text-zinc-50 rounded-lg hover:bg-zinc-300 dark:hover:bg-zinc-600 transition-colors"
                                >
                                    Cancel
                                </button>
                                <button
                                    onClick={handleDeleteConfirm}
                                    className="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
                                >
                                    Delete Form
                                </button>
                            </div>
                        </div>
                    </div>
                )}
            </main>
        </div>

    );
}
