import { Link, useNavigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import type { Poem } from '@/types/poem';
import { PoemFormatBadge } from './PoemFormatBadge';

interface PoemCardProps {
  poem: Poem;
}

export function PoemCard({ poem }: PoemCardProps) {
  var navigate = useNavigate();

  var previewText = poem.stanzas?.[0]?.text ?? '';
  var previewLines = previewText.split('\n').slice(0, 4);
  var isTruncated = previewText.split('\n').length > 4;

  return (
    <motion.div
      whileHover={{ scale: 1.02 }}
      transition={{ duration: 0.2 }}
      onClick={() => navigate(`/poems/${poem.id}`)}
      className="card cursor-pointer transition-shadow hover:shadow-md"
    >
      <div className="mb-3 flex items-start justify-between gap-2">
        <Link
          to={`/poems/${poem.id}`}
          onClick={(e) => e.stopPropagation()}
          className="font-serif text-lg font-bold text-ink no-underline hover:text-accent"
        >
          {poem.title}
        </Link>
        <PoemFormatBadge format={poem.format} />
      </div>

      {poem.author && (
        <div className="mb-3 flex items-center gap-2">
          <img
            src={poem.author.avatarUrl || '/default-avatar.png'}
            alt=""
            className="h-6 w-6 rounded-full object-cover"
          />
          <span className="font-sans text-sm text-feather">{poem.author.displayName}</span>
        </div>
      )}

      {previewLines.length > 0 && (
        <div className="mb-4 text-base leading-relaxed text-ink-light" style={{ fontFamily: 'var(--font-body)' }}>
          {previewLines.map((line, i) => (
            <p key={i} className="my-0">{line}</p>
          ))}
          {isTruncated && <p className="my-0 text-feather">...</p>}
        </div>
      )}

      <div className="flex items-center gap-4 font-sans text-xs text-feather">
        <span className="flex items-center gap-1">
          <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
          </svg>
          {poem.likeCount}
        </span>
        <span className="flex items-center gap-1">
          <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
          {poem.stanzaCount}
        </span>
        <span className="flex items-center gap-1">
          <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
          {poem.commentCount}
        </span>
      </div>
    </motion.div>
  );
}
