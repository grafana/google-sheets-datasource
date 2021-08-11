import React, { useEffect, ChangeEvent } from 'react';
import { QueryEditorProps } from '@grafana/data';
import { LinkButton, InlineFormLabel, Segment, SegmentAsync, LegacyForms, Button } from '@grafana/ui';
import { DataSource } from '../DataSource';
import { SheetsQuery, SheetsSourceOptions } from '../types';
import { usePicker } from './usePicker';

type Props = QueryEditorProps<DataSource, SheetsQuery, SheetsSourceOptions>;

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

const defaultCacheDuration = 300;

export const formatCacheTimeLabel = (s: number = defaultCacheDuration) => {
  if (s < 60) {
    return s + 's';
  } else if (s < 3600) {
    return s / 60 + 'm';
  }

  return s / 3600 + 'h';
};

export function QueryEditor({ onChange, onRunQuery, query, datasource }: Props) {
  const { openPicker } = usePicker({
    pickerCallback: (data) => {
      if (data.action === google.picker.Action.PICKED) {
        onChange({ ...query, spreadsheet: data.docs[0].id });
        onRunQuery();
      }
    },
  });

  useEffect(() => {
    if (!query.hasOwnProperty('cacheDurationSeconds')) {
      // We should probably not mutate the query like this
      query.cacheDurationSeconds = defaultCacheDuration; // um :(
    }
  }, [query]);

  const onRangeChange = (event: ChangeEvent<HTMLInputElement>) => {
    onChange({
      ...query,
      range: event.target.value,
    });
  };

  const onSpreadsheetIDChange = (item: any) => {
    if (!item.value) {
      return; // ignore delete?
    }

    const v = item.value;
    // Check for pasted full URLs
    if (/(.*)\/spreadsheets\/d\/(.*)/.test(v)) {
      onChange({ ...query, ...getGoogleSheetRangeInfoFromURL(v) });
    } else {
      onChange({ ...query, spreadsheet: v });
    }
    onRunQuery();
  };

  const toggleUseTimeFilter = (event?: React.SyntheticEvent<HTMLInputElement>) => {
    onChange({
      ...query,
      useTimeFilter: !query.useTimeFilter,
    });
    onRunQuery();
  };

  return (
    <>
      <div className="gf-form-inline">
        <Button
          type="button"
          onClick={() => {
            openPicker();
          }}
        >
          Open picker
        </Button>
      </div>
      <div className="gf-form-inline">
        <InlineFormLabel
          width={10}
          className="query-keyword"
          tooltip={
            <p>
              The <code>spreadsheetId</code> is used to identify which spreadsheet is to be accessed or altered. This ID
              is the value between the &quot/d/&quot and the &quot/edit&quot in the URL of your spreadsheet.
            </p>
          }
        >
          Spreadsheet ID
        </InlineFormLabel>
        <SegmentAsync
          loadOptions={() => datasource.getSpreadSheets()}
          placeholder="Enter SpreadsheetID"
          value={query.spreadsheet}
          allowCustomValue={true}
          onChange={onSpreadsheetIDChange}
        ></SegmentAsync>
        {query.spreadsheet && (
          <LinkButton
            style={{ marginTop: 1 }}
            variant="link"
            icon="link"
            href={toGoogleURL(query)}
            target="_blank"
          ></LinkButton>
        )}
        <div className="gf-form gf-form--grow">
          <div className="gf-form-label gf-form-label--grow" />
        </div>
      </div>
      <div className="gf-form-inline">
        <InlineFormLabel
          width={10}
          className="query-keyword"
          tooltip={
            <p>
              A string like <code>Sheet1!A1:B2</code>, that refers to a group of cells in the spreadsheet, and is
              typically used in formulas. Named ranges are also supported. When a named range conflicts with a
              sheet&rsquo;s name, the named range is preferred.
            </p>
          }
        >
          Range
        </InlineFormLabel>
        <input
          className="gf-form-input width-14"
          value={query.range || ''}
          placeholder="ie: Class Data!A2:E"
          onChange={onRangeChange}
          onBlur={onRunQuery}
        ></input>
        <div className="gf-form gf-form--grow">
          <div className="gf-form-label gf-form-label--grow" />
        </div>
      </div>
      <div className="gf-form-inline">
        <InlineFormLabel
          width={10}
          className="query-keyword"
          tooltip="Time in seconds that the spreadsheet will be cached in Grafana after receiving a response from the spreadsheet API"
        >
          Cache Time
        </InlineFormLabel>
        <Segment
          value={{ label: formatCacheTimeLabel(query.cacheDurationSeconds), value: query.cacheDurationSeconds }}
          options={[0, 5, 10, 30, 60, 60 * 2, 60 * 5, 60 * 10, 60 * 30, 3600, 3600 * 2, 3600 * 5].map((value) => ({
            label: formatCacheTimeLabel(value),
            value,
            description: value ? '' : 'Response is not cached at all',
          }))}
          onChange={({ value }) => onChange({ ...query, cacheDurationSeconds: value! })}
        />
        <div className="gf-form gf-form--grow">
          <div className="gf-form-label gf-form-label--grow" />
        </div>
      </div>
      <div className="gf-form-inline">
        <LegacyForms.Switch
          label="Use Time Filter"
          labelClass={'width-10  query-keyword'}
          tooltip="Apply the dashboard time range to the first time field"
          checked={query.useTimeFilter === true}
          onChange={toggleUseTimeFilter}
        />
        <div className="gf-form gf-form--grow">
          <div className="gf-form-label gf-form-label--grow" />
        </div>
      </div>
    </>
  );
}
