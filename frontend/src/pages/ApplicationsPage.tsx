import { useState, useEffect } from 'react';
import { ApplicationList, ApplicationForm } from '@/components/features/applications';
import { useApplicationStore } from '@/stores';
import { LoadingSpinner } from '@/components/ui';
import { Plus } from 'lucide-react';

export function ApplicationsPage() {
  const [isFormOpen, setIsFormOpen] = useState(false);
  const { applications, isLoading, fetchApplications } = useApplicationStore();

  useEffect(() => {
    fetchApplications();
  }, [fetchApplications]);

  const activeApps = applications.filter(app => app.is_active);
  const inactiveApps = applications.filter(app => !app.is_active);

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Applications</h1>
            <p className="mt-2 text-sm text-gray-600">
              Manage OAuth2 client applications and API access
            </p>
          </div>
          <button
            onClick={() => setIsFormOpen(true)}
            className="btn-primary px-4 py-2 text-sm"
          >
            <Plus className="w-4 h-4 mr-2" />
            Create Application
          </button>
        </div>
      </div>

      {/* Statistics */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div className="card">
          <div className="text-sm font-medium text-gray-600">Total Apps</div>
          <div className="mt-2 text-3xl font-bold text-gray-900">{applications.length}</div>
        </div>
        <div className="card">
          <div className="text-sm font-medium text-gray-600">Active</div>
          <div className="mt-2 text-3xl font-bold text-success-600">{activeApps.length}</div>
        </div>
        <div className="card">
          <div className="text-sm font-medium text-gray-600">Inactive</div>
          <div className="mt-2 text-3xl font-bold text-gray-400">{inactiveApps.length}</div>
        </div>
        <div className="card">
          <div className="text-sm font-medium text-gray-600">This Month</div>
          <div className="mt-2 text-sm text-gray-600">
            {applications.filter(a => {
              const createdDate = new Date(a.created_at);
              const now = new Date();
              return createdDate.getMonth() === now.getMonth() && 
                     createdDate.getFullYear() === now.getFullYear();
            }).length} new
          </div>
        </div>
      </div>

      {/* Application List */}
      {isLoading ? (
        <div className="flex items-center justify-center py-12">
          <LoadingSpinner size="lg" />
        </div>
      ) : (
        <ApplicationList />
      )}

      {/* Create Application Modal */}
      {isFormOpen && (
        <ApplicationForm onClose={() => setIsFormOpen(false)} />
      )}
    </div>
  );
}
