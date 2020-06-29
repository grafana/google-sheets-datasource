/// <reference path="../../node_modules/@grafana/e2e/cypress/support/index.d.ts" />
import { e2e } from '@grafana/e2e';

const addGoogleSheetsDataSource = (config: { apiKey: string } | { jwtPath: string }) => {
  const fillApiKey = (apiKey: string) => getByPlaceholder('Enter API Key').scrollIntoView().type(apiKey);

  const fillJwt = (jwtPath: string) => {
    e2e().contains('.gf-form-group', 'Auth').within(() => {
      e2e.flows.selectOption(e2e.components.Select.option(), 'Google JWT File');
      e2e().get('input[type=file]').attachFile(jwtPath);
    });
  };

  // This gets auto-removed within `afterEach` of @grafana/e2e
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
  const fillSpreadsheetID = () => {
    e2e.components.QueryTab.content().within(() => {
      e2e()
        .contains('.gf-form-label', 'Enter SpreadsheetID')
        .parent('.gf-form') // the <Label/>
        .click({ force: true }); // https://github.com/cypress-io/cypress/issues/7306

      e2e()
        .contains('.gf-form-input', 'Choose')
        .find('.gf-form-select-box__input input')
        .scrollIntoView()
        .type(`${spreadsheetId}{enter}`);
    });
  };

  // This gets auto-removed within `afterEach` of @grafana/e2e
  e2e.flows.addPanel({
    queriesForm: () => fillSpreadsheetID(),
  }).then(({ config }: any) => {
    e2e.components.Panels.Panel.containerByTitle(config.panelTitle)
      .find('.panel-content')
      .screenshot('chart');
    e2e().compareScreenshots('chart');
  });
};

export const getByPlaceholder = (placeholder: string) => e2e().get(`[placeholder="${placeholder}"]`);

e2e.scenario({
  describeName: 'Smoke tests',
  itName: 'Login, create data source, dashboard and panel',
  scenario: () => {
    e2e().readProvisions([
      // Paths are relative to <project-root>/provisioning
      'datasources/google-sheets-datasource-API-key.yml',
      'datasources/google-sheets-datasource-jwt.yml',
    ]).then(([apiKeyProvision, jwtProvision]) => {
      const { apiKey } = apiKeyProvision.datasources[0].secureJsonData;
      //const { jwt } = jwtProvision.datasources[0].secureJsonData;
      const sheetId = '1Kn_9WKsuT-H0aJL3fvqukt27HlizMLd-KQfkNgeWj4U';

      //const jwtPath = 'jwt.json';
      //e2e().writeFile(`${Cypress.config('fixturesFolder')}/${jwtPath}`, jwt);

      addGoogleSheetsDataSource({ apiKey });
      e2e.flows.addDashboard();
      addGoogleSheetsPanel(sheetId);

      //addGoogleSheetsDataSource({ jwtPath });
      //e2e.flows.addDashboard();
      //addGoogleSheetsPanel(sheetId);
    });
  },
});
