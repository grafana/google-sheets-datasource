import { DataSourceSettings, onUpdateDatasourceJsonDataOption } from '@grafana/data';
import { Button, InlineFormLabel, Input } from '@grafana/ui';
import { useLoadGapi } from 'components/useLoadGapi';
import React, { useCallback, useEffect, useState } from 'react';
import { GoogleSheetsSecureJsonData, SheetsSourceOptions } from 'types';

type Props = {
  options: DataSourceSettings<SheetsSourceOptions, GoogleSheetsSecureJsonData>;
  onOptionsChange: (options: DataSourceSettings<SheetsSourceOptions, GoogleSheetsSecureJsonData>) => void;
};

export function OAuth(props: Props) {
  const {
    options: {
      jsonData: { developerKey, appId, clientId },
    },
  } = props;
  const [currentUser, setCurrentUser] = useState<gapi.auth2.BasicProfile>();

  useLoadGapi(() => {
    gapi.load('client', {
      callback: () => {
        initClient();
      },
    });
  });

  const initClient = useCallback(() => {
    if (!clientId || !appId || !developerKey) {
      return;
    }
    gapi.client
      .init({
        apiKey: developerKey,
        clientId: clientId,
        scope: 'https://www.googleapis.com/auth/drive.file',
      })
      .then(function () {
        if (gapi.auth2.getAuthInstance().isSignedIn.get()) {
          setCurrentUser(gapi.auth2.getAuthInstance().currentUser.get().getBasicProfile());
        } else {
          setCurrentUser(undefined);
        }
      });
  }, [appId, clientId, developerKey]);

  useEffect(() => {
    if (window.gapi) {
      initClient();
    }
  }, [initClient]);

  return (
    <>
      <div className="gf-form">
        <InlineFormLabel className="width-10">Developer key / API key</InlineFormLabel>
        <Input
          css={{}}
          className="width-30"
          value={developerKey || ''}
          onChange={onUpdateDatasourceJsonDataOption(props, 'developerKey')}
        />
      </div>
      <div className="gf-form">
        <InlineFormLabel className="width-10">App ID</InlineFormLabel>
        <Input
          css={{}}
          className="width-30"
          value={appId || ''}
          onChange={onUpdateDatasourceJsonDataOption(props, 'appId')}
        />
      </div>
      <div className="gf-form">
        <InlineFormLabel className="width-10">Client ID</InlineFormLabel>
        <Input
          css={{}}
          className="width-30"
          value={clientId || ''}
          onChange={onUpdateDatasourceJsonDataOption(props, 'clientId')}
        />
      </div>
      <div className="gf-form" style={{ alignItems: 'center' }}>
        <InlineFormLabel className="width-10">Google Sign In</InlineFormLabel>
        {currentUser ? (
          <>
            <div style={{ margin: '0 8px' }}>
              You are logged in with:{' '}
              <img
                referrerPolicy="no-referrer"
                src={currentUser.getImageUrl()}
                alt="avatar icon"
                style={{ borderRadius: 50, width: 24 }}
              />{' '}
              {currentUser.getName()} ({currentUser.getEmail()})
            </div>
            <Button
              type="button"
              onClick={() => {
                gapi.auth2.getAuthInstance().signOut();
                setCurrentUser(undefined);
              }}
            >
              Sign Out
            </Button>
          </>
        ) : developerKey && appId && clientId ? (
          <Button
            type="button"
            onClick={() => {
              gapi.auth2
                .getAuthInstance()
                .signIn()
                .then((user) => {
                  setCurrentUser(user.getBasicProfile());
                });
            }}
          >
            Authorize
          </Button>
        ) : null}
      </div>
      <div className="grafana-info-box" style={{ marginTop: 24 }}>
        <h4>Using OAuth</h4>
      </div>
    </>
  );
}
