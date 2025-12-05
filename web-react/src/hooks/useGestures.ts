/**
 * Touch gesture hooks using @use-gesture/react
 */
import { useGesture } from '@use-gesture/react';
import { useRef, useState } from 'react';

/**
 * Swipe gesture hook
 */
export function useSwipe(
  onSwipeLeft?: () => void,
  onSwipeRight?: () => void,
  onSwipeUp?: () => void,
  onSwipeDown?: () => void
) {
  const [swipeDirection, setSwipeDirection] = useState<string | null>(null);

  const bind = useGesture({
    onDragEnd: ({ direction, velocity }) => {
      const [x, y] = direction;
      const [vx, vy] = velocity;

      // Minimum velocity threshold
      if (Math.abs(vx) > 0.5 || Math.abs(vy) > 0.5) {
        if (Math.abs(x) > Math.abs(y)) {
          // Horizontal swipe
          if (x > 0) {
            setSwipeDirection('right');
            onSwipeRight?.();
          } else {
            setSwipeDirection('left');
            onSwipeLeft?.();
          }
        } else {
          // Vertical swipe
          if (y > 0) {
            setSwipeDirection('down');
            onSwipeDown?.();
          } else {
            setSwipeDirection('up');
            onSwipeUp?.();
          }
        }
      }
    },
  });

  return { bind, swipeDirection };
}

/**
 * Pull to refresh hook
 */
export function usePullToRefresh(onRefresh: () => Promise<void>) {
  const [isRefreshing, setIsRefreshing] = useState(false);
  const [pullDistance, setPullDistance] = useState(0);

  const bind = useGesture({
    onDrag: ({ movement: [, my], dragging }) => {
      if (dragging && my > 0) {
        setPullDistance(Math.min(my, 100));
      }
    },
    onDragEnd: async ({ movement: [, my] }) => {
      if (my > 80 && !isRefreshing) {
        setIsRefreshing(true);
        await onRefresh();
        setIsRefreshing(false);
      }
      setPullDistance(0);
    },
  });

  return { bind, isRefreshing, pullDistance };
}

/**
 * Pinch to zoom hook
 */
export function usePinchZoom(
  onZoom?: (scale: number) => void,
  minScale: number = 0.5,
  maxScale: number = 3
) {
  const [scale, setScale] = useState(1);

  const bind = useGesture({
    onPinch: ({ offset: [scaleValue] }) => {
      const clampedScale = Math.max(minScale, Math.min(maxScale, scaleValue));
      setScale(clampedScale);
      onZoom?.(clampedScale);
    },
  });

  return { bind, scale };
}

