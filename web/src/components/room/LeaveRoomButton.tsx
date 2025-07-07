import { roomApi } from "@/api";

const LeaveRoomButton = ({ roomId }: { roomId: string }) => {
  return <button onClick={() => roomApi.leaveRoom(roomId)}>Leave Room</button>;
};

export default LeaveRoomButton;
