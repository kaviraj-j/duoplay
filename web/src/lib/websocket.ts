// WebSocket connection manager interface
export interface IWebSocketManager {
  connect(roomId: string, token?: string): WebSocket;
  disconnect(roomId: string): void;
  getConnection(roomId: string): WebSocket | undefined;
  isConnected(roomId: string): boolean;
  addConnection(roomId: string, ws: WebSocket): void;
  disconnectAll(): void;
  createRoomConnection(): Promise<{ roomId: string; ws: WebSocket }>;
}

// WebSocket connection manager implementation
class WebSocketManager implements IWebSocketManager {
  private connections: Map<string, WebSocket> = new Map();
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
  }

  // Create a room WebSocket connection (for room creation)
  createRoomConnection(): Promise<{ roomId: string; ws: WebSocket }> {
    const token = localStorage.getItem("duoplay_token");
    const wsUrl = `${this.baseUrl}/room/join?token=${token}`;
    const ws = new WebSocket(wsUrl);
  
    return new Promise((resolve, reject) => {
      // Wait for connection to open FIRST
      ws.addEventListener("open", () => {
        console.log("WebSocket connection opened");
        
        // NOW start listening for messages
        ws.addEventListener("message", (event) => {
          try {
            console.log("Received message:", event.data);
            const data = JSON.parse(event.data);
            if (data.type === "room_created") {
              const roomId = data.data.id;
              this.addConnection(roomId, ws);
              resolve({ roomId, ws });
            } else if (data.type === "error") {
              reject(new Error(data.message));
            }
          } catch (error: unknown) {
            console.error(error)
            reject(new Error("Failed to parse server response"));
          }
        });
      });
  
      ws.addEventListener("error", (error) => {
        console.error("WebSocket error:", error);
        reject(new Error("WebSocket connection failed"));
      });
  
      // Timeout after 10 seconds
      setTimeout(() => {
        reject(new Error("Connection timeout"));
      }, 10000);
    });
  }
}

// Create and export a singleton instance
export const wsManager: IWebSocketManager = new WebSocketManager();

// Export the class for testing or custom instances
export { WebSocketManager };
