import api from './api';

export interface Employee {
  id: string;
  name: string;
  email: string;
  role: 'admin' | 'employee';
  department?: string;
  designation?: string;
  joining_date?: string;
  salary?: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

export interface EmployeeListResponse {
  employees: Employee[];
  total: number;
  page: number;
  limit: number;
}

export interface CreateEmployeeRequest {
  name: string;
  email: string;
  password: string;
  role?: 'admin' | 'employee';
  department?: string;
  designation?: string;
  joining_date?: string;
  salary?: number;
}

export interface UpdateEmployeeRequest {
  name?: string;
  email?: string;
  role?: 'admin' | 'employee';
  department?: string;
  designation?: string;
  joining_date?: string;
  salary?: number;
  is_active?: boolean;
}

export const employeeService = {
  list: async (page = 1, limit = 10): Promise<EmployeeListResponse> => {
    const response = await api.get<EmployeeListResponse>('/employees', {
      params: { page, limit },
    });
    return response.data;
  },

  getById: async (id: string): Promise<Employee> => {
    const response = await api.get<Employee>(`/employees/${id}`);
    return response.data;
  },

  getMe: async (): Promise<Employee> => {
    const response = await api.get<Employee>('/employees/me');
    return response.data;
  },

  create: async (data: CreateEmployeeRequest): Promise<Employee> => {
    const response = await api.post<Employee>('/employees', data);
    return response.data;
  },

  update: async (id: string, data: UpdateEmployeeRequest): Promise<Employee> => {
    const response = await api.put<Employee>(`/employees/${id}`, data);
    return response.data;
  },
};
