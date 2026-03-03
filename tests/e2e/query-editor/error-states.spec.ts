import { test, expect } from '@grafana/plugin-e2e';

test.describe('Query Editor - Error States', () => {
  test('should handle empty Spreadsheet ID gracefully', async ({
    panelEditPage,
    readProvisionedDataSource,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'default.yaml' });
    await panelEditPage.datasource.set(ds.name);

    // Leave Spreadsheet ID empty
    await expect(page.getByRole('button', { name: 'Enter SpreadsheetID' })).toBeVisible();

    // Refresh panel - should not crash
    await panelEditPage.refreshPanel();
  });

  test('should show error for invalid Spreadsheet ID', async ({
    panelEditPage,
    readProvisionedDataSource,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'default.yaml' });
    await panelEditPage.datasource.set(ds.name);

    // Enter invalid spreadsheet ID
    await page.getByRole('button', { name: 'Enter SpreadsheetID' }).click();
    // Scope the combobox to the query editor row to avoid selecting the wrong combobox
    const queryEditorRow = page.locator('[data-testid="query-editor-row"]');
    const combobox = queryEditorRow.getByRole('combobox').last();
    await combobox.fill('INVALID_SHEET_ID_12345');
    // Wait for the async options to load and the "Hit enter to add" create option to appear
    await page.getByRole('option', { name: /Hit enter to add/ }).waitFor({ state: 'visible' });
    await combobox.press('Enter');

    // The spreadsheet ID appears as a button in the SegmentAsync component
    await expect(page.getByRole('button', { name: 'INVALID_SHEET_ID_12345' })).toBeVisible();

    // Refresh panel and expect error
    await panelEditPage.refreshPanel();
    // The plugin shows query errors as an inline message in the query editor row,
    // not as a Grafana Alert component. The panel status error icon confirms the error state.
    await expect(page.locator('[data-testid="data-testid Panel status error"]')).toBeVisible();
  });
});
