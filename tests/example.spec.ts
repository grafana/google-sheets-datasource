import { test, expect } from '@grafana/plugin-e2e';

const GOOGLE_SHEETS_SPREADSHEETS = {
  spreadsheets: {
    sheet1: 'Datasource test spreadsheet',
    sheet2: 'Google Sheets Datasource - Average Temperature',
  },
};

test('should list spreadsheets when clicking on spreadsheet segment', async ({
  panelEditPage,
  page,
  readProvision,
}) => {
  const sheetsDataSource = await readProvision({
    filePath: 'datasources/google-sheets-datasource-jwt.yaml',
  }).then((provision) => provision.datasources?.[0]!);
  await panelEditPage.datasource.set(sheetsDataSource.name!);
  await panelEditPage.mockResourceResponse('spreadsheets', GOOGLE_SHEETS_SPREADSHEETS);
  await panelEditPage.getQueryEditorRow('A').getByText('Enter SpreadsheetID').click();
  await expect(page.getByText(GOOGLE_SHEETS_SPREADSHEETS.spreadsheets.sheet1, { exact: true })).toHaveCount(1);
  await expect(page.getByText(GOOGLE_SHEETS_SPREADSHEETS.spreadsheets.sheet2, { exact: true })).toHaveCount(1);
});
