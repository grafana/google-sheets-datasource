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
  - visualize spreadsheets
labels:
  products:
    - oss
    - enterprise
    - cloud
weight: 103
---

# Provisioning the Google Sheets data source in Grafana

You can define and configure the Google Sheets data source in YAML files with Grafana provisioning. For more information about provisioning a data source, and for available configuration options, refer to [Provision Grafana](https://grafana.com/docs/grafana/latest/administration/provisioning/#data-sources).

You can use any of the 3 types of provisioning:

- API Key
- Google JWT File
- GCE Default Service Account

## Using the API key authentication type

You can use the basic API key authentication type by simply creating a new API Key for the Google Sheets from the Workspace. For more information about the Google Sheets API Key, refer to [Google Sheets API](https://developers.google.com/sheets/api/reference/rest).

**Example**

The following YAML snippet can be used to provision the Google Sheets data source for Grafana if you are using the API key authentication type.

```yaml
apiVersion: 1
datasources:
  - name: GoogleSheetsDatasourceApiKey
    type: google-sheets-datasource
    enabled: true
    jsonData:
      authenticationType: 'key'
    secureJsonData:
      apiKey: ’<YOUR-API-KEY>’
    version: 1
    editable: true
```

## Using the Google JWT service accounts authentication type

You can use the Google JSON Web Tokens (JWT) service accounts authentication type that will allow you to authenticate for server-side applications or backend services that need to access Google APIs on behalf of a user or service account. For more information about the Google Sheets API Key, refer to [Using JWT to authenticate users](https://cloud.google.com/api-gateway/docs/authenticating-users-jwt).

**Example**

The following YAML snippet can be used to provision the Google Sheets data source for Grafana if you are using the JWT (service account) authentication type.

```yaml
apiVersion: 1
datasources:
  - name: GoogleSheetsDatasourceJWT
    type: google-sheets-datasource
    enabled: true
    jsonData:
      authenticationType: 'jwt'
      defaultProject: ’<YOUR_PROJECT_ID>’
      clientEmail: ’<YOUR_CLIENT_EMAIL>’
      tokenUri: 'https://oauth2.googleapis.com/token'
    secureJsonData:
      privateKey: '-----BEGIN PRIVATE KEY-----\nnn-----END PRIVATE KEY-----\n'
    version: 1
    editable: true
```

## Using the GCE authentication type

You can also use the Google Compute Engine (GCE) authentication type if you running applications or services on Google Compute Engine virtual machines as it provides a default service account that is associated with each virtual machin which can also be be used to authenticate and authorize access to Google services and APIs from within the virtual machine. For more information about the Google Sheets API Key, refer to [Authenticate to Compute Engine](https://cloud.google.com/compute/docs/authentication).

**Example**

The following YAML snippet can be used to provision the Google Sheets data source for Grafana if you are using the GCE authentication type.

```yaml
apiVersion: 1
datasources:
  - name: GoogleSheetsDatasourceJWT
    type: google-sheets-datasource
    enabled: true
    jsonData:
      authenticationType: 'gce'
      defaultProject: ’<YOUR_PROJECT_ID>’
    version: 1
    editable: true
```
