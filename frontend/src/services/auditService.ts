import api from './api';

export interface AuditLog {
  id: string;
  user_id?: string;
  user_name?: string;
  action: string;
  entity_type: string;
  entity_id?: string;
  description?: string;
  metadata?: Record<string, any>;
  ip_address?: string;
  user_agent?: string;
  created_at: string;
}

export interface AuditLogListResponse {
  audit_logs: AuditLog[];
}

export const auditService = {
  getAll: async (page = 1, limit = 10, userId?: string, action?: string): Promise<AuditLogListResponse> => {
    const response = await api.get<AuditLogListResponse>('/audit-logs', {
      params: { page, limit, user_id: userId, action },
    });
    return response.data;
  },
};
