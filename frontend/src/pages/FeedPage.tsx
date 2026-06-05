import { useState } from 'react';
import { Navigate, Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import { useFeed } from '@/hooks/usePoems';
import { useAuthStore } from '@/stores/authStore';
import { PoemCard } from '@/components/poem/PoemCard';

export function FeedPage() {
  var { isAuthenticated, isLoading: authLoading } = useAuthStore();
  var [page, setPage] = useState(1);
  var { data, isLoading } = useFeed({ page, pageSize: 12 });
  var totalPages = data ? Math.ceil(data.totalCount / data.pageSize) : 0;

  if (authLoading) {
    return null;
  }

  if (!isAuthenticated) {
    return <Navigate to="/" replace />;
  }

  return (
    <div>
      <h1 className="mb-6 font-serif text-3xl font-bold text-ink">
        Your Feed
      </h1>

      {isLoading ? (
        <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="card animate-pulse">
              <div className="mb-3 h-6 w-3/4 rounded bg-parchment-dark" />
              <div className="mb-2 h-4 w-full rounded bg-parchment-dark" />
              <div className="h-4 w-2/3 rounded bg-parchment-dark" />
            </div>
          ))}
        </div>
      ) : data?.items.length === 0 ? (
        <div className="py-16 text-center">
          <p className="mb-4 text-lg text-feather">
            Follow poets to see their work here. Explore poems to discover new
            voices.
          </p>
          <Link to="/explore" className="btn-primary">
            Explore Poems
          </Link>
        </div>
      ) : (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3"
        >
          {data?.items.map((poem) => (
            <PoemCard key={poem.id} poem={poem} />
          ))}
        </motion.div>
      )}

      {totalPages > 1 && (
        <div className="mt-8 flex justify-center gap-4">
          <button
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page <= 1}
            className="btn-secondary disabled:opacity-40"
          >
            Previous
          </button>
          <span className="flex items-center font-sans text-sm text-feather">
            Page {page} of {totalPages}
          </span>
          <button
            onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
            disabled={page >= totalPages}
            className="btn-secondary disabled:opacity-40"
          >
            Next
          </button>
        </div>
      )}
    </div>
  );
}
