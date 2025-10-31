import { create } from 'zustand';
import { filesService } from '@/services';
import { FileItem, FileFilters, FileUploadProgress } from '@/types';
import { FILE_CONFIG } from '@/constants';

interface FileState {
  // State
  files: FileItem[];
  selectedFiles: string[];
  uploadProgress: Record<string, FileUploadProgress>;
  filters: FileFilters;
  isLoading: boolean;
  error: string | null;
  
  // Pagination
  currentPage: number;
  totalPages: number;
  totalFiles: number;
  
  // Actions
  fetchFiles: (page?: number) => Promise<void>;
  uploadFile: (file: File) => Promise<void>;
  uploadMultipleFiles: (files: File[]) => Promise<void>;
  deleteFile: (fileId: string) => Promise<void>;
  bulkDeleteFiles: (fileIds: string[]) => Promise<void>;
  downloadFile: (fileId: string, filename: string) => Promise<void>;
  updateFilters: (filters: Partial<FileFilters>) => void;
  selectFile: (fileId: string) => void;
  selectAllFiles: (selected: boolean) => void;
  clearSelection: () => void;
  clearError: () => void;
  resetUploadProgress: (fileId: string) => void;
}

export const useFileStore = create<FileState>((set, get) => ({
  // Initial state
  files: [],
  selectedFiles: [],
  uploadProgress: {},
  filters: {
    sortBy: 'uploadedAt',
    sortOrder: 'desc',
  },
  isLoading: false,
  error: null,
  currentPage: 1,
  totalPages: 1,
  totalFiles: 0,
  
  // Actions
  fetchFiles: async (page = 1) => {
    set({ isLoading: true, error: null });
    
    try {
      const { filters } = get();
      const offset = (page - 1) * 20;
      const files = await filesService.listFiles({ 
        offset, 
        limit: 20,
        sort_by: filters.sortBy,
        order: filters.sortOrder
      });
      
      set({
        files: files,
        currentPage: page,
        totalFiles: files.length,
        isLoading: false,
      });
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : 'Failed to fetch files',
        isLoading: false,
      });
    }
  },
  
  uploadFile: async (file: File) => {
    // Validate file
    if (file.size > FILE_CONFIG.MAX_FILE_SIZE) {
      throw new Error('File size exceeds maximum limit');
    }
    
    const fileId = `upload_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
    
    // Initialize upload progress
    set((state) => ({
      uploadProgress: {
        ...state.uploadProgress,
        [fileId]: {
          fileId,
          progress: 0,
          status: 'uploading',
        },
      },
    }));
    
    try {
      const uploadedFile = await filesService.uploadFile(file, (progress) => {
        set((state) => ({
          uploadProgress: {
            ...state.uploadProgress,
            [fileId]: {
              ...state.uploadProgress[fileId],
              progress,
            },
          },
        }));
      });
      
      // Update progress to complete
      set((state) => ({
        uploadProgress: {
          ...state.uploadProgress,
          [fileId]: {
            ...state.uploadProgress[fileId],
            progress: 100,
            status: 'complete',
          },
        },
        files: [uploadedFile, ...state.files],
        totalFiles: state.totalFiles + 1,
      }));
      
      // Remove progress after delay
      setTimeout(() => {
        set((state) => {
          const newProgress = { ...state.uploadProgress };
          delete newProgress[fileId];
          return { uploadProgress: newProgress };
        });
      }, 2000);
    } catch (error) {
      set((state) => ({
        uploadProgress: {
          ...state.uploadProgress,
          [fileId]: {
            ...state.uploadProgress[fileId],
            status: 'error',
            error: error instanceof Error ? error.message : 'Upload failed',
          },
        },
      }));
      throw error;
    }
  },
  
  uploadMultipleFiles: async (files: File[]) => {
    const uploads = files.map(file => get().uploadFile(file));
    await Promise.allSettled(uploads);
  },
  
  deleteFile: async (fileId: string) => {
    try {
      await filesService.deleteFile(fileId);
      
      set((state) => ({
        files: state.files.filter(file => file.id !== fileId),
        selectedFiles: state.selectedFiles.filter(id => id !== fileId),
        totalFiles: state.totalFiles - 1,
      }));
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : 'Failed to delete file',
      });
      throw error;
    }
  },
  
  bulkDeleteFiles: async (fileIds: string[]) => {
    try {
      // Delete each file individually since backend doesn't have bulk delete
      await Promise.all(fileIds.map(id => filesService.deleteFile(id)));
      
      set((state) => ({
        files: state.files.filter(file => !fileIds.includes(file.id)),
        selectedFiles: [],
        totalFiles: state.totalFiles - fileIds.length,
      }));
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : 'Failed to delete files',
      });
      throw error;
    }
  },
  
  downloadFile: async (fileId: string, filename: string) => {
    try {
      await filesService.downloadFile(fileId, filename);
    } catch (error) {
      set({
        error: error instanceof Error ? error.message : 'Failed to download file',
      });
      throw error;
    }
  },
  
  updateFilters: (newFilters: Partial<FileFilters>) => {
    set((state) => ({
      filters: { ...state.filters, ...newFilters },
      currentPage: 1, // Reset to first page when filters change
    }));
    
    // Fetch files with new filters
    get().fetchFiles(1);
  },
  
  selectFile: (fileId: string) => {
    set((state) => {
      const isSelected = state.selectedFiles.includes(fileId);
      return {
        selectedFiles: isSelected
          ? state.selectedFiles.filter(id => id !== fileId)
          : [...state.selectedFiles, fileId],
      };
    });
  },
  
  selectAllFiles: (selected: boolean) => {
    set((state) => ({
      selectedFiles: selected ? state.files.map(file => file.id) : [],
    }));
  },
  
  clearSelection: () => {
    set({ selectedFiles: [] });
  },
  
  clearError: () => {
    set({ error: null });
  },
  
  resetUploadProgress: (fileId: string) => {
    set((state) => {
      const newProgress = { ...state.uploadProgress };
      delete newProgress[fileId];
      return { uploadProgress: newProgress };
    });
  },
}));
