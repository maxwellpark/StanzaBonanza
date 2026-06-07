import { useReviewStanza } from '@/hooks/usePoems';
import type { Stanza } from '@/types/poem';

interface Props {
  poemId: string;
  pendingStanzas: Stanza[];
}

export function StanzaReviewPanel({ poemId, pendingStanzas }: Props) {
  var reviewStanza = useReviewStanza(poemId);

  if (pendingStanzas.length === 0) {
    return null;
  }

  return (
    <section className="mb-8 rounded-xl border border-amber-200 bg-amber-50 p-4">
      <h3 className="mb-4 font-serif text-lg font-bold text-ink">
        Pending Stanzas ({pendingStanzas.length})
      </h3>
      <ul className="space-y-4">
        {pendingStanzas.map((stanza) => (
          <li key={stanza.id} className="rounded-lg border border-parchment-dark bg-white p-4">
            <div className="mb-2 flex items-center gap-2">
              {stanza.author && (
                <>
                  <img
                    src={stanza.author.avatarUrl || '/default-avatar.png'}
                    alt=""
                    className="h-6 w-6 rounded-full object-cover"
                  />
                  <span className="font-sans text-sm text-feather">
                    {stanza.author.displayName}
                  </span>
                </>
              )}
              {stanza.literaryDevice && (
                <span className="ml-auto rounded-full bg-parchment-dark px-2 py-0.5 font-sans text-xs text-feather">
                  {stanza.literaryDevice}
                </span>
              )}
            </div>

            <p className="poem-text mb-4 whitespace-pre-wrap leading-relaxed text-ink">
              {stanza.text}
            </p>

            <div className="flex gap-2">
              <button
                onClick={() => reviewStanza.mutate({ stanzaId: stanza.id, approved: true })}
                disabled={reviewStanza.isPending}
                className="btn-primary text-sm disabled:opacity-50"
              >
                Approve
              </button>
              <button
                onClick={() => reviewStanza.mutate({ stanzaId: stanza.id, approved: false })}
                disabled={reviewStanza.isPending}
                className="btn-secondary text-sm disabled:opacity-50"
              >
                Reject
              </button>
            </div>
          </li>
        ))}
      </ul>
    </section>
  );
}
