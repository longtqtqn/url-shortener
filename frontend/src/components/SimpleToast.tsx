import { useEffect, useState } from 'react';

interface SimpleToastProps {
  message: string;
  type?: 'error' | 'success' | 'info';
  onClose: () => void;
  duration?: number;
}

const SimpleToast: React.FC<SimpleToastProps> = ({ 
  message, 
  type = 'error', 
  onClose, 
  duration = 4000 
}) => {
  const [isVisible, setIsVisible] = useState(false);

  useEffect(() => {
    // Show toast
    setTimeout(() => setIsVisible(true), 10);
    
    // Auto hide
    const timer = setTimeout(() => {
      setIsVisible(false);
      setTimeout(onClose, 300);
    }, duration);

    return () => clearTimeout(timer);
  }, [duration, onClose]);

  const getStyles = () => {
    switch (type) {
      case 'success':
        return 'bg-green-100 border-green-400 text-green-700';
      case 'info':
        return 'bg-blue-100 border-blue-400 text-blue-700';
      default:
        return 'bg-red-100 border-red-400 text-red-700';
    }
  };

  const getIcon = () => {
    switch (type) {
      case 'success':
        return '✅';
      case 'info':
        return 'ℹ️';
      default:
        return '❌';
    }
  };

  return (
    <div
      className={`fixed top-4 right-4 z-50 max-w-sm p-4 border rounded-lg shadow-lg transition-all duration-300 ${
        isVisible ? 'opacity-100 translate-x-0' : 'opacity-0 translate-x-full'
      } ${getStyles()}`}
    >
      <div className="flex items-start">
        <span className="mr-2 text-lg">{getIcon()}</span>
        <div className="flex-1">
          <p className="text-sm font-medium">{message}</p>
        </div>
        <button
          onClick={() => {
            setIsVisible(false);
            setTimeout(onClose, 300);
          }}
          className="ml-2 text-gray-500 hover:text-gray-700 text-lg leading-none"
        >
          ×
        </button>
      </div>
    </div>
  );
};

export default SimpleToast;
