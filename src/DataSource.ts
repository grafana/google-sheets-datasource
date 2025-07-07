import {
  DataQueryRequest,
  DataQueryResponse,
  DataSourceInstanceSettings,
  ScopedVars,
  SelectableValue,
} from '@grafana/data';
import { DataSourceOptions } from '@grafana/google-sdk';
import { DataSourceWithBackend, getTemplateSrv, TemplateSrv } from '@grafana/runtime';
import { SheetsQuery, SheetsVariableQuery } from './types';
import { Observable } from 'rxjs';
import { trackRequest } from 'tracking';
import { SheetsVariableSupport } from 'variables';

export class DataSource extends DataSourceWithBackend<SheetsQuery, DataSourceOptions> {
  authType: string;
  constructor(
    instanceSettings: DataSourceInstanceSettings<DataSourceOptions>,
    private readonly templateSrv: TemplateSrv = getTemplateSrv()
  ) {
    super(instanceSettings);
    this.authType = instanceSettings.jsonData.authenticationType;
    this.variables = new SheetsVariableSupport(this);
  }

  query(request: DataQueryRequest<SheetsQuery>): Observable<DataQueryResponse> {
    trackRequest(request);
    return super.query(request);
  }

  // Enables default annotation support for 7.2+
  annotations = {};

  // Support template variables for spreadsheet and range
  applyTemplateVariables(query: SheetsQuery, scopedVars: ScopedVars) {
    return {
      ...query,
      spreadsheet: this.interpolateVariable(query.spreadsheet, scopedVars) ?? '',
      range: this.interpolateVariable(query.range, scopedVars),
    };
  }

  interpolateVariableQuery(query: SheetsVariableQuery, scopedVars: ScopedVars) {
    return {
      // Interpolate query
      ...this.applyTemplateVariables(query, scopedVars),
      // Interpolate additional fields in variable query
      filterValue: this.interpolateVariable(query.filterValue, scopedVars),
      filterField: this.interpolateVariable(query.filterField, scopedVars),
      valueField: this.interpolateVariable(query.valueField, scopedVars),
      labelField: this.interpolateVariable(query.labelField, scopedVars),
    };
  }

  interpolateVariable(value: string | undefined, item: ScopedVars) {
    // If we don't have value or value is empty string, return it
    if (!value) {
      return value;
    }
    return this.templateSrv.replace(value, item);
  }

  async getSpreadSheets(): Promise<Array<SelectableValue<string>>> {
    return this.getResource('spreadsheets').then(({ spreadsheets }) =>
      spreadsheets
        ? Object.entries(spreadsheets).map(([value, label]) => ({ label, value }) as SelectableValue<string>)
        : []
    );
  }
}
