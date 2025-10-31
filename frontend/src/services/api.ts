import axios, { AxiosInstance, AxiosResponse, AxiosError } from 'axios';
import { API_CONFIG, API_ENDPOINTS, STORAGE_KEYS, ERROR_MESSAGES } from '@/constants';
import { AuthTokens, ApiResponse } from '@/types';

/**
 * Custom error class for API errors
 */
export class ApiError extends Error {
  constructor(
    message: string,
    public status?: number,
    public code?: string
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

/**
 * Main API client class
 */
class ApiClient {
  private client: AxiosInstance;
  
  constructor() {
    this.client = axios.create({
      baseURL: API_CONFIG.BASE_URL,
      timeout: API_CONFIG.TIMEOUT,
      headers: {
        'Content-Type': 'application/json',
      },
    });
    
    this.setupInterceptors();
  }
  
  /**
   * Setup request and response interceptors
   */
  private setupInterceptors(): void {
    // Request interceptor for auth token
    this.client.interceptors.request.use(
      (config) => {
        const token = this.getAccessToken();
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      (error) => Promise.reject(error)
    );
    
    // Response interceptor for error handling and token refresh
    this.client.interceptors.response.use(
      (response: AxiosResponse) => response,
      async (error: AxiosError) => {
        const originalRequest = error.config;
        
        if (error.response?.status === 401 && originalRequest && !originalRequest._retry) {
          originalRequest._retry = true;
          
          try {
            const refreshToken = this.getRefreshToken();
            if (refreshToken) {
              await this.refreshAccessToken();
              
              // Update the authorization header and retry
              const newToken = this.getAccessToken();
              if (newToken && originalRequest.headers) {
                originalRequest.headers.Authorization = `Bearer ${newToken}`;
                return this.client(originalRequest);
              }
            }
          } catch (refreshError) {
            // Refresh failed, clear tokens and redirect to login
            this.clearTokens();
            window.location.href = '/login';
            return Promise.reject(refreshError);
          }
        }
        
        return Promise.reject(this.handleApiError(error));
      }
    );
  }
  
  /**
   * Handle API errors and convert to standardized format
   */
  private handleApiError(error: AxiosError): ApiError {
    if (!error.response) {
      // Network error
      return new ApiError(ERROR_MESSAGES.NETWORK_ERROR);
    }
    
    const { status, data } = error.response;
    let message: string = ERROR_MESSAGES.SERVER_ERROR;
    
    switch (status) {
      case 400:
        message = (data as any)?.message || 'Bad request';
        break;
      case 401:
        message = ERROR_MESSAGES.UNAUTHORIZED;
        break;
      case 403:
        message = ERROR_MESSAGES.FORBIDDEN;
        break;
      case 404:
        message = ERROR_MESSAGES.NOT_FOUND;
        break;
      case 500:
        message = ERROR_MESSAGES.SERVER_ERROR;
        break;
      default:
        message = (data as any)?.message || `HTTP ${status} Error`;
    }
    
    return new ApiError(message, status);
  }
  
  /**
   * Token management methods
   */
  private getAccessToken(): string | null {
    return localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN);
  }
  
  private getRefreshToken(): string | null {
    return localStorage.getItem(STORAGE_KEYS.REFRESH_TOKEN);
  }
  
  private setTokens(accessToken: string, refreshToken?: string): void {
    localStorage.setItem(STORAGE_KEYS.ACCESS_TOKEN, accessToken);
    if (refreshToken) {
      localStorage.setItem(STORAGE_KEYS.REFRESH_TOKEN, refreshToken);
    }
  }
  
  private clearTokens(): void {
    localStorage.removeItem(STORAGE_KEYS.ACCESS_TOKEN);
    localStorage.removeItem(STORAGE_KEYS.REFRESH_TOKEN);
    localStorage.removeItem(STORAGE_KEYS.USER_DATA);
  }
  
  /**
   * Refresh access token
   */
  private async refreshAccessToken(): Promise<void> {
    const response = await this.client.get<ApiResponse<AuthTokens>>(
      API_ENDPOINTS.REFRESH_TOKEN
    );
    
    if (response.data.success && response.data.data) {
      const { access_token, refresh_token } = response.data.data;
      this.setTokens(access_token, refresh_token);
    } else {
      throw new ApiError('Token refresh failed');
    }
  }
  
  /**
   * Generic HTTP methods
   */
  async get<T>(url: string, params?: any): Promise<T> {
    const response = await this.client.get<T>(url, { params });
    return response.data;
  }
  
  async post<T>(url: string, data?: any, config?: any): Promise<T> {
    const response = await this.client.post<T>(url, data, config);
    return response.data;
  }
  
  async put<T>(url: string, data?: any, config?: any): Promise<T> {
    const response = await this.client.put<T>(url, data, config);
    return response.data;
  }
  
  async patch<T>(url: string, data?: any): Promise<T> {
    const response = await this.client.patch<T>(url, data);
    return response.data;
  }
  
  async delete<T>(url: string): Promise<T> {
    const response = await this.client.delete<T>(url);
    return response.data;
  }
  
  /**
   * File upload with progress tracking
   */
  async uploadFile(
    url: string,
    file: File,
    onProgress?: (progress: number) => void
  ): Promise<any> {
    const formData = new FormData();
    formData.append('file', file);
    
    const response = await this.client.post(url, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
      onUploadProgress: (progressEvent) => {
        if (onProgress && progressEvent.total) {
          const progress = Math.round((progressEvent.loaded * 100) / progressEvent.total);
          onProgress(progress);
        }
      },
    });
    
    return response.data;
  }
  
  /**
   * File download
   */
  async downloadFile(url: string, filename?: string): Promise<Blob> {
    const response = await this.client.get(url, {
      responseType: 'blob',
    });
    
    // If filename is provided, trigger download
    if (filename) {
      const blob = response.data;
      const downloadUrl = URL.createObjectURL(blob);
      const link = document.createElement('a');
      link.href = downloadUrl;
      link.download = filename;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(downloadUrl);
    }
    
    return response.data;
  }
  
  /**
   * Set authentication tokens
   */
  setAuthTokens(accessToken: string, refreshToken?: string): void {
    this.setTokens(accessToken, refreshToken);
  }
  
  /**
   * Clear authentication tokens
   */
  clearAuthTokens(): void {
    this.clearTokens();
  }
  
  /**
   * Check if user is authenticated
   */
  isAuthenticated(): boolean {
    return !!this.getAccessToken();
  }
}

// Export singleton instance
export const apiClient = new ApiClient();

// Type declarations for axios
declare module 'axios' {
  interface AxiosRequestConfig {
    _retry?: boolean;
  }
}
