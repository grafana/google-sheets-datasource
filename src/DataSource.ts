import { DataSourceInstanceSettings, ScopedVars, SelectableValue } from '@grafana/data';
import { DataSourceOptions } from '@grafana/google-sdk';
import { DataSourceWithBackend, getTemplateSrv } from '@grafana/runtime';
import { SheetsQuery } from './types';

export class DataSource extends DataSourceWithBackend<SheetsQuery, DataSourceOptions> {
  authType: string;
  constructor(instanceSettings: DataSourceInstanceSettings<DataSourceOptions>) {
    super(instanceSettings);
    this.authType = instanceSettings.jsonData.authenticationType;
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
}
