import api from "@/lib/axios";
import type { User, NewUserPayload, GameListPayload, Room } from "@/types";

export const healthCheck = async () => {
  const response = await api.get("/health");
  return response.data;
};

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
  createRoom: async (): Promise<{
    data: Room;
  }> => {
    const response = await api.post<{ data: Room }>("/room");
    return response.data;
  },
  joinRoom: (roomID: string) => api.post(`/room/${roomID}/join`),
  getRoom: (roomID: string) => api.get(`/room/${roomID}`),

  joinQueue: () => api.post("/room/queue"),
  leaveQueue: () => api.post("/room/leaveQueue"),
  leaveRoom: (roomID: string) => api.post(`/room/leave/${roomID}`),
};
