import { AuthProvider } from "@/contexts/AuthContext";
import "./App.css";
import { BrowserRouter, Route, Routes } from "react-router-dom";
import HomePage from "@/pages/HomePage";
import GamePage from "@/pages/GamePage";
import { ProtectedRoute } from "@/components/auth/ProtectedRoutes";
import { Layout } from "@/components/layout/Layout";
function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route element={<Layout />}>
            <Route path="/" element={<HomePage />}></Route>
            <Route
              path="/game/:roomId"
              element={
                <ProtectedRoute>
                  <GamePage />
                </ProtectedRoute>
              }
            ></Route>
          </Route>
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}

export default App;
