import { DataSourcePluginOptionsEditorProps, onUpdateDatasourceSecureJsonDataOption } from '@grafana/data';
import { AuthConfig, DataSourceOptions } from '@grafana/google-sdk';
import { Field, SecretInput, Divider } from '@grafana/ui';
import { t } from '@grafana/i18n';
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
    placeholder: t('config.apiKey.placeholder', 'Enter API key'),
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
        <h5>{t('config.authType.title', 'Choosing an authentication type')}</h5>
        <ul>
          <li>
            <strong>{t('config.authType.jwt.label', 'Google JWT File')}</strong>:{' '}
            {t(
              'config.authType.jwt.description',
              'provides access to private spreadsheets and works in all environments where Grafana is running.'
            )}
          </li>
          <li>
            <strong>{t('config.authType.apiKey.label', 'API Key')}</strong>:{' '}
            {t('config.authType.apiKey.description', 'simpler configuration, but requires spreadsheets to be public.')}
          </li>
          <li>
            <strong>{t('config.authType.gce.label', 'GCE Default Service Account')}</strong>:{' '}
            {t(
              'config.authType.gce.description',
              'automatically retrieves default credentials. Requires Grafana to be running on a Google Compute Engine virtual machine.'
            )}
          </li>
        </ul>
        <br />
        <p>
          <strong>
            {t(
              'config.authType.instructions',
              'Select an Authentication type below and expand Configure Google Sheets Authentication for detailed guidance on configuration'
            )}
            .
          </strong>
        </p>
      </div>
      <ConfigurationHelp authenticationType={options.jsonData.authenticationType} />

      <Divider />

      <AuthConfig authOptions={googleSheetsAuthTypes} onOptionsChange={props.onOptionsChange} options={options} />

      {options.jsonData.authenticationType === GoogleSheetsAuth.API && (
        <Field label={t('config.apiKey.label', 'API Key')}>
          <SecretInput {...apiKeyProps} label={t('config.apiKey.label', 'API key')} width={40} />
        </Field>
      )}
    </>
  );
}
