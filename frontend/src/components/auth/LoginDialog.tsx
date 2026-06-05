import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { useUIStore } from '@/stores/uiStore';
import { api } from '@/lib/api';
import { cn } from '@/lib/utils';

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
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('Something went wrong. Please try again.');
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
                    {error && (
                      <p className="text-sm text-error">{error}</p>
                    )}
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
                <input
                  type="email"
                  placeholder="your@email.com"
                  disabled
                  className="rounded-lg border border-parchment-dark bg-white px-4 py-2 font-sans text-sm text-ink outline-none opacity-50"
                />
                <div className="relative">
                  <button
                    disabled
                    className="btn-primary w-full disabled:opacity-50"
                    title="Coming soon"
                  >
                    Sign in with Passkey
                  </button>
                  <span className="mt-1 block text-center font-sans text-xs text-feather">
                    Coming soon
                  </span>
                </div>
              </div>
            )}
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
}
