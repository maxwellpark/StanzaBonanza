import { Link } from 'react-router-dom';
import type { Stanza } from '@/types/poem';

interface StanzaBlockProps {
  stanza: Stanza;
  isLast?: boolean;
}

var deviceColors: Record<string, string> = {
  metaphor: 'bg-purple-100 text-purple-700',
  simile: 'bg-blue-100 text-blue-700',
  alliteration: 'bg-emerald-100 text-emerald-700',
  enjambment: 'bg-amber-100 text-amber-700',
  imagery: 'bg-rose-100 text-rose-700',
  personification: 'bg-cyan-100 text-cyan-700',
};

export function StanzaBlock({ stanza, isLast }: StanzaBlockProps) {
  return (
    <div className="py-4">
      <div className="poem-text">{stanza.text}</div>

      <div className="mt-3 flex flex-wrap items-center gap-2">
        {stanza.author && (
          <Link
            to={`/profile/${stanza.authorId}`}
            className="font-sans text-sm text-feather no-underline hover:text-accent"
          >
            &mdash; @{stanza.author.displayName}
          </Link>
        )}

        {stanza.literaryDevice && (
          <span
            className={`inline-block rounded-full px-2 py-0.5 font-sans text-xs font-medium ${
              deviceColors[stanza.literaryDevice] ?? 'bg-gray-100 text-gray-700'
            }`}
          >
            {stanza.literaryDevice}
          </span>
        )}

        {stanza.status === 'pending' && (
          <span className="inline-block rounded-full bg-warning/15 px-2 py-0.5 font-sans text-xs font-medium text-warning">
            Awaiting approval
          </span>
        )}
      </div>

      {!isLast && (
        <div className="mt-6 text-center text-feather/40 select-none">&#10022;</div>
      )}
    </div>
  );
}
