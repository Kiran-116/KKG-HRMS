import api from './api';

export interface HRQueryRequest {
  query: string;
}

export interface HRQueryResponse {
  answer: string;
}

export const aiService = {
  query: async (query: string): Promise<string> => {
    const response = await api.post<HRQueryResponse>('/ai/hr-assistant', { query });
    return response.data.answer;
  },
};
