import { DataSourcePluginOptionsEditorProps, onUpdateDatasourceSecureJsonDataOption } from '@grafana/data';
import { AuthConfig, DataSourceOptions } from '@grafana/google-sdk';
import { Field, SecretInput } from '@grafana/ui';
import React from 'react';
import { GoogleSheetsAuth, GoogleSheetsSecureJSONData, googleSheetsAuthTypes } from '../types';
import { getBackwardCompatibleOptions } from '../utils';
import { ConfigurationHelp } from './ConfigurationHelp';
import { DataSourceDescription } from '@grafana/experimental';
import { Divider } from './Divider';

export type Props = DataSourcePluginOptionsEditorProps<DataSourceOptions, GoogleSheetsSecureJSONData>;

export function ConfigEditor(props: Props) {
  const options = getBackwardCompatibleOptions(props.options);

  const apiKeyProps = {
    isConfigured: Boolean(options.secureJsonFields.apiKey),
    value: options.secureJsonData?.apiKey || '',
    placeholder: 'Enter API key',
    id: 'apiKey',
    onReset: () =>
      props.onOptionsChange({
        ...options,
        secureJsonFields: { ...options.secureJsonFields, apiKey: false },
        secureJsonData: { apiKey: '' },
        jsonData: options.jsonData,
      }),
    onChange: onUpdateDatasourceSecureJsonDataOption(props, 'apiKey'),
  };

  const loadSheetIDs = async () => {
    const { options } = props;
    try {
      const ds = (await getDataSourceSrv().get(options.uid)) as DataSource;
      return ds.getSpreadSheets();
    } catch {
      return [];
    }
  };

  return (
    <>
      <DataSourceDescription
        dataSourceName="Google Sheets"
        docsLink="https://grafana.com/grafana/plugins/grafana-googlesheets-datasource/"
        hasRequiredFields={false}
      />

      <div className="gf-form">
        <InlineFormLabel
          className="width-10"
          tooltip="The id of a default google sheet. The datasource must be saved before this can be set."
        >
          Default SheetID
        </InlineFormLabel>
        <SegmentAsync
          className="width-30"
          loadOptions={() => loadSheetIDs()}
          placeholder="Select Spreadsheet ID"
          value={options.jsonData.defaultSheetID}
          allowCustomValue={true}
          onChange={onUpdateDatasourceJsonDataOptionSelect(props, 'defaultSheetID')}
          disabled={
            (props.options.jsonData.authType === googleSheetsAuthTypes && (!secureJsonFields || !secureJsonFields.apiKey)) ||
            (jsonData.authType === googleSheetsAuthTypes.JWT && (!secureJsonFields || !secureJsonFields.jwt))
          }
        />
      </div>

      <Divider />

      <ConfigurationHelp authenticationType={options.jsonData.authenticationType} />

      <Divider />

      <AuthConfig authOptions={googleSheetsAuthTypes} onOptionsChange={props.onOptionsChange} options={options} />

      {options.jsonData.authenticationType === GoogleSheetsAuth.API && (
        <Field label="API Key">
          <SecretInput {...apiKeyProps} label="API key" width={40} />
        </Field>
      )}
    </>
  );
}
