import { useState, useEffect } from 'react';
import { apiService, getErrorMessage } from '../services/api';
import { KeyIcon, DocumentDuplicateIcon, TrashIcon } from '@heroicons/react/24/outline';
import ErrorDisplay from './ErrorDisplay';
import { useToast } from '../contexts/ToastContext';

interface DashboardProps {
  activeTab: 'create' | 'list' | 'api-keys';
  setActiveTab: (tab: 'create' | 'list' | 'api-keys') => void;
}

interface ApiKeyData {
  key: string;
  createdAt: string;
}

const Dashboard: React.FC<DashboardProps> = ({ activeTab, setActiveTab }) => {
  const [newApiKey, setNewApiKey] = useState<string | null>(null);
  const [apiKeys, setApiKeys] = useState<ApiKeyData[]>([]);
  const [isCreatingApiKey, setIsCreatingApiKey] = useState(false);
  const [deletingKeys, setDeletingKeys] = useState<Set<string>>(new Set());
  const [copiedApiKey, setCopiedApiKey] = useState<string | null>(null);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const { showError, showSuccess } = useToast();

  // Load API keys when component mounts or active tab changes
  useEffect(() => {
    if (activeTab === 'api-keys') {
      loadApiKeys();
    }
  }, [activeTab]);

  const loadApiKeys = () => {
    const storedApiKeys = apiService.getApiKeys();
    setApiKeys(storedApiKeys);
  };

  const handleCreateApiKey = async () => {
    setIsCreatingApiKey(true);
    setError('');
    setSuccess('');
    setNewApiKey(null);

    try {
      const response = await apiService.createApiKey();
      setNewApiKey(response.api_key);
      showSuccess('API key created successfully!');
      
      // Add the new API key to the stored list
      const newApiKeyData: ApiKeyData = {
        key: response.api_key,
        createdAt: new Date().toISOString()
      };
      const updatedKeys = [...apiKeys, newApiKeyData];
      setApiKeys(updatedKeys);
      apiService.storeApiKeys(updatedKeys);
      
    } catch (err: any) {
      const errorMessage = getErrorMessage(err);
      showError(errorMessage);
    } finally {
      setIsCreatingApiKey(false);
    }
  };

  const handleDeleteApiKey = async (apiKeyToDelete: string) => {
    if (!confirm('Are you sure you want to delete this API key? This will also delete all associated links.')) {
      return;
    }

    setDeletingKeys(prev => new Set(prev).add(apiKeyToDelete));
    setError('');
    setSuccess('');

    try {
      await apiService.deleteApiKey(apiKeyToDelete);
      
      // Remove from local storage
      apiService.removeApiKeyFromList(apiKeyToDelete);
      
      // Update local state
      setApiKeys(prev => prev.filter(key => key.key !== apiKeyToDelete));
      
      showSuccess('API key and associated links deleted successfully!');
    } catch (err: any) {
      const errorMessage = getErrorMessage(err);
      showError(errorMessage);
    } finally {
      setDeletingKeys(prev => {
        const newSet = new Set(prev);
        newSet.delete(apiKeyToDelete);
        return newSet;
      });
    }
  };

  const copyApiKey = async (key: string) => {
    try {
      await navigator.clipboard.writeText(key);
      setCopiedApiKey(key);
      setTimeout(() => setCopiedApiKey(null), 2000); // Clear after 2 seconds
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
        <div className="max-w-4xl mx-auto bg-white rounded-lg shadow-md p-6">
          <div className="flex items-center mb-6">
            <KeyIcon className="h-6 w-6 text-blue-500 mr-3" />
            <h2 className="text-2xl font-bold text-gray-900">API Key Management</h2>
          </div>

          <div className="space-y-6">
            <p className="text-gray-600">
              Generate and manage API keys to access the URL shortener programmatically. 
              You can use these keys to create and manage links through our API.
            </p>

            {/* Usage instructions */}
            <div className="bg-blue-50 border border-blue-200 text-blue-700 px-4 py-3 rounded-md">
              <h4 className="font-medium mb-3">How to use your API keys:</h4>
              
              <div className="text-sm space-y-3">
                <div>
                  <p className="font-medium mb-1">Authentication:</p>
                  <p>• Include the header: <code className="bg-blue-100 px-1 rounded">X-API-KEY: your_api_key_here</code></p>
                </div>
                
                <div>
                  <p className="font-medium mb-1">Available endpoints:</p>
                  <ul className="space-y-1 ml-2">
                    <li>• <code className="bg-blue-100 px-1 rounded">POST /api/v1/links</code> - Create a new short link</li>
                    <li>• <code className="bg-blue-100 px-1 rounded">GET /api/v1/links</code> - Get all your links</li>
                    <li>• <code className="bg-blue-100 px-1 rounded">DELETE /api/v1/links/:shortCode</code> - Delete a specific link</li>
                  </ul>
                </div>                
                <p>• Each API key can create and manage its own set of links</p>
              </div>
            </div>

            {/* Create new API key section */}
            <div className="border-t pt-4">
              <button
                onClick={handleCreateApiKey}
                disabled={isCreatingApiKey}
                className="inline-flex items-center px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <KeyIcon className="h-4 w-4 mr-2" />
                {isCreatingApiKey ? 'Creating...' : 'Generate New API Key'}
              </button>
            </div>

            {/* Show newly created API key */}
            {newApiKey && (
              <div className="bg-green-50 border border-green-200 rounded-md p-4">
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm font-medium text-green-800">Your new API Key:</span>
                  <button
                    onClick={() => copyApiKey(newApiKey)}
                    className={`inline-flex items-center text-sm transition-colors ${
                      copiedApiKey === newApiKey 
                        ? 'text-green-600' 
                        : 'text-green-600 hover:text-green-800'
                    }`}
                    title={copiedApiKey === newApiKey ? "Copied!" : "Copy API key to clipboard"}
                  >
                    {copiedApiKey === newApiKey ? (
                      <>
                        <svg className="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                        </svg>
                        Copied!
                      </>
                    ) : (
                      <>
                        <DocumentDuplicateIcon className="h-4 w-4 mr-1" />
                        Copy
                      </>
                    )}
                  </button>
                </div>
                <code className="block bg-green-100 text-green-800 p-2 rounded text-sm font-mono break-all">
                  {newApiKey}
                </code>
                <p className="text-xs text-green-600 mt-2">
                  ⚠️ Store this key securely - it won't be shown again after you refresh the page!
                </p>
              </div>
            )}

            {/* Existing API keys list */}
            {apiKeys.length > 0 && (
              <div className="border-t pt-4">
                <h3 className="text-lg font-medium text-gray-900 mb-4">Your API Keys</h3>
                <div className="space-y-3">
                  {apiKeys.map((apiKeyData, index) => (
                    <div key={apiKeyData.key} className="bg-gray-50 border border-gray-200 rounded-md p-4">
                      <div className="flex items-center space-x-2 mb-2">
                        <span className="text-sm font-medium text-gray-900">
                          API Key #{index + 1}
                        </span>
                        <span className="text-xs text-gray-500">
                          Created: {new Date(apiKeyData.createdAt).toLocaleDateString()}
                        </span>
                      </div>
                      
                      <div className="grid grid-cols-20 items-center justify-between space-x-10">
                        <code className="col-span-15 bg-white text-gray-800 p-2 rounded text-sm font-mono break-all border">
                          {apiKeyData.key.substring(0, 8)}***
                        </code>
                        
                        <div className="items-center space-x-2 flex-shrink-0 col-span-5">
                          <button
                            onClick={() => copyApiKey(apiKeyData.key)}
                            className={`inline-flex items-center px-3 py-1 text-sm rounded border transition-colors ${
                              copiedApiKey === apiKeyData.key 
                                ? 'text-green-600 bg-green-50 border-green-200' 
                                : 'text-blue-600 hover:text-blue-800 hover:bg-blue-50 border-blue-200 hover:border-blue-300'
                            }`}
                            title={copiedApiKey === apiKeyData.key ? "Copied!" : "Copy full API key to clipboard"}
                          >
                            {copiedApiKey === apiKeyData.key ? (
                              <>
                                <svg className="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                                </svg>
                                Copied!
                              </>
                            ) : (
                              <>
                                <DocumentDuplicateIcon className="h-4 w-4 mr-1" />
                                Copy
                              </>
                            )}
                          </button>
                          <button
                            onClick={() => handleDeleteApiKey(apiKeyData.key)}
                            disabled={deletingKeys.has(apiKeyData.key)}
                            className="inline-flex items-center px-3 py-1 text-red-600 hover:text-red-800 hover:bg-red-50 text-sm rounded border border-red-200 hover:border-red-300 transition-colors disabled:opacity-50 disabled:cursor-not-allowed disabled:hover:bg-transparent"
                            title="Delete API key and all associated links"
                          >
                            {deletingKeys.has(apiKeyData.key) ? (
                              <>
                                <div className="animate-spin h-4 w-4 border-2 border-red-600 border-t-transparent rounded-full mr-1"></div>
                                Deleting...
                              </>
                            ) : (
                              <>
                                <TrashIcon className="h-4 w-4 mr-1" />
                                Delete
                              </>
                            )}
                          </button>
                        </div>
                      </div>
                      
                      <p className="text-xs text-gray-500 mt-2">
                        ⚠️ Deleting this API key will also delete all links created with it.
                      </p>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Empty state */}
            {apiKeys.length === 0 && !newApiKey && (
              <div className="text-center py-8">
                <KeyIcon className="mx-auto h-12 w-12 text-gray-400" />
                <h3 className="mt-2 text-sm font-medium text-gray-900">No API keys</h3>
                <p className="mt-1 text-sm text-gray-500">
                  Get started by creating your first API key.
                </p>
              </div>
            )}

            {/* Error and success messages */}
            {error && (
              <ErrorDisplay 
                error={error} 
                onDismiss={() => setError('')}
              />
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
