---
title: Template variables
menuTitle: Template variables
description: Learn how to create and use variables with the Google Sheets data source for Grafana.
keywords:
  - data source
  - google sheets
  - spreadsheets
  - variables
  - template variables
  - dynamic dashboards
  - dashboard variables
labels:
  products:
    - oss
    - enterprise
    - cloud
aliases:
  - /docs/plugins/grafana-googlesheets-datasource/latest/variables/
last_reviewed: 2025-02-11
weight: 300
---

# Template variables

Instead of hard-coding details such as spreadsheet ID or range in your queries, you can use variables. This helps you create more interactive, dynamic, and reusable dashboards. Grafana refers to such variables as template variables. They typically appear as controls (such as drop-down lists) at the top of the dashboard so you can change what data is shown without editing each panel.

For an introduction to templating and template variables, refer to [Variables](https://grafana.com/docs/grafana/latest/dashboards/variables/) and [Add and manage variables](https://grafana.com/docs/grafana/latest/dashboards/variables/add-variable/).

The Google Sheets data source supports **query variables**: variables whose options are loaded from a Google Sheet. You can then use those variables in panel queries (for example, in **Spreadsheet ID** or **Range**) to make dashboards dynamic.

## Before you begin

- [Configure the Google Sheets data source](configure.md) and ensure **Save & test** shows **Success**.
- Have a sheet that contains the data for your variable (e.g. a column of values and optionally a column of labels).

## Supported variable types

| Variable type | Supported |
|---------------|-----------|
| Query | Yes. Options are loaded from a Google Sheet. |
| Custom | Yes. Use Grafana's built-in Custom type; values are not from Google Sheets. |
| Data source | Yes. |

## Query variables

Query variables are populated from a Google Sheet. You choose which spreadsheet and range to use, then which column is the value and which (if any) is the label. You can optionally filter rows by a column and value.

### Create a query variable

To create a query variable:

1. Open a dashboard and go to **Dashboard settings** (gear icon) > **Variables**.
1. Click **Add variable**.
1. Set **Name** and **Type**. For **Type**, select **Query**.
1. In **Data source**, select your Google Sheets data source.
1. In the variable query editor you will see **Spreadsheet ID**, **Range**, **Cache Time**, and **Use Time Filter** (same as in the [query editor](query-editor.md)). Set **Spreadsheet ID** and **Range** to the sheet and range that contain your variable data (e.g. `Sheet1!A1:B10`). Adjust **Cache Time** or **Use Time Filter** if needed.
1. Set **Value Field** to the column that holds the values used in queries.
1. Optionally set **Label Field** to the column that holds the text shown in the drop-down.
1. Optionally use **Optional filtering** to limit rows: set **Filter Field** and **Filter Value**.

### Value and Label fields

- **Value Field**: The column that contains the actual values to be used in queries
- **Label Field**: The column that contains the display text shown in the drop-down (if different from the value)

If you don't specify a label field, the value field will be used for both the value and display text.

### Filtering

You can filter your variable data by specifying additional filter criteria:

- **Filter Field**: The column to use for filtering
- **Filter Value**: The value to match in the filter field

Only rows where the filter field matches the specified filter value will be included in the variable drop-down.

### Example

Consider a Google Sheet with this data:

| Country Code | Country Name | Region |
|-------------|--------------|--------|
| US          | United States | North America |
| GB          | United Kingdom | Europe |
| CA          | Canada | North America |
| FR          | France | Europe |

To create a country variable showing only North American countries:

1. Set **Value Field** to `Country Code`
2. Set **Label Field** to `Country Name`
3. Set **Filter Field** to `Region`
4. Set **Filter Value** to `North America`

This creates a drop-down showing "United States", "Canada" but using the values "US", "CA" in your queries.

## Use variables in queries

After you create a variable, use it in panel queries by referencing its name with a `$` prefix (for example, `$country`). The Google Sheets data source interpolates variables in:

- **Spreadsheet ID** – e.g. use `$spreadsheet` if the variable holds a spreadsheet ID
- **Range** – e.g. use `$range` or build a range like `Sheet1!A1:$column` to make the range depend on the selected value

When the user changes the variable in the dashboard drop-down, queries that reference it are re-run with the new value.

## Next steps

- Use variables in panel queries: [Query editor](query-editor.md)
- Build a full example: [Quick start: create a sample dashboard](_index.md#quick-start-create-a-sample-dashboard)

## Related topics

- [Query editor](query-editor.md)
- [Configure the data source](configure.md)
- [Grafana Variables documentation](https://grafana.com/docs/grafana/latest/dashboards/variables/)
