import { describe, it, expect, vi, beforeEach } from 'vitest';
import { useAuthStore } from './authStore';
import type { User } from '@/types/user';

// Mock the api module before store imports resolve.
vi.mock('@/lib/api', () => ({
  api: {
    get: vi.fn(),
    post: vi.fn(),
  },
}));

import { api } from '@/lib/api';

var mockUser: User = {
  id: 'user-1',
  displayName: 'Ada Lovelace',
  email: 'ada@example.com',
  bio: 'Poet and mathematician',
  avatarUrl: 'https://example.com/ada.jpg',
  isVerified: true,
  createdAt: '2024-01-01T00:00:00Z',
  updatedAt: '2024-01-01T00:00:00Z',
};

beforeEach(() => {
  // Reset store state between tests.
  useAuthStore.setState({ user: null, isLoading: true, isAuthenticated: false });
  vi.clearAllMocks();
});

describe('fetchUser', () => {
  it('sets user and isAuthenticated on success', async () => {
    vi.mocked(api.get).mockResolvedValueOnce(mockUser);

    await useAuthStore.getState().fetchUser();

    var state = useAuthStore.getState();
    expect(state.user).toEqual(mockUser);
    expect(state.isAuthenticated).toBe(true);
    expect(state.isLoading).toBe(false);
  });

  it('clears user and sets isAuthenticated false on failure', async () => {
    vi.mocked(api.get).mockRejectedValueOnce(new Error('Unauthorized'));

    await useAuthStore.getState().fetchUser();

    var state = useAuthStore.getState();
    expect(state.user).toBeNull();
    expect(state.isAuthenticated).toBe(false);
    expect(state.isLoading).toBe(false);
  });

  it('calls the /auth/me endpoint', async () => {
    vi.mocked(api.get).mockResolvedValueOnce(mockUser);

    await useAuthStore.getState().fetchUser();

    expect(api.get).toHaveBeenCalledWith('/auth/me');
  });
});

describe('logout', () => {
  it('clears user and isAuthenticated after logout', async () => {
    useAuthStore.setState({ user: mockUser, isAuthenticated: true, isLoading: false });
    vi.mocked(api.post).mockResolvedValueOnce(undefined);

    await useAuthStore.getState().logout();

    var state = useAuthStore.getState();
    expect(state.user).toBeNull();
    expect(state.isAuthenticated).toBe(false);
  });

  it('calls /auth/logout', async () => {
    vi.mocked(api.post).mockResolvedValueOnce(undefined);

    await useAuthStore.getState().logout();

    expect(api.post).toHaveBeenCalledWith('/auth/logout');
  });

  it('clears user even when the logout request fails', async () => {
    useAuthStore.setState({ user: mockUser, isAuthenticated: true, isLoading: false });
    vi.mocked(api.post).mockRejectedValueOnce(new Error('Network error'));

    // The promise rejects but the finally block still clears state.
    await useAuthStore.getState().logout().catch(() => undefined);

    var state = useAuthStore.getState();
    expect(state.user).toBeNull();
    expect(state.isAuthenticated).toBe(false);
  });
});

describe('setUser', () => {
  it('sets user and marks isAuthenticated true when given a user', () => {
    useAuthStore.getState().setUser(mockUser);

    var state = useAuthStore.getState();
    expect(state.user).toEqual(mockUser);
    expect(state.isAuthenticated).toBe(true);
  });

  it('clears user and marks isAuthenticated false when given null', () => {
    useAuthStore.setState({ user: mockUser, isAuthenticated: true, isLoading: false });
    useAuthStore.getState().setUser(null);

    var state = useAuthStore.getState();
    expect(state.user).toBeNull();
    expect(state.isAuthenticated).toBe(false);
  });
});
