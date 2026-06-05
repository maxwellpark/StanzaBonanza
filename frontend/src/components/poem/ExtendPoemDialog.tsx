import { useState } from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import type { PoemFormat } from '@/types/poem';
import { useSubmitStanza } from '@/hooks/usePoems';
import { formatPoemFormat } from '@/lib/utils';

interface ExtendPoemDialogProps {
  poemId: string;
  format: PoemFormat;
  isOpen: boolean;
  onClose: () => void;
}

var literaryDevices = [
  'metaphor',
  'simile',
  'alliteration',
  'enjambment',
  'imagery',
  'personification',
  'hyperbole',
  'onomatopoeia',
  'irony',
];

export function ExtendPoemDialog({ poemId, format, isOpen, onClose }: ExtendPoemDialogProps) {
  var [text, setText] = useState('');
  var [device, setDevice] = useState('');
  var mutation = useSubmitStanza(poemId);

  function handleBackdropClick(e: React.MouseEvent) {
    if (e.target === e.currentTarget) {
      onClose();
    }
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!text.trim()) {
      return;
    }

    try {
      await mutation.mutateAsync({
        text: text.trim(),
        literaryDevice: device || undefined,
      });
      setText('');
      setDevice('');
      onClose();
      alert('Stanza submitted successfully!');
    } catch (err) {
      if (err instanceof Error) {
        alert(`Error: ${err.message}`);
      } else {
        alert('Something went wrong. Please try again.');
      }
    }
  }

  return (
    <AnimatePresence>
      {isOpen && (
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          transition={{ duration: 0.2 }}
          className="fixed inset-0 z-[100] flex items-center justify-center bg-black/50 px-4"
          onClick={handleBackdropClick}
        >
          <motion.div
            initial={{ opacity: 0, scale: 0.95, y: 10 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.95, y: 10 }}
            transition={{ duration: 0.2 }}
            className="relative w-full max-w-lg rounded-xl border border-parchment-dark bg-parchment p-6 shadow-lg"
          >
            <button
              onClick={onClose}
              className="absolute right-4 top-4 text-feather transition-colors hover:text-ink"
              aria-label="Close"
            >
              <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>

            <h2 className="mb-1 font-serif text-xl font-bold text-ink">Add a Stanza</h2>
            <p className="mb-4 font-sans text-sm text-feather">
              Format: {formatPoemFormat(format)}
            </p>

            <form onSubmit={handleSubmit} className="flex flex-col gap-4">
              <textarea
                value={text}
                onChange={(e) => setText(e.target.value)}
                placeholder="Write your stanza..."
                rows={6}
                required
                className="w-full resize-none rounded-lg border border-parchment-dark bg-white px-4 py-3 text-base leading-relaxed text-ink outline-none transition-colors focus:border-accent"
                style={{ fontFamily: 'var(--font-body)' }}
              />

              <div>
                <label className="mb-1 block font-sans text-sm text-feather">
                  Literary Device (optional)
                </label>
                <select
                  value={device}
                  onChange={(e) => setDevice(e.target.value)}
                  className="w-full rounded-lg border border-parchment-dark bg-white px-4 py-2 font-sans text-sm text-ink outline-none transition-colors focus:border-accent"
                >
                  <option value="">None</option>
                  {literaryDevices.map((d) => (
                    <option key={d} value={d}>
                      {d.charAt(0).toUpperCase() + d.slice(1)}
                    </option>
                  ))}
                </select>
              </div>

              <button
                type="submit"
                disabled={mutation.isPending || !text.trim()}
                className="btn-primary disabled:opacity-50"
              >
                {mutation.isPending ? 'Submitting...' : 'Submit Stanza'}
              </button>
            </form>
          </motion.div>
        </motion.div>
      )}
    </AnimatePresence>
  );
}
