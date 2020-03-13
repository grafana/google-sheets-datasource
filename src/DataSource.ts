import { DataSourceInstanceSettings, SelectableValue } from '@grafana/data';
import { DataSourceWithBackend, getBackendSrv } from '@grafana/runtime';

import { SheetsQuery, SheetsSourceOptions } from './types';


export class DataSource extends DataSourceWithBackend<SheetsQuery, SheetsSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<SheetsSourceOptions>) {
    super(instanceSettings);
  }

  async getSpreadSheets(): Promise<Array<SelectableValue<string>>> {
    return this.getResource('spreadsheets').then(({ spreadsheets }) =>
      spreadsheets ? Object.entries(spreadsheets).map(([value, label]) => ({ label, value } as SelectableValue<string>)) : []
    );
  }

  /**
   * Checks the plugin health
   */
  async testDatasource(): Promise<any> {
    return getBackendSrv().get(`/api/datasources/${this.id}/health`).then( res => {
      console.log( 'TEST', res );
      return {
        status: 'success',
        message: 'Success',
      };
    });
  }

  // async testDatasource() {
  //   return this.getResource('test').then((rsp: any) => {
  //     if (rsp.error) {
  //       return {
  //         status: 'fail',
  //         message: rsp.error,
  //       };
  //     }

  //     return {
  //       status: 'success',
  //       message: 'Success',
  //     };
  //   });
  // }
}
