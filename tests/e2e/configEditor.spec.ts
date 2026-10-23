import { expect, test } from '@grafana/plugin-e2e';
import { GoogleSheetsDataSourceOptions, GoogleSheetsSecureJSONData } from '../../src/types';

test.describe('Config editor', () => {
  test('"Save & test" should be successful for the provisioned JWT datasource', async ({
    readProvisionedDataSource,
    gotoDataSourceConfigPage,
    request,
  }) => {
    const ds = await readProvisionedDataSource<GoogleSheetsDataSourceOptions, GoogleSheetsSecureJSONData>({
      fileName: 'datasources.yml',
    });

    const configPage = await gotoDataSourceConfigPage(ds.uid);

    await expect(configPage.saveAndTest()).toBeOK();
    await expect(configPage).toHaveAlert('success');
  });
});
