import { useEffect, useCallback } from "react";
import { createRoomMessageHandler } from "@/handlers/roomHandler";
import { wsManager } from "@/lib/websocket";
import type { Room } from "@/types";


export const useRoomWebSocketHandler = ({
  removeRoom,
  saveRoom,
  updateRoom,
  room,
}: {
  removeRoom: () => void;
  saveRoom: (room: Room) => void;
  updateRoom: (updates: Partial<Room>) => void;
  room: Room | null;
}) => {
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

