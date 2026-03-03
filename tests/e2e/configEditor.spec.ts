import { test, expect } from '@grafana/plugin-e2e';

test.describe('Config Editor', () => {
  test('should render config page with default JWT auth type', async ({ createDataSourceConfigPage, page }) => {
    await createDataSourceConfigPage({ type: 'grafana-googlesheets-datasource' });

    await expect(page.getByText('Type: Google Sheets', { exact: true })).toBeVisible();
    await expect(page.locator('legend', { hasText: 'Authentication' })).toBeVisible();
    await expect(page.getByRole('radio', { name: 'JWT button' })).toBeChecked();
    await expect(page.locator('legend', { hasText: 'JWT Key Details' })).toBeVisible();
    await expect(page.getByText('Default Spreadsheet ID')).toBeVisible();
    await expect(page.getByRole('button', { name: 'Configure Google Sheets Authentication' })).toBeVisible();
  });

  test('should switch to API Key auth and show API key field', async ({ createDataSourceConfigPage, page }) => {
    await createDataSourceConfigPage({ type: 'grafana-googlesheets-datasource' });

    await page.getByRole('radio', { name: 'API Key' }).check({ force: true });

    await expect(page.getByPlaceholder('Enter API key')).toBeVisible();
    await expect(page.locator('legend', { hasText: 'JWT Key Details' })).not.toBeVisible();
  });

  test('should switch to GCE auth and hide JWT/API Key fields', async ({ createDataSourceConfigPage, page }) => {
    await createDataSourceConfigPage({ type: 'grafana-googlesheets-datasource' });

    await page.getByRole('radio', { name: 'GCE button' }).check({ force: true });

    await expect(page.locator('legend', { hasText: 'JWT Key Details' })).not.toBeVisible();
    await expect(page.getByPlaceholder('Enter API key')).not.toBeVisible();
  });

  test('should expand help section and show JWT help by default', async ({ createDataSourceConfigPage, page }) => {
    await createDataSourceConfigPage({ type: 'grafana-googlesheets-datasource' });

    await page.getByRole('button', { name: 'Configure Google Sheets Authentication' }).click();

    await expect(page.getByRole('heading', { name: 'Generate a JWT file' })).toBeVisible();
  });

  test('should show API Key help when API Key auth is selected', async ({ createDataSourceConfigPage, page }) => {
    await createDataSourceConfigPage({ type: 'grafana-googlesheets-datasource' });

    await page.getByRole('radio', { name: 'API Key' }).check({ force: true });
    await page.getByRole('button', { name: 'Configure Google Sheets Authentication' }).click();

    await expect(page.getByRole('heading', { name: 'Generate an API key' })).toBeVisible();
  });

  test('should show GCE help when GCE auth is selected', async ({ createDataSourceConfigPage, page }) => {
    await createDataSourceConfigPage({ type: 'grafana-googlesheets-datasource' });

    await page.getByRole('radio', { name: 'GCE button' }).check({ force: true });
    await page.getByRole('button', { name: 'Configure Google Sheets Authentication' }).click();

    await expect(page.getByRole('heading', { name: 'Configure GCE Service Account' })).toBeVisible();
  });

  test('should save and test successfully with mocked health check', async ({
    createDataSourceConfigPage,
    page,
  }) => {
    const configPage = await createDataSourceConfigPage({ type: 'grafana-googlesheets-datasource' });

    await page.getByRole('radio', { name: 'API Key' }).check({ force: true });
    await page.getByPlaceholder('Enter API key').fill('test-api-key');

    await configPage.mockHealthCheckResponse({ message: 'Data source is working' });
    await expect(configPage.saveAndTest()).toBeOK();
    await expect(page.getByRole('alert')).toContainText('Data source is working');
  });

  test('should show error when save and test fails', async ({ createDataSourceConfigPage, page }) => {
    const configPage = await createDataSourceConfigPage({ type: 'grafana-googlesheets-datasource' });

    await configPage.mockHealthCheckResponse({ message: 'Authentication failed' }, 500);
    await expect(configPage.saveAndTest()).not.toBeOK();
    await expect(page.getByRole('alert')).toBeVisible();
  });
});
