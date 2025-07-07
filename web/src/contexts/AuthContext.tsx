import React, {
  createContext,
  useContext,
  useState,
  useEffect,
  type ReactNode,
} from "react";
import type { AuthContextType, User } from "@/types";
import { userApi } from "@/api";
const AuthContext = createContext<AuthContextType | undefined>(undefined);

interface AuthProviderProps {
  children: ReactNode;
}

export const AuthProvider: React.FC<AuthProviderProps> = ({ children }) => {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  // Check for existing user session on mount
  useEffect(() => {
    (async () => {
      const savedUser = localStorage.getItem("duoplay_user");
      const token = localStorage.getItem("duoplay_token");
      const IsAuthenticated = !!(savedUser && token);

      if (!IsAuthenticated) {
        setUser(null);
        setIsLoading(false);
        return;
      }

      try {
        const validatedUser = await userApi.getCurrentUser();
        setUser(validatedUser);
      } catch (error) {
        console.error("Error validating user:", error);
        // Clear invalid session
        localStorage.removeItem("duoplay_user");
        localStorage.removeItem("duoplay_token");
        setUser(null);
      } finally {
        setIsLoading(false);
      }
    })();
  }, []);

  const login = (user: User, token: string) => {
    setUser(user);
    localStorage.setItem("duoplay_user", JSON.stringify(user));
    localStorage.setItem("duoplay_token", token);
  };

  const logout = () => {
    setUser(null);
    localStorage.removeItem("duoplay_user");
  };

  const value: AuthContextType = {
    user,
    login,
    logout,
    isLoading,
  };

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
};

export const useAuthContext = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuthContext must be used within an AuthProvider");
  }
  return context;
};
