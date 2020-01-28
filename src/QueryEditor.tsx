import React, { PureComponent, ChangeEvent } from 'react';
import { QueryEditorProps } from '@grafana/data';
import { FormField } from '@grafana/ui';
import { DataSource } from './DataSource';
import { SheetsQuery, SheetsSourceOptions, GoogleSheetRangeInfo } from './types';

type Props = QueryEditorProps<DataSource, SheetsQuery, SheetsSourceOptions>;

interface State {}

export function getGoogleSheetRangeInfoFromURL(url: string): Partial<GoogleSheetRangeInfo> {
  let idx = url?.indexOf('/d/');
  if (!idx) {
    // The original value
    return { spreadsheetId: url };
  }

  let id = url.substring(idx + 3);
  idx = id.indexOf('/');
  if (idx) {
    id = id.substring(0, idx);
  }

  idx = url.indexOf('range=');
  if (idx > 0) {
    const sub = url.substring(idx + 'range='.length);
    return { spreadsheetId: id, range: sub };
  }
  return { spreadsheetId: id };
}

export function toGoogleURL(info: GoogleSheetRangeInfo): string {
  let url = `https://docs.google.com/spreadsheets/d/${info.spreadsheetId}/view`;
  if (info.range) {
    url += '#range=' + info.range;
  }
  return url;
}

const PASTE_SEPERATOR = '»';

export class QueryEditor extends PureComponent<Props, State> {
  onComponentDidMount() {}

  onSpreadsheetIdPasted = (e: any) => {
    const v = e.clipboardData.getData('text/plain');
    if (v) {
      const info = getGoogleSheetRangeInfoFromURL(v);
      if (info.spreadsheetId) {
        console.log('PASTED', v, info);
        info.spreadsheetId = info.spreadsheetId + PASTE_SEPERATOR;
        this.props.onChange({
          ...this.props.query,
          ...info,
        });
        console.log('UPDATED', info);
      }
    }
  };

  onSpreadsheetIdChange = (event: ChangeEvent<HTMLInputElement>) => {
    console.log('CHANGE', event.target.value);
    let v = event.target.value;
    const idx = v.indexOf(PASTE_SEPERATOR);
    if (idx > 0) {
      v = v.substring(0, idx);
    }
    this.props.onChange({
      ...this.props.query,
      spreadsheetId: v,
    });
  };

  onRangeChange = (event: ChangeEvent<HTMLInputElement>) => {
    this.props.onChange({
      ...this.props.query,
      range: event.target.value,
    });
  };

  render() {
    const { query, onRunQuery } = this.props;
    return (
      <>
        <div className="gf-form">
          <div className="form-field">
            <FormField
              inputWidth={30}
              labelWidth={8}
              label="Spreadsheet ID"
              value={query.spreadsheetId || ''}
              placeholder="Enter ID from URL"
              onPaste={this.onSpreadsheetIdPasted}
              onChange={this.onSpreadsheetIdChange}
              onBlur={onRunQuery}
            />
            <a href={toGoogleURL(query)}>link</a>
          </div>
        </div>
        <div className="gf-form">
          <div className="form-field">
            <FormField
              inputWidth={30}
              labelWidth={8}
              label="Range"
              value={query.range || ''}
              placeholder="ie: Class Data!A2:E"
              onChange={this.onRangeChange}
              onBlur={onRunQuery}
            />
          </div>
        </div>
      </>
    );
  }
}
