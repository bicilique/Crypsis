// Rekey Form Component
import { useState } from 'react';
import { X } from 'lucide-react';
import { securityService } from '@/services';

interface RekeyFormProps {
  onClose: () => void;
}

export function RekeyForm({ onClose }: RekeyFormProps) {
  const [keyUID, setKeyUID] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!confirm('This will re-encrypt all files with a new key. This operation cannot be undone. Continue?')) {
      return;
    }

    setLoading(true);
    setError(null);

    try {
      await securityService.rekeyFiles(keyUID);
      alert('Re-keying completed successfully!');
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to re-key files');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl max-w-md w-full mx-4">
        <div className="flex items-center justify-between p-6 border-b">
          <h3 className="text-lg font-semibold text-gray-900">Re-key All Files</h3>
          <button onClick={onClose} className="text-gray-400 hover:text-gray-600">
            <X className="w-5 h-5" />
          </button>
        </div>

        <form onSubmit={handleSubmit} className="p-6">
          <div className="space-y-4">
            <div className="bg-yellow-50 border border-yellow-200 rounded p-4">
              <p className="text-sm text-yellow-800">
                <strong>Warning:</strong> This operation will re-encrypt all files with a new encryption key. 
                Make sure you have the correct Key UID before proceeding.
              </p>
            </div>

            <div>
              <label className="label">New Key UID</label>
              <input
                type="text"
                value={keyUID}
                onChange={(e) => setKeyUID(e.target.value)}
                className="input"
                required
                placeholder="Enter new key UID"
              />
              <p className="mt-1 text-xs text-gray-500">
                This is the unique identifier for your new encryption key
              </p>
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
              className="btn-danger px-4 py-2 text-sm disabled:opacity-50"
            >
              {loading ? 'Re-keying...' : 'Start Re-keying'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
