import {
  DataQueryRequest,
  DataQueryResponse,
  DataSourceApi,
  DataSourceInstanceSettings,
  SelectableValue,
} from '@grafana/data';
import { DataSourceOptions } from '@grafana/google-sdk';
import { SheetsQuery } from './types';
import { getBackendSrv, BackendSrvRequest, BackendDataSourceResponse, toDataQueryResponse } from '@grafana/runtime';

import { Observable, of, lastValueFrom } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { DataQuery } from '@grafana/schema';


export class DataSource2 extends DataSourceApi<SheetsQuery, DataSourceOptions> {
  baseUrl: string;

  constructor(instanceSettings: DataSourceInstanceSettings<DataSourceOptions>) {
    super(instanceSettings);

    // K8s!
    this.baseUrl = 'https://localhost:6443/apis/googlesheets.ext.grafana.com/v1/namespaces/default/datasources/12345/';
  }

  query(request: DataQueryRequest<SheetsQuery>): Observable<DataQueryResponse> {
    const { range, requestId, hideFromInspector = false } = request;
    const queries = request.targets; // TODO, cleanup
    const body: any = { queries };
    if (range) {
      body.from = range.from.valueOf().toString();
      body.to = range.to.valueOf().toString();
    }

    const headers: Record<string, string> = {};
    if (request.dashboardUID) {
      headers[PluginRequestHeaders.DashboardUID] = request.dashboardUID;
    }
    if (request.panelId) {
      headers[PluginRequestHeaders.PanelID] = `${request.panelId}`;
    }
    // if (request.queryGroupId) {
    //   headers[PluginRequestHeaders.QueryGroupID] = `${request.queryGroupId}`;
    // }
    return getBackendSrv()
      .fetch<BackendDataSourceResponse>({
        url: this.baseUrl + 'query',
        method: 'POST',
        data: body,
        requestId,
        hideFromInspector,
        headers,
      })
      .pipe(
        map((raw) => {
          return toDataQueryResponse(raw, queries as DataQuery[]);
        }),
        catchError((err) => {
          return of(toDataQueryResponse(err));
        })
      );
  }

  async getSpreadSheets(): Promise<Array<SelectableValue<string>>> {
    return this.getResource('spreadsheets').then(({ spreadsheets }) =>
      spreadsheets
        ? Object.entries(spreadsheets).map(([value, label]) => ({ label, value } as SelectableValue<string>))
        : []
    );
  }

  /**
   * Make a GET request to the datasource resource path
   */
  async getResource<T = any>(
    path: string,
    params?: BackendSrvRequest['params'],
    options?: Partial<BackendSrvRequest>
  ): Promise<T> {
    const headers: string[] = []; //this.getRequestHeaders();
    const result = await lastValueFrom(
      getBackendSrv().fetch<T>({
        ...options,
        method: 'GET',
        headers: options?.headers ? { ...options.headers, ...headers } : headers,
        params: params ?? options?.params,
        url: `${this.baseUrl}/resources/${path}`,
      })
    );
    return result.data;
  }

  /**
   * Checks whether we can connect to the API.
   */
  async testDatasource() {
    // TODO... real healthcheck call
    return Promise.resolve({
      status: 'success',
      message: 'Success',
    });
  }
}



// Internal for now
enum PluginRequestHeaders {
    PluginID = 'X-Plugin-Id', // can be used for routing
    DatasourceUID = 'X-Datasource-Uid', // can be used for routing/ load balancing
    DashboardUID = 'X-Dashboard-Uid', // mainly useful for debugging slow queries
    PanelID = 'X-Panel-Id', // mainly useful for debugging slow queries
    QueryGroupID = 'X-Query-Group-Id', // mainly useful to find related queries with query splitting
    FromExpression = 'X-Grafana-From-Expr', // used by datasources to identify expression queries
  }
