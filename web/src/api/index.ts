import api from "@/lib/axios";
import type { User, NewUserPayload, GameListPayload } from "@/types";

export const healthCheck = async () => {
  const response = await api.get("/health");
  return response.data;
};

export const userApi = {
  register: async (
    payload: NewUserPayload
  ): Promise<{ user: User; token: string }> => {
    const response = await api.post("/user", payload);
    return response.data as { user: User; token: string };
  },

  getCurrentUser: async (): Promise<User> => {
    const response = await api.get("/user/me");
    return response.data as User;
  },
};

export const gameApi = {
  getGamesList: async (): Promise<{
    type: string;
    data: GameListPayload[];
  }> => {
    const response = await api.get("/game/list");
    return response.data as { type: string; data: GameListPayload[] };
  },
};

export const roomApi = {
  createRoom: () => api.post("/room"),
  joinRoom: (roomID: string) => api.post(`/room/${roomID}/join`),
  getRoom: (roomID: string) => api.get(`/room/${roomID}`),
};
