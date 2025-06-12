const HomePage = () => {
  return <div>HomePage

    {["tic-tac-toe"].map((game) => (
      <div key={game} className="game-card">
        <h2 className="game-title">{game}</h2>
        <p className="game-description">Play {game} with your friends!</p>
        <button className="join-game-button"><a href={`/game/${game}`}>Join Game</a></button>
      </div>
    ))}
  </div>;
};

export default HomePage;
