import { useState } from 'react';
import { apiService, type RegisterRequest, type LoginRequest } from '../services/api';
import type { AuthMode } from '../types';

interface AuthProps {
  onAuthSuccess: (token: string, apiKey?: string) => void;
}

const Auth: React.FC<AuthProps> = ({ onAuthSuccess }) => {
  const [mode, setMode] = useState<AuthMode>('login');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');
    setSuccess('');

    // Validate passwords match for registration
    if (mode === 'register' && password !== confirmPassword) {
      setError('Passwords do not match');
      setIsLoading(false);
      return;
    }

    try {
      let response;
      
      if (mode === 'login') {
        const data: LoginRequest = { email, password };
        response = await apiService.login(data);
        setSuccess('Login successful!');
      } else {
        const data: RegisterRequest = { email, password };
        response = await apiService.register(data);
        setSuccess('Registration successful!');
      }
      
      // Store the JWT token
      apiService.storeToken(response.token);
      
      if (mode === 'login' && response.api_keys) {
        // Login: Store the list of API keys
        apiService.storeApiKeys(response.api_keys);
        setSuccess(`Login successful! Found ${response.api_keys.length} API key(s).`);
        onAuthSuccess(response.token, response.api_keys[0]?.key);
      } else if (mode === 'register' && response.api_keys) {
        // Register: Store the API keys array
        apiService.storeApiKeys(response.api_keys);
        setSuccess('Registration successful! API key created.');
        onAuthSuccess(response.token, response.api_keys[0]?.key);
      } else if (mode === 'register' && response.api_key) {
        // Register: Handle single API key (backward compatibility)
        apiService.storeApiKey(response.api_key);
        setSuccess('Registration successful! API key created.');
        onAuthSuccess(response.token, response.api_key);
      } else {
        // Fallback
        setSuccess(`${mode === 'login' ? 'Login' : 'Registration'} successful!`);
        onAuthSuccess(response.token);
      }
    } catch (err: any) {
      setError(err.response?.data?.error || `Failed to ${mode}`);
    } finally {
      setIsLoading(false);
    }
  };

  const toggleMode = () => {
    setMode(mode === 'login' ? 'register' : 'login');
    setError('');
    setSuccess('');
    setPassword('');
    setConfirmPassword('');
  };

  return (
    <div className="w-full max-w-md mx-auto">
      <h2 className="text-2xl font-bold text-gray-900 mb-6 text-center">
        {mode === 'login' ? 'Welcome Back' : 'Create Account'}
      </h2>

      <form onSubmit={handleSubmit} className="flex flex-col items-center space-y-4">
        <div className="w-[70%] grid grid-cols-11">
          <label htmlFor="email" className="col-span-5 text-xs font-medium text-gray-700">
            Email Address
          </label>
          <input
            type="email"  
            id="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
            className="col-span-6 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors duration-200 text-sm"
            placeholder="Enter your email"
          />
        </div>

        <div className="w-[70%] grid grid-cols-11">
          <label htmlFor="password" className="col-span-5 text-xs font-medium text-gray-700">
            Password
          </label>
          <input
            type="password"
            id="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            required
            minLength={6}
            className="col-span-6 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors duration-200 text-sm"
            placeholder="Enter your password"
          />
        </div>

        {mode === 'register' && (
          <div className="w-[70%] grid grid-cols-11">
            <label htmlFor="confirmPassword" className="col-span-5 text-xs font-medium text-gray-700 mb-1">
              Confirm Password
            </label>
            <input
              type="password"
              id="confirmPassword"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              required
              minLength={6}
              className="col-span-6 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors duration-200 text-sm"
              placeholder="Confirm your password"
            />
          </div>
        )}

        {error && (
          <div className="bg-red-50 border border-red-200 text-red-700 px-3 py-2 rounded-md flex items-center space-x-2">
            <svg className="w-4 h-4 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span className="text-sm">{error}</span>
          </div>
        )}

        {success && (
          <div className="bg-green-50 border border-green-200 text-green-700 px-3 py-2 rounded-md flex items-center space-x-2">
            <svg className="w-4 h-4 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            <span className="text-sm">{success}</span>
          </div>
        )}

        <button
          type="submit"
          disabled={isLoading}
          className="w-[20%] bg-blue-600 text-white py-2 px-4 rounded-md font-medium hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200 flex items-center justify-center space-x-2 text-sm"
        >
          {isLoading && (
            <svg className="animate-spin h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
              <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
              <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          )}
          <span>
            {isLoading ? 
              (mode === 'login' ? 'Signing In...' : 'Creating Account...') : 
              (mode === 'login' ? 'Sign In' : 'Create Account')
            }
          </span>
        </button>
      </form>

      <div className="mt-4 text-center">
        <p className="text-xs text-gray-600">
          {mode === 'login' ? "Don't have an account?" : "Already have an account?"}
          <button
            type="button"
            onClick={toggleMode}
            className="ml-2 text-blue-600 hover:text-blue-800 underline text-xs"
          >
            {mode === 'login' ? 'Create one' : 'Sign in'}
          </button>
        </p>
      </div>

      <div className="mt-4 pt-4 border-t border-gray-200">
        <p className="text-xs text-gray-500 text-center">
          {mode === 'register' ? 
            'By creating an account, you can manage your links and generate API keys.' :
            'Sign in to access your dashboard and manage your shortened URLs.'
          }
        </p>
      </div>
    </div>
  );
};

export default Auth;
