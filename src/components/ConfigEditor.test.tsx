import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import { ConfigEditor } from './ConfigEditor';
import { DataSourcePluginOptionsEditorProps, DataSourceSettings } from '@grafana/data';
import { SheetsSourceOptions, GoogleSheetsSecureJsonData, GoogleAuthType } from '../types';

//
// const datasource = new DataSource(DataSourceSettings);
// const loadDataSourceMock = jest.fn();
// jest.mock('app/features/plugins/datasource_srv', () => ({
//   getDatasourceSrv: () => ({
//     loadDatasource: loadDataSourceMock,
//   }),
// }));

// const putMock = jest.fn();
// const getMock = jest.fn();
//

// describe('Render', () => {
//   beforeEach(() => {
//     (window as any).grafanaBootData = {
//       settings: {},
//     };
//     jest.resetAllMocks();
//     putMock.mockImplementation(async () => ({datasource: setupMockedDataSource().datasource}));
//     getMock.mockImplementation(async () => ({datasource: setupMockedDataSource().datasource}));
//     loadDataSourceMock.mockResolvedValue(datasource);
//   });
//   it('should display log group selector field', async () => {
//     setup();
//     await waitFor(async () => expect(await screen.getByText('Select Log Groups')).toBeInTheDocument());
//   });
// });

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
};

// jsonData.authType === GoogleAuthType.JWT && (!secureJsonFields || !secureJsonFields.jwt);

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

  it('some name', async () => {
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

    // userEvent.click(selectEl);

    //
    //   const selectEl = screen.getByLabelText(selectors.components.ConfigEditor.workgroup.input);
    //   expect(selectEl).toBeInTheDocument();
    //
    //   await select(selectEl, resourceName, { container: document.body });
    //
    //   expect(onChange).toHaveBeenCalledWith({
    //     ...props.options,
    //     jsonData: { ...props.options.jsonData, workgroup: resourceName },
    //   });
  });

  // it('should save and request spreadsheets', async () => {
  //   const onChange = jest.fn();
  //   render(<ConfigEditor {...props} onOptionsChange={onChange} />);
  //
  //   const d = screen.getByTestId(selectors.components.ConfigEditor.workgroup.wrapper);
  //   expect(d).toBeInTheDocument();
  //   d.click();
  //
  //   const selectEl = screen.getByLabelText(selectors.components.ConfigEditor.workgroup.input);
  //   expect(selectEl).toBeInTheDocument();
  //
  //   await select(selectEl, resourceName, { container: document.body });
  //
  //   expect(onChange).toHaveBeenCalledWith({
  //     ...props.options,
  //     jsonData: { ...props.options.jsonData, workgroup: resourceName },
  //   });
  // });
});
