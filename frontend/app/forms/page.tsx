'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { formsClient } from '@/lib/api/forms';
import { ListFormsResponse, FormViewResponse, AuthError } from '@/lib/api/types';

export default function ListFormsPage() {
    const router = useRouter();

    const [formviewList, setFormviewList] = useState<ListFormsResponse | null>(null);
    const [error, setError] = useState<Error | AuthError | null>(null);
    const [isLoading, setIsLoading] = useState(true);

    //query string params
    const [limit, setLimit] = useState<number>(10);
    const [currentPage, setCurrentPage] = useState<number>(1);
    const [offset, setOffset] = useState<number>(0);

    const [formType, setFormType] = useState<string>('');
    const [formTypeInput, setFormTypeInput] = useState<string>('');

    const [searchName, setSearchName] = useState<string>('');
    const [searchNameInput, setSearchNameInput] = useState<string>('');

    const [sortBy, setSortBy] = useState<string>('created_at');
    const [sortByInput, setSortByInput] = useState<string>('created_at');

    const [order, setOrder] = useState<string>('DESC');
    const [orderInput, setOrderInput] = useState<string>('DESC');

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
    }, [limit, offset, formType, searchName, sortBy, order]);

    const handleApplyFilters = () => {
        setSearchName(searchNameInput);
        setFormType(formTypeInput);
        setSortBy(sortByInput);
        setOrder(orderInput);
        setCurrentPage(1);
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
        setCurrentPage(1);
        setOffset(0);
    };

    const handlePageChange = (newPage: number) => {
        setCurrentPage(newPage);
        setOffset((newPage - 1) * limit);
    };

    const handleLimitChange = (newLimit: number) => {
        setLimit(newLimit);
        setCurrentPage(1);
        setOffset(0);
    };

    const totalPages = formviewList ? Math.ceil(formviewList.count / limit) : 0;

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
        return (<div>Not available! FormviewList is null</div>);
    }
    return (
        <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950">
            <header className="bg-white dark:bg-zinc-900 shadow">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
                    <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">
                        Your Forms
                    </h1>
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
                                        <option value="pesticide">Pesticide</option>
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
