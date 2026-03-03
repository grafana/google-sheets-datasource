import { test, expect } from '@grafana/plugin-e2e';

async function enterSpreadsheetId(page: any, id: string) {
  await page.getByRole('button', { name: 'Enter SpreadsheetID' }).click();
  // Scope the combobox to the query editor row to avoid selecting the wrong combobox
  const queryEditorRow = page.locator('[data-testid="query-editor-row"]');
  const combobox = queryEditorRow.getByRole('combobox').last();
  await combobox.fill(id);
  // Wait for the async options to load and the "Hit enter to add" create option to appear
  await page.getByRole('option', { name: /Hit enter to add/ }).waitFor({ state: 'visible' });
  await combobox.press('Enter');
}

test.describe('Query Editor - Spreadsheet ID Field', () => {
  test('should render Spreadsheet ID field with placeholder text on fresh query', async ({
    panelEditPage,
    readProvisionedDataSource,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'default.yaml' });
    await panelEditPage.datasource.set(ds.name);

    await expect(page.getByText('Spreadsheet ID')).toBeVisible();
    await expect(page.getByRole('button', { name: 'Enter SpreadsheetID' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Open link' })).not.toBeVisible();
  });

  test('should accept a raw Spreadsheet ID and show Open link button', async ({
    panelEditPage,
    readProvisionedDataSource,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'default.yaml' });
    const spreadsheetId = '1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgVE2upms';

    await panelEditPage.datasource.set(ds.name);
    await enterSpreadsheetId(page, spreadsheetId);

    await expect(page.getByText(spreadsheetId)).toBeVisible();
    await expect(page.getByRole('link', { name: 'Open link' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Open link' })).toHaveAttribute(
      'href',
      `https://docs.google.com/spreadsheets/d/${spreadsheetId}/view`
    );
  });

  test('should extract spreadsheet ID and range from a pasted Google Sheets URL', async ({
    panelEditPage,
    readProvisionedDataSource,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'default.yaml' });
    const spreadsheetId = '1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgVE2upms';
    const range = 'Sheet1!A1:B5';
    const sheetsUrl = `https://docs.google.com/spreadsheets/d/${spreadsheetId}/edit#range=${range}`;

    await panelEditPage.datasource.set(ds.name);
    await enterSpreadsheetId(page, sheetsUrl);

    await expect(page.getByText(spreadsheetId)).toBeVisible();
    await expect(page.getByRole('link', { name: 'Open link' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Open link' })).toHaveAttribute(
      'href',
      `https://docs.google.com/spreadsheets/d/${spreadsheetId}/view#range=${range}`
    );
  });
});
