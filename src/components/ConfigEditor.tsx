import React, { PureComponent } from 'react';
import { SecretFormField, FormLabel, Select } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps, onUpdateDatasourceSecureJsonDataOption, onUpdateDatasourceJsonDataOptionSelect } from '@grafana/data';
import { SheetsSourceOptions, GoogleSheetsSecureJsonData, GoogleAuthType, googleAuthTypes } from '../types';
import { JWTConfig } from './';

export type Props = DataSourcePluginOptionsEditorProps<SheetsSourceOptions, GoogleSheetsSecureJsonData>;

export class ConfigEditor extends PureComponent<Props> {
  componentWillMount() {
    // Set the default values
    if (!this.props.options.jsonData.hasOwnProperty('authType')) {
      this.props.options.jsonData.authType = GoogleAuthType.KEY;
    }
  }

  onResetApiKey = () => {
    const { options } = this.props;
    this.props.onOptionsChange({
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

  render() {
    const { options, onOptionsChange } = this.props;
    const { secureJsonFields, jsonData } = options;
    const secureJsonData = options.secureJsonData as GoogleSheetsSecureJsonData;
    return (
      <div className="gf-form-group">
        <div className="gf-form">
          <FormLabel className="width-10">Auth</FormLabel>
          <Select
            className="width-30"
            value={googleAuthTypes.find(x => x.value === jsonData.authType) || googleAuthTypes[0]}
            options={googleAuthTypes}
            defaultValue={options.jsonData.authType}
            onChange={onUpdateDatasourceJsonDataOptionSelect(this.props, 'authType')}
          />
        </div>
        {jsonData.authType === GoogleAuthType.KEY && (
          <>
            <div className="gf-form">
              <SecretFormField
                isConfigured={(secureJsonFields && secureJsonFields.apiKey) as boolean}
                value={secureJsonData?.apiKey || ''}
                label="API Key"
                labelWidth={10}
                inputWidth={30}
                placeholder="Enter API Key"
                onReset={this.onResetApiKey}
                onChange={onUpdateDatasourceSecureJsonDataOption(this.props, 'apiKey')}
              />
            </div>
          </>
        )}
        {jsonData.authType === GoogleAuthType.JWT && (
          <JWTConfig
            isConfigured={(secureJsonFields && !!secureJsonFields.jwt) as boolean}
            onChange={jwt => {
              onOptionsChange({
                ...options,
                secureJsonData: {
                  ...secureJsonData,
                  jwt,
                },
              });
            }}
          ></JWTConfig>
        )}
        <div className="grafana-info-box" style={{ marginTop: 24 }}>
          {jsonData.authType === GoogleAuthType.JWT ? (
            <>
              <h4>Generate a JWT file</h4>
              <ol style={{ listStylePosition: 'inside' }}>
                <li>
                  Open the <a href="https://console.developers.google.com/apis/credentials">Credentials</a> page in the Google API Console.
                </li>
                <li>
                  Click <strong>Create Credentials</strong> then click <strong>Service account</strong>.
                </li>
                <li>On the Create service account page, enter the Service account details.</li>
                <li>
                  On the <code>Create service account</code> page, fill in the <code>Service account details</code> and then click <code>Create</code>
                </li>
                <li>
                  On the <code>Service account permissions</code> page, don't add a role to the service account. Just click <code>Continue</code>
                </li>
                <li>
                  In the next step, click <code>Create Key</code>. Choose key type <code>JSON</code> and click <code>Create</code>. A JSON key file
                  will be created and downloaded to your computer
                </li>
                <li>
                  Open the <a href="https://console.cloud.google.com/apis/library/sheets.googleapis.com?q=sheet">Google Sheets</a> in API Library and
                  enable access for your account
                </li>
                <li>
                  Open the <a href="https://console.cloud.google.com/apis/library/drive.googleapis.com?q=drive">Google Drive</a> in API Library and
                  enable access for your account. Access to the Google Drive API is used to list all spreadsheets that you have access to.
                </li>
                <li>
                  Drag the file to the dotted zone above. Then click <code>Save & Test</code>. The file contents will be encrypted and saved in the
                  Grafana database.
                </li>
              </ol>
            </>
          ) : (
            <>
              <h4>Generate an API key</h4>
              <ol style={{ listStylePosition: 'inside' }}>
                <li>
                  Open the <a href="https://console.developers.google.com/apis/credentials">Credentials page</a> in the Google API Console.
                </li>
                <li>
                  Click <strong>Create Credentials</strong> and then click <strong>API key</strong>.
                </li>
                <li>Copy the key and paste it in the API Key field above. The file contents are encrypted and saved in the Grafana database.</li>
              </ol>
            </>
          )}
        </div>
      </div>
    );
  }
}
