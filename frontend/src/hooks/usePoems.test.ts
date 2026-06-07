import { describe, it, expect, beforeEach } from 'vitest';
import { renderHook, waitFor } from '@testing-library/react';
import { http, HttpResponse } from 'msw';
import { server } from '@/test/server';
import { Wrapper } from '@/test/helpers';
import { usePoem, usePoems, useReviewStanza } from './usePoems';
import type { Poem, Stanza } from '@/types/poem';
import type { PaginatedResponse } from '@/types/api';

var mockPoem: Poem = {
  id: 'poem-1',
  authorId: 'user-1',
  title: 'Ode to a Nightingale',
  description: 'A classic ode',
  format: 'free_verse',
  approvalMode: 'approval_required',
  isHallOfFame: false,
  likeCount: 12,
  stanzaCount: 3,
  commentCount: 2,
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
};

var mockPoems: PaginatedResponse<Poem> = {
  items: [mockPoem],
  totalCount: 1,
  page: 1,
  pageSize: 20,
};

function apiJson<T>(data: T) {
  return HttpResponse.json({ data });
}

beforeEach(() => {
  server.use(
    http.get('/api/v1/poems', () => apiJson(mockPoems)),
    http.get('/api/v1/poems/:id', ({ params }) => {
      if (params.id === 'poem-1') {
        return apiJson(mockPoem);
      }
      return new HttpResponse(null, { status: 404 });
    }),
  );
});

describe('usePoems', () => {
  it('returns a list of poems', async () => {
    var { result } = renderHook(() => usePoems(), { wrapper: Wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data?.items).toHaveLength(1);
    expect(result.current.data?.items[0].title).toBe('Ode to a Nightingale');
  });

  it('exposes totalCount from the paginated response', async () => {
    var { result } = renderHook(() => usePoems(), { wrapper: Wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data?.totalCount).toBe(1);
  });

  it('builds query params when format is supplied', async () => {
    var capturedUrl = '';
    server.use(
      http.get('/api/v1/poems', ({ request }) => {
        capturedUrl = request.url;
        return apiJson(mockPoems);
      }),
    );

    var { result } = renderHook(() => usePoems({ format: 'haiku', page: 2 }), {
      wrapper: Wrapper,
    });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(capturedUrl).toContain('format=haiku');
    expect(capturedUrl).toContain('page=2');
  });
});

describe('usePoem', () => {
  it('returns the poem matching the supplied id', async () => {
    var { result } = renderHook(() => usePoem('poem-1'), { wrapper: Wrapper });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(result.current.data?.id).toBe('poem-1');
    expect(result.current.data?.title).toBe('Ode to a Nightingale');
  });

  it('does not fetch when id is empty', () => {
    var { result } = renderHook(() => usePoem(''), { wrapper: Wrapper });

    // enabled: !!id means query stays idle.
    expect(result.current.fetchStatus).toBe('idle');
  });
});

describe('useReviewStanza', () => {
  it('sends PUT with approved: true when approving a stanza', async () => {
    var capturedBody: unknown;
    server.use(
      http.put('/api/v1/poems/poem-1/stanzas/stanza-1', async ({ request }) => {
        capturedBody = await request.json();
        var stanza: Partial<Stanza> = { id: 'stanza-1', status: 'approved' };
        return apiJson(stanza);
      }),
    );

    var { result } = renderHook(() => useReviewStanza('poem-1'), { wrapper: Wrapper });

    result.current.mutate({ stanzaId: 'stanza-1', approved: true });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(capturedBody).toEqual({ approved: true });
  });

  it('sends PUT with approved: false when rejecting a stanza', async () => {
    var capturedBody: unknown;
    server.use(
      http.put('/api/v1/poems/poem-1/stanzas/stanza-2', async ({ request }) => {
        capturedBody = await request.json();
        var stanza: Partial<Stanza> = { id: 'stanza-2', status: 'rejected' };
        return apiJson(stanza);
      }),
    );

    var { result } = renderHook(() => useReviewStanza('poem-1'), { wrapper: Wrapper });

    result.current.mutate({ stanzaId: 'stanza-2', approved: false });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(capturedBody).toEqual({ approved: false });
  });

  it('hits the correct URL including poemId and stanzaId', async () => {
    var capturedUrl = '';
    server.use(
      http.put('/api/v1/poems/:poemId/stanzas/:stanzaId', ({ request, params }) => {
        capturedUrl = request.url;
        var stanza: Partial<Stanza> = { id: params.stanzaId as string };
        return apiJson(stanza);
      }),
    );

    var { result } = renderHook(() => useReviewStanza('poem-42'), { wrapper: Wrapper });

    result.current.mutate({ stanzaId: 'stanza-99', approved: true });

    await waitFor(() => expect(result.current.isSuccess).toBe(true));

    expect(capturedUrl).toContain('/api/v1/poems/poem-42/stanzas/stanza-99');
  });
});
