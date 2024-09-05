---
title: Creating a sample Dashboard using the Google Sheets data source plugin for Grafana
menuTitle: Creating a sample Dashboard
description: Creating a sample Dashboard using the Google Sheets data source plugin to visualize Google Spreadsheets data in Grafana.
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
weight: 400
---

# Creating a sample Dashboard

We are going to create a sample Dashboard by using this publicly available [demo spreadsheet](https://docs.google.com/spreadsheets/d/1TZlZX67Y0s4CvRro_3pCYqRCKuXer81oFp_xcsjPpe8/edit?usp=sharing) that is suitable for visualization in graphs and in tables.

## Example

1. Copy the sheet-ID from the demo spreadsheet which should look similiar to `1TZlZX67Y0s4CvRro_3pCYqRCKuXer81oFp_xcsjPpe8` in your web browser.

1. Paste in the Query editor
   ![Paste the Spreadsheet ID in the query editor](/media/docs/plugins/google-sheets-example-1.png)

1. After that you should be able to see a similar Time series visualization as Grafana automatically detect this data as Time series and uses the Time Series panel visualziation to display it.
   ![View Spreadsheet data in Time Series panel visualization](/media/docs/plugins/google-sheets-example-2.png)

1. You can also use other visualization panel options for e.g. Bar Gauge:
   ![View Spreadsheet data in Bar Gauge panel visualization ](/media/docs/plugins/google-sheets-example-3.png)

## Play demo

The Play demo dashboards provides a reference the above sample dashboard and also allows you to modify and customized it on the fly.

{{< docs/play title="GitHub data source plugin demo" url="https://play.grafana.org/dashboards/f/bb613d16-7ee5-4cf4-89ac-60dd9405fdd7/demo-github" >}}
