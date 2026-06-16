import { test, expect } from '@grafana/plugin-e2e';

test('Smoke test: plugin loads', async ({ createDataSourceConfigPage, page }) => {
  await createDataSourceConfigPage({ type: 'grafana-googlesheets-datasource' });

  await expect(page.getByRole('group', { name: 'Authentication', exact: true })).toBeVisible();
});
