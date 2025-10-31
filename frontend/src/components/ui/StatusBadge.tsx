import React from 'react';
import { cn } from '@/utils';
import { StatusBadgeProps } from '@/types';
import { STATUS_COLORS } from '@/constants';

/**
 * Status Badge Component
 * 
 * Displays status information with appropriate colors and optional icons.
 * Used for file status, user status, system alerts, etc.
 */
export const StatusBadge: React.FC<StatusBadgeProps> = ({
  status,
  text,
  icon,
  size = 'md',
}) => {
  const baseClasses = 'inline-flex items-center border rounded-full font-medium';
  
  const sizeClasses = {
    sm: 'px-2 py-1 text-xs',
    md: 'px-2.5 py-0.5 text-sm',
    lg: 'px-3 py-1 text-sm',
  };
  
  return (
    <span
      className={cn(
        baseClasses,
        sizeClasses[size],
        STATUS_COLORS[status]
      )}
    >
      {icon && (
        <span className={cn('flex-shrink-0', text ? 'mr-1' : '')}>
          {icon}
        </span>
      )}
      {text}
    </span>
  );
};
