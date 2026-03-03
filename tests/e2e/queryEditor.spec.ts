import { test, expect } from '@grafana/plugin-e2e';

test.describe('Query Editor', () => {
  test('should render all query editor fields in explore', async ({ createDataSourceConfigPage, explorePage, page }) => {
    const configPage = await createDataSourceConfigPage({ type: 'grafana-googlesheets-datasource' });
    await explorePage.goto();
    await explorePage.datasource.set(configPage.datasource.name);
    await explorePage.mockResourceResponse('spreadsheets', []);

    await expect(page.getByText('Spreadsheet ID')).toBeVisible();
    await expect(page.getByRole('button', { name: 'Enter SpreadsheetID' })).toBeVisible();
    await expect(page.getByRole('textbox', { name: 'Class Data!A2:E' })).toBeVisible();
    await expect(page.getByText('Cache Time')).toBeVisible();
    await expect(page.getByText('5m')).toBeVisible();
    await expect(page.getByText('Use Time Filter')).toBeVisible();
  });

  test('should run query with mocked response', async ({ createDataSourceConfigPage, explorePage }) => {
    const configPage = await createDataSourceConfigPage({ type: 'grafana-googlesheets-datasource' });
    await explorePage.goto();
    await explorePage.datasource.set(configPage.datasource.name);
    await explorePage.mockResourceResponse('spreadsheets', []);
    await explorePage.mockQueryDataResponse({
      results: {
        A: {
          status: 200,
          frames: [
            {
              schema: {
                fields: [
                  { name: 'time', type: 'time' },
                  { name: 'value', type: 'number' },
                ],
              },
              data: {
                values: [
                  [1609459200000, 1609545600000],
                  [10, 20],
                ],
              },
            },
          ],
        },
      },
    });

    await expect(explorePage.runQuery()).toBeOK();
  });

  test('should display error when query returns error', async ({ createDataSourceConfigPage, explorePage }) => {
    const configPage = await createDataSourceConfigPage({ type: 'grafana-googlesheets-datasource' });
    await explorePage.goto();
    await explorePage.datasource.set(configPage.datasource.name);
    await explorePage.mockResourceResponse('spreadsheets', []);
    await explorePage.mockQueryDataResponse({}, 400);

    await expect(explorePage.runQuery()).not.toBeOK();
  });
});
