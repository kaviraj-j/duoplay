import { useRoom } from "@/contexts/RoomContext";

const tictactoe = () => {
  const { room } = useRoom();
  if (!room) {
    return <>No room found</>;
  }
  return (
    <div>tictactoe</div>
  )
}

export default tictactoe