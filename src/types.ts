import { DataQuery } from '@grafana/schema';
import { GoogleAuthType, GOOGLE_AUTH_TYPE_OPTIONS, DataSourceSecureJsonData } from '@grafana/google-sdk';

export const GoogleSheetsAuth = {
  ...GoogleAuthType,
  API: 'key',
} as const;

export const googleSheetsAuthTypes = [{ label: 'API Key', value: GoogleSheetsAuth.API }, ...GOOGLE_AUTH_TYPE_OPTIONS];

export interface GoogleSheetsSecureJSONData extends DataSourceSecureJsonData {
  apiKey?: string;
}

export interface CacheInfo {
  hit: boolean;
  count: number;
  expires: string;
}

export interface SheetResponseMeta {
  spreadsheetId: string;
  range: string;
  majorDimension: string;
  cache: CacheInfo;
  warnings: string[];
}

//-------------------------------------------------------------------------------
// The Sheets specific types
//-------------------------------------------------------------------------------

export interface SheetsQuery extends DataQuery {
  spreadsheet: string;
  range?: string;
  cacheDurationSeconds?: number;
  useTimeFilter?: boolean;
}

export interface SheetsVariableQuery extends SheetsQuery {
  valueField?: string;
  labelField?: string;
  filterField?: string;
  filterValue?: string;
}
