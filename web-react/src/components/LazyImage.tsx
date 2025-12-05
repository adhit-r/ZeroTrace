/**
 * Optimized LazyImage component with Intersection Observer
 * Supports WebP/AVIF with fallbacks, responsive images, and blur placeholders
 */
import React, { useState, useRef, useEffect, ImgHTMLAttributes } from 'react';

interface LazyImageProps extends Omit<ImgHTMLAttributes<HTMLImageElement>, 'src' | 'srcSet'> {
  src: string;
  srcSet?: string;
  alt: string;
  placeholder?: string; // Base64 blur placeholder
  webp?: string; // WebP version
  avif?: string; // AVIF version
  sizes?: string; // Responsive image sizes
  fallback?: string; // Fallback image
  className?: string;
  onLoad?: () => void;
  onError?: () => void;
}

export const LazyImage: React.FC<LazyImageProps> = ({
  src,
  srcSet,
  alt,
  placeholder,
  webp,
  avif,
  sizes,
  fallback,
  className = '',
  onLoad,
  onError,
  ...props
}) => {
  const [isLoaded, setIsLoaded] = useState(false);
  const [isInView, setIsInView] = useState(false);
  const [hasError, setHasError] = useState(false);
  const imgRef = useRef<HTMLImageElement>(null);
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            setIsInView(true);
            observer.disconnect();
          }
        });
      },
      {
        rootMargin: '50px', // Start loading 50px before image enters viewport
        threshold: 0.01,
      }
    );

    if (containerRef.current) {
      observer.observe(containerRef.current);
    }

    return () => {
      observer.disconnect();
    };
  }, []);

  const handleLoad = () => {
    setIsLoaded(true);
    onLoad?.();
  };

  const handleError = () => {
    setHasError(true);
    onError?.();
  };

  const imageSrc = hasError && fallback ? fallback : src;

  return (
    <div
      ref={containerRef}
      className={`relative overflow-hidden ${className}`}
      style={{ aspectRatio: props.width && props.height ? `${props.width}/${props.height}` : undefined }}
    >
      {/* Blur placeholder */}
      {placeholder && !isLoaded && (
        <div
          className="absolute inset-0 bg-cover bg-center filter blur-sm scale-110"
          style={{
            backgroundImage: `url(${placeholder})`,
            opacity: isLoaded ? 0 : 1,
            transition: 'opacity 0.3s ease-in-out',
          }}
          aria-hidden="true"
        />
      )}

      {/* Actual image */}
      {isInView && (
        <picture>
          {/* AVIF format (best compression) */}
          {avif && (
            <source srcSet={avif} type="image/avif" sizes={sizes} />
          )}
          {/* WebP format (good compression) */}
          {webp && (
            <source srcSet={webp} type="image/webp" sizes={sizes} />
          )}
          {/* Fallback to original */}
          <img
            ref={imgRef}
            src={imageSrc}
            srcSet={srcSet}
            sizes={sizes}
            alt={alt}
            className={`w-full h-full object-cover transition-opacity duration-300 ${
              isLoaded ? 'opacity-100' : 'opacity-0'
            }`}
            onLoad={handleLoad}
            onError={handleError}
            loading="lazy"
            decoding="async"
            {...props}
          />
        </picture>
      )}

      {/* Loading state */}
      {!isInView && (
        <div className="absolute inset-0 bg-gray-200 animate-pulse" aria-hidden="true" />
      )}
    </div>
  );
};

export default LazyImage;

