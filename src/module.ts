import { DataSourcePlugin } from '@grafana/data';
import { initPluginTranslations } from '@grafana/i18n';
import pluginJson from 'plugin.json';
import { config } from '@grafana/runtime';
import semver from 'semver';
import { loadResources } from './loadResources';
import { DataSourceOptions } from '@grafana/google-sdk';
import { ConfigEditor, MetaInspector, QueryEditor } from './components';
import { DataSource } from './DataSource';
import { SheetsQuery } from './types';

// Before Grafana version 12.1.0 the plugin is responsible for loading translation resources
// In Grafana version 12.1.0 and later Grafana is responsible for loading translation resources
const loaders = semver.lt(config?.buildInfo?.version, '12.1.0') ? [loadResources] : [];

await initPluginTranslations(pluginJson.id, loaders);

export const plugin = new DataSourcePlugin<DataSource, SheetsQuery, DataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor)
  .setMetadataInspector(MetaInspector);
