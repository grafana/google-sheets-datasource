# Configuring the Google Sheets data source

The Google Sheets data source is using the [Google Sheet API](https://developers.google.com/sheets/api) to access spreadsheets. The data source supports two ways of authenticating against the Google Sheets API. **API Key** auth is used to access public spreadsheets, and **Google JWT File** auth using a service account is used to access private files.

## API Key

If a spreadsheet is shared publicly on the Internet, it can be accessed in the Google Sheets data source using **API Key** auth. When accessing public spreadsheets using the Google Sheets API, the request doesn't need to be authorized, but does need to be accompanied by an identifier, such as an API key.

To generate an API Key, follow the steps in the Google Sheets data source configuration page.

If you want to know how to share a file or folder, read about that in the [official Google drive documentation](https://support.google.com/drive/answer/2494822?co=GENIE.Platform%3DDesktop&hl=en#share_publicly).

## Google JWT File

Whenever access to private spreadsheets is necessary, service account auth using a Google JWT File should be used. A Google service account is an account that belongs to a project within an account or organization instead of to an individual end user. Your application calls Google APIs on behalf of the service account, so users aren't directly involved.

The project that the service account is associated with needs to be granted access to the [Google Sheets API](https://console.cloud.google.com/apis/library/sheets.googleapis.com?q=sheet) and the [Google Drive API](https://console.cloud.google.com/apis/library/drive.googleapis.com?q=drive).

The Google Sheets data source uses the scope `https://www.googleapis.com/auth/spreadsheets.readonly` to get read-only access to spreadsheets. It also uses the scope `https://www.googleapis.com/auth/drive.metadata.readonly` to list all spreadsheets that the service account has access to in Google Drive.

To create a service account, generate a Google JWT file and enable the APIs. For more detailed instructions, refer to the steps documented for the Google Sheets data source in the "Add a data source" page in Grafana.

### Sharing

By default, the service account doesn't have access to any spreadsheets within the account/organization that it is associated with. To grant the service account access to files and/or folders in Google Drive, you need to share the file/folder with the service account's email address. The email is specified in the Google JWT File. If you want to know how to share a file or folder, please refer to the [official Google drive documentation](https://support.google.com/drive/answer/2494822?co=GENIE.Platform%3DDesktop&hl=en#share_publicly).

> **_:warning:_** Beware that once a file/folder is shared with the service account, all users in Grafana will be able to see the spreadsheet/spreadsheets.

## Configure a GCE Default Service Account

When Grafana is running on a Google Compute Engine (GCE) virtual machine, Grafana can automatically retrieve default credentials from the metadata server. As a result, there is no need to generate a private key file for the service account. You also do not need to upload the file to Grafana. The following preconditions must be met before Grafana can retrieve default credentials.

- You must create a Service Account for use by the GCE virtual machine. For more information, refer to [Create new service account](https://cloud.google.com/compute/docs/access/create-enable-service-accounts-for-instances#createanewserviceaccount).
- Verify that the GCE virtual machine instance is running as the service account that you created. For more information, refer to [setting up an instance to run as a service account](https://cloud.google.com/compute/docs/access/create-enable-service-accounts-for-instances#using).
- Allow access to the specified API scope (`"https://www.googleapis.com/auth/spreadsheets.readonly"`).

For more information about creating and enabling service accounts for GCE instances, refer to [enabling service accounts for instances in Google documentation](https://cloud.google.com/compute/docs/access/create-enable-service-accounts-for-instances).
