import {
  DataSourcePluginOptionsEditorProps,
  onUpdateDatasourceSecureJsonDataOption,
} from '@grafana/data';
import { AuthConfig, DataSourceOptions } from '@grafana/google-sdk';
import { Field, SecretInput, Divider, SegmentAsync } from '@grafana/ui';
import React from 'react';
import { GoogleSheetsAuth, GoogleSheetsSecureJSONData, googleSheetsAuthTypes, SheetsSourceOptions } from '../types';
import { getBackwardCompatibleOptions } from '../utils';
import { ConfigurationHelp } from './ConfigurationHelp';
import { getDataSourceSrv } from '@grafana/runtime';
import { DataSource } from '../DataSource';

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
    try {
      const ds = (await getDataSourceSrv().get(options.uid)) as DataSource;
      return ds.getSpreadSheets();
    } catch {
      return [];
    }
  };

  return (
    <>
      <Divider />
      <div className="grafana-info-box">
        <h5>Choosing an authentication type</h5>
        <ul>
          <li><strong>Google JWT File</strong>: provides access to private spreadsheets and works in all environments where Grafana is running.</li> 
          <li><strong>API Key</strong>: simpler configuration, but requires spreadsheets to be public.</li>
          <li><strong>GCE Default Service Account</strong>: automatically retrieves default credentials. Requires Grafana to be running on a Google Compute Engine virtual machine.</li>
        </ul>
        <br/>
        <p><strong>Select an Authentication type below and expand <strong>Configure Google Sheets Authentication</strong> for 
          detailed guidance on configuration</strong>.
        </p>
      </div>
      <ConfigurationHelp authenticationType={options.jsonData.authenticationType} />

      <Divider />

      <AuthConfig authOptions={googleSheetsAuthTypes} onOptionsChange={props.onOptionsChange} options={options} />

      {options.jsonData.authenticationType === GoogleSheetsAuth.API && (
        <Field label="API Key">
          <SecretInput {...apiKeyProps} label="API key" width={40} />
        </Field>
      )}

      <Divider />

      <Field
        label="Default spreadsheet ID"
        description="The ID of a default Google spreadsheet. The datasource must be saved before this can be set."
      >
        <SegmentAsync
          loadOptions={loadSheetIDs}
          placeholder="Select Spreadsheet ID"
          value={(options.jsonData as SheetsSourceOptions).defaultSheetID}
          allowCustomValue={true}
          onChange={(value) => {
            props.onOptionsChange({
              ...options,
              jsonData: {
                ...options.jsonData,
                defaultSheetID: value?.value || value,
              } as DataSourceOptions & { defaultSheetID?: string },
            });
          }}
          disabled={
            (options.jsonData.authenticationType === GoogleSheetsAuth.API &&
              (!options.secureJsonFields || !options.secureJsonFields.apiKey)) ||
            (options.jsonData.authenticationType !== GoogleSheetsAuth.API &&
              (!options.secureJsonFields || !options.secureJsonFields.jwt))
          }
        />
      </Field>
    </>
  );
}
