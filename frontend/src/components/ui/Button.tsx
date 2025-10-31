import { cn } from '@/utils';
import { ButtonProps } from '@/types';

/**
 * Reusable Button Component
 * 
 * A flexible button component that supports different variants, sizes, and states.
 * Includes loading state with spinner and proper accessibility attributes.
 */
export const Button: React.FC<ButtonProps> = ({
  variant = 'primary',
  size = 'md',
  loading = false,
  disabled = false,
  icon,
  children,
  onClick,
  type = 'button',
  className,
  ...props
}) => {
  const baseClasses = 'inline-flex items-center justify-center rounded-md font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed';
  
  const variantClasses = {
    primary: 'bg-primary-600 hover:bg-primary-700 focus:ring-primary-500 text-white',
    secondary: 'bg-gray-200 hover:bg-gray-300 focus:ring-gray-500 text-gray-900',
    outline: 'border border-gray-300 hover:bg-gray-50 focus:ring-primary-500 text-gray-700',
    ghost: 'hover:bg-gray-100 focus:ring-gray-500 text-gray-700',
    danger: 'bg-red-600 hover:bg-red-700 focus:ring-red-500 text-white',
  };
  
  const sizeClasses = {
    xs: 'px-2.5 py-1.5 text-xs',
    sm: 'px-3 py-2 text-sm',
    md: 'px-4 py-2 text-sm',
    lg: 'px-4 py-2 text-base',
    xl: 'px-6 py-3 text-base',
  };
  
  const isDisabled = disabled || loading;
  
  return (
    <button
      type={type}
      onClick={onClick}
      disabled={isDisabled}
      className={cn(
        baseClasses,
        variantClasses[variant],
        sizeClasses[size],
        className
      )}
      aria-disabled={isDisabled}
      {...props}
    >
      {loading && (
        <div className="w-4 h-4 mr-2 border-2 border-current border-t-transparent rounded-full animate-spin" />
      )}
      {icon && !loading && (
        <span className={cn('flex-shrink-0', children ? 'mr-2' : '')}>
          {icon}
        </span>
      )}
      {children}
    </button>
  );
};
