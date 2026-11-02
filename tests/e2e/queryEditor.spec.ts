import { expect, test } from '@grafana/plugin-e2e';

test.describe('Query editor in provisioned dashboard', () => {
  const DASHBOARD_UID = '_gFwRQuZks';

  test('should display live spreadsheet data in a gauge panel', async ({ gotoPanelEditPage }) => {
    const panelEditPage = await gotoPanelEditPage({ dashboard: { uid: DASHBOARD_UID }, id: '3' });
    await panelEditPage.timeRange.set({ from: '2019-01-01 00:00:00', to: '2019-12-31 23:59:59' });
    await panelEditPage.setVisualization('Gauge');
    await expect(panelEditPage.refreshPanel()).toBeOK();
  });

  test('should display live spreadsheet data in a table panel', async ({ gotoPanelEditPage }) => {
    const panelEditPage = await gotoPanelEditPage({ dashboard: { uid: DASHBOARD_UID }, id: '2' });
    await panelEditPage.timeRange.set({ from: '2019-01-01 00:00:00', to: '2019-12-31 23:59:59' });
    await panelEditPage.setVisualization('Table');
    await expect(panelEditPage.refreshPanel()).toBeOK();
  });
});
