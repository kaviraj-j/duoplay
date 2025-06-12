import Header from "@/components/layout/Header";
import { Outlet } from "react-router-dom";

export function Layout() {
  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <Header />
      {/* Main content */}
      <main className="container mx-auto px-4 py-8">
        <Outlet />
      </main>
    </div>
  );
}
