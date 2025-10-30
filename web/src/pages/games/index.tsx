import { useNavigate, useParams } from "react-router-dom";
import TicTacToePage from "@/pages/games/tictactoe";
import { useRoom } from "@/contexts/RoomContext";
import { useEffect } from "react";

const gameComponents: Record<string, React.FC> = {
  tictactoe: TicTacToePage,
};

const GamesPage = () => {
  const navigate = useNavigate();
  const { gameName } = useParams();
  const { room } = useRoom();

  useEffect(() => {
    if (!room) {
      navigate("/");
      return;
    }

    const GameComponent = gameComponents[gameName as keyof typeof gameComponents];
    if (!GameComponent) {
      navigate("/");
      return;
    }
  }, [room, gameName, navigate]);

  if (!room) {
    return null;
  }

  const GameComponent = gameComponents[gameName as keyof typeof gameComponents];
  if (!GameComponent) {
    return null;
  }

  return <GameComponent />;
};

export default GamesPage;
