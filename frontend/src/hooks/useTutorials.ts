import { useQuery } from '@tanstack/react-query';
import { api } from '@/lib/api';
import type { Tutorial } from '@/types/tutorial';

export function useTutorials() {
  return useQuery({
    queryKey: ['tutorials'],
    queryFn: () => api.get<Tutorial[]>('/tutorials'),
  });
}

export function useTutorial(slug: string) {
  return useQuery({
    queryKey: ['tutorial', slug],
    queryFn: () => api.get<Tutorial>(`/tutorials/${slug}`),
    enabled: !!slug,
  });
}
