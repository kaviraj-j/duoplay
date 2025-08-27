import api from "@/lib/axios";
import type { User, NewUserPayload } from "@/types";

export const userApi = {
  register: async (
    payload: NewUserPayload
  ): Promise<{ data: User; token: string }> => {
    const response = await api.post("/user", payload);
    return response.data as { data: User; token: string };
  },

  getCurrentUser: async (): Promise<User> => {
    const response = await api.get<{ data: User }>("/user/me");
    return response.data.data as User;
  },
};
