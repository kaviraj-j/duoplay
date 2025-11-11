import { useAuthContext } from "@/contexts/AuthContext";
import { LoadingSpinner } from "@/components/common/LoadingSpinner";
import { Navigate } from "react-router-dom";
import { toastService } from "@/services/toastService";

interface ProtectedRouteProps {
  children: React.ReactNode;
}

export const ProtectedRoute = ({ children }: ProtectedRouteProps) => {
  const { user, isLoading } = useAuthContext();

  if (isLoading) {
    return <LoadingSpinner />;
  }

  if (!user) {
    console.log("ProtectedRoute: user not authenticated");
    toastService.show({
      message: "Please log in to access this page.",
      type: "error",
    });
    return <Navigate to="/" replace />;
  }

  return <>{children}</>;
};
