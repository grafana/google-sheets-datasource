import { DataSourceSettings } from '@grafana/data';
import React from 'react';
import { GoogleSheetsSecureJsonData, SheetsSourceOptions } from 'types';
import { JWTConfig } from './JWTConfig';

type Props = {
  options: DataSourceSettings<SheetsSourceOptions, GoogleSheetsSecureJsonData>;
  onOptionsChange: (options: DataSourceSettings<SheetsSourceOptions, GoogleSheetsSecureJsonData>) => void;
};

export function JWTAuth({ options, onOptionsChange }: Props) {
  const { secureJsonFields, secureJsonData } = options;

  return (
    <>
      <JWTConfig
        isConfigured={(secureJsonFields && !!secureJsonFields.jwt) as boolean}
        onChange={(jwt) => {
          onOptionsChange({
            ...options,
            secureJsonData: {
              ...secureJsonData,
              jwt,
            },
          });
        }}
      />
      <div className="grafana-info-box" style={{ marginTop: 24 }}>
        <h4>Generate a JWT file</h4>
        <ol style={{ listStylePosition: 'inside' }}>
          <li>
            Open the <a href="https://console.developers.google.com/apis/credentials">Credentials</a> page in the Google
            API Console.
          </li>
          <li>
            Click <strong>Create Credentials</strong> then click <strong>Service account</strong>.
          </li>
          <li>On the Create service account page, enter the Service account details.</li>
          <li>
            On the <code>Create service account</code> page, fill in the <code>Service account details</code> and then
            click <code>Create</code>
          </li>
          <li>
            On the <code>Service account permissions</code> page, don&rsquo;t add a role to the service account. Just
            click <code>Continue</code>
          </li>
          <li>
            In the next step, click <code>Create Key</code>. Choose key type <code>JSON</code> and click{' '}
            <code>Create</code>. A JSON key file will be created and downloaded to your computer
          </li>
          <li>
            Open the{' '}
            <a href="https://console.cloud.google.com/apis/library/sheets.googleapis.com?q=sheet">Google Sheets</a> in
            API Library and enable access for your account
          </li>
          <li>
            Open the{' '}
            <a href="https://console.cloud.google.com/apis/library/drive.googleapis.com?q=drive">Google Drive</a> in API
            Library and enable access for your account. Access to the Google Drive API is used to list all spreadsheets
            that you have access to.
          </li>
          <li>
            Drag the file to the dotted zone above. Then click <code>Save & Test</code>. The file contents will be
            encrypted and saved in the Grafana database.
          </li>
        </ol>
      </div>
    </>
  );
}
