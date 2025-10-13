import React from 'react';
import { useClerk } from '@clerk/clerk-react';
import { Terminal, LogIn } from 'lucide-react';

const Login: React.FC = () => {
  const { openSignIn } = useClerk();

  const handleSignIn = () => {
    openSignIn({
      redirectUrl: '/',
    });
  };

  return (
    <div className="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8 bg-gradient-to-br from-gray-900 via-gray-800 to-black">
      <div className="max-w-md w-full space-y-8">
        {/* Logo and Title */}
        <div className="text-center">
          <div className="mx-auto h-20 w-20 flex items-center justify-center rounded-lg bg-gradient-to-br from-orange-500 to-red-600 border-3 border-black shadow-[8px_8px_0_0_rgba(0,0,0,1)] mb-6">
            <Terminal className="h-12 w-12 text-white" />
          </div>
          <h2 className="text-4xl font-black text-white mb-2 tracking-tight">
            ZEROTRACE
          </h2>
          <p className="text-sm text-gray-400 uppercase tracking-widest font-bold">
            ENTERPRISE VULNERABILITY MANAGEMENT
          </p>
        </div>

        {/* Sign In Card */}
        <div className="bg-white border-3 border-black rounded-lg shadow-[8px_8px_0_0_rgba(0,0,0,1)] p-8">
          <div className="space-y-6">
            <div className="text-center">
              <h3 className="text-xl font-bold text-black uppercase tracking-wide mb-2">
                SECURE LOGIN
              </h3>
              <p className="text-sm text-gray-600">
                Access your security dashboard
              </p>
            </div>

            <button
              onClick={handleSignIn}
              className="w-full flex items-center justify-center gap-3 px-6 py-4 bg-gradient-to-r from-orange-500 to-red-600 text-white font-bold uppercase tracking-wide rounded-lg border-3 border-black shadow-[4px_4px_0_0_rgba(0,0,0,1)] hover:shadow-[6px_6px_0_0_rgba(0,0,0,1)] hover:translate-x-[-2px] hover:translate-y-[-2px] transition-all"
            >
              <LogIn className="h-5 w-5" />
              SIGN IN WITH CLERK
            </button>

            <div className="text-center pt-4 border-t-2 border-gray-200">
              <p className="text-xs text-gray-500 uppercase tracking-wider">
                Protected by end-to-end encryption
              </p>
            </div>
          </div>
        </div>

        {/* Footer */}
        <div className="text-center">
          <p className="text-xs text-gray-500 uppercase tracking-widest">
            SECURE ENTERPRISE VULNERABILITY DETECTION
          </p>
        </div>
      </div>
    </div>
  );
};

export default Login;
