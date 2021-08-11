import { useEffect } from 'react';

export function useLoadGapi(onLoad?: () => void) {
  useEffect(() => {
    if (window.gapi) {
      return;
    }
    const script = document.createElement('script');
    script.src = 'https://apis.google.com/js/api.js';
    script.async = true;
    script.onload = onLoad || null;
    document.body.appendChild(script);
    return () => {
      document.body.removeChild(script);
    };
  }, [onLoad]);
}
