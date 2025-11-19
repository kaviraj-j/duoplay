import { useEffect, useCallback, useRef } from "react";
import { createRoomMessageHandler } from "@/handlers/roomHandler";
import { wsManager } from "@/lib/websocket";
import { roomApi } from "@/api";
import type { Room } from "@/types";
import type { GameChoiceData } from "@/contexts/RoomContext";


export const useRoomWebSocketHandler = ({
  removeRoom,
  saveRoom,
  updateRoom,
  room,
  setPendingGameChoice,
  navigate,
}: {
  removeRoom: () => void;
  saveRoom: (room: Room) => void;
  updateRoom: (updates: Partial<Room>) => void;
  room: Room | null;
  setPendingGameChoice?: (choice: GameChoiceData | null) => void;
  navigate?: (path: string) => void;
}) => {
  const createHandler = useCallback(() => {
    return createRoomMessageHandler({
      removeRoom,
      saveRoom,
      updateRoom,
      setPendingGameChoice,
      navigate,
    });
  }, [removeRoom, saveRoom, updateRoom, setPendingGameChoice, navigate]);

  const isReconnectingRef = useRef(false);

  useEffect(() => {
    const handler = createHandler();

    if (room?.id) {
      const ws = wsManager.getConnection(room.id);
      const isConnected = wsManager.isConnected(room.id);
      
      if (isConnected && ws) {
        // Connection exists and is open, just update the handler
        wsManager.setMessageHandler(room.id, handler);
      } else if (!isReconnectingRef.current) {
        // No connection or connection is closed - reconnect
        isReconnectingRef.current = true;
        console.log("Reconnecting to room:", room.id);
        
        roomApi.joinRoom(room.id, handler)
          .then(() => {
            console.log("Successfully reconnected to room:", room.id);
            isReconnectingRef.current = false;
          })
          .catch((error) => {
            console.error("Failed to reconnect to room:", error);
            isReconnectingRef.current = false;
            // If reconnection fails, remove the room from localStorage
            removeRoom();
          });
      }
    }
  }, [createHandler, room?.id, removeRoom]);

  return createHandler;
};

