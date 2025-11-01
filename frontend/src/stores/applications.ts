import { create } from 'zustand';
import { applicationsService } from '@/services';

export interface Application {
  id: string;
  app_name: string;
  client_id: string;
  client_secret?: string;
  is_active: boolean;
  uri?: string;
  redirect_uri?: string;
  created_at: string;
  updated_at: string;
}

interface ApplicationState {
  applications: Application[];
  isLoading: boolean;
  error: string | null;
  
  fetchApplications: (params?: { offset?: number; limit?: number; sort_by?: string; order?: string }) => Promise<void>;
  getApplication: (id: string) => Promise<Application>;
  createApplication: (data: { name: string; uri: string; redirectUri: string }) => Promise<void>;
  deleteApplication: (id: string) => Promise<void>;
  recoverApplication: (id: string) => Promise<void>;
  rotateSecret: (id: string) => Promise<void>;
  clearError: () => void;
}

export const useApplicationStore = create<ApplicationState>((set, get) => ({
  applications: [],
  isLoading: false,
  error: null,

  fetchApplications: async (params = {}) => {
    set({ isLoading: true, error: null });
    try {
      const response = await applicationsService.listApplications(params);
      set({ applications: response.data, isLoading: false });
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : 'Failed to fetch applications',
        isLoading: false,
      });
    }
  },

  getApplication: async (id: string) => {
    set({ isLoading: true, error: null });
    try {
      const app = await applicationsService.getApplication(id);
      set({ isLoading: false });
      return app;
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : 'Failed to fetch application',
        isLoading: false,
      });
      throw error;
    }
  },

  createApplication: async (data) => {
    set({ isLoading: true, error: null });
    try {
      await applicationsService.createApplication(data);
      await get().fetchApplications();
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : 'Failed to create application',
        isLoading: false,
      });
      throw error;
    }
  },

  deleteApplication: async (id: string) => {
    set({ isLoading: true, error: null });
    try {
      await applicationsService.deleteApplication(id);
      await get().fetchApplications();
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : 'Failed to delete application',
        isLoading: false,
      });
      throw error;
    }
  },

  recoverApplication: async (id: string) => {
    set({ isLoading: true, error: null });
    try {
      await applicationsService.recoverApplication(id);
      await get().fetchApplications();
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : 'Failed to recover application',
        isLoading: false,
      });
      throw error;
    }
  },

  rotateSecret: async (id: string) => {
    set({ isLoading: true, error: null });
    try {
      await applicationsService.rotateSecret(id);
      await get().fetchApplications();
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : 'Failed to rotate secret',
        isLoading: false,
      });
      throw error;
    }
  },

  clearError: () => set({ error: null }),
}));
