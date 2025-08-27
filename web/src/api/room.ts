import { api, wsManager } from "@/lib";
import type { Room } from "@/types";

export const roomApi = {
  getRoom: async (roomID: string): Promise<Room> => {
    const response = await api.get<{ data: Room }>(`/room/${roomID}`);
    return response.data.data;
  },

  createRoom: async (): Promise<{ roomId: string; ws: WebSocket }> => {    
    return wsManager.createRoomConnection();
  },

  joinRoom: (roomID: string) => api.post(`/room/${roomID}/join`),

  joinQueue: () => api.post("/room/queue"),
  leaveQueue: () => api.post("/room/leaveQueue"),
  leaveRoom: (roomID: string) => {
    // Disconnect WebSocket when leaving room
    wsManager.disconnect(roomID);
    return api.post(`/room/leave/${roomID}`);
  },

  // WebSocket utility methods
  connectToRoom: (roomId: string, token?: string): WebSocket => {
    return wsManager.connect(roomId, token);
  },

  disconnectFromRoom: (roomId: string): void => {
    wsManager.disconnect(roomId);
  },

  getRoomConnection: (roomId: string): WebSocket | undefined => {
    return wsManager.getConnection(roomId);
  },

  isRoomConnected: (roomId: string): boolean => {
    return wsManager.isConnected(roomId);
  },

  // Cleanup method for when the app unmounts
  cleanup: (): void => {
    wsManager.disconnectAll();
  }
};
