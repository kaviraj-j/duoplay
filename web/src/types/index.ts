export interface User {
  id: string;
  name: string;
}

export interface NewUserPayload {
  name: string;
}

export interface AuthContextType {
  user: User | null;
  login: (user: User, token: string) => void;
  logout: () => void;
  isLoading: boolean;
}

export interface GameListPayload {
  name: string;
  display_name: string;
}
export type Player = {
  user: User;
};

export type GameType = "tictactoe";

export type GameStatus = "not_started" | "in_progress" | "over";

export interface GameState {
  type: GameType;
  status: GameStatus;
  state: TicTacToeState | Record<string, unknown>;
}

export type Room = {
  id: string;
  players: Record<string, Player>;
  gameName?: string;
  game?: GameState;
  status?: string;
  game_selection?: Record<string, GameType>;
};

// Game interface for frontend game implementations
export interface Game {
  type: GameType;
  render: () => React.ReactNode;
  handleMove: (move: TicTacToeMove | Record<string, unknown>) => void;
  getState: () => TicTacToeState | Record<string, unknown>;
  updateState: (state: TicTacToeState | Record<string, unknown>) => void;
}

// TicTacToe specific types
export interface TicTacToeState {
  Board: string[][];
  CurrentPlayer: string;
}

export interface TicTacToeMove {
  row: number;
  col: number;
}
