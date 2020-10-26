import { DataSourceInstanceSettings, SelectableValue } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';

import { SheetsQuery, SheetsSourceOptions } from './types';

export class DataSource extends DataSourceWithBackend<SheetsQuery, SheetsSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<SheetsSourceOptions>) {
    super(instanceSettings);
  }

  // Enables default annotation support for 7.2+
  annotations = {};

  // Support template variables for spreadsheet and range
  applyTemplateVariables(query: SheetsQuery) {
    const templateSrv = getTemplateSrv();
    return {
      ...query,
      spreadsheet: templateSrv.replace(query.spreadsheet),
      range: query.range ? templateSrv.replace(query.range) : '',
    };
  }

  async getSpreadSheets(): Promise<Array<SelectableValue<string>>> {
    return this.getResource('spreadsheets').then(({ spreadsheets }) =>
      spreadsheets
        ? Object.entries(spreadsheets).map(([value, label]) => ({ label, value } as SelectableValue<string>))
        : []
    );
  }
}
