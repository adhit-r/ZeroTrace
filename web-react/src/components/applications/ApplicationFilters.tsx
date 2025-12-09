import React from 'react';
import { Filter, Search } from 'lucide-react';
import { Input } from '@/components/ui/input';
import { Card } from '@/components/ui/card';

export type GroupingType = 'agent' | 'risk' | 'classification';
export type FilterType = 'all' | 'vulnerable' | 'safe';

interface ApplicationFiltersProps {
  searchTerm: string;
  onSearchChange: (term: string) => void;
  grouping: GroupingType;
  onGroupingChange: (grouping: GroupingType) => void;
  filterType: FilterType;
  onFilterChange: (filter: FilterType) => void;
}

const ApplicationFilters: React.FC<ApplicationFiltersProps> = ({
  searchTerm,
  onSearchChange,
  grouping,
  onGroupingChange,
  filterType,
  onFilterChange,
}) => {
  return (
    <Card className="p-4 bg-white border-4 border-black rounded-lg">
      <div className="space-y-4">
        {/* Search */}
        <div className="relative">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-5 w-5 text-gray-400" />
          <Input
            type="text"
            placeholder="Search applications..."
            value={searchTerm}
            onChange={(e) => onSearchChange(e.target.value)}
            className="pl-10 border-4 border-black"
          />
        </div>

        {/* Grouping Options */}
        <div>
          <div className="flex items-center gap-2 mb-2">
            <Filter className="h-4 w-4 text-gray-600" />
            <label className="text-sm font-bold text-gray-700 uppercase">Group By:</label>
          </div>
          <div className="flex gap-2 flex-wrap">
            <button
              onClick={() => onGroupingChange('agent')}
              className={`px-4 py-2 font-bold uppercase border-4 border-black rounded-lg transition-colors ${
                grouping === 'agent'
                  ? 'bg-orange-500 text-white'
                  : 'bg-white text-black hover:bg-gray-50'
              }`}
            >
              Agent
            </button>
            <button
              onClick={() => onGroupingChange('risk')}
              className={`px-4 py-2 font-bold uppercase border-4 border-black rounded-lg transition-colors ${
                grouping === 'risk'
                  ? 'bg-orange-500 text-white'
                  : 'bg-white text-black hover:bg-gray-50'
              }`}
            >
              Risk Level
            </button>
            <button
              onClick={() => onGroupingChange('classification')}
              className={`px-4 py-2 font-bold uppercase border-4 border-black rounded-lg transition-colors ${
                grouping === 'classification'
                  ? 'bg-orange-500 text-white'
                  : 'bg-white text-black hover:bg-gray-50'
              }`}
            >
              Classification
            </button>
          </div>
        </div>

        {/* Status Filter */}
        <div>
          <label className="text-sm font-bold text-gray-700 uppercase mb-2 block">Status:</label>
          <div className="flex gap-2 flex-wrap">
            <button
              onClick={() => onFilterChange('all')}
              className={`px-4 py-2 font-bold uppercase border-4 border-black rounded-lg transition-colors ${
                filterType === 'all'
                  ? 'bg-orange-500 text-white'
                  : 'bg-white text-black hover:bg-gray-50'
              }`}
            >
              All
            </button>
            <button
              onClick={() => onFilterChange('vulnerable')}
              className={`px-4 py-2 font-bold uppercase border-4 border-black rounded-lg transition-colors ${
                filterType === 'vulnerable'
                  ? 'bg-red-500 text-white'
                  : 'bg-white text-black hover:bg-gray-50'
              }`}
            >
              Vulnerable
            </button>
            <button
              onClick={() => onFilterChange('safe')}
              className={`px-4 py-2 font-bold uppercase border-4 border-black rounded-lg transition-colors ${
                filterType === 'safe'
                  ? 'bg-green-500 text-white'
                  : 'bg-white text-black hover:bg-gray-50'
              }`}
            >
              Safe
            </button>
          </div>
        </div>
      </div>
    </Card>
  );
};

export default ApplicationFilters;

