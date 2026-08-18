import { expect, test, type DataSourceSettings } from '@grafana/plugin-e2e';
import { GoogleSheetsDataSourceOptions, GoogleSheetsSecureJSONData } from '../../src/types';

const PROVISIONING_FILE = 'datasources.yml';

test.describe('Config editor', () => {
  test.describe('rendering', () => {
    test('smoke: should render config editor', { tag: '@plugins' }, async ({ createDataSourceConfigPage, page }) => {
      await createDataSourceConfigPage({ type: 'grafana-googlesheets-datasource' });

      await expect(page.getByRole('group', { name: 'Authentication', exact: true })).toBeVisible();
    });

    test('should render authentication settings', async ({ createDataSourceConfigPage, page }) => {
      await createDataSourceConfigPage({ type: 'grafana-googlesheets-datasource' });

      await expect(page.getByRole('heading', { name: 'Choosing an authentication type' })).toBeVisible();
      await expect(page.getByRole('radio', { name: 'API Key' })).toBeVisible();
      await expect(page.getByRole('radio', { name: 'JWT button' })).toBeVisible();
      await expect(page.getByRole('radio', { name: 'GCE button' })).toBeVisible();
    });

    test('should render the default spreadsheet setting', async ({ createDataSourceConfigPage, page }) => {
      await createDataSourceConfigPage({ type: 'grafana-googlesheets-datasource' });

      await expect(page.getByText('Default Spreadsheet ID', { exact: true })).toBeVisible();
      await expect(page.getByText('Optional spreadsheet ID to use as default when creating new queries')).toBeVisible();
    });
  });

  test.describe('provisioned datasource', () => {
    let dataSource: DataSourceSettings<GoogleSheetsDataSourceOptions, GoogleSheetsSecureJSONData>;
    let runtimeUid = '';
    let shouldDeleteRuntimeDataSource = false;

    test.beforeEach(async ({ createDataSource, readProvisionedDataSource, request }, testInfo) => {
      dataSource = await readProvisionedDataSource<GoogleSheetsDataSourceOptions, GoogleSheetsSecureJSONData>({
        fileName: PROVISIONING_FILE,
      });

      const provisionedResponse = await request.get(`/api/datasources/uid/${dataSource.uid}`);
      if (provisionedResponse.ok()) {
        runtimeUid = dataSource.uid;
        shouldDeleteRuntimeDataSource = false;
        return;
      }

      const runtimeDataSource = await createDataSource({
        type: dataSource.type,
        name: `${dataSource.name} E2E ${testInfo.workerIndex}-${Date.now()}`,
        uid: `google-sheets-config-e2e-${testInfo.workerIndex}-${Date.now()}`,
        jsonData: dataSource.jsonData,
      });
      runtimeUid = runtimeDataSource.uid;
      shouldDeleteRuntimeDataSource = true;
    });

    test.afterEach(async ({ request }) => {
      if (shouldDeleteRuntimeDataSource) {
        await request.delete(`/api/datasources/uid/${runtimeUid}`);
      }
    });

    test('should load provisioned authentication settings', async ({ gotoDataSourceConfigPage, page }) => {
      await gotoDataSourceConfigPage(runtimeUid);

      await expect(page.getByRole('radio', { name: 'API Key' })).toBeChecked();
    });

    test('should load the provisioned default spreadsheet', async ({ gotoDataSourceConfigPage, page }) => {
      await gotoDataSourceConfigPage(runtimeUid);

      await expect(page.getByRole('button', { name: dataSource.jsonData.defaultSheetID })).toBeVisible();
    });
  });

  test.describe('save & test', () => {
    test('should pass health check for provisioned datasource', async ({
      gotoDataSourceConfigPage,
      page,
      readProvisionedDataSource,
      request,
    }) => {
      const dataSource = await readProvisionedDataSource<GoogleSheetsDataSourceOptions, GoogleSheetsSecureJSONData>({
        fileName: PROVISIONING_FILE,
      });
      const healthResponse = await request.get(`/api/datasources/uid/${dataSource.uid}/health`);
      test.skip(!healthResponse.ok(), 'Google Sheets credentials are not configured');

      const configPage = await gotoDataSourceConfigPage(dataSource.uid);
      await page.getByRole('button', { name: /^(Save & test|Test)$/ }).click();

      await expect(configPage).toHaveAlert('success');
    });

    test('should show error alert when health check fails', async ({ createDataSourceConfigPage, page }) => {
      const configPage = await createDataSourceConfigPage({ type: 'grafana-googlesheets-datasource' });
      await configPage.mockHealthCheckResponse({ status: 'ERROR', message: 'mocked failure' }, 400);

      await page.getByRole('button', { name: 'Save & test' }).click();

      await expect(configPage).toHaveAlert('error', { hasText: 'mocked failure' });
    });

    test('should show error alert when Google Sheets is unreachable', async ({ createDataSourceConfigPage, page }) => {
      const configPage = await createDataSourceConfigPage({ type: 'grafana-googlesheets-datasource' });
      await configPage.mockHealthCheckResponse({ status: 'ERROR', message: 'Google Sheets API is unreachable' }, 503);

      await page.getByRole('button', { name: 'Save & test' }).click();

      await expect(configPage).toHaveAlert('error', { hasText: 'Google Sheets API is unreachable' });
    });
  });
});
