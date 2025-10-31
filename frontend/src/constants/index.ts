// API Configuration
export const API_CONFIG = {
  BASE_URL: import.meta.env.VITE_API_URL || 'http://host.docker.internal:8080',
  TIMEOUT: 30000,
  RETRY_ATTEMPTS: 3,
  RETRY_DELAY: 1000,
} as const;

// API Endpoints
export const API_ENDPOINTS = {
  // Authentication
  LOGIN: '/api/admin/login',
  LOGOUT: '/api/admin/logout',
  REFRESH_TOKEN: '/api/admin/refresh-token',
  
  // File Operations
  FILES: '/api/files',
  FILE_UPLOAD: '/api/files',
  FILE_DOWNLOAD: (id: string) => `/api/files/${id}/download`,
  FILE_UPDATE: (id: string) => `/api/files/${id}/update`,
  FILE_DELETE: (id: string) => `/api/files/${id}/delete`,
  FILE_METADATA: (id: string) => `/api/files/${id}/metadata`,
  FILE_ENCRYPT: '/api/files/encrypt',
  FILE_DECRYPT: '/api/files/decrypt',
  FILE_LIST: '/api/files/list',
  
  // Admin Operations
  ADMIN_LIST: '/api/admin/list',
  ADMIN_ADD: '/api/admin/add',
  ADMIN_UPDATE_USERNAME: '/api/admin/username',
  ADMIN_UPDATE_PASSWORD: '/api/admin/password',
  ADMIN_DELETE: '/api/admin',
  
  // Application Management
  APPS: '/api/admin/apps',
  APP_DETAIL: (id: string) => `/api/admin/apps/${id}`,
  APP_DELETE: (id: string) => `/api/admin/apps/${id}`,
  APP_ROTATE_SECRET: (id: string) => `/api/admin/apps/${id}/rotate-secret`,
  
  // Admin File Operations
  ADMIN_FILES: '/api/admin/files',
  ADMIN_LOGS: '/api/admin/logs',
  ADMIN_REKEY: '/api/admin/files/re-key',
} as const;

// Application Routes
export const ROUTES = {
  HOME: '/',
  LOGIN: '/login',
  DASHBOARD: '/dashboard',
  FILES: '/files',
  ADMIN: '/admin',
  ADMIN_USERS: '/admin/users',
  ADMIN_APPS: '/admin/apps',
  ADMIN_LOGS: '/admin/logs',
  SECURITY: '/security',
  SETTINGS: '/settings',
  NOT_FOUND: '*',
} as const;

// Local Storage Keys
export const STORAGE_KEYS = {
  ACCESS_TOKEN: 'ecrypt_access_token',
  REFRESH_TOKEN: 'ecrypt_refresh_token',
  USER_DATA: 'ecrypt_user_data',
  PREFERENCES: 'ecrypt_preferences',
  THEME: 'ecrypt_theme',
} as const;

// File Upload Configuration
export const FILE_CONFIG = {
  MAX_FILE_SIZE: 100 * 1024 * 1024, // 100MB
  CHUNK_SIZE: 1024 * 1024, // 1MB chunks
  ALLOWED_TYPES: [
    'image/*',
    'application/pdf',
    'text/*',
    'application/json',
    'application/zip',
    'application/x-zip-compressed',
    '.doc', '.docx',
    '.xls', '.xlsx',
    '.ppt', '.pptx',
  ],
  CONCURRENT_UPLOADS: 3,
} as const;

// UI Constants
export const UI_CONFIG = {
  // Pagination
  DEFAULT_PAGE_SIZE: 20,
  PAGE_SIZE_OPTIONS: [10, 20, 50, 100],
  
  // Debounce timing
  SEARCH_DEBOUNCE: 300,
  
  // Animation durations
  ANIMATION_FAST: 150,
  ANIMATION_NORMAL: 300,
  ANIMATION_SLOW: 500,
  
  // Toast display duration
  TOAST_DURATION: 5000,
} as const;

// Security Settings
export const SECURITY_CONFIG = {
  // Token expiry buffer (refresh before actual expiry)
  TOKEN_REFRESH_BUFFER: 5 * 60 * 1000, // 5 minutes
  
  // Session timeout
  SESSION_TIMEOUT: 30 * 60 * 1000, // 30 minutes
  
  // Password requirements
  PASSWORD_MIN_LENGTH: 8,
  PASSWORD_REQUIRE_UPPERCASE: true,
  PASSWORD_REQUIRE_LOWERCASE: true,
  PASSWORD_REQUIRE_NUMBERS: true,
  PASSWORD_REQUIRE_SYMBOLS: true,
} as const;

// Status Mappings
export const STATUS_COLORS = {
  success: 'text-green-700 bg-green-50 border-green-200',
  warning: 'text-yellow-700 bg-yellow-50 border-yellow-200',
  error: 'text-red-700 bg-red-50 border-red-200',
  info: 'text-blue-700 bg-blue-50 border-blue-200',
  neutral: 'text-gray-700 bg-gray-50 border-gray-200',
} as const;

// File Size Formatting
export const FILE_SIZE_UNITS = ['B', 'KB', 'MB', 'GB', 'TB'] as const;

// Date Format Options
export const DATE_FORMATS = {
  SHORT: 'MMM dd, yyyy',
  LONG: 'MMMM dd, yyyy',
  WITH_TIME: 'MMM dd, yyyy HH:mm',
  TIME_ONLY: 'HH:mm:ss',
  ISO: 'yyyy-MM-dd',
} as const;

// Error Messages
export const ERROR_MESSAGES = {
  NETWORK_ERROR: 'Network error. Please check your connection and try again.',
  UNAUTHORIZED: 'Your session has expired. Please log in again.',
  FORBIDDEN: 'You do not have permission to perform this action.',
  NOT_FOUND: 'The requested resource was not found.',
  SERVER_ERROR: 'Internal server error. Please try again later.',
  FILE_TOO_LARGE: 'File size exceeds the maximum limit.',
  INVALID_FILE_TYPE: 'File type is not supported.',
  UPLOAD_FAILED: 'File upload failed. Please try again.',
  LOGIN_FAILED: 'Invalid username or password.',
  TOKEN_EXPIRED: 'Your session has expired. Please log in again.',
} as const;

// Success Messages
export const SUCCESS_MESSAGES = {
  LOGIN_SUCCESS: 'Login successful!',
  LOGOUT_SUCCESS: 'Logged out successfully.',
  FILE_UPLOADED: 'File uploaded successfully.',
  FILE_DELETED: 'File deleted successfully.',
  FILE_ENCRYPTED: 'File encrypted successfully.',
  FILE_DECRYPTED: 'File decrypted successfully.',
  USER_CREATED: 'User created successfully.',
  USER_UPDATED: 'User updated successfully.',
  USER_DELETED: 'User deleted successfully.',
  APP_CREATED: 'Application created successfully.',
  APP_DELETED: 'Application deleted successfully.',
  SECRET_ROTATED: 'Client secret rotated successfully.',
} as const;

// Feature Flags
export const FEATURES = {
  KMS_INTEGRATION: import.meta.env.VITE_KMS_ENABLED === 'true',
  BULK_OPERATIONS: import.meta.env.VITE_BULK_OPS_ENABLED === 'true',
  ADMIN_PANEL: import.meta.env.VITE_ADMIN_PANEL_ENABLED === 'true',
  FILE_PREVIEW: import.meta.env.VITE_FILE_PREVIEW_ENABLED === 'true',
  NOTIFICATIONS: import.meta.env.VITE_NOTIFICATIONS_ENABLED === 'true',
} as const;

// Environment
export const ENV = {
  NODE_ENV: import.meta.env.MODE,
  IS_DEVELOPMENT: import.meta.env.DEV,
  IS_PRODUCTION: import.meta.env.PROD,
  APP_VERSION: import.meta.env.VITE_APP_VERSION || '1.0.0',
} as const;
