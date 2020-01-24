import React, { PureComponent } from 'react';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { SheetsSourceOptions } from './types';

export type Props = DataSourcePluginOptionsEditorProps<SheetsSourceOptions>;

export class ConfigEditor extends PureComponent<Props> {
  render() {
    return <div className="gf-form-group">TODO... enter JWT etc</div>;
  }
}
