import { useState } from 'react';
import { useLoadGapi } from './useLoadGapi';

export function usePicker({ pickerCallback }: { pickerCallback: (data: google.picker.ResponseObject) => void }) {
  useLoadGapi();
  const [oAuthToken, setOAuthToken] = useState<string>();
  const [isPickerApiLoaded, setIsPickerApiLoaded] = useState(false);

  function handleAuthResult(authResult: GoogleApiOAuth2TokenObject) {
    if (authResult && !authResult.error) {
      setOAuthToken(authResult.access_token);
      createPicker();
    }
  }

  function createPicker() {
    if (isPickerApiLoaded && oAuthToken) {
      const view = new google.picker.DocsView(google.picker.ViewId.SPREADSHEETS);
      const picker = new google.picker.PickerBuilder()
        .enableFeature(google.picker.Feature.MINE_ONLY)
        .setAppId('appid')
        .setOAuthToken(oAuthToken)
        .addView(view)
        .setDeveloperKey('developerkey')
        .setCallback(pickerCallback)
        .build();
      picker.setVisible(true);
    }
  }

  const openPicker = () => {
    if (!window.gapi) {
      throw new Error('Google api is not loaded.');
    }
    gapi.load('auth', {
      callback: () => {
        window.gapi.auth.authorize(
          {
            client_id: 'clientid',
            scope: ['https://www.googleapis.com/auth/drive.file'],
            immediate: false,
          },
          handleAuthResult
        );
      },
    });
    gapi.load('picker', {
      callback: () => {
        setIsPickerApiLoaded(true);
      },
    });
  };

  return { openPicker };
}
