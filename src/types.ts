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

//-------------------------------------------------------------------------------
// Sheet metadata (returned in custom)
//-------------------------------------------------------------------------------

export interface CacheInfo {
  hit: boolean;
  count: number;
  time: number;
}

export interface SheetResponseMeta {
  spreadsheetId: string;
  range: string;
  cache: CacheInfo;
  warnings: string[];
}

//-------------------------------------------------------------------------------
// The Sheets specific types
//-------------------------------------------------------------------------------

export interface SheetsQuery extends DataQuery {
  spreadsheetId: string;
  range?: string; // without a range it should get all sheets
}

export interface SheetsSourceOptions extends GoogleCloudOptions {
  authType: GoogleAuthType;
  apiKey: string;
  jwtFile: string;
  cacheDurationSeconds: number;
}

export interface GoogleSheetsSecureJsonData {
  apiKey: string;
}
