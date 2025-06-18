import { DataSourcePluginOptionsEditorProps, onUpdateDatasourceSecureJsonDataOption } from '@grafana/data';
import { AuthConfig, DataSourceOptions } from '@grafana/google-sdk';
import { Field, SecretInput, Divider } from '@grafana/ui';
import React from 'react';
import { GoogleSheetsAuth, GoogleSheetsSecureJSONData, googleSheetsAuthTypes } from '../types';
import { getBackwardCompatibleOptions } from '../utils';
import { ConfigurationHelp } from './ConfigurationHelp';
import { DataSourceDescription } from '@grafana/plugin-ui';

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

  return (
    <>
      <DataSourceDescription
        dataSourceName="Google Sheets"
        docsLink="https://grafana.com/docs/plugins/grafana-googlesheets-datasource/latest/"
        hasRequiredFields={false}
      />

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
