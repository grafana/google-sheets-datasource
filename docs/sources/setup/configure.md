---
title: Configure the Google Sheets data source plugin
menuTitle: Configure
description: Learn how to configure the Google Sheets data source plugin
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

# Configure the Google Sheets data source plugin

Before configuring the the Google Sheets data source plugin, you must:
- [Install the plugin](https://grafana.com/docs/grafana/latest/administration/plugin-management/#install-a-plugin)
- [Add a new data source](https://grafana.com/docs/grafana/latest/datasources/#add-a-data-source)

## Authentication

The Google Sheets data source supports the following three ways of authentication:

- [Google JWT File](#authenticate-with-a-service-account-jwt): uses a service account and can access private spreadsheets. Works in all environments where Grafana is running. 
- [API key](#authenticate-with-an-api-key): offers simpler configuration, but requires spreadsheets to be public.
- [GCE Default Service Account](#authenticate-with-the-default-gce-service-account) automatically retrieves default credentials. Requires Grafana to be running on a Google Compute Engine virtual machine.

## Authenticate with a service account JWT

If you want to access private spreadsheets, you must use a service account authentication.
A Google service account is an account that belongs to a project within an account or organization instead of to an individual end user. The application, in this case Grafana, calls Google APIs on behalf of the service account, so users aren't directly involved.

The project that the service account is associated with needs to be granted access to the [Google Sheets API](https://console.cloud.google.com/apis/library/sheets.googleapis.com?q=sheet) and the [Google Drive API](https://console.cloud.google.com/apis/library/drive.googleapis.com?q=drive).

The Google Sheets data source uses the scope `https://www.googleapis.com/auth/spreadsheets.readonly` to get read-only access to spreadsheets. It also uses the scope `https://www.googleapis.com/auth/drive.metadata.readonly` to list all spreadsheets that the service account has access to in Google Drive.

To create a service account, generate a Google JWT file and enable the APIs:

1. Before you can use the Google APIs, you need to enable them in your Google Cloud project.
  1. Open the [Google Sheets API page](https://console.cloud.google.com/apis/library/sheets.googleapis.com?q=sheet) and click enable.
  1. Open the [Google Drive API page](https://console.cloud.google.com/apis/library/drive.googleapis.com?q=drive) and click enable.
1. Open the [Credentials](https://console.developers.google.com/apis/credentials) page in the Google API Console.
1. Click **Create Credentials** then **Service account**.
1. Fill out the service account details form and then click **Create**.
1. Ignore the **Service account permissions** and **Principals with access** sections, just click **Continue**.
1. Click into the details for the service account, navigate to the **Keys** tab, and click **Add Key**. Choose key type **JSON** and click **Create**. A JSON key file will be created and downloaded to your computer.
1. Upload or drag this file into the **JWT Key Details** section of the data source configuration.

### Sharing

By default, the service account doesn't have access to any spreadsheets within the account or organization that it's associated with.
To grant the service account access to files and or folders in Google Drive, you need to share the file or folder with the service account's email address.
The service account's email address is the `client_email` field in the JWT file
To share a file or folder, refer to the [official Google drive documentation](https://support.google.com/drive/answer/2494822?co=GENIE.Platform%3DDesktop&hl=en#share_publicly).

{{< admonition type="caution" >}}
Beware that after you share a file or folder with the service account, all users in Grafana with permissions on the data source are able to see the spreadsheets.
{{< /admonition >}}

## Authenticate with an API key

If a spreadsheet is shared publicly on the internet the request doesn't need to be authorized, but does need to be accompanied by an identifier - which is the API key.

To generate an API key:

1. Before you can use the Google APIs, you need to enable them in your Google Cloud project.
  1. Open the [Google Sheets API page](https://console.cloud.google.com/apis/library/sheets.googleapis.com?q=sheet) and click enable.
1. Open the [Credentials page](https://console.developers.google.com/apis/credentials) in the Google API Console.
1. Click **Create Credentials** and then **API key**.
1. Paste the value in the **API Key** field of the data source configuration.

{{< admonition type="note" >}}
Learn how to share a file or folder publicly in the [official Google Sheet documentation](https://support.google.com/a/users/answer/13309904#sheets_share_link).
{{< /admonition >}}

## Authenticate with the default GCE service account

When Grafana is running on a Google Compute Engine (GCE) virtual machine, Grafana can automatically retrieve default credentials from the metadata server.
As a result, there is no need to generate a private key file for the service account.
You also don't need to upload the file to Grafana.

To authenticate with the default GCE service account:

1. You must create a service account for use by the GCE virtual machine.
   For more information, refer to [Create new service account](https://cloud.google.com/compute/docs/access/create-enable-service-accounts-for-instances#createanewserviceaccount).
1. Verify that the GCE virtual machine instance is running as the service account that you created.
   For more information, refer to [setting up an instance to run as a service account](https://cloud.google.com/compute/docs/access/create-enable-service-accounts-for-instances#using).
1. Allow access to the specified API scope.
1. Enter the project name in the **Default project** field of the data source configuration