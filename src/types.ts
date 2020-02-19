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

//-------------------------------------------------------------------------------
// The Sheets specicif types
//-------------------------------------------------------------------------------
export interface GoogleSheetRangeInfo {
  spreadsheet: SelectableValue<string>;
  range: string;
}

export interface SheetsQuery extends DataQuery, GoogleSheetRangeInfo {
  cacheDurationSeconds: number;
  queryType: string;
}

export interface JWTFile {
  type: string;
  project_id: string;
  private_key_id: string;
  private_key: string;
  client_email: string;
  client_id: string;
  auth_uri: string;
  token_uri: string;
  auth_provider_x509_cert_url: string;
  client_x509_cert_url: string;
}

export interface SheetsSourceOptions extends GoogleCloudOptions {
  authType: GoogleAuthType;
  apiKey: string;
  jwt: JWTFile;
}

export interface GoogleSheetsSecureJsonData {
  apiKey: string;
}
