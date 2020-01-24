import { DataSourceInstanceSettings } from '@grafana/data';
import { getBackendSrv, DataSourceWithBackend } from '@grafana/runtime';

import { MyQuery, MyDataSourceOptions } from './types';

export class DataSource extends DataSourceWithBackend<MyQuery, MyDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<MyDataSourceOptions>) {
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
