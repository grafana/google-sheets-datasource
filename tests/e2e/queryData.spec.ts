import { expect, test } from '@grafana/plugin-e2e';

test.describe('Query data in provisioned dashboard', () => {
  const DASHBOARD_UID = 'e2e-sheets-data';

  test('should return the expected data for all rows', async ({ gotoDashboardPage }) => {
    const dashboardPage = await gotoDashboardPage({ uid: DASHBOARD_UID });
    await expect(dashboardPage.getPanelByTitle('All rows')).toMatchDataSnapshot('all-rows');
  });

  test('should return only the rows inside the time range when the time filter is on', async ({
    gotoDashboardPage,
  }) => {
    const dashboardPage = await gotoDashboardPage({ uid: DASHBOARD_UID });
    await expect(dashboardPage.getPanelByTitle('Time filtered')).toMatchDataSnapshot('time-filtered');
  });
});
