import api from './api';

export interface RegisterRequest {
  name: string;
  email: string;
  password: string;
  role?: 'admin' | 'employee';
  department?: string;
  designation?: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface User {
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

export interface LoginResponse {
  user: User;
  access_token: string;
  refresh_token: string;
  token_type: string;
  expires_in: number;
  must_change_password?: boolean;
}

export interface MagicLoginRequest {
  token: string;
}

export interface SetPasswordRequest {
  new_password: string;
  confirm_password: string;
}

export const authService = {
  register: async (data: RegisterRequest): Promise<LoginResponse> => {
    const response = await api.post<LoginResponse>('/auth/register', data);
    return response.data;
  },

  login: async (data: LoginRequest): Promise<LoginResponse> => {
    const response = await api.post<LoginResponse>('/auth/login', data);
    return response.data;
  },

  getMe: async (): Promise<User> => {
    const response = await api.get<User>('/auth/me');
    return response.data;
  },

  magicLogin: async (token: string): Promise<LoginResponse> => {
    const response = await api.post<LoginResponse>('/auth/magic-login', { token });
    return response.data;
  },

  setPassword: async (data: SetPasswordRequest): Promise<void> => {
    await api.post('/auth/set-password', data);
  },

  logout: (): void => {
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    window.location.href = '/login';
  },
};
