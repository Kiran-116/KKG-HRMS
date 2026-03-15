import api from './api';

export interface Leave {
  id: string;
  user_id: string;
  start_date: string;
  end_date: string;
  reason: string;
  status: 'pending' | 'approved' | 'rejected';
  approved_by?: string;
  created_at: string;
  updated_at: string;
}

export interface LeaveListResponse {
  leaves: Leave[];
}

export interface ApplyLeaveRequest {
  start_date: string;
  end_date: string;
  reason: string;
}

export const leaveService = {
  apply: async (data: ApplyLeaveRequest): Promise<Leave> => {
    const response = await api.post<Leave>('/leaves/apply', data);
    return response.data;
  },

  getMyLeaves: async (page = 1, limit = 10): Promise<LeaveListResponse> => {
    const response = await api.get<LeaveListResponse>('/leaves/me', {
      params: { page, limit },
    });
    return response.data;
  },

  getAll: async (page = 1, limit = 10, status?: string): Promise<LeaveListResponse> => {
    const response = await api.get<LeaveListResponse>('/leaves', {
      params: { page, limit, status },
    });
    return response.data;
  },

  approve: async (id: string): Promise<Leave> => {
    const response = await api.put<Leave>(`/leaves/${id}/approve`);
    return response.data;
  },

  reject: async (id: string): Promise<Leave> => {
    const response = await api.put<Leave>(`/leaves/${id}/reject`);
    return response.data;
  },
};
