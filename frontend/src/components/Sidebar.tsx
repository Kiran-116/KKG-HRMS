import React from 'react';
import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

const Sidebar: React.FC = () => {
  const location = useLocation();
  const { isAdmin } = useAuth();

  const isActive = (path: string) => location.pathname === path;

  const menuItems = [
    { path: '/dashboard', label: 'Dashboard', icon: '📊' },
    ...(isAdmin
      ? [
          { path: '/employees', label: 'Employees', icon: '👥' },
          { path: '/attendance', label: 'All Attendance', icon: '⏰' },
          { path: '/attendance/checkin', label: 'Check In/Out', icon: '✅' },
          { path: '/attendance/history', label: 'My Attendance', icon: '📊' },
          { path: '/leaves', label: 'Leaves', icon: '📅' },
          { path: '/payroll', label: 'Payroll', icon: '💰' },
          { path: '/audit-logs', label: 'Audit Logs', icon: '📝' },
        ]
      : [
          { path: '/attendance/checkin', label: 'Check In/Out', icon: '⏰' },
          { path: '/attendance/history', label: 'My Attendance', icon: '📊' },
          { path: '/leaves/me', label: 'My Leaves', icon: '📅' },
          { path: '/salary/me', label: 'My Salary', icon: '💰' },
        ]),
    { path: '/documents', label: 'Documents', icon: '📄' },
    { path: '/notifications', label: 'Notifications', icon: '🔔' },
    { path: '/ai-assistant', label: 'AI Assistant', icon: '🤖' },
  ];

  return (
    <aside className="w-64 bg-white shadow-lg min-h-screen fixed left-0 top-16 z-10">
      <nav className="p-4">
        <ul className="space-y-2">
          {menuItems.map((item) => (
            <li key={item.path}>
              <Link
                to={item.path}
                className={`flex items-center space-x-3 px-4 py-3 rounded-lg transition-colors ${
                  isActive(item.path)
                    ? 'bg-indigo-100 text-indigo-700 font-semibold'
                    : 'text-gray-700 hover:bg-gray-100'
                }`}
              >
                <span className="text-xl">{item.icon}</span>
                <span>{item.label}</span>
              </Link>
            </li>
          ))}
        </ul>
      </nav>
    </aside>
  );
};

export default Sidebar;
