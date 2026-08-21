import { expect, test, type ExplorePage } from '@grafana/plugin-e2e';

const PLUGIN_ID = 'grafana-googlesheets-datasource';
const DEFAULT_SPREADSHEET_ID = process.env.DS_INSTANCE_SPREADSHEET_ID ?? '1TZlZX67Y0s4CvRro_3pCYqRCKuXer81oFp_xcsjPpe8';
const API_KEY = process.env.DS_INSTANCE_API_KEY ?? process.env.GOOGLE_SHEETS_API_KEY;
const CLIENT_EMAIL = process.env.DS_INSTANCE_CLIENT_EMAIL;
const PRIVATE_KEY = process.env.DS_INSTANCE_PRIVATE_KEY;
const DEFAULT_PROJECT = process.env.DS_INSTANCE_DEFAULT_PROJECT;
const TOKEN_URI = process.env.DS_INSTANCE_TOKEN_URI ?? 'https://oauth2.googleapis.com/token';
const HAS_LIVE_CREDENTIALS = Boolean(API_KEY || (CLIENT_EMAIL && PRIVATE_KEY));

// GRAFANA_URL is injected only by the Cloud cron workflow (playwright-cloud). Its presence
// means missing live credentials indicate a broken Vault injection, not a local/PR run.
const isCloudRun = Boolean(process.env.GRAFANA_URL);

function skipOrFailIfNoLiveCredentials(hasCredentials: boolean) {
  if (hasCredentials) {
    return;
  }

  if (isCloudRun) {
    throw new Error('Cloud e2e is missing Google Sheets credentials from Vault');
  }

  test.skip(true, 'Google Sheets credentials are not configured');
}

const MOCK_QUERY_RESPONSE = {
  results: {
    A: {
      status: 200,
      frames: [
        {
          schema: {
            refId: 'A',
            fields: [
              { name: 'Time', type: 'time' },
              { name: 'Value', type: 'number' },
            ],
          },
          data: {
            values: [
              [1787032800000, 1787032860000],
              [10, 12],
            ],
          },
        },
      ],
    },
  },
};

interface QueryResponse {
  results?: {
    A?: {
      frames?: unknown[];
    };
  };
}

function exploreUrl(
  datasourceUid: string,
  query: { spreadsheet?: string; range?: string; cacheDurationSeconds?: number; useTimeFilter?: boolean } = {}
) {
  const panes = JSON.stringify({
    explore: {
      datasource: datasourceUid,
      queries: [
        {
          refId: 'A',
          datasource: { type: PLUGIN_ID, uid: datasourceUid },
          spreadsheet: query.spreadsheet ?? DEFAULT_SPREADSHEET_ID,
          range: query.range ?? 'A1:E5',
          cacheDurationSeconds: query.cacheDurationSeconds ?? 300,
          useTimeFilter: query.useTimeFilter ?? false,
        },
      ],
      range: { from: 'now-6h', to: 'now' },
    },
  });

  return `/explore?orgId=1&schemaVersion=1&panes=${encodeURIComponent(panes)}`;
}

function waitForQueryDataResponseWithBody(explorePage: ExplorePage) {
  let body: QueryResponse | undefined;
  const responsePromise = explorePage.waitForQueryDataResponse(async (response) => {
    if (!response.ok()) {
      return false;
    }

    const responseBody = (await response.json().catch(() => undefined)) as QueryResponse | undefined;
    if (!Array.isArray(responseBody?.results?.A?.frames)) {
      return false;
    }

    body = responseBody;
    return true;
  });

  return { responsePromise, getBody: () => body };
}

function liveDataSourceConfig() {
  if (CLIENT_EMAIL && PRIVATE_KEY) {
    return {
      jsonData: {
        authenticationType: 'jwt',
        clientEmail: CLIENT_EMAIL,
        defaultProject: DEFAULT_PROJECT,
        defaultSheetID: DEFAULT_SPREADSHEET_ID,
        tokenUri: TOKEN_URI,
      },
      secureJsonData: {
        privateKey: PRIVATE_KEY,
      },
    };
  }

  return {
    jsonData: {
      authenticationType: 'key',
      defaultSheetID: DEFAULT_SPREADSHEET_ID,
    },
    secureJsonData: {
      apiKey: API_KEY,
    },
  };
}

test.describe('Query editor', () => {
  let datasourceUid = '';

  test.beforeEach(async ({ createDataSource }, testInfo) => {
    const dataSource = await createDataSource({
      type: PLUGIN_ID,
      name: `Google Sheets E2E ${testInfo.workerIndex}-${Date.now()}`,
      uid: `google-sheets-e2e-${testInfo.workerIndex}-${Date.now()}`,
      jsonData: {
        authenticationType: 'key',
        defaultSheetID: DEFAULT_SPREADSHEET_ID,
      },
      secureJsonData: {
        apiKey: API_KEY ?? 'e2e-placeholder',
      },
    });
    datasourceUid = dataSource.uid;
  });

  test.afterEach(async ({ request }) => {
    await request.delete(`/api/datasources/uid/${datasourceUid}`);
  });

  test.describe('rendering', () => {
    test('smoke: should render all query fields', { tag: '@plugins' }, async ({ explorePage, page }) => {
      await explorePage.mockQueryDataResponse(MOCK_QUERY_RESPONSE);
      await page.goto(exploreUrl(datasourceUid));

      const queryRow = explorePage.getQueryEditorRow('A');
      await expect(queryRow.getByText('Spreadsheet ID', { exact: true })).toBeVisible();
      await expect(queryRow.getByText('Range', { exact: true })).toBeVisible();
      await expect(queryRow.getByText('Cache Time', { exact: true })).toBeVisible();
      await expect(queryRow.getByText('Use Time Filter', { exact: true })).toBeVisible();
    });

    test('should render the configured spreadsheet and range', async ({ explorePage, page }) => {
      await explorePage.mockQueryDataResponse(MOCK_QUERY_RESPONSE);
      await page.goto(exploreUrl(datasourceUid, { range: 'Sheet1!A2:E' }));

      const queryRow = explorePage.getQueryEditorRow('A');
      await expect(queryRow.getByRole('button', { name: DEFAULT_SPREADSHEET_ID })).toBeVisible();
      await expect(queryRow.getByPlaceholder('Class Data!A2:E')).toHaveValue('Sheet1!A2:E');
    });
  });

  test.describe('query options', () => {
    test('should allow editing the range', async ({ explorePage, page }) => {
      await explorePage.mockQueryDataResponse(MOCK_QUERY_RESPONSE);
      await page.goto(exploreUrl(datasourceUid));

      const rangeInput = explorePage.getQueryEditorRow('A').getByPlaceholder('Class Data!A2:E');
      await rangeInput.fill('Sheet1!A1:B10');

      await expect(rangeInput).toHaveValue('Sheet1!A1:B10');
    });

    test('should toggle the time filter', async ({ explorePage, page }) => {
      await explorePage.mockQueryDataResponse(MOCK_QUERY_RESPONSE);
      await page.goto(exploreUrl(datasourceUid));

      const timeFilter = explorePage.getQueryEditorRow('A').getByRole('switch');
      await timeFilter.check({ force: true });

      await expect(timeFilter).toBeChecked();
    });

    test('should receive query data', async ({ explorePage, page }) => {
      await explorePage.mockQueryDataResponse(MOCK_QUERY_RESPONSE);
      const { responsePromise, getBody } = waitForQueryDataResponseWithBody(explorePage);

      await page.goto(exploreUrl(datasourceUid));
      await responsePromise;

      expect(getBody()?.results?.A?.frames).toHaveLength(1);
    });
  });
});

test.describe('Query editor with live Google Sheets data', () => {
  test.describe.configure({ mode: 'serial' });
  let datasourceUid = '';

  test.beforeEach(async ({ createDataSource, request }, testInfo) => {
    skipOrFailIfNoLiveCredentials(HAS_LIVE_CREDENTIALS);

    const dataSource = await createDataSource({
      type: PLUGIN_ID,
      name: `Google Sheets live E2E ${testInfo.workerIndex}-${Date.now()}`,
      uid: `google-sheets-live-e2e-${testInfo.workerIndex}-${Date.now()}`,
      ...liveDataSourceConfig(),
    });
    datasourceUid = dataSource.uid;

    const healthResponse = await request.get(`/api/datasources/uid/${datasourceUid}/health`);
    skipOrFailIfNoLiveCredentials(healthResponse.ok());
  });

  test.afterEach(async ({ request }) => {
    if (datasourceUid) {
      await request.delete(`/api/datasources/uid/${datasourceUid}`);
    }
  });

  test('should return rows from the public test spreadsheet', async ({ explorePage, page }) => {
    const { responsePromise, getBody } = waitForQueryDataResponseWithBody(explorePage);

    await page.goto(
      exploreUrl(datasourceUid, {
        spreadsheet: DEFAULT_SPREADSHEET_ID,
        range: 'A1:E20',
        cacheDurationSeconds: 0,
      })
    );
    await responsePromise;

    expect(getBody()?.results?.A?.frames?.length).toBeGreaterThan(0);
  });
});
