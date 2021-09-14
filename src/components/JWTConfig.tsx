import { css } from '@emotion/css';
import { GrafanaTheme2 } from '@grafana/data';
import { Button, FileDropzone, FileListItem, InlineFormLabel, useStyles2 } from '@grafana/ui';
import { isObject, startCase } from 'lodash';
import React, { useState } from 'react';

const configKeys = [
  'type',
  'project_id',
  'private_key_id',
  'private_key',
  'client_email',
  'client_id',
  'auth_uri',
  'token_uri',
  'auth_provider_x509_cert_url',
  'client_x509_cert_url',
];

export interface Props {
  onChange: (jwt: string) => void;
  isConfigured: boolean;
}

const validateJson = (json: { [key: string]: string }) => isObject(json) && configKeys.every((key) => !!json[key]);

export function JWTConfig({ onChange, isConfigured }: Props) {
  const [enableUpload, setEnableUpload] = useState<boolean>(!isConfigured);
  const [error, setError] = useState<string>();
  const styles = useStyles2(getStyles);

  return enableUpload ? (
    <div className={styles.dropzone}>
      <FileDropzone
        options={{
          multiple: false,
        }}
        onLoad={(result) => {
          const json = JSON.parse(result as string);
          if (validateJson(json)) {
            onChange(result as string);
            setEnableUpload(false);
          } else {
            setError('Invalid JWT file');
          }
        }}
        fileListRenderer={(file, removeFile) => {
          if (error) {
            return <FileListItem file={{ ...file, error: new DOMException(error) }} removeFile={removeFile} />;
          }
          return <FileListItem file={file} removeFile={removeFile} />;
        }}
      />
    </div>
  ) : (
    <>
      {configKeys.map((key) => (
        <div className="gf-form" key={key}>
          <InlineFormLabel width={10}>{startCase(key)}</InlineFormLabel>
          <input disabled className="gf-form-input width-30" value="configured" />
        </div>
      ))}
      <Button style={{ marginTop: 12 }} variant="secondary" onClick={() => setEnableUpload(true)}>
        Upload another JWT file
      </Button>
    </>
  );
}

function getStyles(theme: GrafanaTheme2) {
  return {
    dropzone: css`
      margin-top: ${theme.spacing(2)};
    `,
  };
}
