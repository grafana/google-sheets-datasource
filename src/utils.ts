import { GoogleSheetsAuth } from './types';
import { Props } from './components/ConfigEditor';

export function getBackwardCompatibleOptions(options: Props['options']): Props['options'] {
  const changedOptions = {
    ...options,
    jsonData: { ...options.jsonData },
    secureJsonFields: { ...options.secureJsonFields },
  };
  // Make sure we support the old authType property
  changedOptions.jsonData.authenticationType = options.jsonData.authenticationType || options.jsonData.authType!;

  // Show a configured message for the JWT fields when JWT is set in secureJsonFields
  if (changedOptions.jsonData.authenticationType === GoogleSheetsAuth.JWT && options.secureJsonFields?.jwt) {
    changedOptions.jsonData.clientEmail = 'configured';
    changedOptions.jsonData.defaultProject = 'configured';
    changedOptions.jsonData.tokenUri = 'configured';
    changedOptions.secureJsonFields.privateKey = true;
  }

  return changedOptions;
}
