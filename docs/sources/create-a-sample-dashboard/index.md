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

In this task you're going to create a sample dashboard using a publicly available [demonstration spreadsheet](https://docs.google.com/spreadsheets/d/1TZlZX67Y0s4CvRro_3pCYqRCKuXer81oFp_xcsjPpe8/edit?usp=sharing).

## Before you begin

- Ensure that you have the permissions to create a dashboard and add a data source.
  For more information about permissions, refer to [About users and permissions](https://grafana.com/docs/grafana/latest/administration/roles-and-permissions/).
- Configure the Google Sheets data source plugin.

  You can authenticate with an API key to query public Google Sheets.
  To create an API key, refer to [Authenticate with an API key](../setup/authenticate/#authenticate-with-an-api-key).

  To configure the plugin, refer to [Configure the Google Sheets Plugin](../setup/configure/).

## Create a sample dashboard

To create a sample dashboard:

1. Navigate to the main menu and click on **Dashboards**.
1. Click on the **New** button and select **New Dashboard**.
1. Click on the **Add visualization** button.
1. Select the Google Sheets data source plugin.
1. Browse to the [demonstration spreadsheet](https://docs.google.com/spreadsheets/d/1TZlZX67Y0s4CvRro_3pCYqRCKuXer81oFp_xcsjPpe8/edit?usp=sharing) on your browser.
1. Copy the spreadsheet ID. It should look similar to `1TZlZX67Y0s4CvRro_3pCYqRCKuXer81oFp_xcsjPpe8`.

1. Paste the spreadsheet ID into the query editor

   {{< figure alt="Paste the spreadsheet ID into the query editor" src="/media/docs/plugins/google-sheets-example-1.png" >}}

   Grafana automatically detects this data as time series data and uses the time series panel visualization to display it.

   {{< figure alt="Spreadsheet data visualized in the time series panel visualization" src="/media/docs/plugins/google-sheets-example-2.png" >}}

   You can also use other visualizations like the bar gauge visualization:
   {{< figure alt="Spreadsheet data visualized in the bar gauge panel visualization" src="/media/docs/plugins/google-sheets-example-3.png" >}}

## Grafana Play demonstration

Grafana Play provides a reference dashboard and lets you to modify and create your own custom dashboards.

{{< docs/play title="Google Sheets data source plugin demo" url="https://play.grafana.org/d/ddkar8yanj56oa/visualizing-google-sheets-data" >}}

<!-- cSpell:ignore xcsj -->
