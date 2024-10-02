---
title: Google Sheets data source plugin for Grafana
menuTitle: Google Sheets data source
description: The Google Sheets data source lets you visualize Google spreadsheet data in Grafana dashboards.
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

The Google Sheets data source plugin for Grafana lets you to visualize your Google spreadsheets in Grafana.
It uses the Google Sheets API to read the data which you can view in dashboard panels or [Explore](https://grafana.com/docs/grafana/latest/explore/).

{{< docs/play title="Google Sheets data source plugin demo" url="https://play.grafana.org/d/ddkar8yanj56oa/visualizing-google-sheets-data" >}}

## Requirements

To use the Google Sheets data source plugin, you need:

- A [Google account](https://support.google.com/accounts/answer/27441?hl=en) to be able to use Google Sheets.
- Any of the following Grafana editions:
  - A [Grafana OSS](https://grafana.com/oss/grafana/) server.
  - A [Grafana Cloud](https://grafana.com/pricing/) stack.
  - A self-managed Grafana Enterprise server with an [activated license](/docs/grafana/latest/administration/enterprise-licensing/).

## Get started

- To start using the plugin, you need to generate an access token, then install and configure the plugin. To do this, refer to [Setup](./setup/).
- To learn how to use the query editor, refer to [Query Editor](./query-editor/).
- To quickly visualize spreadsheet data in Grafana, refer to [Create a sample dashboard](./create-a-sample-dashboard/).

## Get the most out of the plugin

- Add [annotations](/docs/grafana/latest/dashboards/build-dashboards/annotate-visualizations/)
- Configure and use [variables](https://grafana.com/docs/grafana/latest/dashboards/variables/)
- Apply [transformations](/docs/grafana/latest/panels-visualizations/query-transform-data/transform-data/)

## Quota

The Google Sheets API has per-minute quotas, and they're refilled every minute.
To understand the API quotas, refer to the [Google Sheets API Usage limits documentation](https://developers.google.com/sheets/api/limits).

## Report issues

Report issues, bugs, and feature requests in the [official Google Sheets data source repository](https://github.com/grafana/google-sheets-datasource/issues).
