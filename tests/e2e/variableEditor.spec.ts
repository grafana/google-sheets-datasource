import { test, expect } from '@grafana/plugin-e2e';

test.describe('Variable Query Editor', () => {
  test('should render variable query editor fields', async ({
    createDataSourceConfigPage,
    page,
  }) => {
    await createDataSourceConfigPage({ type: 'grafana-googlesheets-datasource' });

    // Navigate to a new dashboard's variable settings
    await page.goto('/dashboard/new?editview=variables');
    await page.getByRole('button', { name: 'Add variable' }).click();

    // Embedded query editor fields
    await expect(page.getByText('Spreadsheet ID')).toBeVisible();
    await expect(page.getByText('Range', { exact: true })).toBeVisible();
    await expect(page.getByText('Cache Time')).toBeVisible();
    await expect(page.getByText('Use Time Filter')).toBeVisible();

    // Variable-specific fields
    await expect(page.getByText('Value Field')).toBeVisible();
    await expect(page.getByText('Label Field')).toBeVisible();
    await expect(page.getByText('Filter Field')).toBeVisible();
    await expect(page.getByText('Filter Value')).toBeVisible();
    await expect(page.getByPlaceholder('Enter filter value')).toBeVisible();
  });
});
