import { DataSourceSettings, onUpdateDatasourceSecureJsonDataOption } from '@grafana/data';
import { LegacyForms } from '@grafana/ui';
import React from 'react';
import { SheetsSourceOptions, GoogleSheetsSecureJsonData } from 'types';

type Props = {
  options: DataSourceSettings<SheetsSourceOptions, GoogleSheetsSecureJsonData>;
  onOptionsChange: (options: DataSourceSettings<SheetsSourceOptions, GoogleSheetsSecureJsonData>) => void;
};

export function APIAuth(props: Props) {
  const { options, onOptionsChange } = props;
  const onResetApiKey = () => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        ...options.secureJsonData,
        apiKey: '',
      },
      secureJsonFields: {
        ...options.secureJsonFields,
        apiKey: false,
      },
    });
  };

  return (
    <>
      <div className="gf-form">
        <LegacyForms.SecretFormField
          isConfigured={options.secureJsonFields?.apiKey}
          value={options.secureJsonData?.apiKey || ''}
          label="API Key"
          labelWidth={10}
          inputWidth={30}
          placeholder="Enter API Key"
          onReset={onResetApiKey}
          onChange={onUpdateDatasourceSecureJsonDataOption(props, 'apiKey')}
        />
      </div>
      <div className="grafana-info-box" style={{ marginTop: 24 }}>
        <h4>Generate an API key</h4>
        <ol style={{ listStylePosition: 'inside' }}>
          <li>
            Open the <a href="https://console.developers.google.com/apis/credentials">Credentials page</a> in the Google
            API Console.
          </li>
          <li>
            Click <strong>Create Credentials</strong> and then click <strong>API key</strong>.
          </li>
          <li>
            Copy the key and paste it in the API Key field above. The file contents are encrypted and saved in the
            Grafana database.
          </li>
        </ol>
      </div>
    </>
  );
}
