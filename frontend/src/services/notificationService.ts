import api from './api';

export interface Notification {
  id: string;
  user_id: string;
  title: string;
  message: string;
  type: 'info' | 'success' | 'warning' | 'error';
  is_read: boolean;
  created_at: string;
  updated_at: string;
}

export interface NotificationListResponse {
  notifications: Notification[];
}

export const notificationService = {
  getNotifications: async (page = 1, limit = 10): Promise<NotificationListResponse> => {
    const response = await api.get<NotificationListResponse>('/notifications', {
      params: { page, limit },
    });
    return response.data;
  },

  getUnreadCount: async (): Promise<number> => {
    const response = await api.get<{ count: number }>('/notifications/unread-count');
    return response.data.count;
  },

  markAsRead: async (id: string): Promise<void> => {
    await api.put(`/notifications/${id}/read`);
  },
};
