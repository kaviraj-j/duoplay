import { useParams } from "react-router-dom";
import TicTacToePage from "@/pages/games/tictactoe";

const gameComponents: Record<string, React.FC> = {
  tictactoe: TicTacToePage,
};

const GamesPage = () => {
  

  const { gameName } = useParams();

  if (!gameName) {
    return <div>Invalid game</div>;
  }

  const GameComponent = gameComponents[gameName];

  if (!GameComponent) {
    return <div>Invalid game</div>;
  }

  return <GameComponent />;
};

export default GamesPage;
