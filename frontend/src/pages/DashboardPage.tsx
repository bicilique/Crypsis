import React, { useEffect } from 'react';
import {
  FolderIcon,
  UsersIcon,
  CubeIcon,
  CloudArrowUpIcon,
  ShieldCheckIcon,
  ExclamationTriangleIcon,
} from '@heroicons/react/24/outline';
import { Card, StatusBadge, LoadingSpinner } from '@/components/ui';
import { useAdminStore, useFileStore, useApplicationStore } from '@/stores';
import { formatFileSize } from '@/utils';

/**
 * Dashboard Page Component
 * 
 * Overview page showing system statistics and recent activity.
 * Displays key metrics for files, users, applications, and security status.
 */
export const DashboardPage: React.FC = () => {
  const { admins, fetchAdmins, isLoading: adminsLoading } = useAdminStore();
  const { files, fetchFiles, isLoading: filesLoading } = useFileStore();
  const { applications, fetchApplications, isLoading: appsLoading } = useApplicationStore();
  
  const isLoading = adminsLoading || filesLoading || appsLoading;
  
  useEffect(() => {
    fetchAdmins();
    fetchFiles();
    fetchApplications();
  }, [fetchAdmins, fetchFiles, fetchApplications]);
  
  if (isLoading && files.length === 0) {
    return (
      <div className="flex items-center justify-center h-64">
        <LoadingSpinner size="lg" text="Loading dashboard..." />
      </div>
    );
  }
  
  // Calculate statistics from available data
  const totalFiles = files.length;
  const totalUsers = admins.length;
  const totalApplications = applications.length;
  const storageUsed = files.reduce((sum, file) => sum + (file.file_size || 0), 0);
  
  const stats = [
    {
      name: 'Total Files',
      value: totalFiles,
      icon: FolderIcon,
      change: '+4.75%',
      changeType: 'positive' as const,
    },
    {
      name: 'Active Users',
      value: totalUsers,
      icon: UsersIcon,
      change: '+54.02%',
      changeType: 'positive' as const,
    },
    {
      name: 'Applications',
      value: totalApplications,
      icon: CubeIcon,
      change: '-1.39%',
      changeType: 'negative' as const,
    },
    {
      name: 'Storage Used',
      value: formatFileSize(storageUsed),
      icon: CloudArrowUpIcon,
      change: '+10.18%',
      changeType: 'positive' as const,
    },
  ];
  
  return (
    <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      {/* Page header */}
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
        <p className="mt-2 text-gray-600">
          Welcome to Crypsis secure file management system
        </p>
      </div>
      
      {/* Stats */}
      <div className="grid grid-cols-1 gap-6 sm:grid-cols-2 lg:grid-cols-4 mb-8">
        {stats.map((stat) => (
          <Card key={stat.name} className="bg-white">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <stat.icon className="h-6 w-6 text-gray-400" />
              </div>
              <div className="ml-4 w-0 flex-1">
                <dl>
                  <dt className="text-sm font-medium text-gray-500 truncate">
                    {stat.name}
                  </dt>
                  <dd className="flex items-baseline">
                    <div className="text-2xl font-semibold text-gray-900">
                      {stat.value}
                    </div>
                    <div
                      className={`ml-2 flex items-baseline text-sm font-semibold ${
                        stat.changeType === 'positive'
                          ? 'text-green-600'
                          : 'text-red-600'
                      }`}
                    >
                      {stat.change}
                    </div>
                  </dd>
                </dl>
              </div>
            </div>
          </Card>
        ))}
      </div>
      
      {/* Content grid */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Security Status */}
        <Card
          title="Security Status"
          subtitle="Current system security overview"
        >
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center">
                <ShieldCheckIcon className="h-5 w-5 text-green-500 mr-2" />
                <span className="text-sm text-gray-700">Encryption Status</span>
              </div>
              <StatusBadge status="success" text="Active" />
            </div>
            
            <div className="flex items-center justify-between">
              <div className="flex items-center">
                <ShieldCheckIcon className="h-5 w-5 text-green-500 mr-2" />
                <span className="text-sm text-gray-700">KMS Integration</span>
              </div>
              <StatusBadge status="success" text="Connected" />
            </div>
            
            <div className="flex items-center justify-between">
              <div className="flex items-center">
                <ExclamationTriangleIcon className="h-5 w-5 text-yellow-500 mr-2" />
                <span className="text-sm text-gray-700">Security Alerts</span>
              </div>
              <StatusBadge status="warning" text="2 Pending" />
            </div>
          </div>
        </Card>
        
        {/* System Health */}
        <Card
          title="System Health"
          subtitle="Current system performance metrics"
        >
          <div className="space-y-4">
            <div>
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm text-gray-700">CPU Usage</span>
                <span className="text-sm font-medium text-gray-900">45%</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div className="bg-blue-600 h-2 rounded-full" style={{ width: '45%' }}></div>
              </div>
            </div>
            
            <div>
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm text-gray-700">Memory Usage</span>
                <span className="text-sm font-medium text-gray-900">62%</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div className="bg-green-600 h-2 rounded-full" style={{ width: '62%' }}></div>
              </div>
            </div>
            
            <div>
              <div className="flex items-center justify-between mb-2">
                <span className="text-sm text-gray-700">Storage Usage</span>
                <span className="text-sm font-medium text-gray-900">78%</span>
              </div>
              <div className="w-full bg-gray-200 rounded-full h-2">
                <div className="bg-yellow-600 h-2 rounded-full" style={{ width: '78%' }}></div>
              </div>
            </div>
          </div>
        </Card>
        
        {/* Recent Activity */}
        <Card
          title="Recent Activity"
          subtitle="Latest system activities"
          className="lg:col-span-2"
        >
          <div className="flow-root">
            <ul className="-mb-8">
              {[
                {
                  id: 1,
                  type: 'upload',
                  content: 'New file uploaded: document.pdf',
                  timestamp: '2 minutes ago',
                  user: 'admin',
                },
                {
                  id: 2,
                  type: 'user',
                  content: 'New admin user created: john.doe',
                  timestamp: '1 hour ago',
                  user: 'admin',
                },
                {
                  id: 3,
                  type: 'security',
                  content: 'Security key rotated successfully',
                  timestamp: '3 hours ago',
                  user: 'system',
                },
                {
                  id: 4,
                  type: 'application',
                  content: 'OAuth2 application registered: MyApp',
                  timestamp: '1 day ago',
                  user: 'admin',
                },
              ].map((item, itemIdx, items) => (
                <li key={item.id}>
                  <div className="relative pb-8">
                    {itemIdx !== items.length - 1 ? (
                      <span
                        className="absolute left-4 top-4 -ml-px h-full w-0.5 bg-gray-200"
                        aria-hidden="true"
                      />
                    ) : null}
                    <div className="relative flex space-x-3">
                      <div>
                        <span className="h-8 w-8 rounded-full bg-primary-500 flex items-center justify-center ring-8 ring-white">
                          <FolderIcon className="h-4 w-4 text-white" aria-hidden="true" />
                        </span>
                      </div>
                      <div className="flex min-w-0 flex-1 justify-between space-x-4 pt-1.5">
                        <div>
                          <p className="text-sm text-gray-500">
                            {item.content}{' '}
                            <span className="font-medium text-gray-900">
                              by {item.user}
                            </span>
                          </p>
                        </div>
                        <div className="whitespace-nowrap text-right text-sm text-gray-500">
                          {item.timestamp}
                        </div>
                      </div>
                    </div>
                  </div>
                </li>
              ))}
            </ul>
          </div>
        </Card>
      </div>
    </div>
  );
};
