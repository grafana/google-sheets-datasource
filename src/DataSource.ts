import { DataSourceInstanceSettings, SelectableValue, ScopedVars } from '@grafana/data';
import { DataSourceWithBackend, getTemplateSrv, HealthStatus } from '@grafana/runtime';

import { GoogleAuthType, SheetsQuery, SheetsSourceOptions } from './types';

export class DataSource extends DataSourceWithBackend<SheetsQuery, SheetsSourceOptions> {
  constructor(public instanceSettings: DataSourceInstanceSettings<SheetsSourceOptions>) {
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

  async callHealthCheck() {
    if (this.instanceSettings.jsonData.authType === GoogleAuthType.OAUTH && window.gapi) {
      if (gapi.auth2.getAuthInstance().isSignedIn.get()) {
        return Promise.resolve({ status: HealthStatus.OK, message: 'Data source is working.' });
      }

      return Promise.resolve({ status: HealthStatus.Error, message: 'You need to sign in here.' });
    }
    return super.callHealthCheck();
  }
}
