import { DataSourceInstanceSettings, SelectableValue } from '@grafana/data';
import { getBackendSrv, DataSourceWithBackend } from '@grafana/runtime'; //DataSourceWithBackend

import { SheetsQuery, SheetsSourceOptions } from './types';

export class DataSource extends DataSourceWithBackend<SheetsQuery, SheetsSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<SheetsSourceOptions>) {
    super(instanceSettings);
  }

  async getSpreadSheets(): Promise<SelectableValue<string>[]> {
    return this.getResource('spreadsheets').then(({ spreadsheets }) =>
      Object.entries(spreadsheets).map(([value, label]) => ({ label, value } as SelectableValue<string>))
    );
  }

  async getResource(path: string, body?: any): Promise<{ [key: string]: any }> {
    return getBackendSrv().post(`/api/datasources/${this.id}/resources/${path}`, { ...body });
  }

  async testDatasource() {
    return getBackendSrv()
      .post('/api/ds/query', {
        from: '5m',
        to: 'now',
        queries: [
          {
            refId: 'A',
            datasource: this.name,
            datasourceId: this.id,
            queryType: 'testAPI',
          },
        ],
      })
      .then((rsp: any) => {
        if (rsp.results[''].meta && rsp.results[''].meta.error) {
          return {
            status: 'fail',
            message: rsp.results[''].meta.error,
          };
        }

        return {
          status: 'success',
          message: 'Success',
        };
      });
  }
}
