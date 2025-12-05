/**
 * Agent store with optimized selectors
 */
import { create } from 'zustand';
import { devtools } from 'zustand/middleware';

export interface Agent {
  id: string;
  name: string;
  status: 'online' | 'offline' | 'error';
  last_seen: string;
  version: string;
  organization_id: string;
}

interface AgentState {
  agents: Agent[];
  selectedAgent: Agent | null;
  filters: {
    status: string[];
    search: string;
  };
  
  // Actions
  setAgents: (agents: Agent[]) => void;
  addAgent: (agent: Agent) => void;
  updateAgent: (id: string, updates: Partial<Agent>) => void;
  setSelectedAgent: (agent: Agent | null) => void;
  setFilters: (filters: Partial<AgentState['filters']>) => void;
  
  // Computed values
  onlineAgents: () => Agent[];
  offlineAgents: () => Agent[];
}

export const useAgentStore = create<AgentState>()(
  devtools(
    (set, get) => ({
      agents: [],
      selectedAgent: null,
      filters: {
        status: [],
        search: '',
      },
      
      setAgents: (agents) => set({ agents }),
      
      addAgent: (agent) =>
        set((state) => ({
          agents: [...state.agents, agent],
        })),
      
      updateAgent: (id, updates) =>
        set((state) => ({
          agents: state.agents.map((a) => (a.id === id ? { ...a, ...updates } : a)),
          selectedAgent:
            state.selectedAgent?.id === id
              ? { ...state.selectedAgent, ...updates }
              : state.selectedAgent,
        })),
      
      setSelectedAgent: (agent) => set({ selectedAgent: agent }),
      
      setFilters: (filters) =>
        set((state) => ({
          filters: { ...state.filters, ...filters },
        })),
      
      // Computed selectors
      onlineAgents: () => get().agents.filter((a) => a.status === 'online'),
      offlineAgents: () => get().agents.filter((a) => a.status === 'offline'),
    }),
    { name: 'AgentStore' }
  )
);

// Optimized selectors
export const useAgents = () => useAgentStore((state) => state.agents);
export const useSelectedAgent = () => useAgentStore((state) => state.selectedAgent);
export const useAgentFilters = () => useAgentStore((state) => state.filters);
export const useOnlineAgents = () => useAgentStore((state) => state.onlineAgents());
export const useOfflineAgents = () => useAgentStore((state) => state.offlineAgents());

