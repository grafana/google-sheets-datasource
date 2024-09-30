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
weight: 104
---

# Provisioning the Google Sheets data source in Grafana

You can define and configure the Google Sheets data source in YAML files with Grafana provisioning. For more information about provisioning a data source, and for available configuration options, refer to [Provision Grafana](https://grafana.com/docs/grafana/latest/administration/provisioning/#data-sources).

You can provision the data source using any of the following authentication mechanisms:

- Using an API Key
- Using a service account JWT
- Using the default GCE service account

## Using an API key

To create the API key, refer to [Create an API key](../authentication/#using-an-api-key).

**Example**

The following YAML snippet provisions the Google Sheets data source using API key authentication.
Replace _`<API KEY>`_ with your API key, and replace _`<DATA SOURCE NAME>`_ with the name you want to give the data source.

```yaml
apiVersion: 1
datasources:
  - name: <DATA SOURCE NAME>
    type: grafana-googlesheets-datasource
    enabled: true
    jsonData:
      authenticationType: 'key'
    secureJsonData:
      apiKey: <API KEY>
    version: 1
    editable: true
```

## Using a service account JWT

To create a service account and its JWT, refer to [Create a service account](../authentication/#using-a-service-account-jwt).

**Example**

The following YAML snippet provisions the Google Sheets data source using service account JWT authentication.
Replace _`<PROJECT ID>`_, _`<CLIENT EMAIL>`_ with your service account details, _`<PRIVATE KEY DATA>`_ with your JWT key data, and replace _`<DATA SOURCE NAME>`_ with the name you want to give the data source.

```yaml
apiVersion: 1
datasources:
  - name: <DATA SOURCE NAME>
    type: grafana-googlesheets-datasource
    enabled: true
    jsonData:
      authenticationType: 'jwt'
      defaultProject: <PROJECT ID>
      clientEmail: <CLIENT EMAIL>
      tokenUri: 'https://oauth2.googleapis.com/token'
    secureJsonData:
      privateKey: <PRIVATE KEY DATA>
    version: 1
    editable: true
```

## Using the default GCE service account

You can use the Google Compute Engine (GCE) default service account to authenticate data source requests if you're running Grafana on GCE.

**Example**

The following YAML snippet provisions the Google Sheets data source using the default GCE service account for authentication.
Replace _`<PROJECT ID>` with your GCE project ID and replace _`<DATA SOURCE NAME>`_ with the name you want to give the data source.

```yaml
apiVersion: 1
datasources:
  - name: <DATA SOURCE NAME>
    type: grafana-googlesheets-datasource
    enabled: true
    jsonData:
      authenticationType: 'gce'
      defaultProject: '<PROJECT ID>'
    version: 1
    editable: true
```
