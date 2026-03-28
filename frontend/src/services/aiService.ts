import api from './api';

export interface HRQueryRequest {
  query: string;
}

export interface AIResponse {
  type: string;
  message: string;
  data: Array<Record<string, any>>;
}

export const aiService = {
  query: async (query: string): Promise<AIResponse> => {
    const response = await api.post<AIResponse>('/ai/hr-assistant', { query });
    return response.data;
  },
};
