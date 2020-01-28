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
  { label: 'Google JWT File', value: GoogleAuthType.JWT },
  { label: 'GCE Default Service Account', value: GoogleAuthType.GCE },
  { label: 'None (public)', value: GoogleAuthType.NONE },
];

export interface GoogleCloundOptions extends DataSourceJsonData {
  authenticationType: GoogleAuthType;
}

//-------------------------------------------------------------------------------
// The Sheets specicif types
//-------------------------------------------------------------------------------
export interface GoogleSheetRangeInfo {
  spreadsheetId: string;
  range: string;
}

export interface SheetsQuery extends DataQuery, GoogleSheetRangeInfo {}

export interface SheetsSourceOptions extends GoogleCloundOptions {}

export interface GoogleSheetsSecureJsonData {
  apiKey: string;
}
