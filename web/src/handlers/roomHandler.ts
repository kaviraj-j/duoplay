import { toastService } from "@/services/toastService";
import type { Room } from "@/types";

export type RoomMessageHandlerCallbacks = {
  removeRoom: () => void;
  saveRoom?: (room: Room) => void;
  updateRoom?: (updates: Partial<Room>) => void;
};

export const createRoomMessageHandler = (
  callbacks: RoomMessageHandlerCallbacks
) => {
  return (event: MessageEvent) => {
    try {
      const data = JSON.parse(event.data);
      console.log("WebSocket message received:", data);

      // Display message via toast if available
      if (data.message) {
        toastService.show({
          message: data.message,
          type: data.type || "info",
          dismissable: true,
          durationMs: 5000,
        });
      }

      // Handle different message types
      switch (data.type) {
        case "opponent_left":
          callbacks.removeRoom();
          break;

        case "room_updated":
          if (data.data && callbacks.updateRoom) {
            callbacks.updateRoom(data.data);
          } else if (data.data && callbacks.saveRoom) {
            callbacks.saveRoom(data.data as Room);
          }
          break;

        case "room_joined":
          if (data.data && callbacks.saveRoom) {
            callbacks.saveRoom(data.data as Room);
          }
          break;

        case "game_started":
          if (data.data && callbacks.updateRoom) {
            callbacks.updateRoom({ gameName: data.data.gameName });
          }
          break;

        default:
          // Handle other message types here
          break;
      }
    } catch (error: unknown) {
      console.error("Failed to parse WebSocket message:", error);
    }
  };
};
