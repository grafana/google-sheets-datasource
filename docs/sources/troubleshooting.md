---
title: Troubleshoot Google Sheets data Source issues
menuTitle: Troubleshooting
description: Fix common issues with the Google Sheets data source.
keywords:
  - data source
  - google sheets
  - troubleshoot
  - errors
  - permissions
labels:
  products:
    - oss
    - enterprise
    - cloud
last_reviewed: 2025-02-11
weight: 500
---

# Troubleshoot Google Sheets data source issues

This guide helps you fix common issues when configuring or using the Google Sheets data source.

## Save & test and connection errors

These messages appear when you click **Save & test** on the data source configuration page.

### Unable to load settings

**Cause:** Grafana could not read the data source configuration (for example, invalid or corrupted settings).

**Solution:**

- Ensure the data source configuration was saved correctly. Re-enter authentication details if needed.
- If you use provisioning, check that the YAML or API payload is valid and that secrets are available.
- Check Grafana server logs for the underlying error (e.g. "error reading settings").

### Unable to create client

**Cause:** The plugin could not create a Google API client. Usually an authentication or configuration problem.

**Common sub-messages and fixes:**

| Message or cause | Solution |
|------------------|----------|
| **missing AuthenticationType setting** | In the data source config, select an authentication method: **Google JWT File**, **API key**, or **GCE Default Service Account**. Save and test again. |
| **missing API Key** | For API key authentication, paste your API key in the **API Key** field (or ensure the provisioned secret is set). |
| **datasource is missing authentication details** | For **Google JWT File**, you must provide either a JWT file (upload/paste) or **Client email**, **Private key**, and **Default project**. See [Configure the data source](configure.md#authenticate-with-a-service-account-jwt). |
| **error parsing JWT file** | The pasted or uploaded JWT is not valid JSON or is malformed. Re-download the service account key from Google Cloud Console and paste or upload it again. Ensure the full JSON is used with no extra characters. |
| **Failed to create http client** / **unable to retrieve Sheets client** | Check Grafana logs for details. Often related to JWT format, network, or TLS. |

### Permissions check failed

**Cause:** **Save & test** successfully created a client but the test request to Google failed. The plugin tests access by reading a public test spreadsheet and (for JWT) listing Drive files.

**Solution:**

- **Network:** Ensure the Grafana server can reach Google APIs (`https://sheets.googleapis.com`, `https://www.googleapis.com`). If you use a proxy or firewall, allow these endpoints.
- **API key:** If using an API key, ensure the Google Sheets API (and Drive API if you use “Select Spreadsheet ID”) is enabled for the key and that key restrictions (e.g. IP, referrer) allow requests from Grafana.
- **JWT / service account:** Ensure the [Google Sheets API](https://console.cloud.google.com/apis/library/sheets.googleapis.com) and [Google Drive API](https://console.cloud.google.com/apis/library/drive.googleapis.com) are enabled for the project. Ensure the service account has access to at least one spreadsheet (e.g. share the sheet with the service account email). For “Invalid grant” or “account not found”, verify the service account key is correct and that the account has not been deleted or disabled.
- **GCE default account:** If using GCE Default Service Account, ensure Grafana runs on a Google Compute Engine VM and that the default service account has the required scopes and access to the sheet.

## Query and panel errors

These can appear in the panel, in the query editor, or in the response.

### Spreadsheet not found

**Cause:** The Google Sheets API returned 404 for the given Spreadsheet ID.

**Solution:**

- Verify the **Spreadsheet ID** in the query. In Google Sheets, the ID is in the URL: `https://docs.google.com/spreadsheets/d/<SPREADSHEET_ID>/edit`.
- Ensure the spreadsheet has not been deleted.
- For JWT/service account: share the spreadsheet with the service account email (e.g. `something@project.iam.gserviceaccount.com`) with at least **Viewer** access. For API key: the spreadsheet must be **published to the web** or otherwise publicly readable if your key only allows public data.

### Google API Error 403

**Cause:** Google rejected the request (forbidden). Usually permissions or API configuration.

**Solution:**

- **Sharing (JWT/service account):** Share the spreadsheet with the service account email with **Viewer** (or **Editor** if you need write; the plugin only reads).
- **API key:** Ensure the spreadsheet is shared so that “Anyone with the link” can view, or use a key that has access to the sheet. Check [API key restrictions](https://console.cloud.google.com/apis/credentials) so the key is allowed for the Sheets API (and Drive API if listing spreadsheets).
- **APIs not enabled:** In Google Cloud Console, enable [Google Sheets API](https://console.cloud.google.com/apis/library/sheets.googleapis.com) and [Google Drive API](https://console.cloud.google.com/apis/library/drive.googleapis.com) for the project.
- **Quotas:** If you hit rate limits, you may see errors; see [Quota](_index.md#quota) and consider increasing **Cache Time** to reduce requests.

### No data or empty panel

**Cause:** The query returned no rows or the range/spreadsheet is wrong.

**Solution:**

- Confirm **Spreadsheet ID** and **Range** (e.g. `Sheet1!A1:D10`). The range must exist and the sheet name must match (case-sensitive).
- If **Use Time Filter** is enabled, the panel only shows rows where the time column falls within the dashboard time range. Widen the dashboard time range or ensure the sheet has a column the plugin detects as time (date/datetime format) and that its values are inside the selected range.
- Check that the sheet has data in the specified range and that the first row is the header row.

### Invalid time column / error while parsing date

**Cause:** **Use Time Filter** is enabled but the chosen time column has invalid or non-parsable values.

**Solution:**

- Use a column that contains real dates or times. In Google Sheets, format the column as a date or date-time so the plugin can detect it. Avoid mixed types or text that does not look like a date in the same column.
- If the error mentions “error while parsing date”, fix or remove invalid cells in that column so every value is a valid date/time or leave the cell empty (empty may be skipped depending on behavior).

### Unable to create Google API client (in panel)

**Cause:** Same as [Unable to create client](#unable-to-create-client) but occurring when a panel runs a query (e.g. after config change or on load).

**Solution:** Fix authentication and configuration as in [Save & test and connection errors](#save--test-and-connection-errors), then re-run the query or reload the dashboard.

## Template variables

### Variable returns no options

**Cause:** The variable query returns no rows or the value/label columns are wrong.

**Solution:**

- Confirm **Spreadsheet ID** and **Range** in the variable query. Ensure **Value Field** is set to a column that exists in the range.
- If you use **Optional filtering** (**Filter Field** and **Filter Value**), ensure at least some rows match; otherwise the list will be empty.
- Test the same range in a panel query to confirm the sheet returns data.

## Annotations

### Annotations do not appear

**Cause:** Annotation query returns no data or Grafana cannot map time/text columns.

**Solution:**

- Use column headers **time** and **text** in your sheet (see [Annotations](annotations.md#query-requirements)). Ensure the **time** column is formatted as date/datetime in Google Sheets.
- Enable **Use Time Filter** on the annotation query and ensure the dashboard time range covers the events in the sheet.
- Verify **Spreadsheet ID** and **Range** and that the annotation query is enabled (toggle on) in Dashboard settings → Annotations.

## API quotas and rate limits

Google enforces [per-minute quotas](https://developers.google.com/sheets/api/limits) for the Sheets API. If you see rate-limit or quota errors:

- Increase **Cache Time** in the query (or data source default) so the same range is not requested too often.
- Reduce the number of panels or variables that query the same or many spreadsheets in a short time.
- For heavy use, consider requesting a quota increase in Google Cloud Console.

## Get additional help

If you've tried the solutions above and still encounter issues:

- Check the [Grafana community forums](https://community.grafana.com/) for similar issues.
- Review the [Google Sheets data source GitHub issues](https://github.com/grafana/google-sheets-datasource/issues) for known bugs.
- Enable debug logging in Grafana to capture detailed error information.
- Contact [Grafana Support](https://grafana.com/help/) if you're an Enterprise, Cloud Pro or Cloud Contracted user.

When reporting issues, include:

- Grafana version
- Error messages (redact sensitive information)
- Steps to reproduce
- Relevant configuration (redact credentials)

## Related documentation

- [Configure the data source](configure.md)
- [Query editor](query-editor.md)
- [Template variables](template-variables.md)
- [Annotations](annotations.md)
