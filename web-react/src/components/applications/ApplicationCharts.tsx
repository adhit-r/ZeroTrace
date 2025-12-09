import React from 'react';
import { Bar, Pie } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  BarElement,
  ArcElement,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
} from 'chart.js';
import type { Application, ApplicationStats } from '@/services/applicationService';

ChartJS.register(
  CategoryScale,
  LinearScale,
  BarElement,
  ArcElement,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

interface ApplicationChartsProps {
  applications: Application[];
  stats: ApplicationStats;
}

const ApplicationCharts: React.FC<ApplicationChartsProps> = ({ applications, stats }) => {
  // Bar chart: Applications by risk level
  const riskLevelChartData = {
    labels: ['Critical', 'High', 'Medium', 'Low', 'Safe'],
    datasets: [
      {
        label: 'Applications',
        data: [
          stats.byRiskLevel.critical,
          stats.byRiskLevel.high,
          stats.byRiskLevel.medium,
          stats.byRiskLevel.low,
          stats.byRiskLevel.safe,
        ],
        backgroundColor: [
          'rgba(239, 68, 68, 0.8)',
          'rgba(249, 115, 22, 0.8)',
          'rgba(234, 179, 8, 0.8)',
          'rgba(59, 130, 246, 0.8)',
          'rgba(34, 197, 94, 0.8)',
        ],
        borderColor: [
          'rgb(239, 68, 68)',
          'rgb(249, 115, 22)',
          'rgb(234, 179, 8)',
          'rgb(59, 130, 246)',
          'rgb(34, 197, 94)',
        ],
        borderWidth: 3,
      },
    ],
  };

  // Pie chart: Application distribution by classification
  const classificationLabels = Object.keys(stats.byClassification);
  const classificationChartData = {
    labels: classificationLabels,
    datasets: [
      {
        label: 'Applications',
        data: classificationLabels.map(key => stats.byClassification[key]),
        backgroundColor: [
          'rgba(249, 115, 22, 0.8)',
          'rgba(59, 130, 246, 0.8)',
          'rgba(34, 197, 94, 0.8)',
          'rgba(234, 179, 8, 0.8)',
          'rgba(168, 85, 247, 0.8)',
          'rgba(236, 72, 153, 0.8)',
        ],
        borderColor: [
          'rgb(249, 115, 22)',
          'rgb(59, 130, 246)',
          'rgb(34, 197, 94)',
          'rgb(234, 179, 8)',
          'rgb(168, 85, 247)',
          'rgb(236, 72, 153)',
        ],
        borderWidth: 3,
      },
    ],
  };

  // Stacked bar: Vulnerable vs Safe apps per agent
  const agentIds = Object.keys(stats.byAgent);
  const agentNames = agentIds.map(id => {
    const app = applications.find(a => a.agentId === id);
    return app?.agentName || id.slice(0, 8);
  });
  
  const vulnerableByAgent = agentIds.map(id => 
    applications.filter(a => a.agentId === id && a.status === 'vulnerable').length
  );
  const safeByAgent = agentIds.map(id => 
    applications.filter(a => a.agentId === id && a.status === 'safe').length
  );

  const agentChartData = {
    labels: agentNames,
    datasets: [
      {
        label: 'Vulnerable',
        data: vulnerableByAgent,
        backgroundColor: 'rgba(239, 68, 68, 0.8)',
        borderColor: 'rgb(239, 68, 68)',
        borderWidth: 3,
      },
      {
        label: 'Safe',
        data: safeByAgent,
        backgroundColor: 'rgba(34, 197, 94, 0.8)',
        borderColor: 'rgb(34, 197, 94)',
        borderWidth: 3,
      },
    ],
  };

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'top' as const,
        labels: {
          font: {
            weight: 'bold' as const,
            size: 12,
          },
        },
      },
      title: {
        display: true,
        font: {
          weight: 'bold' as const,
          size: 16,
        },
      },
    },
    scales: {
      y: {
        beginAtZero: true,
        ticks: {
          font: {
            weight: 'bold' as const,
          },
        },
      },
      x: {
        ticks: {
          font: {
            weight: 'bold' as const,
          },
        },
      },
    },
  };

  const pieOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'right' as const,
        labels: {
          font: {
            weight: 'bold' as const,
            size: 12,
          },
        },
      },
      title: {
        display: true,
        font: {
          weight: 'bold' as const,
          size: 16,
        },
      },
    },
  };

  return (
    <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
      {/* Applications by Risk Level */}
      <div className="bg-white border-4 border-black rounded-lg p-6">
        <h3 className="text-xl font-black text-black uppercase mb-4">Applications by Risk Level</h3>
        <div className="h-64">
          <Bar data={riskLevelChartData} options={chartOptions} />
        </div>
      </div>

      {/* Application Distribution by Classification */}
      <div className="bg-white border-4 border-black rounded-lg p-6">
        <h3 className="text-xl font-black text-black uppercase mb-4">Distribution by Classification</h3>
        <div className="h-64">
          <Pie data={classificationChartData} options={pieOptions} />
        </div>
      </div>

      {/* Vulnerable vs Safe Apps per Agent */}
      <div className="bg-white border-4 border-black rounded-lg p-6 lg:col-span-2">
        <h3 className="text-xl font-black text-black uppercase mb-4">Vulnerable vs Safe Apps per Agent</h3>
        <div className="h-64">
          <Bar 
            data={agentChartData} 
            options={{
              ...chartOptions,
              scales: {
                ...chartOptions.scales,
                x: {
                  ...chartOptions.scales.x,
                  stacked: true,
                },
                y: {
                  ...chartOptions.scales.y,
                  stacked: true,
                },
              },
            }} 
          />
        </div>
      </div>
    </div>
  );
};

export default ApplicationCharts;

