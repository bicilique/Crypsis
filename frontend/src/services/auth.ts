import { apiClient } from './api';
import { API_ENDPOINTS } from '@/constants';
import { LoginCredentials, User, AuthTokens, ApiResponse } from '@/types';

/**
 * Authentication service
 */
export const authService = {
  /**
   * Login user with credentials
   */
  async login(credentials: LoginCredentials): Promise<{ user: User; tokens: AuthTokens }> {
    const response = await apiClient.post<ApiResponse<{ access_token: string }>>(
      API_ENDPOINTS.LOGIN,
      credentials
    );
    
    if (!response.success || !response.data) {
      throw new Error(response.message || 'Login failed');
    }
    
    // Backend returns { access_token: string }
    const accessToken = response.data.access_token;
    
    // Store token in API client (no refresh token in current backend implementation)
    apiClient.setAuthTokens(accessToken);
    
    // Create user object from credentials
    const user: User = {
      id: credentials.username, // Use username as ID for now
      username: credentials.username,
      role: 'admin',
      createdAt: new Date().toISOString(),
    };
    
    return { 
      user, 
      tokens: { 
        access_token: accessToken, 
        refresh_token: '' // Not used by current backend
      } 
    };
  },
  
  /**
   * Logout user
   */
  async logout(): Promise<void> {
    try {
      await apiClient.get<ApiResponse>(API_ENDPOINTS.LOGOUT);
    } catch (error) {
      // Even if logout API fails, clear local tokens
      console.warn('Logout API failed:', error);
    } finally {
      apiClient.clearAuthTokens();
    }
  },
  
  /**
   * Refresh access token
   */
  async refreshToken(): Promise<AuthTokens> {
    const response = await apiClient.get<ApiResponse<{ access_token: string }>>(
      API_ENDPOINTS.REFRESH_TOKEN
    );
    
    if (!response.success || !response.data) {
      throw new Error('Token refresh failed');
    }
    
    const accessToken = response.data.access_token;
    apiClient.setAuthTokens(accessToken);
    
    return {
      access_token: accessToken,
      refresh_token: ''
    };
  },
  
  /**
   * Check if user is currently authenticated
   */
  isAuthenticated(): boolean {
    return apiClient.isAuthenticated();
  },
  
  /**
   * Clear authentication tokens
   */
  clearAuth(): void {
    apiClient.clearAuthTokens();
  },
};
