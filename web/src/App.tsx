import { AuthProvider } from "@/contexts/AuthContext";
import "./App.css";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import HomePage from "@/pages/HomePage";
import GamePage from "@/pages/games";
import { ProtectedRoute } from "@/components/auth/ProtectedRoutes";
import { Layout } from "@/components/layout/Layout";
import { RoomProvider } from "@/contexts/RoomContext";
function App() {
  return (
    <AuthProvider>
      <RoomProvider>
        <BrowserRouter>
          <Routes>
            <Route element={<Layout />}>
              <Route path="/" element={<HomePage />}></Route>
              <Route
                path="/game/:gameName"
                element={
                  <ProtectedRoute>
                    <GamePage />
                  </ProtectedRoute>
                }
              ></Route>
              <Route path="/*" element={<HomePage />}></Route>
            </Route>
          </Routes>
        </BrowserRouter>
      </RoomProvider>
    </AuthProvider>
  );
}

export default App;
