import React from 'react';
import { cn } from '@/utils';

interface LoadingSpinnerProps {
  size?: 'sm' | 'md' | 'lg' | 'xl';
  className?: string;
  text?: string;
}

/**
 * Loading Spinner Component
 * 
 * A reusable loading indicator with different sizes and optional text.
 * Used throughout the app for loading states.
 */
export const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({
  size = 'md',
  className,
  text,
}) => {
  const sizeClasses = {
    sm: 'w-4 h-4',
    md: 'w-6 h-6',
    lg: 'w-8 h-8',
    xl: 'w-12 h-12',
    xxl: 'w-16 h-16',
  };
  
  return (
    <div className={cn('flex items-center justify-center', className)}>
      <div className="flex flex-col items-center">
        <div
          className={cn(
            'border-2 border-gray-200 border-t-primary-600 rounded-full animate-spin',
            sizeClasses[size]
          )}
        />
        {text && (
          <p className="mt-2 text-sm text-gray-600">
            {text}
          </p>
        )}
      </div>
    </div>
  );
};
