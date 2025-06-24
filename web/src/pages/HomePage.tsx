import { gameApi } from "@/api";
import type { GameListPayload } from "@/types";
import { useEffect, useState } from "react";

const HomePage = () => {
  const [games, setGames] = useState<GameListPayload[]>([]);

  useEffect(() => {
    gameApi.getGamesList().then((gamesResponse) => {
      console.log("games", gamesResponse);
      setGames(gamesResponse.data);
    });
  }, []);

  return (
    <div>
      HomePage
      {games.map((game) => (
        <div key={game.name} className="game-card">
          <h2 className="game-title">{game.display_name}</h2>
          <p className="game-description">
            Play {game.display_name} with your friends!
          </p>
          <button className="join-game-button">
            <a href={`/game/${game.name}`}>Join Game</a>
          </button>
        </div>
      ))}
    </div>
  );
};

export default HomePage;
