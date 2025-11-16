import { createContext, useContext, useEffect, useState } from "react";
import type { ReactNode } from "react";
import { roomApi } from "@/api"; // You should have these API functions implemented
import type { Room } from "@/types";
import { useRoomWebSocketHandler } from "@/hooks/useRoomWebSocketHandler";

export interface GameChoiceData {
  game_type: string;
  player_id: string;
  player_name: string;
}

type RoomContextType = {
  room: Room | null;
  saveRoom: (room: Room) => void;
  removeRoom: () => void;
  updateRoom: (updates: Partial<Room>) => void;
  pendingGameChoice: GameChoiceData | null;
  setPendingGameChoice: (choice: GameChoiceData | null) => void;
};


const RoomContext = createContext<RoomContextType | undefined>(undefined);

export const RoomProvider = ({ children }: { children: ReactNode }) => {
  const [room, setRoom] = useState<Room | null>(null);
  const [pendingGameChoice, setPendingGameChoice] = useState<GameChoiceData | null>(null);


  useEffect(() => {
    const fetchRoom = async (roomID: string): Promise<Room | null> => {
      try {
        const response = await roomApi.getRoom(roomID);
        const roomData = response as Room;
        return roomData;
      } catch (err: unknown) {
        console.error(err);
        return null;
      }
    };

    const roomDetails = localStorage.getItem("duoplay_room");
    if (roomDetails) {
      const roomData = JSON.parse(roomDetails) as Room;
      // validate room
      fetchRoom(roomData.id).then((room) => {
        if (room) {
          setRoom(room);
        } else {
          localStorage.removeItem("duoplay_room");
        }
      });
    } else {
      localStorage.removeItem("duoplay_room");
    }
  }, []);

  const saveRoom = (room: Room) => {
    localStorage.setItem("duoplay_room", JSON.stringify(room));
    setRoom(room);
  };

  const removeRoom = () => {
    localStorage.removeItem("duoplay_room");
    setRoom(null);
  };

  const updateRoom = (updates: Partial<Room>) => {
    setRoom((currentRoom) => {
      if (!currentRoom) return currentRoom;
      const updatedRoom = { ...currentRoom, ...updates };
      localStorage.setItem("duoplay_room", JSON.stringify(updatedRoom));
      return updatedRoom;
    });
  };

  useRoomWebSocketHandler({ removeRoom, saveRoom, updateRoom, room, setPendingGameChoice });

  return (
    <RoomContext.Provider value={{ room, saveRoom, removeRoom, updateRoom, pendingGameChoice, setPendingGameChoice }}>
      {children}
    </RoomContext.Provider>
  );
};

export const useRoom = () => {
  const context = useContext(RoomContext);
  if (!context) {
    throw new Error("useRoom must be used within a RoomProvider");
  }
  return context;
};
