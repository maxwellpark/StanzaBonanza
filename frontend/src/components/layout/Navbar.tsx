import { useState } from 'react';
import { Link } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { useAuthStore } from '@/stores/authStore';
import { useUIStore } from '@/stores/uiStore';
import { cn } from '@/lib/utils';

var navLinks = [
  { to: '/explore', label: 'Explore' },
  { to: '/hall-of-fame', label: 'Hall of Fame' },
  { to: '/tutorials', label: 'Tutorials' },
];

export function Navbar() {
  var [mobileOpen, setMobileOpen] = useState(false);
  var { user, isAuthenticated } = useAuthStore();
  var { openLogin } = useUIStore();

  return (
    <nav className="sticky top-0 z-50 border-b border-parchment-dark bg-white/80 backdrop-blur">
      <div className="mx-auto flex max-w-6xl items-center justify-between px-4 py-3">
        <Link to="/" className="font-serif text-xl font-bold text-ink no-underline">
          Stanza Bonanza
        </Link>

        <div className="hidden items-center gap-6 md:flex">
          {navLinks.map((link) => (
            <Link
              key={link.to}
              to={link.to}
              className="font-sans text-sm text-feather transition-colors hover:text-ink"
            >
              {link.label}
            </Link>
          ))}
        </div>

        <div className="hidden items-center gap-4 md:flex">
          {isAuthenticated && user ? (
            <div className="flex items-center gap-3">
              <button className="relative text-feather transition-colors hover:text-ink">
                <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
                </svg>
              </button>
              <Link to={`/profile/${user.id}`} className="flex items-center gap-2 no-underline">
                <img
                  src={user.avatarUrl || '/default-avatar.png'}
                  alt=""
                  className="h-8 w-8 rounded-full object-cover"
                />
                <span className="font-sans text-sm text-ink">{user.displayName}</span>
              </Link>
            </div>
          ) : (
            <button onClick={openLogin} className="btn-primary text-sm">
              Sign In
            </button>
          )}
        </div>

        <button
          className="flex flex-col gap-1 md:hidden"
          onClick={() => setMobileOpen(!mobileOpen)}
          aria-label="Toggle menu"
        >
          <span className={cn('block h-0.5 w-5 bg-ink transition-transform', mobileOpen && 'translate-y-1.5 rotate-45')} />
          <span className={cn('block h-0.5 w-5 bg-ink transition-opacity', mobileOpen && 'opacity-0')} />
          <span className={cn('block h-0.5 w-5 bg-ink transition-transform', mobileOpen && '-translate-y-1.5 -rotate-45')} />
        </button>
      </div>

      <AnimatePresence>
        {mobileOpen && (
          <motion.div
            initial={{ height: 0, opacity: 0 }}
            animate={{ height: 'auto', opacity: 1 }}
            exit={{ height: 0, opacity: 0 }}
            transition={{ duration: 0.2 }}
            className="overflow-hidden border-t border-parchment-dark md:hidden"
          >
            <div className="flex flex-col gap-2 px-4 py-3">
              {navLinks.map((link) => (
                <Link
                  key={link.to}
                  to={link.to}
                  onClick={() => setMobileOpen(false)}
                  className="font-sans text-sm text-feather transition-colors hover:text-ink"
                >
                  {link.label}
                </Link>
              ))}
              {isAuthenticated && user ? (
                <Link
                  to={`/profile/${user.id}`}
                  onClick={() => setMobileOpen(false)}
                  className="flex items-center gap-2 pt-2 no-underline"
                >
                  <img
                    src={user.avatarUrl || '/default-avatar.png'}
                    alt=""
                    className="h-8 w-8 rounded-full object-cover"
                  />
                  <span className="font-sans text-sm text-ink">{user.displayName}</span>
                </Link>
              ) : (
                <button
                  onClick={() => {
                    setMobileOpen(false);
                    openLogin();
                  }}
                  className="btn-primary mt-2 text-sm"
                >
                  Sign In
                </button>
              )}
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </nav>
  );
}
