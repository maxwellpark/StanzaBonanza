import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { MemoryRouter } from 'react-router-dom';
import { QueryClient, QueryClientProvider } from '@tanstack/react-query';
import { NotificationsPanel } from './NotificationsPanel';
import type { Notification } from '@/types/social';
import type { PaginatedResponse } from '@/types/api';

var mockMarkRead = vi.fn();

vi.mock('@/hooks/useSocial', () => ({
  useNotifications: vi.fn(),
  useMarkNotificationsRead: vi.fn(() => ({ mutate: mockMarkRead })),
}));

import { useNotifications } from '@/hooks/useSocial';

function makeNotification(overrides: Partial<Notification> = {}): Notification {
  return {
    id: 'notif-1',
    recipientId: 'user-1',
    actorId: 'user-2',
    actor: {
      id: 'user-2',
      displayName: 'Wordsworth',
      email: 'w@example.com',
      bio: '',
      avatarUrl: '',
      isVerified: false,
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-01T00:00:00Z',
    },
    type: 'like',
    poemId: 'poem-1',
    poem: {
      id: 'poem-1',
      authorId: 'user-1',
      title: 'The Prelude',
      description: '',
      format: 'free_verse',
      approvalMode: 'open',
      isHallOfFame: false,
      likeCount: 5,
      stanzaCount: 2,
      commentCount: 0,
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-01T00:00:00Z',
    },
    read: false,
    createdAt: new Date(Date.now() - 60 * 1000).toISOString(),
    ...overrides,
  };
}

function makeResponse(items: Notification[]): PaginatedResponse<Notification> {
  return { items, totalCount: items.length, page: 1, pageSize: 20 };
}

function renderPanel() {
  var client = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  return render(
    <MemoryRouter>
      <QueryClientProvider client={client}>
        <NotificationsPanel />
      </QueryClientProvider>
    </MemoryRouter>,
  );
}

beforeEach(() => {
  vi.clearAllMocks();
});

describe('NotificationsPanel', () => {
  it('renders the bell button', () => {
    vi.mocked(useNotifications).mockReturnValue({
      data: makeResponse([]),
      isLoading: false,
    } as ReturnType<typeof useNotifications>);

    renderPanel();

    expect(screen.getByRole('button', { name: /notifications/i })).toBeInTheDocument();
  });

  it('shows no badge when there are no unread notifications', () => {
    var read = makeNotification({ read: true });
    vi.mocked(useNotifications).mockReturnValue({
      data: makeResponse([read]),
      isLoading: false,
    } as ReturnType<typeof useNotifications>);

    renderPanel();

    // Badge span should not exist.
    expect(screen.queryByText('1')).not.toBeInTheDocument();
  });

  it('shows an unread badge with the correct count', () => {
    var unread1 = makeNotification({ id: 'n-1', read: false });
    var unread2 = makeNotification({ id: 'n-2', read: false });
    vi.mocked(useNotifications).mockReturnValue({
      data: makeResponse([unread1, unread2]),
      isLoading: false,
    } as ReturnType<typeof useNotifications>);

    renderPanel();

    expect(screen.getByText('2')).toBeInTheDocument();
  });

  it('caps the badge at "9+" when there are more than 9 unread', () => {
    var unread = Array.from({ length: 10 }, (_, i) =>
      makeNotification({ id: `n-${i}`, read: false }),
    );
    vi.mocked(useNotifications).mockReturnValue({
      data: makeResponse(unread),
      isLoading: false,
    } as ReturnType<typeof useNotifications>);

    renderPanel();

    expect(screen.getByText('9+')).toBeInTheDocument();
  });

  it('opens the dropdown panel when the bell is clicked', async () => {
    var user = userEvent.setup();
    var notif = makeNotification();
    vi.mocked(useNotifications).mockReturnValue({
      data: makeResponse([notif]),
      isLoading: false,
    } as ReturnType<typeof useNotifications>);

    renderPanel();

    await user.click(screen.getByRole('button', { name: /notifications/i }));

    expect(screen.getByRole('heading', { name: /notifications/i })).toBeInTheDocument();
  });

  it('shows notification text in the dropdown', async () => {
    var user = userEvent.setup();
    var notif = makeNotification();
    vi.mocked(useNotifications).mockReturnValue({
      data: makeResponse([notif]),
      isLoading: false,
    } as ReturnType<typeof useNotifications>);

    renderPanel();

    await user.click(screen.getByRole('button', { name: /notifications/i }));

    // notifLabel for a 'like' type: "<actor> liked "<poem title>""
    expect(screen.getByText(/Wordsworth liked/)).toBeInTheDocument();
  });

  it('shows "No notifications yet." when the list is empty', async () => {
    var user = userEvent.setup();
    vi.mocked(useNotifications).mockReturnValue({
      data: makeResponse([]),
      isLoading: false,
    } as ReturnType<typeof useNotifications>);

    renderPanel();

    await user.click(screen.getByRole('button', { name: /notifications/i }));

    expect(screen.getByText('No notifications yet.')).toBeInTheDocument();
  });

  it('shows an unread indicator inside the dropdown for unread notifications', async () => {
    var user = userEvent.setup();
    var notif = makeNotification({ read: false });
    vi.mocked(useNotifications).mockReturnValue({
      data: makeResponse([notif]),
      isLoading: false,
    } as ReturnType<typeof useNotifications>);

    renderPanel();

    await user.click(screen.getByRole('button', { name: /notifications/i }));

    // The "X unread" label appears in the dropdown header.
    expect(screen.getByText(/1 unread/i)).toBeInTheDocument();
  });

  it('calls markRead with unread notification ids when the panel opens', async () => {
    var user = userEvent.setup();
    var notif = makeNotification({ id: 'notif-unread', read: false });
    vi.mocked(useNotifications).mockReturnValue({
      data: makeResponse([notif]),
      isLoading: false,
    } as ReturnType<typeof useNotifications>);

    renderPanel();

    await user.click(screen.getByRole('button', { name: /notifications/i }));

    await waitFor(() => {
      expect(mockMarkRead).toHaveBeenCalledWith(['notif-unread']);
    });
  });

  it('does not call markRead when all notifications are already read', async () => {
    var user = userEvent.setup();
    var notif = makeNotification({ read: true });
    vi.mocked(useNotifications).mockReturnValue({
      data: makeResponse([notif]),
      isLoading: false,
    } as ReturnType<typeof useNotifications>);

    renderPanel();

    await user.click(screen.getByRole('button', { name: /notifications/i }));

    expect(mockMarkRead).not.toHaveBeenCalled();
  });
});
