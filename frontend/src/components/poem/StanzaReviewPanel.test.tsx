import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { StanzaReviewPanel } from './StanzaReviewPanel';
import type { Stanza } from '@/types/poem';

// useReviewStanza is mocked so the component never hits the network.
var mockMutate = vi.fn();

vi.mock('@/hooks/usePoems', () => ({
  useReviewStanza: vi.fn(() => ({
    mutate: mockMutate,
    isPending: false,
  })),
}));

import { useReviewStanza } from '@/hooks/usePoems';

var pendingStanzas: Stanza[] = [
  {
    id: 'stanza-1',
    poemId: 'poem-1',
    authorId: 'user-2',
    author: {
      id: 'user-2',
      displayName: 'Keats',
      email: 'keats@example.com',
      bio: '',
      avatarUrl: '',
      isVerified: false,
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-01T00:00:00Z',
    },
    text: 'Season of mists and mellow fruitfulness',
    position: 1,
    literaryDevice: 'personification',
    status: 'pending',
    createdAt: '2024-01-01T00:00:00Z',
  },
  {
    id: 'stanza-2',
    poemId: 'poem-1',
    authorId: 'user-3',
    author: {
      id: 'user-3',
      displayName: 'Shelley',
      email: 'shelley@example.com',
      bio: '',
      avatarUrl: '',
      isVerified: false,
      createdAt: '2024-01-01T00:00:00Z',
      updatedAt: '2024-01-01T00:00:00Z',
    },
    text: 'O wild West Wind, thou breath of Autumn\'s being',
    position: 2,
    status: 'pending',
    createdAt: '2024-01-01T00:00:00Z',
  },
];

beforeEach(() => {
  vi.clearAllMocks();
});

describe('StanzaReviewPanel', () => {
  it('renders nothing when there are no pending stanzas', () => {
    var { container } = render(
      <StanzaReviewPanel poemId="poem-1" pendingStanzas={[]} />,
    );
    expect(container.firstChild).toBeNull();
  });

  it('shows the pending stanza count in the heading', () => {
    render(<StanzaReviewPanel poemId="poem-1" pendingStanzas={pendingStanzas} />);
    expect(screen.getByRole('heading', { name: /Pending Stanzas \(2\)/i })).toBeInTheDocument();
  });

  it('renders the stanza text for each pending stanza', () => {
    render(<StanzaReviewPanel poemId="poem-1" pendingStanzas={pendingStanzas} />);
    expect(screen.getByText('Season of mists and mellow fruitfulness')).toBeInTheDocument();
    expect(screen.getByText("O wild West Wind, thou breath of Autumn's being")).toBeInTheDocument();
  });

  it('shows the author display name for each stanza', () => {
    render(<StanzaReviewPanel poemId="poem-1" pendingStanzas={pendingStanzas} />);
    expect(screen.getByText('Keats')).toBeInTheDocument();
    expect(screen.getByText('Shelley')).toBeInTheDocument();
  });

  it('shows the literary device tag when present', () => {
    render(<StanzaReviewPanel poemId="poem-1" pendingStanzas={pendingStanzas} />);
    expect(screen.getByText('personification')).toBeInTheDocument();
  });

  it('renders Approve and Reject buttons for each stanza', () => {
    render(<StanzaReviewPanel poemId="poem-1" pendingStanzas={pendingStanzas} />);
    expect(screen.getAllByRole('button', { name: /approve/i })).toHaveLength(2);
    expect(screen.getAllByRole('button', { name: /reject/i })).toHaveLength(2);
  });

  it('calls mutate with approved: true when Approve is clicked', async () => {
    var user = userEvent.setup();
    render(<StanzaReviewPanel poemId="poem-1" pendingStanzas={[pendingStanzas[0]]} />);

    await user.click(screen.getByRole('button', { name: /approve/i }));

    expect(mockMutate).toHaveBeenCalledWith({ stanzaId: 'stanza-1', approved: true });
  });

  it('calls mutate with approved: false when Reject is clicked', async () => {
    var user = userEvent.setup();
    render(<StanzaReviewPanel poemId="poem-1" pendingStanzas={[pendingStanzas[0]]} />);

    await user.click(screen.getByRole('button', { name: /reject/i }));

    expect(mockMutate).toHaveBeenCalledWith({ stanzaId: 'stanza-1', approved: false });
  });

  it('passes the correct poemId to useReviewStanza', () => {
    render(<StanzaReviewPanel poemId="poem-42" pendingStanzas={pendingStanzas} />);
    expect(useReviewStanza).toHaveBeenCalledWith('poem-42');
  });

  it('disables buttons while a review is pending', () => {
    vi.mocked(useReviewStanza).mockReturnValue({
      mutate: mockMutate,
      isPending: true,
    } as unknown as ReturnType<typeof useReviewStanza>);

    render(<StanzaReviewPanel poemId="poem-1" pendingStanzas={[pendingStanzas[0]]} />);

    for (var btn of screen.getAllByRole('button')) {
      expect(btn).toBeDisabled();
    }
  });
});
