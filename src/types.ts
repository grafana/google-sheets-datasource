import { DataQuery, DataSourceJsonData, SelectableValue } from '@grafana/data';

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
  GCE = 'gce',
  NONE = 'none',
}

export const googleAuthTypes = [
  { label: 'None (public)', value: GoogleAuthType.NONE },
  { label: 'Google JWT File', value: GoogleAuthType.JWT },
  { label: 'GCE Default Service Account', value: GoogleAuthType.GCE },
];

export interface GoogleCloudOptions extends DataSourceJsonData {
  authenticationType: GoogleAuthType;
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
  spreadsheet: SelectableValue<string>;
  range: string;
  cacheDurationSeconds: number;
}

export interface SheetsSourceOptions extends GoogleCloudOptions {
  authType: GoogleAuthType;
}

export interface GoogleSheetsSecureJsonData {
  apiKey: string;
  jwt: string;
}
