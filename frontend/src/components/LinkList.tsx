import { useState, useEffect } from 'react';
import { apiService, type Link } from '../services/api';
import { TrashIcon, LinkIcon, CalendarIcon, EyeIcon } from '@heroicons/react/24/outline';

const LinkList: React.FC = () => {
  const [links, setLinks] = useState<Link[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [copiedUrl, setCopiedUrl] = useState<string | null>(null);

  const fetchLinks = async () => {
    try {
      setIsLoading(true);
      const userLinks = await apiService.getUserLinks();
      setLinks(userLinks);
      setError('');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to fetch links');
    } finally {
      setIsLoading(false);
    }
  };

  const handleDelete = async (shortCode: string) => {
    if (!confirm('Are you sure you want to delete this link?')) {
      return;
    }

    try {
      await apiService.deleteLink(shortCode);
      setLinks(links.filter(link => !link.shortURL.includes(shortCode)));
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to delete link');
    }
  };

  const copyToClipboard = async (url: string) => {
    try {
      await navigator.clipboard.writeText(url);
      setCopiedUrl(url);
      setTimeout(() => setCopiedUrl(null), 2000); // Clear after 2 seconds
    } catch (err) {
      setError('Failed to copy to clipboard');
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  useEffect(() => {
    fetchLinks();
  }, []);

  if (isLoading) {
    return (
      <div className="flex justify-center items-center py-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded-md">
        {error}
        <button
          onClick={fetchLinks}
          className="ml-4 text-red-600 hover:text-red-800 underline"
        >
          Retry
        </button>
      </div>
    );
  }

  if (links.length === 0) {
    return (
      <div className="text-center py-8">
        <LinkIcon className="mx-auto h-12 w-12 text-gray-400" />
        <h3 className="mt-2 text-sm font-medium text-gray-900">No links yet</h3>
        <p className="mt-1 text-sm text-gray-500">
          Create your first short URL to get started!
        </p>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow-md mx-[1%]">
      <div className="px-6 py-4 border-b border-gray-200">
        <h2 className="text-lg font-medium text-gray-900">Your Short URLs</h2>
        <p className="text-sm text-gray-500">
          Manage and track your shortened links
        </p>
      </div>
      
      {/* Table Header */}
      <div className="px-6 py-3 bg-gray-50 border-b border-gray-200">
        <div className="grid grid-cols-14 gap-3 text-xs font-medium text-gray-500 uppercase tracking-wide">
          <div className="col-span-2">API Key</div>
          <div className="col-span-3">Short URL</div>
          <div className="col-span-3">Original URL</div>
          <div className="col-span-1">Clicks</div>
          <div className="col-span-2">Created At</div>
          <div className="col-span-2">Last Clicked</div>
          <div className="col-span-1 text-center">Actions</div>
        </div>
      </div>
      
      {/* Table Body */}
      <div className="divide-y divide-gray-200">
        {links.map((link, index) => (
          <div key={index} className="px-6 py-4 hover:bg-gray-50">
            <div className="grid grid-cols-14 gap-3 items-center">
              {/* API Key */}
              <div className="col-span-2">
                <span className="text-xs font-mono bg-gray-100 px-2 py-1 rounded">
                  {link.apiKey ? `${link.apiKey.substring(0, 8)}****` : 'N/A'}
                </span>
              </div>
              
              {/* Short URL */}
              <div className="col-span-3">
                <a 
                  href={link.shortURL} 
                  target="_blank" 
                  rel="noopener noreferrer"
                  className="text-sm font-medium text-blue-600 hover:text-blue-800 underline break-all"
                >
                  {link.shortURL}
                </a>
              </div>
              
              {/* Original URL */}
              <div className="col-span-3">
                <span className="text-sm text-gray-600 break-all">
                  {link.longURL.length > 50 ? `${link.longURL.substring(0, 50)}...` : link.longURL}
                </span>
              </div>
              
              {/* Click Count */}
              <div className="col-span-1 flex items-center">
                <div className="grid grid-cols-6 items-center gap-1 w-full">
                  <svg className="!h-6 !w-6 text-gray-400 col-span-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                  </svg>
                  <span className="text-sm font-medium col-span-4">{link.clickCount}</span>
                </div>
              </div>
              
              {/* Created At */}
              <div className="col-span-2">
                <div className="grid grid-cols-7 items-center gap-1">
                  <CalendarIcon className="h-6 w-6 text-gray-400 col-span-2" />
                  <span className="text-xs text-gray-500 col-span-5">
                    {formatDate(link.createdAt)}
                  </span>
                </div>
              </div>
              
              {/* Last Clicked */}
              <div className="col-span-2">
                {link.lastClicked ? (
                  <div className="grid grid-cols-4 items-center gap-1">
                    <svg className="h-6 w-6 col-span-1 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    <span className="text-xs text-gray-500 col-span-3">
                      {formatDate(link.lastClicked)}
                    </span>
                  </div>
                ) : (
                  <span className="text-xs text-gray-400 italic">Never</span>
                )}
              </div>
              
              {/* Actions */}
              <div className="col-span-1">
                <div className="grid grid-cols-2">
                  <button
                    onClick={() => copyToClipboard(link.shortURL)}
                    className={`p-1 rounded transition-colors ${
                      copiedUrl === link.shortURL 
                        ? 'text-green-600 bg-green-50' 
                        : 'text-blue-600 hover:text-blue-800 hover:bg-blue-50'
                    }`}
                    title={copiedUrl === link.shortURL ? "Copied!" : "Copy short URL"}
                  >
                    {copiedUrl === link.shortURL ? (
                      <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                      </svg>
                    ) : (
                      <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                      </svg>
                    )}
                  </button>
                  <button
                    onClick={() => handleDelete(link.shortURL.split('/').pop()!)}
                    className="text-red-600 hover:text-red-800 p-1 rounded hover:bg-red-50"
                    title="Delete link"
                  >
                    <TrashIcon className="h-4 w-4" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
};

export default LinkList;
