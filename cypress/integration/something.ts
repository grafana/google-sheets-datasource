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

const addGoogleSheetsPanel = (spreadsheetId: string) => {
  const fillSpreadsheetID = () =>
    e2e.components.QueryEditorRows.rows().within(() => {
      e2e()
        .get('.gf-form:has(.gf-form-label:contains("Enter SpreadsheetID"))') // the <Label/>
        .click({ force: true }); // https://github.com/cypress-io/cypress/issues/7306
      e2e()
        .contains('.gf-form-input', 'Choose')
        .find('.gf-form-select-box__input input')
        .scrollIntoView()
        .type(`${spreadsheetId}{enter}`);
    });

  e2e.flows.addPanel({
    matchScreenshot: true,
    queriesForm: () => fillSpreadsheetID(),
    visualizationName: e2e.flows.VISUALIZATION_TABLE,
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

      //const jwtPath = 'jwt.json';
      //e2e().writeFile(`${Cypress.config('fixturesFolder')}/${jwtPath}`, jwt);

      // These gets auto-removed within `afterEach` of @grafana/e2e
      addGoogleSheetsDataSource({ apiKey });
      e2e.flows.addDashboard();
      addGoogleSheetsPanel(sheetId);

      // These gets auto-removed within `afterEach` of @grafana/e2e
      //addGoogleSheetsDataSource({ jwtPath });
      //e2e.flows.addDashboard();
      //addGoogleSheetsPanel(sheetId);
    });
  },
});
