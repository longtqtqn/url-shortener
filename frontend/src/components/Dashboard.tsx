import { useState } from 'react';
import { apiService } from '../services/api';
import { KeyIcon, DocumentDuplicateIcon } from '@heroicons/react/24/outline';

interface DashboardProps {
  activeTab: 'create' | 'list' | 'api-keys';
  setActiveTab: (tab: 'create' | 'list' | 'api-keys') => void;
}

const Dashboard: React.FC<DashboardProps> = ({ activeTab, setActiveTab }) => {
  const [apiKey, setApiKey] = useState<string | null>(null);
  const [isCreatingApiKey, setIsCreatingApiKey] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleCreateApiKey = async () => {
    setIsCreatingApiKey(true);
    setError('');
    setSuccess('');

    try {
      const response = await apiService.createApiKey();
      setApiKey(response.api_key);
      setSuccess('API key created successfully!');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create API key');
    } finally {
      setIsCreatingApiKey(false);
    }
  };

  const copyApiKey = async () => {
    if (!apiKey) return;

    try {
      await navigator.clipboard.writeText(apiKey);
      setSuccess('API key copied to clipboard!');
    } catch (err) {
      setError('Failed to copy API key to clipboard');
    }
  };

  return (
    <div className="space-y">
      {/* Tab Navigation */}
      <div className="border-b border-gray-200">
        <nav className="flex space-x-8">
          <button
            onClick={() => setActiveTab('create')}
            className={`py-2 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'create'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            }`}
          >
            Create Link
          </button>
          <button
            onClick={() => setActiveTab('list')}
            className={`py-2 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'list'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            }`}
          >
            My Links
          </button>
          <button
            onClick={() => setActiveTab('api-keys')}
            className={`py-2 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'api-keys'
                ? 'border-blue-500 text-blue-600'
                : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
            }`}
          >
            API Keys
          </button>
        </nav>
      </div>

      {/* API Keys Tab Content */}
      {activeTab === 'api-keys' && (
        <div className="max-w-2xl mx-auto bg-white rounded-lg shadow-md p-6">
          <div className="grid items-center grid-cols-17">
            <KeyIcon className="col-start-5 col-span-1 h-6 w-6 text-blue-500 mr-3" />
            <h2 className="col-span-5 text-2xl font-bold text-gray-900">API Key Management</h2>
          </div>

          <div className="space-y-4">
            <p className="text-gray-600">
              Generate an API key to access the URL shortener programmatically. 
              You can use this key to create and manage links through our API.
            </p>

            {!apiKey ? (
              <div className="space-y-4">
                <button
                  onClick={handleCreateApiKey}
                  disabled={isCreatingApiKey}
                  className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <KeyIcon className="h-4 w-4 mr-2" />
                  {isCreatingApiKey ? 'Creating...' : 'Generate API Key'}
                </button>

                {apiService.hasApiKey() && (
                  <div className="bg-yellow-50 border border-yellow-200 text-yellow-700 px-4 py-3 rounded-md">
                    <p className="text-sm">
                      You already have an API key stored locally. If you've lost it, 
                      you can generate a new one above.
                    </p>
                  </div>
                )}
              </div>
            ) : (
              <div className="space-y-4">
                <div className="bg-green-50 border border-green-200 rounded-md p-4">
                  <div className="flex items-center justify-between mb-2">
                    <span className="text-sm font-medium text-green-800">Your API Key:</span>
                    <button
                      onClick={copyApiKey}
                      className="inline-flex items-center text-green-600 hover:text-green-800 text-sm"
                    >
                      <DocumentDuplicateIcon className="h-4 w-4 mr-1" />
                      Copy
                    </button>
                  </div>
                  <code className="block bg-green-100 text-green-800 p-2 rounded text-sm font-mono break-all">
                    {apiKey}
                  </code>
                </div>

                <div className="bg-blue-50 border border-blue-200 text-blue-700 px-4 py-3 rounded-md">
                  <h4 className="font-medium mb-2">How to use your API key:</h4>
                  <ul className="text-sm space-y-1">
                    <li>• Include the header: <code className="bg-blue-100 px-1 rounded">X-API-KEY: {apiKey}</code></li>
                    <li>• Use endpoints like: <code className="bg-blue-100 px-1 rounded">POST /api/v1/links</code></li>
                    <li>• Store this key securely - it won't be shown again</li>
                  </ul>
                </div>

                <button
                  onClick={handleCreateApiKey}
                  disabled={isCreatingApiKey}
                  className="text-sm text-gray-600 hover:text-gray-800 underline"
                >
                  Generate a new API key
                </button>
              </div>
            )}

            {error && (
              <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md">
                {error}
              </div>
            )}

            {success && (
              <div className="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded-md">
                {success}
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default Dashboard;
