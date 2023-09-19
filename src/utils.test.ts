import { GoogleSheetsAuth } from './types';
import { getBackwardCompatibleOptions } from './utils';

describe('getBackwardCompatibleOptions', () => {
  it('should not mutate the option object', () => {
    const options: any = Object.freeze({
      jsonData: Object.freeze({
        authenticationType: GoogleSheetsAuth.JWT,
      }),
      secureJsonFields: Object.freeze({}),
    });
    expect(getBackwardCompatibleOptions(options)).toEqual(options);
  });

  it('should set authenticationType to authType if authType is set', () => {
    const options: any = {
      jsonData: {
        authType: GoogleSheetsAuth.API,
      },
      secureJsonFields: {},
    };
    const expectedOptions = {
      jsonData: {
        authenticationType: GoogleSheetsAuth.API,
        authType: GoogleSheetsAuth.API,
      },
      secureJsonFields: {},
    };
    expect(getBackwardCompatibleOptions(options)).toEqual(expectedOptions);
  });

  it('should set JWT fields to "configured" if JWT is set in secureJsonFields', () => {
    const options: any = {
      jsonData: {
        authenticationType: GoogleSheetsAuth.JWT,
        clientEmail: '',
        defaultProject: '',
        tokenUri: '',
      },
      secureJsonFields: {
        jwt: true,
      },
    };
    const expectedOptions = {
      jsonData: {
        authenticationType: GoogleSheetsAuth.JWT,
        clientEmail: 'configured',
        defaultProject: 'configured',
        tokenUri: 'configured',
      },
      secureJsonFields: {
        jwt: true,
        privateKey: true,
      },
    };
    expect(getBackwardCompatibleOptions(options)).toEqual(expectedOptions);
  });
});
