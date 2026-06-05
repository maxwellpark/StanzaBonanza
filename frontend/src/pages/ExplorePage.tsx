import { useState } from 'react';
import { motion } from 'framer-motion';
import { usePoems } from '@/hooks/usePoems';
import { PoemCard } from '@/components/poem/PoemCard';
import { formatPoemFormat } from '@/lib/utils';
import type { PoemFormat } from '@/types/poem';

var formats: (PoemFormat | 'all')[] = [
  'all',
  'free_verse',
  'haiku',
  'sonnet',
  'limerick',
  'iambic_pentameter',
  'rhyming_couplets',
];

type SortOption = 'recent' | 'popular';

export function ExplorePage() {
  var [page, setPage] = useState(1);
  var [format, setFormat] = useState<PoemFormat | 'all'>('all');
  var [sort, setSort] = useState<SortOption>('recent');

  var { data, isLoading } = usePoems({
    page,
    pageSize: 12,
    format: format === 'all' ? undefined : format,
    sort,
  });

  var totalPages = data ? Math.ceil(data.totalCount / data.pageSize) : 0;

  return (
    <div>
      <h1 className="mb-6 font-serif text-3xl font-bold text-ink">Explore</h1>

      <div className="mb-6 flex flex-wrap gap-2">
        {formats.map((f) => (
          <button
            key={f}
            onClick={() => {
              setFormat(f);
              setPage(1);
            }}
            className={`rounded-full px-4 py-1.5 font-sans text-sm transition-colors ${
              format === f
                ? 'bg-accent text-white'
                : 'bg-white text-feather hover:bg-parchment-dark'
            }`}
          >
            {f === 'all' ? 'All' : formatPoemFormat(f)}
          </button>
        ))}
      </div>

      <div className="mb-6 flex gap-2">
        {(['recent', 'popular'] as SortOption[]).map((s) => (
          <button
            key={s}
            onClick={() => {
              setSort(s);
              setPage(1);
            }}
            className={`rounded-lg px-3 py-1 font-sans text-sm transition-colors ${
              sort === s
                ? 'bg-ink text-parchment'
                : 'bg-white text-feather hover:bg-parchment-dark'
            }`}
          >
            {s.charAt(0).toUpperCase() + s.slice(1)}
          </button>
        ))}
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
            <PoemCard key={poem.id} poem={poem} />
          ))}
          {data?.items.length === 0 && (
            <p className="col-span-full text-center text-feather">
              No poems found. Try a different filter.
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
