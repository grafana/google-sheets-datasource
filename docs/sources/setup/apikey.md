---
title: Create a Google Sheets API Key
menuTitle: Create a Google Sheets API Key
description: Create a Google Sheets API Key
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
weight: 102
---

# Create a Google Sheets API Key

The Google Sheets data source plugin uses the Google Sheet API to access the spreadsheets. It supports the following three ways of authentication:

- API Key
- Google JWT File
- GCE Default Service Account

## API Key

If a spreadsheet is shared publicly on the Internet, it can be accessed in the Google Sheets data source using **API Key** auth. When accessing public spreadsheets using the Google Sheets API, the request doesn't need to be authorized, but does need to be accompanied by an identifier, such as an API key.

To generate an API Key, follow the steps:

1. Open the [Credentials page](https://console.developers.google.com/apis/credentials) in the Google API Console.
1. Click Create Credentials and then click API key.
1. Before using Google APIs, you need to turn them on in a Google Cloud project. [Enable the API](https://console.cloud.google.com/apis/library/sheets.googleapis.com)
1. Copy the key and paste it in the API Key field above. The file contents are encrypted and saved in the Grafana database.

{{< admonition type="note" >}}
If you want to know how to share a file or folder, read about that in the [official Google drive documentation](https://support.google.com/drive/answer/2494822?co=GENIE.Platform%3DDesktop&hl=en#share_publicly).
{{< /admonition >}}
