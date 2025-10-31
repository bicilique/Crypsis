import React from 'react';
import { cn } from '@/utils';
import { CardProps } from '@/types';

/**
 * Reusable Card Component
 * 
 * A container component that provides consistent styling for content sections.
 * Supports title, subtitle, actions, and optional hover effects.
 */
export const Card: React.FC<CardProps> = ({
  title,
  subtitle,
  action,
  children,
  className,
  hover = false,
  onClick,
}) => {
  const baseClasses = 'bg-white border border-gray-200 rounded-lg shadow-sm';
  const hoverClasses = hover ? 'hover:shadow-md transition-shadow cursor-pointer' : '';
  
  return (
    <div
      className={cn(baseClasses, hoverClasses, className)}
      onClick={onClick}
    >
      {(title || subtitle || action) && (
        <div className="px-6 py-4 border-b border-gray-200">
          <div className="flex items-center justify-between">
            <div>
              {title && (
                <h3 className="text-lg font-medium text-gray-900">
                  {title}
                </h3>
              )}
              {subtitle && (
                <p className="mt-1 text-sm text-gray-500">
                  {subtitle}
                </p>
              )}
            </div>
            {action && (
              <div className="flex-shrink-0">
                {action}
              </div>
            )}
          </div>
        </div>
      )}
      
      <div className="px-6 py-4">
        {children}
      </div>
    </div>
  );
};
