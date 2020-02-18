import React, { PureComponent, ChangeEvent } from 'react';
import { QueryEditorProps } from '@grafana/data';
import { LinkButton, FormLabel, Segment, SegmentAsync } from '@grafana/ui';
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
  componentWillMount() {
    if (!this.props.query.queryType) {
      this.props.query.queryType = 'query';
    }

    if (!this.props.query.hasOwnProperty('cacheDurationSeconds')) {
      this.props.query.cacheDurationSeconds = 300;
    }
  }

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

  onRangeChange = (event: ChangeEvent<HTMLInputElement>) => {
    this.props.onChange({
      ...this.props.query,
      range: event.target.value,
    });
  };

  render() {
    const { query, onRunQuery, onChange, datasource } = this.props;
    return (
      <>
        <div className="gf-form-inline">
          <FormLabel
            width={10}
            className="query-keyword"
            tooltip={
              <p>
                The <code>spreadsheetId</code> is used to identify which spreadsheet is to be accessed or altered. This ID is the value between the
                "/d/" and the "/edit" in the URL of your spreadsheet.
              </p>
            }
          >
            Spreadsheet ID
          </FormLabel>
          <SegmentAsync
            loadOptions={() => datasource.metricFindQuery(query, 'getSpreadsheets')}
            placeholder="Enter SpreadsheetID"
            value={query.spreadsheetId || ''}
            allowCustomValue={true}
            onChange={({ value }) => {
              /(.*)\/spreadsheets\/d\/(.*)/.test(value!)
                ? onChange({ ...query, ...getGoogleSheetRangeInfoFromURL(value!) })
                : onChange({ ...query, spreadsheetId: value! });
              onRunQuery();
            }}
          ></SegmentAsync>
          <LinkButton disabled={!query.spreadsheetId} variant="secondary" icon="fa fa-link" href={toGoogleURL(query)} target="_blank"></LinkButton>
          <div className="gf-form gf-form--grow">
            <div className="gf-form-label gf-form-label--grow" />
          </div>
        </div>
        <div className="gf-form-inline">
          <FormLabel
            width={10}
            className="query-keyword"
            tooltip={
              <p>
                A string like <code>Sheet1!A1:B2</code>, that refers to a group of cells in the spreadsheet, and is typically used in formulas.Named
                ranges are also supported. When a named range conflicts with a sheet's name, the named range is preferred.
              </p>
            }
          >
            Range
          </FormLabel>
          <input
            className="gf-form-input width-14"
            value={query.range || ''}
            placeholder="ie: Class Data!A2:E"
            onChange={this.onRangeChange}
            onBlur={onRunQuery}
          ></input>
          <div className="gf-form gf-form--grow">
            <div className="gf-form-label gf-form-label--grow" />
          </div>
        </div>
        <div className="gf-form-inline">
          <FormLabel
            width={10}
            className="query-keyword"
            tooltip="Time in seconds that the spreadsheet will be cached in Grafana after receiving a response from the spreadsheet API"
          >
            Cache Time
          </FormLabel>
          <Segment
            value={{ label: `${query.cacheDurationSeconds}s`, value: query.cacheDurationSeconds }}
            options={[0, 5, 10, 30, 60, 120, 300, 600, 3600].map(value => ({
              label: `${value}s`,
              value,
              description: value ? '' : 'Response is not cached at all',
            }))}
            onChange={({ value }) => onChange({ ...query, cacheDurationSeconds: value! })}
          />
          <div className="gf-form gf-form--grow">
            <div className="gf-form-label gf-form-label--grow" />
          </div>
        </div>
      </>
    );
  }
}
