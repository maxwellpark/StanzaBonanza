import { Navigate } from 'react-router-dom';
import { useAuthStore } from '@/stores/authStore';
import { PoemEditor } from '@/components/poem/PoemEditor';

export function CreatePoemPage() {
  var { isAuthenticated, isLoading } = useAuthStore();

  if (isLoading) {
    return null;
  }

  if (!isAuthenticated) {
    return <Navigate to="/" replace />;
  }

  return (
    <div className="mx-auto max-w-2xl">
      <PoemEditor />
    </div>
  );
}
