import { create } from 'zustand';
import { adminService } from '@/services';

export interface Admin {
  id: string;
  username: string;
  created_at: string;
  updated_at: string;
  is_active?: boolean;
}

interface AdminState {
  admins: Admin[];
  isLoading: boolean;
  error: string | null;
  
  fetchAdmins: (params?: { offset?: number; limit?: number; sort_by?: string; order?: string }) => Promise<void>;
  addAdmin: (username: string, password: string) => Promise<void>;
  updateUsername: (username: string) => Promise<void>;
  updatePassword: (password: string) => Promise<void>;
  deleteAdmin: (id: string) => Promise<void>;
  clearError: () => void;
}

export const useAdminStore = create<AdminState>((set, get) => ({
  admins: [],
  isLoading: false,
  error: null,

  fetchAdmins: async (params = {}) => {
    set({ isLoading: true, error: null });
    try {
      const data = await adminService.listAdmins(params);
      set({ admins: data, isLoading: false });
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : 'Failed to fetch admins',
        isLoading: false,
      });
    }
  },

  addAdmin: async (username: string, password: string) => {
    set({ isLoading: true, error: null });
    try {
      await adminService.addAdmin({ username, password });
      await get().fetchAdmins();
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : 'Failed to add admin',
        isLoading: false,
      });
      throw error;
    }
  },

  updateUsername: async (username: string) => {
    set({ isLoading: true, error: null });
    try {
      await adminService.updateUsername(username);
      set({ isLoading: false });
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : 'Failed to update username',
        isLoading: false,
      });
      throw error;
    }
  },

  updatePassword: async (password: string) => {
    set({ isLoading: true, error: null });
    try {
      await adminService.updatePassword(password);
      set({ isLoading: false });
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : 'Failed to update password',
        isLoading: false,
      });
      throw error;
    }
  },

  deleteAdmin: async (id: string) => {
    set({ isLoading: true, error: null });
    try {
      await adminService.deleteAdmin(id);
      await get().fetchAdmins();
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : 'Failed to delete admin',
        isLoading: false,
      });
      throw error;
    }
  },

  clearError: () => set({ error: null }),
}));
