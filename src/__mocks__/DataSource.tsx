// import {DataSourceInstanceSettings, PluginType} from "@grafana/data";
// import {GoogleAuthType, SheetsSourceOptions} from "../types";
// import {getTemplateSrv} from '@grafana/runtime';
//
// import {CustomVariableModel} from '@grafana/data';
// import {DataSource} from "../DataSource";
//
// // export function setupMockedTemplateService(variables: CustomVariableModel[]) {
// //   const templateService = getTemplateSrv();
// //   templateService.init(variables);
// //   templateService.getVariables = jest.fn().mockReturnValue(variables);
// //   return templateService;
// // }
//
// export const DataSourceSettings: DataSourceInstanceSettings<SheetsSourceOptions> = {
//   jsonData: {authType: GoogleAuthType.JWT},
//   id: 0,
//   uid: '',
//   type: '',
//   meta: {
//     id: '',
//     name: '',
//     type: PluginType.datasource,
//     info: {
//       author: {
//         name: '',
//       },
//       description: '',
//       links: [],
//       logos: {
//         large: '',
//         small: '',
//       },
//       screenshots: [],
//       updated: '',
//       version: '',
//     },
//     module: '',
//     baseUrl: '',
//   },
//   access: 'direct',
//   name: 'Google Sheets Test Datasource'
// }
//
// export function setupMockedDataSource({
//                                         variables,
//                                         mockGetVariableName = true,
//                                         getMock = jest.fn(),
//                                         customInstanceSettings = DataSourceSettings,
//                                       }: {
//   getMock?: jest.Func;
//   variables?: CustomVariableModel[];
//   mockGetVariableName?: boolean;
//   customInstanceSettings?: DataSourceInstanceSettings<SheetsSourceOptions>;
// } = {}) {
//
//   const datasource = new DataSource(customInstanceSettings);
//   const fetchMock = jest.fn().mockReturnValue(of({}));
//   setBackendSrv({
//     ...getBackendSrv(),
//     fetch: fetchMock,
//     get: getMock,
//   });
//
//   return {datasource, fetchMock, templateService, timeSrv};
// }
