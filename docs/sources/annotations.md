---
title: Annotations
menuTitle: Annotations
description: Use a Google Sheet as an annotation source to overlay events on dashboard panels.
keywords:
  - data source
  - google sheets
  - annotations
  - events
  - deploy markers
labels:
  products:
    - oss
    - enterprise
    - cloud
last_reviewed: 2025-02-11
weight: 400
---

# Annotations

Annotations overlay event information on top of graphs. You can use a Google Sheet as the source for annotations to mark deployments, incidents, releases, or other events on your dashboards.

For general information about annotations, refer to [Annotate visualizations](https://grafana.com/docs/grafana/latest/dashboards/build-dashboards/annotate-visualizations/).

## Before you begin

- [Configure the Google Sheets data source](configure.md) and ensure **Save & test** shows **Success**.
- Have a Google Sheet that contains your event data, with at least one column that holds a date or time (so the plugin can detect it as a time field). Optionally include columns for the annotation text and tags.

## Create an annotation query

To add a Google Sheets annotation to a dashboard:

1. Open the dashboard and click **Dashboard settings** (gear icon).
1. Select **Annotations** in the left menu.
1. Click **Add annotation query** (or **+ New query** if you already have one).
1. Enter a **Name** for the annotation (e.g. "Deployments", "Incidents").
1. Select your **Google Sheets** data source.
1. In the query editor, set **Spreadsheet ID** and **Range** to the sheet and range that contain your annotation data (e.g. `Sheet1!A1:C100`). The first row should be the header row with column names.
1. Enable **Use Time Filter** so only rows within the dashboard time range are included.
1. Optionally adjust **Cache Time**.
1. Click **Apply** to save, then **Save dashboard**.

## Query requirements

The plugin uses the first row of the range as column headers. Those headers become the field names Grafana expects for annotations. Use the following column names in your sheet:

| Column   | Required | Description |
|----------|----------|-------------|
| `time`   | Yes      | The timestamp for the annotation. Use a column with date or datetime values. The plugin detects columns formatted as date/time in Google Sheets and treats them as time. |
| `text`   | No       | The description shown when you hover over the annotation. |
| `tags`   | No       | Comma-separated tags to categorize and filter annotations. |
| `timeend`| No       | End timestamp for range annotations (shaded region instead of a vertical line). |

Always enable **Use Time Filter** in the annotation query so only events in the dashboard time range are fetched.

## Example sheet layout

A simple sheet for annotations might look like this:

| time                | text              | tags        |
|---------------------|-------------------|-------------|
| 2025-01-15 10:00:00 | Deployed v2.1     | deploy,prod |
| 2025-01-15 14:30:00 | Incident resolved | incident    |
| 2025-01-16 09:00:00 | Release v2.2      | deploy      |

- **time**: Format the column as a date or date-time in Google Sheets so the plugin detects it as a time field.
- **text**: Short description for the annotation.
- **tags**: Optional; use commas to separate multiple tags.

Your range could be something like `Annotations!A1:C100` (adjust the sheet name and row count to match your data).

## Customize annotation appearance

After creating the annotation query, you can change how it appears:

- **Color** – Choose a color for the annotation markers.
- **Show in** – Choose which panels display the annotations (all panels, selected panels, or all except selected).
- **Filter by** – Add filters to limit when annotations are shown.

For details, see [Annotate visualizations](https://grafana.com/docs/grafana/latest/dashboards/build-dashboards/annotate-visualizations/).

## Next steps

- [Query editor](query-editor.md) – Spreadsheet ID, Range, and Use Time Filter
- [Configure the data source](configure.md)
