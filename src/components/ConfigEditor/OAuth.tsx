import {
  DataSourceSettings,
  onUpdateDatasourceJsonDataOption,
  onUpdateDatasourceSecureJsonDataOption,
} from '@grafana/data';
import { Button, InlineFormLabel, Input, LegacyForms } from '@grafana/ui';
import { useLoadGapi } from 'components/useLoadGapi';
import React, { useCallback, useEffect, useState } from 'react';
import { GoogleSheetsSecureJsonData, SheetsSourceOptions } from 'types';

type Props = {
  options: DataSourceSettings<SheetsSourceOptions, GoogleSheetsSecureJsonData>;
  onOptionsChange: (options: DataSourceSettings<SheetsSourceOptions, GoogleSheetsSecureJsonData>) => void;
};

export function OAuth(props: Props) {
  const { options, onOptionsChange } = props;
  const [isSignedIn, setIsSignedIn] = useState(false);

  useLoadGapi(() => {
    gapi.load('client:auth2', {
      callback: () => {
        initClient();
      },
    });
  });

  const onResetApiKey = () => {
    onOptionsChange({
      ...options,
      secureJsonData: {
        ...options.secureJsonData,
        apiKey: '',
      },
      secureJsonFields: {
        ...options.secureJsonFields,
        apiKey: false,
      },
    });
  };

  const initClient = useCallback(() => {
    if (!options.jsonData.clientId || !options.jsonData.appId || !options.secureJsonFields?.apiKey) {
      return;
    }
    gapi.client
      .init({
        apiKey: options.secureJsonData?.apiKey,
        clientId: options.jsonData.clientId,
        scope: 'https://www.googleapis.com/auth/drive.file',
      })
      .then(function () {
        // Listen for sign-in state changes.
        gapi.auth2.getAuthInstance().isSignedIn.listen((isSigned) => {
          setIsSignedIn(isSigned);
        });

        // Handle the initial sign-in state.
        setIsSignedIn(gapi.auth2.getAuthInstance().isSignedIn.get());
      });
  }, [
    options.jsonData.appId,
    options.jsonData.clientId,
    options.secureJsonData?.apiKey,
    options.secureJsonFields?.apiKey,
  ]);

  useEffect(() => {
    if (window.gapi) {
      initClient();
    }
  }, [initClient]);

  return (
    <>
      <div className="gf-form">
        <LegacyForms.SecretFormField
          isConfigured={options.secureJsonFields?.apiKey}
          value={options.secureJsonData?.apiKey || ''}
          label="API Key"
          labelWidth={10}
          inputWidth={30}
          placeholder="Enter API Key"
          onReset={onResetApiKey}
          onChange={onUpdateDatasourceSecureJsonDataOption(props, 'apiKey')}
        />
      </div>
      <div className="gf-form">
        <InlineFormLabel className="width-10">App ID</InlineFormLabel>
        <Input
          css={{}}
          className="width-30"
          value={options.jsonData.appId}
          onChange={onUpdateDatasourceJsonDataOption(props, 'appId')}
        />
      </div>
      <div className="gf-form">
        <InlineFormLabel className="width-10">Client ID</InlineFormLabel>
        <Input
          css={{}}
          className="width-30"
          value={options.jsonData.clientId}
          onChange={onUpdateDatasourceJsonDataOption(props, 'clientId')}
        />
      </div>
      <div className="gf-form">
        <InlineFormLabel className="width-10">Google Sign In</InlineFormLabel>
        {!isSignedIn ? (
          <Button
            type="button"
            onClick={() => {
              gapi.auth2.getAuthInstance().signIn();
            }}
          >
            Authorize
          </Button>
        ) : (
          <Button
            type="button"
            onClick={() => {
              gapi.auth2.getAuthInstance().signOut();
            }}
          >
            Sign Out
          </Button>
        )}
      </div>
      <div className="grafana-info-box" style={{ marginTop: 24 }}>
        <h4>Using OAuth</h4>
        {isSignedIn ? (
          <div>
            You are logged in with: {gapi.auth2.getAuthInstance().currentUser.get().getBasicProfile().getName()}
          </div>
        ) : null}
      </div>
    </>
  );
}
