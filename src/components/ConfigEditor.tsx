import {
  DataSourcePluginOptionsEditorProps,
  onUpdateDatasourceSecureJsonDataOption,
  SelectableValue,
} from '@grafana/data';
import { AuthConfig } from '@grafana/google-sdk';
import { DataSourceDescription } from '@grafana/plugin-ui';
import { Field, SecretInput, SegmentAsync, Divider } from '@grafana/ui';
import React, { useState, useEffect } from 'react';
import { GoogleSheetsSecureJSONData, googleSheetsAuthTypes, GoogleSheetsAuth, GoogleSheetsDataSourceOptions } from '../types';
import { getBackwardCompatibleOptions } from '../utils';
import { ConfigurationHelp } from './ConfigurationHelp';
import { getDataSourceSrv } from '@grafana/runtime';
import { DataSource } from '../DataSource';

export type Props = DataSourcePluginOptionsEditorProps<GoogleSheetsDataSourceOptions, GoogleSheetsSecureJSONData>;

export function ConfigEditor(props: Props) {
  const options = getBackwardCompatibleOptions(props.options);
  const [selectedSheetOption, setSelectedSheetOption] = useState<SelectableValue<string> | string | undefined>(
    options.jsonData.defaultSheetID
  );

  const apiKeyProps = {
    isConfigured: Boolean(options.secureJsonFields.apiKey),
    value: options.secureJsonData?.apiKey || '',
    placeholder: 'Enter API key',
    id: 'apiKey',
    onReset: () =>
      props.onOptionsChange({
        ...options,
        secureJsonFields: { ...options.secureJsonFields, apiKey: false },
        secureJsonData: { apiKey: '' },
        jsonData: options.jsonData,
      }),
    onChange: onUpdateDatasourceSecureJsonDataOption(props, 'apiKey'),
  };

  const loadSheetIDs = async () => {
    if (!options.uid) {
      return [];
    }
    try {
      const ds = (await getDataSourceSrv().get(options.uid)) as DataSource;
      return ds.getSpreadSheets();
    } catch {
      return [];
    }
  };

  useEffect(() => {
    const currentValue = options.jsonData.defaultSheetID;
    if (!currentValue || !options.uid) {
      setSelectedSheetOption(currentValue);
      return;
    }
    const updateSelectedOption = async () => {
      try {
        const ds = (await getDataSourceSrv().get(options.uid!)) as DataSource;
        const sheetOptions = await ds.getSpreadSheets();
        const matchingOption = sheetOptions.find((opt) => opt.value === currentValue);
        setSelectedSheetOption(matchingOption || currentValue);
      } catch {
        setSelectedSheetOption(currentValue);
      }
    };
    updateSelectedOption();
  }, [options.jsonData.defaultSheetID, options.uid]);
  return (
    <>
      <DataSourceDescription
        dataSourceName="Google Sheets"
        docsLink="https://grafana.com/docs/plugins/grafana-googlesheets-datasource/latest/"
        hasRequiredFields={false}
      />

      <Divider />
      <div className="grafana-info-box">
        <h5>Choosing an authentication type</h5>
        <ul>
          <li><strong>Google JWT File</strong>: provides access to private spreadsheets and works in all environments where Grafana is running.</li> 
          <li><strong>API Key</strong>: simpler configuration, but requires spreadsheets to be public.</li>
          <li><strong>GCE Default Service Account</strong>: automatically retrieves default credentials. Requires Grafana to be running on a Google Compute Engine virtual machine.</li>
        </ul>
        <br/>
        <p><strong>Select an Authentication type below and expand <strong>Configure Google Sheets Authentication</strong> for 
          detailed guidance on configuration</strong>.
        </p>
      </div>
      <ConfigurationHelp authenticationType={options.jsonData.authenticationType} />

      <Divider />

      <AuthConfig authOptions={googleSheetsAuthTypes} onOptionsChange={props.onOptionsChange} options={options} />

      {options.jsonData.authenticationType === GoogleSheetsAuth.API && (
        <Field label="API Key">
          <SecretInput {...apiKeyProps} label="API key" width={40} />
        </Field>
      )}

      <Divider />

      <Field
        label="Default Spreadsheet ID"
        description="Optional spreadsheet ID to use as default when creating new queries"
      >
        <SegmentAsync
          loadOptions={loadSheetIDs}
          placeholder="Select Spreadsheet ID"
          value={selectedSheetOption}
          allowCustomValue={true}
          onChange={(value) => {
            const sheetId = typeof value === 'string' ? value : value?.value;
            setSelectedSheetOption(value);
            props.onOptionsChange({
              ...options,
              jsonData: {
                ...options.jsonData,
                defaultSheetID: sheetId,
              },
            });
          }}
        />
      </Field>
    </>
  );
}
