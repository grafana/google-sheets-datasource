import { DataSourcePlugin } from '@grafana/data';
import { DataSourceOptions } from '@grafana/google-sdk';
import { ConfigEditor, MetaInspector, QueryEditor } from './components';
import { DataSource } from './DataSource';
import { SheetsQuery } from './types';

export const plugin = new DataSourcePlugin<DataSource, SheetsQuery, DataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor)
  .setMetadataInspector(MetaInspector);
