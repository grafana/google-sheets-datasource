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

  it('filters response data based on filterField and filterValue', () => {
    const query: SheetsVariableQuery = {
      valueField: 'id',
      labelField: 'name',
      filterField: 'squad',
      filterValue: 'frontend',
      spreadsheet: 'test-sheet',
      range: 'A1:D4',
      cacheDurationSeconds: 300,
      refId: 'A',
    };

    const response = {
      data: [
        toDataFrame({
          fields: [
            { name: 'id', values: ['1', '2', '3', '4'] },
            { name: 'name', values: ['Plugin A', 'Plugin B', 'Plugin C', 'Plugin D'] },
            { name: 'squad', values: ['frontend', 'backend', 'frontend', 'infra'] },
          ],
        }),
      ],
    };

    const result = queryResponseToVariablesFrame(query, response);

    expect(result.data).toEqual([
      { value: '1', text: 'Plugin A' },
      { value: '3', text: 'Plugin C' },
    ]);
  });

  test('returns all data when no filter is applied', () => {
    const query: SheetsVariableQuery = {
      valueField: 'id',
      labelField: 'name',
      spreadsheet: 'test-sheet',
      range: 'A1:D4',
      cacheDurationSeconds: 300,
      refId: 'A',
    };

    const response = {
      data: [
        toDataFrame({
          fields: [
            { name: 'id', values: ['1', '2', '3', '4'] },
            { name: 'name', values: ['Plugin A', 'Plugin B', 'Plugin C', 'Plugin D'] },
            { name: 'squad', values: ['frontend', 'backend', 'frontend', 'infra'] },
          ],
        }),
      ],
    };

    const result = queryResponseToVariablesFrame(query, response);

    expect(result.data).toEqual([
      { value: '1', text: 'Plugin A' },
      { value: '2', text: 'Plugin B' },
      { value: '3', text: 'Plugin C' },
      { value: '4', text: 'Plugin D' },
    ]);
  });

  test('returns empty data when filter matches no rows', () => {
    const query: SheetsVariableQuery = {
      valueField: 'id',
      labelField: 'name',
      filterField: 'squad',
      filterValue: 'nonexistent',
      spreadsheet: 'test-sheet',
      range: 'A1:D4',
      cacheDurationSeconds: 300,
      refId: 'A',
    };

    const response = {
      data: [
        toDataFrame({
          fields: [
            { name: 'id', values: ['1', '2', '3', '4'] },
            { name: 'name', values: ['Plugin A', 'Plugin B', 'Plugin C', 'Plugin D'] },
            { name: 'squad', values: ['frontend', 'backend', 'frontend', 'infra'] },
          ],
        }),
      ],
    };

    const result = queryResponseToVariablesFrame(query, response);

    expect(result.data).toEqual([]);
  });

  test('ignores filter when filterField is specified but filterValue is empty', () => {
    const query: SheetsVariableQuery = {
      valueField: 'id',
      labelField: 'name',
      filterField: 'squad',
      filterValue: '',
      spreadsheet: 'test-sheet',
      range: 'A1:D4',
      cacheDurationSeconds: 300,
      refId: 'A',
    };

    const response = {
      data: [
        toDataFrame({
          fields: [
            { name: 'id', values: ['1', '2', '3', '4'] },
            { name: 'name', values: ['Plugin A', 'Plugin B', 'Plugin C', 'Plugin D'] },
            { name: 'squad', values: ['frontend', 'backend', 'frontend', 'infra'] },
          ],
        }),
      ],
    };

    const result = queryResponseToVariablesFrame(query, response);

    expect(result.data).toEqual([
      { value: '1', text: 'Plugin A' },
      { value: '2', text: 'Plugin B' },
      { value: '3', text: 'Plugin C' },
      { value: '4', text: 'Plugin D' },
    ]);
  });
});
