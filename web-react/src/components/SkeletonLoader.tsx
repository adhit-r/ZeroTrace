/**
 * Reusable skeleton loader components with shimmer animation
 * Provides better perceived performance than spinners
 */
import React from 'react';
import { cn } from '@/lib/utils';

interface SkeletonProps {
  className?: string;
  variant?: 'text' | 'circular' | 'rectangular';
  width?: string | number;
  height?: string | number;
  animation?: 'pulse' | 'wave' | 'none';
}

/**
 * Base skeleton component with shimmer effect
 */
export const Skeleton: React.FC<SkeletonProps> = ({
  className,
  variant = 'rectangular',
  width,
  height,
  animation = 'wave',
}) => {
  const baseClasses = 'bg-gray-200 dark:bg-gray-700';
  
  const variantClasses = {
    text: 'rounded',
    circular: 'rounded-full',
    rectangular: 'rounded-lg',
  };

  const animationClasses = {
    pulse: 'animate-pulse',
    wave: 'animate-shimmer',
    none: '',
  };

  return (
    <div
      className={cn(
        baseClasses,
        variantClasses[variant],
        animationClasses[animation],
        className
      )}
      style={{
        width: width || '100%',
        height: height || '1rem',
      }}
      aria-label="Loading..."
      role="status"
    />
  );
};

/**
 * Skeleton for table rows
 */
interface TableSkeletonProps {
  rows?: number;
  columns?: number;
  className?: string;
}

export const TableSkeleton: React.FC<TableSkeletonProps> = ({
  rows = 5,
  columns = 4,
  className,
}) => {
  return (
    <div className={cn('space-y-2', className)}>
      {Array.from({ length: rows }).map((_, rowIndex) => (
        <div key={rowIndex} className="flex gap-4">
          {Array.from({ length: columns }).map((_, colIndex) => (
            <Skeleton
              key={colIndex}
              variant="text"
              height="2rem"
              className="flex-1"
            />
          ))}
        </div>
      ))}
    </div>
  );
};

/**
 * Skeleton for cards
 */
interface CardSkeletonProps {
  count?: number;
  showImage?: boolean;
  showActions?: boolean;
  className?: string;
}

export const CardSkeleton: React.FC<CardSkeletonProps> = ({
  count = 1,
  showImage = false,
  showActions = false,
  className,
}) => {
  return (
    <div className={cn('space-y-4', className)}>
      {Array.from({ length: count }).map((_, index) => (
        <div
          key={index}
          className="border border-gray-200 dark:border-gray-700 rounded-lg p-4 space-y-3"
        >
          {showImage && (
            <Skeleton variant="rectangular" height="200px" width="100%" />
          )}
          <div className="space-y-2">
            <Skeleton variant="text" height="1.5rem" width="60%" />
            <Skeleton variant="text" height="1rem" width="100%" />
            <Skeleton variant="text" height="1rem" width="80%" />
          </div>
          {showActions && (
            <div className="flex gap-2">
              <Skeleton variant="rectangular" height="2.5rem" width="100px" />
              <Skeleton variant="rectangular" height="2.5rem" width="100px" />
            </div>
          )}
        </div>
      ))}
    </div>
  );
};

/**
 * Skeleton for list items
 */
interface ListSkeletonProps {
  items?: number;
  showAvatar?: boolean;
  className?: string;
}

export const ListSkeleton: React.FC<ListSkeletonProps> = ({
  items = 5,
  showAvatar = false,
  className,
}) => {
  return (
    <div className={cn('space-y-3', className)}>
      {Array.from({ length: items }).map((_, index) => (
        <div key={index} className="flex items-center gap-3">
          {showAvatar && (
            <Skeleton variant="circular" width="40px" height="40px" />
          )}
          <div className="flex-1 space-y-2">
            <Skeleton variant="text" height="1rem" width="70%" />
            <Skeleton variant="text" height="0.875rem" width="50%" />
          </div>
        </div>
      ))}
    </div>
  );
};

/**
 * Skeleton for dashboard stats
 */
interface StatsSkeletonProps {
  count?: number;
  className?: string;
}

export const StatsSkeleton: React.FC<StatsSkeletonProps> = ({
  count = 4,
  className,
}) => {
  return (
    <div className={cn('grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4', className)}>
      {Array.from({ length: count }).map((_, index) => (
        <div
          key={index}
          className="border border-gray-200 dark:border-gray-700 rounded-lg p-4 space-y-2"
        >
          <Skeleton variant="text" height="0.875rem" width="60%" />
          <Skeleton variant="text" height="2rem" width="40%" />
          <Skeleton variant="text" height="0.75rem" width="80%" />
        </div>
      ))}
    </div>
  );
};

// Add shimmer animation to global CSS (will be added to index.css)
export const shimmerKeyframes = `
@keyframes shimmer {
  0% {
    background-position: -1000px 0;
  }
  100% {
    background-position: 1000px 0;
  }
}

.animate-shimmer {
  background: linear-gradient(
    90deg,
    #f0f0f0 0%,
    #e0e0e0 20%,
    #f0f0f0 40%,
    #f0f0f0 100%
  );
  background-size: 1000px 100%;
  animation: shimmer 2s infinite;
}

.dark .animate-shimmer {
  background: linear-gradient(
    90deg,
    #374151 0%,
    #4b5563 20%,
    #374151 40%,
    #374151 100%
  );
  background-size: 1000px 100%;
}
`;

export default Skeleton;

