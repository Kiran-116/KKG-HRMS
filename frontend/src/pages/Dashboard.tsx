import React from 'react';
import { useAuth } from '../contexts/AuthContext';
import AdminDashboard from './dashboard/AdminDashboard';
import EmployeeDashboard from './dashboard/EmployeeDashboard';

const Dashboard: React.FC = () => {
  const { isAdmin } = useAuth();
  return isAdmin ? <AdminDashboard /> : <EmployeeDashboard />;
};

export default Dashboard;
