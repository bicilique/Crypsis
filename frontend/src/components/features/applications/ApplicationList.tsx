// Application List Component
import { Trash2, RotateCw, Eye, EyeOff, Copy } from 'lucide-react';
import { useState } from 'react';
import { useApplicationStore, type Application } from '@/stores';
import { formatDate } from '@/utils';

export function ApplicationList() {
  const { applications, deleteApplication, rotateSecret } = useApplicationStore();
  const [showSecret, setShowSecret] = useState<Record<string, boolean>>({});

  const handleDelete = async (id: string) => {
    if (confirm('Are you sure you want to delete this application?')) {
      try {
        await deleteApplication(id);
      } catch (error) {
        alert('Failed to delete application');
      }
    }
  };

  const handleRotateSecret = async (id: string) => {
    if (confirm('This will generate a new secret. The old secret will stop working. Continue?')) {
      try {
        await rotateSecret(id);
        alert('Secret rotated successfully. Please update your application configuration.');
      } catch (error) {
        alert('Failed to rotate secret');
      }
    }
  };

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text);
    alert('Copied to clipboard!');
  };

  if (applications.length === 0) {
    return (
      <div className="text-center py-12">
        <p className="text-gray-500">No applications found. Create your first application to get started.</p>
      </div>
    );
  }

  return (
    <div className="bg-white shadow-sm rounded-lg overflow-hidden">
      <div className="grid gap-4 p-4">
        {applications.map((app) => (
          <div key={app.id} className="border rounded-lg p-4 hover:shadow-md transition-shadow">
            <div className="flex items-start justify-between">
              <div className="flex-1">
                <h3 className="text-lg font-semibold text-gray-900">{app.app_name}</h3>
                <div className="mt-2 space-y-2">
                  <div className="flex items-center gap-2">
                    <span className="text-sm text-gray-500">Client ID:</span>
                    <code className="text-sm bg-gray-100 px-2 py-1 rounded">{app.client_id}</code>
                    <button
                      onClick={() => copyToClipboard(app.client_id)}
                      className="text-gray-400 hover:text-gray-600"
                    >
                      <Copy className="w-4 h-4" />
                    </button>
                  </div>
                  
                  {app.client_secret && (
                    <div className="flex items-center gap-2">
                      <span className="text-sm text-gray-500">Client Secret:</span>
                      <code className="text-sm bg-gray-100 px-2 py-1 rounded">
                        {showSecret[app.id] ? app.client_secret : '••••••••••••••••'}
                      </code>
                      <button
                        onClick={() => setShowSecret({ ...showSecret, [app.id]: !showSecret[app.id] })}
                        className="text-gray-400 hover:text-gray-600"
                      >
                        {showSecret[app.id] ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                      </button>
                      {showSecret[app.id] && (
                        <button
                          onClick={() => copyToClipboard(app.client_secret!)}
                          className="text-gray-400 hover:text-gray-600"
                        >
                          <Copy className="w-4 h-4" />
                        </button>
                      )}
                    </div>
                  )}

                  <div className="text-sm text-gray-500">
                    Created: {formatDate(app.created_at)}
                  </div>

                  <div>
                    {app.is_active ? (
                      <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                        Active
                      </span>
                    ) : (
                      <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                        Inactive
                      </span>
                    )}
                  </div>
                </div>
              </div>

              <div className="flex items-center gap-2">
                <button
                  onClick={() => handleRotateSecret(app.id)}
                  className="text-primary-600 hover:text-primary-900"
                  title="Rotate Secret"
                >
                  <RotateCw className="w-4 h-4" />
                </button>
                <button
                  onClick={() => handleDelete(app.id)}
                  className="text-red-600 hover:text-red-900"
                  title="Delete"
                >
                  <Trash2 className="w-4 h-4" />
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
