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
aliases:
  - /docs/plugins/grafana-googlesheets-datasource/latest/create-a-sample-dashboard/
weight: 10
---

# Google Sheets data source

The Google Sheets data source lets you visualize your Google spreadsheets in Grafana. It uses the [Google Sheets API](https://developers.google.com/workspace/sheets/api/guides/concepts) to read data that you can view in dashboard panels or [Explore](https://grafana.com/docs/grafana/latest/explore/).

This video shows how to use the plugin: {{< youtube id="hqeqeQFrtSA">}}

{{< docs/play title="Google Sheets data source demo" url="https://play.grafana.org/d/ddkar8yanj56oa/visualizing-google-sheets-data" >}}

## Supported features

| Feature | Supported | Notes |
|---------|-----------|-------|
| Query data | Yes | Query spreadsheets by ID and range; cache by (spreadsheet + range). |
| Annotations | Yes | Use a sheet as an annotation source (time, text, tags). |
| Template variables | Yes | Query variables from a sheet; use variables in Spreadsheet ID or Range. |
| Explore | Yes | Ad-hoc queries without building a dashboard. |
| Alerting | No | Use a data source that supports alerting for alert rules. |

## Requirements

Before you start, ensure you have:

- **Grafana 11.6.0 or newer** (plugin requirement)
- A [Google account](https://support.google.com/accounts/answer/27441?hl=en) to use Google Sheets
- One of these Grafana editions: [Grafana OSS](https://grafana.com/oss/grafana/), [Grafana Cloud](https://grafana.com/pricing/), or self-managed [Grafana Enterprise](https://grafana.com/docs/grafana/latest/administration/enterprise-licensing/) with an activated license

## Get started

The following documents help you get started:

- [Configure the data source](configure.md) – Set up authentication and connect to Google Sheets.
- [Query editor](query-editor.md) – Query spreadsheet data and build panels.
- [Template variables](template-variables.md) – Create dynamic dashboards with variables.
- [Troubleshooting](troubleshooting.md) – Solve common configuration and query errors.

[Install the plugin](https://grafana.com/docs/grafana/latest/administration/plugin-management/#install-a-plugin) if you haven’t already. Try the [Quick start](#quick-start-create-a-sample-dashboard) below to build a sample dashboard in a few steps.

## Key capabilities

- With JWT (service account) authentication: choose spreadsheets from a drop-down list of sheets shared with the service account
- Paste a Google Sheets URL (including **Get link to this range**) and the editor extracts spreadsheet ID and range automatically
- Use [template variables](template-variables.md) in **Spreadsheet ID** or **Range** for dynamic dashboards
- Cache responses by (spreadsheet ID + range) with configurable **Cache Time**; apply the dashboard time range with **Use Time Filter**

## Quick start: create a sample dashboard

You can try the data source using a [public demonstration spreadsheet](https://docs.google.com/spreadsheets/d/1TZlZX67Y0s4CvRro_3pCYqRCKuXer81oFp_xcsjPpe8/edit?usp=sharing).

1. [Configure the data source](configure.md) (API key is enough for this public sheet).
1. Go to **Dashboards** → **New** → **New Dashboard** → **Add visualization**.
1. Select the **Google Sheets** data source.
1. In the query editor, paste the spreadsheet ID: `1TZlZX67Y0s4CvRro_3pCYqRCKuXer81oFp_xcsjPpe8`. You can leave **Range** empty to use the first sheet, or set a range (e.g. `Sheet1!A1:E100`).
1. Run the query. Grafana will detect time series data and suggest a time series panel; you can switch to other visualizations (e.g. bar gauge, table) from the panel.

The embedded Grafana Play dashboard at the top of this page shows a full example you can open and edit.

## Additional features

After you configure the data source, you can:

- Add [annotations](annotations.md) to overlay events on panels
- Use [template variables](template-variables.md) in queries
- Use [Explore](https://grafana.com/docs/grafana/latest/explore/) for ad-hoc queries without building a dashboard
- Apply [transformations](https://grafana.com/docs/grafana/latest/panels-visualizations/query-transform-data/transform-data/) to query results

## Quota

The Google Sheets API uses per-minute quotas that refill every minute. For limits and usage, refer to the [Google Sheets API usage limits](https://developers.google.com/sheets/api/limits).

## Known limitations

- **Read-only:** The data source only reads from spreadsheets; it does not write or edit data.
- **Alerting:** Alert rules cannot use this data source. Use a data source that supports alerting (e.g. Prometheus, MySQL) for alerts.
- **GCE authentication:** [GCE Default Service Account](configure.md#authenticate-with-the-default-gce-service-account) is only supported when Grafana runs on a Google Compute Engine VM. It is not supported in Grafana Cloud or other hosted environments.
- **API key:** With API key authentication, spreadsheets must be publicly viewable (e.g. “Anyone with the link”). The **Select Spreadsheet ID** drop-down is only available when using JWT (service account) authentication.
- **Provisioning:** Provisioning the data source using a local private key file (`privateKeyPath`) is not supported in hosted environments such as Grafana Cloud.

## Related resources

- [Google Sheets API documentation](https://developers.google.com/sheets/api/guides/concepts)
- [Grafana community forums](https://community.grafana.com/)
- [Report bugs and request features](https://github.com/grafana/google-sheets-datasource/issues) (GitHub)

## Plugin updates

Ensure your plugin version is up to date so you have access to all current features and improvements. Navigate to **Administration** > **Plugins and data** > **Plugins** to check for updates.

{{< admonition type="note" >}}
Plugins are automatically updated in Grafana Cloud.
{{< /admonition >}}
