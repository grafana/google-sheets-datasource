import { DataSourceInstanceSettings } from '@grafana/data';
import { getBackendSrv, DataSourceWithBackend } from '@grafana/runtime';

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
