import { useEffect, useState } from 'react';
import { useSearchParams, useNavigate } from 'react-router-dom';
import { api } from '@/lib/api';
import { useAuthStore } from '@/stores/authStore';

export function MagicLinkVerify() {
  var [searchParams] = useSearchParams();
  var navigate = useNavigate();
  var { fetchUser } = useAuthStore();
  var [error, setError] = useState('');
  var [isVerifying, setIsVerifying] = useState(true);

  useEffect(() => {
    var token = searchParams.get('token');

    if (!token) {
      setError('Missing verification token.');
      setIsVerifying(false);
      return;
    }

    var cancelled = false;

    async function verify() {
      try {
        await api.get(`/auth/magic-link/verify?token=${encodeURIComponent(token!)}`);
        if (cancelled) {
          return;
        }
        await fetchUser();
        navigate('/feed', { replace: true });
      } catch (err) {
        if (cancelled) {
          return;
        }
        if (err instanceof Error) {
          setError(err.message);
        } else {
          setError('Verification failed. The link may have expired.');
        }
        setIsVerifying(false);
      }
    }

    verify();

    return () => {
      cancelled = true;
    };
  }, [searchParams, fetchUser, navigate]);

  return (
    <div className="flex min-h-[50vh] items-center justify-center">
      <div className="card max-w-md text-center">
        {isVerifying ? (
          <>
            <div className="mx-auto mb-4 h-8 w-8 animate-spin rounded-full border-2 border-parchment-dark border-t-accent" />
            <p className="font-serif text-lg text-ink">Verifying your login...</p>
          </>
        ) : (
          <>
            <div className="mb-4 text-3xl text-error">!</div>
            <h2 className="mb-2 font-serif text-lg font-bold text-ink">Verification Failed</h2>
            <p className="text-sm text-feather">{error}</p>
          </>
        )}
      </div>
    </div>
  );
}
