import { create } from 'zustand';
import { logsService } from '@/services';

export interface FileLog {
  file_id: string;
  actor_id: string;
  actor_type: string;
  action: string;
  ip: string;
  timestamp: string;
  user_agent: string;
  metadata?: Record<string, any>;
}

interface LogsState {
  logs: FileLog[];
  isLoading: boolean;
  error: string | null;
  
  fetchLogs: (params?: { offset?: number; limit?: number; sort_by?: string; order?: string }) => Promise<void>;
  clearError: () => void;
}

export const useLogsStore = create<LogsState>((set) => ({
  logs: [],
  isLoading: false,
  error: null,

  fetchLogs: async (params = {}) => {
    set({ isLoading: true, error: null });
    try {
      const response = await logsService.listLogs(params);
      set({ logs: response.data, isLoading: false });
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : 'Failed to fetch logs',
        isLoading: false,
      });
    }
  },

  clearError: () => set({ error: null }),
}));
