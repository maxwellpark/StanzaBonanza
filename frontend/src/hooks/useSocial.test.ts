import { describe, it, expect, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { http, HttpResponse } from 'msw';
import { server } from '@/test/server';
import { Wrapper } from '@/test/helpers';
import { useNotifications, useToggleFollow } from './useSocial';
import type { Notification } from '@/types/social';
import type { PaginatedResponse } from '@/types/api';

function apiJson<T>(data: T) {
  return HttpResponse.json({ data });
}

var mockNotifications: PaginatedResponse<Notification> = {
  items: [
    {
      id: 'notif-1',
      recipientId: 'user-1',
      actorId: 'user-2',
      actor: {
        id: 'user-2',
        displayName: 'Byron',
        email: 'byron@example.com',
        bio: '',
        avatarUrl: '',
        isVerified: false,
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z',
      },
      type: 'like',
      poemId: 'poem-1',
      read: false,
      createdAt: '2024-05-01T10:00:00Z',
    },
    {
      id: 'notif-2',
      recipientId: 'user-1',
      actorId: 'user-3',
      type: 'follow',
      read: true,
      createdAt: '2024-04-28T08:00:00Z',
    },
  ],
  totalCount: 2,
  page: 1,
  pageSize: 20,
};

beforeEach(() => {
  server.use(
    http.get('/api/v1/notifications', () => apiJson(mockNotifications)),
    http.post('/api/v1/users/:userId/follow', ({ params }) => {
      return apiJson({ following: true, userId: params.userId });
    }),
  );
});

describe('useNotifications', () => {
  it('returns paginated notification items', async () => {
    var { result } = renderHook(() => useNotifications(), { wrapper: Wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data?.items).toHaveLength(2);
    expect(result.current.data?.totalCount).toBe(2);
  });

  it('includes the notification type and read flag', async () => {
    var { result } = renderHook(() => useNotifications(), { wrapper: Wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    var first = result.current.data!.items[0];
    expect(first.type).toBe('like');
    expect(first.read).toBe(false);

    var second = result.current.data!.items[1];
    expect(second.type).toBe('follow');
    expect(second.read).toBe(true);
  });

  it('passes the page param in the query string', async () => {
    var capturedUrl = '';
    server.use(
      http.get('/api/v1/notifications', ({ request }) => {
        capturedUrl = request.url;
        return apiJson(mockNotifications);
      }),
    );

    var { result } = renderHook(() => useNotifications({ page: 3 }), { wrapper: Wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(capturedUrl).toContain('page=3');
  });
});

describe('useToggleFollow', () => {
  it('calls POST /api/v1/users/:userId/follow', async () => {
    var capturedUrl = '';
    server.use(
      http.post('/api/v1/users/:userId/follow', ({ request }) => {
        capturedUrl = request.url;
        return apiJson({ following: true });
      }),
    );

    var { result } = renderHook(() => useToggleFollow('user-99'), { wrapper: Wrapper });

    result.current.mutate();

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(capturedUrl).toContain('/api/v1/users/user-99/follow');
  });

  it('returns the following status from the response', async () => {
    server.use(
      http.post('/api/v1/users/:userId/follow', () => apiJson({ following: false })),
    );

    var { result } = renderHook(() => useToggleFollow('user-5'), { wrapper: Wrapper });

    result.current.mutate();

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data).toEqual({ following: false });
  });
});
