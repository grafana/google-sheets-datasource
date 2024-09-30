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

- [with an API Key](#authenticate-with-an-api-key)
- [with a service account JWT](#authenticate-with-a-service-account-jwt)
- [with the default GCE service account](#authenticate-with-the-default-gce-service-account)

## Authenticate with an API Key

If a spreadsheet is shared publicly on the Internet, it can be accessed in the Google Sheets data source using **API Key** auth. When accessing public spreadsheets using the Google Sheets API, the request doesn't need to be authorized, but does need to be accompanied by an identifier, such as an API key.

To generate an API Key, refer to the following steps:

1. Open the [Credentials page](https://console.developers.google.com/apis/credentials) in the Google API Console.
1. Click Create Credentials and then click API key.
1. Before using Google APIs, you need to turn them on in a Google Cloud project. [Enable the API](https://console.cloud.google.com/apis/library/sheets.googleapis.com)
1. Copy and Paste the API Key to an editor which you will use it later when configuring the plugin.

{{< admonition type="note" >}}
If you want to know how to share a file or folder, read about that in the [official Google drive documentation](https://support.google.com/drive/answer/2494822?co=GENIE.Platform%3DDesktop&hl=en#share_publicly).
{{< /admonition >}}

## Authenticate with a service account JWT

Whenever access to private spreadsheets is necessary, service account auth using a Google JWT File should be used. A Google service account is an account that belongs to a project within an account or organization instead of to an individual end user. Your application calls Google APIs on behalf of the service account, so users aren't directly involved.

The project that the service account is associated with needs to be granted access to the [Google Sheets API](https://console.cloud.google.com/apis/library/sheets.googleapis.com?q=sheet) and the [Google Drive API](https://console.cloud.google.com/apis/library/drive.googleapis.com?q=drive).

The Google Sheets data source uses the scope `https://www.googleapis.com/auth/spreadsheets.readonly` to get read-only access to spreadsheets. It also uses the scope `https://www.googleapis.com/auth/drive.metadata.readonly` to list all spreadsheets that the service account has access to in Google Drive.

To create a service account, generate a Google JWT file and enable the APIs, refer to the following steps:

1. Open the [Credentials](https://console.developers.google.com/apis/credentials) page in the Google API Console.
1. Click **Create Credentials** then click Service account.
1. On the Create service account page, enter the Service account details.
1. On the `Create service account` page, fill in the `Service account details` and then click `Create`.
1. On the `Service account permissions` page, don’t add a role to the service account. Just click `Continue`.
1. In the next step, click `Create Key`. Choose key type `JSON` and click `Create`. A JSON key file will be created and downloaded to your computer
1. Open the [Google Sheets](https://console.cloud.google.com/apis/library/sheets.googleapis.com?q=sheet) in API Library and enable access for your account
1. Open the [Google Drive](https://console.cloud.google.com/apis/library/drive.googleapis.com?q=drive) in API Library and enable access for your account. Access to the Google Drive API is used to list all spreadsheets that you have access to.
1. Share any private files/folders you want to access with the service account's email address. The email is specified as `client_email` in the Google JWT File.
1. Save this file on your machine as you will use it later when configuring the plugin.

### Sharing

By default, the service account doesn't have access to any spreadsheets within the account/organization that it is associated with. To grant the service account access to files and/or folders in Google Drive, you need to share the file/folder with the service account's email address. The email is specified in the Google JWT File. If you want to know how to share a file or folder, please refer to the [official Google drive documentation](https://support.google.com/drive/answer/2494822?co=GENIE.Platform%3DDesktop&hl=en#share_publicly).

{{< admonition type="caution" >}}
Beware that once a file/folder is shared with the service account, all users in Grafana will be able to see the spreadsheet/spreadsheets.
{{< /admonition >}}

## Authenticate with the default GCE service account

When Grafana is running on a Google Compute Engine (GCE) virtual machine, Grafana can automatically retrieve default credentials from the metadata server. As a result, there is no need to generate a private key file for the service account. You also do not need to upload the file to Grafana.

For creating and enabling service accounts for GCE instances, refer to the following steps:

1. Select the **GCE Default Service Account** in the Authentication type.
1. You must create a Service Account for use by the GCE virtual machine. For more information, refer to [Create new service account](https://cloud.google.com/compute/docs/access/create-enable-service-accounts-for-instances#createanewserviceaccount).
1. Verify that the GCE virtual machine instance is running as the service account that you created. For more information, refer to [setting up an instance to run as a service account](https://cloud.google.com/compute/docs/access/create-enable-service-accounts-for-instances#using).
1. Allow access to the specified API scope.
1. Copy and Paste the project name to an editor which you will use it later when configuring the plugin.

{{< admonition type="note" >}}
For more information about creating and enabling service accounts for GCE instances, refer to enabling service accounts for instances in Google documentation.
{{< /admonition >}}
