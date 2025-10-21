import React, { useRef, useEffect } from 'react';
import * as d3 from 'd3';

export interface HeatmapData {
  id: string;
  name: string;
  lat: number;
  lng: number;
  riskScore: number;
  criticalVulns: number;
  totalAssets: number;
  complianceScore: number;
  lastScan: string;
}

interface RiskHeatmapProps {
  data: HeatmapData[];
  selectedBranchId: string | null | undefined;
  onBranchSelect: (branchId: string | null) => void;
  className?: string;
}

const RiskHeatmap: React.FC<RiskHeatmapProps> = ({ data, selectedBranchId, onBranchSelect }) => {
  const svgRef = useRef<SVGSVGElement | null>(null);

  useEffect(() => {
    if (!svgRef.current) return;

    const width = 500;
    const height = 300;

    const svg = d3.select(svgRef.current)
      .attr('viewBox', `0 0 ${width} ${height}`);

    // Create a simple background
    svg.append('rect')
      .attr('width', width)
      .attr('height', height)
      .style('fill', '#34495e')
      .style('stroke', '#2c3e50')
      .style('stroke-width', 2);

    const getRiskColor = (riskScore: number): string => {
      if (riskScore >= 80) return '#e74c3c'; // red
      if (riskScore >= 60) return '#e67e22'; // orange
      if (riskScore >= 40) return '#f1c40f'; // yellow
      return '#2ecc71'; // green
    };

    const markers = svg.selectAll('.marker')
      .data(data, (d: any) => d.id);

    markers.enter().append('circle')
      .attr('class', 'marker')
      .attr('cx', (d: any) => width * 0.1 + (d.lng + 180) / 360 * width * 0.8)
      .attr('cy', (d: any) => height * 0.1 + (90 - d.lat) / 180 * height * 0.8)
      .attr('r', 8)
      .style('fill', d => getRiskColor(d.riskScore))
      .style('stroke', 'white')
      .style('stroke-width', 1.5)
      .style('cursor', 'pointer')
            .on('click', (_, d: any) => onBranchSelect(d.id))
            .on('mouseover', function () {
              d3.select(this).transition().duration(200).attr('r', 12);
            })
      .on('mouseout', function () {
        d3.select(this).transition().duration(200).attr('r', 8);
      });

    markers.exit().remove();

    // Update selection styles
    svg.selectAll('.marker')
      .style('stroke', (d: any) => d.id === selectedBranchId ? 'cyan' : 'white')
      .style('stroke-width', (d: any) => d.id === selectedBranchId ? 3 : 1.5);


  }, [data, selectedBranchId, onBranchSelect]);

  return <svg ref={svgRef} className="w-full h-full"></svg>;
};

export default RiskHeatmap;
