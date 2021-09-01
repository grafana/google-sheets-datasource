import { useState, useEffect } from 'react';
import { useLoadGapi } from './useLoadGapi';

type UsePickerOptions = {
  pickerCallback?: (data: google.picker.ResponseObject) => void;
  appId?: string;
  apiKey?: string;
  clientId?: string;
};

export function usePicker({ pickerCallback, apiKey, appId, clientId }: UsePickerOptions) {
  const [picker, setPicker] = useState<google.picker.Picker>();
  const [isPickerLoaded, setIsPickerLoaded] = useState(false);
  const [isSignedIn, setIsSignedIn] = useState(false);

  useLoadGapi(() => {
    gapi.load('client', {
      callback: () => {
        gapi.client
          .init({
            apiKey,
            clientId: clientId,
            scope: 'https://www.googleapis.com/auth/drive.file',
          })
          .then(() => {
            setIsSignedIn(gapi.auth2.getAuthInstance().isSignedIn.get());
          });
      },
    });
    gapi.load('picker', {
      callback: () => {
        setIsPickerLoaded(true);
      },
    });
  }, Boolean(apiKey && clientId));

  // Create picker
  useEffect(() => {
    if (isPickerLoaded && isSignedIn && appId && apiKey && pickerCallback) {
      const { access_token } = gapi.auth2.getAuthInstance().currentUser.get().getAuthResponse();
      const pickerBuilder = new google.picker.PickerBuilder()
        .setAppId(appId)
        .setOAuthToken(access_token)
        .addView(google.picker.ViewId.SPREADSHEETS)
        .setDeveloperKey(apiKey)
        .setCallback(pickerCallback)
        .build();

      setPicker(pickerBuilder);
    }
  }, [apiKey, appId, isPickerLoaded, isSignedIn, pickerCallback]);

  // Dispose picker
  useEffect(() => {
    return () => {
      picker?.dispose();
    };
  }, [picker]);

  const openPicker = () => {
    picker?.setVisible(true);
  };

  return { openPicker, isSignedIn };
}
