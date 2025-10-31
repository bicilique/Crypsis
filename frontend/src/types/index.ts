// Authentication types
export interface LoginCredentials {
  username: string;
  password: string;
  rememberMe?: boolean;
}

export interface User {
  id: string;
  username: string;
  role: 'admin' | 'user';
  createdAt: string;
  lastLogin?: string;
}

export interface AuthTokens {
  access_token: string;
  refresh_token?: string;
  expires_in?: number;
}

// File management types - matches backend FileResponse
export interface FileItem {
  id: string;
  file_name: string;
  file_size: number;
  file_type: string;
  app_id?: string;
  updated_at: string;
  deleted?: boolean;
}

// Extended file metadata - matches backend FileMetadataResponse
export interface FileMetadata {
  id?: string;
  file_name?: string;
  file_size?: number;
  file_type?: string;
  version_id?: string;
  hash?: string;
  bucket?: string;
  location?: string;
  created_at?: string;
  updated_at?: string;
}

export interface FileUploadProgress {
  fileId: string;
  progress: number;
  status: 'uploading' | 'processing' | 'complete' | 'error';
  error?: string;
}

export interface FileUploadProgress {
  fileId: string;
  progress: number;
  status: 'uploading' | 'processing' | 'complete' | 'error';
  error?: string;
}

export interface FileFilters {
  search?: string;
  type?: string;
  encrypted?: boolean;
  sortBy?: 'name' | 'size' | 'uploadedAt';
  sortOrder?: 'asc' | 'desc';
}

// Admin types - matches backend AdminResponse
export interface AdminUser {
  id: string;
  username: string;
  created_at: string;
  updated_at: string;
}

// Application types - matches backend AppDetailResponse
export interface Application {
  id: string;
  app_name: string;
  client_id: string;
  client_secret?: string;
  is_active: boolean;
  uri?: string;
  redirectUri?: string;
  created_at: string;
  updated_at: string;
}

// Application list item - matches backend AppResponse  
export interface AppListItem {
  id: string;
  app_name: string;
  is_active: boolean;
}

// File log - matches backend FileLogResponse
export interface FileLogEntry {
  file_id: string;
  actor_id: string;
  actor_type: string;
  action: string;
  ip: string;
  timestamp: string;
  user_agent: string;
  metadata?: Record<string, any>;
}

// System monitoring types
export interface SystemStats {
  totalFiles: number;
  totalUsers: number;
  totalApplications: number;
  storageUsed: number;
  storageLimit: number;
  encryptedFiles: number;
  recentUploads: number;
}

export interface SecurityAlert {
  id: string;
  type: 'critical' | 'high' | 'medium' | 'low' | 'info';
  title: string;
  description: string;
  timestamp: string;
  dismissed: boolean;
  actionRequired?: boolean;
}

// API Response types
export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  message?: string;
  error?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}

// UI Component types
export interface ButtonProps {
  variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger';
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl';
  loading?: boolean;
  disabled?: boolean;
  icon?: React.ReactNode;
  children: React.ReactNode;
  onClick?: () => void;
  type?: 'button' | 'submit' | 'reset';
  className?: string;
}

export interface CardProps {
  title?: string;
  subtitle?: string;
  action?: React.ReactNode;
  children: React.ReactNode;
  className?: string;
  hover?: boolean;
  onClick?: () => void;
}

export interface StatusBadgeProps {
  status: 'success' | 'warning' | 'error' | 'info' | 'neutral';
  text: string;
  icon?: React.ReactNode;
  size?: 'sm' | 'md' | 'lg';
}

// Form types
export interface FormState {
  loading: boolean;
  error?: string;
  success?: boolean;
}

// Navigation types
export interface NavItem {
  name: string;
  href: string;
  icon: React.ComponentType<any>;
  badge?: number;
  children?: NavItem[];
}

// Theme types
export interface Theme {
  colors: {
    primary: Record<string, string>;
    secondary: Record<string, string>;
    success: string;
    warning: string;
    error: string;
  };
  spacing: Record<string, string>;
  typography: {
    fontFamily: Record<string, string>;
    fontSize: Record<string, string>;
  };
}
