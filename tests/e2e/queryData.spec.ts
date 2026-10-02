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

  test('should hand the panel the data from the same request when refreshing', async ({ gotoPanelEditPage }) => {
    const panelEditPage = await gotoPanelEditPage({ dashboard: { uid: DASHBOARD_UID }, id: '1' });
    const { response, body, data } = await panelEditPage.refreshPanelWithData();

    await expect(response).toBeOK();
    expect(body).toHaveProperty('results.A.frames');
    expect(data.state).toBe('Done');
    expect(data.errors).toEqual([]);
    expect(data.series).toHaveLength(1);
    expect(data.series[0].fields.length).toBeGreaterThan(1);
  });
});
