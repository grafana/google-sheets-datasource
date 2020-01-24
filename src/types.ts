import { DataQuery, DataSourceJsonData, QueryEditorProps, SelectableValue } from '@grafana/data';
import { ComponentType } from 'react';

export enum Scenario {
  wave = 'wave',
  noise = 'noise',
  arrowFile = 'arrowFile',
  csvWave = 'csvWave',
}

export interface MyQuery extends DataQuery {
  scenario: Scenario;
}

// NOTE: the actual implementations are defined in

export type ScenarioEditorProps<T extends MyQuery> = QueryEditorProps<any, T, MyDataSourceOptions>;

export interface ScenarioProvider<T extends MyQuery> extends SelectableValue<Scenario> {
  defaultQuery: Partial<T>;
  editor: ComponentType<ScenarioEditorProps<T>>;
}

/**
 * These are options configured for each DataSource instance
 */
export interface MyDataSourceOptions extends DataSourceJsonData {}
