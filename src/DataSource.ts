import { DataSourceInstanceSettings } from '@grafana/data';
import { getBackendSrv, DataSourceWithBackend } from '@grafana/runtime'; //DataSourceWithBackend

import { SheetsQuery, SheetsSourceOptions } from './types';

export class DataSource extends DataSourceWithBackend<SheetsQuery, SheetsSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<SheetsSourceOptions>) {
    super(instanceSettings);
  }

  metricFindQuery(query: SheetsQuery, queryType: string): Promise<any[]> {
    // TODO: convert to resource endpoit
    return getBackendSrv()
      .post('/api/ds/query', {
        from: '5m',
        to: 'now',
        queries: [
          {
            ...query,
            datasource: this.name,
            datasourceId: this.id,
            queryType,
          },
        ],
      })
      .then((rsp: any) => {
        console.log({ rsp });
        return Object.entries(rsp.results[''].meta).map(([value, label]) => ({ label, value: /^\d+$/.test(value) ? Number(value) : value }));
      });
  }

  async getResource(body?: any): Promise<any[]> {
    return getBackendSrv()
      .post(`/api/datasources/${this.id}/resources`, {
        queryType: 'getHeaders',
        ...body,
      })
      .then((rsp: any) => {
        console.log({ rsp });
        return rsp;
      });
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
