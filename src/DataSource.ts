import { DataSourceInstanceSettings, SelectableValue } from '@grafana/data';
import { DataSourceWithBackend, getDataSourceSrv } from '@grafana/runtime';

import { SheetsQuery, SheetsSourceOptions } from './types';

export class DataSource extends DataSourceWithBackend<SheetsQuery, SheetsSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<SheetsSourceOptions>) {
    super(instanceSettings);
  }

  // Support template variables for spreadsheet and range
  applyTemplateVariables(query: SheetsQuery) {
    let templateSrv = (getDataSourceSrv() as any).templateSrv;
    return {
      ...query,
      spreadsheet: templateSrv.replace(query.spreadsheet),
      range: templateSrv.replace(query.range),
    };
  }

  async getSpreadSheets(): Promise<Array<SelectableValue<string>>> {
    return this.getResource('spreadsheets').then(({ spreadsheets }) =>
      spreadsheets ? Object.entries(spreadsheets).map(([value, label]) => ({ label, value } as SelectableValue<string>)) : []
    );
  }
}
