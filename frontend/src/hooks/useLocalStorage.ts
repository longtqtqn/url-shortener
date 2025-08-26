import { useState, useEffect } from "react";

export default function useLocalStorage(key: string, initialValue: string) {
  const [storedValue, setStoredValue] = useState(() => {
    try {
      const item = localStorage.getItem(key);
      return item ?? initialValue;
    } catch {
      return initialValue;
    }
  });

  useEffect(() => {
    try {
      localStorage.setItem(key, storedValue);
    } catch {}
  }, [key, storedValue]);

  return [storedValue, setStoredValue] as const;
}
