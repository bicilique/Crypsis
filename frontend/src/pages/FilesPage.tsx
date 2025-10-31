import { useState, useEffect } from 'react';
import { FileList, FileUpload } from '@/components/features/files';
import { useFileStore } from '@/stores/files';
import { LoadingSpinner } from '@/components/ui';
import { Upload } from 'lucide-react';
import { formatBytes } from '@/utils';

export function FilesPage() {
  const [isUploadModalOpen, setIsUploadModalOpen] = useState(false);
  const { files, isLoading, fetchFiles } = useFileStore();

  useEffect(() => {
    fetchFiles();
  }, [fetchFiles]);

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-bold text-gray-900">Files</h1>
            <p className="mt-2 text-sm text-gray-600">
              Manage encrypted files with enterprise-level security
            </p>
          </div>
          <button
            onClick={() => setIsUploadModalOpen(true)}
            className="btn-primary px-4 py-2 text-sm"
          >
            <Upload className="w-4 h-4 mr-2" />
            Upload File
          </button>
        </div>
      </div>

      {/* File Statistics */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6 mb-8">
        <div className="card">
          <div className="text-sm font-medium text-gray-600">Total Files</div>
          <div className="mt-2 text-3xl font-bold text-gray-900">{files.length}</div>
        </div>
        <div className="card">
          <div className="text-sm font-medium text-gray-600">Total Size</div>
          <div className="mt-2 text-3xl font-bold text-gray-900">
            {formatBytes(files.reduce((acc, f) => acc + (f.file_size || 0), 0))}
          </div>
        </div>
        <div className="card">
          <div className="text-sm font-medium text-gray-600">Active</div>
          <div className="mt-2 text-3xl font-bold text-success-600">
            {files.filter(f => !f.deleted).length}
          </div>
        </div>
        <div className="card">
          <div className="text-sm font-medium text-gray-600">Deleted</div>
          <div className="mt-2 text-3xl font-bold text-gray-500">
            {files.filter(f => f.deleted).length}
          </div>
        </div>
      </div>

      {/* Files List */}
      {isLoading ? (
        <div className="flex items-center justify-center py-12">
          <LoadingSpinner size="lg" />
        </div>
      ) : (
        <FileList />
      )}

      {/* Upload Modal */}
      {isUploadModalOpen && (
        <FileUpload onClose={() => setIsUploadModalOpen(false)} />
      )}
    </div>
  );
}
