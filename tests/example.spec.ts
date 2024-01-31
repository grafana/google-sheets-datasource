import { test, expect } from '@grafana/plugin-e2e';

test('should return data and not display panel error when a valid query is provided', async ({
  explorePage,
  page,
  readProvision,
}) => {
  const provision = await readProvision({ filePath: 'datasources/google-sheets-datasource-jwt.yaml' });
  await explorePage.datasource.set(provision.datasources?.[0]!.name!);
  await explorePage.timeRange.set({ from: '2019-01-11', to: '2019-12-15' });
  await explorePage.getQueryEditorRow('A').getByText('Enter SpreadsheetID').click();
  await page.keyboard.insertText('1TZlZX67Y0s4CvRro_3pCYqRCKuXer81oFp_xcsjPpe8');
  const responsePromise = page.waitForResponse((resp) => resp.url().includes('/api/ds/query'));
  await page.keyboard.press('Tab');
  await responsePromise;
  await expect(explorePage.runQuery()).toBeOK();
});
