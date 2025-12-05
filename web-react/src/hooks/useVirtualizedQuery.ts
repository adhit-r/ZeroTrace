/**
 * React Query hook optimized for virtualized lists
 * Fetches data in pages for efficient rendering
 */
import { useInfiniteQuery, UseInfiniteQueryOptions } from '@tanstack/react-query';

interface VirtualizedQueryOptions<TData, TPageParam = number> {
  queryKey: string[];
  queryFn: (pageParam: TPageParam) => Promise<{ data: TData[]; nextCursor?: TPageParam }>;
  getNextPageParam?: (lastPage: { data: TData[]; nextCursor?: TPageParam }) => TPageParam | undefined;
  pageSize?: number;
  enabled?: boolean;
}

export function useVirtualizedQuery<TData, TPageParam = number>(
  options: VirtualizedQueryOptions<TData, TPageParam>
) {
  return useInfiniteQuery({
    queryKey: options.queryKey,
    queryFn: ({ pageParam = 0 as TPageParam }) => options.queryFn(pageParam),
    getNextPageParam: options.getNextPageParam || ((lastPage) => lastPage.nextCursor),
    initialPageParam: 0 as TPageParam,
    enabled: options.enabled !== false,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

