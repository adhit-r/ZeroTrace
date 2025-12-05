/**
 * Main application store using Zustand
 * Optimized with selectors to prevent unnecessary re-renders
 */
import { create } from 'zustand';
import { devtools, persist } from 'zustand/middleware';

interface AppState {
  // UI State
  sidebarOpen: boolean;
  theme: 'light' | 'dark' | 'system';
  selectedOrganizationId: string | null;
  
  // User preferences
  preferences: {
    itemsPerPage: number;
    autoRefresh: boolean;
    refreshInterval: number;
    notificationsEnabled: boolean;
  };
  
  // Actions
  toggleSidebar: () => void;
  setSidebarOpen: (open: boolean) => void;
  setTheme: (theme: 'light' | 'dark' | 'system') => void;
  setSelectedOrganizationId: (id: string | null) => void;
  updatePreferences: (prefs: Partial<AppState['preferences']>) => void;
}

export const useAppStore = create<AppState>()(
  devtools(
    persist(
      (set) => ({
        // Initial state
        sidebarOpen: true,
        theme: 'dark',
        selectedOrganizationId: null,
        preferences: {
          itemsPerPage: 25,
          autoRefresh: true,
          refreshInterval: 30000, // 30 seconds
          notificationsEnabled: true,
        },
        
        // Actions
        toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
        setSidebarOpen: (open) => set({ sidebarOpen: open }),
        setTheme: (theme) => set({ theme }),
        setSelectedOrganizationId: (id) => set({ selectedOrganizationId: id }),
        updatePreferences: (prefs) =>
          set((state) => ({
            preferences: { ...state.preferences, ...prefs },
          })),
      }),
      {
        name: 'zerotrace-app-storage',
        partialize: (state) => ({
          theme: state.theme,
          preferences: state.preferences,
          sidebarOpen: state.sidebarOpen,
        }),
      }
    ),
    { name: 'AppStore' }
  )
);

// Optimized selectors - components only re-render when selected state changes
export const useSidebarOpen = () => useAppStore((state) => state.sidebarOpen);
export const useTheme = () => useAppStore((state) => state.theme);
export const useSelectedOrganizationId = () => useAppStore((state) => state.selectedOrganizationId);
export const usePreferences = () => useAppStore((state) => state.preferences);

