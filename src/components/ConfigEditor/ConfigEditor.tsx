import {
  DataSourcePluginOptionsEditorProps,
  onUpdateDatasourceJsonDataOptionChecked,
  onUpdateDatasourceJsonDataOptionSelect,
} from '@grafana/data';
import { InlineFormLabel, LegacyForms, Select } from '@grafana/ui';
import React from 'react';
import { GoogleAuthType, googleAuthTypes, GoogleSheetsSecureJsonData, SheetsSourceOptions } from '../../types';
import { APIAuth } from './APIAuth';
import { JWTAuth } from './JWTAuth';
import { OAuth } from './OAuth';

export type Props = DataSourcePluginOptionsEditorProps<SheetsSourceOptions, GoogleSheetsSecureJsonData>;

export function ConfigEditor(props: Props) {
  const { jsonData } = props.options;

  if (!jsonData.hasOwnProperty('authType')) {
    jsonData.authType = GoogleAuthType.KEY;
  }

  const renderBody = () => {
    switch (jsonData.authType) {
      case GoogleAuthType.JWT:
        return <JWTAuth {...props} />;
      case GoogleAuthType.KEY:
        return <APIAuth {...props} />;
      case GoogleAuthType.OAUTH:
        return <OAuth {...props} />;
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
      <div className="gf-form">
        <LegacyForms.Switch
          label="Forward OAuth Identity"
          labelClass="width-13"
          checked={jsonData.oauthPassThru || false}
          onChange={onUpdateDatasourceJsonDataOptionChecked(props, 'oauthPassThru')}
          tooltip="Forward the user's upstream OAuth identity to the data source (Their access token gets passed along)."
        />
      </div>
      {renderBody()}
    </div>
  );
}
