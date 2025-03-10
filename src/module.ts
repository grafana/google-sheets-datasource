import { DataSourcePlugin } from '@grafana/data';
import { DataSourceOptions } from '@grafana/google-sdk';
import { ConfigEditor, MetaInspector, QueryEditor } from './components';
import { GoogleSheetsDataSource } from './datasource';
import { SheetsQuery } from './types';

export const plugin = new DataSourcePlugin<GoogleSheetsDataSource, SheetsQuery, DataSourceOptions>(
  GoogleSheetsDataSource
)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor)
  .setMetadataInspector(MetaInspector);
