import React, { useState } from 'react';
import { Link, useLocation } from 'react-router-dom';
import { 
  Home, 
  Server, 
  Shield, 
  Search, 
  Settings, 
  User,
  Menu,
  X
} from 'lucide-react';

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const [sidebarOpen, setSidebarOpen] = useState(false);
  const location = useLocation();

  const navigation = [
    { name: 'Dashboard', href: '/', icon: Home },
    { name: 'Agents', href: '/agents', icon: Server },
    { name: 'Vulnerabilities', href: '/vulnerabilities', icon: Shield },
    { name: 'Security', href: '/security', icon: Shield },
    { name: 'Scans', href: '/scans', icon: Search },
    { name: 'Topology', href: '/topology', icon: Server },
    { name: 'Settings', href: '/settings', icon: Settings },
    { name: 'Profile', href: '/profile', icon: User },
  ];

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Mobile sidebar overlay */}
      {sidebarOpen && (
        <div 
          className="fixed inset-0 bg-black bg-opacity-50 z-40 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      {/* Desktop Layout */}
      <div className="hidden lg:flex h-screen">
        {/* Desktop Sidebar */}
        <div className="w-64 bg-white border-3 border-black shadow-neubrutalist-lg flex-shrink-0">
          <div className="flex items-center justify-between h-16 px-4 border-b-3 border-black bg-orange-500">
            <h1 className="text-xl font-black text-white uppercase tracking-wider">ZeroTrace</h1>
          </div>
          <nav className="mt-4 p-2">
            {navigation.map((item) => {
              const isActive = location.pathname === item.href;
              return (
                <Link
                  key={item.name}
                  to={item.href}
                  className={`flex items-center px-4 py-3 mb-2 text-black font-bold uppercase tracking-wider border-3 border-black rounded shadow-neubrutalist hover:shadow-neubrutalist-hover hover:translate-x-0.5 hover:translate-y-0.5 transition-all duration-150 ease-in-out ${
                    isActive 
                      ? 'bg-orange-500 text-white' 
                      : 'bg-white hover:bg-orange-100'
                  }`}
                >
                  <item.icon className="w-5 h-5 mr-3" />
                  {item.name}
                </Link>
              );
            })}
          </nav>
        </div>

        {/* Desktop Main Content */}
        <div className="flex-1 flex flex-col overflow-hidden">
          {/* Top bar */}
          <div className="bg-white border-b-3 border-black shadow-neubrutalist">
            <div className="flex items-center justify-between h-16 px-4">
              <div className="flex items-center space-x-4">
                <span className="text-sm font-bold text-black uppercase tracking-wider">Enterprise Security Management Platform</span>
              </div>
            </div>
          </div>

          {/* Page content */}
          <main className="flex-1 overflow-y-auto bg-gray-100">
            {children}
          </main>
        </div>
      </div>

      {/* Mobile Layout */}
      <div className="lg:hidden">
        {/* Mobile Top bar */}
        <div className="bg-white border-b-3 border-black shadow-neubrutalist">
          <div className="flex items-center justify-between h-16 px-4">
            <button
              onClick={() => setSidebarOpen(true)}
              className="p-2 bg-orange-500 text-white border-3 border-black rounded shadow-neubrutalist hover:shadow-neubrutalist-hover hover:translate-x-0.5 hover:translate-y-0.5 transition-all duration-150 ease-in-out"
            >
              <Menu className="w-6 h-6" />
            </button>
            <h1 className="text-xl font-black text-black uppercase tracking-wider">ZeroTrace</h1>
            <div></div>
          </div>
        </div>

        {/* Mobile Sidebar */}
        <div className={`
          ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'} 
          fixed inset-y-0 left-0 z-50 w-64 bg-white border-3 border-black shadow-neubrutalist-lg 
          transform transition-transform duration-300 ease-in-out
        `}>
          <div className="flex items-center justify-between h-16 px-4 border-b-3 border-black bg-orange-500">
            <h1 className="text-xl font-black text-white uppercase tracking-wider">ZeroTrace</h1>
            <button
              onClick={() => setSidebarOpen(false)}
              className="p-2 bg-white border-2 border-black rounded hover:bg-gray-100 transition-colors"
            >
              <X className="w-5 h-5 text-black" />
            </button>
          </div>
          <nav className="mt-4 p-2">
            {navigation.map((item) => {
              const isActive = location.pathname === item.href;
              return (
                <Link
                  key={item.name}
                  to={item.href}
                  onClick={() => setSidebarOpen(false)}
                  className={`flex items-center px-4 py-3 mb-2 text-black font-bold uppercase tracking-wider border-3 border-black rounded shadow-neubrutalist hover:shadow-neubrutalist-hover hover:translate-x-0.5 hover:translate-y-0.5 transition-all duration-150 ease-in-out ${
                    isActive 
                      ? 'bg-orange-500 text-white' 
                      : 'bg-white hover:bg-orange-100'
                  }`}
                >
                  <item.icon className="w-5 h-5 mr-3" />
                  {item.name}
                </Link>
              );
            })}
          </nav>
        </div>

        {/* Mobile Page content */}
        <main className="min-h-screen bg-gray-100">
          {children}
        </main>
      </div>
    </div>
  );
};

export default Layout;