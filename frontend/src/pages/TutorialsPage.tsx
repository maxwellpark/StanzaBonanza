import { Link } from 'react-router-dom';
import { motion } from 'framer-motion';
import { useTutorials } from '@/hooks/useTutorials';
import { PoemFormatBadge } from '@/components/poem/PoemFormatBadge';
import type { Difficulty } from '@/types/tutorial';

var difficultyColors: Record<Difficulty, string> = {
  beginner: 'bg-green-100 text-green-700',
  intermediate: 'bg-yellow-100 text-yellow-700',
  advanced: 'bg-red-100 text-red-700',
};

export function TutorialsPage() {
  var { data: tutorials, isLoading } = useTutorials();

  return (
    <div>
      <div className="mb-8">
        <h1 className="font-serif text-3xl font-bold text-ink">Poetry Tutorials</h1>
        <p className="mt-2 font-sans text-base text-feather">
          Learn the art of different poetic forms, from haiku to sonnets.
        </p>
      </div>

      {isLoading ? (
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3">
          {Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="card animate-pulse space-y-3">
              <div className="h-5 w-3/4 rounded bg-parchment-dark" />
              <div className="h-4 w-full rounded bg-parchment-dark" />
              <div className="h-4 w-2/3 rounded bg-parchment-dark" />
            </div>
          ))}
        </div>
      ) : !tutorials?.length ? (
        <div className="py-16 text-center">
          <p className="font-sans text-feather">No tutorials yet. Check back soon.</p>
        </div>
      ) : (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-3"
        >
          {tutorials.map((tutorial) => (
            <motion.div
              key={tutorial.id}
              whileHover={{ scale: 1.02 }}
              transition={{ duration: 0.15 }}
            >
              <Link
                to={`/tutorials/${tutorial.slug}`}
                className="card block no-underline transition-shadow hover:shadow-md"
              >
                <div className="mb-3 flex flex-wrap items-start gap-2">
                  <h2 className="flex-1 font-serif text-lg font-bold text-ink">
                    {tutorial.title}
                  </h2>
                  <PoemFormatBadge format={tutorial.format} />
                </div>
                <div className="flex items-center gap-2">
                  <span className={`rounded-full px-2 py-0.5 font-sans text-xs font-medium capitalize ${difficultyColors[tutorial.difficulty]}`}>
                    {tutorial.difficulty}
                  </span>
                </div>
              </Link>
            </motion.div>
          ))}
        </motion.div>
      )}
    </div>
  );
}
