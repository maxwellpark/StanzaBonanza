import type { User } from './user';

export type PoemFormat =
  | 'free_verse'
  | 'haiku'
  | 'sonnet'
  | 'limerick'
  | 'iambic_pentameter'
  | 'rhyming_couplets'
  | 'custom';

export type ApprovalMode = 'open' | 'approval_required' | 'closed';

export type StanzaStatus = 'approved' | 'pending' | 'rejected';

export interface Poem {
  id: string;
  authorId: string;
  author?: User;
  title: string;
  description: string;
  format: PoemFormat;
  formatRules?: string;
  approvalMode: ApprovalMode;
  maxStanzas?: number;
  isHallOfFame: boolean;
  likeCount: number;
  stanzaCount: number;
  commentCount: number;
  stanzas?: Stanza[];
  createdAt: string;
  updatedAt: string;
}

export interface Stanza {
  id: string;
  poemId: string;
  authorId: string;
  author?: User;
  text: string;
  position: number;
  literaryDevice?: string;
  status: StanzaStatus;
  createdAt: string;
}

export interface CreatePoemInput {
  title: string;
  description?: string;
  text: string;
  format: PoemFormat;
  approvalMode: ApprovalMode;
  maxStanzas?: number;
}

export interface SubmitStanzaInput {
  text: string;
  literaryDevice?: string;
}
