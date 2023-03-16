import { DataSourceInstanceSettings, SelectableValue, ScopedVars, CoreApp } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';

import { SheetsQuery, SheetsSourceOptions } from './types';

export class DataSource extends DataSourceWithBackend<SheetsQuery, SheetsSourceOptions> {
  constructor(private instanceSettings: DataSourceInstanceSettings<SheetsSourceOptions>) {
    super(instanceSettings);
  }

  // Enables default annotation support for 7.2+
  annotations = {};

  // Support template variables for spreadsheet and range
  applyTemplateVariables(query: SheetsQuery, scopedVars: ScopedVars) {
    const templateSrv = getTemplateSrv();
    return {
      ...query,
      spreadsheet: templateSrv.replace(query.spreadsheet, scopedVars),
      range: query.range ? templateSrv.replace(query.range, scopedVars) : '',
    };
  }

  async getSpreadSheets(): Promise<Array<SelectableValue<string>>> {
    return this.getResource('spreadsheets').then(({ spreadsheets }) =>
      spreadsheets
        ? Object.entries(spreadsheets).map(([value, label]) => ({ label, value } as SelectableValue<string>))
        : []
    );
  }

  getDefaultQuery(app: CoreApp): Partial<SheetsQuery> {
    return { spreadsheet: this.instanceSettings.jsonData.defaultSheetID };
  }
}
