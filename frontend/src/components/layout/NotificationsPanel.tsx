import { useState, useRef, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { motion, AnimatePresence } from 'framer-motion';
import { useNotifications, useMarkNotificationsRead } from '@/hooks/useSocial';
import type { Notification, NotificationType } from '@/types/social';
import { timeAgo, cn } from '@/lib/utils';

function notifLabel(n: Notification): string {
  var actor = n.actor?.displayName ?? 'Someone';
  var title = n.poem?.title ? `"${n.poem.title}"` : 'your poem';
  switch (n.type) {
    case 'like': return `${actor} liked ${title}`;
    case 'comment': return `${actor} commented on ${title}`;
    case 'follow': return `${actor} followed you`;
    case 'stanza_submitted': return `${actor} submitted a stanza to ${title}`;
    case 'stanza_approved': return `Your stanza in ${title} was approved`;
    case 'stanza_rejected': return `Your stanza in ${title} was rejected`;
    case 'poem_featured': return `${title} was featured in Hall of Fame`;
    default: return 'New notification';
  }
}

function notifHref(n: Notification): string {
  if (n.type === 'follow' && n.actorId) {
    return `/profile/${n.actorId}`;
  }
  if (n.poemId) {
    return `/poems/${n.poemId}`;
  }
  return '/';
}

const iconPaths: Record<NotificationType, string> = {
  like: 'M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z',
  comment: 'M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z',
  follow: 'M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z',
  stanza_submitted: 'M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z',
  stanza_approved: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z',
  stanza_rejected: 'M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z',
  poem_featured: 'M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z',
};

function iconColorClass(type: NotificationType) {
  if (type === 'stanza_approved') {
    return 'bg-green-100 text-green-600';
  }
  if (type === 'stanza_rejected') {
    return 'bg-red-100 text-red-600';
  }
  if (type === 'poem_featured') {
    return 'bg-yellow-100 text-yellow-600';
  }
  return 'bg-parchment-dark text-feather';
}

export function NotificationsPanel() {
  var [isOpen, setIsOpen] = useState(false);
  var containerRef = useRef<HTMLDivElement>(null);
  var navigate = useNavigate();
  var { data, isLoading } = useNotifications({ pageSize: 20 });
  var markRead = useMarkNotificationsRead();

  var notifications = data?.items ?? [];
  var unreadCount = notifications.filter((n) => !n.read).length;

  useEffect(() => {
    if (!isOpen) {
      return;
    }
    var ids = notifications.filter((n) => !n.read).map((n) => n.id);
    if (ids.length > 0) {
      markRead.mutate(ids);
    }
  }, [isOpen]);

  useEffect(() => {
    function onMouseDown(e: MouseEvent) {
      if (containerRef.current && !containerRef.current.contains(e.target as Node)) {
        setIsOpen(false);
      }
    }
    if (isOpen) {
      document.addEventListener('mousedown', onMouseDown);
    }
    return () => document.removeEventListener('mousedown', onMouseDown);
  }, [isOpen]);

  function handleNotifClick(n: Notification) {
    setIsOpen(false);
    navigate(notifHref(n));
  }

  return (
    <div ref={containerRef} className="relative">
      <button
        onClick={() => setIsOpen((o) => !o)}
        className="relative text-feather transition-colors hover:text-ink"
        aria-label="Notifications"
      >
        <svg className="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
          <path strokeLinecap="round" strokeLinejoin="round" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" />
        </svg>
        {unreadCount > 0 && (
          <span className="absolute -right-1 -top-1 flex h-4 w-4 items-center justify-center rounded-full bg-accent font-sans text-[10px] font-bold text-white">
            {unreadCount > 9 ? '9+' : unreadCount}
          </span>
        )}
      </button>

      <AnimatePresence>
        {isOpen && (
          <motion.div
            initial={{ opacity: 0, y: -8 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -8 }}
            transition={{ duration: 0.15 }}
            className="absolute right-0 top-full z-50 mt-2 w-80 rounded-xl border border-parchment-dark bg-white shadow-lg"
          >
            <div className="flex items-center justify-between border-b border-parchment-dark px-4 py-3">
              <h3 className="font-sans text-sm font-semibold text-ink">Notifications</h3>
              {unreadCount > 0 && (
                <span className="rounded-full bg-accent px-2 py-0.5 font-sans text-xs text-white">
                  {unreadCount} unread
                </span>
              )}
            </div>

            <div className="max-h-96 overflow-y-auto">
              {isLoading ? (
                <div className="space-y-3 p-4">
                  {Array.from({ length: 4 }).map((_, i) => (
                    <div key={i} className="flex animate-pulse gap-3">
                      <div className="h-8 w-8 flex-shrink-0 rounded-full bg-parchment-dark" />
                      <div className="flex-1 space-y-2">
                        <div className="h-3 w-full rounded bg-parchment-dark" />
                        <div className="h-3 w-2/3 rounded bg-parchment-dark" />
                      </div>
                    </div>
                  ))}
                </div>
              ) : notifications.length === 0 ? (
                <p className="px-4 py-8 text-center font-sans text-sm text-feather">
                  No notifications yet.
                </p>
              ) : (
                <ul>
                  {notifications.map((n) => (
                    <li key={n.id}>
                      <button
                        onClick={() => handleNotifClick(n)}
                        className={cn(
                          'flex w-full items-start gap-3 px-4 py-3 text-left transition-colors hover:bg-parchment',
                          !n.read && 'bg-parchment/60'
                        )}
                      >
                        <span className={cn(
                          'flex h-8 w-8 flex-shrink-0 items-center justify-center rounded-full',
                          iconColorClass(n.type)
                        )}>
                          <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                            <path strokeLinecap="round" strokeLinejoin="round" d={iconPaths[n.type]} />
                          </svg>
                        </span>
                        <div className="min-w-0 flex-1">
                          <p className="font-sans text-sm leading-snug text-ink">{notifLabel(n)}</p>
                          <p className="mt-0.5 font-sans text-xs text-feather">{timeAgo(n.createdAt)}</p>
                        </div>
                        {!n.read && (
                          <span className="mt-2 h-2 w-2 flex-shrink-0 rounded-full bg-accent" />
                        )}
                      </button>
                    </li>
                  ))}
                </ul>
              )}
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
