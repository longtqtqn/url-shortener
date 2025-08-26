import React, { useEffect } from 'react';

interface CenteredModalProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  children: React.ReactNode;
  maxWidth?: 'sm' | 'md' | 'lg' | 'xl';
}

const CenteredModal: React.FC<CenteredModalProps> = ({ 
  isOpen, 
  onClose, 
  title, 
  children, 
  maxWidth = 'md' 
}) => {
  // Get window dimensions for responsive behavior
  const [windowSize, setWindowSize] = React.useState({
    width: typeof window !== 'undefined' ? window.innerWidth : 1024,
    height: typeof window !== 'undefined' ? window.innerHeight : 768,
  });

  React.useEffect(() => {
    const handleResize = () => {
      setWindowSize({
        width: window.innerWidth,
        height: window.innerHeight,
      });
    };

    if (typeof window !== 'undefined') {
      window.addEventListener('resize', handleResize);
      return () => window.removeEventListener('resize', handleResize);
    }
  }, []);

  // Use responsive sizing: 70% on desktop, 90% on mobile
  const isSmallScreen = windowSize.width < 768;
  const modalWidthPercent = isSmallScreen ? 90 : 70;
  const modalHeightPercent = isSmallScreen ? 85 : 70;
  const horizontalOffsetPercent = (100 - modalWidthPercent) / 2;
  
  // Position modal higher on screen (closer to top)
  const topOffsetPercent = isSmallScreen ? 5 : 8; // Start much higher
  // Handle escape key press
  useEffect(() => {
    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === 'Escape' && isOpen) {
        onClose();
      }
    };

    if (isOpen) {
      document.addEventListener('keydown', handleEscape);
    }
    
    return () => document.removeEventListener('keydown', handleEscape);
  }, [isOpen, onClose]);

  // Prevent body scroll when modal is open
  useEffect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
      document.body.style.paddingRight = '0px';
    } else {
      document.body.style.overflow = '';
      document.body.style.paddingRight = '';
    }

    return () => {
      document.body.style.overflow = '';
      document.body.style.paddingRight = '';
    };
  }, [isOpen]);

  if (!isOpen) return null;

  const maxWidthClasses = {
    sm: 'max-w-sm',
    md: 'max-w-md',
    lg: 'max-w-lg',
    xl: 'max-w-xl'
  };

  return (
    <>
      {/* Modal backdrop covering full screen */}
      <div
        style={{
          position: 'fixed',
          top: 0,
          left: 0,
          right: 0,
          bottom: 0,
          zIndex: 9999,
          backgroundColor: 'rgba(0, 0, 0, 0.5)',
          backdropFilter: 'blur(2px)',
        }}
        onClick={onClose}
      />
      
      {/* Modal container - responsive sizing, positioned higher */}
      <div
        style={{
          position: 'fixed',
          top: `${topOffsetPercent}%`,
          left: `${horizontalOffsetPercent}%`,
          width: `${modalWidthPercent}%`,
          height: `${modalHeightPercent}%`,
          zIndex: 10000,
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          pointerEvents: 'none', // Allow clicks to pass through to backdrop
        }}
      >
        <div
          className={`
            bg-white rounded-xl shadow-2xl w-full ${maxWidthClasses[maxWidth]}
            transform transition-all duration-300 max-h-full overflow-y-auto
          `}
          style={{
            position: 'relative',
            maxWidth: '60%',
            maxHeight: '80%',
            pointerEvents: 'auto', // Re-enable clicks for modal content
            backgroundColor: 'rgba(197, 227, 217, 0.7)',
          }}
          onClick={(e) => e.stopPropagation()}
          role="dialog"
          aria-modal="true"
          aria-labelledby="modal-title"
        >
          {/* Modal Header */}
          <div className="flex items-center justify-between border-b border-gray-200 px-6 py-4 bg-white sticky top-0 z-10">
            <div className="flex items-center space-x-3">
              <div className="flex h-10 w-10 items-center justify-center rounded-full bg-blue-600">
                <svg className="h-5 w-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
              </div>
              <h3 
                className="text-lg font-semibold text-gray-900" 
                id="modal-title"
              >
                {title}
              </h3>
            </div>
            
            <button
              type="button"
              className="rounded-full p-3 text-gray-400 hover:bg-gray-100 hover:text-gray-600 focus:outline-none 
                    focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors duration-200 -mr-2 -mt-2"
              onClick={onClose}
              aria-label="Close modal"
            >
              <svg className="h-7 w-7" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          
          {/* Modal Content */}
          <div className="px-6 py-6 bg-white">
            {children}
          </div>
        </div>
      </div>
    </>
  );
};

export default CenteredModal;
