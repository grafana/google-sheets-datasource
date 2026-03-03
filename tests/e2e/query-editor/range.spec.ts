import { test, expect } from '@grafana/plugin-e2e';

test.describe('Query Editor - Range Field', () => {
  test('should render Range input with correct placeholder', async ({
    panelEditPage,
    readProvisionedDataSource,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'default.yaml' });
    await panelEditPage.datasource.set(ds.name);

    await expect(page.getByText('Range', { exact: true })).toBeVisible();

    const rangeInput = page.getByPlaceholder('Class Data!A2:E');
    await expect(rangeInput).toBeVisible();
    await expect(rangeInput).toHaveValue('');
  });

  test('should accept a range value', async ({
    panelEditPage,
    readProvisionedDataSource,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'default.yaml' });
    await panelEditPage.datasource.set(ds.name);

    const rangeInput = page.getByPlaceholder('Class Data!A2:E');
    await rangeInput.click();
    await rangeInput.fill('Sheet1!A1:D50');

    await expect(rangeInput).toHaveValue('Sheet1!A1:D50');
  });

  test('should accept a named range', async ({
    panelEditPage,
    readProvisionedDataSource,
    page,
  }) => {
    const ds = await readProvisionedDataSource({ fileName: 'default.yaml' });
    await panelEditPage.datasource.set(ds.name);

    const rangeInput = page.getByPlaceholder('Class Data!A2:E');
    await rangeInput.click();
    await rangeInput.fill('SalesData');
    await rangeInput.press('Tab');

    await expect(rangeInput).toHaveValue('SalesData');
  });
});
