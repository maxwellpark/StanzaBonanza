import type { PoemFormat } from '@/types/poem';
import { formatPoemFormat, cn } from '@/lib/utils';

var formatColors: Record<PoemFormat, string> = {
  free_verse: 'bg-blue-100 text-blue-800',
  haiku: 'bg-emerald-100 text-emerald-800',
  sonnet: 'bg-purple-100 text-purple-800',
  limerick: 'bg-amber-100 text-amber-800',
  iambic_pentameter: 'bg-rose-100 text-rose-800',
  rhyming_couplets: 'bg-cyan-100 text-cyan-800',
  custom: 'bg-gray-100 text-gray-800',
};

interface PoemFormatBadgeProps {
  format: PoemFormat;
  className?: string;
}

export function PoemFormatBadge({ format, className }: PoemFormatBadgeProps) {
  return (
    <span
      className={cn(
        'inline-block rounded-full px-2.5 py-0.5 font-sans text-xs font-medium',
        formatColors[format],
        className
      )}
    >
      {formatPoemFormat(format)}
    </span>
  );
}
