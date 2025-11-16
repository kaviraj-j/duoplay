import { LoginModal } from "@/components/auth/LoginModal";
import { LogoutModal } from "@/components/auth/LogoutModal";
import { LoadingSpinner } from "@/components/common/LoadingSpinner";
import ChooseOpponent from "@/components/room/ChooseOpponent";
import LeaveRoomButton from "@/components/room/LeaveRoomButton";
import { useAuthContext } from "@/contexts/AuthContext";
import { useRoom } from "@/contexts/RoomContext";

const Header = () => {
  const { user, isLoading } = useAuthContext();
  const { room } = useRoom();

  if (isLoading) {
    return <LoadingSpinner />;
  }

  return (
    <header className="bg-white shadow-sm">
      <nav className="container mx-auto px-4 py-3 flex justify-between items-center">
        <h1 className="text-xl font-bold">DuoPlay</h1>
        {user && (
          <div>
            {room ? (
              <div className="flex items-center gap-4">
                <LeaveRoomButton />
              </div>
            ) : (
              <ChooseOpponent />
            )}
          </div>
        )}
        <div>
          {user ? (
            <div className="flex items-center gap-4">
              <span>Welcome, {user.name}</span>
              <LogoutModal />
            </div>
          ) : (
            <LoginModal />
          )}
        </div>
      </nav>
    </header>
  );
};

export default Header;
