import { AuthProvider } from "@/contexts/AuthContext";
import "./App.css";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import HomePage from "@/pages/HomePage";
import GamePage from "@/pages/games";
import { ProtectedRoute } from "@/components/auth/ProtectedRoutes";
import { Layout } from "@/components/layout/Layout";
import { RoomProvider, useRoom } from "@/contexts/RoomContext";
import JoinRoom from "@/components/room/JoinRoom";
import { Toaster } from "react-hot-toast";
import { GameChoiceModal } from "@/components/room/GameChoiceModal";

function AppContent() {
  const { pendingGameChoice, setPendingGameChoice } = useRoom();

  return (
    <>
      <BrowserRouter>
        <Routes>
          <Route element={<Layout />}>
            <Route path="/" element={<HomePage />}></Route>
            <Route
              path="/room/:roomUid/join"
              element={
                <ProtectedRoute>
                  <JoinRoom />
                </ProtectedRoute>
              }
            ></Route>
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
      <GameChoiceModal
        open={pendingGameChoice !== null}
        gameChoice={pendingGameChoice}
        onClose={() => setPendingGameChoice(null)}
      />
    </>
  );
}

function App() {
  return (
    <>
      <Toaster position="top-center" />
      <AuthProvider>
        <RoomProvider>
          <AppContent />
        </RoomProvider>
      </AuthProvider>
    </>
  );
}

export default App;
