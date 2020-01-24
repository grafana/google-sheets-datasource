import React, { PureComponent } from 'react';
import { DataSourcePluginOptionsEditorProps } from '@grafana/data';
import { MyDataSourceOptions } from './types';

export type Props = DataSourcePluginOptionsEditorProps<MyDataSourceOptions>;

export class ConfigEditor extends PureComponent<Props> {
  render() {
    // No options yet... may a global seed?

    return <div className="gf-form-group">&nbsp;</div>;
  }
}
