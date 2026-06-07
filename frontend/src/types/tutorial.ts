import type { PoemFormat } from './poem';

export type Difficulty = 'beginner' | 'intermediate' | 'advanced';

export interface Tutorial {
  id: string;
  title: string;
  slug: string;
  format: PoemFormat;
  contentMd: string;
  difficulty: Difficulty;
  displayOrder: number;
  createdAt: string;
}
