import { createContext, useContext, useEffect, useState } from "react";
import type { ReactNode } from "react";
import { roomApi } from "@/api"; // You should have these API functions implemented
import type { Room } from "@/types";

type RoomContextType = {
  room: Room | null;
  createRoom: (roomData: Room) => Promise<Room | null>;
  joinRoom: (roomID: string) => Promise<Room | null>;
  leaveRoom: () => Promise<void>;
};

const RoomContext = createContext<RoomContextType | undefined>(undefined);

export const RoomProvider = ({ children }: { children: ReactNode }) => {
  const [room, setRoom] = useState<Room | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchRoom = async (roomID: string): Promise<Room | null> => {
      try {
        const response = await roomApi.getRoom(roomID);
        const roomData = response.data as Room;
        return roomData;
      } catch (err: any) {
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
        }
      });
    }
  }, []);

  const createRoom = async (roomData: Room): Promise<Room | null> => {
    setLoading(true);
    setError(null);
    try {
      setRoom(roomData);
      localStorage.setItem("duoplay_room", JSON.stringify(roomData));
      return roomData;
    } catch (err: any) {
      setError((err as Error)?.message || "Failed to create room");
      return null;
    } finally {
      setLoading(false);
    }
  };

  const joinRoom = async (roomID: string): Promise<Room | null> => {
    setLoading(true);
    setError(null);
    try {
      const response = await roomApi.joinRoom(roomID);
      const roomData = response.data as Room;
      setRoom(roomData);
      return roomData;
    } catch (err: any) {
      setError((err as Error)?.message || "Failed to join room");
      return null;
    } finally {
      setLoading(false);
    }
  };

  const leaveRoom = async (): Promise<void> => {
    const roomID = room?.id;
    if (!roomID) {
      return;
    }
    setLoading(true);
    setError(null);

    setRoom(null);
    localStorage.removeItem("duoplay_room");
  };

  return (
    <RoomContext.Provider value={{ room, createRoom, joinRoom, leaveRoom }}>
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
