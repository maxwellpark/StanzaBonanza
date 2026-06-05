import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { useToggleLike } from '@/hooks/useSocial';
import { useAuthStore } from '@/stores/authStore';
import { useUIStore } from '@/stores/uiStore';

interface LikeButtonProps {
  poemId: string;
  likeCount: number;
  isLiked?: boolean;
}

export function LikeButton({ poemId, likeCount, isLiked: initialLiked }: LikeButtonProps) {
  var [optimisticLiked, setOptimisticLiked] = useState(initialLiked ?? false);
  var [optimisticCount, setOptimisticCount] = useState(likeCount);
  var [animKey, setAnimKey] = useState(0);
  var mutation = useToggleLike(poemId);
  var { isAuthenticated } = useAuthStore();
  var { openLogin } = useUIStore();

  function handleClick() {
    if (!isAuthenticated) {
      openLogin();
      return;
    }

    var nextLiked = !optimisticLiked;
    setOptimisticLiked(nextLiked);
    setOptimisticCount((c) => c + (nextLiked ? 1 : -1));

    if (nextLiked) {
      setAnimKey((k) => k + 1);
    }

    mutation.mutate(undefined, {
      onError: () => {
        setOptimisticLiked(!nextLiked);
        setOptimisticCount((c) => c + (nextLiked ? -1 : 1));
      },
    });
  }

  return (
    <button
      onClick={handleClick}
      className="relative flex items-center gap-1 font-sans text-sm transition-colors"
    >
      <span className="relative">
        <svg
          className={`h-5 w-5 transition-colors ${
            optimisticLiked ? 'fill-error text-error' : 'fill-none text-feather hover:text-ink'
          }`}
          viewBox="0 0 24 24"
          stroke="currentColor"
          strokeWidth={2}
        >
          <path strokeLinecap="round" strokeLinejoin="round" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
        </svg>

        <AnimatePresence>
          {optimisticLiked && (
            <motion.svg
              key={animKey}
              initial={{ opacity: 1, y: 0, scale: 1 }}
              animate={{ opacity: 0, y: -16, scale: 0.6 }}
              exit={{ opacity: 0 }}
              transition={{ duration: 0.5, ease: 'easeOut' }}
              className="absolute inset-0 h-5 w-5 fill-error text-error"
              viewBox="0 0 24 24"
              stroke="currentColor"
              strokeWidth={2}
            >
              <path strokeLinecap="round" strokeLinejoin="round" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
            </motion.svg>
          )}
        </AnimatePresence>
      </span>

      <span className={optimisticLiked ? 'text-error' : 'text-feather'}>
        {optimisticCount}
      </span>
    </button>
  );
}
