import React, { PureComponent, ChangeEvent } from 'react';
import { MyQuery, ScenarioProvider, Scenario, ScenarioEditorProps } from 'types';
import { toIntegerOrUndefined } from '@grafana/data';
import { FormField } from '@grafana/ui';

interface CSVWaveQuery extends MyQuery {
  csvValues: string;
  timeStep: number;
  shift: number;
  phase: number;
}

type Props = ScenarioEditorProps<CSVWaveQuery>;
type State = {};

class CSVWaveEditor extends PureComponent<Props, State> {
  state = {};

  onCSVChange = (event: ChangeEvent<HTMLInputElement>) => {
    const v = event.target.value;
    if (v === undefined) {
      return; // skip undefined
    }
    this.props.onChange({
      ...this.props.query,
      csvValues: v,
    });
  };

  onTimeStepChange = (event: ChangeEvent<HTMLInputElement>) => {
    const v = toIntegerOrUndefined(event.target.value);
    if (v === undefined) {
      return; // skip undefined
    }
    this.props.onChange({
      ...this.props.query,
      timeStep: v,
    });
  };

  onShiftChange = (event: ChangeEvent<HTMLInputElement>) => {
    const v = toIntegerOrUndefined(event.target.value);
    if (v === undefined) {
      return; // skip undefined
    }
    this.props.onChange({
      ...this.props.query,
      shift: v,
    });
  };

  onPhaseChange = (event: ChangeEvent<HTMLInputElement>) => {
    const v = toIntegerOrUndefined(event.target.value);
    if (v === undefined) {
      return; // skip undefined
    }
    this.props.onChange({
      ...this.props.query,
      phase: v,
    });
  };

  render() {
    const { query } = this.props;
    return (
      <div className="gf-form">
        <FormField label="Time Step" labelWidth={7} inputWidth={4} onChange={this.onTimeStepChange} value={query.timeStep || 5} type="number" />
        <FormField label="Shift" labelWidth={7} inputWidth={4} onChange={this.onShiftChange} value={query.shift || 0} type="number" />
        <FormField label="Phase" labelWidth={7} inputWidth={4} onChange={this.onPhaseChange} value={query.phase || 0} type="number" />
        <FormField label="CSV Values" labelWidth={6} inputWidth={30} onChange={this.onCSVChange} value={query.csvValues || ''} type="string" />
      </div>
    );
  }
}

export const csvWaveProvider: ScenarioProvider<CSVWaveQuery> = {
  value: Scenario.csvWave,
  label: 'CSV Wave',
  description: 'Create a repeating wave of the CSVValues',
  defaultQuery: {
    csvValues: '1, 2.2, NaN, Null, 3',
    timeStep: 1,
    shift: 0,
    phase: 0,
  },
  editor: CSVWaveEditor,
};
