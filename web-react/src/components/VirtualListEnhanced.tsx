/**
 * Enhanced Virtual List using @tanstack/react-virtual
 * Better performance and more features than custom implementation
 */
import React, { useRef } from 'react';
import { useVirtualizer } from '@tanstack/react-virtual';
import { cn } from '@/lib/utils';

interface VirtualListEnhancedProps<T> {
  items: T[];
  renderItem: (item: T, index: number, virtualItem: any) => React.ReactNode;
  containerHeight?: number | string;
  itemHeight?: number | ((index: number) => number);
  overscan?: number;
  className?: string;
  containerClassName?: string;
  estimateSize?: number;
  horizontal?: boolean;
  getScrollElement?: () => HTMLElement | null;
}

export function VirtualListEnhanced<T>({
  items,
  renderItem,
  containerHeight = '100%',
  itemHeight = 50,
  overscan = 5,
  className = '',
  containerClassName = '',
  estimateSize,
  horizontal = false,
  getScrollElement,
}: VirtualListEnhancedProps<T>) {
  const parentRef = useRef<HTMLDivElement>(null);

  const virtualizer = useVirtualizer({
    count: items.length,
    getScrollElement: getScrollElement || (() => parentRef.current),
    estimateSize: estimateSize || (typeof itemHeight === 'number' ? itemHeight : undefined),
    overscan,
    horizontal,
    ...(typeof itemHeight === 'function' && { getItemSize: itemHeight }),
  });

  return (
    <div
      ref={parentRef}
      className={cn('overflow-auto', containerClassName)}
      style={{
        height: typeof containerHeight === 'number' ? `${containerHeight}px` : containerHeight,
        width: horizontal ? '100%' : undefined,
      }}
    >
      <div
        style={{
          height: horizontal ? '100%' : `${virtualizer.getTotalSize()}px`,
          width: horizontal ? `${virtualizer.getTotalSize()}px` : '100%',
          position: 'relative',
        }}
      >
        {virtualizer.getVirtualItems().map((virtualItem) => (
          <div
            key={virtualItem.key}
            data-index={virtualItem.index}
            ref={virtualizer.measureElement}
            style={{
              position: 'absolute',
              top: 0,
              left: 0,
              width: horizontal ? `${virtualItem.size}px` : '100%',
              height: horizontal ? '100%' : `${virtualItem.size}px`,
              transform: horizontal
                ? `translateX(${virtualItem.start}px)`
                : `translateY(${virtualItem.start}px)`,
            }}
            className={className}
          >
            {renderItem(items[virtualItem.index], virtualItem.index, virtualItem)}
          </div>
        ))}
      </div>
    </div>
  );
}

export default VirtualListEnhanced;

