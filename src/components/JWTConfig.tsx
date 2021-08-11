import { Button } from '@grafana/ui';
import React, { useState } from 'react';
import { DropZone } from './';

export interface Props {
  onChange: (jwt: string) => void;
  isConfigured: boolean;
}

export function JWTConfig({ onChange, isConfigured }: Props) {
  const [enableUpload, setEnableUpload] = useState<boolean>(!isConfigured);
  const [error, setError] = useState<string>();

  return enableUpload ? (
    <>
      <DropZone
        baseStyle={{ marginTop: '24px' }}
        accept="application/json"
        onDrop={(acceptedFiles) => {
          const reader = new FileReader();
          if (acceptedFiles.length === 1) {
            reader.onloadend = (e: any) => {
              onChange(e.target.result);
              setEnableUpload(false);
            };
            reader.readAsText(acceptedFiles[0]);
          } else if (acceptedFiles.length > 1) {
            setError('You can only upload one file');
          }
        }}
      >
        <p style={{ margin: 0, fontSize: 18 }}>Drop the file here, or click to use the file explorer</p>
      </DropZone>

      {error && (
        <pre style={{ margin: '12px 0 0' }} className="gf-form-pre alert alert-error">
          {error}
        </pre>
      )}
    </>
  ) : (
    <>
      {/* {configKeys.map((key) => (
        <div className="gf-form" key={key}>
          <InlineFormLabel width={10}>{startCase(key)}</InlineFormLabel>
          <input disabled className="gf-form-input width-30" value="configured" />
        </div>
      ))} */}
      <Button style={{ marginTop: 12 }} variant="secondary" onClick={() => setEnableUpload(true)}>
        Upload another JWT file
      </Button>
    </>
  );
}
