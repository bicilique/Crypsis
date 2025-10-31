import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { authService } from '@/services';
import { User, LoginCredentials } from '@/types';
import { STORAGE_KEYS } from '@/constants';

interface AuthState {
  // State
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
  
  // Actions
  login: (credentials: LoginCredentials) => Promise<void>;
  logout: () => Promise<void>;
  refreshToken: () => Promise<void>;
  clearError: () => void;
  setUser: (user: User | null) => void;
}

export const useAuthStore = create<AuthState>()(
  persist(
    (set) => ({
      // Initial state
      user: null,
      isAuthenticated: false,
      isLoading: false,
      error: null,
      
      // Actions
      login: async (credentials: LoginCredentials) => {
        set({ isLoading: true, error: null });
        
        try {
          const { user } = await authService.login(credentials);
          
          set({
            user,
            isAuthenticated: true,
            isLoading: false,
            error: null,
          });
        } catch (error) {
          set({
            error: error instanceof Error ? error.message : 'Login failed',
            isLoading: false,
            isAuthenticated: false,
            user: null,
          });
          throw error;
        }
      },
      
      logout: async () => {
        set({ isLoading: true });
        
        try {
          await authService.logout();
        } catch (error) {
          console.warn('Logout API failed:', error);
        } finally {
          set({
            user: null,
            isAuthenticated: false,
            isLoading: false,
            error: null,
          });
        }
      },
      
      refreshToken: async () => {
        try {
          await authService.refreshToken();
        } catch (error) {
          // If refresh fails, logout user
          set({
            user: null,
            isAuthenticated: false,
            error: 'Session expired',
          });
          throw error;
        }
      },
      
      clearError: () => {
        set({ error: null });
      },
      
      setUser: (user: User | null) => {
        set({
          user,
          isAuthenticated: !!user,
        });
      },
    }),
    {
      name: STORAGE_KEYS.USER_DATA,
      partialize: (state) => ({
        user: state.user,
        isAuthenticated: state.isAuthenticated,
      }),
    }
  )
);

// Helper to check if user is authenticated on app start
export const initializeAuth = () => {
  const isAuthenticated = authService.isAuthenticated();
  
  if (!isAuthenticated) {
    useAuthStore.getState().setUser(null);
  }
};
