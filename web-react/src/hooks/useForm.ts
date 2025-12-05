/**
 * Form utilities with react-hook-form and zod validation
 * Provides type-safe form handling with reduced re-renders
 */
import { useForm as useReactHookForm, UseFormProps, UseFormReturn } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';

/**
 * Typed form hook with Zod validation
 */
export function useForm<T extends z.ZodType<any, any>>(
  schema: T,
  options?: Omit<UseFormProps<z.infer<T>>, 'resolver'>
): UseFormReturn<z.infer<T>> {
  return useReactHookForm<z.infer<T>>({
    resolver: zodResolver(schema),
    mode: 'onChange', // Validate on change for better UX
    ...options,
  });
}

/**
 * Common form validation schemas
 */
export const formSchemas = {
  email: z.string().email('Invalid email address'),
  password: z.string().min(8, 'Password must be at least 8 characters'),
  url: z.string().url('Invalid URL'),
  number: z.coerce.number(),
  positiveNumber: z.coerce.number().positive('Must be a positive number'),
  nonEmptyString: z.string().min(1, 'This field is required'),
  optionalString: z.string().optional(),
  date: z.coerce.date(),
  uuid: z.string().uuid('Invalid UUID format'),
};

/**
 * Helper to create form schemas
 */
export function createFormSchema<T extends z.ZodRawShape>(shape: T) {
  return z.object(shape);
}

/**
 * Example usage:
 * 
 * const schema = createFormSchema({
 *   name: formSchemas.nonEmptyString,
 *   email: formSchemas.email,
 *   age: formSchemas.positiveNumber,
 * });
 * 
 * const form = useForm(schema);
 */

