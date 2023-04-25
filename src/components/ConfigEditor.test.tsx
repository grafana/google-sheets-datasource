import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { ConfigEditor } from './ConfigEditor';
import {
  DataSourcePluginOptionsEditorProps,
  DataSourceSettings,
  onUpdateDatasourceJsonDataOptionSelect,
} from '@grafana/data';
import { SheetsSourceOptions, GoogleSheetsSecureJsonData, GoogleAuthType } from '../types';
import userEvent from '@testing-library/user-event';
import exp from 'constants';

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
  it('should render default spreadsheet ID field', async () => {
    const onChange = jest.fn();
    const props = {
      options: dataSourceSettings,
      onOptionsChange: onChange,
    } as DataSourcePluginOptionsEditorProps<SheetsSourceOptions, GoogleSheetsSecureJsonData>;
    render(<ConfigEditor {...props} onOptionsChange={onChange} />);
    expect(screen.getByText('Default spreadsheet ID')).toBeInTheDocument();
  });

  it('should update default spreadsheet after selecting it', async () => {
    const onChange = jest.fn();
    const props = {
      options: dataSourceSettings,
      onOptionsChange: onChange,
    } as DataSourcePluginOptionsEditorProps<SheetsSourceOptions, GoogleSheetsSecureJsonData>;

    render(<ConfigEditor {...props} onOptionsChange={onChange} />);

    const selectEl = screen.getByText('Select Spreadsheet ID');
    expect(selectEl).toBeInTheDocument();
    userEvent.click(selectEl);
    const spreadsheetOption = await screen.findByText('label1');
    userEvent.click(spreadsheetOption);
    await waitFor(() => expect(mockedSelect).toHaveBeenCalledWith({ label: 'label1', value: 'value1' }));
  });
});
