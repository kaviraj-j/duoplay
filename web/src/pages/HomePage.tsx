import { gameApi } from "@/api";
import { useRoom } from "@/contexts/RoomContext";
import type { GameListPayload } from "@/types";
import { useEffect, useState } from "react";

const HomePage = () => {
  const [games, setGames] = useState<GameListPayload[]>([]);
  const { room } = useRoom();

  useEffect(() => {
    gameApi.getGamesList().then((gamesResponse) => {
      console.log("games", gamesResponse);
      setGames(gamesResponse.data);
    });
  }, []);

  return (
    <div>
      {games.map((game) => (
        <div key={game.name} className="game-card">
          <h2 className="game-title">{game.display_name}</h2>
          {room && (
            <button className="join-game-button">
              <a href={`/game/${game.name}`}>Play Game</a>
            </button>
          )}
        </div>
      ))}
    </div>
  );
};

export default HomePage;
