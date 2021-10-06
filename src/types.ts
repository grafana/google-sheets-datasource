import { DataQuery, DataSourceJsonData } from '@grafana/data';

//-------------------------------------------------------------------------------
// General google cloud auth types
// same as stackdriver etc
//-------------------------------------------------------------------------------

export interface JWT {
  private_key: any;
  token_uri: any;
  client_email: any;
  project_id: any;
}

export enum GoogleAuthType {
  JWT = 'jwt',
  KEY = 'key',
}

export const googleAuthTypes = [
  { label: 'API Key (public spreadsheets)', value: GoogleAuthType.KEY },
  { label: 'Google JWT File (public and private spreadsheets)', value: GoogleAuthType.JWT },
];

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

export interface SheetsSourceOptions extends DataSourceJsonData {
  authType: GoogleAuthType;
}

export interface GoogleSheetsSecureJsonData {
  apiKey?: string;
  jwt?: string;
}
