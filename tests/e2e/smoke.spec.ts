import { test, expect } from '@grafana/plugin-e2e';

test('Smoke test: plugin loads', async ({ createDataSourceConfigPage, page }) => {
  await createDataSourceConfigPage({ type: 'grafana-googlesheets-datasource' });

  await expect(await page.getByText('Type: Google Sheets', { exact: true })).toBeVisible();
  await expect(page.getByRole('group', { name: 'Authentication', exact: true })).toBeVisible();
});
