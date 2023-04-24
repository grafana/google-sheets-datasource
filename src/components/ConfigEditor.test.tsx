import React from 'react';
import { render, screen, fireEvent } from '@testing-library/react';
import { ConfigEditor } from './ConfigEditor';

jest.mock('@grafana/runtime', () => ({
  getDataSourceSrv: () => ({
    get: Promise.resolve({
      getSpreadSheets: jest.fn().mockImplementation(() =>
        Promise.resolve([
          { label: 'label1', value: 'value1' },
          { label: 'label2', value: 'value2' },
        ])
      ),
    }),
  }),
}));

describe('ConfigEditor', () => {
  it('should support old authType property', () => {
    const onOptionsChange = jest.fn();
    // Render component with old authType property
    render(
      <ConfigEditor
        onOptionsChange={onOptionsChange}
        options={{ jsonData: { authType: 'key', authenticationType: '' }, secureJsonFields: {} } as any}
      />
    );

    // Check that the correct auth type is selected
    expect(screen.getByRole('radio', { name: 'API Key' })).toBeChecked();

    // Make sure that the user can still change the auth type
    fireEvent.click(screen.getByLabelText('Google JWT File'));

    // Check onOptionsChange is called with the correct value
    expect(onOptionsChange).toHaveBeenCalledWith({
      jsonData: { authType: 'key', authenticationType: 'jwt' },
      secureJsonFields: {},
    });
  });

  it('should be backward compatible with API Key', () => {
    render(
      <ConfigEditor
        onOptionsChange={jest.fn()}
        options={{ jsonData: { authType: 'key', authenticationType: '' }, secureJsonFields: { apiKey: true } } as any}
      />
    );

    // Check that the correct auth type is selected
    expect(screen.getByRole('radio', { name: 'API Key' })).toBeChecked();

    // Check that the API key is configured
    expect(screen.getByPlaceholderText('Enter API key')).toHaveAttribute('value', 'configured');
  });

  it('should be backward compatible with JWT auth type', () => {
    render(
      <ConfigEditor
        onOptionsChange={jest.fn()}
        options={{ jsonData: { authType: 'jwt', authenticationType: '' }, secureJsonFields: { jwt: true } } as any}
      />
    );

    // Check that the correct auth type is selected
    expect(screen.getByLabelText('Google JWT File')).toBeChecked();

    // Check that the Private key input is configured
    expect(screen.getByTestId('Private Key Input')).toHaveAttribute('value', 'configured');
  });

  it('should render default spreadsheet ID field', async () => {
    const onChange = jest.fn();
    render(
      <ConfigEditor
        onOptionsChange={onChange}
        options={{ jsonData: { authenticationType: 'key' }, secureJsonFields: {} } as any}
      />
    );
    expect(screen.getByText('Default spreadsheet ID')).toBeInTheDocument();
  });

  it('should save default spreadsheet ID', async () => {
    // TODO: Implement test
  });
});
