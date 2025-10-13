import React, { useState } from 'react';
import { Outlet, Link, useLocation } from 'react-router-dom';
import { useUser, UserButton, useClerk } from '@clerk/clerk-react';
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
  Terminal,
  Building2,
  MapPin
} from 'lucide-react';

const Layout: React.FC = () => {
  const { user } = useUser();
  const { signOut } = useClerk();
  const location = useLocation();
  const [sidebarOpen, setSidebarOpen] = useState(false);

  // Navigation structure matching ZeroTrace design system
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
      description: 'Active agent monitoring'
    },
    { 
      name: 'VULNERABILITIES', 
      href: '/vulnerabilities', 
      icon: AlertTriangle,
      description: 'Software vulnerability management'
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
    <div className="min-h-screen bg-gray-50">
      {/* Mobile sidebar */}
      <div className={`fixed inset-0 z-50 lg:hidden ${sidebarOpen ? 'block' : 'hidden'}`}>
        <div className="fixed inset-0 bg-black bg-opacity-50" onClick={() => setSidebarOpen(false)} />
        <div className="fixed inset-y-0 left-0 flex w-64 flex-col bg-white border-r-3 border-black shadow-xl">
          <div className="flex h-16 items-center justify-between px-4 border-b-3 border-black bg-orange-50">
            <div className="flex items-center">
              <Terminal className="h-8 w-8 text-orange-600 mr-3" />
              <h1 className="text-xl font-bold text-black uppercase tracking-wider">ZeroTrace</h1>
            </div>
            <button
              onClick={() => setSidebarOpen(false)}
              className="p-2 hover:bg-gray-200 rounded border-2 border-black"
            >
              <X className="h-5 w-5" />
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
                  className={`block p-3 rounded border-2 border-black transition-all duration-150 ${
                    isActive 
                      ? 'bg-orange-100 text-orange-800 border-orange-300 shadow-[4px_4px_0px_0px_rgba(255,107,0,1)]' 
                      : 'bg-white text-black hover:bg-gray-50 hover:shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]'
                  }`}
                >
                  <div className="flex items-center">
                    <item.icon className="mr-3 h-5 w-5" />
                    <div className="flex-1">
                      <div className="flex items-center justify-between">
                        <span className="font-bold text-sm uppercase tracking-wider">{item.name}</span>
                        {item.badge && (
                          <span className="px-2 py-1 bg-orange-100 text-orange-800 border-2 border-orange-300 rounded text-xs font-bold uppercase tracking-wider">
                            {item.badge}
                          </span>
                        )}
                      </div>
                      <p className="text-xs mt-1 text-gray-600">
                        {item.description}
                      </p>
                    </div>
                  </div>
                </Link>
              );
            })}
          </nav>
          <div className="border-t-3 border-black p-4 bg-gray-50">
            <div className="flex items-center mb-3">
              <div className="flex-shrink-0">
                <div className="h-8 w-8 rounded-full bg-orange-100 border-2 border-orange-300 flex items-center justify-center">
                  <Users className="h-4 w-4 text-orange-800" />
                </div>
              </div>
              <div className="ml-3">
                <p className="text-sm font-bold text-black">{user?.fullName || user?.username || 'USER'}</p>
                <p className="text-xs text-gray-600">{user?.primaryEmailAddress?.emailAddress || 'user@zerotrace.com'}</p>
              </div>
            </div>
            <button
              onClick={() => signOut()}
              className="w-full p-2 bg-red-100 text-red-800 border-2 border-red-300 rounded text-sm font-bold uppercase tracking-wider hover:bg-red-200 transition-colors"
            >
              <LogOut className="mr-2 h-4 w-4 inline-block" />
              SIGN OUT
            </button>
          </div>
        </div>
      </div>

      {/* Desktop sidebar */}
      <div className="hidden lg:fixed lg:inset-y-0 lg:flex lg:w-64 lg:flex-col">
        <div className="flex flex-col w-full bg-white border-r-3 border-black shadow-xl">
          {/* Logo */}
          <div className="flex h-16 items-center px-4 border-b-3 border-black bg-orange-50">
            <div className="flex items-center">
              <Terminal className="h-8 w-8 text-orange-600 mr-3" />
              <h1 className="text-xl font-bold text-black uppercase tracking-wider">ZeroTrace</h1>
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
                  className={`block p-3 rounded border-2 border-black transition-all duration-150 ${
                    isActive 
                      ? 'bg-orange-100 text-orange-800 border-orange-300 shadow-[4px_4px_0px_0px_rgba(255,107,0,1)]' 
                      : 'bg-white text-black hover:bg-gray-50 hover:shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]'
                  }`}
                >
                  <div className="flex items-center">
                    <item.icon className="mr-3 h-5 w-5" />
                    <div className="flex-1">
                      <div className="flex items-center justify-between">
                        <span className="font-bold text-sm uppercase tracking-wider">{item.name}</span>
                        {item.badge && (
                          <span className="px-2 py-1 bg-orange-100 text-orange-800 border-2 border-orange-300 rounded text-xs font-bold uppercase tracking-wider">
                            {item.badge}
                          </span>
                        )}
                      </div>
                      <p className="text-xs mt-1 text-gray-600">
                        {item.description}
                      </p>
                    </div>
                  </div>
                </Link>
              );
            })}
          </nav>

          {/* User section */}
          <div className="border-t-3 border-black p-4 bg-gray-50">
            <div className="flex items-center mb-3">
              <div className="flex-shrink-0">
                <div className="h-8 w-8 rounded-full bg-orange-100 border-2 border-orange-300 flex items-center justify-center">
                  <Users className="h-4 w-4 text-orange-800" />
                </div>
              </div>
              <div className="ml-3">
                <p className="text-sm font-bold text-black">{user?.fullName || user?.username || 'USER'}</p>
                <p className="text-xs text-gray-600">{user?.primaryEmailAddress?.emailAddress || 'user@zerotrace.com'}</p>
              </div>
            </div>
            <button
              onClick={() => signOut()}
              className="w-full p-2 bg-red-100 text-red-800 border-2 border-red-300 rounded text-sm font-bold uppercase tracking-wider hover:bg-red-200 transition-colors"
            >
              <LogOut className="mr-2 h-4 w-4 inline-block" />
              SIGN OUT
            </button>
          </div>
        </div>
      </div>

      {/* Main content */}
      <div className="lg:pl-64">
        {/* Top bar */}
        <div className="bg-white border-b-3 border-black shadow-lg">
          <div className="px-6 py-4">
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-x-4">
                <button
                  type="button"
                  className="p-2 hover:bg-gray-200 rounded border-2 border-black lg:hidden"
                  onClick={() => setSidebarOpen(true)}
                >
                  <Menu className="h-5 w-5" />
                </button>

                {/* Search bar */}
                <div className="relative flex flex-1 max-w-md">
                  <div className="pointer-events-none absolute inset-y-0 left-0 flex items-center pl-3">
                    <Search className="h-5 w-5 text-gray-400" />
                  </div>
                  <input
                    type="text"
                    placeholder="Search assets, vulnerabilities..."
                    className="w-full pl-10 pr-4 py-2 border-3 border-black rounded focus:outline-none focus:border-orange-500 focus:shadow-[4px_4px_0px_0px_rgba(255,107,0,1)] transition-all duration-150"
                  />
                </div>
              </div>

              {/* Right side actions */}
              <div className="flex items-center gap-x-4">
                {/* Notifications */}
                <button className="relative p-2 hover:bg-gray-200 rounded border-2 border-black transition-colors">
                  <Bell className="h-5 w-5" />
                  <span className="absolute top-1 right-1 h-2 w-2 bg-red-500 rounded-full"></span>
                </button>

                {/* Divider */}
                <div className="hidden lg:block lg:h-6 lg:w-px lg:bg-gray-300" />

                {/* User info */}
                <div className="flex items-center gap-x-4">
                  <span className="text-sm font-bold text-black uppercase tracking-wider">
                    WELCOME, {user?.name || 'ADMIN'}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Page content */}
        <main className="p-6 min-h-screen">
          <Outlet />
        </main>
      </div>
    </div>
  );
};

export default Layout;