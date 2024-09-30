---
title: Authenticate
menuTitle: Authenticate
description: Authenticate the Google Sheets data source plugin
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
weight: 102
---

# Authenticate the Google Sheets data source plugin

The Google Sheets data source plugin uses the Google Sheet API to access the spreadsheets.
It supports the following three ways of authentication:

- [with an API key](#authenticate-with-an-api-key)
- [with a service account JWT](#authenticate-with-a-service-account-jwt)
- [with the default GCE service account](#authenticate-with-the-default-gce-service-account)

## Authenticate with an API key

If a spreadsheet is shared publicly on the internet, you can access in the Google Sheets data source with an API key.
The request doesn't need to be authorized, but does need to be accompanied by an identifier which is the API key.

To generate an API key:

1. Before you can use the Google APIs, you need to turn them on in a Google Cloud project.
   To enable the Google Sheets API, refer to the [Google Sheets API page](https://console.cloud.google.com/apis/library/sheets.googleapis.com).
1. Open the [Credentials page](https://console.developers.google.com/apis/credentials) in the Google API Console.
1. Click **Create Credentials** and then **API key**.
1. Copy the API key to an editor as you will use it later to configure the plugin.

{{< admonition type="note" >}}
If you want to know how to share a file or folder, read about that in the [official Google Drive documentation](https://support.google.com/drive/answer/2494822?co=GENIE.Platform%3DDesktop&hl=en#share_publicly).
{{< /admonition >}}

## Authenticate with a service account JWT

If you want to access private spreadsheets, you must use a service account authentication.
A Google service account is an account that belongs to a project within an account or organization instead of to an individual end user. Your application calls Google APIs on behalf of the service account, so users aren't directly involved.

The project that the service account is associated with needs to be granted access to the [Google Sheets API](https://console.cloud.google.com/apis/library/sheets.googleapis.com?q=sheet) and the [Google Drive API](https://console.cloud.google.com/apis/library/drive.googleapis.com?q=drive).

The Google Sheets data source uses the scope `https://www.googleapis.com/auth/spreadsheets.readonly` to get read-only access to spreadsheets. It also uses the scope `https://www.googleapis.com/auth/drive.metadata.readonly` to list all spreadsheets that the service account has access to in Google Drive.

To create a service account, generate a Google JWT file and enable the APIs:

1. Before you can use the Google APIs, you need to turn them on in a Google Cloud project.
   To enable the Google Sheets API, refer to the [Google Sheets API page](https://console.cloud.google.com/apis/library/sheets.googleapis.com).
1. Open the [Credentials](https://console.developers.google.com/apis/credentials) page in the Google API Console.
1. Click **Create Credentials** then **Service account**.
1. Fill out the service account details form and then click **Create**.
1. On the **Service account permissions** page, don't add a role to the service account, just click **Continue**.
1. In the next step, click **Create Key**.

   Choose key type `JSON` and click **Create**.

   It creates a JSON key file that's downloaded to your computer

1. Open the [Google Sheets API page](https://console.cloud.google.com/apis/library/sheets.googleapis.com?q=sheet) and enable access for your account.
1. Open the [Google Drive API page](https://console.cloud.google.com/apis/library/drive.googleapis.com?q=drive) and enable access for your account.
   You need access to the Google Drive API to list all spreadsheets that you have access to.
1. Share any private files and folders you want to access with the service account's email address.
   The service account's email address is the `client_email` field in the JWT file.
1. Keep the JWT file on your machine as you will use it later to configure the plugin.

### Sharing

By default, the service account doesn't have access to any spreadsheets within the account or organization that it's associated with.
To grant the service account access to files and or folders in Google Drive, you need to share the file or folder with the service account's email address.
The service account's email address is the `client_email` field in the JWT file
To share a file or folder, refer to the [official Google drive documentation](https://support.google.com/drive/answer/2494822?co=GENIE.Platform%3DDesktop&hl=en#share_publicly).

{{< admonition type="caution" >}}
Beware that after you share a file or folder with the service account, all users in Grafana with permissions on the data source are able to see the spreadsheets.
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
1. Copy the project name as you will use it later to configure the plugin.
