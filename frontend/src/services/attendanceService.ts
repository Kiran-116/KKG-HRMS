import api from './api';

export interface Attendance {
  id: string;
  user_id: string;
  date: string;
  check_in?: string;
  check_out?: string;
  status: 'present' | 'absent' | 'late' | 'half_day';
  created_at: string;
  updated_at: string;
}

export interface AttendanceListResponse {
  attendances: Attendance[];
}

export const attendanceService = {
  checkIn: async (date?: string): Promise<Attendance> => {
    const response = await api.post<Attendance>('/attendance/checkin', { date });
    return response.data;
  },

  checkOut: async (date?: string): Promise<Attendance> => {
    const response = await api.post<Attendance>('/attendance/checkout', { date });
    return response.data;
  },

  getMyAttendance: async (page = 1, limit = 10): Promise<AttendanceListResponse> => {
    const response = await api.get<AttendanceListResponse>('/attendance/me', {
      params: { page, limit },
    });
    return response.data;
  },

  getAll: async (page = 1, limit = 10, userId?: string, date?: string): Promise<AttendanceListResponse> => {
    const response = await api.get<AttendanceListResponse>('/attendance', {
      params: { page, limit, user_id: userId, date },
    });
    return response.data;
  },
};
