import { CustomVariableSupport, DataQueryRequest } from '@grafana/data';
import { GoogleSheetsDataSource } from './datasource';
import { QueryEditor } from './components/QueryEditor';
import { SheetsQuery } from './types';

export class GoogleSheetsVariableSupport extends CustomVariableSupport<GoogleSheetsDataSource, SheetsQuery> {
  editor = QueryEditor;

  constructor(private datasource: GoogleSheetsDataSource) {
    super();
  }

  getDefaultQuery(): Partial<SheetsQuery> {
    return {
      refId: 'tempvar',
    };
  }

  query(request: DataQueryRequest<SheetsQuery>) {
    if (!this.datasource) {
      throw new Error('Datasource not initialized');
    }

    return this.datasource.query(request);
  }
}
