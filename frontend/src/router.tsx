import { createBrowserRouter } from 'react-router-dom';
import { AppShell } from '@/components/layout/AppShell';
import { HomePage } from '@/pages/HomePage';
import { ExplorePage } from '@/pages/ExplorePage';
import { PoemPage } from '@/pages/PoemPage';
import { CreatePoemPage } from '@/pages/CreatePoemPage';
import { HallOfFamePage } from '@/pages/HallOfFamePage';
import { FeedPage } from '@/pages/FeedPage';
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
      { path: '/auth/verify', element: <MagicLinkVerify /> },
      { path: '*', element: <NotFoundPage /> },
    ],
  },
]);
