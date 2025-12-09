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
        <div className="w-64 bg-white flex-shrink-0 flex flex-col">
          {/* Orange Header */}
          <div className="bg-orange-500 px-4 py-4">
            <h1 className="text-xl font-black text-white uppercase">ZEROTRACE</h1>
          </div>
          
          {/* Navigation Menu */}
          <nav className="flex-1 p-4 space-y-3">
            {navigation.map((item) => {
              const isActive = location.pathname === item.href;
              return (
                <Link
                  key={item.name}
                  to={item.href}
                  className={`flex items-center px-4 py-3 rounded border-4 border-black font-black uppercase text-sm transition-all ${
                    isActive 
                      ? 'bg-orange-500 text-white shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]' 
                      : 'bg-white text-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:shadow-[6px_6px_0px_0px_rgba(0,0,0,1)] hover:-translate-x-0.5 hover:-translate-y-0.5 active:shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] active:translate-x-0 active:translate-y-0'
                  }`}
                  style={{
                    boxShadow: isActive 
                      ? '4px 4px 0px 0px rgba(0,0,0,1)' 
                      : '4px 4px 0px 0px rgba(0,0,0,1)'
                  }}
                >
                  <item.icon className="w-5 h-5 mr-3 flex-shrink-0" />
                  <span>{item.name}</span>
                </Link>
              );
            })}
          </nav>
        </div>

        {/* Desktop Main Content */}
        <div className="flex-1 flex flex-col overflow-hidden">
          {/* Top bar */}
          <div className="bg-white border-b-4 border-black">
            <div className="flex items-center justify-between h-16 px-4">
              <div className="flex items-center space-x-4">
                <span className="text-sm font-bold text-black uppercase">Enterprise Security Management Platform</span>
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
        <div className="bg-white border-b-4 border-black">
          <div className="flex items-center justify-between h-16 px-4">
            <button
              onClick={() => setSidebarOpen(true)}
              className="p-2 bg-orange-500 text-white border-4 border-black rounded shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:shadow-[6px_6px_0px_0px_rgba(0,0,0,1)] hover:-translate-x-0.5 hover:-translate-y-0.5 active:shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] active:translate-x-0 active:translate-y-0 transition-all"
            >
              <Menu className="w-6 h-6" />
            </button>
            <h1 className="text-xl font-black text-black uppercase">ZEROTRACE</h1>
            <div></div>
          </div>
        </div>

        {/* Mobile Sidebar */}
        <div className={`
          ${sidebarOpen ? 'translate-x-0' : '-translate-x-full'} 
          fixed inset-y-0 left-0 z-50 w-64 bg-white
          transform transition-transform duration-300 ease-in-out flex flex-col
        `}>
          {/* Orange Header */}
          <div className="bg-orange-500 px-4 py-4 flex items-center justify-between">
            <h1 className="text-xl font-black text-white uppercase">ZEROTRACE</h1>
            <button
              onClick={() => setSidebarOpen(false)}
              className="p-2 bg-white border-4 border-black rounded shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:shadow-[6px_6px_0px_0px_rgba(0,0,0,1)] hover:-translate-x-0.5 hover:-translate-y-0.5 active:shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] active:translate-x-0 active:translate-y-0 transition-all"
            >
              <X className="w-5 h-5 text-black" />
            </button>
          </div>
          
          {/* Navigation Menu */}
          <nav className="flex-1 p-4 space-y-3 overflow-y-auto">
            {navigation.map((item) => {
              const isActive = location.pathname === item.href;
              return (
                <Link
                  key={item.name}
                  to={item.href}
                  onClick={() => setSidebarOpen(false)}
                  className={`flex items-center px-4 py-3 rounded border-4 border-black font-black uppercase text-sm transition-all ${
                    isActive 
                      ? 'bg-orange-500 text-white shadow-[4px_4px_0px_0px_rgba(0,0,0,1)]' 
                      : 'bg-white text-black shadow-[4px_4px_0px_0px_rgba(0,0,0,1)] hover:shadow-[6px_6px_0px_0px_rgba(0,0,0,1)] hover:-translate-x-0.5 hover:-translate-y-0.5 active:shadow-[2px_2px_0px_0px_rgba(0,0,0,1)] active:translate-x-0 active:translate-y-0'
                  }`}
                  style={{
                    boxShadow: isActive 
                      ? '4px 4px 0px 0px rgba(0,0,0,1)' 
                      : '4px 4px 0px 0px rgba(0,0,0,1)'
                  }}
                >
                  <item.icon className="w-5 h-5 mr-3 flex-shrink-0" />
                  <span>{item.name}</span>
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