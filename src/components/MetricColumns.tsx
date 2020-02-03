import React, { FunctionComponent } from 'react';
import { SelectableValue } from '@grafana/data';
import { SegmentAsync } from '@grafana/ui';

export interface Props {
  values: Array<SelectableValue<number>>;
  onChange: (values: Array<SelectableValue<number>>) => void;
  loadColumns: () => Promise<Array<SelectableValue<number>>>;
}

const removeText = '-- remove metric --';
const removeOption: SelectableValue<number> = { label: removeText, value: -1 };

export const MetricColumns: FunctionComponent<Props> = ({ loadColumns, values, onChange }) => (
  <>
    {values &&
      values.map((value, index) => (
        <SegmentAsync
          allowCustomValue
          key={index}
          value={value}
          loadOptions={() => loadColumns().then(options => [removeOption, ...options.filter(({ label }) => label !== removeText)])}
          onChange={option =>
            onChange(option.label === removeText ? values.filter((_, i) => i !== index) : values.map((v, i) => (i === index ? option : v)))
          }
        />
      ))}
    <SegmentAsync
      Component={
        <a className="gf-form-label query-part">
          <i className="fa fa-plus" />
        </a>
      }
      allowCustomValue
      onChange={option => onChange([...values, option])}
      loadOptions={loadColumns}
    />
  </>
);
