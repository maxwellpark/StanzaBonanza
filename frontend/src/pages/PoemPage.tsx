import { useParams } from 'react-router-dom';
import { motion } from 'framer-motion';
import { usePoem } from '@/hooks/usePoems';
import { PoemDetail } from '@/components/poem/PoemDetail';
import { Link } from 'react-router-dom';

export function PoemPage() {
  var { poemId } = useParams<{ poemId: string }>();
  var { data: poem, isLoading, isError } = usePoem(poemId ?? '');

  if (isLoading) {
    return (
      <div className="mx-auto max-w-2xl">
        <div className="card animate-pulse">
          <div className="mb-4 h-8 w-1/2 rounded bg-parchment-dark" />
          <div className="mb-3 h-4 w-3/4 rounded bg-parchment-dark" />
          <div className="space-y-2">
            {Array.from({ length: 4 }).map((_, i) => (
              <div key={i} className="h-4 w-full rounded bg-parchment-dark" />
            ))}
          </div>
        </div>
      </div>
    );
  }

  if (isError || !poem) {
    return (
      <div className="flex flex-col items-center py-16 text-center">
        <h2 className="mb-2 font-serif text-2xl font-bold text-ink">
          Poem not found
        </h2>
        <p className="mb-4 text-feather">
          This poem may have been removed or doesn't exist.
        </p>
        <Link to="/explore" className="btn-primary">
          Explore Poems
        </Link>
      </div>
    );
  }

  return (
    <motion.div
      initial={{ opacity: 0, y: 10 }}
      animate={{ opacity: 1, y: 0 }}
      className="mx-auto max-w-2xl"
    >
      <PoemDetail poem={poem} />
    </motion.div>
  );
}
