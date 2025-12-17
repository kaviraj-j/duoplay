import { useEffect, useState } from "react";
import { useRoom } from "@/contexts/RoomContext";
import { roomApi } from "@/api/room";
import TicTacToeGame from "@/components/games/TicTacToeGame";
import type { TicTacToeState, TicTacToeMove } from "@/types";

const TicTacToePage = () => {
  const { room } = useRoom();
  const [gameState, setGameState] = useState<TicTacToeState | null>(null);
  const [currentPlayerId, setCurrentPlayerId] = useState<string>("");
  const [playerIds, setPlayerIds] = useState<string[]>([]);
  const [gameStatus, setGameStatus] = useState<"not_started" | "in_progress" | "over">("not_started");

  useEffect(() => {
    if (room?.game?.state) {
      // Initialize game state from room
      const state = room.game.state as TicTacToeState;
      setGameState(state);
      setCurrentPlayerId(state.CurrentPlayer || "");

      if (room.game.status) {
        setGameStatus(room.game.status);
      }

      // Get player IDs from room
      const ids = Object.keys(room.players);
      setPlayerIds(ids);
    }
  }, [room]);

  const handleMove = (move: TicTacToeMove) => {
    if (!room) return;
    
    // Send move to server via WebSocket
    roomApi.sendMove(room.id, move);
  };

  if (!room) {
    return <>No room found</>;
  }

  if (!room.game) {
    return <>Game not started yet...</>;
  }

  if (!gameState) {
    return <>Loading game...</>;
  }

  return (
    <TicTacToeGame
      gameState={gameState}
      onMove={handleMove}
      currentPlayerId={currentPlayerId}
      playerIds={playerIds}
        gameStatus={gameStatus}
    />
  );
};

export default TicTacToePage;