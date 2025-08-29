import { ExclamationTriangleIcon, XCircleIcon, InformationCircleIcon } from '@heroicons/react/24/outline';

interface ErrorDisplayProps {
  error: string;
  type?: 'error' | 'warning' | 'info';
  className?: string;
  onDismiss?: () => void;
}

const ErrorDisplay: React.FC<ErrorDisplayProps> = ({ 
  error, 
  type = 'error', 
  className = '',
  onDismiss 
}) => {
  if (!error) return null;

  const getIcon = () => {
    switch (type) {
      case 'warning':
        return <ExclamationTriangleIcon className="h-5 w-5 text-yellow-500" />;
      case 'info':
        return <InformationCircleIcon className="h-5 w-5 text-blue-500" />;
      default:
        return <XCircleIcon className="h-5 w-5 text-red-500" />;
    }
  };

  const getStyles = () => {
    switch (type) {
      case 'warning':
        return 'bg-yellow-50 border-yellow-200 text-yellow-800';
      case 'info':
        return 'bg-blue-50 border-blue-200 text-blue-800';
      default:
        return 'bg-red-50 border-red-200 text-red-800';
    }
  };

  // Check if error contains multiple lines (validation errors)
  const isMultiline = error.includes('\n');

  return (
    <div className={`border rounded-md p-4 ${getStyles()} ${className}`}>
      <div className="flex items-start">
        <div className="flex-shrink-0">
          {getIcon()}
        </div>
        <div className="ml-3 flex-1">
          {isMultiline ? (
            <div className="whitespace-pre-line text-sm">
              {error}
            </div>
          ) : (
            <div className="text-sm">
              {error}
            </div>
          )}
        </div>
        {onDismiss && (
          <div className="ml-auto pl-3">
            <div className="-mx-1.5 -my-1.5">
              <button
                type="button"
                onClick={onDismiss}
                className={`inline-flex rounded-md p-1.5 focus:outline-none focus:ring-2 focus:ring-offset-2 ${
                  type === 'warning' 
                    ? 'text-yellow-500 hover:bg-yellow-100 focus:ring-yellow-600' 
                    : type === 'info'
                    ? 'text-blue-500 hover:bg-blue-100 focus:ring-blue-600'
                    : 'text-red-500 hover:bg-red-100 focus:ring-red-600'
                }`}
              >
                <XCircleIcon className="h-5 w-5" />
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default ErrorDisplay;
