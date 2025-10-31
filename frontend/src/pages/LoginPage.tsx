import React, { useState } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
import { ShieldCheckIcon, EyeIcon, EyeSlashIcon } from '@heroicons/react/24/outline';
import { Button, Input, Alert } from '@/components/ui';
import { useAuthStore } from '@/stores';
import { LoginCredentials } from '@/types';
import { ROUTES } from '@/constants';

// Validation schema
const loginSchema = z.object({
  username: z.string().min(1, 'Username is required'),
  password: z.string().min(1, 'Password is required'),
  rememberMe: z.boolean().optional(),
});

/**
 * Login Page Component
 * 
 * Provides authentication interface with form validation and error handling.
 * Redirects to dashboard after successful login.
 */
export const LoginPage: React.FC = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const { login, isLoading, error, clearError } = useAuthStore();
  
  const [showPassword, setShowPassword] = useState(false);
  
  const from = (location.state as any)?.from?.pathname || ROUTES.DASHBOARD;
  
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginCredentials>({resolver: zodResolver(loginSchema),
      defaultValues: {
        username: '',
        password: '',
        rememberMe: false,
      },
  });
  
  const onSubmit = async (data: LoginCredentials) => {
    try {
      clearError();
      await login(data);
      navigate(from, { replace: true });
    } catch (error) {
      // Error is handled by the store
      console.error('Login failed:', error);
    }
  };
  
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        {/* Header */}
        <div className="text-center">
          <div className="mx-auto h-16 w-16 flex items-center justify-center rounded-full bg-primary-100">
            <ShieldCheckIcon className="h-10 w-10 text-primary-600" />
          </div>
          <h2 className="mt-6 text-3xl font-bold text-gray-900">
            Crypsis Admin
          </h2>
          <p className="mt-2 text-sm text-gray-600">
            Secure file encryption and storage management
          </p>
        </div>
        
        {/* Login Form */}
        <div className="bg-white py-8 px-6 shadow-lg rounded-lg">
          <form className="space-y-6" onSubmit={handleSubmit(onSubmit)}>
            {error && (
              <Alert
                type="error"
                message={error}
                dismissible
                onDismiss={clearError}
              />
            )}
            
            <Input
              label="Username"
              type="text"
              autoComplete="username"
              {...register('username')}
              error={errors.username?.message}
              placeholder="Enter your username"
              className="h-14 text-lg"
            />
            
            <Input
              label="Password"
              type={showPassword ? 'text' : 'password'}
              autoComplete="current-password"
              {...register('password')}
              error={errors.password?.message}
              placeholder="Enter your password"
              className="h-14 text-lg"
              rightIcon={
                <button
                  type="button"
                  onClick={() => setShowPassword(!showPassword)}
                  className="text-gray-400 hover:text-gray-600"
                  aria-label={showPassword ? 'Hide password' : 'Show password'}
                >
                  {showPassword ? (
                    <EyeSlashIcon className="h-5 w-5" />
                  ) : (
                    <EyeIcon className="h-5 w-5" />
                  )}
                </button>
              }
            />
            
            {/* <div className="flex items-center">
              <input
                id="remember-me"
                type="checkbox"
                {...register('rememberMe')}
                className="h-4 w-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500"
              />
              <label htmlFor="remember-me" className="ml-2 block text-sm text-gray-900">
                Remember me
              </label>
            </div> */}
            
            <Button
              type="submit"
              variant="primary"
              size="lg"
              loading={isLoading}
              className="w-full"
            >
              Sign in
            </Button>
          </form>
        </div>
        
        {/* Security Notice */}
        <div className="text-center">
          <p className="text-xs text-gray-500">
            This is a secure area. All activities are monitored and logged.
          </p>
        </div>
      </div>
    </div>
  );
};
