import React, { useEffect, useState } from 'react';
import { dashboardService, AdminDashboard } from '../../services/dashboardService';
import StatCard from '../../components/dashboard/StatCard';
import AttendanceChart from '../../components/dashboard/AttendanceChart';
import PayrollChart from '../../components/dashboard/PayrollChart';

const AdminDashboardPage: React.FC = () => {
  const [dashboard, setDashboard] = useState<AdminDashboard | null>(null);
  const [loading, setLoading] = useState(true);
  const [range, setRange] = useState<'day' | 'month' | 'year'>('month');

  useEffect(() => {
    loadDashboard();
  }, [range]);

  const loadDashboard = async () => {
    try {
      const data = await dashboardService.getAdminDashboard(range);
      setDashboard(data);
    } catch (error) {
      console.error('Failed to load dashboard:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return <div className="text-center py-12">Loading dashboard...</div>;
  }

  if (!dashboard) {
    return <div className="text-center py-12 text-red-600">Failed to load dashboard</div>;
  }

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-900">Admin Dashboard</h1>
        <p className="text-gray-600 mt-2">Overview of your HR management system</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <StatCard
          title="Total Employees"
          value={dashboard.total_employees}
          icon="👥"
          color="blue"
        />
        <StatCard
          title="Present Today"
          value={dashboard.present_today}
          icon="✅"
          color="green"
        />
        <StatCard
          title="Absent Today"
          value={dashboard.absent_today}
          icon="❌"
          color="red"
        />
        <StatCard
          title="Pending Leaves"
          value={dashboard.pending_leaves}
          icon="📅"
          color="yellow"
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white rounded-lg shadow p-6">
          <div className="flex items-center justify-between mb-4">
            <h2 className="text-xl font-semibold text-gray-900">
              Attendance Trend (
              {range === 'day' ? 'Last 7 Days' : range === 'year' ? 'Last 1 Year' : 'Last 30 Days'})
            </h2>
            <select
              value={range}
              onChange={(e) => setRange(e.target.value as 'day' | 'month' | 'year')}
              className="px-3 py-1 border rounded-md text-sm"
            >
              <option value="day">Day</option>
              <option value="month">Month</option>
              <option value="year">Year</option>
            </select>
          </div>
          <AttendanceChart
            data={(dashboard.attendance_trend || []).slice(
              -(range === 'day' ? 7 : range === 'year' ? 365 : 30)
            )}
          />
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">Monthly Payroll Trend</h2>
          <PayrollChart data={dashboard.payroll_trend || []} />
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">Payroll Summary</h2>
          <div className="space-y-2">
            {dashboard.payroll_summary ? (
              <>
                <p className="text-sm text-gray-600">
                  Month: {dashboard.payroll_summary.month}/{dashboard.payroll_summary.year}
                </p>
                <p className="text-2xl font-bold text-gray-900">
                  ${dashboard.payroll_summary.total.toLocaleString()}
                </p>
              </>
            ) : (
              <p className="text-sm text-gray-500">No payroll data available</p>
            )}
          </div>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-xl font-semibold text-gray-900 mb-4">Recent Activities</h2>
          <div className="space-y-3 max-h-64 overflow-y-auto">
            {dashboard.recent_activities && dashboard.recent_activities.length > 0 ? (
              dashboard.recent_activities.map((activity, index) => (
                <div key={index} className="border-b pb-2 last:border-0">
                  <p className="text-sm font-medium text-gray-900">{activity.action}</p>
                  <p className="text-xs text-gray-500">{activity.entity_type}</p>
                  <p className="text-xs text-gray-400">
                    {new Date(activity.created_at).toLocaleString()}
                  </p>
                </div>
              ))
            ) : (
              <p className="text-sm text-gray-500">No recent activities</p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default AdminDashboardPage;
