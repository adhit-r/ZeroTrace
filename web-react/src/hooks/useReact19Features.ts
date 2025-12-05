/**
 * React 19 feature hooks
 * useOptimistic, useActionState, useFormStatus
 */

import { useOptimistic, useActionState, useFormStatus } from 'react';

/**
 * Wrapper for useOptimistic with better TypeScript support
 */
export function useOptimisticUpdate<T>(
  currentState: T,
  updateFn: (currentState: T, optimisticValue: T) => T
) {
  return useOptimistic(currentState, updateFn);
}

/**
 * Hook for form actions with React 19 useActionState
 * Provides pending state and form action handling
 */
export function useFormAction<TFormData, TActionResult>(
  action: (prevState: TActionResult | null, formData: FormData) => Promise<TActionResult>,
  initialState: TActionResult | null = null
) {
  return useActionState(action, initialState);
}

/**
 * Hook to get form status from a form context
 * Must be used within a <form> element
 */
export function useFormPending() {
  const { pending } = useFormStatus();
  return pending;
}

/**
 * Example usage component for React 19 form features
 */
export function FormStatusButton({ children, ...props }: React.ButtonHTMLAttributes<HTMLButtonElement>) {
  const { pending } = useFormStatus();
  
  return (
    <button
      {...props}
      disabled={pending || props.disabled}
      aria-busy={pending}
    >
      {pending ? 'Loading...' : children}
    </button>
  );
}

