---
title: Create a sample dashboard using the Google Sheets data source plugin for Grafana
menuTitle: Create a sample dashboard
description: Create a sample Dashboard using the Google Sheets data source plugin to visualize Google Spreadsheets data in Grafana.
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

# Create a sample dashboard using the Google Sheets data source plugin for Grafana

We are going to create a sample Dashboard by using this publicly available [demo spreadsheet](https://docs.google.com/spreadsheets/d/1TZlZX67Y0s4CvRro_3pCYqRCKuXer81oFp_xcsjPpe8/edit?usp=sharing) that is suitable for visualization in graphs and in tables.

## Before you begin

- Ensure that you have the proper permissions. For more information about permissions, refer to [About users and permissions](https://grafana.com/docs/grafana/latest/administration/roles-and-permissions/).
- Ensure that the Google Sheets data source plugin is correctly setup on the machine, refer to [Setup](./setup/) if you need instructions.

## To create a dashboard

1. Copy the spreadsheet ID from the demo spreadsheet which should look similar to `1TZlZX67Y0s4CvRro_3pCYqRCKuXer81oFp_xcsjPpe8` in your web browser.

1. Paste in the Query editor
   {{< figure alt="Paste the Spreadsheet ID in the query editor" src="/media/docs/plugins/google-sheets-example-1.png"  >}}

1. After that you should be able to see a similar Time series visualization as Grafana automatically detect this data as Time series and uses the Time Series panel visualziation to display it.
   {{< figure alt="View Spreadsheet data in Time Series panel visualization" src="/media/docs/plugins/google-sheets-example-2.png" >}}

1. You can also use other visualization panel options for e.g. Bar gauge:
   {{< figure alt="View Spreadsheet data in Bar Gauge panel visualization" src="/media/docs/plugins/google-sheets-example-3.png" >}}

## Play demo

The Play demo dashboards provides a reference dashboard and allows you to modify and create your own custom dashboards.

{{< docs/play title="Google Sheets data source plugin demo" url="https://play.grafana.org/d/ddkar8yanj56oa/visualizing-google-sheets-data" >}}
