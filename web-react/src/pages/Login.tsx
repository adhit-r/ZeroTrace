import React from 'react';

const Login: React.FC = () => {
  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <div className="max-w-md w-full space-y-8">
        <div>
          <h2 className="mt-6 text-center text-3xl font-extrabold text-gray-900">
            Sign in to ZeroTrace
          </h2>
          <p className="mt-2 text-center text-sm text-gray-600">
            Enterprise vulnerability management platform
          </p>
        </div>
        <div className="bg-white p-8 rounded-lg shadow">
          <p className="text-gray-600 text-center">
            Authentication system coming soon...
          </p>
        </div>
      </div>
    </div>
  );
};

export default Login;