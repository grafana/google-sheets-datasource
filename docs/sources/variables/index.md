---
title: Variables
description: Learn how to create and use variables with the Google Sheets data source plugin for Grafana.
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
weight: 300
---

# Variables

A variable is a placeholder for a value that you can use in dashboard queries. Variables allow you to create more interactive and dynamic dashboards by replacing hard-coded values with dynamic options. They are displayed as dropdown lists at the top of the dashboard, making it easy to change the data being displayed.

The Google Sheets data source plugin supports two types of variables:

- **Query variables**: Create variables populated with data from your Google Sheets
- **Template variables**: Use variables in your queries to make them dynamic

## Query variables

Query variables allow you to create dropdown lists populated with data from your Google Sheets. These variables can be used in other queries to create dynamic dashboards.

### Create a query variable

To create a query variable:

1. Go to your dashboard settings.
1. Click **Variables**.
1. Click **Add variable**.
1. In the **Query options** section, select your Google Sheets data source.
1. Configure your variable query:
   - **Spreadsheet ID**: Enter the ID of the spreadsheet containing your variable data
   - **Range**: Specify the range containing your data (e.g., `Sheet1!A1:B10`)
   - **Value Field**: Select the column to use as the variable value
   - **Label Field**: Select the column to use as the display text (optional)

### Value and Label fields

- **Value Field**: The column that contains the actual values to be used in queries
- **Label Field**: The column that contains the display text shown in the dropdown (if different from the value)

If you don't specify a label field, the value field will be used for both the value and display text.

### Example

Consider a Google Sheet with the following data:

| Country Code | Country Name |
|-------------|--------------|
| US          | United States |
| GB          | United Kingdom |
| CA          | Canada |

To create a country variable:

1. Set **Value Field** to `Country Code`
2. Set **Label Field** to `Country Name`

This creates a dropdown showing "United States", "United Kingdom", "Canada" but using the values "US", "GB", "CA" in your queries.

## Related topics

- [Query Editor](../query-editor/)
- [Setup](../setup/)
- [Grafana Variables documentation](https://grafana.com/docs/grafana/latest/dashboards/variables/) 