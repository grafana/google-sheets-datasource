import React, { PureComponent, ChangeEvent } from 'react';
import { MyQuery, ScenarioProvider, Scenario, ScenarioEditorProps } from 'types';
import { FormField } from '@grafana/ui';

interface ArrowFileQuery extends MyQuery {
  url: string;
}

type Props = ScenarioEditorProps<ArrowFileQuery>;
type State = {};

class ArrowFileEditor extends PureComponent<Props, State> {
  state = {};

  onURLChange = (event: ChangeEvent<HTMLInputElement>) => {
    const v = event.target.value;
    if (v === undefined) {
      return; // skip undefined
    }
    this.props.onChange({
      ...this.props.query,
      url: v,
    });
  };

  render() {
    const { query } = this.props;
    return (
      <div className="gf-form">
        <FormField label="url" labelWidth={5} inputWidth={30} onChange={this.onURLChange} value={query.url || ''} type="string" />
      </div>
    );
  }
}

export const arrowFileProvider: ScenarioProvider<ArrowFileQuery> = {
  value: Scenario.arrowFile,
  label: 'Arrow File',
  description: 'Get an Arrow file via a url',
  defaultQuery: {
    url: 'https://github.com/grafana/grafana-plugin-sdk-go/blob/master/dataframe/testdata/all_types.golden.arrow?raw=true', // SDK's Golden Test File
  },
  editor: ArrowFileEditor,
};
