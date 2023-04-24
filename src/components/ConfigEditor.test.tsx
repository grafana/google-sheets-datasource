import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { ConfigEditor } from './ConfigEditor';
import { DataSourcePluginOptionsEditorProps, DataSourceSettings } from '@grafana/data';
import { SheetsSourceOptions, GoogleSheetsSecureJsonData, GoogleAuthType } from '../types';

const dataSourceSettings: DataSourceSettings<SheetsSourceOptions, GoogleSheetsSecureJsonData> = {
  jsonData: {
    authType: GoogleAuthType.JWT,
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
    get: Promise.resolve({
      getSpreadSheets: jest.fn().mockImplementation(() =>
        Promise.resolve([
          { label: 'label1', value: 'value1' },
          { label: 'label2', value: 'value2' },
        ])
      ),
    }),
  }),
}));

describe('ConfigEditor', () => {
  it('should render default spreadsheet ID field', async () => {
    const onChange = jest.fn();
    const props = {
      options: dataSourceSettings,
      onOptionsChange: onChange,
    } as DataSourcePluginOptionsEditorProps<SheetsSourceOptions, GoogleSheetsSecureJsonData>;
    render(<ConfigEditor {...props} onOptionsChange={onChange} />);
    expect(screen.getByText('Default spreadsheet ID')).toBeInTheDocument();
  });

  it('should display available spreadsheets in selector', async () => {
    const onChange = jest.fn();
    const props = {
      options: dataSourceSettings,
      onOptionsChange: onChange,
    } as DataSourcePluginOptionsEditorProps<SheetsSourceOptions, GoogleSheetsSecureJsonData>;

    render(<ConfigEditor {...props} onOptionsChange={onChange} />);

    waitFor(() => {
      const selectEl = screen.getByText('Default spreadsheet ID');
      expect(selectEl).not.toBeInTheDocument();
    });
  });
});
