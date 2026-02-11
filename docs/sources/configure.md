---
title: Configure the Google Sheets data source
menuTitle: Configure
description: Learn how to configure and provision the Google Sheets data source plugin.
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
aliases:
  - /docs/plugins/grafana-googlesheets-datasource/latest/setup/
  - /docs/plugins/grafana-googlesheets-datasource/latest/setup/configure/
  - /docs/plugins/grafana-googlesheets-datasource/latest/setup/provisioning/
  - /docs/plugins/grafana-googlesheets-datasource/latest/setup/authenticate/
  - /docs/plugins/grafana-googlesheets-datasource/latest/setup/install/
weight: 100
---

# Configure the Google Sheets data source

This document explains how to configure and provision the Google Sheets data source plugin.

## Before you begin

Before configuring the Google Sheets data source plugin, you must:

- [Install the plugin](https://grafana.com/docs/grafana/latest/administration/plugin-management/#install-a-plugin)
- [Add a new data source](https://grafana.com/docs/grafana/latest/datasources/#add-a-data-source)

## Add the data source

To add the data source:

1. Click **Connections** in the left-side menu.
1. Click **Add new connection**.
1. Type **Google Sheets** in the search bar.
1. Select **Google Sheets**.
1. Click **Add new data source**.

## Authentication

The Google Sheets data source supports three authentication methods:

- [Google JWT File](#authenticate-with-a-service-account-jwt): uses a service account and can access private spreadsheets. Works in all environments where Grafana is running.
- [API key](#authenticate-with-an-api-key): offers simpler configuration, but requires spreadsheets to be public.
- [GCE Default Service Account](#authenticate-with-the-default-gce-service-account): automatically retrieves default credentials. Requires Grafana to be running on a Google Compute Engine virtual machine.

Depending on your authentication type, you may need to configure appropriate permissions on your spreadsheets. Refer to the [Sharing](#sharing) section for more guidance.

### Authenticate with a service account JWT

If you want to access private spreadsheets, you must use a service account authentication.
A Google service account belongs to a project within an account or organization instead of to an individual end user. The application, in this case Grafana, calls Google APIs on behalf of the service account, so users aren't directly involved.

The project that the service account is associated with needs to be granted access to the [Google Sheets API](https://console.cloud.google.com/apis/library/sheets.googleapis.com?q=sheet) and the [Google Drive API](https://console.cloud.google.com/apis/library/drive.googleapis.com?q=drive).

The Google Sheets data source uses the scope `https://www.googleapis.com/auth/spreadsheets.readonly` to get read-only access to spreadsheets. It also uses the scope `https://www.googleapis.com/auth/drive.metadata.readonly` to list all spreadsheets that the service account has access to in Google Drive.

To create a service account, generate a Google JWT file and enable the APIs:

1. Before you can use the Google APIs, you need to enable them in your Google Cloud project.
   1. Open the [Google Sheets API page](https://console.cloud.google.com/apis/library/sheets.googleapis.com?q=sheet) and click enable.
   1. Open the [Google Drive API page](https://console.cloud.google.com/apis/library/drive.googleapis.com?q=drive) and click enable.
1. Open the [Credentials](https://console.developers.google.com/apis/credentials) page in the Google API Console.
1. Click **Create Credentials** then **Service account**.
1. Fill out the service account details form and then click **Create and continue**.
1. Ignore the **Service account permissions** and **Principals with access** sections, just click **Done**.
1. Click into the details for the service account, navigate to the **Keys** tab, and click **Add Key**. Choose key type **JSON** and click **Create**. A JSON key file will be created and downloaded to your computer.
1. Upload or drag this file into the **JWT Key Details** section of the data source configuration.
1. Grant the service account [access to resources](#granting-access-to-the-service-account-used-with-jwt-authentication) as appropriate.

### Authenticate with an API key

If a spreadsheet is [shared publicly](#sharing) on the internet the request doesn't need to be authorized, but does need to be accompanied by an identifier—which is the API key.

To generate an API key:

1. Before you can use the Google APIs, you need to enable them in your Google Cloud project.
   1. Open the [Google Sheets API page](https://console.cloud.google.com/apis/library/sheets.googleapis.com?q=sheet) and click enable.
1. Open the [Credentials page](https://console.developers.google.com/apis/credentials) in the Google API Console.
1. Click **Create Credentials** and then **API key**.
1. Paste the value in the **API Key** field of the data source configuration.

### Authenticate with the default GCE service account

{{< admonition type="note" >}}
This is **only** compatible when running Grafana on a Google Compute Engine (GCE) virtual machine. It is **not supported** in on-premise deployments, Grafana Cloud or other hosted environments.
{{< /admonition >}}

When Grafana is running on a Google Compute Engine (GCE) virtual machine, Grafana can automatically retrieve default credentials from the metadata server.
As a result, there is no need to generate a private key file for the service account.
You also don't need to upload the file to Grafana.

To authenticate with the default GCE service account:

1. You must create a service account for use by the GCE virtual machine.
   For more information, refer to [Create new service account](https://cloud.google.com/compute/docs/access/create-enable-service-accounts-for-instances#createanewserviceaccount).
1. Verify that the GCE virtual machine instance is running as the service account that you created.
   For more information, refer to [setting up an instance to run as a service account](https://cloud.google.com/compute/docs/access/create-enable-service-accounts-for-instances#using).
1. Allow access to the specified API scope.
1. Enter the project name in the **Default project** field of the data source configuration.

## Configure settings

| Setting | Description |
|---------|-------------|
| **Name** | The name used to refer to the data source in panels and queries. |
| **Default** | Toggle to make this the default data source for new panels. |
| **Default project** | (GCE authentication only) The GCE project ID. |
| **Default Spreadsheet ID** | Optional. When set, this spreadsheet ID is pre-filled in new queries. See [Default Spreadsheet ID](#default-spreadsheet-id). |

## Default Spreadsheet ID

You can optionally configure a **Default Spreadsheet ID** in the data source settings. When set, this spreadsheet ID will be automatically populated in new queries, making it faster to create queries that use the same spreadsheet.

To configure a default spreadsheet ID:

1. Navigate to the data source configuration page.
1. Scroll to the **Default Spreadsheet ID** field.
1. Choose one of these options:
   - **Select from drop-down** (JWT authentication only): If you're using Google JWT File authentication, you can select a spreadsheet from the drop-down menu. The drop-down shows all spreadsheets that the service account has access to.
   - **Enter a spreadsheet ID**: Manually enter the spreadsheet ID from the spreadsheet URL.
   - **Paste a spreadsheet URL**: Paste the full spreadsheet URL, and the ID will be automatically extracted.

When you create a new query, the default spreadsheet ID will be pre-filled in the **Spreadsheet ID** field of the query editor.

{{< admonition type="note" >}}
The default spreadsheet ID is optional. If not set, you'll need to specify the spreadsheet ID for each query manually.
{{< /admonition >}}

## Sharing

Refer to the official guidance from Google on how to share resources:

- [Google Sheets](https://support.google.com/a/users/answer/13309904#sheets_share_link)
- [Google Drive](https://support.google.com/drive/answer/2494822?co=GENIE.Platform%3DDesktop#share_publicly)

### Granting access to the service account used with JWT authentication

By default, the service account doesn't have access to any spreadsheets within the account or organization that it's associated with.
To grant the service account access to files or folders in Google Drive, you need to share the file or folder with the service account's email address.
The service account's email address is the `client_email` field in the JWT file.

{{< admonition type="caution" >}}
Beware that after you share a file or folder with the service account, all users in Grafana with permissions on the data source are able to see the spreadsheets.
{{< /admonition >}}

## Verify the connection

Click **Save & test** to verify the connection.

## Provision the data source

You can define the Google Sheets data source in YAML files as part of Grafana's provisioning system.
For more information, refer to [Provision Grafana](https://grafana.com/docs/grafana/latest/administration/provisioning/#data-sources).

You can provision the data source using any of these authentication mechanisms:

- [With an API key](#with-an-api-key)
- [With a service account JWT](#with-a-service-account-jwt)
- [With the default GCE service account](#with-the-default-gce-service-account)

### With an API key

To create the API key, refer to [Authenticate with an API key](#authenticate-with-an-api-key).

Example:

```yaml
apiVersion: 1
datasources:
  - name: <DATA_SOURCE_NAME>
    type: grafana-googlesheets-datasource
    jsonData:
      authenticationType: 'key'
      defaultSheetID: <SPREADSHEET_ID> # Optional: default spreadsheet ID for new queries
    secureJsonData:
      apiKey: <API_KEY>
```

Replace `<API_KEY>`, `<DATA_SOURCE_NAME>`, and optionally `<SPREADSHEET_ID>` with your values.

### With a service account JWT

To create a service account and its JWT file, refer to [Authenticate with a service account JWT](#authenticate-with-a-service-account-jwt).

Example:

```yaml
apiVersion: 1
datasources:
  - name: <DATA_SOURCE_NAME>
    type: grafana-googlesheets-datasource
    jsonData:
      authenticationType: 'jwt'
      defaultProject: <PROJECT_ID>
      clientEmail: <CLIENT_EMAIL>
      tokenUri: 'https://oauth2.googleapis.com/token'
      defaultSheetID: <SPREADSHEET_ID> # Optional: default spreadsheet ID for new queries
    secureJsonData:
      privateKey: <PRIVATE_KEY_DATA>
```

Replace `<PROJECT_ID>`, `<CLIENT_EMAIL>`, `<PRIVATE_KEY_DATA>`, `<DATA_SOURCE_NAME>`, and optionally `<SPREADSHEET_ID>` with your values.

#### Private key from local file

The following example shows provisioning the Google Sheets data source using a private key file stored locally.

{{< admonition type="note" >}}
This is not supported in hosted environments such as Grafana Cloud.
{{< /admonition >}}

```yaml
apiVersion: 1
datasources:
  - name: <DATA_SOURCE_NAME>
    type: grafana-googlesheets-datasource
    jsonData:
      authenticationType: 'jwt'
      defaultProject: <PROJECT_ID>
      clientEmail: <CLIENT_EMAIL>
      privateKeyPath: '/path/to/privateKey'
      tokenUri: 'https://oauth2.googleapis.com/token'
      defaultSheetID: <SPREADSHEET_ID> # Optional: default spreadsheet ID for new queries
```

### With the default GCE service account

You can use the Google Compute Engine (GCE) default service account to authenticate data source requests if you're running Grafana on GCE.

Example:

```yaml
apiVersion: 1
datasources:
  - name: <DATA_SOURCE_NAME>
    type: grafana-googlesheets-datasource
    jsonData:
      authenticationType: 'gce'
      defaultProject: <PROJECT_ID>
      defaultSheetID: <SPREADSHEET_ID> # Optional: default spreadsheet ID for new queries
```

Replace `<PROJECT_ID>`, `<DATA_SOURCE_NAME>`, and optionally `<SPREADSHEET_ID>` with your values.
