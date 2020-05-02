import { DataSourceInstanceSettings, SelectableValue } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';

import { SheetsQuery, SheetsSourceOptions } from './types';

export class DataSource extends DataSourceWithBackend<SheetsQuery, SheetsSourceOptions> {
  /** @ngInject */
  constructor(instanceSettings: DataSourceInstanceSettings<SheetsSourceOptions>, private templateSrv: any) {
    super(instanceSettings);
  }

  // Support template variables for spreadsheet and range
  applyTemplateVariables(query: SheetsQuery) {
    return {
      ...query,
      spreadsheet: this.templateSrv.replace(query.spreadsheet),
      range: this.templateSrv.replace(query.range),
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
