import { useState } from 'react';
import { RekeyForm, SecuritySettings } from '@/components/features/security';
import { Shield, Key, Lock, AlertTriangle } from 'lucide-react';

export function SecurityPage() {
  const [isRekeyModalOpen, setIsRekeyModalOpen] = useState(false);

  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Security Center</h1>
        <p className="mt-2 text-sm text-gray-600">
          Manage encryption keys and security settings
        </p>
      </div>

      {/* Security Overview */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div className="card">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <Shield className="w-8 h-8 text-success-600" />
            </div>
            <div className="ml-4">
              <div className="text-sm font-medium text-gray-600">Encryption Status</div>
              <div className="mt-1 text-lg font-semibold text-success-600">Active</div>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <Key className="w-8 h-8 text-primary-600" />
            </div>
            <div className="ml-4">
              <div className="text-sm font-medium text-gray-600">Encryption Method</div>
              <div className="mt-1 text-lg font-semibold text-gray-900">AES-256-GCM</div>
            </div>
          </div>
        </div>

        <div className="card">
          <div className="flex items-center">
            <div className="flex-shrink-0">
              <Lock className="w-8 h-8 text-primary-600" />
            </div>
            <div className="ml-4">
              <div className="text-sm font-medium text-gray-600">KMS Integration</div>
              <div className="mt-1 text-lg font-semibold text-gray-900">Enabled</div>
            </div>
          </div>
        </div>
      </div>

      {/* Security Alert */}
      <div className="card bg-yellow-50 border-yellow-200 mb-8">
        <div className="flex items-start">
          <AlertTriangle className="w-6 h-6 text-yellow-600 mt-0.5" />
          <div className="ml-3">
            <h3 className="text-sm font-medium text-yellow-800">Security Notice</h3>
            <p className="mt-2 text-sm text-yellow-700">
              Regular key rotation is recommended for enhanced security. Last rotation was 30 days ago.
            </p>
          </div>
        </div>
      </div>

      {/* Key Management */}
      <div className="card mb-8">
        <h2 className="text-xl font-semibold text-gray-900 mb-4">Key Management</h2>
        <div className="space-y-4">
          <div className="flex items-center justify-between py-4 border-b border-gray-200">
            <div>
              <h3 className="text-sm font-medium text-gray-900">Re-encrypt Files</h3>
              <p className="mt-1 text-sm text-gray-600">
                Re-encrypt all files with a new encryption key
              </p>
            </div>
            <button
              onClick={() => setIsRekeyModalOpen(true)}
              className="btn-primary px-4 py-2 text-sm"
            >
              <Key className="w-4 h-4 mr-2" />
              Start Re-keying
            </button>
          </div>

          <div className="flex items-center justify-between py-4 border-b border-gray-200">
            <div>
              <h3 className="text-sm font-medium text-gray-900">Hash Method</h3>
              <p className="mt-1 text-sm text-gray-600">
                File integrity verification method
              </p>
            </div>
            <span className="text-sm font-medium text-gray-900">SHA-256</span>
          </div>

          <div className="flex items-center justify-between py-4">
            <div>
              <h3 className="text-sm font-medium text-gray-900">Hash Encrypted Files</h3>
              <p className="mt-1 text-sm text-gray-600">
                Calculate hash after encryption
              </p>
            </div>
            <span className="text-sm font-medium text-success-600">Enabled</span>
          </div>
        </div>
      </div>

      {/* Security Settings */}
      <SecuritySettings />

      {/* Rekey Modal */}
      {isRekeyModalOpen && (
        <RekeyForm onClose={() => setIsRekeyModalOpen(false)} />
      )}
    </div>
  );
}
