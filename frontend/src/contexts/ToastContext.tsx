import React, { createContext, useContext } from 'react';

interface ToastContextType {
  showError: (message: string) => void;
  showSuccess: (message: string) => void;
  showInfo: (message: string) => void;
}

const ToastContext = createContext<ToastContextType | undefined>(undefined);

export const useToast = () => {
  const context = useContext(ToastContext);
  if (context === undefined) {
    throw new Error('useToast must be used within a ToastProvider');
  }
  return context;
};

interface ToastProviderProps {
  children: React.ReactNode;
  showError: (message: string) => void;
  showSuccess: (message: string) => void;
  showInfo: (message: string) => void;
}

export const ToastProvider: React.FC<ToastProviderProps> = ({ 
  children, 
  showError, 
  showSuccess, 
  showInfo 
}) => {
  return (
    <ToastContext.Provider value={{ showError, showSuccess, showInfo }}>
      {children}
    </ToastContext.Provider>
  );
};
