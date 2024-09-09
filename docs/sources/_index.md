---
title: Google Sheets data source plugin for Grafana
menuTitle: Google Sheets data source
description: The Google Sheets data source lets you visualize Google Spreadsheets data in Grafana dashboards.
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
weight: 10
---

# Google Sheets data source plugin for Grafana

The Google Sheets data source plugin for Grafana lets you to visualize your Google Spreadsheets in Grafana. It uses the Google Sheets API to read the data and allow you to define the query inside the editor to view in a Dashboard panel or Explore.

{{< docs/play title="Google Sheets data source plugin demo" url="https://play.grafana.org/d/ddkar8yanj56oa/visualizing-google-sheets-data" >}}

## Requirements

To use the Google Sheets data source plugin, you will need:

- A [Google account](https://support.google.com/accounts/answer/27441?hl=en) to be able to use Google Sheets.
- Any of the following Grafana editions:
  - A [Grafana OSS](https://grafana.com/oss/grafana/) server.
  - A [Grafana Cloud](https://grafana.com/pricing/) stack.
  - An on-premise Grafana Enterprise server with an [activated license](https://grafana.com/docs/grafana/latest/enterprise/license/activate-license/).

## Get started

- To start using the plugin, you need to generate an access token, then install and configure the plugin. To do this, refer to [Setup](./setup).
- To learn how to use the query editor, refer to [Query Editor](./query-editor/).
- To quickly visualize spreadsheet data in Grafana, refer to [Sample dashboards](./sample-dashboard/).

## Get the most out of the plugin

- Add [Annotations](https://grafana.com/docs/grafana/latest/dashboards/annotations/)
- Configure and use [Templates and variables](https://grafana.com/docs/grafana/latest/variables/)
- Add [Transformations](https://grafana.com/docs/grafana/latest/panels/transformations/)

## Quota

Please refer to the [official docs](https://developers.google.com/sheets/api/limits).

## Report issues

Use the [official Google Sheets repository](https://github.com/grafana/google-sheets-datasource/issues) to report issues, bugs, and feature requests.
