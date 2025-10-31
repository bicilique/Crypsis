// Security Settings Component
export function SecuritySettings() {
  return (
    <div className="card">
      <h2 className="text-xl font-semibold text-gray-900 mb-4">Security Settings</h2>
      <div className="space-y-4">
        <div className="flex items-center justify-between py-3 border-b">
          <div>
            <h3 className="text-sm font-medium text-gray-900">Encryption Algorithm</h3>
            <p className="text-sm text-gray-500">AES-256-GCM (Galois/Counter Mode)</p>
          </div>
          <span className="text-sm font-medium text-green-600">Active</span>
        </div>

        <div className="flex items-center justify-between py-3 border-b">
          <div>
            <h3 className="text-sm font-medium text-gray-900">Key Storage</h3>
            <p className="text-sm text-gray-500">File-based key management</p>
          </div>
          <span className="text-sm font-medium text-green-600">Configured</span>
        </div>

        <div className="flex items-center justify-between py-3 border-b">
          <div>
            <h3 className="text-sm font-medium text-gray-900">KMS Integration</h3>
            <p className="text-sm text-gray-500">External Key Management Service</p>
          </div>
          <span className="text-sm font-medium text-gray-500">Optional</span>
        </div>

        <div className="flex items-center justify-between py-3">
          <div>
            <h3 className="text-sm font-medium text-gray-900">OAuth2 Authentication</h3>
            <p className="text-sm text-gray-500">Hydra OAuth2 Server</p>
          </div>
          <span className="text-sm font-medium text-green-600">Enabled</span>
        </div>
      </div>
    </div>
  );
}
