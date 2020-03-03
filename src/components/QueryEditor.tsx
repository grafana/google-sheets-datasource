import React, { PureComponent, ChangeEvent } from 'react';
import { QueryEditorProps } from '@grafana/data';
import { LinkButton, FormLabel, Segment, SegmentAsync } from '@grafana/ui';
import { DataSource } from '../DataSource';
import { SheetsQuery, SheetsSourceOptions } from '../types';

type Props = QueryEditorProps<DataSource, SheetsQuery, SheetsSourceOptions>;

interface State {}

export function getGoogleSheetRangeInfoFromURL(url: string): Partial<SheetsQuery> {
  let idx = url?.indexOf('/d/');
  if (!idx) {
    // The original value
    return { spreadsheet: url };
  }

  let id = url.substring(idx + 3);
  idx = id.indexOf('/');
  if (idx) {
    id = id.substring(0, idx);
  }

  idx = url.indexOf('range=');
  if (idx > 0) {
    const sub = url.substring(idx + 'range='.length);
    return { spreadsheet: id, range: sub };
  }
  return { spreadsheet: id };
}

export function toGoogleURL(info: SheetsQuery): string {
  let url = `https://docs.google.com/spreadsheets/d/${info.spreadsheet}/view`;
  if (info.range) {
    url += '#range=' + info.range;
  }
  return url;
}

export class QueryEditor extends PureComponent<Props, State> {
  componentWillMount() {
    if (!this.props.query.hasOwnProperty('cacheDurationSeconds')) {
      this.props.query.cacheDurationSeconds = 300;
    }
  }

  onRangeChange = (event: ChangeEvent<HTMLInputElement>) => {
    this.props.onChange({
      ...this.props.query,
      range: event.target.value,
    });
  };

  onSpreadsheetIDChange = (item: any) => {
    const { query, onRunQuery, onChange } = this.props;

    /(.*)\/spreadsheets\/d\/(.*)/.test(item.value!)
      ? onChange({ ...query, ...getGoogleSheetRangeInfoFromURL(item.value!) })
      : onChange({ ...query, spreadsheet: item });

    onRunQuery();
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
            loadOptions={() => datasource.getSpreadSheets()}
            placeholder="Enter SpreadsheetID"
            value={query.spreadsheet}
            allowCustomValue={true}
            onChange={this.onSpreadsheetIDChange}
          ></SegmentAsync>
          <LinkButton
            style={{ marginTop: 1 }}
            disabled={!query.spreadsheet}
            variant="secondary"
            icon="fa fa-link"
            href={toGoogleURL(query)}
            target="_blank"
          ></LinkButton>
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
