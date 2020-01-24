import React, { PureComponent, ChangeEvent } from 'react';

import { MyQuery, ScenarioProvider, Scenario, ScenarioEditorProps } from 'types';
import { toFloatOrUndefined, toNumberString } from '@grafana/data';
import { FormField } from '@grafana/ui';

interface NoiseQuery extends MyQuery {
  start: number;
}

type Props = ScenarioEditorProps<NoiseQuery>;
type State = {};

class NoiseEditor extends PureComponent<Props, State> {
  state = {};

  onStartChange = (event: ChangeEvent<HTMLInputElement>) => {
    const v = toFloatOrUndefined(event.target.value);
    if (v === undefined) {
      return; // skip undefined
    }
    this.props.onChange({
      ...this.props.query,
      start: v,
    });
  };

  render() {
    const { query } = this.props;
    return (
      <div className="gf-form">
        <FormField label="Start" labelWidth={5} onChange={this.onStartChange} value={toNumberString(query.start) || ''} type="number" />
      </div>
    );
  }
}

export const noiseProvider: ScenarioProvider<NoiseQuery> = {
  value: Scenario.noise,
  label: 'Noise',
  description: 'random noise',
  defaultQuery: {
    start: 0, // start at zero
  },
  editor: NoiseEditor,
};
