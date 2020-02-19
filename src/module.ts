import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './DataSource';
import { ConfigEditor, QueryEditor } from './components';
import { SheetsQuery, SheetsSourceOptions } from './types';
import { MetaInspector } from './MetaInspector';

export const plugin = new DataSourcePlugin<DataSource, SheetsQuery, SheetsSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor)
  .setMetadataInspector(MetaInspector);
