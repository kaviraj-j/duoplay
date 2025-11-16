import { roomApi } from "@/api";
import { useRoom } from "@/contexts/RoomContext";

const LeaveRoomButton = ({ roomId }: { roomId: string }) => {
  const { removeRoom } = useRoom();
  const handleLeaveRoom = () => {
    roomApi
      .leaveRoom(roomId)
      .then(() => {
        removeRoom();
      })
      .catch((err) => {
        console.error("Failed to leave room:", err);
      });
  };
  return <button onClick={handleLeaveRoom}>Leave Room</button>;
};

export default LeaveRoomButton;
