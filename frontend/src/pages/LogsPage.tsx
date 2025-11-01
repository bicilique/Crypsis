import { useState, useEffect } from 'react';
import { LogsList, LogsFilter } from '@/components/features/logs';
import { useLogsStore } from '@/stores';
import { LoadingSpinner } from '@/components/ui';
import { FileText, Download } from 'lucide-react';

export function LogsPage() {
  const { logs, isLoading, fetchLogs } = useLogsStore();
  const [filters, setFilters] = useState({
    action: '',
    actor_type: '',
    date_from: '',
    date_to: '',
  });

  useEffect(() => {
    fetchLogs();
  }, [fetchLogs]);

  const handleExport = () => {
    const dataStr = JSON.stringify(logs, null, 2);
    const blob = new Blob([dataStr], { type: 'application/json' });
    const url = URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    link.download = `audit-logs-${new Date().toISOString()}.json`;
    link.click();
    URL.revokeObjectURL(url);
  };

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Audit Logs</h1>
            <p className="mt-2 text-sm text-gray-600">
              Complete activity tracking and monitoring
            </p>
          </div>
          <button
            onClick={handleExport}
            className="btn-secondary px-4 py-2 text-sm"
            disabled={logs.length === 0}
          >
            <Download className="w-4 h-4 mr-2" />
            Export Logs
          </button>
        </div>
      </div>

      {/* Statistics */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div className="card">
          <div className="text-sm font-medium text-gray-600">Total Events</div>
          <div className="mt-2 text-3xl font-bold text-gray-900">{logs.length}</div>
        </div>
        <div className="card">
          <div className="text-sm font-medium text-gray-600">Today</div>
          <div className="mt-2 text-3xl font-bold text-primary-600">
            {logs.filter(log => {
              const logDate = new Date(log.timestamp);
              const today = new Date();
              return logDate.toDateString() === today.toDateString();
            }).length}
          </div>
        </div>
        <div className="card">
          <div className="text-sm font-medium text-gray-600">File Operations</div>
          <div className="mt-2 text-3xl font-bold text-gray-900">
            {logs.filter(log => ['upload', 'download', 'delete', 'update'].includes(log.action)).length}
          </div>
        </div>
        <div className="card">
          <div className="text-sm font-medium text-gray-600">Admin Actions</div>
          <div className="mt-2 text-3xl font-bold text-gray-900">
            {logs.filter(log => log.actor_type === 'admin').length}
          </div>
        </div>
      </div>

      {/* Filters */}
      <div className="mb-6">
        <LogsFilter filters={filters} onFilterChange={setFilters} />
      </div>

      {/* Logs List */}
      {isLoading ? (
        <div className="flex items-center justify-center py-12">
          <LoadingSpinner size="lg" />
        </div>
      ) : (
        <LogsList logs={logs} />
      )}
    </div>
  );
}
