import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { api } from '@/lib/api';
import type { Comment, Notification } from '@/types/social';
import type { User } from '@/types/user';
import type { PaginatedResponse, PaginationParams } from '@/types/api';

export function useToggleLike(poemId: string) {
  var queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => api.post<{ liked: boolean }>(`/poems/${poemId}/like`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['poem', poemId] });
      queryClient.invalidateQueries({ queryKey: ['poems'] });
    },
  });
}

export function useComments(poemId: string, params?: PaginationParams) {
  return useQuery({
    queryKey: ['comments', poemId, params],
    queryFn: () => api.get<PaginatedResponse<Comment>>(`/poems/${poemId}/comments?page=${params?.page ?? 1}`),
    enabled: !!poemId,
  });
}

export function useAddComment(poemId: string) {
  var queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: { text: string; parentId?: string }) =>
      api.post<Comment>(`/poems/${poemId}/comments`, input),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comments', poemId] });
      queryClient.invalidateQueries({ queryKey: ['poem', poemId] });
    },
  });
}

export function useToggleFollow(userId: string) {
  var queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => api.post<{ following: boolean }>(`/users/${userId}/follow`),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['user', userId] });
    },
  });
}

export function useFollowers(userId: string, params?: PaginationParams) {
  return useQuery({
    queryKey: ['followers', userId, params],
    queryFn: () => api.get<PaginatedResponse<User>>(`/users/${userId}/followers?page=${params?.page ?? 1}`),
    enabled: !!userId,
  });
}

export function useNotifications(params?: PaginationParams) {
  return useQuery({
    queryKey: ['notifications', params],
    queryFn: () => api.get<PaginatedResponse<Notification>>(`/notifications?page=${params?.page ?? 1}`),
    refetchInterval: 30000,
  });
}

export function useMarkNotificationsRead() {
  var queryClient = useQueryClient();
  return useMutation({
    mutationFn: (ids: string[]) => api.post('/notifications/read', { ids }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['notifications'] });
    },
  });
}
