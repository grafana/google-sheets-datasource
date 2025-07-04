import { map, Observable } from 'rxjs';
import { CustomVariableSupport, DataQueryResponse, DataQueryRequest, DataFrameView } from '@grafana/data';
import { DataSource } from './DataSource';
import VariableQueryEditor from './components/VariableQueryEditor';
import type { SheetsVariableQuery } from './types';

export class SheetsVariableSupport extends CustomVariableSupport<DataSource, SheetsVariableQuery> {
  constructor(private readonly datasource: DataSource) {
    super();
    this.datasource = datasource;
    this.query = this.query.bind(this);
  }
  editor = VariableQueryEditor;
  query(request: DataQueryRequest<SheetsVariableQuery>): Observable<DataQueryResponse> {
    let query = { ...request?.targets[0], refId: 'metricFindQuery' };
    return this.datasource
      .query({ ...request, targets: [query] })
      .pipe(map((response) => ({ ...response, data: response.data || [] })))
      .pipe(map((response) => queryResponseToVariablesFrame(query, response)));
  }
}

export const queryResponseToVariablesFrame = (query: SheetsVariableQuery, response: DataQueryResponse) => {
  if (response?.data?.length < 1) {
    return { ...response, data: [] };
  }
  const view = new DataFrameView(response.data[0] || {});
  const data = view
    .map((item) => {
      const value = item[query.valueField || ''];
      const text = item[query.labelField || ''] || value;
      return { value, text };
    })
    .filter((item) => !!item.value && !!item.text);
  return { ...response, data };
};
