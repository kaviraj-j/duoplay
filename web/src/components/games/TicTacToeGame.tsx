import React from "react";
import { Box, Typography, Paper } from "@mui/material";
import { useAuthContext } from "@/contexts/AuthContext";
import type { TicTacToeState, TicTacToeMove, GameStatus } from "@/types";

interface TicTacToeGameProps {
  gameState: TicTacToeState;
  onMove: (move: TicTacToeMove) => void;
  currentPlayerId: string;
  playerIds: string[];
  gameStatus: GameStatus;
}

const TicTacToeGame: React.FC<TicTacToeGameProps> = ({
  gameState,
  onMove,
  currentPlayerId,
  playerIds,
  gameStatus,
}) => {
  const { user } = useAuthContext();
  const isGameOver = gameStatus === "over";
  const isMyTurn = !isGameOver && user?.id === currentPlayerId;
  const myPlayerIndex = playerIds.indexOf(user?.id || "");
  const mySymbol = myPlayerIndex === 0 ? "X" : "O";

  const handleCellClick = (row: number, col: number) => {
    if (isGameOver || !isMyTurn) {
      return;
    }
    if (gameState.Board[row][col] !== "") {
      return;
    }
    onMove({ row, col });
  };

  const getCellValue = (row: number, col: number): string => {
    const cellValue = gameState.Board[row][col];
    if (!cellValue) return "";
    // Map player ID to symbol
    if (cellValue === playerIds[0]) {
      return "X";
    } else if (cellValue === playerIds[1]) {
      return "O";
    }
    return "";
  };

  return (
    <Box sx={{ display: "flex", flexDirection: "column", alignItems: "center", p: 3 }}>
      <Typography variant="h4" gutterBottom>
        Tic Tac Toe
      </Typography>

      <Typography variant="h6" gutterBottom sx={{ mb: 2 }}>
        {isGameOver
          ? "Game over"
          : isMyTurn
            ? "Your turn"
            : "Opponent's turn"}
      </Typography>

      <Typography variant="body1" sx={{ mb: 2 }}>
        You are: {mySymbol}
      </Typography>

      <Box
        sx={{
          display: "grid",
          gridTemplateColumns: "repeat(3, 1fr)",
          gap: 1,
          width: "fit-content",
          mb: 2,
        }}
      >
        {[0, 1, 2].map((row) =>
          [0, 1, 2].map((col) => (
            <Paper
              key={`${row}-${col}`}
              elevation={3}
              sx={{
                width: 100,
                height: 100,
                display: "flex",
                alignItems: "center",
                justifyContent: "center",
                cursor:
                  !isGameOver && isMyTurn && gameState.Board[row][col] === ""
                    ? "pointer"
                    : "default",
                backgroundColor:
                  !isGameOver && isMyTurn && gameState.Board[row][col] === ""
                    ? "#e3f2fd"
                    : "#fff",
                "&:hover": {
                  backgroundColor:
                    !isGameOver && isMyTurn && gameState.Board[row][col] === ""
                      ? "#bbdefb"
                      : "#fff",
                },
              }}
              onClick={() => handleCellClick(row, col)}
            >
              <Typography variant="h3" sx={{ fontWeight: "bold" }}>
                {getCellValue(row, col)}
              </Typography>
            </Paper>
          ))
        )}
      </Box>
    </Box>
  );
};

export default TicTacToeGame;

