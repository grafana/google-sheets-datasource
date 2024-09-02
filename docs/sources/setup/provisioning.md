---
title: Provisioning the Google Sheets data source in Grafana
menuTitle: Provisioning
description: Provisioning the Google Sheets source plugin
keywords:
  - data source
  - google sheets
  - spreadsheets
  - xls data
  - xlsx data
  - excel sheets
  - excel data
  - csv data
labels:
  products:
    - oss
    - enterprise
    - cloud
weight: 104
---

# Provisioning the Google Sheets data source in Grafana

You can define and configure the Google Sheets data source in YAML files with Grafana provisioning. For more information about provisioning a data source, and for available configuration options, refer to [Provision Grafana](https://grafana.com/docs/grafana/latest/administration/provisioning/#data-sources).

You can use any of the 3 types of provisioning:

- API key authentication type
- Google JWT file (service account) authentication type
- GCE authentication type

## Using the API key authentication type

You can use the basic API key authentication type by simply creating a new API Key for the Google Sheets from the Workspace. For more information about creating the API Key, refer to [Google Workspace guides](https://developers..google.com/workspace/guides/enable-apis#google-cloud-console)

**Example**

The following YAML snippet can be used to provision the Google Sheets data source for Grafana if you are using the API key authentication type.:

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

## Using the Google JWT file (service account) authentication type

You can also use the Google JWT file authentication type that will allow you to authenticate for server-side applications or backend services that need to access Google APIs on behalf of a user or service account.

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

## Using the GCE authentication type


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