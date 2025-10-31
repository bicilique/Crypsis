import React, { Fragment } from 'react';
import { Menu, Transition } from '@headlessui/react';
import {
  Bars3Icon,
  BellIcon,
  ShieldCheckIcon,
  UserCircleIcon,
  ArrowRightOnRectangleIcon,
  CogIcon,
} from '@heroicons/react/24/outline';
import { useAuthStore } from '@/stores';
import { cn } from '@/utils';

interface HeaderProps {
  onMenuClick: () => void;
}

/**
 * Header Component
 * 
 * Top navigation bar with menu toggle, logo, and user menu.
 * Includes notifications and user profile dropdown.
 */
export const Header: React.FC<HeaderProps> = ({ onMenuClick }) => {
  const { user, logout } = useAuthStore();
  
  const handleLogout = async () => {
    try {
      await logout();
    } catch (error) {
      console.error('Logout failed:', error);
    }
  };
  
  return (
    <div className="sticky top-0 z-40 flex h-16 shrink-0 items-center gap-x-4 border-b border-gray-200 bg-white px-4 shadow-sm sm:gap-x-6 sm:px-6 lg:px-8">
      {/* Mobile menu button */}
      <button
        type="button"
        className="-m-2.5 p-2.5 text-gray-700 lg:hidden"
        onClick={onMenuClick}
      >
        <span className="sr-only">Open sidebar</span>
        <Bars3Icon className="h-6 w-6" aria-hidden="true" />
      </button>
      
      {/* Separator */}
      <div className="h-6 w-px bg-gray-200 lg:hidden" aria-hidden="true" />
      
      {/* Logo and title for mobile */}
      <div className="flex items-center lg:hidden">
        <ShieldCheckIcon className="h-8 w-8 text-primary-600" />
        <span className="ml-2 text-lg font-semibold text-gray-900">
          Crypsis
        </span>
      </div>
      
      <div className="flex flex-1 gap-x-4 self-stretch lg:gap-x-6">
        {/* Search bar - placeholder for future implementation */}
        <div className="flex flex-1"></div>
        
        <div className="flex items-center gap-x-4 lg:gap-x-6">
          {/* Notifications */}
          <button
            type="button"
            className="-m-2.5 p-2.5 text-gray-400 hover:text-gray-500"
          >
            <span className="sr-only">View notifications</span>
            <BellIcon className="h-6 w-6" aria-hidden="true" />
          </button>
          
          {/* Separator */}
          <div
            className="hidden lg:block lg:h-6 lg:w-px lg:bg-gray-200"
            aria-hidden="true"
          />
          
          {/* Profile dropdown */}
          <Menu as="div" className="relative">
            <Menu.Button className="-m-1.5 flex items-center p-1.5">
              <span className="sr-only">Open user menu</span>
              <UserCircleIcon className="h-8 w-8 text-gray-400" aria-hidden="true" />
              <span className="hidden lg:flex lg:items-center">
                <span className="ml-4 text-sm font-semibold leading-6 text-gray-900">
                  {user?.username}
                </span>
              </span>
            </Menu.Button>
            
            <Transition
              as={Fragment}
              enter="transition ease-out duration-100"
              enterFrom="transform opacity-0 scale-95"
              enterTo="transform opacity-100 scale-100"
              leave="transition ease-in duration-75"
              leaveFrom="transform opacity-100 scale-100"
              leaveTo="transform opacity-0 scale-95"
            >
              <Menu.Items className="absolute right-0 z-10 mt-2.5 w-32 origin-top-right rounded-md bg-white py-2 shadow-lg ring-1 ring-gray-900/5 focus:outline-none">
                <Menu.Item>
                  {({ active }) => (
                    <button
                      className={cn(
                        active ? 'bg-gray-50' : '',
                        'flex w-full items-center px-3 py-1 text-sm leading-6 text-gray-900'
                      )}
                    >
                      <CogIcon className="mr-2 h-4 w-4" />
                      Settings
                    </button>
                  )}
                </Menu.Item>
                <Menu.Item>
                  {({ active }) => (
                    <button
                      onClick={handleLogout}
                      className={cn(
                        active ? 'bg-gray-50' : '',
                        'flex w-full items-center px-3 py-1 text-sm leading-6 text-gray-900'
                      )}
                    >
                      <ArrowRightOnRectangleIcon className="mr-2 h-4 w-4" />
                      Sign out
                    </button>
                  )}
                </Menu.Item>
              </Menu.Items>
            </Transition>
          </Menu>
        </div>
      </div>
    </div>
  );
};
