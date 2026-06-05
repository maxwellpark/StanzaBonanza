import type { User } from './user';
import type { Poem } from './poem';

export interface Comment {
  id: string;
  poemId: string;
  authorId: string;
  author?: User;
  parentId?: string;
  text: string;
  createdAt: string;
  updatedAt: string;
}

export type NotificationType =
  | 'like'
  | 'comment'
  | 'follow'
  | 'stanza_submitted'
  | 'stanza_approved'
  | 'stanza_rejected'
  | 'poem_featured';

export interface Notification {
  id: string;
  recipientId: string;
  actorId?: string;
  actor?: User;
  type: NotificationType;
  poemId?: string;
  poem?: Poem;
  read: boolean;
  createdAt: string;
}
