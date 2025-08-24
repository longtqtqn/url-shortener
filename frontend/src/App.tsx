import { useState, useEffect } from 'react';
import Header from './components/Header';
import Auth from './components/Auth';
import CenteredModal from './components/CenteredModal';
import CreateLink from './components/CreateLink';
import LinkList from './components/LinkList';
import Dashboard from './components/Dashboard';
import { apiService } from './services/api';

type AuthState = 'public' | 'authenticated';

function App() {
  const [authState, setAuthState] = useState<AuthState>('public');
  const [activeTab, setActiveTab] = useState<'create' | 'list' | 'api-keys'>('create');
  const [userEmail, setUserEmail] = useState<string>('');
  const [showAuthModal, setShowAuthModal] = useState(false);

  useEffect(() => {
    checkAuthentication();
  }, []);

  const checkAuthentication = async () => {
    // Check if user has JWT token
    if (apiService.hasToken()) {
      try {
        // Validate token by trying to fetch user data
        await apiService.getUserLinks();
        setAuthState('authenticated');
        return;
      } catch (error) {
        // Token is invalid, remove it
        console.log('Token validation failed, removing token');
        apiService.removeToken();
      }
    }

    // Check if user has API key (for backward compatibility)
    if (apiService.hasApiKey()) {
      try {
        // Validate API key by trying to fetch user data
        await apiService.getUserLinks();
        setAuthState('authenticated');
        return;
      } catch (error) {
        // API key is invalid, remove it
        console.log('API key validation failed, removing API key');
        apiService.removeApiKey();
      }
    }

    // Default to public mode (no authentication required)
    setAuthState('public');
  };

  const handleAuthSuccess = (token: string, apiKey?: string) => {
    apiService.storeToken(token);
    if (apiKey) {
      apiService.storeApiKey(apiKey);
    }
    setAuthState('authenticated');
    setActiveTab('create');
    setShowAuthModal(false);
  };

  const handleLogout = () => {
    apiService.logout();
    setAuthState('public');
    setActiveTab('create');
    setUserEmail('');
    setShowAuthModal(false);
  };

  const handleLinkCreated = () => {
    // Refresh the link list if we're on that tab and authenticated
    if (activeTab === 'list' && authState === 'authenticated') {
      // Force a re-render of LinkList
      setActiveTab('create');
      setTimeout(() => setActiveTab('list'), 100);
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header with login/logout functionality */}
      {authState === 'authenticated' ? (
        <Header onLogout={handleLogout} userEmail={userEmail} />
      ) : (
        <header className="bg-white shadow-sm border-b border-gray-200">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex justify-between items-center h-16">
              <div className="flex items-center">
                <h1 className="text-xl font-bold text-gray-900">
                  URL Shortener
                </h1>
              </div>
              
              <div className="flex items-center space-x-4">
                <button
                  onClick={() => setShowAuthModal(true)}
                  className="bg-blue-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-blue-700"
                >
                  Login / Register
                </button>
              </div>
            </div>
          </div>
        </header>
      )}
      
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {authState === 'authenticated' ? (
          // Authenticated user interface
          <div className="space-y-8">
            <Dashboard activeTab={activeTab} setActiveTab={setActiveTab} />
            
            {activeTab === 'create' && (
              <CreateLink onLinkCreated={handleLinkCreated} isAuthenticated={true} />
            )}

            {activeTab === 'list' && (
              <LinkList />
            )}

            {/* API Keys tab is handled within Dashboard component */}
          </div>
        ) : (
          // Public interface (no authentication required)
          <div className="space-y-8">
            <div className="text-center mb-8">
              <h1 className="text-4xl font-bold text-gray-900 mb-4">
                Welcome to URL Shortener
              </h1>
              <p className="text-lg text-gray-600 mb-6">
                Create short, memorable URLs for your long links
              </p>
              <p className="text-sm text-gray-500">
                No registration required! Create short links instantly, or{' '}
                <button
                  onClick={() => setShowAuthModal(true)}
                  className="text-blue-600 hover:text-blue-800 underline"
                >
                  login to manage your links
                </button>
              </p>
            </div>

            {/* Public link creation */}
            <CreateLink onLinkCreated={handleLinkCreated} isAuthenticated={false} />
          </div>
        )}
      </div>

      {/* Authentication Modal */}
      <CenteredModal
        isOpen={showAuthModal}
        onClose={() => setShowAuthModal(false)}
        title="Join URL Shortener"
        maxWidth="md"
      >
        <Auth onAuthSuccess={handleAuthSuccess} />
      </CenteredModal>
    </div>
  );
}

export default App;
