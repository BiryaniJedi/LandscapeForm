'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { chemicalsClient } from '@/lib/api/chemicals';
import { Chemical, CreateChemicalRequest } from '@/lib/api/types';

export default function AdminChemicalsPage() {
    const router = useRouter();
    const [chemicals, setChemicals] = useState<Chemical[]>([]);
    const [error, setError] = useState<string | null>(null);
    const [isLoadingChemicals, setIsLoadingChemicals] = useState(true);
    const [isLoading, setIsLoading] = useState(true);
    const [actionInProgress, setActionInProgress] = useState<number | null>(null);
    const [categoryFilter, setCategoryFilter] = useState<'lawn' | 'shrub' | ''>('');
    const [showCreateForm, setShowCreateForm] = useState(false);
    const [editingChemical, setEditingChemical] = useState<Chemical | null>(null);
    const [formData, setFormData] = useState<CreateChemicalRequest>({
        category: 'lawn',
        brand_name: '',
        chemical_name: '',
        epa_reg_no: '',
        recipe: '',
        unit: '',
    });

    const fetchChemicals = async (category?: 'lawn' | 'shrub') => {
        setIsLoadingChemicals(true);
        try {
            const data = await chemicalsClient.listChemicals(category);
            setChemicals(data.chemicals);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to load chemicals');
        } finally {
            setIsLoading(false)
            setIsLoadingChemicals(false);
        }
    };

    useEffect(() => {
        fetchChemicals();
    }, []);

    const handleCategoryFilterChange = (category: 'lawn' | 'shrub' | '') => {
        setCategoryFilter(category);
        fetchChemicals(category || undefined);
    };

    const handleCreateChemical = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        try {
            await chemicalsClient.createChemical(formData);
            setShowCreateForm(false);
            resetForm();
            fetchChemicals(categoryFilter || undefined);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to create chemical');
        }
    };

    const handleUpdateChemical = async (e: React.FormEvent) => {
        e.preventDefault();
        if (!editingChemical) return;

        setError(null);
        try {
            await chemicalsClient.updateChemical(editingChemical.id, formData);
            setEditingChemical(null);
            resetForm();
            fetchChemicals(categoryFilter || undefined);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to update chemical');
        }
    };

    const handleDeleteChemical = async (id: number) => {
        if (!confirm('Are you sure you want to delete this chemical? This action cannot be undone.')) {
            return;
        }

        setActionInProgress(id);
        setError(null);
        try {
            await chemicalsClient.deleteChemical(id);
            fetchChemicals(categoryFilter || undefined);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to delete chemical');
        } finally {
            setActionInProgress(null);
            setIsLoading(false);
        }
    };

    const startEdit = (chemical: Chemical) => {
        setEditingChemical(chemical);
        setFormData({
            category: chemical.category,
            brand_name: chemical.brand_name,
            chemical_name: chemical.chemical_name,
            epa_reg_no: chemical.epa_reg_no,
            recipe: chemical.recipe,
            unit: chemical.unit,
        });
        setShowCreateForm(false);
    };

    const resetForm = () => {
        setFormData({
            category: 'lawn',
            brand_name: '',
            chemical_name: '',
            epa_reg_no: '',
            recipe: '',
            unit: '',
        });
        setEditingChemical(null);
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

    if (error) {
        return (
            <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950">
                <header className="bg-white dark:bg-zinc-900 shadow">
                    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
                        <div className="flex justify-between items-center">
                            <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">
                                Error Loading Chemicals
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
                <main className="max-w-2xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                    <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4">
                        <p className="text-red-800 dark:text-red-200">{error}</p>
                    </div>
                </main>
            </div>
        );
    }

    return (
        <div className="min-h-screen bg-zinc-50 dark:bg-zinc-950">
            <header className="bg-white dark:bg-zinc-900 shadow">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex items-center justify-between">
                    <div>
                        <h1 className="text-2xl font-bold text-zinc-900 dark:text-zinc-50">
                            Chemicals Management
                        </h1>
                        <p className="text-sm text-zinc-600 dark:text-zinc-400 mt-1">
                            Manage lawn and shrub chemicals
                        </p>
                    </div>
                    <button
                        onClick={() => router.push('/dashboard')}
                        className="px-4 py-2 bg-zinc-200 dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50 rounded-lg hover:bg-zinc-300 dark:hover:bg-zinc-700 transition-colors"
                    >
                        Back to Dashboard
                    </button>
                </div>
            </header>

            <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
                {error && (
                    <div className="bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-lg p-4 mb-6">
                        <p className="text-red-800 dark:text-red-200">{error}</p>
                    </div>
                )}

                {/* Controls */}
                <div className="bg-white dark:bg-zinc-900 rounded-lg shadow p-4 mb-6">
                    <div className="flex items-center justify-between gap-4">
                        <div className="flex items-center gap-4">
                            <label className="text-sm font-medium text-zinc-900 dark:text-zinc-50">
                                Filter by Category:
                            </label>
                            <select
                                value={categoryFilter}
                                onChange={(e) => handleCategoryFilterChange(e.target.value as 'lawn' | 'shrub' | '')}
                                className="px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50"
                            >
                                <option value="">All Categories</option>
                                <option value="lawn">Lawn</option>
                                <option value="shrub">Shrub</option>
                            </select>
                        </div>
                        <button
                            onClick={() => {
                                setShowCreateForm(true);
                                setEditingChemical(null);
                                resetForm();
                            }}
                            className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
                        >
                            Add New Chemical
                        </button>
                    </div>
                </div>

                {/* Create/Edit Form */}
                {(showCreateForm || editingChemical) && (
                    <div className="bg-white dark:bg-zinc-900 rounded-lg shadow p-6 mb-6">
                        <h2 className="text-xl font-semibold text-zinc-900 dark:text-zinc-50 mb-4">
                            {editingChemical ? 'Edit Chemical' : 'Create New Chemical'}
                        </h2>
                        <form onSubmit={editingChemical ? handleUpdateChemical : handleCreateChemical} className="space-y-4">
                            <div>
                                <label className="block text-sm font-medium text-zinc-900 dark:text-zinc-50 mb-2">
                                    Category *
                                </label>
                                <select
                                    value={formData.category}
                                    onChange={(e) => setFormData({ ...formData, category: e.target.value as 'lawn' | 'shrub' })}
                                    required
                                    className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50"
                                >
                                    <option value="lawn">Lawn</option>
                                    <option value="shrub">Shrub</option>
                                </select>
                            </div>

                            <div className="grid grid-cols-2 gap-4">
                                <div>
                                    <label className="block text-sm font-medium text-zinc-900 dark:text-zinc-50 mb-2">
                                        Brand Name *
                                    </label>
                                    <input
                                        type="text"
                                        value={formData.brand_name}
                                        onChange={(e) => setFormData({ ...formData, brand_name: e.target.value })}
                                        required
                                        className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50"
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-zinc-900 dark:text-zinc-50 mb-2">
                                        Chemical Name *
                                    </label>
                                    <input
                                        type="text"
                                        value={formData.chemical_name}
                                        onChange={(e) => setFormData({ ...formData, chemical_name: e.target.value })}
                                        required
                                        className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50"
                                    />
                                </div>
                            </div>

                            <div className="grid grid-cols-3 gap-4">
                                <div>
                                    <label className="block text-sm font-medium text-zinc-900 dark:text-zinc-50 mb-2">
                                        EPA Reg No
                                    </label>
                                    <input
                                        type="text"
                                        value={formData.epa_reg_no}
                                        onChange={(e) => setFormData({ ...formData, epa_reg_no: e.target.value })}
                                        className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50"
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-zinc-900 dark:text-zinc-50 mb-2">
                                        Recipe
                                    </label>
                                    <input
                                        type="text"
                                        value={formData.recipe}
                                        onChange={(e) => setFormData({ ...formData, recipe: e.target.value })}
                                        className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50"
                                    />
                                </div>
                                <div>
                                    <label className="block text-sm font-medium text-zinc-900 dark:text-zinc-50 mb-2">
                                        Unit
                                    </label>
                                    <input
                                        type="text"
                                        value={formData.unit}
                                        onChange={(e) => setFormData({ ...formData, unit: e.target.value })}
                                        className="w-full px-3 py-2 border border-zinc-300 dark:border-zinc-700 rounded-lg bg-white dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50"
                                    />
                                </div>
                            </div>

                            <div className="flex gap-4">
                                <button
                                    type="submit"
                                    className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
                                >
                                    {editingChemical ? 'Update Chemical' : 'Create Chemical'}
                                </button>
                                <button
                                    type="button"
                                    onClick={resetForm}
                                    className="px-4 py-2 bg-zinc-200 dark:bg-zinc-800 text-zinc-900 dark:text-zinc-50 rounded-lg hover:bg-zinc-300 dark:hover:bg-zinc-700 transition-colors"
                                >
                                    Cancel
                                </button>
                            </div>
                        </form>
                    </div>
                )}

                {/* Chemicals List */}
                {isLoadingChemicals ? (
                    <div className="flex items-center justify-center py-12">
                        <div className="text-center">
                            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto"></div>
                            <p className="mt-4 text-zinc-600 dark:text-zinc-400">Loading chemicals...</p>
                        </div>
                    </div>
                ) : chemicals.length === 0 ? (
                    <div className="bg-white dark:bg-zinc-900 rounded-lg shadow p-12 text-center">
                        <p className="text-zinc-600 dark:text-zinc-400 text-lg">No chemicals found.</p>
                    </div>
                ) : (
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
                                        <th className="px-6 py-3 text-right text-xs font-medium text-zinc-500 dark:text-zinc-400 uppercase tracking-wider">
                                            Actions
                                        </th>
                                    </tr>
                                </thead>
                                <tbody className="bg-white dark:bg-zinc-900 divide-y divide-zinc-200 dark:divide-zinc-800">
                                    {chemicals.map((chem) => (
                                        <tr key={chem.id}>
                                            <td className="px-6 py-4 whitespace-nowrap text-sm text-zinc-900 dark:text-zinc-50">
                                                {chem.id}
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
                                            <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium space-x-2">
                                                <button
                                                    onClick={() => startEdit(chem)}
                                                    className="text-blue-600 hover:text-blue-900 dark:text-blue-400 dark:hover:text-blue-300"
                                                >
                                                    Edit
                                                </button>
                                                <button
                                                    onClick={() => handleDeleteChemical(chem.id)}
                                                    disabled={actionInProgress === chem.id}
                                                    className="text-red-600 hover:text-red-900 dark:text-red-400 dark:hover:text-red-300 disabled:opacity-50 disabled:cursor-not-allowed"
                                                >
                                                    {actionInProgress === chem.id ? 'Deleting...' : 'Delete'}
                                                </button>
                                            </td>
                                        </tr>
                                    ))}
                                </tbody>
                            </table>
                        </div>
                    </div>
                )}
            </main>
        </div>
    );
}
