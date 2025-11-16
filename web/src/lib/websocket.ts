// WebSocket connection manager interface
export interface IWebSocketManager {
  connect(roomId: string, token?: string): WebSocket;
  disconnect(roomId: string): void;
  getConnection(roomId: string): WebSocket | undefined;
  isConnected(roomId: string): boolean;
  addConnection(roomId: string, ws: WebSocket): void;
  disconnectAll(): void;
  createRoomConnection(messageHandler?: (event: MessageEvent) => void): Promise<{ roomId: string; ws: WebSocket }>;
  joinRoom(roomId: string, messageHandler?: (event: MessageEvent) => void): Promise<{ roomId: string; ws: WebSocket }>;
  setMessageHandler(roomId: string, handler: (event: MessageEvent) => void): void;
  sendMessage(roomId: string, message: any): void;
}

// WebSocket connection manager implementation
class WebSocketManager implements IWebSocketManager {
  private connections: Map<string, WebSocket> = new Map();
  private messageHandlers: Map<string, (event: MessageEvent) => void> = new Map();
  private baseUrl: string;

  constructor() {
    // Use the same base URL as the API but with ws:// or wss://
    const apiUrl = import.meta.env.VITE_API_URL || "http://localhost:8080";
    this.baseUrl = apiUrl.replace(/^http/, "ws");
  }

  connect(roomId: string, token?: string): WebSocket {
    // Close existing connection if any
    this.disconnect(roomId);

    // Create new WebSocket connection
    const wsUrl = `${this.baseUrl}/ws/room/${roomId}`;
    const ws = new WebSocket(wsUrl);

    // Add authentication header if token is provided
    if (token) {
      ws.addEventListener("open", () => {
        ws.send(JSON.stringify({ type: "auth", token }));
      });
    }

    // Store the connection
    this.connections.set(roomId, ws);

    // Handle connection close
    ws.addEventListener("close", () => {
      this.connections.delete(roomId);
    });

    return ws;
  }

  disconnect(roomId: string): void {
    const ws = this.connections.get(roomId);
    if (ws) {
      ws.close();
      this.connections.delete(roomId);
      this.messageHandlers.delete(roomId);
    }
  }

  getConnection(roomId: string): WebSocket | undefined {
    return this.connections.get(roomId);
  }

  isConnected(roomId: string): boolean {
    const ws = this.connections.get(roomId);
    return ws ? ws.readyState === WebSocket.OPEN : false;
  }

  // Add connection to the manager (used by createRoom)
  addConnection(roomId: string, ws: WebSocket): void {
    this.connections.set(roomId, ws);
  }

  // Close all connections
  disconnectAll(): void {
    this.connections.forEach((ws) => {
      ws.close();
    });
    this.connections.clear();
    this.messageHandlers.clear();
  }

  // Set message handler for a room connection
  setMessageHandler(roomId: string, handler: (event: MessageEvent) => void): void {
    const ws = this.connections.get(roomId);
    if (ws) {
      // Remove old handler if exists
      const oldHandler = this.messageHandlers.get(roomId);
      if (oldHandler) {
        ws.removeEventListener("message", oldHandler);
      }
      // Add new handler
      ws.addEventListener("message", handler);
      this.messageHandlers.set(roomId, handler);
    }
  }

  createRoomConnection(messageHandler?: (event: MessageEvent) => void): Promise<{ roomId: string; ws: WebSocket }> {
    const token = localStorage.getItem("duoplay_token");
    const wsUrl = `${this.baseUrl}/room/join?token=${token}`;
    const ws = new WebSocket(wsUrl);

    return new Promise((resolve, reject) => {
      const timeout = setTimeout(() => {
        reject(new Error("Connection timeout"));
      }, 10000);

      ws.addEventListener("open", () => {
        console.log("WebSocket connection opened");

        const tmpListener = (event: MessageEvent) => {
          try {
            const data = JSON.parse(event.data);

            if (data.type === "room_created") {
              const roomId = data.data.id;

              ws.removeEventListener("message", tmpListener);

              // Set the message handler if provided
              if (messageHandler) {
                ws.addEventListener("message", messageHandler);
                this.messageHandlers.set(roomId, messageHandler);
              }

              this.addConnection(roomId, ws);
              clearTimeout(timeout);
              resolve({ roomId, ws });
            } else if (data.type === "error") {
              clearTimeout(timeout);
              ws.removeEventListener("message", tmpListener);
              reject(new Error(data.message));
            }
          } catch (error) {
            console.error(error);
            clearTimeout(timeout);
            ws.removeEventListener("message", tmpListener);
            reject(new Error("Failed to parse server response"));
          }
        };

        ws.addEventListener("message", tmpListener);
      });

      ws.addEventListener("error", (error) => {
        console.error("WebSocket error:", error);
        reject(new Error("WebSocket connection failed"));
      });
    });
  }

  joinRoom(roomId: string, messageHandler?: (event: MessageEvent) => void): Promise<{ roomId: string; ws: WebSocket }> {
    const token = localStorage.getItem("duoplay_token");
    const wsUrl = `${this.baseUrl}/room/${roomId}/join?token=${token}`;
    const ws = new WebSocket(wsUrl);

    return new Promise((resolve, reject) => {
      let isResolved = false;
      const timeout = setTimeout(() => {
        if (!isResolved) {
          isResolved = true;
          reject(new Error("Connection timeout"));
        }
      }, 10000);

      ws.addEventListener("open", () => {

        const tmpListener = (event: MessageEvent) => {
          if (isResolved) return;
          
          try {
            const data = JSON.parse(event.data);

            if (data.type === "joined_room" || data.type === "room_joined") {
              // Server confirmed the join was successful
              isResolved = true;
              ws.removeEventListener("message", tmpListener);

              // Set the message handler if provided
              if (messageHandler) {
                ws.addEventListener("message", messageHandler);
                this.messageHandlers.set(roomId, messageHandler);
              }

              this.addConnection(roomId, ws);
              clearTimeout(timeout);
              resolve({ roomId, ws });
            } else if (data.type === "error") {
              isResolved = true;
              clearTimeout(timeout);
              ws.removeEventListener("message", tmpListener);
              ws.close();
              reject(new Error(data.message || "Failed to join room"));
            }
          } catch (error) {
            if (isResolved) return;
            console.error(error);
            isResolved = true;
            clearTimeout(timeout);
            ws.removeEventListener("message", tmpListener);
            reject(new Error("Failed to parse server response"));
          }
        };

        ws.addEventListener("message", tmpListener);
      });

      ws.addEventListener("error", (error) => {
        if (isResolved) return;
        console.error("WebSocket error:", error);
        isResolved = true;
        clearTimeout(timeout);
        reject(new Error("WebSocket connection failed"));
      });

      ws.addEventListener("close", (event) => {
        // If connection closes before we receive confirmation, reject the promise
        if (!isResolved) {
          isResolved = true;
          clearTimeout(timeout);
          reject(new Error(`Connection closed unexpectedly: ${event.code} ${event.reason || "Unknown reason"}`));
        }
      });
    });
  }

  sendMessage(roomId: string, message: any): void {
    const ws = this.connections.get(roomId);
        if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify(message));
    } else {
      console.error(`WebSocket for room ${roomId} is not connected.`);
    }
  }
}

// Create and export a singleton instance
export const wsManager: IWebSocketManager = new WebSocketManager();

// Export the class for testing or custom instances
export { WebSocketManager };
