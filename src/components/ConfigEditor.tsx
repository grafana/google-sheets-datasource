import React, { PureComponent } from 'react';
import { SecretFormField, FormLabel, Select } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps, onUpdateDatasourceSecureJsonDataOption, onUpdateDatasourceJsonDataOptionSelect } from '@grafana/data';
import { SheetsSourceOptions, GoogleSheetsSecureJsonData, GoogleAuthType, googleAuthTypes } from '../types';
import { JWTConfig } from './';

export type Props = DataSourcePluginOptionsEditorProps<SheetsSourceOptions>;

export class ConfigEditor extends PureComponent<Props> {
  onResetApiKey = () => {
    // :( TODO: typings do not let me call the standard function!!!
    // :( updateDatasourcePluginResetOption(this.props, 'apiKey');

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
    // HACK till after: https://github.com/grafana/grafana/pull/21772
    const secureJsonData = options.secureJsonData as GoogleSheetsSecureJsonData;
    return (
      <div className="gf-form-group">
        <div className="gf-form">
          <FormLabel className="width-10">Auth Provider</FormLabel>
          <Select
            className="width-30"
            value={googleAuthTypes.find(x => x.value === jsonData.authType) || googleAuthTypes[0]}
            options={googleAuthTypes}
            defaultValue={options.jsonData.authType}
            onChange={onUpdateDatasourceJsonDataOptionSelect(this.props, 'authType')}
          />
        </div>
        {jsonData.authType === GoogleAuthType.NONE && (
          <>
            <div className="gf-form">
              {console.log({ tjennnna: secureJsonData })}
              <SecretFormField
                isConfigured={(secureJsonFields && secureJsonFields.apiKey) as boolean}
                value={secureJsonData?.apiKey || ''}
                label="API Key"
                labelWidth={10}
                inputWidth={25}
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
                  ...secureJsonData!,
                  jwt,
                },
              });
            }}
          ></JWTConfig>
        )}
      </div>
    );
  }
}
