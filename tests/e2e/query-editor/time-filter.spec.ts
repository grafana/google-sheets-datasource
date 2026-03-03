import { test, expect } from '@grafana/plugin-e2e';

test.describe('Query Editor - Use Time Filter Toggle', () => {
  test('should be off by default on a fresh query', async ({ panelEditPage, readProvisionedDataSource, page }) => {
    const ds = await readProvisionedDataSource({ fileName: 'default.yaml' });
    await panelEditPage.datasource.set(ds.name);

    await expect(page.getByText('Use Time Filter')).toBeVisible();
    // Use data-testid to scope to the query editor row
    const queryRow = page.locator('[data-testid="query-editor-row"]');
    await expect(queryRow.locator('[role="switch"]').last()).not.toBeChecked();
  });

  test('should toggle on when clicked', async ({ panelEditPage, readProvisionedDataSource, page }) => {
    const ds = await readProvisionedDataSource({ fileName: 'default.yaml' });
    await panelEditPage.datasource.set(ds.name);

    // Wait for query editor to be fully rendered
    await expect(page.getByText('Use Time Filter')).toBeVisible();

    const queryRow = page.locator('[data-testid="query-editor-row"]');
    const timeFilterSwitch = queryRow.locator('[role="switch"]').last();
    await expect(timeFilterSwitch).not.toBeChecked();
    await timeFilterSwitch.click({ force: true });

    await expect(timeFilterSwitch).toBeChecked();
  });

  test('should toggle off after being toggled on', async ({ panelEditPage, readProvisionedDataSource, page }) => {
    const ds = await readProvisionedDataSource({ fileName: 'default.yaml' });
    await panelEditPage.datasource.set(ds.name);

    // Wait for query editor to be fully rendered
    await expect(page.getByText('Use Time Filter')).toBeVisible();

    const queryRow = page.locator('[data-testid="query-editor-row"]');
    const timeFilterSwitch = queryRow.locator('[role="switch"]').last();

    await expect(timeFilterSwitch).not.toBeChecked();
    await timeFilterSwitch.click({ force: true });
    await expect(timeFilterSwitch).toBeChecked();

    await timeFilterSwitch.click({ force: true });
    await expect(timeFilterSwitch).not.toBeChecked();
  });
});
