import { useState } from 'react';
import { useParams, Navigate } from 'react-router-dom';
import { motion } from 'framer-motion';
import { startRegistration } from '@simplewebauthn/browser';
import { useAuthStore } from '@/stores/authStore';
import { api } from '@/lib/api';
import {
  useUserProfile,
  useFollowers,
  useFollowing,
  useToggleFollow,
  useUpdateProfile,
} from '@/hooks/useSocial';
import { useUserPoems } from '@/hooks/usePoems';
import { PoemCard } from '@/components/poem/PoemCard';
import { formatDate } from '@/lib/utils';
import type { User } from '@/types/user';

type Tab = 'poems' | 'followers' | 'following';

function RegisterPasskeyButton() {
  var [status, setStatus] = useState<'idle' | 'pending' | 'done' | 'error'>('idle');
  var [errorMsg, setErrorMsg] = useState('');

  async function handleRegister() {
    setStatus('pending');
    setErrorMsg('');
    try {
      var options = await api.post<Parameters<typeof startRegistration>[0]>('/auth/register/begin');
      var result = await startRegistration({ optionsJSON: options as any });
      await api.post('/auth/register/finish', result);
      setStatus('done');
    } catch (err) {
      setStatus('error');
      setErrorMsg(err instanceof Error ? err.message : 'Registration failed.');
    }
  }

  if (status === 'done') {
    return <p className="font-sans text-sm text-success">Passkey registered successfully.</p>;
  }

  return (
    <div className="flex flex-col gap-1">
      <button
        onClick={handleRegister}
        disabled={status === 'pending'}
        className="btn-secondary text-sm disabled:opacity-50"
      >
        {status === 'pending' ? 'Follow device prompts...' : 'Register a Passkey'}
      </button>
      {status === 'error' && <p className="font-sans text-xs text-error">{errorMsg}</p>}
    </div>
  );
}

function EditProfileForm({ user, onDone }: { user: User; onDone: () => void }) {
  var [displayName, setDisplayName] = useState(user.displayName);
  var [bio, setBio] = useState(user.bio ?? '');
  var [avatarUrl, setAvatarUrl] = useState(user.avatarUrl ?? '');
  var [error, setError] = useState('');
  var updateProfile = useUpdateProfile();
  var { setUser } = useAuthStore();

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError('');
    if (!displayName.trim()) {
      setError('Display name is required.');
      return;
    }
    try {
      await updateProfile.mutateAsync({ displayName, bio, avatarUrl });
      setUser({ ...user, displayName, bio, avatarUrl });
      onDone();
    } catch {
      setError('Failed to update profile. Please try again.');
    }
  }

  return (
    <form onSubmit={handleSubmit} className="mt-4 flex flex-col gap-3">
      <div>
        <label className="mb-1 block font-sans text-xs text-feather">Display name</label>
        <input
          value={displayName}
          onChange={(e) => setDisplayName(e.target.value)}
          maxLength={50}
          className="w-full rounded-lg border border-parchment-dark bg-white px-3 py-2 font-sans text-sm text-ink outline-none transition-colors focus:border-accent"
        />
      </div>
      <div>
        <label className="mb-1 block font-sans text-xs text-feather">Bio</label>
        <textarea
          value={bio}
          onChange={(e) => setBio(e.target.value)}
          rows={3}
          className="w-full resize-none rounded-lg border border-parchment-dark bg-white px-3 py-2 font-sans text-sm text-ink outline-none transition-colors focus:border-accent"
        />
      </div>
      <div>
        <label className="mb-1 block font-sans text-xs text-feather">Avatar URL</label>
        <input
          value={avatarUrl}
          onChange={(e) => setAvatarUrl(e.target.value)}
          placeholder="https://..."
          className="w-full rounded-lg border border-parchment-dark bg-white px-3 py-2 font-sans text-sm text-ink outline-none transition-colors focus:border-accent"
        />
      </div>
      {error && <p className="font-sans text-sm text-error">{error}</p>}
      <div className="flex gap-2">
        <button type="submit" disabled={updateProfile.isPending} className="btn-primary text-sm disabled:opacity-50">
          {updateProfile.isPending ? 'Saving...' : 'Save'}
        </button>
        <button type="button" onClick={onDone} className="btn-secondary text-sm">
          Cancel
        </button>
      </div>
    </form>
  );
}

function UserList({ users }: { users: User[] }) {
  if (users.length === 0) {
    return <p className="py-8 text-center font-sans text-sm text-feather">Nobody here yet.</p>;
  }
  return (
    <ul className="divide-y divide-parchment-dark">
      {users.map((u) => (
        <li key={u.id} className="flex items-center gap-3 py-3">
          <img
            src={u.avatarUrl || '/default-avatar.png'}
            alt=""
            className="h-10 w-10 rounded-full object-cover"
          />
          <div>
            <p className="font-sans text-sm font-medium text-ink">{u.displayName}</p>
            {u.bio && <p className="font-sans text-xs text-feather">{u.bio}</p>}
          </div>
        </li>
      ))}
    </ul>
  );
}

export function ProfilePage() {
  var { userId } = useParams<{ userId: string }>();
  var { user: me, isAuthenticated } = useAuthStore();
  var [activeTab, setActiveTab] = useState<Tab>('poems');
  var [editing, setEditing] = useState(false);

  var { data: profile, isLoading } = useUserProfile(userId ?? '');
  var { data: poemsData } = useUserPoems(userId ?? '');
  var { data: followersData } = useFollowers(userId ?? '', { pageSize: 50 });
  var { data: followingData } = useFollowing(userId ?? '', { pageSize: 50 });
  var toggleFollow = useToggleFollow(userId ?? '');

  if (!userId) {
    return <Navigate to="/" replace />;
  }

  if (isLoading) {
    return (
      <div className="mx-auto max-w-2xl animate-pulse space-y-4 py-8">
        <div className="flex gap-4">
          <div className="h-20 w-20 rounded-full bg-parchment-dark" />
          <div className="flex-1 space-y-3">
            <div className="h-6 w-1/3 rounded bg-parchment-dark" />
            <div className="h-4 w-2/3 rounded bg-parchment-dark" />
          </div>
        </div>
      </div>
    );
  }

  if (!profile) {
    return (
      <div className="py-16 text-center">
        <p className="font-sans text-feather">User not found.</p>
      </div>
    );
  }

  var isOwnProfile = isAuthenticated && me?.id === userId;
  var followerCount = followersData?.totalCount ?? 0;
  var followingCount = followingData?.totalCount ?? 0;
  var poemCount = poemsData?.totalCount ?? 0;
  var amFollowing = followersData?.items.some((u) => u.id === me?.id) ?? false;

  var tabs: { id: Tab; label: string; count: number }[] = [
    { id: 'poems', label: 'Poems', count: poemCount },
    { id: 'followers', label: 'Followers', count: followerCount },
    { id: 'following', label: 'Following', count: followingCount },
  ];

  return (
    <motion.div
      initial={{ opacity: 0 }}
      animate={{ opacity: 1 }}
      className="mx-auto max-w-2xl"
    >
      <div className="card mb-6">
        <div className="flex flex-col gap-4 sm:flex-row sm:items-start">
          <img
            src={profile.avatarUrl || '/default-avatar.png'}
            alt=""
            className="h-20 w-20 rounded-full object-cover"
          />
          <div className="flex-1">
            <div className="flex flex-wrap items-center gap-3">
              <h1 className="font-serif text-2xl font-bold text-ink">{profile.displayName}</h1>
              {profile.isVerified && (
                <span className="rounded-full bg-accent/10 px-2 py-0.5 font-sans text-xs text-accent">
                  Verified
                </span>
              )}
            </div>
            {profile.bio && (
              <p className="mt-1 font-sans text-sm leading-relaxed text-feather">{profile.bio}</p>
            )}
            <p className="mt-1 font-sans text-xs text-feather">
              Joined {formatDate(profile.createdAt)}
            </p>

            <div className="mt-4 flex flex-wrap gap-2">
              {isOwnProfile ? (
                !editing && (
                  <>
                    <button onClick={() => setEditing(true)} className="btn-secondary text-sm">
                      Edit Profile
                    </button>
                    <RegisterPasskeyButton />
                  </>
                )
              ) : (
                isAuthenticated && (
                  <button
                    onClick={() => toggleFollow.mutate()}
                    disabled={toggleFollow.isPending}
                    className={amFollowing ? 'btn-secondary text-sm' : 'btn-primary text-sm'}
                  >
                    {toggleFollow.isPending ? '...' : amFollowing ? 'Unfollow' : 'Follow'}
                  </button>
                )
              )}
            </div>

            {isOwnProfile && editing && (
              <EditProfileForm user={profile} onDone={() => setEditing(false)} />
            )}
          </div>
        </div>
      </div>

      <div className="mb-6 flex gap-1 rounded-lg border border-parchment-dark bg-parchment p-1">
        {tabs.map((tab) => (
          <button
            key={tab.id}
            onClick={() => setActiveTab(tab.id)}
            className={`flex-1 rounded-md px-3 py-1.5 font-sans text-sm transition-colors ${
              activeTab === tab.id
                ? 'bg-white text-ink shadow-sm'
                : 'text-feather hover:text-ink'
            }`}
          >
            {tab.label}
            <span className="ml-1.5 font-sans text-xs text-feather">{tab.count}</span>
          </button>
        ))}
      </div>

      {activeTab === 'poems' && (
        <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
          {poemsData?.items.map((poem) => (
            <PoemCard key={poem.id} poem={poem} />
          ))}
          {!poemsData?.items.length && (
            <p className="col-span-2 py-8 text-center font-sans text-sm text-feather">
              No poems yet.
            </p>
          )}
        </div>
      )}

      {activeTab === 'followers' && (
        <UserList users={followersData?.items ?? []} />
      )}

      {activeTab === 'following' && (
        <UserList users={followingData?.items ?? []} />
      )}
    </motion.div>
  );
}
