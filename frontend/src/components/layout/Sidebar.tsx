import React, { Fragment } from 'react';
import { NavLink } from 'react-router-dom';
import { Dialog, Transition } from '@headlessui/react';
import {
  XMarkIcon,
  ShieldCheckIcon,
  FolderIcon,
  HomeIcon,
  UsersIcon,
  CubeIcon,
  ClipboardDocumentListIcon,
  ExclamationTriangleIcon,
} from '@heroicons/react/24/outline';
import { cn } from '@/utils';
import { ROUTES } from '@/constants';
import { NavItem } from '@/types';

interface SidebarProps {
  open: boolean;
  onClose: () => void;
}

// Navigation items
const navigation: NavItem[] = [
  {
    name: 'Dashboard',
    href: ROUTES.DASHBOARD,
    icon: HomeIcon,
  },
  {
    name: 'Files',
    href: ROUTES.FILES,
    icon: FolderIcon,
  },
  {
    name: 'Admin Panel',
    href: ROUTES.ADMIN,
    icon: UsersIcon,
    children: [
      {
        name: 'Users',
        href: ROUTES.ADMIN_USERS,
        icon: UsersIcon,
      },
      {
        name: 'Applications',
        href: ROUTES.ADMIN_APPS,
        icon: CubeIcon,
      },
      {
        name: 'Audit Logs',
        href: ROUTES.ADMIN_LOGS,
        icon: ClipboardDocumentListIcon,
      },
    ],
  },
  {
    name: 'Security',
    href: ROUTES.SECURITY,
    icon: ExclamationTriangleIcon,
  },
];

/**
 * Sidebar Component
 * 
 * Navigation sidebar with responsive mobile overlay.
 * Displays navigation items with icons and active states.
 */
export const Sidebar: React.FC<SidebarProps> = ({ open, onClose }) => {
  const SidebarContent = () => (
    <div className="flex grow flex-col gap-y-5 overflow-y-auto bg-white px-6 pb-4">
      {/* Logo */}
      <div className="flex h-16 shrink-0 items-center">
        <ShieldCheckIcon className="h-8 w-8 text-primary-600" />
        <span className="ml-2 text-xl font-bold text-gray-900">
          Crypsis
        </span>
      </div>
      
      {/* Navigation */}
      <nav className="flex flex-1 flex-col">
        <ul role="list" className="flex flex-1 flex-col gap-y-7">
          <li>
            <ul role="list" className="-mx-2 space-y-1">
              {navigation.map((item) => (
                <li key={item.name}>
                  <NavLink
                    to={item.href}
                    className={({ isActive }) =>
                      cn(
                        isActive
                          ? 'bg-primary-50 text-primary-700'
                          : 'text-gray-700 hover:text-primary-700 hover:bg-gray-50',
                        'group flex gap-x-3 rounded-md p-2 text-sm leading-6 font-semibold'
                      )
                    }
                    onClick={onClose}
                  >
                    <item.icon
                      className="h-6 w-6 shrink-0"
                      aria-hidden="true"
                    />
                    {item.name}
                    {item.badge && (
                      <span className="ml-auto inline-block px-2 py-0.5 text-xs font-medium bg-primary-100 text-primary-700 rounded-full">
                        {item.badge}
                      </span>
                    )}
                  </NavLink>
                  
                  {/* Sub-navigation */}
                  {item.children && (
                    <ul className="mt-1 px-2">
                      {item.children.map((subItem) => (
                        <li key={subItem.name}>
                          <NavLink
                            to={subItem.href}
                            className={({ isActive }) =>
                              cn(
                                isActive
                                  ? 'bg-primary-50 text-primary-700'
                                  : 'text-gray-600 hover:text-primary-700 hover:bg-gray-50',
                                'group flex gap-x-3 rounded-md p-2 pl-8 text-sm leading-6 font-medium'
                              )
                            }
                            onClick={onClose}
                          >
                            <subItem.icon
                              className="h-5 w-5 shrink-0"
                              aria-hidden="true"
                            />
                            {subItem.name}
                          </NavLink>
                        </li>
                      ))}
                    </ul>
                  )}
                </li>
              ))}
            </ul>
          </li>
          
          {/* System status */}
          <li className="mt-auto">
            <div className="rounded-md bg-green-50 p-3">
              <div className="flex">
                <div className="flex-shrink-0">
                  <div className="h-2 w-2 bg-green-400 rounded-full"></div>
                </div>
                <div className="ml-3">
                  <p className="text-sm font-medium text-green-800">
                    System Online
                  </p>
                  <p className="text-xs text-green-600">
                    All services operational
                  </p>
                </div>
              </div>
            </div>
          </li>
        </ul>
      </nav>
    </div>
  );
  
  return (
    <>
      {/* Mobile sidebar */}
      <Transition.Root show={open} as={Fragment}>
        <Dialog as="div" className="relative z-50 lg:hidden" onClose={onClose}>
          <Transition.Child
            as={Fragment}
            enter="transition-opacity ease-linear duration-300"
            enterFrom="opacity-0"
            enterTo="opacity-100"
            leave="transition-opacity ease-linear duration-300"
            leaveFrom="opacity-100"
            leaveTo="opacity-0"
          >
            <div className="fixed inset-0 bg-gray-900/80" />
          </Transition.Child>
          
          <div className="fixed inset-0 flex">
            <Transition.Child
              as={Fragment}
              enter="transition ease-in-out duration-300 transform"
              enterFrom="-translate-x-full"
              enterTo="translate-x-0"
              leave="transition ease-in-out duration-300 transform"
              leaveFrom="translate-x-0"
              leaveTo="-translate-x-full"
            >
              <Dialog.Panel className="relative mr-16 flex w-full max-w-xs flex-1">
                <Transition.Child
                  as={Fragment}
                  enter="ease-in-out duration-300"
                  enterFrom="opacity-0"
                  enterTo="opacity-100"
                  leave="ease-in-out duration-300"
                  leaveFrom="opacity-100"
                  leaveTo="opacity-0"
                >
                  <div className="absolute left-full top-0 flex w-16 justify-center pt-5">
                    <button
                      type="button"
                      className="-m-2.5 p-2.5"
                      onClick={onClose}
                    >
                      <span className="sr-only">Close sidebar</span>
                      <XMarkIcon
                        className="h-6 w-6 text-white"
                        aria-hidden="true"
                      />
                    </button>
                  </div>
                </Transition.Child>
                
                <SidebarContent />
              </Dialog.Panel>
            </Transition.Child>
          </div>
        </Dialog>
      </Transition.Root>
      
      {/* Desktop sidebar */}
      <div className="hidden lg:fixed lg:inset-y-0 lg:z-50 lg:flex lg:w-64 lg:flex-col">
        <div className="flex grow flex-col gap-y-5 overflow-y-auto border-r border-gray-200 bg-white">
          <SidebarContent />
        </div>
      </div>
    </>
  );
};
