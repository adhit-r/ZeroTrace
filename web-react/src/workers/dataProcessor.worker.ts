/**
 * Web Worker for heavy data processing
 * Offloads computations from main thread
 */

// Handle messages from main thread
self.addEventListener('message', (event: MessageEvent) => {
  const { type, data, id } = event.data;

  try {
    let result;

    switch (type) {
      case 'PROCESS_CHART_DATA':
        result = processChartData(data);
        break;

      case 'FILTER_DATA':
        result = filterData(data);
        break;

      case 'PARSE_CSV':
        result = parseCSV(data);
        break;

      case 'PARSE_JSON':
        result = parseJSON(data);
        break;

      case 'TRANSFORM_DATA':
        result = transformData(data);
        break;

      case 'AGGREGATE_DATA':
        result = aggregateData(data);
        break;

      default:
        throw new Error(`Unknown message type: ${type}`);
    }

    // Send result back to main thread
    self.postMessage({ id, type, result, success: true });
  } catch (error) {
    // Send error back to main thread
    self.postMessage({
      id,
      type,
      error: error instanceof Error ? error.message : 'Unknown error',
      success: false,
    });
  }
});

/**
 * Process chart data for visualization
 */
function processChartData(data: {
  rawData: any[];
  options?: {
    groupBy?: string;
    aggregate?: 'sum' | 'avg' | 'count' | 'max' | 'min';
    dateFormat?: string;
  };
}): any {
  const { rawData, options = {} } = data;

  if (!rawData || rawData.length === 0) {
    return { labels: [], datasets: [] };
  }

  // Group data if needed
  let processedData = rawData;
  if (options.groupBy) {
    const grouped = rawData.reduce((acc, item) => {
      const key = item[options.groupBy!];
      if (!acc[key]) {
        acc[key] = [];
      }
      acc[key].push(item);
      return acc;
    }, {} as Record<string, any[]>);

    processedData = Object.entries(grouped).map(([key, values]) => ({
      [options.groupBy!]: key,
      values,
    }));
  }

  // Aggregate if needed
  if (options.aggregate) {
    processedData = processedData.map((item) => {
      const values = item.values || [item];
      let aggregated: number;

      switch (options.aggregate) {
        case 'sum':
          aggregated = values.reduce((sum: number, v: any) => sum + (v.value || 0), 0);
          break;
        case 'avg':
          aggregated =
            values.reduce((sum: number, v: any) => sum + (v.value || 0), 0) / values.length;
          break;
        case 'count':
          aggregated = values.length;
          break;
        case 'max':
          aggregated = Math.max(...values.map((v: any) => v.value || 0));
          break;
        case 'min':
          aggregated = Math.min(...values.map((v: any) => v.value || 0));
          break;
        default:
          aggregated = 0;
      }

      return { ...item, aggregated };
    });
  }

  return processedData;
}

/**
 * Filter large datasets
 */
function filterData(data: {
  items: any[];
  filters: Record<string, any>;
}): any[] {
  const { items, filters } = data;

  return items.filter((item) => {
    return Object.entries(filters).every(([key, value]) => {
      if (value === undefined || value === null || value === '') {
        return true;
      }

      const itemValue = item[key];

      if (typeof value === 'string') {
        return String(itemValue).toLowerCase().includes(String(value).toLowerCase());
      }

      if (Array.isArray(value)) {
        return value.includes(itemValue);
      }

      return itemValue === value;
    });
  });
}

/**
 * Parse CSV data
 */
function parseCSV(data: { csv: string; delimiter?: string }): any[] {
  const { csv, delimiter = ',' } = data;
  const lines = csv.split('\n');
  const headers = lines[0].split(delimiter).map((h) => h.trim());

  return lines.slice(1).map((line) => {
    const values = line.split(delimiter).map((v) => v.trim());
    const obj: Record<string, string> = {};
    headers.forEach((header, index) => {
      obj[header] = values[index] || '';
    });
    return obj;
  });
}

/**
 * Parse JSON data
 */
function parseJSON(data: { json: string }): any {
  return JSON.parse(data.json);
}

/**
 * Transform data structure
 */
function transformData(data: {
  items: any[];
  transform: (item: any) => any;
}): any[] {
  const { items, transform } = data;
  return items.map(transform);
}

/**
 * Aggregate data
 */
function aggregateData(data: {
  items: any[];
  groupBy: string;
  aggregations: Record<string, 'sum' | 'avg' | 'count' | 'max' | 'min'>;
}): any[] {
  const { items, groupBy, aggregations } = data;

  const grouped = items.reduce((acc, item) => {
    const key = item[groupBy];
    if (!acc[key]) {
      acc[key] = [];
    }
    acc[key].push(item);
    return acc;
  }, {} as Record<string, any[]>);

  return (Object.entries(grouped) as [string, any[]][]).map(([key, groupItems]) => {
    const result: any = { [groupBy]: key };

    Object.entries(aggregations).forEach(([field, operation]) => {
      const values = groupItems.map((item: any) => item[field]).filter((v: any) => v != null);

      switch (operation) {
        case 'sum':
          result[field] = values.reduce((sum: number, v: any) => sum + Number(v), 0);
          break;
        case 'avg':
          result[field] = values.length > 0 ? values.reduce((sum: number, v: any) => sum + Number(v), 0) / values.length : 0;
          break;
        case 'count':
          result[field] = values.length;
          break;
        case 'max':
          result[field] = values.length > 0 ? Math.max(...values.map((v: any) => Number(v))) : 0;
          break;
        case 'min':
          result[field] = values.length > 0 ? Math.min(...values.map((v: any) => Number(v))) : 0;
          break;
      }
    });

    return result;
  });
}

// Export types for main thread
export type WorkerMessage =
  | { type: 'PROCESS_CHART_DATA'; data: { rawData: any[]; options?: any }; id: string }
  | { type: 'FILTER_DATA'; data: { items: any[]; filters: Record<string, any> }; id: string }
  | { type: 'PARSE_CSV'; data: { csv: string; delimiter?: string }; id: string }
  | { type: 'PARSE_JSON'; data: { json: string }; id: string }
  | { type: 'TRANSFORM_DATA'; data: { items: any[]; transform: (item: any) => any }; id: string }
  | { type: 'AGGREGATE_DATA'; data: { items: any[]; groupBy: string; aggregations: Record<string, string> }; id: string };

