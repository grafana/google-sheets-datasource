import React, { useEffect, useState } from 'react';
import { InlineFieldRow, InlineFormLabel, Select, useTheme2, Input, Label } from '@grafana/ui';
import { QueryEditor } from './QueryEditor';
import { DataSource } from '../DataSource';
import { SheetsVariableQuery } from '../types';
import { CoreApp, Field, getDefaultTimeRange, GrafanaTheme2, SelectableValue } from '@grafana/data';
import { lastValueFrom } from 'rxjs';
import { css } from '@emotion/css';

interface Props {
  query: SheetsVariableQuery;
  onChange: (query: SheetsVariableQuery) => void;
  onRunQuery: () => void;
  datasource: DataSource;
}

const VariableQueryEditor = (props: Props) => {
  const [choices, setChoices] = useState<string[]>([]);
  const [loading, setLoading] = useState(false);

  const { query, datasource, onChange, onRunQuery } = props;

  const theme = useTheme2();
  const styles = getStyles(theme);

  useEffect(() => {
    const fetchChoices = async (query: SheetsVariableQuery) => {
      try {
        setLoading(true);
        const res = await lastValueFrom(datasource.query(createRequest(query)));
        const columns = (res?.data[0]?.fields ?? []).map((f: Field) => f.name);
        setChoices(columns);
      } catch (err) {
        setChoices([]);
      } finally {
        setLoading(false);
      }
    };

    fetchChoices(query);
  }, [query, datasource]);

  return (
    <>
      <QueryEditor query={query} datasource={datasource} onChange={onChange} onRunQuery={onRunQuery} />
      <InlineFieldRow className={styles.rowSpacing}>
        <InlineFormLabel
          width={10}
          tooltip="This field determines the value used for the variable"
          className="query-keyword"
        >
          Value Field
        </InlineFormLabel>
        <Select
          data-testid="value-field-select"
          allowCustomValue
          value={query.valueField}
          onChange={(opt: SelectableValue<string>) => onChange({ ...query, valueField: opt.value ?? '' })}
          width={64}
          placeholder={loading ? 'Loading...' : 'Select'}
          options={choices.map((opt) => ({ label: opt, value: opt }))}
        />
      </InlineFieldRow>
      <InlineFieldRow className={styles.rowSpacing}>
        <InlineFormLabel
          width={10}
          tooltip="This field determines the text used for the variable"
          className="query-keyword"
        >
          Label Field
        </InlineFormLabel>
        <Select
          data-testid="label-field-select"
          allowCustomValue
          value={query.labelField}
          onChange={(opt: SelectableValue<string>) => onChange({ ...query, labelField: opt.value ?? '' })}
          width={64}
          placeholder={loading ? 'Loading...' : 'Select'}
          options={choices.map((opt) => ({ label: opt, value: opt }))}
        />
      </InlineFieldRow>
      <>
        <Label className={styles.filterSection}>Optional filtering</Label>
        <InlineFieldRow className={styles.rowSpacing}>
          <InlineFormLabel width={10} tooltip="Select the column to filter on" className="query-keyword">
            Filter Field
          </InlineFormLabel>
          <Select
            data-testid="filter-field-select"
            value={query.filterField}
            onChange={(opt: SelectableValue<string>) => onChange({ ...query, filterField: opt.value })}
            width={64}
            placeholder={loading ? 'Loading...' : 'Select'}
            options={choices.map((opt) => ({ label: opt, value: opt }))}
            allowCustomValue
          />
        </InlineFieldRow>
        <InlineFieldRow className={styles.rowSpacing}>
          <InlineFormLabel
            width={10}
            tooltip="Enter the value to filter for in the selected column"
            className="query-keyword"
          >
            Filter Value
          </InlineFormLabel>
          <Input
            data-testid="filter-value-input"
            value={query.filterValue || ''}
            onChange={(e) => onChange({ ...query, filterValue: e.currentTarget.value })}
            width={64}
            placeholder="Enter filter value"
          />
        </InlineFieldRow>
      </>
    </>
  );
};

export default VariableQueryEditor;

const getStyles = (theme: GrafanaTheme2) => {
  return {
    rowSpacing: css({
      marginBottom: theme.spacing(0.5),
    }),
    filterSection: css({
      marginTop: theme.spacing(2),
      marginBottom: theme.spacing(1),
    }),
  };
};

// This is used to create a request for the variable query editor
// We need to add this to satisfy the type checker, but it is not important for executed queries
const createRequest = (query: SheetsVariableQuery) => {
  return {
    targets: [{ ...query, refId: 'metricFindQuery' }],
    range: getDefaultTimeRange(),
    requestId: 'metricFindQuery',
    interval: '1s',
    intervalMs: 1000,
    timezone: 'browser',
    panelId: 1,
    dashboardUID: '1',
    scopedVars: {},
    startTime: Date.now(),
    app: CoreApp.Unknown,
  };
};
