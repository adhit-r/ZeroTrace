export const mockLineChartData = {
  data: {
    labels: ['January', 'February', 'March', 'April', 'May', 'June', 'July'],
    datasets: [
      {
        label: 'Vulnerabilities Found',
        data: [65, 59, 80, 81, 56, 55, 40],
        fill: false,
        backgroundColor: 'rgb(255, 99, 132)',
        borderColor: 'rgba(255, 99, 132, 0.2)',
      },
    ],
  },
  options: {
    responsive: true,
    maintainAspectRatio: false,
  },
};

export const mockDoughnutChartData = {
  data: {
    labels: ['Critical', 'High', 'Medium', 'Low'],
    datasets: [
      {
        label: '# of Vulnerabilities',
        data: [12, 19, 3, 5],
        backgroundColor: [
          'rgba(255, 99, 132, 0.8)',
          'rgba(255, 159, 64, 0.8)',
          'rgba(255, 205, 86, 0.8)',
          'rgba(75, 192, 192, 0.8)',
        ],
        borderColor: [
          'rgba(0, 0, 0, 1)',
          'rgba(0, 0, 0, 1)',
          'rgba(0, 0, 0, 1)',
          'rgba(0, 0, 0, 1)',
        ],
        borderWidth: 3,
      },
    ],
  },
  options: {
    responsive: true,
    maintainAspectRatio: false,
  },
};

export const mockBarChartData = {
    data: {
      labels: ['Week 1', 'Week 2', 'Week 3', 'Week 4'],
      datasets: [
        {
          label: 'Patched',
          data: [20, 25, 15, 30],
          backgroundColor: 'rgba(75, 192, 192, 0.8)',
          borderColor: 'rgba(0, 0, 0, 1)',
          borderWidth: 3,
        },
        {
          label: 'Outstanding',
          data: [10, 5, 15, 10],
          backgroundColor: 'rgba(255, 99, 132, 0.8)',
          borderColor: 'rgba(0, 0, 0, 1)',
          borderWidth: 3,
        },
      ],
    },
    options: {
      responsive: true,
      maintainAspectRatio: false,
      scales: {
        y: {
          beginAtZero: true,
        },
      },
    },
  };

