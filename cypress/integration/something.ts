/// <reference path="../../node_modules/@grafana/e2e/cypress/support/index.d.ts" />
import { e2e } from '@grafana/e2e';

const addGoogleSheetsDataSource = (apiKey: string) => {
  const fillApiKey = () => getByPlaceholder('Enter API Key').scrollIntoView().type(apiKey);

  // This gets auto-removed within `afterEach` of @grafana/e2e
  e2e.flows.addDataSource({
    checkHealth: true,
    expectedAlertMessage: 'Success',
    form: () => fillApiKey(),
    name: 'Google Sheets',
  });
};

const addGoogleSheetsPanel = (spreadsheetId: string) => {
  const fillSpreadsheetID = () => {
    const getContainer = e2e.components.QueryTab.content;

    e2e().server();
    e2e().route('POST', '/api/ds/query').as('chartData');

    getContainer()
      .contains('.gf-form-label', 'Enter SpreadsheetID')
      .parent('.gf-form') // the <Label/>
      .click({ force: true }); // https://github.com/cypress-io/cypress/issues/7306

    getContainer()
      .contains('.gf-form-input', 'Choose')
      .find('.gf-form-select-box__input input')
      .scrollIntoView()
      .type(spreadsheetId);

    // Persist the value
    getContainer().click();

    e2e().wait('@chartData');
  };

  // @todo remove `@ts-ignore` when possible
  // @ts-ignore
  e2e.getScenarioContext().then(({ lastAddedDataSource }) => {
    // This gets auto-removed within `afterEach` of @grafana/e2e
    e2e.flows.addPanel({
      dataSourceName: lastAddedDataSource,
      queriesForm: () => fillSpreadsheetID(),
      visualizationName: 'Table',
    }).then((panelTitle: string) => {
      e2e.components.PanelEditor.OptionsPane.close().click();
      e2e.components.Panels.Panel.containerByTitle(panelTitle).find('.panel-content').screenshot('chart');
      e2e().compareScreenshots('chart');
    });
  });
};

export const getByPlaceholder = (placeholder: string) => e2e().get(`[placeholder="${placeholder}"]`);

e2e.scenario({
  describeName: 'Smoke tests',
  itName: 'Login, create data source, dashboard and panel',
  scenario: () => {
    // Paths are relative to <project-root>/provisioning
    const provisionPaths = [
      'datasources/google-sheets-datasource-API-key.yml',
      'datasources/google-sheets-datasource-jwt.yml',
    ];

    e2e().readProvisions(provisionPaths).then(([apiKeyProvision, jwtProvision]) => {
      addGoogleSheetsDataSource(apiKeyProvision.datasources[0].secureJsonData.apiKey);
      e2e.flows.addDashboard();
      addGoogleSheetsPanel('1Kn_9WKsuT-H0aJL3fvqukt27HlizMLd-KQfkNgeWj4U');
    });
  },
});
