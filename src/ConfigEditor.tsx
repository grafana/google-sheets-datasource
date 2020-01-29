import React, { PureComponent } from 'react';
import { SecretFormField, FormField, FormLabel, Select } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps, onUpdateDatasourceSecureJsonDataOption, onUpdateDatasourceJsonDataOptionSelect } from '@grafana/data';
import { SheetsSourceOptions, GoogleSheetsSecureJsonData, GoogleAuthType, googleAuthTypes } from './types';

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
    const { options } = this.props;
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
            <div className="gf-form">
              <FormField
                label="API Key (temp)"
                labelWidth={10}
                inputWidth={25}
                placeholder="Enter API Key"
                onChange={e =>
                  this.props.onOptionsChange({
                    ...this.props,
                    ...this.props.options,
                    jsonData: {
                      ...this.props.options.jsonData,
                      apiKey: e.target.value,
                    },
                  })
                }
              />
            </div>
          </>
        )}
        {jsonData.authType === GoogleAuthType.JWT && (
          <textarea
            onPaste={e => onUpdateDatasourceJsonDataOptionSelect(this.props, 'jwtFile')({ value: e.clipboardData.getData('text/plain') })}
            placeholder="Paste your Google JWT file content here"
            className="gf-form-input"
            style={{ height: 200 }}
          >
            {jsonData.jwtFile}
          </textarea>
        )}
      </div>
    );
  }
}
