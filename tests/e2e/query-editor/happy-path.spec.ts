import { test, expect } from '@grafana/plugin-e2e';

const SPREADSHEET_ID = '1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgVE2upms';
const BASE_SHEETS_URL = `https://docs.google.com/spreadsheets/d/${SPREADSHEET_ID}/view`;

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

test.describe('Query Editor - Full Happy Path', () => {
  test('should configure all query fields successfully', async ({ panelEditPage, readProvisionedDataSource, page }) => {
    const ds = await readProvisionedDataSource({ fileName: 'default.yaml' });
    await panelEditPage.datasource.set(ds.name);

    // Set Spreadsheet ID via SegmentAsync combobox
    await enterSpreadsheetId(page, SPREADSHEET_ID);

    await expect(page.getByRole('link', { name: 'Open link' })).toBeVisible();

    // Set Range
    await page.getByPlaceholder('Class Data!A2:E').fill('Sheet1!A1:E10');
    await page.keyboard.press('Tab');

    // Set Cache Time to 10m
    await page.getByRole('button', { name: '5m' }).click();
    await page.getByRole('option', { name: '10m' }).click();

    // Enable Use Time Filter
    const queryRow = page.locator('[data-testid="query-editor-row"]');
    const timeFilterSwitch = queryRow.locator('[role="switch"]').last();
    await expect(timeFilterSwitch).not.toBeChecked();
    await timeFilterSwitch.click({ force: true });

    // Refresh panel - no crash expected
    await panelEditPage.refreshPanel();
  });

  test('should construct correct Open link URL with range', async ({
    panelEditPage,
    readProvisionedDataSource,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'default.yaml' });
    await panelEditPage.datasource.set(ds.name);

    // Enter spreadsheet ID
    await enterSpreadsheetId(page, SPREADSHEET_ID);

    // Enter range
    await page.getByPlaceholder('Class Data!A2:E').fill('Sheet1!A1:B5');
    await page.keyboard.press('Tab');

    // Verify href includes #range=
    const openLinkButton = page.getByRole('link', { name: 'Open link' });
    await expect(openLinkButton).toBeVisible();
    await expect(openLinkButton).toHaveAttribute('href', `${BASE_SHEETS_URL}#range=Sheet1!A1:B5`);
  });

  test('should construct Open link URL without range fragment when range is empty', async ({
    panelEditPage,
    readProvisionedDataSource,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'default.yaml' });
    await panelEditPage.datasource.set(ds.name);

    // Enter spreadsheet ID only (no range)
    await enterSpreadsheetId(page, SPREADSHEET_ID);

    // Verify href has no #range= fragment
    const openLinkButton = page.getByRole('link', { name: 'Open link' });
    await expect(openLinkButton).toBeVisible();
    await expect(openLinkButton).toHaveAttribute('href', BASE_SHEETS_URL);
  });
});
