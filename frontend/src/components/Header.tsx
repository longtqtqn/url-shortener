import React from 'react';
import { LinkIcon, UserIcon } from '@heroicons/react/24/outline';
import { apiService } from '../services/api';

interface HeaderProps {
  onLogout: () => void;
  userEmail?: string;
}

const Header: React.FC<HeaderProps> = ({ onLogout, userEmail }) => {
  const handleLogout = () => {
    apiService.logout();
    onLogout();
  };

  const getAuthInfo = () => {
    if (apiService.hasToken()) {
      return userEmail ? `Signed in as ${userEmail}` : 'Signed in';
    } else if (apiService.hasApiKey()) {
      return `API Key: ${apiService.getApiKey()?.substring(0, 8)}...`;
    }
    return 'Not authenticated';
  };

  return (
    <header className="bg-white shadow-sm border-b border-gray-200">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-16">
          <div className="flex items-center">
            <LinkIcon className="h-8 w-8 text-blue-600" />
            <h1 className="ml-2 text-xl font-bold text-gray-900">
              URL Shortener
            </h1>
          </div>
          
          <div className="flex items-center space-x-4">
            <div className="flex items-center text-sm text-gray-500">
              <UserIcon className="h-4 w-4 mr-1" />
              <span>{getAuthInfo()}</span>
            </div>
            <button
              onClick={handleLogout}
              className="text-gray-600 hover:text-gray-800 px-3 py-2 rounded-md text-sm font-medium hover:bg-gray-100"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;
