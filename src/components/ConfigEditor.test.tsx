import React from 'react';

import { render, screen } from '@testing-library/react';
import { ConfigEditor } from './ConfigEditor';

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
    expect(screen.getByLabelText('API Key')).toBeChecked();

    // Make sure that the user can still change the auth type
    screen.getByLabelText('Google JWT File').click();

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
    expect(screen.getByLabelText('API Key')).toBeChecked();

    // Check that the API key is configured
    expect(screen.getByLabelText('API key')).toHaveAttribute('value', 'configured');
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
});
