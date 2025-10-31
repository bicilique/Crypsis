import { useState, useEffect } from 'react';
import { Outlet } from 'react-router-dom';
import {Header} from './Header';
import { Sidebar } from './Sidebar';

/**
 * App Layout Component
 * 
 * Main layout structure with sidebar navigation and header.
 * Provides responsive design with collapsible sidebar on mobile.
 */
export const AppLayout = () => {
  const [sidebarOpen, setSidebarOpen] = useState(false);

  // Prevent background scroll when sidebar is open (mobile)
  useEffect(() => {
    if (sidebarOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }
    return () => {
      document.body.style.overflow = '';
    };
  }, [sidebarOpen]);

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Sidebar overlay for mobile */}
      <div
        className={`fixed inset-0 z-40 bg-black bg-opacity-50 transition-opacity lg:hidden ${sidebarOpen ? 'block' : 'hidden'}`}
        aria-hidden={!sidebarOpen}
        onClick={() => setSidebarOpen(false)}
      />
      {/* Sidebar */}
      <Sidebar
        open={sidebarOpen}
        onClose={() => setSidebarOpen(false)}
        aria-label="Sidebar navigation"
      />
      {/* Main content */}
      <div className="lg:pl-64 flex flex-col flex-1">
        {/* Header */}
        <Header onMenuClick={() => setSidebarOpen(true)} />
        {/* Page content */}
        <main className="flex-1">
          <div className="py-6">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  );
};
