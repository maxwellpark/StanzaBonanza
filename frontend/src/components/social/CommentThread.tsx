import { useState } from 'react';
import { Link } from 'react-router-dom';
import { useComments, useAddComment } from '@/hooks/useSocial';
import { useAuthStore } from '@/stores/authStore';
import { useUIStore } from '@/stores/uiStore';
import { timeAgo } from '@/lib/utils';

interface CommentThreadProps {
  poemId: string;
}

export function CommentThread({ poemId }: CommentThreadProps) {
  var [text, setText] = useState('');
  var { data, isLoading } = useComments(poemId);
  var addComment = useAddComment(poemId);
  var { isAuthenticated } = useAuthStore();
  var { openLogin } = useUIStore();

  var comments = data?.items ?? [];

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!text.trim()) {
      return;
    }

    if (!isAuthenticated) {
      openLogin();
      return;
    }

    try {
      await addComment.mutateAsync({ text: text.trim() });
      setText('');
    } catch {
      alert('Failed to post comment. Please try again.');
    }
  }

  if (isLoading) {
    return <p className="font-sans text-sm text-feather">Loading comments...</p>;
  }

  return (
    <div>
      {comments.length === 0 && (
        <p className="mb-4 font-sans text-sm text-feather">No comments yet. Start the conversation!</p>
      )}

      <div className="flex flex-col gap-4">
        {comments.map((comment) => (
          <div key={comment.id} className="rounded-lg border border-parchment-dark bg-white px-4 py-3">
            <div className="mb-1 flex items-center gap-2">
              {comment.author && (
                <Link
                  to={`/profile/${comment.authorId}`}
                  className="flex items-center gap-2 no-underline"
                >
                  <img
                    src={comment.author.avatarUrl || '/default-avatar.png'}
                    alt=""
                    className="h-6 w-6 rounded-full object-cover"
                  />
                  <span className="font-sans text-sm font-medium text-ink hover:text-accent">
                    {comment.author.displayName}
                  </span>
                </Link>
              )}
              <span className="font-sans text-xs text-feather">{timeAgo(comment.createdAt)}</span>
            </div>
            <p className="font-sans text-sm leading-relaxed text-ink-light">{comment.text}</p>
          </div>
        ))}
      </div>

      <form onSubmit={handleSubmit} className="mt-4 flex gap-2">
        <textarea
          value={text}
          onChange={(e) => setText(e.target.value)}
          placeholder={isAuthenticated ? 'Write a comment...' : 'Sign in to comment'}
          rows={2}
          className="flex-1 resize-none rounded-lg border border-parchment-dark bg-white px-4 py-2 font-sans text-sm text-ink outline-none transition-colors focus:border-accent"
        />
        <button
          type="submit"
          disabled={addComment.isPending || !text.trim()}
          className="btn-primary self-end disabled:opacity-50"
        >
          {addComment.isPending ? '...' : 'Post'}
        </button>
      </form>
    </div>
  );
}
