import React, { useState } from 'react';
import { Outlet, Link, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { 
  Home, 
  Server, 
  AlertTriangle, 
  Target,
  Activity,
  Settings, 
  LogOut, 
  Menu,
  X,
  Network,
  Shield,
  FileText,
  Users,
  Search,
  Bell,
  Terminal
} from 'lucide-react';

const Layout: React.FC = () => {
  const { user, logout } = useAuth();
  const location = useLocation();
  const [sidebarOpen, setSidebarOpen] = useState(false);

  // Terminal-inspired navigation structure
  const navigation = [
    { 
      name: 'DASHBOARD', 
      href: '/', 
      icon: Home,
      description: 'Security overview and metrics'
    },
    { 
      name: 'AGENTS', 
      href: '/agents', 
      icon: Server,
      description: 'Active agent monitoring',
      badge: '12'
    },
    { 
      name: 'VULNERABILITIES', 
      href: '/vulnerabilities', 
      icon: AlertTriangle,
      description: 'Software vulnerability management',
      badge: '89'
    },
    { 
      name: 'SCANS', 
      href: '/scans', 
      icon: Activity,
      description: 'Scan management and history'
    },
    { 
      name: 'TOPOLOGY', 
      href: '/topology', 
      icon: Network,
      description: 'Network asset visualization'
    },
    { 
      name: 'REPORTS', 
      href: '/reports', 
      icon: FileText,
      description: 'Security reports and analytics'
    },
    { 
      name: 'SETTINGS', 
      href: '/settings', 
      icon: Settings,
      description: 'System configuration'
    },
  ];

  return (
    <div className="min-h-screen">
      {/* Mobile sidebar */}
      <div className={`fixed inset-0 z-50 lg:hidden ${sidebarOpen ? 'block' : 'hidden'}`}>
        <div className="fixed inset-0 bg-black bg-opacity-80" onClick={() => setSidebarOpen(false)} />
        <div className="fixed inset-y-0 left-0 flex w-64 flex-col sidebar">
          <div className="flex h-16 items-center justify-between px-4 border-b border-light-gray">
            <div className="flex items-center">
              <Terminal className="h-8 w-8 text-gold mr-3" />
              <h1 className="text-xl font-semibold text-gold">VULNDETECT</h1>
            </div>
            <button
              onClick={() => setSidebarOpen(false)}
              className="text-text-muted hover:text-gold"
            >
              <X className="h-6 w-6" />
            </button>
          </div>
          <nav className="flex-1 space-y-1 px-2 py-4 overflow-y-auto">
            {navigation.map((item) => {
              const isActive = location.pathname === item.href;
              return (
                <Link
                  key={item.name}
                  to={item.href}
                  onClick={() => setSidebarOpen(false)}
                  className={`menu-item ${isActive ? 'active' : ''}`}
                >
                  <div className="flex items-center">
                    <item.icon className="mr-3 h-5 w-5" />
                    <div className="flex-1">
                      <div className="flex items-center justify-between">
                        <span>{item.name}</span>
                        {item.badge && (
                          <span className="badge badge-info">
                            {item.badge}
                          </span>
                        )}
                      </div>
                      <p className="text-xs mt-1 text-text-secondary">
                        {item.description}
                      </p>
                    </div>
                  </div>
                </Link>
              );
            })}
          </nav>
          <div className="border-t border-light-gray p-4">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <div className="h-8 w-8 rounded-full bg-medium-gray flex items-center justify-center border border-gold">
                  <Users className="h-4 w-4 text-gold" />
                </div>
              </div>
              <div className="ml-3">
                <p className="text-sm font-medium text-text-primary">{user?.name || 'ADMIN'}</p>
                <p className="text-xs text-text-secondary">{user?.email || 'admin@vulndetect.com'}</p>
              </div>
            </div>
            <button
              onClick={logout}
              className="mt-3 w-full btn btn-ghost btn-sm"
            >
              <LogOut className="mr-3 h-4 w-4" />
              SIGN OUT
            </button>
          </div>
        </div>
      </div>

      {/* Desktop sidebar */}
      <div className="hidden lg:fixed lg:inset-y-0 lg:flex lg:w-64 lg:flex-col">
        <div className="sidebar">
          {/* Logo */}
          <div className="flex h-16 items-center px-4 border-b border-light-gray">
            <div className="flex items-center">
              <Terminal className="h-8 w-8 text-gold mr-3" />
              <h1 className="text-xl font-semibold text-gold">VULNDETECT</h1>
            </div>
          </div>

          {/* Navigation */}
          <nav className="flex-1 space-y-1 px-2 py-4 overflow-y-auto">
            {navigation.map((item) => {
              const isActive = location.pathname === item.href;
              return (
                <Link
                  key={item.name}
                  to={item.href}
                  className={`menu-item ${isActive ? 'active' : ''}`}
                >
                  <div className="flex items-center">
                    <item.icon className="mr-3 h-5 w-5" />
                    <div className="flex-1">
                      <div className="flex items-center justify-between">
                        <span>{item.name}</span>
                        {item.badge && (
                          <span className="badge badge-info">
                            {item.badge}
                          </span>
                        )}
                      </div>
                      <p className="text-xs mt-1 text-text-secondary">
                        {item.description}
                      </p>
                    </div>
                  </div>
                </Link>
              );
            })}
          </nav>

          {/* User section */}
          <div className="border-t border-light-gray p-4">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <div className="h-8 w-8 rounded-full bg-medium-gray flex items-center justify-center border border-gold">
                  <Users className="h-4 w-4 text-gold" />
                </div>
              </div>
              <div className="ml-3">
                <p className="text-sm font-medium text-text-primary">{user?.name || 'ADMIN'}</p>
                <p className="text-xs text-text-secondary">{user?.email || 'admin@vulndetect.com'}</p>
              </div>
            </div>
            <button
              onClick={logout}
              className="mt-3 w-full btn btn-ghost btn-sm"
            >
              <LogOut className="mr-3 h-4 w-4" />
              SIGN OUT
            </button>
          </div>
        </div>
      </div>

      {/* Main content */}
      <div className="main-content">
        {/* Top bar */}
        <div className="header">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-x-4">
              <button
                type="button"
                className="btn btn-ghost btn-sm lg:hidden"
                onClick={() => setSidebarOpen(true)}
              >
                <Menu className="h-6 w-6" />
              </button>

              {/* Search bar */}
              <div className="relative flex flex-1 max-w-md">
                <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                  <Search className="h-5 w-5 text-text-muted" />
                </div>
                <input
                  type="text"
                  placeholder="SEARCH ASSETS, VULNERABILITIES..."
                  className="input pl-10"
                />
              </div>
            </div>

            {/* Right side actions */}
            <div className="flex items-center gap-x-4">
              {/* Notifications */}
              <button className="btn btn-ghost btn-sm relative">
                <Bell className="h-6 w-6" />
                <span className="absolute top-1 right-1 h-2 w-2 bg-critical rounded-full"></span>
              </button>

              {/* Divider */}
              <div className="hidden lg:block lg:h-6 lg:w-px lg:bg-light-gray" />

              {/* User info */}
              <div className="flex items-center gap-x-4">
                <span className="text-sm text-text-primary">
                  WELCOME, {user?.name || 'ADMIN'}
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* Page content */}
        <main className="py-6">
          <div className="container">
            <Outlet />
          </div>
        </main>
      </div>
    </div>
  );
};

export default Layout;
