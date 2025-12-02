import React from 'react';
import { render, screen, fireEvent, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { ConfigEditor } from './ConfigEditor';
import { DataSourceSettings } from '@grafana/data';
import { GoogleSheetsSecureJSONData } from '../types';
import { GoogleAuthType, DataSourceOptions } from '@grafana/google-sdk';

jest.mock('@grafana/plugin-ui', () => ({
  DataSourceDescription: ({ dataSourceName }: any) => <div>{dataSourceName}</div>,
}));

const dataSourceSettings: DataSourceSettings<DataSourceOptions, GoogleSheetsSecureJSONData> = {
  jsonData: {
    authenticationType: GoogleAuthType.JWT,
  },
  secureJsonFields: {
    jwt: true,
  },
  id: 0,
  uid: '',
  type: '',
  access: 'direct',
  name: 'Google Sheets Test Datasource',
  basicAuth: false,
  basicAuthUser: '',
  database: '',
  isDefault: false,
  orgId: 0,
  readOnly: false,
  secureJsonData: undefined,
  typeLogoUrl: '',
  typeName: '',
  url: '',
  user: '',
  withCredentials: false,
};

jest.mock('@grafana/runtime', () => ({
  getDataSourceSrv: () => ({
    get: (_: string) =>
      Promise.resolve({
        getSpreadSheets: () =>
          Promise.resolve([
            { label: 'label1', value: 'value1' },
            { label: 'label2', value: 'value2' },
          ]),
      }),
  }),
}));

const mockedSelect = jest.fn();
jest.mock('@grafana/data', () => ({
  ...jest.requireActual('@grafana/data'),
  // onUpdateDatasourceJsonDataOptionSelect will be called anyway on mount, and twice, one for spreadsheetId and one for authType.
  // The one we are actually testing for is the returned function (in this case mockedSelect) and checking that it is called with the value of the selected option
  onUpdateDatasourceJsonDataOptionSelect: jest.fn(() => mockedSelect),
}));

describe('ConfigEditor', () => {
  afterEach(() => {
    jest.clearAllMocks();
  });

  it('should support old authType property', () => {
    const onOptionsChange = jest.fn();
    // Render component with old authType property
    render(
      <ConfigEditor
        onOptionsChange={onOptionsChange}
        options={{ jsonData: { authType: 'key', authenticationType: '' }, secureJsonFields: {} } as any}
      />
    );

    // Check that the correct auth type is selected
    expect(screen.getByRole('radio', { name: 'API Key' })).toBeChecked();

    // Make sure that the user can still change the auth type
    fireEvent.click(screen.getByLabelText('Google JWT File'));

    // Check onOptionsChange is called with the correct value
    expect(onOptionsChange).toHaveBeenCalledWith({
      jsonData: { authType: 'key', authenticationType: 'jwt' },
      secureJsonFields: {},
    });
  });

  it('should be backward compatible with API Key', () => {
    render(
      <ConfigEditor
        onOptionsChange={jest.fn()}
        options={{ jsonData: { authType: 'key', authenticationType: '' }, secureJsonFields: { apiKey: true } } as any}
      />
    );

    // Check that the correct auth type is selected
    expect(screen.getByRole('radio', { name: 'API Key' })).toBeChecked();

    // Check that the API key is configured
    expect(screen.getByPlaceholderText('Enter API key')).toHaveAttribute('value', 'configured');
  });

  it('should be backward compatible with JWT auth type', () => {
    render(
      <ConfigEditor
        onOptionsChange={jest.fn()}
        options={{ jsonData: { authType: 'jwt', authenticationType: '' }, secureJsonFields: { jwt: true } } as any}
      />
    );

    // Check that the correct auth type is selected
    expect(screen.getByLabelText('Google JWT File')).toBeChecked();

    // Check that the Private key input is configured
    expect(screen.getByTestId('Private Key Input')).toHaveAttribute('value', 'configured');
  });
  it('should render default spreadsheet ID field', async () => {
    const onChange = jest.fn();
    render(
      <ConfigEditor
        onOptionsChange={onChange}
        options={{ jsonData: { authenticationType: 'key' }, secureJsonFields: {} } as any}
      />
    );
    expect(screen.getByText('Default Spreadsheet ID')).toBeInTheDocument();
  });

  it('should display default spreadsheet ID when set', () => {
    const onChange = jest.fn();
    render(
      <ConfigEditor
        onOptionsChange={onChange}
        options={{
          jsonData: { authenticationType: 'jwt', defaultSheetID: '1234567890' },
          secureJsonFields: {},
        } as any}
      />
    );
    expect(screen.getByText('Default Spreadsheet ID')).toBeInTheDocument();
  });
});
