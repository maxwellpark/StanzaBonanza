import { createBrowserRouter } from 'react-router-dom';
import { AppShell } from '@/components/layout/AppShell';
import { HomePage } from '@/pages/HomePage';
import { ExplorePage } from '@/pages/ExplorePage';
import { PoemPage } from '@/pages/PoemPage';
import { CreatePoemPage } from '@/pages/CreatePoemPage';
import { HallOfFamePage } from '@/pages/HallOfFamePage';
import { FeedPage } from '@/pages/FeedPage';
import { ProfilePage } from '@/pages/ProfilePage';
import { TutorialsPage } from '@/pages/TutorialsPage';
import { TutorialDetailPage } from '@/pages/TutorialDetailPage';
import { MagicLinkVerify } from '@/components/auth/MagicLinkVerify';
import { NotFoundPage } from '@/pages/NotFoundPage';

export var router = createBrowserRouter([
  {
    element: <AppShell />,
    children: [
      { path: '/', element: <HomePage /> },
      { path: '/explore', element: <ExplorePage /> },
      { path: '/poems/new', element: <CreatePoemPage /> },
      { path: '/poems/:poemId', element: <PoemPage /> },
      { path: '/hall-of-fame', element: <HallOfFamePage /> },
      { path: '/feed', element: <FeedPage /> },
      { path: '/profile/:userId', element: <ProfilePage /> },
      { path: '/tutorials', element: <TutorialsPage /> },
      { path: '/tutorials/:slug', element: <TutorialDetailPage /> },
      { path: '/auth/verify', element: <MagicLinkVerify /> },
      { path: '*', element: <NotFoundPage /> },
    ],
  },
]);
