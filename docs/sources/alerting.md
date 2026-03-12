---
title: Alerting
menuTitle: Alerting
description: Create alert rules from Google Sheets queries to get notified when your data meets conditions.
keywords:
  - data source
  - google sheets
  - alerting
  - alert rules
  - notifications
labels:
  products:
    - oss
    - enterprise
    - cloud
last_reviewed: 2025-02-11
weight: 500
---

# Alerting

You can create [alert rules](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/alerting/) that use the Google Sheets data source. When your query results meet the conditions you define, Grafana can send notifications to your configured contact points.

For an overview of alerting in Grafana, see [Alerting](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/alerting/).

## Before you begin

- [Configure the Google Sheets data source](configure.md) and ensure **Save & test** shows **Success**.
- Have a dashboard panel that uses a Google Sheets query whose data you want to alert on (e.g. a metric that should stay above or below a threshold).

## Create an alert rule from a panel

1. Open a dashboard that has a panel using the **Google Sheets** data source.
1. Click the panel title and select **Edit** (or create a new panel and add a Google Sheets query).
1. In the query editor, set **Spreadsheet ID** and **Range** (and **Use Time Filter** if needed) so the panel shows the data you want to evaluate.
1. Save the panel and the dashboard if you have made changes.
1. Click the panel title and select **Edit** → open the **Alert** tab (or **Alert** in the panel editor).
1. Click **Create alert rule from this panel**.
1. Define the rule: name, folder, condition (e.g. when a value is above or below a threshold), evaluation group, and contact points. For details, see [Create an alert rule](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/alerting/alerting-rules/create-alert-rules/).
1. Save the rule.

After the rule is created, Grafana will evaluate it on the schedule you configured and send notifications when the condition is met.

## Query and data considerations

- Use the same [query editor](query-editor.md) options as in panels: **Spreadsheet ID**, **Range**, **Use Time Filter**, and **Cache Time**. The alert evaluation runs your query and then applies the condition to the result.
- For threshold-style rules, your sheet should return data that Grafana can treat as numeric or time series (e.g. a column with numbers or a time column plus value column). The exact condition you can set depends on your panel type and how the data is shaped.
- If your sheet data or range changes (e.g. new rows), ensure the alert rule’s query still matches the range you intend (e.g. a fixed range like `Sheet1!A1:E100` or a range that includes new rows).

## Next steps

- [Query editor](query-editor.md) – Spreadsheet ID, Range, and Use Time Filter
- [Configure the data source](configure.md)
- [Alerting](https://grafana.com/docs/grafana/<GRAFANA_VERSION>/alerting/) – Grafana alerting documentation
