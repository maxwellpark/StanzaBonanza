import { Link } from 'react-router-dom';

export function NotFoundPage() {
  return (
    <div className="flex flex-col items-center py-24 text-center">
      <h1 className="mb-2 font-serif text-5xl font-bold text-ink">404</h1>
      <p className="mb-6 text-lg text-feather">Page not found</p>
      <Link to="/" className="btn-primary">
        Back to Home
      </Link>
    </div>
  );
}
