import { test, expect } from '@grafana/plugin-e2e';

test('data query should be OK when URL is valid', async ({ panelEditPage, page }) => {
  const API_URL = 'https://jsonplaceholder.typicode.com/users';
  await panelEditPage.datasource.set('Infinity E2E');
  await page.getByTestId('infinity-query-url-input').fill(API_URL);
  await expect(panelEditPage.refreshPanel()).toBeOK();
});
