import { CoreApp, DataQueryRequest } from '@grafana/data';
import { reportInteraction } from '@grafana/runtime';
import { SheetsQuery } from 'types';

export const trackRequest = (request: DataQueryRequest<SheetsQuery>) => {
  if (request.app === CoreApp.Dashboard || request.app === CoreApp.PanelViewer) {
    return;
  }

  request.targets.forEach((target) => {
    reportInteraction('grafana_google_sheets_query_executed', {
      app: request.app,
      useTimeFilter: target.useTimeFilter ?? false,
      cacheDurationSeconds: target.cacheDurationSeconds ?? 0,
    });
  });
};
