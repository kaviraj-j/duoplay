import { useEffect, useCallback } from "react";
import { useRoom } from "@/contexts/RoomContext";
import { createRoomMessageHandler } from "@/handlers/roomHandler";
import { wsManager } from "@/lib/websocket";

/**
 * Custom hook that sets up WebSocket message handlers for room connections.
 * This hook ensures that all existing and future room WebSocket connections
 * use the message handler with the current room context functions.
 * 
 * Returns a function to create a message handler that can be used when
 * creating new room connections.
 */
export const useRoomWebSocketHandler = () => {
  const { removeRoom, saveRoom, updateRoom, room } = useRoom();

  const createHandler = useCallback(() => {
    return createRoomMessageHandler({
      removeRoom,
      saveRoom,
      updateRoom,
    });
  }, [removeRoom, saveRoom, updateRoom]);

  useEffect(() => {
    const handler = createHandler();

    if (room?.id) {
      const ws = wsManager.getConnection(room.id);
      if (ws) {
        wsManager.setMessageHandler(room.id, handler);
      }
    }
  }, [createHandler, room?.id]);

  return createHandler;
};

