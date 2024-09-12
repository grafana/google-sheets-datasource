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

![The Google Sheets data source query editor configured to query a Google Sheet](/media/docs/plugins/google-sheets-query-editor-1.png)

## Spreadsheet ID

Once the **Spreadsheet ID** field is clicked, you have the following options:

- Enter a spreadsheet ID
- Enter a spreadsheet URL. The query editor will then extract the spreadsheet ID from the URL.
- Select a spreadsheet from the dropdown. The dropdown will only be populated if [Google JWT File](./setup/configure.md/) auth is used and as long as spreadsheets are shared with the service account. Read about configuring JWT Auth [here](./setup/configure.md).
  ![Available spreadsheets listed in a dropdown](/media/docs/plugins/google-sheets-query-editor-2.png)
- Enter a link to a certain range. The query editor will then extract both spreadsheet ID and range from the URL. To copy a range, open the Spreadsheet and select the cells that you want to include. Then right click and select `Get link to this range`. The link will be stored in the clipboard.  
  ![Available spreadsheets listed in a dropdown](/media/docs/plugins/google-sheets-query-editor-3.png)

Right next to the Spreadsheet ID input field there's button. If you click on that button, the spreadsheet will be opened in Google Sheets in a separate tab.

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

The **Time Filter** toggle controls whether to filter rows containing cells with time fields using the dashboard timepicker time.
