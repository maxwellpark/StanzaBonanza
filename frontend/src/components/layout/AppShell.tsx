import { Outlet } from 'react-router-dom';
import { Navbar } from './Navbar';
import { Footer } from './Footer';
import { LoginDialog } from '@/components/auth/LoginDialog';

export function AppShell() {
  return (
    <div className="flex min-h-screen flex-col">
      <Navbar />
      <LoginDialog />
      <main className="mx-auto w-full max-w-6xl flex-1 px-4 py-6">
        <Outlet />
      </main>
      <Footer />
    </div>
  );
}
