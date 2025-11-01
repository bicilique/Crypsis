// Application Form Component
import { useState } from 'react';
import { X } from 'lucide-react';
import { useApplicationStore } from '@/stores';

interface ApplicationFormProps {
  onClose: () => void;
}

export function ApplicationForm({ onClose }: ApplicationFormProps) {
  const [formData, setFormData] = useState({
    name: '',
    uri: '',
    redirectUri: '',
  });
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const { createApplication } = useApplicationStore();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);

    try {
      await createApplication(formData);
      alert('Application created successfully! Please save the Client ID and Secret.');
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create application');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl max-w-md w-full mx-4">
        <div className="flex items-center justify-between p-6 border-b">
          <h3 className="text-lg font-semibold text-gray-900">Create Application</h3>
          <button onClick={onClose} className="text-gray-400 hover:text-gray-600">
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="p-6">
          <div className="space-y-4">
            <div>
              <label className="label">Application Name</label>
              <input
                type="text"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                className="input"
                required
                placeholder="My Application"
              />
            </div>

            <div>
              <label className="label">Homepage URI</label>
              <input
                type="url"
                value={formData.uri}
                onChange={(e) => setFormData({ ...formData, uri: e.target.value })}
                className="input"
                required
                placeholder="https://myapp.example.com"
              />
            </div>

            <div>
              <label className="label">Redirect URI</label>
              <input
                type="url"
                value={formData.redirectUri}
                onChange={(e) => setFormData({ ...formData, redirectUri: e.target.value })}
                className="input"
                required
                placeholder="https://myapp.example.com/callback"
              />
              <p className="mt-1 text-xs text-gray-500">OAuth2 callback URL</p>
            </div>

            {error && (
              <div className="p-3 bg-red-50 border border-red-200 rounded text-sm text-red-600">
                {error}
              </div>
            )}
          </div>

          <div className="flex items-center justify-end gap-3 mt-6">
            <button type="button" onClick={onClose} className="btn-secondary px-4 py-2 text-sm">
              Cancel
            </button>
            <button
              type="submit"
              disabled={loading}
              className="btn-primary px-4 py-2 text-sm disabled:opacity-50"
            >
              {loading ? 'Creating...' : 'Create Application'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
