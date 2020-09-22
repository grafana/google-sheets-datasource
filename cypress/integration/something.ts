import { e2e } from '@grafana/e2e';

const addGoogleSheetsDataSource = (config: { apiKey: string } | { jwtPath: string }) => {
  const fillApiKey = (newValue: string) =>
    e2e()
      .get('[placeholder="Enter API Key"]')
      .scrollIntoView()
      .type(newValue);

  const fillJwt = (newValue: string) =>
    e2e()
      .contains('.gf-form-group', 'Auth')
      .within(() => {
        e2e.flows.selectOption({
          container: e2e.components.Select.option(),
          optionText: 'Google JWT File',
        });

        e2e()
          .get('input[type=file]')
          .attachFile(newValue);
      });

  e2e.flows.addDataSource({
    checkHealth: true,
    expectedAlertMessage: 'Success',
    form: () => {
      if ('apiKey' in config) {
        fillApiKey(config.apiKey);
      } else if ('jwtPath' in config) {
        fillJwt(config.jwtPath);
      }
    },
    type: 'Google Sheets',
  });
};

interface GoogleSheetsPanelConfig {
  sheetId: string;
  sheetIdVariable: string;
}

const addGoogleSheetsPanel = ({ sheetId, sheetIdVariable }: GoogleSheetsPanelConfig) => {
  const fillSpreadsheetID = (newValue: string, previousValue = 'Enter SpreadsheetID', avoidFlakiness = false) =>
    e2e.components.QueryEditorRows.rows()
      .contains('.gf-form-inline', 'Spreadsheet ID')
      .within(() => {
        e2e()
          .get(`.gf-form:has(.gf-form-label:contains("${previousValue}"))`) // the <Label/>
          .click({ force: avoidFlakiness });
        e2e()
          .get('input')
          .scrollIntoView()
          .type(`${newValue}{enter}`);
      });

  // Assert sheet id
  e2e.flows.addPanel({
    matchScreenshot: true,
    queriesForm: () => fillSpreadsheetID(sheetId),
    visualizationName: e2e.flows.VISUALIZATION_TABLE,
  }).then(({ config: { panelTitle } }: any) => {
    // Assert template variable as sheet id
    e2e.flows.editPanel({
      matchScreenshot: true,
      panelTitle,
      queriesForm: () => fillSpreadsheetID(`$${sheetIdVariable}`, sheetId, true),
      visitDashboardAtStart: false,
    });
  });
};

e2e.scenario({
  describeName: 'Smoke tests',
  itName: 'Login, create data source, dashboard and panel',
  scenario: () => {
    e2e().readProvisions([
      // Paths are relative to <project-root>/provisioning
      'datasources/google-sheets-datasource-API-key.yaml',
      'datasources/google-sheets-datasource-jwt.yaml',
    ]).then(([apiKeyProvision, jwtProvision]) => {
      const { apiKey } = apiKeyProvision.datasources[0].secureJsonData;
      //const { jwt } = jwtProvision.datasources[0].secureJsonData;
      const sheetId = '1Kn_9WKsuT-H0aJL3fvqukt27HlizMLd-KQfkNgeWj4U';
      const sheetIdVariable = 'sheetId';

      //const jwtPath = 'jwt.json';
      //e2e().writeFile(`${Cypress.config('fixturesFolder')}/${jwtPath}`, jwt);

      // This gets auto-removed within `afterEach` of @grafana/e2e
      addGoogleSheetsDataSource({ apiKey });

      // This gets auto-removed within `afterEach` of @grafana/e2e
      e2e.flows.addDashboard({
        variables: [
          {
            constantValue: sheetId,
            label: 'Template Variable',
            name: sheetIdVariable,
            type: e2e.flows.VARIABLE_TYPE_CONSTANT,
          },
        ],
      });

      // This gets auto-removed within `afterEach` of @grafana/e2e
      addGoogleSheetsPanel({ sheetId, sheetIdVariable });

      // These get auto-removed within `afterEach` of @grafana/e2e
      //addGoogleSheetsDataSource({ jwtPath });
      //addGoogleSheetsPanel({ sheetId, sheetIdVariable });
    });
  },
});
