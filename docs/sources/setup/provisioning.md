---
title: Google Sheets data source provisioning
menuTitle: Provisioning
description: Learn how to provision the Google Sheets data source plugin.
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

# Google Sheets data source provisioning

You can define and configure the Google Sheets data source in YAML files with Grafana provisioning.
For more information about provisioning a data source, and for available configuration options, refer to [Provision Grafana](https://grafana.com/docs/grafana/latest/administration/provisioning/#data-sources).

You can provision the data source using any of the following authentication mechanisms:

- [With an API key](#with-an-api-key)
- [With a service account JWT](#with-a-service-account-jwt)
- [With the default GCE service account](#with-the-default-gce-service-account)

## With an API key

To create the API key, refer to [Authenticate with an API key](../configure/#authenticate-with-an-api-key).

**Example**

The following YAML snippet provisions the Google Sheets data source using API key authentication.
Replace _`<API KEY>`_ with your API key, and replace _`<DATA SOURCE NAME>`_ with the name you want to give the data source.

```yaml
apiVersion: 1
datasources:
  - name: '<DATA SOURCE NAME>'
    type: grafana-googlesheets-datasource
    jsonData:
      authenticationType: 'key'
    secureJsonData:
      apiKey: '<API KEY>'
```

## With a service account JWT

To create a service account and its JWT file, refer to [Authenticate with a service account JWT](../configure/#authenticate-with-a-service-account-jwt).

**Example**

The following YAML snippet provisions the Google Sheets data source using service account JWT authentication.
Replace _`<PROJECT ID>`_, _`<CLIENT EMAIL>`_ with your service account details, _`<PRIVATE KEY DATA>`_ with your JWT key data, and replace _`<DATA SOURCE NAME>`_ with the name you want to give the data source.

```yaml
apiVersion: 1
datasources:
  - name: '<DATA SOURCE NAME>'
    type: grafana-googlesheets-datasource
    jsonData:
      authenticationType: 'jwt'
      defaultProject: '<PROJECT ID>'
      clientEmail: '<CLIENT EMAIL>'
      tokenUri: 'https://oauth2.googleapis.com/token'
    secureJsonData:
      privateKey: <PRIVATE KEY DATA>
```

### Private key from local file

The Following example shows the provisioning of google sheets datasource plugin instance using a private key file stored locally.

{{< admonition type="note" >}}
This is not supported in hosted environments such as Grafana Cloud.
{{< /admonition >}}

**Example**

```yaml
apiVersion: 1
datasources:
  - name: '<DATA SOURCE NAME>'
    type: grafana-googlesheets-datasource
    jsonData:
      authenticationType: 'jwt'
      defaultProject: '<PROJECT ID>'
      clientEmail: '<CLIENT EMAIL>'
      privateKeyPath: '/path/to/privateKey'
      tokenUri: 'https://oauth2.googleapis.com/token'
```

## With the default GCE service account

You can use the Google Compute Engine (GCE) default service account to authenticate data source requests if you're running Grafana on GCE.

**Example**

The following YAML snippet provisions the Google Sheets data source using the default GCE service account for authentication.
Replace _`<PROJECT ID>`_ with your GCE project ID and replace _`<DATA SOURCE NAME>`_ with the name you want to give the data source.

```yaml
apiVersion: 1
datasources:
  - name: '<DATA SOURCE NAME>'
    type: grafana-googlesheets-datasource
    jsonData:
      authenticationType: 'gce'
      defaultProject: '<PROJECT ID>'
```
