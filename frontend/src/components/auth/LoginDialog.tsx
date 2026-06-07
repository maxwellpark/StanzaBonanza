import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { startAuthentication } from '@simplewebauthn/browser';
import { useUIStore } from '@/stores/uiStore';
import { useAuthStore } from '@/stores/authStore';
import { api } from '@/lib/api';
import { cn } from '@/lib/utils';
import type { User } from '@/types/user';

type Tab = 'magic-link' | 'passkey';

export function LoginDialog() {
  var { isLoginOpen, closeLogin } = useUIStore();

  var [activeTab, setActiveTab] = useState<Tab>('magic-link');
  var [email, setEmail] = useState('');
  var [isSubmitting, setIsSubmitting] = useState(false);
  var [successMessage, setSuccessMessage] = useState('');
  var [error, setError] = useState('');

  function resetState() {
    setEmail('');
    setIsSubmitting(false);
    setSuccessMessage('');
    setError('');
    setActiveTab('magic-link');
  }

  function handleClose() {
    closeLogin();
    resetState();
  }

  function handleBackdropClick(e: React.MouseEvent) {
    if (e.target === e.currentTarget) {
      handleClose();
    }
  }

  async function handleMagicLink(e: React.FormEvent) {
    e.preventDefault();
    setError('');
    setIsSubmitting(true);
    try {
      await api.post('/auth/magic-link', { email });
      setSuccessMessage('Check your email for a login link');
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Something went wrong. Please try again.');
    } finally {
      setIsSubmitting(false);
    }
  }

  async function handlePasskeyLogin() {
    setError('');
    setIsSubmitting(true);
    try {
      var options = await api.post<Parameters<typeof startAuthentication>[0]>('/auth/login/begin');
      var result = await startAuthentication({ optionsJSON: options as any });
      var user = await api.post<User>('/auth/login/finish', result);
      var { setUser } = useAuthStore.getState();
      setUser(user);
      handleClose();
    } catch (err) {
      if (err instanceof Error && err.name === 'NotAllowedError') {
        setError('Passkey authentication was cancelled.');
      } else {
        setError(err instanceof Error ? err.message : 'Passkey sign-in failed. Try magic link instead.');
      }
    } finally {
      setIsSubmitting(false);
    }
  }

  return (
    <AnimatePresence>
      {isLoginOpen && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.2 }}
          className="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 px-4"
          onClick={handleBackdropClick}
        >
          <motion.div
            initial={{ opacity: 0, scale: 0.95 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.95 }}
            transition={{ duration: 0.2 }}
            className="relative w-full max-w-md rounded-xl border border-parchment-dark bg-parchment p-6 shadow-lg"
          >
            <button
              onClick={handleClose}
              className="absolute right-4 top-4 text-feather transition-colors hover:text-ink"
              aria-label="Close"
            >
              <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>

            <h2 className="mb-4 font-serif text-xl font-bold text-ink">Welcome Back</h2>

            <div className="mb-4 flex gap-1 rounded-lg bg-parchment-dark p-1">
              {(['magic-link', 'passkey'] as Tab[]).map((tab) => (
                <button
                  key={tab}
                  onClick={() => {
                    setActiveTab(tab);
                    setError('');
                    setSuccessMessage('');
                  }}
                  className={cn(
                    'flex-1 rounded-md px-3 py-1.5 font-sans text-sm transition-colors',
                    activeTab === tab
                      ? 'bg-white text-ink shadow-sm'
                      : 'text-feather hover:text-ink'
                  )}
                >
                  {tab === 'magic-link' ? 'Magic Link' : 'Passkey'}
                </button>
              ))}
            </div>

            {activeTab === 'magic-link' && (
              <>
                {successMessage ? (
                  <div className="rounded-lg bg-success/10 p-4 text-center text-sm text-success">
                    {successMessage}
                  </div>
                ) : (
                  <form onSubmit={handleMagicLink} className="flex flex-col gap-3">
                    <input
                      type="email"
                      required
                      placeholder="your@email.com"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      className="rounded-lg border border-parchment-dark bg-white px-4 py-2 font-sans text-sm text-ink outline-none transition-colors focus:border-accent"
                    />
                    {error && <p className="text-sm text-error">{error}</p>}
                    <button
                      type="submit"
                      disabled={isSubmitting}
                      className="btn-primary disabled:opacity-50"
                    >
                      {isSubmitting ? 'Sending...' : 'Send Magic Link'}
                    </button>
                  </form>
                )}
              </>
            )}

            {activeTab === 'passkey' && (
              <div className="flex flex-col gap-3">
                <p className="font-sans text-sm text-feather">
                  Sign in instantly using a passkey saved on this device.
                </p>
                {error && <p className="text-sm text-error">{error}</p>}
                <button
                  onClick={handlePasskeyLogin}
                  disabled={isSubmitting}
                  className="btn-primary disabled:opacity-50"
                >
                  {isSubmitting ? 'Waiting for passkey...' : 'Sign in with Passkey'}
                </button>
                <p className="text-center font-sans text-xs text-feather">
                  No passkey yet? Sign in with magic link first, then register one from your profile.
                </p>
              </div>
            )}
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
}
