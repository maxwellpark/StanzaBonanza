import { useParams, Link } from 'react-router-dom';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import { useTutorial } from '@/hooks/useTutorials';
import { PoemFormatBadge } from '@/components/poem/PoemFormatBadge';
import { formatDate } from '@/lib/utils';

export function TutorialDetailPage() {
  var { slug } = useParams<{ slug: string }>();
  var { data: tutorial, isLoading, isError } = useTutorial(slug ?? '');

  if (isLoading) {
    return (
      <div className="mx-auto max-w-2xl animate-pulse space-y-4 py-8">
        <div className="h-8 w-2/3 rounded bg-parchment-dark" />
        <div className="space-y-3">
          {Array.from({ length: 8 }).map((_, i) => (
            <div key={i} className="h-4 rounded bg-parchment-dark" style={{ width: `${70 + Math.random() * 30}%` }} />
          ))}
        </div>
      </div>
    );
  }

  if (isError || !tutorial) {
    return (
      <div className="py-16 text-center">
        <p className="mb-4 font-sans text-feather">Tutorial not found.</p>
        <Link to="/tutorials" className="btn-primary">Back to Tutorials</Link>
      </div>
    );
  }

  return (
    <article className="mx-auto max-w-2xl">
      <Link
        to="/tutorials"
        className="mb-6 flex items-center gap-1 font-sans text-sm text-feather no-underline hover:text-ink"
      >
        <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
        </svg>
        All Tutorials
      </Link>

      <header className="mb-8">
        <div className="mb-3 flex flex-wrap items-center gap-3">
          <PoemFormatBadge format={tutorial.format} />
          <span className="font-sans text-sm capitalize text-feather">{tutorial.difficulty}</span>
          <span className="font-sans text-sm text-feather">{formatDate(tutorial.createdAt)}</span>
        </div>
        <h1 className="font-serif text-3xl font-bold text-ink md:text-4xl">{tutorial.title}</h1>
      </header>

      <div className="prose prose-ink max-w-none">
        <ReactMarkdown remarkPlugins={[remarkGfm]}>{tutorial.contentMd}</ReactMarkdown>
      </div>

      <div className="mt-12 border-t border-parchment-dark pt-8">
        <p className="mb-4 font-serif text-lg text-ink">Ready to try it?</p>
        <Link to="/poems/new" className="btn-primary">
          Write a {tutorial.title} Poem
        </Link>
      </div>
    </article>
  );
}
