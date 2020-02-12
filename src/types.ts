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

export enum ResultFormatType {
  TABLE = 'TABLE',
  TIME_SERIES = 'TIME_SERIES',
}

export const resultFormats = [
  { label: 'Table', value: ResultFormatType.TABLE },
  { label: 'Time Series', value: ResultFormatType.TIME_SERIES },
];

export enum MajorDimensionType {
  DIMENSION_UNSPECIFIED = 'DIMENSION_UNSPECIFIED',
  ROWS = 'ROWS',
  COLUMNS = 'COLUMNS',
}

export const majorDimensions = [
  { label: 'Rows', value: MajorDimensionType.ROWS },
  { label: 'Columns', value: MajorDimensionType.COLUMNS },
];

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
  spreadsheetId: string;
  range: string;
  majorDimension: string;
}

export interface SheetsQuery extends DataQuery, GoogleSheetRangeInfo {
  queryType: string;
  resultFormat: ResultFormatType;
  timeColumn: SelectableValue<number>;
  metricColumns: Array<SelectableValue<number>>;
}

export interface SheetsSourceOptions extends GoogleCloudOptions {
  authType: GoogleAuthType;
  apiKey: string;
  jwtFile: string;
}

export interface GoogleSheetsSecureJsonData {
  apiKey: string;
}
