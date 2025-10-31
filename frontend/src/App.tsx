import { useEffect } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { LoginPage } from '@/pages/LoginPage';
import { DashboardPage } from '@/pages/DashboardPage';
import { FilesPage } from '@/pages/FilesPage';
import { AdminsPage } from '@/pages/AdminsPage';
import { ApplicationsPage } from '@/pages/ApplicationsPage';
import { LogsPage } from '@/pages/LogsPage';
import { SecurityPage } from '@/pages/SecurityPage';
import { AppLayout, ProtectedRoute } from '@/components/layout';
import { initializeAuth } from '@/stores';
import { ROUTES } from '@/constants';

/**
 * Main App Component
 * 
 * Sets up routing and authentication initialization.
 * Provides the main application structure with protected routes.
 */
function App() {
  useEffect(() => {
    // Initialize authentication state on app start
    initializeAuth();
  }, []);
  
  return (
    <Router>
      <Routes>
        {/* Public routes */}
        <Route path={ROUTES.LOGIN} element={<LoginPage />} />
        
        {/* Protected routes */}
        <Route
          path="/*"
          element={
            <ProtectedRoute>
              <AppLayout />
            </ProtectedRoute>
          }
        >
          {/* Dashboard */}
          <Route index element={<Navigate to={ROUTES.DASHBOARD} replace />} />
          <Route path={ROUTES.DASHBOARD.slice(1)} element={<DashboardPage />} />
          
          {/* Files Management */}
          <Route path={ROUTES.FILES.slice(1)} element={<FilesPage />} />
          
          {/* Admin Management */}
          <Route path={ROUTES.ADMIN.slice(1)} element={<AdminsPage />} />
          <Route path={ROUTES.ADMIN_USERS.slice(1)} element={<AdminsPage />} />
          <Route path={ROUTES.ADMIN_APPS.slice(1)} element={<ApplicationsPage />} />
          <Route path={ROUTES.ADMIN_LOGS.slice(1)} element={<LogsPage />} />
          
          {/* Security */}
          <Route path={ROUTES.SECURITY.slice(1)} element={<SecurityPage />} />
          
          {/* Catch all - redirect to dashboard */}
          <Route path="*" element={<Navigate to={ROUTES.DASHBOARD} replace />} />
        </Route>
      </Routes>
    </Router>
  );
}

export default App;
