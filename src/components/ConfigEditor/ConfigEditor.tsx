import { DataSourcePluginOptionsEditorProps, onUpdateDatasourceJsonDataOptionSelect } from '@grafana/data';
import { InlineFormLabel, Select } from '@grafana/ui';
import React from 'react';
import { GoogleAuthType, googleAuthTypes, GoogleSheetsSecureJsonData, SheetsSourceOptions } from '../../types';
import { APIAuth } from './APIAuth';
import { JWTAuth } from './JWTAuth';

export type Props = DataSourcePluginOptionsEditorProps<SheetsSourceOptions, GoogleSheetsSecureJsonData>;

export function ConfigEditor(props: Props) {
  const { options, onOptionsChange } = props;
  const { jsonData } = options;

  if (!jsonData.hasOwnProperty('authType')) {
    jsonData.authType = GoogleAuthType.KEY;
  }

  const renderBody = () => {
    switch (jsonData.authType) {
      case GoogleAuthType.JWT:
        return <JWTAuth onOptionsChange={onOptionsChange} options={options} />;
      case GoogleAuthType.KEY:
        return <APIAuth onOptionsChange={onOptionsChange} options={options} />;
    }
  };

  return (
    <div className="gf-form-group">
      <div className="gf-form">
        <InlineFormLabel className="width-10">Auth</InlineFormLabel>
        <Select
          className="width-30"
          value={googleAuthTypes.find((x) => x.value === jsonData.authType) || googleAuthTypes[0]}
          options={googleAuthTypes}
          defaultValue={jsonData.authType}
          onChange={onUpdateDatasourceJsonDataOptionSelect(props, 'authType')}
        />
      </div>
      {renderBody()}
    </div>
  );
}
