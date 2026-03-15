import React, { useEffect, useState } from 'react';
import { notificationService, Notification } from '../../services/notificationService';

const NotificationListPage: React.FC = () => {
  const [notifications, setNotifications] = useState<Notification[]>([]);
  const [loading, setLoading] = useState(true);
  const [page, setPage] = useState(1);

  useEffect(() => {
    loadNotifications();
  }, [page]);

  const loadNotifications = async () => {
    try {
      const data = await notificationService.getNotifications(page, 10);
      setNotifications(data.notifications);
    } catch (error) {
      console.error('Failed to load notifications:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleMarkAsRead = async (id: string) => {
    try {
      await notificationService.markAsRead(id);
      loadNotifications();
    } catch (error) {
      console.error('Failed to mark as read:', error);
    }
  };

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'success':
        return 'bg-green-100 text-green-800';
      case 'warning':
        return 'bg-yellow-100 text-yellow-800';
      case 'error':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-blue-100 text-blue-800';
    }
  };

  if (loading) {
    return <div className="text-center py-12">Loading notifications...</div>;
  }

  return (
    <div className="space-y-6">
      <h1 className="text-3xl font-bold text-gray-900">Notifications</h1>

      <div className="space-y-4">
        {notifications?.map((notification) => (
          <div
            key={notification.id}
            className={`bg-white rounded-lg shadow p-6 ${
              !notification.is_read ? 'border-l-4 border-indigo-500' : ''
            }`}
          >
            <div className="flex justify-between items-start">
              <div className="flex-1">
                <div className="flex items-center space-x-2">
                  <h3 className="text-lg font-semibold text-gray-900">{notification.title}</h3>
                  <span className={`px-2 py-1 text-xs rounded-full ${getTypeColor(notification.type)}`}>
                    {notification.type}
                  </span>
                  {!notification.is_read && (
                    <span className="px-2 py-1 text-xs bg-indigo-100 text-indigo-800 rounded-full">
                      New
                    </span>
                  )}
                </div>
                <p className="text-gray-600 mt-2">{notification.message}</p>
                <p className="text-xs text-gray-400 mt-2">
                  {new Date(notification.created_at).toLocaleString()}
                </p>
              </div>
              {!notification.is_read && (
                <button
                  onClick={() => handleMarkAsRead(notification.id)}
                  className="ml-4 px-3 py-1 text-sm bg-indigo-100 text-indigo-700 rounded-md hover:bg-indigo-200"
                >
                  Mark as Read
                </button>
              )}
            </div>
          </div>
        ))}
      </div>

      {notifications?.length === 0 && (
        <div className="text-center py-12 text-gray-500">No notifications found</div>
      )}
    </div>
  );
};

export default NotificationListPage;
