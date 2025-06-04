import api from '@/lib/axios';
import type { User, NewUserPayload } from '@/types';

export const healthCheck = async () => {
  const response = await api.get('/health');
  return response.data;
};

export const userApi = {
  register: async (payload: NewUserPayload): Promise<{ user: User; token: string }> => {
    const response = await api.post('/user', payload);
    return response.data as { user: User; token: string };
  },

  getCurrentUser: async (): Promise<User> => {
    const response = await api.get('/user/me');
    return response.data as User;
  },
};
