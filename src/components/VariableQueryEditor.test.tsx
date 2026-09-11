import React from 'react';
import { render, screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { of } from 'rxjs';
import { CoreApp } from '@grafana/data';
import VariableQueryEditor from './VariableQueryEditor';
import { DataSource } from '../DataSource';
import { SheetsVariableQuery } from '../types';

// Mock @grafana/ui's Combobox with a plain input + option list.
//
// The real Combobox (@grafana/ui 13.1.5) uses floating-ui to position its
// popover. Under jsdom there's no layout engine, so every position
// measurement comes back as zero, which sends floating-ui's positioning
// middleware into a render loop. With three Comboboxes in this component
// all updating their `options` after mount, that loop never settles and
// pins the process at 100% CPU until Jest is killed. This is a known
// @grafana/ui bug that only reproduces under jsdom; upgrading past 13.1.5
// requires React 19, which this plugin doesn't use yet. Remove this mock
// once @grafana/ui is upgraded past 13.1.5 (or the plugin moves to React 19).
jest.mock('@grafana/ui', () => {
  const actualReact = jest.requireActual('react');
  const actual = jest.requireActual('@grafana/ui');

  const Combobox = ({
    value,
    options,
    onChange,
    placeholder,
    'data-testid': testId,
  }: {
    value?: string;
    options: Array<{ label?: string; value?: string }>;
    onChange: (opt: { label?: string; value?: string }) => void;
    placeholder?: string;
    'data-testid'?: string;
  }) => {
    const [open, setOpen] = actualReact.useState(false);
    return (
      <div>
        <input
          data-testid={testId}
          value={value ?? ''}
          placeholder={placeholder}
          readOnly
          onClick={() => setOpen((o: boolean) => !o)}
        />
        {open && (
          <ul>
            {options.map((opt) => (
              <li
                key={opt.value}
                onClick={() => {
                  onChange(opt);
                  setOpen(false);
                }}
              >
                {opt.label}
              </li>
            ))}
          </ul>
        )}
      </div>
    );
  };

  return { ...actual, Combobox };
});

// Mock the dependencies
jest.mock('@grafana/data', () => ({
  ...jest.requireActual('@grafana/data'),
  getDefaultTimeRange: jest.fn().mockReturnValue({
    from: { valueOf: () => 1000 },
    to: { valueOf: () => 2000 },
    raw: { from: 'now-1h', to: 'now' },
  }),
  CoreApp: {
    Unknown: 'unknown',
  },
}));

// Create a mock datasource factory
const createMockDataSource = (overrides: Partial<DataSource> = {}): DataSource => {
  return {
    query: jest.fn().mockReturnValue(
      of({
        data: [
          {
            fields: [
              { name: 'defaultValueField', values: ['item1', 'item2', 'item3'] },
              { name: 'defaultLabelField', values: ['Item 1', 'Item 2', 'Item 3'] },
              { name: 'baz', values: ['Category 1', 'Category 2', 'Category 3'] },
              { name: 'qux', values: ['Category 1', 'Category 2', 'Category 3'] },
            ],
          },
        ],
      })
    ),
    getSpreadSheets: jest.fn(),
    authType: 'key',
    ...overrides,
  } as DataSource;
};

describe('VariableQueryEditor', () => {
  let props: {
    query: SheetsVariableQuery;
    onChange: jest.Mock;
    onRunQuery: jest.Mock;
    datasource: DataSource;
  };

  beforeEach(() => {
    jest.clearAllMocks();
    props = {
      datasource: createMockDataSource(),
      query: {
        refId: 'A',
        cacheDurationSeconds: 300,
        spreadsheet: 'test-sheet-id',
        valueField: 'defaultValueField',
        labelField: 'defaultLabelField',
        filterField: '',
        filterValue: '',
      },
      onChange: jest.fn(),
      onRunQuery: jest.fn(),
    };
  });

  test('renders without crashing', async () => {
    render(<VariableQueryEditor {...props} />);
    expect(await screen.findByTestId('value-field-select')).toBeInTheDocument();
    expect(await screen.findByTestId('label-field-select')).toBeInTheDocument();
    expect(await screen.findByTestId('filter-field-select')).toBeInTheDocument();
    expect(await screen.findByTestId('filter-value-input')).toBeInTheDocument();
  });

  test('fetches column choices on mount', async () => {
    render(<VariableQueryEditor {...props} />);
    await waitFor(() => {
      expect(props.datasource.query).toHaveBeenCalledWith(
        expect.objectContaining({
          targets: [{ ...props.query, refId: 'metricFindQuery' }],
          requestId: 'metricFindQuery',
          app: CoreApp.Unknown,
        })
      );
    });
    expect(await screen.findByDisplayValue('defaultValueField')).toBeInTheDocument();
    expect(await screen.findByDisplayValue('defaultLabelField')).toBeInTheDocument();
    expect(screen.queryByDisplayValue('baz')).not.toBeInTheDocument();
  });

  test('shows loading state while fetching choices', async () => {
    render(<VariableQueryEditor {...props} />);
    // Should show loading state
    expect(await screen.findAllByPlaceholderText('Loading...')).toHaveLength(3); // All three selects show loading
    // Loading should disappear
    await waitFor(() => {
      expect(screen.queryByPlaceholderText('Loading...')).not.toBeInTheDocument();
    });
  });

  test('shows options for selecting label field', async () => {
    render(<VariableQueryEditor {...props} />);
    await waitFor(() => {
      expect(props.datasource.query).toHaveBeenCalled();
    });

    // bar is the default value
    expect(await screen.findByDisplayValue('defaultLabelField')).toBeInTheDocument();
    expect(screen.queryByDisplayValue('baz')).not.toBeInTheDocument();

    userEvent.click(await screen.findByDisplayValue('defaultLabelField'));
    expect(await screen.findByText('baz')).toBeInTheDocument();
    expect(await screen.findByText('qux')).toBeInTheDocument();
  });

  test('shows options for selecting value field', async () => {
    render(<VariableQueryEditor {...props} />);
    await waitFor(() => {
      expect(props.datasource.query).toHaveBeenCalled();
    });

    // bar is the default value
    expect(await screen.findByDisplayValue('defaultValueField')).toBeInTheDocument();
    expect(screen.queryByDisplayValue('baz')).not.toBeInTheDocument();

    userEvent.click(await screen.findByDisplayValue('defaultValueField'));
    expect(await screen.findByText('baz')).toBeInTheDocument();
    expect(await screen.findByText('qux')).toBeInTheDocument();
  });

  test('updates the query when the value field is changed', async () => {
    render(<VariableQueryEditor {...props} />);
    userEvent.click(await screen.findByDisplayValue('defaultValueField'));
    userEvent.click(await screen.findByText('qux'));
    await waitFor(() => {
      expect(props.onChange).toHaveBeenCalledWith({
        ...props.query,
        valueField: 'qux',
      });
    });
  });

  test('updates the query when the label field is changed', async () => {
    render(<VariableQueryEditor {...props} />);
    userEvent.click(await screen.findByDisplayValue('defaultLabelField'));
    userEvent.click(await screen.findByText('baz'));
    await waitFor(() => {
      expect(props.onChange).toHaveBeenCalledWith({
        ...props.query,
        labelField: 'baz',
      });
    });
  });

  test('shows options for selecting filter field', async () => {
    render(<VariableQueryEditor {...props} />);
    await waitFor(() => {
      expect(props.datasource.query).toHaveBeenCalled();
    });

    // Should show the filter field select
    const filterFieldSelect = await screen.findByTestId('filter-field-select');
    expect(filterFieldSelect).toBeInTheDocument();

    // Click to open dropdown
    userEvent.click(filterFieldSelect);
    expect(await screen.findByText('baz')).toBeInTheDocument();
    expect(await screen.findByText('qux')).toBeInTheDocument();
  });

  test('updates the query when the filter field is changed', async () => {
    render(<VariableQueryEditor {...props} />);
    const filterFieldSelect = await screen.findByTestId('filter-field-select');

    userEvent.click(filterFieldSelect);
    userEvent.click(await screen.findByText('baz'));

    await waitFor(() => {
      expect(props.onChange).toHaveBeenCalledWith({
        ...props.query,
        filterField: 'baz',
      });
    });
  });

  test('updates the query when the filter value is changed', async () => {
    render(<VariableQueryEditor {...props} />);
    const filterValueInput = await screen.findByTestId('filter-value-input');

    userEvent.type(filterValueInput, 'a');
    await waitFor(() => {
      expect(props.onChange).toHaveBeenCalledWith({
        ...props.query,
        filterValue: 'a',
      });
    });
  });

  test('shows loading state for filter field select', async () => {
    render(<VariableQueryEditor {...props} />);
    // Should show loading state for all selects including filter field
    expect(await screen.findAllByPlaceholderText('Loading...')).toHaveLength(3); // Value, Label, and Filter field selects
    // Loading should disappear
    await waitFor(() => {
      expect(screen.queryByPlaceholderText('Loading...')).not.toBeInTheDocument();
    });
  });
});
