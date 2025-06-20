import { DataSourcePluginOptionsEditorProps, onUpdateDatasourceSecureJsonDataOption } from '@grafana/data';
import { AuthConfig, DataSourceOptions } from '@grafana/google-sdk';
import { Field, SecretInput, Divider, TextLink } from '@grafana/ui';
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
      <div className="grafana-info-box">
        <h5>Choosing an Authentication Type</h5>
        <p>
          <TextLink
            href="https://grafana.com/docs/plugins/grafana-googlesheets-datasource/latest/setup/configure/#authenticate-with-a-service-account-jwt"
            external
          >
            Google JWT File
          </TextLink>
          : Provides access to private spreadsheets and works in all environments where Grafana is running. <br />
          <br />
          <TextLink
            href="https://grafana.com/docs/plugins/grafana-googlesheets-datasource/latest/setup/configure/#authenticate-with-an-api-key"
            external
          >
            API Key
          </TextLink>
          : Simpler configuration, but requires spreadsheets to be public.
          <br />
          <br />
          <TextLink
            href="https://grafana.com/docs/plugins/grafana-googlesheets-datasource/latest/setup/configure/#authenticate-with-the-default-gce-service-account"
            external
          >
            GCE Default Service Account
          </TextLink>
          : Automatically retrieves default credentials. Requires Grafana to be running on a Google Compute Engine
          virtual machine.
          <br />
          <br />
          <strong>
            Select an Authentication type below and expand <strong>Configure Google Sheets Authentication</strong> for
            detailed guidance on configuration
          </strong>
          .
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
    </>
  );
}
