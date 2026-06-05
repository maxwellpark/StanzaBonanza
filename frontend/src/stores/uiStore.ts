import { create } from 'zustand';

interface UIState {
  isLoginOpen: boolean;
  openLogin: () => void;
  closeLogin: () => void;
}

export const useUIStore = create<UIState>((set) => ({
  isLoginOpen: false,
  openLogin: () => set({ isLoginOpen: true }),
  closeLogin: () => set({ isLoginOpen: false }),
}));
