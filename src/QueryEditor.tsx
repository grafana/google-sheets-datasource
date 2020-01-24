import React, { PureComponent } from 'react';
import { QueryEditorProps, SelectableValue } from '@grafana/data';
import { FormLabel, Select } from '@grafana/ui';
import { DataSource } from './DataSource';
import { MyQuery, MyDataSourceOptions, Scenario } from './types';
import { scenarios, scenarioList } from './scenarios';
import defaults from 'lodash/defaults';

type Props = QueryEditorProps<DataSource, MyQuery, MyDataSourceOptions>;

interface State {}

export class QueryEditor extends PureComponent<Props, State> {
  onComponentDidMount() {}

  onSelectScenario = (item: SelectableValue<Scenario>) => {
    const scenario = scenarios[item.value!] || scenarios.wave;
    this.props.onChange({
      ...scenario.defaultQuery,
      refId: this.props.query.refId,
      scenario: scenario.value!,
    });
  };

  render() {
    const scenario = scenarios[this.props.query.scenario] || scenarios.wave;
    const query = defaults(this.props.query, scenario.defaultQuery);

    return (
      <>
        <div className="gf-form">
          <div className="form-field">
            <FormLabel width={5}>Scenario</FormLabel>
            <Select options={scenarioList} value={scenario} onChange={this.onSelectScenario} />
          </div>
        </div>
        <scenario.editor {...this.props} query={query} />
      </>
    );
  }
}
