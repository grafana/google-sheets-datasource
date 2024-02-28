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
      authenticationType: 'key'
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
      authenticationType: 'jwt'
      defaultProject: 'your-project-id'
      clientEmail: 'your-client-email'
      tokenUri: 'https://oauth2.googleapis.com/token'
    secureJsonData:
      privateKey: '-----BEGIN PRIVATE KEY-----\nnn-----END PRIVATE KEY-----\n'
    version: 1
    editable: true
```

Here is a provisioning example using a GCE authentication type.

```yaml
apiVersion: 1
datasources:
  - name: GoogleSheetsDatasourceJWT
    type: google-sheets-datasource
    enabled: true
    jsonData:
      authenticationType: 'gce'
      defaultProject: 'your-project-id'
    version: 1
    editable: true
```
