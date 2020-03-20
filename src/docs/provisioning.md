# Provisioning

It's possible configure the Google Sheets data source using config files with Grafana's provisioning system. You can read more about how it works and all the settings you can set for data sources on the [provisioning docs page](https://grafana.com/docs/grafana/latest/administration/provisioning/#datasources).

Here is a provisioning example using API key authentication type.

```yaml
apiVersion: 1
datasources:
  - name: GoogleSheetsDatasourceApiKey
    type: google-sheets-datasource
    enabled: true
    jsonData:
      authType: 'key'
    secureJsonData:
      apiKey: 'your-api-key'
    version: 1
    editable: true
```

Here is a provisioning example using a Google JWT file (service account) authentication type.

```yaml
apiVersion: 1
datasources:
  - name: GoogleSheetsDatasourceJWT
    type: google-sheets-datasource
    enabled: true
    jsonData:
      authType: 'jwt'
    secureJsonData:
      jwt: '{"type":"service_account","project_id":"proj-id","private_key_id":"c4ac...","private_key":"-----BEGIN PRIVATE KEY-----\nnn-----END PRIVATE KEY-----\n","client_email":"nnn@proj.iam.gserviceaccount.com","client_id":"client-id","auth_uri":"https://accounts.google.com/o/oauth2/auth","token_uri":"https://oauth2.googleapis.com/token","auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs","client_x509_cert_url":"cert-url"}'
    version: 1
    editable: true
```
