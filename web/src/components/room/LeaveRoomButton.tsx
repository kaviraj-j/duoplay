import { roomApi } from "@/api";
import { useAuthContext } from "@/contexts/AuthContext";
import { useRoom } from "@/contexts/RoomContext";
import { Button } from "@mui/material";

const LeaveRoomButton = () => {
  const { removeRoom, room } = useRoom();
  const { user } = useAuthContext();
  console.log("User in LeaveRoomButton:", user);
  console.log("LeaveRoomButton rendered for room:", room);
  const opponent: string | null = room?.players
    ? Object.values(room.players).find((p) => p.user.id !== user?.id)?.user
        .name || null
    : null;
  console.log("Opponent in LeaveRoomButton:", opponent);
  const handleLeaveRoom = () => {
    roomApi
      .leaveRoom(room?.id || "")
      .then(() => {
        removeRoom();
      })
      .catch((err) => {
        console.error("Failed to leave room:", err);
      });
  };
  return (
    <Button
      onClick={handleLeaveRoom}
      className="flex items-center gap-3 p-3 rounded-lg bg-gray-100 hover:bg-gray-200 transition text-gray-900"
    >
      <span className="font-medium flex-1">
        {opponent ? `Playing against: ${opponent}` : "Waiting for friend"}
      </span>

      <span className="text-red-600 font-semibold">Leave</span>
    </Button>
  );
};

export default LeaveRoomButton;
