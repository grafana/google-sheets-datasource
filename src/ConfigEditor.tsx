import React, { PureComponent } from 'react';
import { FormLabel, Input, Button } from '@grafana/ui';
import {
  DataSourcePluginOptionsEditorProps,
  onUpdateDatasourceResetOption,
  onUpdateDatasourceSecureJsonDataOption,
} from '@grafana/data';
import { SheetsSourceOptions, GoogleSheetsSecureJsonData } from './types';

export type Props = DataSourcePluginOptionsEditorProps<SheetsSourceOptions, GoogleSheetsSecureJsonData>;

export class ConfigEditor extends PureComponent<Props> {
  render() {
    const {options} = this.props;
    const {secureJsonData, secureJsonFields} = options;

    return <div className="gf-form-group">

{secureJsonFields?.apiKey ? (
                <div className="gf-form-inline">
                  <div className="gf-form">
                    <FormLabel className="width-14">Secret Access Key</FormLabel>
                    <Input className="width-25" placeholder="Configured" disabled={true} />
                  </div>
                  <div className="gf-form">
                    <div className="max-width-30 gf-form-inline">
                      <Button
                        variant="secondary"
                        type="button"
                        onClick={onUpdateDatasourceResetOption(this.props, 'apiKey')}
                      >
                        Reset
                      </Button>
                    </div>
                  </div>
                </div>
              ) : (
                <div className="gf-form-inline">
                  <div className="gf-form">
                    <FormLabel className="width-14">API Key</FormLabel>
                    <div className="width-30">
                      <Input
                        className="width-30"
                        value={secureJsonData?.apiKey || ''}
                        onChange={onUpdateDatasourceSecureJsonDataOption(this.props, 'apiKey')}
                      />
                    </div>
                  </div>
                </div>
              )}
      </div>
}


