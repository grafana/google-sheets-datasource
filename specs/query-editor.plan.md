# Google Sheets Datasource - Query Editor Test Plan

## Application Overview

The Google Sheets datasource plugin (`grafana-googlesheets-datasource`) lets Grafana users visualize data from Google Spreadsheets. The query editor is the primary interface through which users configure which spreadsheet and range to read data from. It exposes four configurable fields: Spreadsheet ID (a searchable async segment picker that can accept a raw ID, a full Google Sheets URL, or a value chosen from a list of available sheets), Range (a free-text cell-range input such as `Sheet1!A1:B10`), Cache Time (a segmented picker for how long responses are cached), and Use Time Filter (an inline switch that applies the dashboard time range to the first time field in the sheet). The plugin is tested against Grafana >= 11.6.0 and uses `@grafana/plugin-e2e` fixtures throughout.

## Test Scenarios

### 1. Query Editor - Spreadsheet ID Field

**Seed:** `tests/e2e/seed.spec.ts`

#### 1.1. Renders Spreadsheet ID field with placeholder text on fresh query

**File:** `tests/e2e/query-editor/spreadsheet-id.spec.ts`

**Steps:**
  1. Use the `panelEditPage` fixture to open a new panel edit page and select the `grafana-googlesheets-datasource` datasource (read from provisioning via `readProvisionedDataSource`).
    - expect: The query editor is visible with the label 'Spreadsheet ID' rendered as an inline form label.
  2. Locate the Spreadsheet ID segment picker button (the element that shows the placeholder 'Enter SpreadsheetID').
    - expect: The segment picker displays the placeholder text 'Enter SpreadsheetID', indicating no spreadsheet has been selected.
  3. Verify the 'Open link' icon button that links to the Google Sheets document is NOT present.
    - expect: The link button is absent because no spreadsheet ID is set.

#### 1.2. Manually entering a raw Spreadsheet ID triggers a query run

**File:** `tests/e2e/query-editor/spreadsheet-id.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource.
    - expect: The query editor is visible.
  2. Click the Spreadsheet ID segment picker to open the input field.
    - expect: An editable input field appears inside the segment.
  3. Type a raw spreadsheet ID such as `1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgVE2upms` into the input and press Enter.
    - expect: The segment displays the typed ID as its selected value.
    - expect: A query run is triggered (network request to the datasource backend is made).
  4. Verify the 'Open link' icon button labeled 'Open link' is now visible next to the Spreadsheet ID field.
    - expect: The link button appears and its `href` attribute points to the corresponding Google Sheets URL containing the entered spreadsheet ID.

#### 1.3. Pasting a full Google Sheets URL extracts the spreadsheet ID and range

**File:** `tests/e2e/query-editor/spreadsheet-id.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource.
    - expect: The query editor is visible.
  2. Click the Spreadsheet ID segment picker and paste the full URL `https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgVE2upms/edit#gid=0&range=B19:F20` into the input, then confirm with Enter.
    - expect: The segment picker updates to show the extracted spreadsheet ID `1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgVE2upms`.
  3. Inspect the Range input field value.
    - expect: The Range field is automatically populated with `B19:F20`, extracted from the URL's `range=` query parameter.

#### 1.4. Spreadsheet ID segment shows loading indicator while fetching available spreadsheets

**File:** `tests/e2e/query-editor/spreadsheet-id.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource.
    - expect: The query editor is visible.
  2. Click the Spreadsheet ID segment picker button to open it and trigger the async load of available spreadsheets.
    - expect: A loading indicator or spinner is shown while the plugin fetches the list of accessible spreadsheets from the backend.

#### 1.5. Tooltip for Spreadsheet ID is accessible via the info icon

**File:** `tests/e2e/query-editor/spreadsheet-id.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource.
    - expect: The query editor is visible with the 'Spreadsheet ID' label and its accompanying info tooltip icon.
  2. Hover over the info icon next to the 'Spreadsheet ID' label.
    - expect: A tooltip appears explaining that the spreadsheetId is the value between '/d/' and '/edit' in the spreadsheet URL.

### 2. Query Editor - Range Field

**Seed:** `tests/e2e/seed.spec.ts`

#### 2.1. Renders Range input with correct placeholder

**File:** `tests/e2e/query-editor/range.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource.
    - expect: The query editor is visible with a 'Range' label and an associated text input.
  2. Inspect the Range input field without interacting with it.
    - expect: The Range input field has placeholder text 'Class Data!A2:E' and is currently empty.

#### 2.2. Typing a range value and blurring triggers a query run

**File:** `tests/e2e/query-editor/range.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource.
    - expect: The query editor is visible.
  2. Click into the Range text input and type `Sheet1!A1:D50`.
    - expect: The input reflects the typed value `Sheet1!A1:D50`.
  3. Click outside the Range input (trigger blur/onBlur).
    - expect: A query run is triggered (network request to the datasource backend is made with the updated range value).

#### 2.3. Range value persists after re-opening the query editor

**File:** `tests/e2e/query-editor/range.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource. Enter a spreadsheet ID and a range value `MySheet!A1:C100`.
    - expect: Both fields reflect the entered values.
  2. Save the panel via the Grafana UI, then reopen the panel edit page using `gotoPanelEditPage` with the saved panel reference.
    - expect: The Range input still contains `MySheet!A1:C100`, confirming the value is persisted to the query model.

#### 2.4. Named range string is accepted in the Range field

**File:** `tests/e2e/query-editor/range.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource.
    - expect: The query editor is visible.
  2. Type a named range string `SalesData` into the Range input and press Tab to blur.
    - expect: The input retains the value `SalesData` and a query run is triggered.

#### 2.5. Range tooltip is accessible via the info icon

**File:** `tests/e2e/query-editor/range.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource.
    - expect: The query editor is visible with the 'Range' label and its info tooltip icon.
  2. Hover over the info icon next to the 'Range' label.
    - expect: A tooltip appears describing the A1 notation format (e.g., `Sheet1!A1:B2`) and that named ranges are also supported.

### 3. Query Editor - Cache Time Field

**Seed:** `tests/e2e/seed.spec.ts`

#### 3.1. Cache Time defaults to 5m (300 seconds) on a fresh query

**File:** `tests/e2e/query-editor/cache-time.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource.
    - expect: The query editor is visible.
  2. Locate the Cache Time segment picker next to the 'Cache Time' label.
    - expect: The Cache Time segment displays the value '5m', reflecting the default of 300 seconds.

#### 3.2. Selecting a different Cache Time option updates the segment and saves the value

**File:** `tests/e2e/query-editor/cache-time.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource.
    - expect: The query editor is visible with Cache Time defaulting to '5m'.
  2. Click the Cache Time segment picker to open the dropdown.
    - expect: A dropdown menu appears listing available cache duration options: '0s', '5s', '10s', '30s', '1m', '2m', '5m', '10m', '30m', '1h', '2h', '5h'.
  3. Select '30m' from the dropdown.
    - expect: The segment displays '30m' as the newly selected Cache Time value.

#### 3.3. Selecting cache time of 0s (no cache) shows the correct label

**File:** `tests/e2e/query-editor/cache-time.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource.
    - expect: The query editor is visible.
  2. Click the Cache Time segment picker and select the '0s' option.
    - expect: The segment displays '0s'. The option description 'Response is not cached at all' was visible in the dropdown.

#### 3.4. All expected Cache Time options are present in the dropdown

**File:** `tests/e2e/query-editor/cache-time.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource.
    - expect: The query editor is visible.
  2. Click the Cache Time segment picker to open its dropdown.
    - expect: The dropdown list contains exactly: '0s', '5s', '10s', '30s', '1m', '2m', '5m', '10m', '30m', '1h', '2h', '5h'.

#### 3.5. Cache Time tooltip is accessible via the info icon

**File:** `tests/e2e/query-editor/cache-time.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource.
    - expect: The query editor is visible with the 'Cache Time' label and its info icon.
  2. Hover over the info icon next to the 'Cache Time' label.
    - expect: A tooltip appears explaining that Cache Time is the number of seconds the spreadsheet will be cached in Grafana after receiving a response from the API.

### 4. Query Editor - Use Time Filter Toggle

**Seed:** `tests/e2e/seed.spec.ts`

#### 4.1. Use Time Filter is off by default on a fresh query

**File:** `tests/e2e/query-editor/time-filter.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource.
    - expect: The query editor is visible.
  2. Locate the InlineSwitch next to the 'Use Time Filter' label and check its state.
    - expect: The toggle is in the unchecked (off) state by default.

#### 4.2. Toggling Use Time Filter on triggers a query run

**File:** `tests/e2e/query-editor/time-filter.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource.
    - expect: The query editor is visible with 'Use Time Filter' off.
  2. Click the 'Use Time Filter' InlineSwitch (using `{ force: true }` in the underlying test because Grafana InlineSwitch requires it).
    - expect: The toggle switches to the checked (on) state.
    - expect: A query run is triggered.

#### 4.3. Toggling Use Time Filter off after it was on triggers a query run

**File:** `tests/e2e/query-editor/time-filter.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource. Enable Use Time Filter by clicking the switch.
    - expect: Use Time Filter is checked.
  2. Click the 'Use Time Filter' InlineSwitch again to disable it.
    - expect: The toggle returns to the unchecked (off) state.
    - expect: A second query run is triggered.

#### 4.4. Use Time Filter tooltip is accessible via the info icon

**File:** `tests/e2e/query-editor/time-filter.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource.
    - expect: The query editor is visible with the 'Use Time Filter' label and its info icon.
  2. Hover over the info icon next to the 'Use Time Filter' label.
    - expect: A tooltip appears with the text 'Apply the dashboard time range to the first time field'.

### 5. Query Editor - Full Happy Path

**Seed:** `tests/e2e/seed.spec.ts`

#### 5.1. Complete query configuration from blank state runs successfully

**File:** `tests/e2e/query-editor/happy-path.spec.ts`

**Steps:**
  1. Use `panelEditPage` to open a new panel edit page and select the `grafana-googlesheets-datasource` datasource read from provisioning via `readProvisionedDataSource`.
    - expect: The query editor appears with all four fields: Spreadsheet ID, Range, Cache Time, Use Time Filter.
  2. Click the Spreadsheet ID segment picker and type in a valid spreadsheet ID `1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgVE2upms`, then press Enter.
    - expect: The segment displays the entered ID.
    - expect: The 'Open link' icon button becomes visible.
  3. Click into the Range input and type `Sheet1!A1:E10`, then press Tab to blur.
    - expect: The Range input shows `Sheet1!A1:E10`.
  4. Click the Cache Time segment picker and select '10m'.
    - expect: The Cache Time segment displays '10m'.
  5. Enable the 'Use Time Filter' toggle.
    - expect: The Use Time Filter switch is in the checked state.
  6. Call `panelEditPage.refreshPanel()` and await the result.
    - expect: The panel refreshes without a Grafana error alert, indicating the query configuration is structurally valid and was sent to the backend.

#### 5.2. Open link button navigates to the correct Google Sheets URL

**File:** `tests/e2e/query-editor/happy-path.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource. Enter spreadsheet ID `1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgVE2upms` and range `Sheet1!A1:B5`.
    - expect: Both fields are filled in and the 'Open link' button is visible.
  2. Inspect the `href` attribute of the 'Open link' button.
    - expect: The `href` equals `https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgVE2upms/view#range=Sheet1!A1:B5`, incorporating both the spreadsheet ID and the range.

#### 5.3. Open link button URL omits range fragment when range is empty

**File:** `tests/e2e/query-editor/happy-path.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource. Enter spreadsheet ID `1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgVE2upms` but leave the Range field empty.
    - expect: The 'Open link' button is visible.
  2. Inspect the `href` attribute of the 'Open link' button.
    - expect: The `href` equals `https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgVE2upms/view` with no `#range=` fragment.

### 6. Query Editor - Variable Template Support

**Seed:** `tests/e2e/seed.spec.ts`

#### 6.1. Template variable syntax is accepted in the Spreadsheet ID field

**File:** `tests/e2e/query-editor/template-variables.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource.
    - expect: The query editor is visible.
  2. Click the Spreadsheet ID segment picker and enter `$sheetId` (a Grafana template variable reference), then press Enter.
    - expect: The segment displays `$sheetId` as its value without error, demonstrating that template variable interpolation syntax is accepted by the field.

#### 6.2. Template variable syntax is accepted in the Range field

**File:** `tests/e2e/query-editor/template-variables.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource.
    - expect: The query editor is visible.
  2. Type `$sheetName!A1:E10` into the Range input field and blur.
    - expect: The Range input retains the value `$sheetName!A1:E10` without error.

### 7. Query Editor - Error States

**Seed:** `tests/e2e/seed.spec.ts`

#### 7.1. Query with empty Spreadsheet ID shows an error or empty data in the panel

**File:** `tests/e2e/query-editor/error-states.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource. Leave the Spreadsheet ID empty.
    - expect: The query editor is visible with no spreadsheet selected.
  2. Call `panelEditPage.refreshPanel()` to attempt a query run without a spreadsheet ID.
    - expect: The panel either shows no data, an error message from the backend, or a Grafana alert indicating the query is incomplete.

#### 7.2. Invalid Spreadsheet ID produces an error response from the backend

**File:** `tests/e2e/query-editor/error-states.spec.ts`

**Steps:**
  1. Open a new panel edit page with the `grafana-googlesheets-datasource` datasource. Enter an obviously invalid spreadsheet ID such as `INVALID_SHEET_ID_12345`.
    - expect: The segment displays the invalid ID.
  2. Call `panelEditPage.refreshPanel()` and await the result.
    - expect: The panel displays an error state or Grafana alert indicating the spreadsheet could not be found or accessed.

### 8. Variable Query Editor

**Seed:** `tests/e2e/seed.spec.ts`

#### 8.1. Variable query editor renders all fields including Value Field, Label Field, Filter Field, and Filter Value

**File:** `tests/e2e/query-editor/variable-query-editor.spec.ts`

**Steps:**
  1. Use the `variableEditPage` fixture to open the variable edit page. Set the datasource to `grafana-googlesheets-datasource` (read from provisioning via `readProvisionedDataSource`).
    - expect: The variable query editor loads and shows all shared query editor fields: Spreadsheet ID, Range, Cache Time, Use Time Filter.
  2. Scroll down or inspect below the standard query fields.
    - expect: Four additional fields are visible: 'Value Field' select, 'Label Field' select, an 'Optional filtering' section header, 'Filter Field' select, and 'Filter Value' text input.

#### 8.2. Value Field and Label Field selects show loading state while fetching column choices

**File:** `tests/e2e/query-editor/variable-query-editor.spec.ts`

**Steps:**
  1. Open the variable edit page with `grafana-googlesheets-datasource`. Enter a spreadsheet ID into the Spreadsheet ID field to trigger a column fetch.
    - expect: The 'Value Field', 'Label Field', and 'Filter Field' selects each display a 'Loading...' placeholder while fetching column names from the spreadsheet.
  2. Wait for the loading to complete.
    - expect: The 'Loading...' placeholder disappears from all three selects.

#### 8.3. Filter Value input accepts free-text input

**File:** `tests/e2e/query-editor/variable-query-editor.spec.ts`

**Steps:**
  1. Open the variable edit page with `grafana-googlesheets-datasource`.
    - expect: The variable query editor is visible with the Filter Value input field.
  2. Locate the Filter Value input (identified by `data-testid="filter-value-input"`) and type `TargetValue` into it.
    - expect: The input retains the typed text `TargetValue`.

#### 8.4. Optional filtering section header is visible and separates filter fields from core fields

**File:** `tests/e2e/query-editor/variable-query-editor.spec.ts`

**Steps:**
  1. Open the variable edit page with `grafana-googlesheets-datasource`.
    - expect: The variable query editor renders all its sections.
  2. Locate the 'Optional filtering' label/header in the editor.
    - expect: The text 'Optional filtering' is present as a section label, visually separating the Value Field and Label Field selects from the Filter Field and Filter Value inputs.
