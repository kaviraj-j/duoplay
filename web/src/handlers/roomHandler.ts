import { toastService } from "@/services/toastService";
import type { Room } from "@/types";
import type { GameChoiceData } from "@/contexts/RoomContext";

export type RoomMessageHandlerCallbacks = {
  removeRoom: () => void;
  saveRoom?: (room: Room) => void;
  updateRoom?: (updates: Partial<Room>) => void;
  setPendingGameChoice?: (choice: GameChoiceData | null) => void;
};

export const createRoomMessageHandler = (
  callbacks: RoomMessageHandlerCallbacks
) => {
  return (event: MessageEvent) => {
    try {
      const data = JSON.parse(event.data);

      // Display message via toast if available
      if (data.message) {
        toastService.show({
          message: data.message,
          type: data.type || "info",
          dismissable: true,
          durationMs: 1500,
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
          // Clear pending game choice when game starts
          if (callbacks.setPendingGameChoice) {
            callbacks.setPendingGameChoice(null);
          }
          break;
        case "game_chosen":
          console.log("Game chosen:", data.data);
          if (data.data && callbacks.setPendingGameChoice) {
            callbacks.setPendingGameChoice(data.data as GameChoiceData);
          }
          break;
        case "game_accepted":
        case "game_rejected":
          // Clear pending game choice when game is accepted or rejected
          if (callbacks.setPendingGameChoice) {
            callbacks.setPendingGameChoice(null);
          }
          break;

        default:
          break;
      }
    } catch (error: unknown) {
      console.error("Failed to parse WebSocket message:", error);
    }
  };
};
