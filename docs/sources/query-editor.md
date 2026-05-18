---
title: Google Sheets query editor
menuTitle: Query editor
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
aliases:
  - /docs/plugins/grafana-googlesheets-datasource/latest/query-editor/
review_date: 2026-05-18
weight: 200
---

# Google Sheets query editor

The Google Sheets data source query editor configures the Google Sheets API query. Use it when building panels in a dashboard or in Explore.

This document walks you through key concepts and a summary of the fields, then how to create a query, then each field in detail, and finally [example use cases](#example-use-cases) for common scenarios.

## Before you begin

- Ensure you have [configured the Google Sheets data source](configure.md) and that **Save & test** shows **Success**.
- Your credentials must have access to the spreadsheets you want to query.

{{< figure alt="The Google Sheets data source query editor configured to query a Google Sheet" src="/media/docs/plugins/google-sheets-query-editor-1.png" caption="The Google Sheets data source query editor configured to query a Google Sheet" >}}

## Key concepts

| Term | Description |
|------|-------------|
| **Spreadsheet ID** | The unique identifier for a Google Sheet. You find it in the sheet URL between `/d/` and `/edit` (for example, `1TZlZX67Y0s4CvRro_3pCYqRCKuXer81oFp_xcsjPpe8`). |
| **Range** | The cells to read, in [A1 notation](https://developers.google.com/sheets/api/guides/concepts#a1_notation) (for example, `Sheet1!A1:E10`) or a named range. Empty means the entire first sheet. |
| **Cache** | Grafana caches each (Spreadsheet ID + Range) response for the duration you set in **Cache Time** to reduce API calls and stay within quota. |

## Query editor fields

| Field | Description |
|-------|-------------|
| **Spreadsheet ID** | Which spreadsheet to query. Can be an ID, a full URL, or a selection from the list (JWT only). |
| **Range** | Which cells to read (A1 notation or named range). Empty = entire first sheet. |
| **Cache Time** | How long to cache the response (default `5m`). Use `0s` to disable cache. |
| **Use Time Filter** | When on, filters rows by the dashboard time range using the first time column in the data. |

The sections below describe each field in detail.

## Create a query

To create a query:

1. Select the **Google Sheets** data source for the panel (or open [Explore](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/explore/) and select it there).
1. In **Spreadsheet ID**, enter the spreadsheet ID or URL, or choose a spreadsheet from the list if you use JWT authentication.
1. Optionally, enter a **Range** in A1 notation (for example, `Sheet1!A2:E`). Leave it empty to use the entire first sheet.
1. Set **Cache Time** if you want something other than the default (5 minutes). Use `0s` to disable caching.
1. Turn **Use Time Filter** on if your data has a time column and you want to filter by the dashboard time range.
1. Run the query. The panel refreshes with the data from your sheet.

## Spreadsheet ID

The **Spreadsheet ID** field controls which spreadsheet to query.

{{< admonition type="tip" >}}
If you've configured a default spreadsheet ID in the data source settings, it will be automatically pre-filled when you create a new query. Refer to [Configure the data source](configure/#default-spreadsheet-id) for more information.
{{< /admonition >}}

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

  To configure a service account with JWT authentication, refer to [Authenticate with a service account JWT](configure/#authenticate-with-a-service-account-jwt).

Next to the **Spreadsheet ID** field there's an external link icon.
Click that icon to open the spreadsheet in Google Sheets in a new tab.

## Range

The **Range** field controls which cells to query. You can use [A1 notation](https://developers.google.com/sheets/api/guides/concepts#a1_notation) or a [named range](https://developers.google.com/sheets/api/guides/concepts#a1_notation). If you leave the range empty, the API returns the whole first sheet in the spreadsheet.

Valid range examples:

| Range | Description |
|-------|-------------|
| *(empty)* | Entire first sheet in the spreadsheet. |
| `Sheet1!A1:E100` | Cells A1 through E100 on the sheet named "Sheet1". |
| `Sheet1!A2:E` | All rows from row 2 onward in columns A–E (open-ended). |
| `Sheet1!A:E` | All rows in columns A–E. |
| `'My Sheet'!A1:D10` | Sheet names that contain spaces must be wrapped in single quotes. |
| `Class Data!A2:E` | Sheet name without spaces — no quotes needed. |
| `MyNamedRange` | A named range defined in the spreadsheet. |

{{< admonition type="tip" >}}
Use a specific range to select relevant data for faster queries and to use less of your Google Sheets API quota.
{{< /admonition >}}

Common range mistakes that cause errors (such as `[sse.readDataError]`):

- Using a colon in the sheet name instead of `!` (for example, `Sheet1:A1:E10` instead of `Sheet1!A1:E10`).
- Misspelling or changing the case of the sheet name (names are case-sensitive).
- Pointing to a range with no data or only a header row.
- Omitting single quotes around sheet names that contain spaces.

## Cache Time

The **Cache Time** field controls how long to cache the Google Sheets API response. The cache key is a combination of spreadsheet ID and range, so changing either results in a different cache key.

Options include `0s`, `5s`, `10s`, `30s`, `1m`, `2m`, `5m`, `10m`, `30m`, `1h`, `2h`, and `5h`. The default is five minutes (`5m`). Set **Cache Time** to `0s` to bypass the cache completely.

## Use Time Filter

The **Use Time Filter** toggle controls whether to apply the dashboard time range to the data. When enabled, the plugin filters rows using the first time field in the data so that only rows within the dashboard time picker range are returned.

The plugin does not use query-language macros (such as `$__timeFilter()`). Use this toggle instead to apply the dashboard time range.

## How data is interpreted

The plugin uses the first row of the range as column headers and infers a type for each column from the cell values:

| Inferred type | Description |
|---------------|-------------|
| **Time** | Columns that Google Sheets formats as date, date-time, or time (or that contain numeric date values) are treated as time. **Use Time Filter** and annotation **time** columns rely on this. |
| **Number** | Numeric cells. |
| **String** | Text or anything that is not consistently time or number. |

- **Mixed types:** If a column has mixed value types (for example, numbers and text), the plugin treats the column as string. A warning is added to the response (for example, “Multiple data types found in column … Using string data type”).
- **Mixed units:** If a column has mixed units (for example, different currencies or formats), the plugin uses the formatted value and may add a warning (for example, “Multiple units found in column … Formatted value will be used”).
- **Parse errors:** Cells that cannot be parsed (for example, invalid date text) produce per-cell warnings. The row may still be returned with a fallback or empty value depending on the column type.

These warnings are included in the query response metadata. If a panel shows unexpected types or empty values, check that the sheet has consistent types in each column and that date columns are formatted as date/time in Google Sheets.

## Example use cases

**Time series or time-based data**

Use a range where the first column (or first column the plugin detects as time) contains dates or timestamps. Enable **Use Time Filter** so the panel only shows rows within the dashboard time picker range. Choose a **Cache Time** that balances freshness with API quota (for example, `5m` for frequently updated data).

**Table or tabular data (no time filter)**

Set **Spreadsheet ID** and **Range** to your data (for example, `Sheet1!A1:D50`). Leave **Use Time Filter** off if there is no time column. Use a longer **Cache Time** if the sheet changes rarely to reduce API usage.

**KPIs or single values**

Use a small range (for example, one row or a few cells, such as `Summary!B2:B5`) and pair it with a Stat or Gauge panel. Set a short **Cache Time** (for example, `30s` or `1m`) if the values update often, or longer if they are static.

**Annotations from a sheet**

Use a Google Sheet as an annotation source to overlay events on panels. Refer to [Annotations](annotations.md) for sheet layout (time, text, tags), steps to add an annotation query, and **Use Time Filter**.

**Same spreadsheet, different panels**

Use one **Spreadsheet ID** across multiple panels and set a different **Range** in each (for example, `Sales!A1:E100`, `Inventory!A1:C50`). To save time, set a default spreadsheet in the [data source configuration](configure.md#default-spreadsheet-id) so new queries are pre-filled.

**Dynamic spreadsheet or range with template variables**

Use [template variables](template-variables.md) in **Spreadsheet ID** or **Range** (for example, `$spreadsheet` or `Sheet1!A1:$region`) so users can switch sheets or ranges from the dashboard. The query runs with the selected variable values.

**Frequently updated vs. static data**

For data that changes often (for example, live status), use a short **Cache Time** (`5s`, `30s`, or `1m`) and be aware of [Google Sheets API quotas](https://developers.google.com/sheets/api/limits). For reference data that rarely changes, use a longer cache (`30m`, `1h`, or `5h`) to reduce API calls.

**Copying a range from Google Sheets**

In Google Sheets, select the cells you need, right-click and choose **Get link to this range**, then paste the link into **Spreadsheet ID**. The editor extracts both the spreadsheet ID and the range from the URL.

## Use SQL expressions with Google Sheets data

You can use Grafana's [SQL expression](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/panels-visualizations/query-transform-data/expression-queries/) transformation to filter or aggregate Google Sheets data. However, there are limitations to be aware of:

- The Google Sheets plugin returns data in **wide format** (one column per field, with the first row as headers). SQL expressions expect this format by default.
- If your sheet layout produces data that Grafana interprets as long format (for example, a single value column with a label column), the SQL expression may fail with errors such as `[sse.readDataError] input data must be a wide series but got type not`.
- To avoid these errors, ensure your sheet has distinct column headers and that each column contains a consistent data type. If a column has mixed types, the plugin treats it as a string, which may cause type mismatches in SQL expressions.
- For complex transformations, consider using Grafana's built-in [transformations](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/panels-visualizations/query-transform-data/transform-data/) (such as **Convert field type**, **Filter by value**, or **Group by**) instead of SQL expressions, as these handle Google Sheets data more reliably.

## Next steps

- [Use template variables](template-variables.md) to make dashboards dynamic.
- [Quick start: create a sample dashboard](_index.md#quick-start-create-a-sample-dashboard) to try the data source with sample data.
- [Configure the data source](configure.md) to change authentication or default spreadsheet.
