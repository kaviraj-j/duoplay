import { useEffect, useState } from "react";
import { useParams, useNavigate } from "react-router-dom";
import { roomApi } from "@/api/room";
import { useRoom } from "@/contexts/RoomContext";
import { LoadingSpinner } from "@/components/common/LoadingSpinner";
import { useRoomWebSocketHandler } from "@/hooks/useRoomWebSocketHandler";

const JoinRoom = () => {
  console.log("JoinRoom");
  const { roomUid } = useParams<{ roomUid: string }>();
  const navigate = useNavigate();
  const { saveRoom, removeRoom, updateRoom, room } = useRoom();
  const createHandler = useRoomWebSocketHandler({ removeRoom, saveRoom, updateRoom, room, navigate });
  const [isJoining, setIsJoining] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const handleJoinRoom = async () => {
      if (!roomUid) {
        setError("Room ID is required");
        return;
      }

      setIsJoining(true);
      setError(null);

      try {
        // Create handler with current room context and pass it to joinRoom
        const handler = createHandler();
        
        // Call the API to join the room with handler
        const response = await roomApi.joinRoom(roomUid, handler);
        
        if (response.roomId) {
          // Use the context to join the room
          const room = await roomApi.getRoom(response.roomId);
          saveRoom(room)
          if (room) {
            // Redirect to the game selection or game page
            navigate("/");
          } else {
            setError("Failed to join room");
          }
        }
      } catch (err: unknown) {
        console.error("Failed to join room:", err);
        setError(err instanceof Error ? err.message : "Failed to join room");
      } finally {
        setIsJoining(false);
      }
    };

    handleJoinRoom();
  }, [roomUid, navigate]);

  if (isJoining) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <LoadingSpinner />
          <p className="mt-4 text-lg">Joining room...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="text-center">
          <h2 className="text-2xl font-bold text-red-600 mb-4">Error</h2>
          <p className="text-lg text-gray-700 mb-6">{error}</p>
          <button
            onClick={() => navigate("/")}
            className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700"
          >
            Go Home
          </button>
        </div>
      </div>
    );
  }

  return null;
};

export default JoinRoom;
