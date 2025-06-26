import { useState } from "react";
import { useRoom } from "@/contexts/RoomContext";
import { useNavigate } from "react-router-dom";
// import { queueApi } from "@/api"; // Uncomment if you have a queue API

const ChooseOpponent = () => {
  const { createRoom, room, loading, error } = useRoom();
  const [inviteLink, setInviteLink] = useState<string | null>(null);
  const [queueStatus, setQueueStatus] = useState<"idle" | "waiting" | "matched">("idle");
  const [queueError, setQueueError] = useState<string | null>(null);
  const navigate = useNavigate();

  // Handler for playing with a friend
  const handlePlayWithFriend = async () => {
    const createdRoom = await createRoom();
    if (createdRoom) {
      const link = `${window.location.origin}/room/${createdRoom.id}`;
      setInviteLink(link);
    }
  };

  // Handler for playing with a random player
  const handlePlayWithRandom = async () => {
    setQueueStatus("waiting");
    setQueueError(null);
    try {
      // TODO: implement actual queue logic
    } catch (err: any) {
      setQueueError("Failed to join queue. Please try again.");
      setQueueStatus("idle");
    }
  };

  return (
    <div className="choose-opponent">
      <h2>Choose Opponent</h2>
      <div style={{ marginBottom: 16 }}>
        <button onClick={handlePlayWithFriend} disabled={loading}>
          Play with a Friend
        </button>
        {inviteLink && (
          <div style={{ marginTop: 8 }}>
            <p>Share this link with your friend:</p>
            <input
              type="text"
              value={inviteLink}
              readOnly
              style={{ width: "100%" }}
              onFocus={e => e.target.select()}
            />
          </div>
        )}
      </div>
      <div>
        <button onClick={handlePlayWithRandom} disabled={queueStatus === "waiting"}>
          Play with Random Player
        </button>
        {queueStatus === "waiting" && <p>Waiting for a random opponent...</p>}
        {queueError && <p style={{ color: "red" }}>{queueError}</p>}
      </div>
      {error && <p style={{ color: "red" }}>{error}</p>}
    </div>
  );
};

export default ChooseOpponent;