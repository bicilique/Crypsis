import { apiClient } from './api';
import { API_ENDPOINTS } from '@/constants';

interface ListResponse<T> {
  success: boolean;
  message: string;
  count: number;
  data: T[];
}

interface StandardResponse<T = any> {
  success: boolean;
  message: string;
  data?: T;
}

/**
 * Admin management service - matches backend API
 */
export const adminService = {
  async listAdmins(params?: { offset?: number; limit?: number; sort_by?: string; order?: string }) {
    const response = await apiClient.get<StandardResponse<any[]>>(API_ENDPOINTS.ADMIN_LIST, { params });
    if (!response.success) throw new Error(response.message);
    return response.data || [];
  },

  async addAdmin(data: { username: string; password: string }) {
    const response = await apiClient.post<StandardResponse>(API_ENDPOINTS.ADMIN_ADD, data);
    if (!response.success) throw new Error(response.message);
    return response.data;
  },

  async updateUsername(username: string) {
    const response = await apiClient.patch<StandardResponse>(API_ENDPOINTS.ADMIN_UPDATE_USERNAME, { username });
    if (!response.success) throw new Error(response.message);
  },

  async updatePassword(password: string) {
    const response = await apiClient.patch<StandardResponse>(API_ENDPOINTS.ADMIN_UPDATE_PASSWORD, { password });
    if (!response.success) throw new Error(response.message);
  },

  async deleteAdmin(id: string) {
    const response = await apiClient.delete<StandardResponse>(`${API_ENDPOINTS.ADMIN_DELETE}?id=${id}`);
    if (!response.success) throw new Error(response.message);
  },
};

/**
 * Applications management service
 */
export const applicationsService = {
  async listApplications(params?: { offset?: number; limit?: number; sort_by?: string; order?: string }) {
    const response = await apiClient.get<ListResponse<any>>(API_ENDPOINTS.APPS, { params });
    if (!response.success) throw new Error(response.message);
    return {
      count: response.count || 0,
      data: response.data || []
    };
  },

  async getApplication(id: string) {
    const response = await apiClient.get<StandardResponse>(API_ENDPOINTS.APP_DETAIL(id));
    if (!response.success) throw new Error(response.message);
    return response.data;
  },

  async createApplication(data: { name: string; uri: string; redirectUri: string }) {
    const response = await apiClient.post<StandardResponse>(API_ENDPOINTS.APPS, data);
    if (!response.success) throw new Error(response.message);
    return response.data;
  },

  async deleteApplication(id: string) {
    const response = await apiClient.delete<StandardResponse>(API_ENDPOINTS.APP_DELETE(id));
    if (!response.success) throw new Error(response.message);
  },

  async recoverApplication(id: string) {
    const response = await apiClient.post<StandardResponse>(`/api/admin/apps/${id}/recover`);
    if (!response.success) throw new Error(response.message);
    return response.data;
  },

  async rotateSecret(id: string) {
    const response = await apiClient.put<StandardResponse>(API_ENDPOINTS.APP_ROTATE_SECRET(id));
    if (!response.success) throw new Error(response.message);
    return response.data;
  },
};

/**
 * Files management service - matches backend file operations
 */
export const filesService = {
  // Client file operations
  async listFiles(params?: { offset?: number; limit?: number; sort_by?: string; order?: string }) {
    const response = await apiClient.get<StandardResponse<{ count: number; files: any[] }>>(
      API_ENDPOINTS.FILE_LIST, 
      { params }
    );
    if (!response.success) throw new Error(response.message);
    return response.data?.files || [];
  },

  async uploadFile(file: File, onProgress?: (progress: number) => void) {
    const formData = new FormData();
    formData.append('file', file);

    const response = await apiClient.uploadFile(
      API_ENDPOINTS.FILE_UPLOAD,
      file,
      onProgress
    );

    if (!response.success) throw new Error(response.message);
    return response.data?.file_id;
  },

  async downloadFile(id: string, filename: string) {
    try {
      const blob = await apiClient.downloadFile(API_ENDPOINTS.FILE_DOWNLOAD(id), filename);
      return blob;
    } catch (error) {
      throw new Error(`Failed to download file: ${error}`);
    }
  },

  async updateFile(id: string, file: File) {
    const formData = new FormData();
    formData.append('file', file);

    const response = await apiClient.put<StandardResponse>(
      API_ENDPOINTS.FILE_UPDATE(id),
      formData,
      { headers: { 'Content-Type': 'multipart/form-data' } }
    );

    if (!response.success) throw new Error(response.message);
  },

  async deleteFile(id: string) {
    const response = await apiClient.delete<StandardResponse>(API_ENDPOINTS.FILE_DELETE(id));
    if (!response.success) throw new Error(response.message);
  },

  async getFileMetadata(id: string) {
    const response = await apiClient.get<StandardResponse>(API_ENDPOINTS.FILE_METADATA(id));
    if (!response.success) throw new Error(response.message);
    return response.data;
  },

  async encryptFile(file: File) {
    const formData = new FormData();
    formData.append('file', file);

    const response = await apiClient.post(API_ENDPOINTS.FILE_ENCRYPT, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
      responseType: 'blob',
    });

    return response as any as Blob;
  },

  async decryptFile(file: File, fileId: string) {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('id', fileId);

    const response = await apiClient.post(API_ENDPOINTS.FILE_DECRYPT, formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
      responseType: 'blob',
    });

    return response as any as Blob;
  },

  // Admin file operations
  async listAdminFiles(params?: { offset?: number; limit?: number; sort_by?: string; order?: string }) {
    const response = await apiClient.get<ListResponse<any>>(API_ENDPOINTS.ADMIN_FILES, { params });
    if (!response.success) throw new Error(response.message);
    return {
      count: response.count || 0,
      data: response.data || []
    };
  },

  async listFilesByApp(appId: string, params?: { offset?: number; limit?: number; sort_by?: string; order?: string }) {
    const response = await apiClient.get<ListResponse<any>>(`/api/admin/apps/${appId}/files`, { params });
    if (!response.success) throw new Error(response.message);
    return {
      count: response.count || 0,
      data: response.data || []
    };
  },
};

/**
 * Logs service
 */
export const logsService = {
  async listLogs(params?: { offset?: number; limit?: number; sort_by?: string; order?: string }) {
    const response = await apiClient.get<ListResponse<any>>(API_ENDPOINTS.ADMIN_LOGS, { params });
    if (!response.success) throw new Error(response.message);
    return {
      count: response.count || 0,
      data: response.data || []
    };
  },
};

/**
 * Security service
 */
export const securityService = {
  async rekeyFiles(keyUID: string) {
    const response = await apiClient.post<StandardResponse>(API_ENDPOINTS.ADMIN_REKEY, { keyUID });
    if (!response.success) throw new Error(response.message);
    return response.data;
  },
};
