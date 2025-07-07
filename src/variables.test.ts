import { queryResponseToVariablesFrame } from './variables';
import { SheetsVariableQuery } from './types';
import { toDataFrame } from '@grafana/data';

describe('queryResponseToVariablesFrame', () => {
  test('transforms response data into value/text format', () => {
    const query: SheetsVariableQuery = {
      valueField: 'id',
      labelField: 'name',
      spreadsheet: 'test-sheet',
      range: 'A1:B3',
      cacheDurationSeconds: 300,
      refId: 'A',
    };

    const response = {
      data: [
        toDataFrame({
          fields: [
            { name: 'id', values: ['1', '2', '3'] },
            { name: 'name', values: ['Item 1', 'Item 2', 'Item 3'] },
          ],
        }),
      ],
    };

    const result = queryResponseToVariablesFrame(query, response);

    expect(result.data).toEqual([
      { value: '1', text: 'Item 1' },
      { value: '2', text: 'Item 2' },
      { value: '3', text: 'Item 3' },
    ]);
  });
});
