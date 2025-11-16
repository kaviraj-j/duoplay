import { api, wsManager } from "@/lib";
import type { Room } from "@/types";

export const roomApi = {
  getRoom: async (roomID: string): Promise<Room> => {
    const response = await api.get<{ data: Room }>(`/room/${roomID}`);
    return response.data.data;
  },

  createRoom: async (
    messageHandler?: (event: MessageEvent) => void
  ): Promise<{ roomId: string; ws: WebSocket }> => {
    return wsManager.createRoomConnection(messageHandler);
  },

  joinRoom: async (
    roomId: string,
    messageHandler?: (event: MessageEvent) => void
  ): Promise<{ roomId: string; ws: WebSocket }> => {
    return wsManager.joinRoom(roomId, messageHandler);
  },
  joinQueue: () => api.post("/room/queue"),
  leaveQueue: () => api.post("/room/leaveQueue"),
  leaveRoom: (roomID: string) => {
    // Disconnect WebSocket when leaving room
    wsManager.disconnect(roomID);
    return api.get(`/room/${roomID}/leave`);
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
  },

  chooseGame: (roomId: string, gameName: string) => {
    return wsManager.sendMessage(roomId, {
      type: "choose_game",
      game_type: gameName ,
    });
  },

  sendMessage: (roomId: string, message: any): void => {
    return wsManager.sendMessage(roomId, message);
  },
};
