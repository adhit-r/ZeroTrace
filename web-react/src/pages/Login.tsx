import React, { useState } from 'react';
import { useAuth } from '../contexts/AuthContext';
import { Terminal, Eye, EyeOff, Shield } from 'lucide-react';

const Login: React.FC = () => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  
  const { login } = useAuth();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');

    try {
      await login(email, password);
    } catch (err: any) {
      setError(err.response?.data?.message || 'LOGIN FAILED');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-md w-full space-y-8">
        {/* Logo and Title */}
        <div className="text-center">
          <div className="mx-auto h-16 w-16 flex items-center justify-center rounded-lg bg-medium-gray border border-gold mb-6">
            <Terminal className="h-10 w-10 text-gold" />
          </div>
          <h2 className="text-3xl font-bold text-gold text-glow">
            VULNDETECT
          </h2>
          <p className="mt-2 text-text-secondary uppercase tracking-wide">
            ENTERPRISE VULNERABILITY MANAGEMENT
          </p>
        </div>
        
        {/* Login Form */}
        <div className="card card-terminal glow-border">
          <form className="space-y-6" onSubmit={handleSubmit}>
            <div>
              <label htmlFor="email" className="block text-sm font-medium text-text-primary mb-2 uppercase tracking-wide">
                EMAIL ADDRESS
              </label>
              <input
                id="email"
                name="email"
                type="email"
                autoComplete="email"
                required
                className="input"
                placeholder="admin@vulndetect.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
            </div>
            
            <div>
              <label htmlFor="password" className="block text-sm font-medium text-text-primary mb-2 uppercase tracking-wide">
                PASSWORD
              </label>
              <div className="relative">
                <input
                  id="password"
                  name="password"
                  type={showPassword ? 'text' : 'password'}
                  autoComplete="current-password"
                  required
                  className="input pr-10"
                  placeholder="••••••••"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                />
                <button
                  type="button"
                  className="absolute inset-y-0 right-0 pr-3 flex items-center"
                  onClick={() => setShowPassword(!showPassword)}
                >
                  {showPassword ? (
                    <EyeOff className="h-5 w-5 text-text-muted hover:text-gold" />
                  ) : (
                    <Eye className="h-5 w-5 text-text-muted hover:text-gold" />
                  )}
                </button>
              </div>
            </div>

            {error && (
              <div className="p-4 bg-critical bg-opacity-20 border border-critical rounded">
                <div className="text-sm text-critical font-medium">{error}</div>
              </div>
            )}

            <div>
              <button
                type="submit"
                disabled={isLoading}
                className="btn btn-primary w-full"
              >
                {isLoading ? (
                  <>
                    <div className="animate-spin h-4 w-4 mr-2 border-2 border-current border-t-transparent rounded-full"></div>
                    AUTHENTICATING...
                  </>
                ) : (
                  'SIGN IN'
                )}
              </button>
            </div>

            <div className="text-center">
              <p className="text-sm text-text-muted">
                DEMO: admin@vulndetect.com / password
              </p>
            </div>
          </form>
        </div>

        {/* Footer */}
        <div className="text-center">
          <p className="text-xs text-text-muted uppercase tracking-wider">
            SECURE ENTERPRISE VULNERABILITY DETECTION
          </p>
        </div>
      </div>
    </div>
  );
};

export default Login;
