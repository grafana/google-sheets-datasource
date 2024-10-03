---
title: Query editor
description: Learn about the query editor for the Google Sheets data source plugin to visualize Google Spreadsheets data in Grafana.
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
weight: 200
---

# Query editor

The Google Sheets data source query editor configures the Google Sheets API query.
Refer to the following sections to understand how to set each configuration option.

{{< figure alt="The Google Sheets data source query editor configured to query a Google Sheet" src="/media/docs/plugins/google-sheets-query-editor-1.png" >}}

## Spreadsheet ID

The **Spreadsheet ID** field controls which spreadsheet to query.

You can:

- Enter a spreadsheet ID.
- Enter a spreadsheet URL.

The query editor automatically extracts the spreadsheet ID from the URL.

- Enter a spreadsheet URL including a range.

  The query editor automatically extracts both spreadsheet ID and range from the URL.
  To copy a range:

  1. Open the spreadsheet.
  1. Select the cells that you want to include.
  1. Right-click one of the cells and choose **Get link to this range**.
     The link is copied to your clipboard.

     {{< figure alt="Google Sheets spreadsheet with selected cells and the right click menu open" src="/media/docs/plugins/google-sheets-query-editor-3.png" caption="Google Sheets spreadsheet with selected cells and the right click menu open" >}}

- Select a spreadsheet from the drop-down menu.

  The drop-down menu is only populated if you are using Google JWT authentication.
  You can only view spreadsheets shared with the service account associated with the token.

  To configure a service account with JWT authentication, refer to [Authenticate with a service account JWT](../setup/authenticate/#authenticate-with-a-service-account-jwt).

Next to the **Spreadsheet ID** field there's an external link icon.
Click that icon to open the spreadsheet in Google Sheets in a new tab.

## Range

The **Range** field controls the range to query.
You use [A1 notation](https://developers.google.com/sheets/api/guides/concepts#a1_notation) to specify the range. If you leave the range field empty, the Google Sheets API returns the whole first sheet in the spreadsheet.

{{< admonition type="tip" >}}
Use a specific range to select relevant data for faster queries and to use less of your Google Sheets API quota.
{{< /admonition >}}

## Cache Time

The **Cache Time** field controls how long to cache the Google Sheets API response.
The cache key is a combination of spreadsheet ID and range.
Changing the spreadsheet ID or range results in a different cache key.

The default cache time is five minutes.
To bypass the cache completely, set **Cache Time** to `0s`.

## Time Filter

The **Time Filter** toggle controls whether to filter rows containing cells with time fields using the dashboard time picker time.
