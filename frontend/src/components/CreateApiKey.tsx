import { useState } from 'react';
import { apiService } from '../services/api';

interface CreateApiKeyProps {
  onApiKeyValidated: () => void;
}

const CreateApiKey: React.FC<CreateApiKeyProps> = ({ onApiKeyValidated }) => {
  const [existingApiKey, setExistingApiKey] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleUseExisting = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');
    setSuccess('');

    try {
      // Store the API key temporarily
      apiService.storeApiKey(existingApiKey);
      
      // Validate the API key by trying to fetch user links
      await apiService.getUserLinks();
      
      setSuccess('API key validated successfully!');
      onApiKeyValidated();
    } catch (err: any) {
      // Remove invalid API key
      apiService.removeApiKey();
      setError('Invalid API key. Please check and try again.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="max-w-md mx-auto bg-white rounded-lg shadow-md p-6">
      <h2 className="text-2xl font-bold text-gray-900 mb-6 text-center">
        Use Existing API Key
      </h2>

      <form onSubmit={handleUseExisting} className="space-y-4">
        <div>
          <label htmlFor="existingApiKey" className="block text-sm font-medium text-gray-700 mb-2">
            Your API Key
          </label>
          <input
            type="text"
            id="existingApiKey"
            value={existingApiKey}
            onChange={(e) => setExistingApiKey(e.target.value)}
            required
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 font-mono text-sm"
            placeholder="Enter your 32-character API key"
          />
        </div>

        <button
          type="submit"
          disabled={isLoading}
          className="w-full bg-green-600 text-white py-2 px-4 rounded-md hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {isLoading ? 'Validating...' : 'Use This API Key'}
        </button>

        <div className="text-sm text-gray-600 text-center">
          <p>Already have an API key? Enter it here to access your account.</p>
        </div>

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
      </form>
    </div>
  );
};

export default CreateApiKey;
