import React, { PureComponent, ChangeEvent } from 'react';

import { MyQuery, ScenarioProvider, Scenario, ScenarioEditorProps } from 'types';
import { toFloatOrUndefined, toNumberString } from '@grafana/data';
import { FormField } from '@grafana/ui';

interface WaveQuery extends MyQuery {
  period: number;
}

type Props = ScenarioEditorProps<WaveQuery>;
type State = {};

class WaveEditor extends PureComponent<Props, State> {
  state = {};

  onPeriodChange = (event: ChangeEvent<HTMLInputElement>) => {
    const v = toFloatOrUndefined(event.target.value);
    if (v === undefined) {
      return; // skip undefined
    }
    this.props.onChange({
      ...this.props.query,
      period: v,
    });
  };

  render() {
    const { query } = this.props;
    return (
      <div className="gf-form">
        <FormField label="Period" labelWidth={5} onChange={this.onPeriodChange} value={toNumberString(query.period) || ''} type="number" />
      </div>
    );
  }
}

export const waveProvider: ScenarioProvider<WaveQuery> = {
  value: Scenario.wave,
  label: 'Wave',
  description: 'This is a wave...',
  defaultQuery: {
    period: 6000, // 1min
  },
  editor: WaveEditor,
};
