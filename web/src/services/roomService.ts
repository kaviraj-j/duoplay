// src/services/roomService.ts
import { roomApi } from "@/api/room";
import type { Room } from "@/types";
import { wsManager } from "@/lib/websocket";

export class RoomService {
  private static instance: RoomService;
  private currentRoom: Room | null = null;

  private constructor() {}

  static getInstance(): RoomService {
    if (!this.instance) {
      this.instance = new RoomService();
    }
    return this.instance;
  }

  async joinRoom(roomId: string): Promise<Room> {
    try {
      // First establish WebSocket connection
      const { ws } = await roomApi.joinRoom(roomId);
      
      // Then get room details
      const room = await roomApi.getRoom(roomId);
      
      // Store room data
      this.currentRoom = room;
      localStorage.setItem("duoplay_room", JSON.stringify(room));
      
      // Setup WebSocket handlers
      this.setupWebSocketHandlers(ws, room);
      
      return room;
    } catch (error) {
      localStorage.removeItem("duoplay_room");
      this.currentRoom = null;
      throw error;
    }
  }

  async createRoom(): Promise<Room> {
    try {
      const { roomId, ws } = await roomApi.createRoom();
      const room = await roomApi.getRoom(roomId);
      
      this.currentRoom = room;
      localStorage.setItem("duoplay_room", JSON.stringify(room));
      
      this.setupWebSocketHandlers(ws, room);
      
      return room;
    } catch (error) {
      localStorage.removeItem("duoplay_room");
      this.currentRoom = null;
      throw error;
    }
  }

  async leaveRoom(roomId: string): Promise<void> {
    try {
      await roomApi.leaveRoom(roomId);
      localStorage.removeItem("duoplay_room");
      this.currentRoom = null;
    } catch (error) {
      console.error("Error leaving room:", error);
      throw error;
    }
  }

  private setupWebSocketHandlers(ws: WebSocket, room: Room) {
    ws.addEventListener("message", (event) => {
      try {
        const data = JSON.parse(event.data);
        // Handle different message types
        switch (data.type) {
          case "room_update":
            this.handleRoomUpdate(data.room);
            break;
          case "game_state":
            this.handleGameState(data.state);
            break;
          // Add more handlers as needed
        }
      } catch (error) {
        console.error("WebSocket message error:", error);
      }
    });
  }

  private handleRoomUpdate(roomData: Room) {
    this.currentRoom = roomData;
    localStorage.setItem("duoplay_room", JSON.stringify(roomData));
    // Emit room update event
  }

  private handleGameState(state: any) {
    // Handle game state updates
  }
}

export const roomService = RoomService.getInstance();