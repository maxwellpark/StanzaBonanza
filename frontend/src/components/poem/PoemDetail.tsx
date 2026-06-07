import { useState } from 'react';
import { Link } from 'react-router-dom';
import type { Poem } from '@/types/poem';
import { useAuthStore } from '@/stores/authStore';
import { useUIStore } from '@/stores/uiStore';
import { formatDate } from '@/lib/utils';
import { PoemFormatBadge } from './PoemFormatBadge';
import { StanzaBlock } from './StanzaBlock';
import { ExtendPoemDialog } from './ExtendPoemDialog';
import { LikeButton } from '@/components/social/LikeButton';
import { CommentThread } from '@/components/social/CommentThread';
import { StanzaReviewPanel } from './StanzaReviewPanel';

interface PoemDetailProps {
  poem: Poem;
}

export function PoemDetail({ poem }: PoemDetailProps) {
  var [extendOpen, setExtendOpen] = useState(false);
  var { isAuthenticated, user } = useAuthStore();
  var { openLogin } = useUIStore();

  var stanzas = poem.stanzas ?? [];
  var pendingStanzas = stanzas.filter((s) => s.status === 'pending');
  var approvedStanzas = stanzas.filter((s) => s.status !== 'pending');
  var isAuthor = isAuthenticated && user?.id === poem.authorId;
  var canExtend =
    poem.approvalMode !== 'closed' &&
    (!poem.maxStanzas || stanzas.length < poem.maxStanzas);

  function handleAddStanza() {
    if (!isAuthenticated) {
      openLogin();
      return;
    }
    setExtendOpen(true);
  }

  return (
    <article className="mx-auto max-w-2xl">
      <header className="mb-8">
        <h1 className="mb-3 font-serif text-3xl font-bold text-ink md:text-4xl">
          {poem.title}
        </h1>

        <div className="mb-4 flex flex-wrap items-center gap-3">
          {poem.author && (
            <Link
              to={`/profile/${poem.authorId}`}
              className="flex items-center gap-2 no-underline"
            >
              <img
                src={poem.author.avatarUrl || '/default-avatar.png'}
                alt=""
                className="h-8 w-8 rounded-full object-cover"
              />
              <span className="font-sans text-sm text-feather hover:text-accent">
                {poem.author.displayName}
              </span>
            </Link>
          )}
          <PoemFormatBadge format={poem.format} />
          <span className="font-sans text-sm text-feather">{formatDate(poem.createdAt)}</span>
        </div>

        {poem.description && (
          <p className="text-base leading-relaxed text-feather">{poem.description}</p>
        )}
      </header>

      {isAuthor && (
        <StanzaReviewPanel poemId={poem.id} pendingStanzas={pendingStanzas} />
      )}

      <div className="card mb-8">
        {approvedStanzas.map((stanza, i) => (
          <StanzaBlock key={stanza.id} stanza={stanza} isLast={i === approvedStanzas.length - 1} />
        ))}

        {approvedStanzas.length === 0 && (
          <p className="py-8 text-center font-sans text-sm text-feather">
            This poem has no stanzas yet. Be the first to contribute!
          </p>
        )}
      </div>

      <div className="mb-8 flex items-center gap-4">
        <LikeButton poemId={poem.id} likeCount={poem.likeCount} />

        <span className="flex items-center gap-1 font-sans text-sm text-feather">
          <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
          {poem.commentCount}
        </span>

        <button
          className="font-sans text-sm text-feather transition-colors hover:text-ink"
          onClick={() => {
            navigator.clipboard.writeText(window.location.href);
            alert('Link copied!');
          }}
        >
          <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
            <path strokeLinecap="round" strokeLinejoin="round" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z" />
          </svg>
        </button>
      </div>

      {canExtend && (
        <div className="mb-8">
          <button onClick={handleAddStanza} className="btn-primary">
            Add a Stanza
          </button>
        </div>
      )}

      <ExtendPoemDialog
        poemId={poem.id}
        format={poem.format}
        isOpen={extendOpen}
        onClose={() => setExtendOpen(false)}
      />

      <section className="mb-12">
        <h2 className="mb-4 font-serif text-xl font-bold text-ink">Comments</h2>
        <CommentThread poemId={poem.id} />
      </section>
    </article>
  );
}
