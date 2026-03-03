import { test, expect } from '@grafana/plugin-e2e';

test.describe('Query Editor - Cache Time Field', () => {
  test('should default to 5m on a fresh query', async ({ panelEditPage, readProvisionedDataSource, page }) => {
    const ds = await readProvisionedDataSource({ fileName: 'default.yaml' });
    await panelEditPage.datasource.set(ds.name);

    await expect(page.getByText('Cache Time')).toBeVisible();
    await expect(page.getByRole('button', { name: '5m' })).toBeVisible();
  });

  test('should allow selecting a different cache time', async ({ panelEditPage, readProvisionedDataSource, page }) => {
    const ds = await readProvisionedDataSource({ fileName: 'default.yaml' });
    await panelEditPage.datasource.set(ds.name);

    const cacheTimeSegment = page.getByRole('button', { name: '5m' });
    await expect(cacheTimeSegment).toBeVisible();

    await cacheTimeSegment.click();
    await page.getByRole('option', { name: '30m' }).click();

    await expect(page.getByRole('button', { name: '30m' })).toBeVisible();
  });

  test('should show all expected cache time options', async ({ panelEditPage, readProvisionedDataSource, page }) => {
    const ds = await readProvisionedDataSource({ fileName: 'default.yaml' });
    await panelEditPage.datasource.set(ds.name);

    const cacheTimeSegment = page.getByRole('button', { name: '5m' });
    await expect(cacheTimeSegment).toBeVisible();

    await cacheTimeSegment.click();

    // The options menu renders as a listbox with option roles
    // '0s' option includes description text, so use exact: true for others
    const optionsList = page.getByRole('listbox', { name: 'Select options menu' });
    await expect(optionsList.getByRole('option')).toHaveCount(12);

    const expectedOptions = ['0s', '5s', '10s', '30s', '1m', '2m', '5m', '10m', '30m', '1h', '2h', '5h'];
    for (const option of expectedOptions) {
      await expect(optionsList.getByText(option, { exact: true })).toBeVisible();
    }
  });
});
