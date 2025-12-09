import React, { useCallback, useEffect, useState, useMemo } from 'react';
import ReactFlow, {
  Controls,
  Background,
  MiniMap,
  useNodesState,
  useEdgesState,
  Panel,
  MarkerType,
  Position,
  Handle,
  type Node,
  type Edge,
} from 'reactflow';
import 'reactflow/dist/style.css';
import { api } from '../../services/api';
import {
  Target,
  ArrowRight,
  Shield,
  AlertTriangle,
  XCircle,
  Zap,
  Lock,
  Database,
  Download,
  RefreshCw,
  Search,
  TrendingUp,
} from 'lucide-react';

// Custom node for attack steps
const AttackStepNode = ({ data, selected }: { data: any; selected: boolean }) => {
  const getStepIcon = () => {
    switch (data.stepType) {
      case 'initial_access':
        return <Target className="w-5 h-5" />;
      case 'lateral_movement':
        return <ArrowRight className="w-5 h-5" />;
      case 'privilege_escalation':
        return <TrendingUp className="w-5 h-5" />;
      case 'data_exfiltration':
        return <Download className="w-5 h-5" />;
      case 'persistence':
        return <Lock className="w-5 h-5" />;
      default:
        return <Zap className="w-5 h-5" />;
    }
  };

  const getStepColor = () => {
    const risk = data.likelihood * data.impact;
    if (risk >= 0.6) return 'bg-red-500 border-red-600';
    if (risk >= 0.3) return 'bg-orange-500 border-orange-600';
    if (risk >= 0.1) return 'bg-yellow-500 border-yellow-600';
    return 'bg-blue-500 border-blue-600';
  };

  return (
    <div
      className={`px-4 py-3 rounded-lg border-2 shadow-lg transition-all min-w-[220px] ${
        selected ? 'ring-4 ring-blue-400' : ''
      } ${getStepColor()} text-white`}
    >
      <Handle type="target" position={Position.Left} className="!bg-white !border-2 !border-gray-800" />
      
      <div className="flex items-start gap-2 mb-2">
        <div className="mt-0.5">{getStepIcon()}</div>
        <div className="flex-1">
          <div className="font-bold text-sm mb-1">{data.label}</div>
          <div className="text-xs opacity-90">{data.technique}</div>
        </div>
        <div className="text-xs font-semibold bg-white/20 px-2 py-1 rounded">
          Step {data.stepNumber}
        </div>
      </div>

      <div className="mt-2 space-y-1 text-xs border-t border-white/20 pt-2">
        <div className="flex justify-between">
          <span className="opacity-80">Target:</span>
          <span className="font-semibold">{data.targetHostname || data.target}</span>
        </div>
        {data.targetIP && (
          <div className="flex justify-between">
            <span className="opacity-80">IP:</span>
            <span className="font-mono text-xs">{data.targetIP}</span>
          </div>
        )}
        {data.cveID && (
          <div className="flex justify-between mt-1 pt-1 border-t border-white/20">
            <span className="opacity-80">CVE:</span>
            <span className="font-mono text-xs font-bold">{data.cveID}</span>
          </div>
        )}
        <div className="flex justify-between">
          <span className="opacity-80">Likelihood:</span>
          <span className="font-semibold">{(data.likelihood * 100).toFixed(0)}%</span>
        </div>
        <div className="flex justify-between">
          <span className="opacity-80">Impact:</span>
          <span className="font-semibold">{(data.impact * 100).toFixed(0)}%</span>
        </div>
        {data.detectionDifficulty && (
          <div className="flex items-center gap-1 mt-1">
            <Shield className="w-3 h-3" />
            <span className="opacity-80">Detection: {data.detectionDifficulty}</span>
          </div>
        )}
      </div>

      <Handle type="source" position={Position.Right} className="!bg-white !border-2 !border-gray-800" />
    </div>
  );
};

// Custom node for entry/exit points
const EntryPointNode = ({ data, selected }: { data: any; selected: boolean }) => {
  return (
    <div
      className={`px-4 py-3 rounded-full border-4 shadow-lg transition-all ${
        selected ? 'ring-4 ring-blue-400' : ''
      } ${data.type === 'entry' ? 'bg-green-500 border-green-600' : 'bg-red-500 border-red-600'} text-white font-bold text-sm min-w-[120px] text-center`}
    >
      <Handle 
        type={data.type === 'entry' ? 'source' : 'target'} 
        position={data.type === 'entry' ? Position.Right : Position.Left}
        className="!bg-white !border-2 !border-gray-800"
      />
      <div className="flex items-center justify-center gap-2">
        {data.type === 'entry' ? (
          <>
            <Target className="w-4 h-4" />
            <span>Entry Point</span>
          </>
        ) : (
          <>
            <Database className="w-4 h-4" />
            <span>Target</span>
          </>
        )}
      </div>
    </div>
  );
};

interface AttackPath {
  path_id: string;
  name: string;
  steps: Array<{
    step_number: number;
    action: string;
    target: string;
    target_ip?: string;
    target_hostname?: string;
    technique: string;
    technique_id?: string;
    likelihood: number;
    impact: number;
    detection_difficulty: string;
    step_type: string;
    cve_id?: string;
    vulnerability_id?: string;
    proof?: string;
    mitigation_controls?: string[];
  }>;
  total_likelihood: number;
  total_impact: number;
  criticality_score: number;
  mitigation_priority: string;
  detection_points?: string[];
  prevention_controls?: string[];
}

interface AttackPathVisualizerProps {
  className?: string;
  onPathSelect?: (path: AttackPath) => void;
}

const AttackPathVisualizer: React.FC<AttackPathVisualizerProps> = ({
  className = '',
  onPathSelect,
}) => {
  const [nodes, setNodes, onNodesChange] = useNodesState([]);
  const [edges, setEdges, onEdgesChange] = useEdgesState([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [attackPaths, setAttackPaths] = useState<AttackPath[]>([]);
  const [selectedPath, setSelectedPath] = useState<string | null>(null);
  const [selectedPaths, setSelectedPaths] = useState<Set<string>>(new Set());
  const [viewMode, setViewMode] = useState<'single' | 'multiple' | 'all'>('multiple');
  const [searchTerm, setSearchTerm] = useState('');
  const [sortBy, setSortBy] = useState<'criticality' | 'likelihood' | 'impact' | 'steps'>('criticality');

  // Fetch attack paths
  useEffect(() => {
    const fetchAttackPaths = async () => {
      try {
        setIsLoading(true);
        setError(null);

        // Fetch from real API
        const response = await api.get('/api/v2/attack-paths');
        const data = response.data;
        
        if (!data.success) {
          throw new Error(data.message || 'Failed to fetch attack paths');
        }

        const paths: AttackPath[] = data.data || [];

        // Fallback to mock data if API returns empty
        const mockPaths: AttackPath[] = [
          {
            path_id: 'ap_001',
            name: 'External to Internal Data Exfiltration',
            steps: [
              {
                step_number: 1,
                action: 'Exploit external vulnerability',
                target: 'Web Server (10.0.1.5)',
                technique: 'SQL Injection',
                likelihood: 0.7,
                impact: 0.6,
                detection_difficulty: 'medium',
                step_type: 'initial_access',
              },
              {
                step_number: 2,
                action: 'Lateral movement',
                target: 'Database Server (10.0.2.10)',
                technique: 'Privilege Escalation',
                likelihood: 0.5,
                impact: 0.8,
                detection_difficulty: 'high',
                step_type: 'lateral_movement',
              },
              {
                step_number: 3,
                action: 'Data exfiltration',
                target: 'Sensitive Data',
                technique: 'Data Export',
                likelihood: 0.4,
                impact: 0.9,
                detection_difficulty: 'high',
                step_type: 'data_exfiltration',
              },
            ],
            total_likelihood: 0.14,
            total_impact: 0.9,
            criticality_score: 0.126,
            mitigation_priority: 'high',
          },
        ];

        const finalPaths = paths.length > 0 ? paths : mockPaths;
        setAttackPaths(finalPaths);
        if (finalPaths.length > 0) {
          // Auto-select top 3 paths by default in multiple mode
          const topPaths = finalPaths.slice(0, 3).map(p => p.path_id);
          setSelectedPaths(new Set(topPaths));
          visualizePaths(finalPaths.slice(0, 3));
          setSelectedPath(finalPaths[0].path_id);
        }
      } catch (err: any) {
        setError(err.message || 'Failed to load attack paths');
      } finally {
        setIsLoading(false);
      }
    };

    fetchAttackPaths();
  }, []);

  // Visualize multiple paths with smart layout
  const visualizePaths = useCallback((paths: AttackPath[]) => {
    if (paths.length === 0) return;

    const newNodes: Node[] = [];
    const newEdges: Edge[] = [];
    const nodeMap = new Map<string, Node>(); // Track shared nodes
    const pathColors = [
      '#ef4444', // red
      '#f97316', // orange
      '#eab308', // yellow
      '#3b82f6', // blue
      '#8b5cf6', // purple
      '#ec4899', // pink
    ];

    // Shared entry point
    const entryNode: Node = {
      id: 'entry',
      type: 'entryPoint',
      position: { x: 50, y: 400 },
      data: { type: 'entry', label: 'External Network' },
    };
    newNodes.push(entryNode);
    nodeMap.set('entry', entryNode);

    // Layout paths vertically with spacing
    const pathSpacing = 250;
    const startY = 100;

    paths.forEach((path, pathIndex) => {
      const pathColor = pathColors[pathIndex % pathColors.length];
      const yOffset = startY + pathIndex * pathSpacing;
      const pathPrefix = `path-${path.path_id}`;

      // Attack steps for this path
      path.steps.forEach((step, stepIndex) => {
        const nodeId = `${pathPrefix}-step-${step.step_number}`;
        const xPos = 300 + stepIndex * 280;

        // Check if we can reuse a node (same target, same step type)
        const nodeKey = `${step.target}-${step.step_type}`;
        let node: Node;

        if (nodeMap.has(nodeKey) && viewMode === 'all') {
          // Reuse existing node for branching
          node = nodeMap.get(nodeKey)!;
        } else {
          node = {
            id: nodeId,
            type: 'attackStep',
            position: { x: xPos, y: yOffset },
            data: {
              label: step.action,
              stepNumber: step.step_number,
              target: step.target,
              targetIP: step.target_ip,
              targetHostname: step.target_hostname,
              technique: step.technique,
              techniqueID: step.technique_id,
              likelihood: step.likelihood,
              impact: step.impact,
              detectionDifficulty: step.detection_difficulty,
              stepType: step.step_type,
              cveID: step.cve_id,
              vulnerabilityID: step.vulnerability_id,
              proof: step.proof,
              mitigationControls: step.mitigation_controls,
              pathId: path.path_id,
              pathColor: pathColor,
            },
          };
          newNodes.push(node);
          if (viewMode === 'all') {
            nodeMap.set(nodeKey, node);
          }
        }

        // Edge from previous node
        const sourceId = stepIndex === 0 
          ? 'entry' 
          : `${pathPrefix}-step-${path.steps[stepIndex - 1].step_number}`;
        
        newEdges.push({
          id: `edge-${pathPrefix}-${stepIndex}`,
          source: sourceId,
          target: node.id,
          type: 'smoothstep',
          animated: true,
          style: { 
            stroke: pathColor, 
            strokeWidth: 3,
            strokeDasharray: pathIndex > 0 ? '5,5' : undefined,
          },
          markerEnd: { type: MarkerType.ArrowClosed, color: pathColor },
          label: `${(step.likelihood * 100).toFixed(0)}%`,
          labelStyle: { fill: pathColor, fontWeight: 'bold', fontSize: 11 },
        });
      });

      // Exit point for this path
      const exitId = `${pathPrefix}-exit`;
      const lastStepId = `${pathPrefix}-step-${path.steps[path.steps.length - 1].step_number}`;
      const exitNode: Node = {
        id: exitId,
        type: 'entryPoint',
        position: { x: 300 + path.steps.length * 280, y: yOffset },
        data: { 
          type: 'exit', 
          label: path.name.includes('Data') ? 'Data Exfiltrated' : 'Target Compromised',
          pathId: path.path_id,
        },
      };
      newNodes.push(exitNode);

      newEdges.push({
        id: `edge-${pathPrefix}-exit`,
        source: lastStepId,
        target: exitId,
        type: 'smoothstep',
        animated: true,
        style: { stroke: pathColor, strokeWidth: 3 },
        markerEnd: { type: MarkerType.ArrowClosed, color: pathColor },
      });
    });

    setNodes(newNodes);
    setEdges(newEdges);
  }, [setNodes, setEdges, viewMode]);

  // Visualize single path (backward compatibility)
  const visualizePath = useCallback((path: AttackPath) => {
    visualizePaths([path]);
  }, [visualizePaths]);

  const nodeTypes = useMemo(
    () => ({
      attackStep: AttackStepNode,
      entryPoint: EntryPointNode,
    }),
    []
  );

  // Sort and filter paths
  const filteredPaths = useMemo(() => {
    let filtered = attackPaths.filter((path) =>
      path.name.toLowerCase().includes(searchTerm.toLowerCase())
    );

    // Sort paths
    filtered.sort((a, b) => {
      switch (sortBy) {
        case 'criticality':
          return b.criticality_score - a.criticality_score;
        case 'likelihood':
          return b.total_likelihood - a.total_likelihood;
        case 'impact':
          return b.total_impact - a.total_impact;
        case 'steps':
          return b.steps.length - a.steps.length;
        default:
          return 0;
      }
    });

    return filtered;
  }, [attackPaths, searchTerm, sortBy]);

  const handlePathSelect = (path: AttackPath, toggle: boolean = false) => {
    if (viewMode === 'single') {
      setSelectedPath(path.path_id);
      setSelectedPaths(new Set([path.path_id]));
      visualizePath(path);
      onPathSelect?.(path);
    } else {
      const newSelected = new Set(selectedPaths);
      if (toggle) {
        if (newSelected.has(path.path_id)) {
          newSelected.delete(path.path_id);
        } else {
          newSelected.add(path.path_id);
        }
      } else {
        newSelected.add(path.path_id);
      }
      setSelectedPaths(newSelected);
      setSelectedPath(path.path_id);
      
      // Visualize all selected paths
      const pathsToShow = filteredPaths.filter(p => newSelected.has(p.path_id));
      if (pathsToShow.length > 0) {
        visualizePaths(pathsToShow);
      }
      onPathSelect?.(path);
    }
  };

  const handleViewModeChange = (mode: 'single' | 'multiple' | 'all') => {
    setViewMode(mode);
    if (mode === 'all') {
      // Show all paths
      setSelectedPaths(new Set(filteredPaths.map(p => p.path_id)));
      visualizePaths(filteredPaths);
    } else if (mode === 'multiple') {
      // Show selected paths
      const pathsToShow = filteredPaths.filter(p => selectedPaths.has(p.path_id));
      if (pathsToShow.length > 0) {
        visualizePaths(pathsToShow);
      }
    } else {
      // Single path mode
      if (selectedPath) {
        const path = filteredPaths.find(p => p.path_id === selectedPath);
        if (path) {
          visualizePath(path);
        }
      }
    }
  };

  if (isLoading) {
    return (
      <div className={`flex items-center justify-center h-full ${className}`}>
        <div className="text-center">
          <RefreshCw className="w-8 h-8 animate-spin mx-auto mb-4 text-red-500" />
          <p className="text-gray-600">Loading attack paths...</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={`flex items-center justify-center h-full ${className}`}>
        <div className="text-center">
          <AlertTriangle className="w-8 h-8 mx-auto mb-4 text-red-500" />
          <p className="text-red-600">{error}</p>
        </div>
      </div>
    );
  }

  return (
    <div className={`w-full h-full relative ${className}`}>
      <ReactFlow
        nodes={nodes}
        edges={edges}
        onNodesChange={onNodesChange}
        onEdgesChange={onEdgesChange}
        nodeTypes={nodeTypes}
        fitView
        attributionPosition="bottom-left"
      >
        <Controls />
        <Background />
        <MiniMap
          nodeColor={(node) => {
            if (node.type === 'entryPoint') {
              return node.data?.type === 'entry' ? '#10b981' : '#ef4444';
            }
            const risk = (node.data?.likelihood || 0) * (node.data?.impact || 0);
            return risk >= 0.6 ? '#ef4444' : risk >= 0.3 ? '#f97316' : risk >= 0.1 ? '#eab308' : '#3b82f6';
          }}
        />

        <Panel position="top-left" className="bg-white p-4 rounded-lg shadow-lg m-4 max-w-sm">
          <div className="space-y-3">
            <div className="flex items-center justify-between mb-3">
              <div className="flex items-center gap-2">
                <Target className="w-5 h-5 text-red-500" />
                <h3 className="font-semibold text-lg">Attack Paths</h3>
              </div>
              <span className="text-xs text-gray-500">{filteredPaths.length} paths</span>
            </div>

            {/* View Mode Toggle */}
            <div className="flex gap-1 p-1 bg-gray-100 rounded-lg">
              <button
                onClick={() => handleViewModeChange('single')}
                className={`flex-1 px-2 py-1 text-xs font-medium rounded transition-colors ${
                  viewMode === 'single' ? 'bg-white shadow-sm text-red-600' : 'text-gray-600'
                }`}
              >
                Single
              </button>
              <button
                onClick={() => handleViewModeChange('multiple')}
                className={`flex-1 px-2 py-1 text-xs font-medium rounded transition-colors ${
                  viewMode === 'multiple' ? 'bg-white shadow-sm text-red-600' : 'text-gray-600'
                }`}
              >
                Multiple
              </button>
              <button
                onClick={() => handleViewModeChange('all')}
                className={`flex-1 px-2 py-1 text-xs font-medium rounded transition-colors ${
                  viewMode === 'all' ? 'bg-white shadow-sm text-red-600' : 'text-gray-600'
                }`}
              >
                All
              </button>
            </div>

            <div className="relative">
              <Search className="absolute left-2 top-2.5 w-4 h-4 text-gray-400" />
              <input
                type="text"
                placeholder="Search paths..."
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="pl-8 pr-3 py-2 border border-gray-300 rounded-md text-sm w-full focus:outline-none focus:ring-2 focus:ring-red-500"
              />
            </div>

            {/* Sort Options */}
            <div>
              <label className="block text-xs font-medium text-gray-700 mb-1">Sort by</label>
              <select
                value={sortBy}
                onChange={(e) => setSortBy(e.target.value as any)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-red-500"
              >
                <option value="criticality">Criticality Score</option>
                <option value="likelihood">Likelihood</option>
                <option value="impact">Impact</option>
                <option value="steps">Number of Steps</option>
              </select>
            </div>

            <div className="space-y-2 max-h-96 overflow-y-auto">
              {filteredPaths.map((path) => {
                const isSelected = viewMode === 'single' 
                  ? selectedPath === path.path_id
                  : selectedPaths.has(path.path_id);
                
                return (
                <div
                  key={path.path_id}
                  onClick={(e) => handlePathSelect(path, e.shiftKey || viewMode !== 'single')}
                  className={`p-3 rounded-lg border-2 cursor-pointer transition-all ${
                    isSelected
                      ? 'border-red-500 bg-red-50'
                      : 'border-gray-200 hover:border-gray-300'
                  }`}
                >
                  <div className="flex items-start justify-between mb-2">
                    <div className="flex-1">
                      <div className="font-semibold text-sm text-gray-900">{path.name}</div>
                      <div className="text-xs text-gray-600 mt-1">
                        {path.steps.length} steps • Score: {(path.criticality_score * 100).toFixed(1)}%
                      </div>
                    </div>
                    <Badge
                      className={
                        path.mitigation_priority === 'high'
                          ? 'bg-red-100 text-red-800'
                          : path.mitigation_priority === 'medium'
                          ? 'bg-yellow-100 text-yellow-800'
                          : 'bg-blue-100 text-blue-800'
                      }
                    >
                      {path.mitigation_priority}
                    </Badge>
                  </div>
                  <div className="flex items-center gap-4 text-xs text-gray-600 mt-2">
                    <div>
                      <span className="font-medium">Likelihood:</span>{' '}
                      {(path.total_likelihood * 100).toFixed(1)}%
                    </div>
                    <div>
                      <span className="font-medium">Impact:</span> {(path.total_impact * 100).toFixed(1)}%
                    </div>
                  </div>
                </div>
                );
              })}
            </div>
          </div>
        </Panel>

        {selectedPath && (
          <Panel position="bottom-right" className="bg-white p-4 rounded-lg shadow-lg m-4 max-w-md">
            <div className="space-y-2">
              <div className="flex items-center justify-between mb-2">
                <h4 className="font-semibold text-gray-900">Path Details</h4>
                <button
                  onClick={() => setSelectedPath(null)}
                  className="text-gray-400 hover:text-gray-600"
                >
                  <XCircle className="w-4 h-4" />
                </button>
              </div>
              {attackPaths
                .find((p) => p.path_id === selectedPath)
                ?.steps.map((step) => (
                  <div key={step.step_number} className="text-sm border-l-2 border-red-500 pl-2 mb-2">
                    <div className="font-medium text-gray-900">
                      Step {step.step_number}: {step.action}
                    </div>
                    <div className="text-xs text-gray-600 mt-1">
                      Target: {step.target_hostname || step.target}
                      {step.target_ip && ` (${step.target_ip})`}
                    </div>
                    <div className="text-xs text-gray-600">
                      Technique: {step.technique}
                      {step.technique_id && ` (${step.technique_id})`}
                    </div>
                    {step.cve_id && (
                      <div className="text-xs text-red-600 mt-1 font-mono">
                        CVE: {step.cve_id}
                      </div>
                    )}
                    {step.proof && (
                      <div className="text-xs text-gray-500 mt-1 italic">
                        {step.proof}
                      </div>
                    )}
                  </div>
                ))}
            </div>
          </Panel>
        )}
      </ReactFlow>
    </div>
  );
};

// Badge component
const Badge: React.FC<{ children: React.ReactNode; className?: string }> = ({
  children,
  className = '',
}) => {
  return (
    <span className={`px-2 py-1 rounded text-xs font-semibold ${className}`}>{children}</span>
  );
};

export default AttackPathVisualizer;

