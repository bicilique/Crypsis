/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_URL: string
  readonly VITE_KMS_ENABLED: string
  readonly VITE_BULK_OPS_ENABLED: string
  readonly VITE_ADMIN_PANEL_ENABLED: string
  readonly VITE_FILE_PREVIEW_ENABLED: string
  readonly VITE_NOTIFICATIONS_ENABLED: string
  readonly VITE_APP_VERSION: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}
