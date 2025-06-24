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

export type Room = {
  id: string;
  players: Record<string, Player>;
  gameName?: string;
};
