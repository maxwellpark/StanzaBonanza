import { Link, useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import { useExplore } from '@/hooks/usePoems';
import { useAuthStore } from '@/stores/authStore';
import { useUIStore } from '@/stores/uiStore';
import { PoemCard } from '@/components/poem/PoemCard';

var container = {
  hidden: { opacity: 0 },
  show: {
    opacity: 1,
    transition: { staggerChildren: 0.1 },
  },
};

var item = {
  hidden: { opacity: 0, y: 20 },
  show: { opacity: 1, y: 0 },
};

export function HomePage() {
  var { isAuthenticated } = useAuthStore();
  var { openLogin } = useUIStore();
  var navigate = useNavigate();
  var { data, isLoading } = useExplore({ page: 1, pageSize: 6 });

  function handleStartWriting() {
    if (isAuthenticated) {
      navigate('/poems/new');
    } else {
      openLogin();
    }
  }

  return (
    <motion.div
      variants={container}
      initial="hidden"
      animate="show"
      className="flex flex-col items-center"
    >
      <motion.section variants={item} className="py-16 text-center">
        <h1 className="font-serif text-5xl font-bold text-ink md:text-6xl">
          Stanza Bonanza
        </h1>
        <p className="mt-4 font-body text-xl text-feather">
          Where poems grow, stanza by stanza
        </p>
        <div className="mt-8 flex flex-wrap justify-center gap-4">
          <button onClick={handleStartWriting} className="btn-primary text-lg">
            Start Writing
          </button>
          <Link to="/explore" className="btn-secondary text-lg">
            Explore Poems
          </Link>
        </div>
      </motion.section>

      <motion.section variants={item} className="w-full">
        <h2 className="mb-6 font-serif text-2xl font-bold text-ink">
          Recent Poems
        </h2>
        {isLoading ? (
          <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
            {Array.from({ length: 6 }).map((_, i) => (
              <div
                key={i}
                className="card animate-pulse"
              >
                <div className="mb-3 h-6 w-3/4 rounded bg-parchment-dark" />
                <div className="mb-2 h-4 w-full rounded bg-parchment-dark" />
                <div className="h-4 w-2/3 rounded bg-parchment-dark" />
              </div>
            ))}
          </div>
        ) : (
          <div className="grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-3">
            {data?.items.map((poem) => (
              <motion.div key={poem.id} variants={item}>
                <PoemCard poem={poem} />
              </motion.div>
            ))}
          </div>
        )}
      </motion.section>
    </motion.div>
  );
}
