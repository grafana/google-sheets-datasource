import { DataSourcePluginOptionsEditorProps, onUpdateDatasourceSecureJsonDataOption } from '@grafana/data';
import { AuthConfig, DataSourceOptions } from '@grafana/google-sdk';
import { Field, LegacyForms, SecretInput } from '@grafana/ui';
import React from 'react';
import { GoogleSheetsAuth, googleSheetsAuthTypes, GoogleSheetsSecureJSONData } from '../types';

const { SecretFormField } = LegacyForms;
export type Props = DataSourcePluginOptionsEditorProps<DataSourceOptions, GoogleSheetsSecureJSONData>;

export function ConfigEditor(props: Props) {
  const apiKeyProps = {
    isConfigured: Boolean(props.options.secureJsonFields.apiKey),
    value: props.options.secureJsonData?.apiKey || '',
    placeholder: 'Enter API key',
    onReset: () =>
      props.onOptionsChange({
        ...props.options,
        secureJsonFields: { ...props.options.secureJsonFields, apiKey: false },
        secureJsonData: { apiKey: '' },
        jsonData: props.options.jsonData,
      }),
    onChange: onUpdateDatasourceSecureJsonDataOption(props, 'apiKey'),
  };

  return (
    <>
      <AuthConfig authOptions={googleSheetsAuthTypes} {...props} />

      {props.options.jsonData.authenticationType === GoogleSheetsAuth.API && (
        <>
          {/* Backward compatibility check. SecretInput was added in 8.5 */}
          {!!SecretInput ? (
            <Field label="API key">
              <SecretInput {...apiKeyProps} width={60} />
            </Field>
          ) : (
            <SecretFormField {...apiKeyProps} label="API key" labelWidth={10} inputWidth={20} />
          )}
        </>
      )}
    </>
    // <div className="gf-form-group">
    //   <div className="grafana-info-box" style={{ marginTop: 24 }}>
    //     {props.options.jsonData.authType === GoogleSheetsAuth.JWT ? (
    //       <>
    //         <h4>Generate a JWT file</h4>
    //         <ol style={{ listStylePosition: 'inside' }}>
    //           <li>
    //             Open the{' '}
    //             <a
    //               href="https://console.developers.google.com/apis/credentials"
    //               target="_blank"
    //               rel="noreferrer noopener"
    //             >
    //               Credentials
    //             </a>{' '}
    //             page in the Google API Console.
    //           </li>
    //           <li>
    //             Click <strong>Create Credentials</strong> then click <strong>Service account</strong>.
    //           </li>
    //           <li>On the Create service account page, enter the Service account details.</li>
    //           <li>
    //             On the <code>Create service account</code> page, fill in the <code>Service account details</code> and
    //             then click <code>Create</code>
    //           </li>
    //           <li>
    //             On the <code>Service account permissions</code> page, don&rsquo;t add a role to the service account.
    //             Just click <code>Continue</code>
    //           </li>
    //           <li>
    //             In the next step, click <code>Create Key</code>. Choose key type <code>JSON</code> and click{' '}
    //             <code>Create</code>. A JSON key file will be created and downloaded to your computer
    //           </li>
    //           <li>
    //             Open the{' '}
    //             <a
    //               href="https://console.cloud.google.com/apis/library/sheets.googleapis.com?q=sheet"
    //               target="_blank"
    //               rel="noreferrer noopener"
    //             >
    //               Google Sheets
    //             </a>{' '}
    //             in API Library and enable access for your account
    //           </li>
    //           <li>
    //             Open the{' '}
    //             <a
    //               href="https://console.cloud.google.com/apis/library/drive.googleapis.com?q=drive"
    //               target="_blank"
    //               rel="noreferrer noopener"
    //             >
    //               Google Drive
    //             </a>{' '}
    //             in API Library and enable access for your account. Access to the Google Drive API is used to list all
    //             spreadsheets that you have access to.
    //           </li>
    //           <li>
    //             Drag the file to the dotted zone above. Then click <code>Save & Test</code>. The file contents will be
    //             encrypted and saved in the Grafana database.
    //           </li>
    //         </ol>
    //       </>
    //     ) : (
    //       <>
    //         <h4>Generate an API key</h4>
    //         <ol style={{ listStylePosition: 'inside' }}>
    //           <li>
    //             Open the{' '}
    //             <a
    //               href="https://console.developers.google.com/apis/credentials"
    //               target="_blank"
    //               rel="noreferrer noopener"
    //             >
    //               Credentials page
    //             </a>{' '}
    //             in the Google API Console.
    //           </li>
    //           <li>
    //             Click <strong>Create Credentials</strong> and then click <strong>API key</strong>.
    //           </li>
    //           <li>
    //             Copy the key and paste it in the API Key field above. The file contents are encrypted and saved in the
    //             Grafana database.
    //           </li>
    //         </ol>
    //       </>
    //     )}
    //   </div>
    // </div>
  );
}
