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
    const response = await api.get<DocumentListResponse>(`/documents/${userId}`, {
      params: { page, limit },
    });
    return response.data;
  },

  delete: async (id: string): Promise<void> => {
    await api.delete(`/documents/${id}`);
  },
};
