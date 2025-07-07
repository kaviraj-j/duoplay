import { useNavigate, useParams } from "react-router-dom";
import TicTacToePage from "@/pages/games/tictactoe";
import { useRoom } from "@/contexts/RoomContext";

const gameComponents: Record<string, React.FC> = {
  tictactoe: TicTacToePage,
};

const GamesPage = () => {
  const navigate = useNavigate();

  const { gameName } = useParams();

  const { room } = useRoom();

  if (!room) {
    navigate("/");
  }

  const GameComponent = gameComponents[gameName as keyof typeof gameComponents];

  if (!GameComponent) {
    return navigate("/");
  }

  return <GameComponent />;
};

export default GamesPage;
