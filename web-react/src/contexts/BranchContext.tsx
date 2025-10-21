import React, { createContext, useContext, useState, useEffect } from 'react';
import type { ReactNode } from 'react';
import type { Branch } from '../types/api';

interface BranchContextType {
  branches: Branch[];
  selectedBranchId: string | null;
  setSelectedBranchId: (branchId: string | null) => void;
  isLoading: boolean;
  error: string | null;
}

const BranchContext = createContext<BranchContextType | undefined>(undefined);

export const useBranch = () => {
  const context = useContext(BranchContext);
  if (context === undefined) {
    throw new Error('useBranch must be used within a BranchProvider');
  }
  return context;
};

interface BranchProviderProps {
  children: ReactNode;
}

export const BranchProvider: React.FC<BranchProviderProps> = ({ children }) => {
  const [branches, setBranches] = useState<Branch[]>([]);
  const [selectedBranchId, setSelectedBranchId] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchBranches = async () => {
      try {
        setIsLoading(true);
        // Mock data for now
        const mockBranches: Branch[] = [
          {
            id: '1',
            name: 'London HQ',
            location: 'London, UK',
            type: 'headquarters',
            status: 'active',
            metrics: {
              totalAssets: 150,
              criticalVulns: 5,
              complianceScore: 85,
              lastScan: '2023-10-27T10:00:00Z'
            },
            coordinates: { lat: 51.5074, lng: -0.1278 }
          },
          {
            id: '2',
            name: 'New York Branch',
            location: 'New York, USA',
            type: 'branch',
            status: 'active',
            metrics: {
              totalAssets: 75,
              criticalVulns: 2,
              complianceScore: 92,
              lastScan: '2023-10-27T09:30:00Z'
            },
            coordinates: { lat: 40.7128, lng: -74.0060 }
          },
          {
            id: '3',
            name: 'Tokyo Datacenter',
            location: 'Tokyo, Japan',
            type: 'datacenter',
            status: 'active',
            metrics: {
              totalAssets: 200,
              criticalVulns: 8,
              complianceScore: 78,
              lastScan: '2023-10-27T08:45:00Z'
            },
            coordinates: { lat: 35.6762, lng: 139.6503 }
          }
        ];
        
        setBranches(mockBranches);
        if (mockBranches.length > 0) {
          setSelectedBranchId(mockBranches[0].id);
        }
      } catch (err) {
        setError('Failed to fetch branches');
        console.error('Error fetching branches:', err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchBranches();
  }, []);

  const value = {
    branches,
    selectedBranchId,
    setSelectedBranchId,
    isLoading,
    error,
  };

  return (
    <BranchContext.Provider value={value}>
      {children}
    </BranchContext.Provider>
  );
};

