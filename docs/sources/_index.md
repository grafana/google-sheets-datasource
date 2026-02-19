---
title: Google Sheets data source
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

# Google Sheets data source

The Google Sheets data source lets you visualize your Google spreadsheets in Grafana. It uses the [Google Sheets API](https://developers.google.com/workspace/sheets/api/guides/concepts) to read data that you can view in dashboard panels or [Explore](https://grafana.com/docs/grafana/latest/explore/).

This video shows how to use the plugin: {{< youtube id="hqeqeQFrtSA">}}

{{< docs/play title="Google Sheets data source demo" url="https://play.grafana.org/d/ddkar8yanj56oa/visualizing-google-sheets-data" >}}

## Requirements

Before you start, ensure you have:

- A [Google account](https://support.google.com/accounts/answer/27441?hl=en) to use Google Sheets
- One of these Grafana editions:
  - [Grafana OSS](https://grafana.com/oss/grafana/)
  - [Grafana Cloud](https://grafana.com/pricing/)
  - Self-managed [Grafana Enterprise](https://grafana.com/docs/grafana/latest/administration/enterprise-licensing/) with an activated license

## Get started

1. [Install the plugin](https://grafana.com/docs/grafana/latest/administration/plugin-management/#install-a-plugin).
1. [Configure the data source](configure.md).
1. Use the [Query editor](query-editor.md) to query spreadsheet data.
1. Add [template variables](template-variables.md) for dynamic dashboards.
1. Try [Create a sample dashboard](create-a-sample-dashboard.md) to get started quickly.

## Additional features

After you configure the data source, you can:

- Add [annotations](annotations.md) to panels
- Use [template variables](template-variables.md) in queries
- Apply [transformations](https://grafana.com/docs/grafana/latest/panels-visualizations/query-transform-data/transform-data/) to query results

## Quota

The Google Sheets API uses per-minute quotas that refill every minute. For limits and usage, refer to the [Google Sheets API usage limits](https://developers.google.com/sheets/api/limits).

## Troubleshooting

If **Save & test** fails, panels show errors, or variables and annotations do not behave as expected, see [Troubleshooting](troubleshooting.md) for common causes and fixes.

## Report issues

Report bugs and request features in the [Google Sheets data source repository](https://github.com/grafana/google-sheets-datasource/issues).
