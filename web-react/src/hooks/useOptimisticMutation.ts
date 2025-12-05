/**
 * React Query hook for optimistic updates
 * Updates UI immediately, rolls back on error
 */
import { useMutation, useQueryClient, UseMutationOptions } from '@tanstack/react-query';

interface OptimisticMutationOptions<TData, TVariables, TContext> {
  mutationFn: (variables: TVariables) => Promise<TData>;
  onMutate?: (variables: TVariables) => Promise<TContext> | TContext;
  onError?: (error: Error, variables: TVariables, context: TContext) => void;
  onSuccess?: (data: TData, variables: TVariables, context: TContext) => void;
  invalidateQueries?: string[];
}

export function useOptimisticMutation<TData, TVariables, TContext = unknown>(
  options: OptimisticMutationOptions<TData, TVariables, TContext>
) {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: options.mutationFn,
    onMutate: async (variables) => {
      // Cancel outgoing refetches
      await queryClient.cancelQueries();

      // Snapshot previous value
      const previousData = queryClient.getQueryData(options.invalidateQueries?.[0] || []);

      // Optimistically update
      const context = options.onMutate ? await options.onMutate(variables) : ({} as TContext);

      return { previousData, context } as TContext;
    },
    onError: (error, variables, context: any) => {
      // Rollback on error
      if (context?.previousData) {
        queryClient.setQueryData(
          options.invalidateQueries?.[0] || [],
          context.previousData
        );
      }
      options.onError?.(error, variables, context?.context || ({} as TContext));
    },
    onSuccess: (data, variables, context: any) => {
      // Invalidate queries to refetch fresh data
      if (options.invalidateQueries) {
        options.invalidateQueries.forEach((queryKey) => {
          queryClient.invalidateQueries({ queryKey: [queryKey] });
        });
      }
      options.onSuccess?.(data, variables, context?.context || ({} as TContext));
    },
  });
}

