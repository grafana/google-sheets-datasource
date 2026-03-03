import { test, expect } from '@grafana/plugin-e2e';

test.describe('Query Editor - Variable Template Support', () => {
  test('should accept template variable syntax in Spreadsheet ID field', async ({
    panelEditPage,
    readProvisionedDataSource,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'default.yaml' });
    await panelEditPage.datasource.set(ds.name);

    await expect(page.getByText('Spreadsheet ID')).toBeVisible();

    await page.getByRole('button', { name: 'Enter SpreadsheetID' }).click();
    const combobox = page.getByRole('combobox').last();
    await combobox.fill('$sheetId');
    await combobox.press('Enter');

    await expect(page.getByText('$sheetId')).toBeVisible();
  });

  test('should accept template variable syntax in Range field', async ({
    panelEditPage,
    readProvisionedDataSource,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'default.yaml' });
    await panelEditPage.datasource.set(ds.name);

    const rangeInput = page.getByPlaceholder('Class Data!A2:E');
    await rangeInput.fill('$sheetName!A1:E10');
    await rangeInput.press('Tab');

    await expect(rangeInput).toHaveValue('$sheetName!A1:E10');
  });
});
