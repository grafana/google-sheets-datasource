import { DataSourceInstanceSettings } from '@grafana/data';
import { getBackendSrv, DataSourceWithBackend } from '@grafana/runtime'; //DataSourceWithBackend

import { SheetsQuery, SheetsSourceOptions } from './types';

export class DataSource extends DataSourceWithBackend<SheetsQuery, SheetsSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<SheetsSourceOptions>) {
    super(instanceSettings);
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
            intervalMs: 1,
            maxDataPoints: 1,
            queryType: 'testAPI',
          },
        ],
      })
      .then((rsp: any) => {
        return {
          status: 'success',
          message: 'Success',
        };
      });
  }
}
