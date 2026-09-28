import { expect, test } from '@grafana/plugin-e2e';

test.describe('Query editor in provisioned dashboard', () => {
  const DASHBOARD_UID = '_gFwRQuZks';
  const PANEL_ID = '3';

  test.beforeEach(async ({ request }) => {
    const dashboardResponse = await request.get(`/api/dashboards/uid/${DASHBOARD_UID}`);
    test.skip(!dashboardResponse.ok(), 'Average Temperature dashboard is not provisioned');
  });

  test('should display live spreadsheet data in a table panel', async ({ gotoPanelEditPage }) => {
    const panelEditPage = await gotoPanelEditPage({ dashboard: { uid: DASHBOARD_UID }, id: PANEL_ID });
    await panelEditPage.timeRange.set({ from: '2019-01-01 00:00:00', to: '2019-12-31 23:59:59' });
    await panelEditPage.setVisualization('Table');

    await expect(panelEditPage.refreshPanel()).toBeOK();
    await expect(panelEditPage.panel.fieldNames).toContainText(['Month', 'Stockholm', 'Berlin', 'Los Angeles']);
    await expect(panelEditPage.panel.data).toContainText(['-1 °C', '2.90 °C', '13 °C']);
  });
});
