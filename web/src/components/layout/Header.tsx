import { LoginModal } from '@/components/auth/LoginModal';
import { LogoutModal } from '@/components/auth/LogoutModal';
import { LoadingSpinner } from '@/components/common/LoadingSpinner';
import { useAuthContext } from '@/contexts/AuthContext';

const Header = () => {
    const { user, isLoading } = useAuthContext();

    if (isLoading) {
      return <LoadingSpinner />;
    }
  
  return (
    <header className="bg-white shadow-sm">
        <nav className="container mx-auto px-4 py-3 flex justify-between items-center">
          <h1 className="text-xl font-bold">DuoPlay</h1>
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
  )
}

export default Header