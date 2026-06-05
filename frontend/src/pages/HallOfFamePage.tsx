import { useState } from 'react';
import { motion } from 'framer-motion';
import { useHallOfFame } from '@/hooks/usePoems';
import { PoemCard } from '@/components/poem/PoemCard';

export function HallOfFamePage() {
  var [page, setPage] = useState(1);
  var { data, isLoading } = useHallOfFame({ page, pageSize: 12 });
  var totalPages = data ? Math.ceil(data.totalCount / data.pageSize) : 0;

  return (
    <div>
      <div className="mb-8 text-center">
        <h1 className="font-serif text-3xl font-bold text-ink">
          Hall of Fame
        </h1>
        <p className="mt-2 text-feather">
          The finest collaborative poems
        </p>
      </div>

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
      ) : (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3"
        >
          {data?.items.map((poem) => (
            <div
              key={poem.id}
              className="rounded-xl border-2 border-quill-gold/40"
            >
              <PoemCard poem={poem} />
            </div>
          ))}
          {data?.items.length === 0 && (
            <p className="col-span-full text-center text-feather">
              No poems in the Hall of Fame yet.
            </p>
          )}
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
