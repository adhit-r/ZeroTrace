import React, { useState, useMemo } from 'react';
import { Package, Search, Monitor, Filter } from 'lucide-react';
import type { Agent } from '../../services/agentService';

interface AssetAppListProps {
    agents: Agent[];
}

interface SoftwareItem {
    name: string;
    version: string;
    vendor?: string;
    type?: string;
    agentName: string;
    agentId: string;
}

export const AssetAppList: React.FC<AssetAppListProps> = ({ agents }) => {
    const [searchTerm, setSearchTerm] = useState('');
    const [selectedAgentId, setSelectedAgentId] = useState<string>('all');
    const [visibleCount, setVisibleCount] = useState(50);

    // Aggregate all software from all agents
    const allSoftware: SoftwareItem[] = useMemo(() => {
        const software: SoftwareItem[] = [];
        agents.forEach(agent => {
            // Check if agent has dependencies/software in metadata
            if (agent.metadata?.dependencies) {
                const deps = agent.metadata.dependencies as any[];
                if (Array.isArray(deps)) {
                    const seen = new Set<string>();
                    deps.forEach(dep => {
                        const key = `${dep.name}:${dep.version}`;
                        if (!seen.has(key)) {
                            seen.add(key);
                            software.push({
                                name: dep.name || 'Unknown',
                                version: dep.version || 'N/A',
                                vendor: dep.description,
                                type: dep.type,
                                agentName: agent.name,
                                agentId: agent.id
                            });
                        }
                    });
                }
            }
        });
        return software;
    }, [agents]);

    // Filter software
    const filteredSoftware = useMemo(() => {
        return allSoftware.filter(item => {
            const matchesSearch = item.name.toLowerCase().includes(searchTerm.toLowerCase()) ||
                item.vendor?.toLowerCase().includes(searchTerm.toLowerCase());
            const matchesAgent = selectedAgentId === 'all' || item.agentId === selectedAgentId;
            return matchesSearch && matchesAgent;
        });
    }, [allSoftware, searchTerm, selectedAgentId]);

    // Get software type distribution
    const typeDistribution = useMemo(() => {
        const dist: Record<string, number> = {};
        allSoftware.forEach(sw => {
            const type = sw.type || 'other';
            dist[type] = (dist[type] || 0) + 1;
        });
        return Object.entries(dist).sort((a, b) => b[1] - a[1]).slice(0, 5);
    }, [allSoftware]);

    const loadMore = () => {
        setVisibleCount(prev => prev + 50);
    };

    return (
        <div className="h-full flex flex-col">
            {/* Header */}
            <div className="p-6 border-b-3 border-black bg-gradient-to-r from-purple-600 to-indigo-600 text-white">
                <div className="flex justify-between items-center">
                    <div className="flex items-center gap-3">
                        <div className="p-2 bg-white/20 rounded-lg">
                            <Package className="w-6 h-6" />
                        </div>
                        <div>
                            <h2 className="text-xl font-black uppercase">Software Inventory</h2>
                            <p className="text-white/70 text-sm">{allSoftware.length.toLocaleString()} total packages across {agents.length} agent(s)</p>
                        </div>
                    </div>
                    <div className="flex gap-2">
                        {typeDistribution.slice(0, 3).map(([type, count]) => (
                            <span key={type} className="px-3 py-1 bg-white/20 rounded-full text-xs font-bold uppercase">
                                {type}: {count.toLocaleString()}
                            </span>
                        ))}
                    </div>
                </div>
            </div>

            {/* Search and Filter */}
            <div className="p-4 bg-gray-50 border-b-2 border-gray-200 flex gap-4">
                <div className="relative flex-1">
                    <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
                    <input
                        type="text"
                        placeholder="Search applications..."
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                        className="w-full bg-white border-2 border-gray-300 rounded-lg pl-10 pr-4 py-2 text-sm focus:outline-none focus:border-purple-500 transition-colors"
                    />
                </div>
                <div className="relative">
                    <Filter className="absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400 w-4 h-4" />
                    <select
                        value={selectedAgentId}
                        onChange={(e) => setSelectedAgentId(e.target.value)}
                        className="bg-white border-2 border-gray-300 rounded-lg pl-10 pr-8 py-2 text-sm focus:outline-none focus:border-purple-500 appearance-none cursor-pointer"
                    >
                        <option value="all">All Agents</option>
                        {agents.map(agent => (
                            <option key={agent.id} value={agent.id}>{agent.name}</option>
                        ))}
                    </select>
                </div>
            </div>

            {/* Results count */}
            <div className="px-6 py-2 bg-gray-100 text-sm text-gray-600 border-b border-gray-200">
                Showing {Math.min(visibleCount, filteredSoftware.length).toLocaleString()} of {filteredSoftware.length.toLocaleString()} packages
            </div>

            {/* Software List */}
            <div className="flex-1 overflow-y-auto">
                {filteredSoftware.length === 0 ? (
                    <div className="text-center text-gray-500 py-16">
                        <Package className="w-12 h-12 mx-auto mb-4 text-gray-300" />
                        <p className="font-bold">No software found</p>
                        <p className="text-sm">Try adjusting your search criteria</p>
                    </div>
                ) : (
                    <>
                        <div className="divide-y divide-gray-100">
                            {filteredSoftware.slice(0, visibleCount).map((item, index) => (
                                <div
                                    key={`${item.agentId}-${item.name}-${index}`}
                                    className="px-6 py-3 hover:bg-gray-50 transition-colors flex justify-between items-center"
                                >
                                    <div className="flex items-center gap-4 flex-1 min-w-0">
                                        <div className="p-2 bg-purple-100 rounded border border-purple-200 flex-shrink-0">
                                            <Package className="w-4 h-4 text-purple-600" />
                                        </div>
                                        <div className="min-w-0 flex-1">
                                            <div className="font-bold text-gray-900 truncate">{item.name}</div>
                                            <div className="text-xs text-gray-500 flex items-center gap-2 mt-0.5">
                                                <span className="bg-gray-200 px-2 py-0.5 rounded font-mono">v{item.version}</span>
                                                {item.type && (
                                                    <span className="text-gray-400">• {item.type}</span>
                                                )}
                                            </div>
                                        </div>
                                    </div>
                                    <div className="flex items-center gap-2 flex-shrink-0">
                                        <span className="text-xs px-2 py-1 bg-indigo-100 text-indigo-700 rounded-full font-bold flex items-center gap-1">
                                            <Monitor className="w-3 h-3" />
                                            {item.agentName}
                                        </span>
                                    </div>
                                </div>
                            ))}
                        </div>

                        {visibleCount < filteredSoftware.length && (
                            <div className="p-4 text-center border-t border-gray-200">
                                <button
                                    onClick={loadMore}
                                    className="px-6 py-2 bg-purple-600 text-white rounded-lg font-bold hover:bg-purple-700 transition-colors"
                                >
                                    Load More ({(filteredSoftware.length - visibleCount).toLocaleString()} remaining)
                                </button>
                            </div>
                        )}
                    </>
                )}
            </div>
        </div>
    );
};
