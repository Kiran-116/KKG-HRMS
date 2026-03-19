import api from './api';

export interface AdminDashboard {
  total_employees: number;
  present_today: number;
  absent_today: number;
  pending_leaves: number;
  recent_activities?: Array<{
    action: string;
    entity_type: string;
    created_at: string;
  }> | null;
  payroll_summary?: {
    month: number;
    year: number;
    total: number;
  } | null;
  attendance_trend?: Array<{
    date: string;
    present: number;
    absent: number;
  }>;
  payroll_trend?: Array<{
    month: number;
    year: number;
    total: number;
  }>;
}

export interface EmployeeDashboard {
  attendance_summary?: {
    present_days: number;
    absent_days: number;
    month: number;
    year: number;
  } | null;
  leave_balance: number;
  salary_summary?: {
    id: string;
    user_id: string;
    base_salary: number;
    bonus: number;
    deductions: number;
    net_salary: number;
    month: number;
    year: number;
  } | null;
  upcoming_holidays?: any[];
  unread_notifications: number;
  attendance_trend?: Array<{
    date: string;
    present: number;
  }>;
}

export const dashboardService = {
  getAdminDashboard: async (range: 'day' | 'month' | 'year' = 'month'): Promise<AdminDashboard> => {
    const response = await api.get<AdminDashboard>('/dashboard/admin', {
      params: { range },
    });
    return response.data;
  },

  getEmployeeDashboard: async (): Promise<EmployeeDashboard> => {
    const response = await api.get<EmployeeDashboard>('/dashboard/employee');
    return response.data;
  },
};
