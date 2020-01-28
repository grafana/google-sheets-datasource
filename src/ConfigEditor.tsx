import React, { PureComponent } from 'react';
import { SecretFormField, FormField } from '@grafana/ui';
import { DataSourcePluginOptionsEditorProps, onUpdateDatasourceSecureJsonDataOption } from '@grafana/data';
import { SheetsSourceOptions, GoogleSheetsSecureJsonData } from './types';

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
            value={jsonData.apiKey}
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
      </div>
    );
  }
}
