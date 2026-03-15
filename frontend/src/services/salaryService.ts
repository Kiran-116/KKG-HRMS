import api from './api';

export interface Salary {
  id: string;
  user_id: string;
  base_salary: number;
  bonus: number;
  deductions: number;
  net_salary: number;
  month: number;
  year: number;
  created_at: string;
  updated_at: string;
}

export interface SalaryListResponse {
  salaries: Salary[];
}

export interface CreateSalaryRequest {
  user_id: string;
  base_salary: number;
  bonus?: number;
  deductions?: number;
  month: number;
  year: number;
}

export const salaryService = {
  getMySalary: async (page = 1, limit = 10): Promise<SalaryListResponse> => {
    const response = await api.get<SalaryListResponse>('/salary/me', {
      params: { page, limit },
    });
    return response.data;
  },

  getByUserId: async (userId: string, page = 1, limit = 10): Promise<SalaryListResponse> => {
    const response = await api.get<SalaryListResponse>(`/salary/${userId}`, {
      params: { page, limit },
    });
    return response.data;
  },

  create: async (data: CreateSalaryRequest): Promise<Salary> => {
    const response = await api.post<Salary>('/salary', data);
    return response.data;
  },
};
