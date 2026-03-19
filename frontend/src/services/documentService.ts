import api from './api';

export interface Document {
  id: string;
  user_id: string;
  file_url: string;
  file_name: string;
  file_size: number;
  document_type: string;
  uploaded_at: string;
  created_at: string;
  updated_at: string;
}

export interface DocumentWithUser extends Document {
  user_name: string;
  user_email: string;
}

export interface DocumentListResponse {
  documents: Document[];
}

export const documentService = {
  upload: async (file: File, documentType: string): Promise<Document> => {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('document_type', documentType);
    const response = await api.post<Document>('/documents', formData, {
      headers: { 'Content-Type': 'multipart/form-data' },
    });
    return response.data;
  },

  getMyDocuments: async (page = 1, limit = 10): Promise<DocumentListResponse> => {
    const response = await api.get<DocumentListResponse>('/documents/me', {
      params: { page, limit },
    });
    return response.data;
  },

  getByUserId: async (userId: string, page = 1, limit = 10): Promise<DocumentListResponse> => {
    const response = await api.get<DocumentListResponse>(`/documents/user/${userId}`, {
      params: { page, limit },
    });
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/documents/${id}`);
  },

  download: async (id: string): Promise<void> => {
    const response = await api.get(`/documents/${id}/download`, {
      responseType: 'blob',
    });
    
    // Create blob URL and trigger download
    const blob = new Blob([response.data]);
    const url = window.URL.createObjectURL(blob);
    const link = document.createElement('a');
    link.href = url;
    
    // Try to get filename from Content-Disposition header
    const contentDisposition = response.headers['content-disposition'];
    let filename = `document-${id}`;
    if (contentDisposition) {
      const filenameMatch = contentDisposition.match(/filename="?(.+)"?/i);
      if (filenameMatch && filenameMatch[1]) {
        filename = filenameMatch[1];
      }
    }
    
    link.setAttribute('download', filename);
    document.body.appendChild(link);
    link.click();
    link.remove();
    window.URL.revokeObjectURL(url);
  },

  getAllDocuments: async (page = 1, limit = 10): Promise<DocumentListResponse & { total: number }> => {
    const response = await api.get<DocumentListResponse & { total: number; page: number; limit: number }>('/documents', {
      params: { page, limit },
    });
    return response.data;
  },
};
