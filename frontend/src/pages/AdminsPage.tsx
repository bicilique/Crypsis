import { useState, useEffect } from 'react';
import { AdminList, AdminForm } from '@/components/features/admin';
import { useAdminStore } from '@/stores';
import { LoadingSpinner } from '@/components/ui';
import { UserPlus } from 'lucide-react';

export function AdminsPage() {
  const [isFormOpen, setIsFormOpen] = useState(false);
  const { admins, isLoading, fetchAdmins } = useAdminStore();

  useEffect(() => {
    fetchAdmins();
  }, [fetchAdmins]);

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Admin Management</h1>
            <p className="mt-2 text-sm text-gray-600">
              Manage administrator accounts and permissions
            </p>
          </div>
          <button
            onClick={() => setIsFormOpen(true)}
            className="btn-primary px-4 py-2 text-sm"
          >
            <UserPlus className="w-4 h-4 mr-2" />
            Add Admin
          </button>
        </div>
      </div>

      {/* Statistics */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div className="card">
          <div className="text-sm font-medium text-gray-600">Total Admins</div>
          <div className="mt-2 text-3xl font-bold text-gray-900">{admins.length}</div>
        </div>
        <div className="card">
          <div className="text-sm font-medium text-gray-600">Active</div>
          <div className="mt-2 text-3xl font-bold text-success-600">
            {admins.filter(a => a.is_active !== false).length}
          </div>
        </div>
        <div className="card">
          <div className="text-sm font-medium text-gray-600">Recent Activity</div>
          <div className="mt-2 text-sm text-gray-600">Last 24 hours</div>
        </div>
      </div>

      {/* Admin List */}
      {isLoading ? (
        <div className="flex items-center justify-center py-12">
          <LoadingSpinner size="lg" />
        </div>
      ) : (
        <AdminList />
      )}

      {/* Add Admin Modal */}
      {isFormOpen && (
        <AdminForm onClose={() => setIsFormOpen(false)} />
      )}
    </div>
  );
}
