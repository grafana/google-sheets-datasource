import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './DataSource';
import { ConfigEditor } from './ConfigEditor';
import { QueryEditor } from './QueryEditor';
import { SheetsQuery, SheetsSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, SheetsQuery, SheetsSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
