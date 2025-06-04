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
