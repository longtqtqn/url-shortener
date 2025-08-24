import { useState } from 'react';
import { apiService, type CreateLinkRequest } from '../services/api';

interface CreateLinkProps {
  onLinkCreated: () => void;
  isAuthenticated: boolean;
}

const CreateLink: React.FC<CreateLinkProps> = ({ onLinkCreated, isAuthenticated }) => {
  const [longUrl, setLongUrl] = useState('');
  const [customShortCode, setCustomShortCode] = useState('');
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const [createdLink, setCreatedLink] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError('');
    setSuccess('');
    setCreatedLink('');

    try {
      const data: CreateLinkRequest = {
        long_url: longUrl,
        ...(customShortCode && { short_code: customShortCode }),
      };
      
      const response = isAuthenticated 
        ? await apiService.createLink(data)
        : await apiService.createLinkPublic(data);
      
      setSuccess('Short URL created successfully!');
      setCreatedLink(response.shortened_url);
      setLongUrl('');
      setCustomShortCode('');
      onLinkCreated();
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create short URL');
    } finally {
      setIsLoading(false);
    }
  };

  const copyToClipboard = async () => {
    try {
      await navigator.clipboard.writeText(createdLink);
      setSuccess('Link copied to clipboard!');
    } catch (err) {
      setError('Failed to copy to clipboard');
    }
  };

  return (
    <div className="flex justify-center pt-16 min-h-[calc(100vh-200px)]">
      <div className="w-[60%] bg-white rounded-lg shadow-md p-6">
        <h2 className="text-2xl font-bold text-gray-900 mb-6 text-center">
          Create Short URL
        </h2>
        
        <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label htmlFor="longUrl" className="block text-sm font-medium text-gray-700 mb-2">
            Long URL *
          </label>
          <input
            type="url"
            id="longUrl"
            value={longUrl}
            onChange={(e) => setLongUrl(e.target.value)}
            required
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            placeholder="https://example.com/very-long-url"
          />
        </div>

        <div>
          <label htmlFor="customShortCode" className="block text-sm font-medium text-gray-700 mb-2">
            Custom Short Code (Optional)
          </label>
          <input
            type="text"
            id="customShortCode"
            value={customShortCode}
            onChange={(e) => setCustomShortCode(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
            placeholder="my-custom-link"
          />
          <p className="mt-1 text-sm text-gray-500">
            Leave empty to auto-generate a 6-character code
          </p>
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

        {createdLink && (
          <div className="bg-blue-50 border border-blue-200 text-blue-700 px-4 py-3 rounded-md">
            <div className="flex items-center justify-between">
              <span className="font-medium">Your short URL:</span>
              <button
                type="button"
                onClick={copyToClipboard}
                className="text-blue-600 hover:text-blue-800 underline"
              >
                Copy
              </button>
            </div>
            <div className="mt-2 break-all font-mono text-sm">
              {createdLink}
            </div>
          </div>
        )}

        <button
          type="submit"
          disabled={isLoading}
          className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {isLoading ? 'Creating...' : 'Create Short URL'}
        </button>
        </form>
      </div>
    </div>
  );
};

export default CreateLink;
