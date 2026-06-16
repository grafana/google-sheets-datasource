import { map, Observable } from 'rxjs';
import { CustomVariableSupport, DataQueryResponse, DataQueryRequest, Field } from '@grafana/data';
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
    let interpolatedQuery = this.datasource.interpolateVariableQuery(query, request.scopedVars);
    return this.datasource
      .query({ ...request, targets: [interpolatedQuery] })
      .pipe(map((response) => ({ ...response, data: response.data || [] })))
      .pipe(map((response) => queryResponseToVariablesFrame(interpolatedQuery, response)));
  }
}

export const queryResponseToVariablesFrame = (query: SheetsVariableQuery, response: DataQueryResponse) => {
  // If no data or no fields, return empty data
  if (response?.data?.length < 1 || !response.data[0]?.fields?.length) {
    return { ...response, data: [] };
  }
  
  // Find the required value field. If not found, return empty data
  const frame = response.data[0];
  const valueField = frame.fields.find((field: Field) => field.name === query.valueField);
  if (!valueField) {
    return { ...response, data: [] };
  }

  // Find the optional label and filter fields.
  const labelField = query.labelField ? frame.fields.find((field: Field) => field.name === query.labelField) : null;
  const filterField = query.filterField ? frame.fields.find((field: Field) => field.name === query.filterField) : null;

  const data = [];

  // Process each row
  for (let i = 0; i < valueField.values.length; i++) {
    const value = valueField.values.get(i);
    const text = labelField ? labelField.values.get(i) : value;

    // Skip rows without value or text
    if (!value || !text) {
      continue;
    }

    // Apply filter if specified
    if (query.filterField && query.filterValue && filterField) {
      const filterValue = filterField.values.get(i);
      if (filterValue !== query.filterValue) {
        continue;
      }
    }

    data.push({ value, text });
  }

  return { ...response, data };
};
