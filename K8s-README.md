# K8s-README

Temporarily, while development is still ongoing on the POC for this plugin as an aggregated API server, this readme
contains information to get started for developers. Once this process is completed, any documnetation changes will
be absorbed into the main README at root.

## Background

1. Until entity-storage is exposed to the plugin in form of an API, `etcd` and standard K8s codegen are prerequisites for the POC.
2. `mage` target has been redefined in `.bra.toml` to provide better debugging capability using appropriate gcflags and ldflags. The current platform binary is generated at root of the repo (not in `dist/`).
3. We have found the ability to run in standalone mode as useful: this means, no Grafana is used and plugin binary is launched directly.
4. Authz is disabled. You can `kubectl` or `curl` against the API without worrying about tokens.

## Setup

1. `brew install etcd` and start `etcd` on standard port `2379`.
2. `mage watch` in one terminal tab.
3. `./gpx_sheets_darwin_arm64 k8s` launches the aggregated API server. NOTE the additional command line argument 
4. `export KUBECONFIG=./data/grafana.kubeconfig`
5. Apply a datasource (google sheet connection) manually. `defaultProject`, `clientEmail` and `privateKey` must be filled from your google service account credentials. Private key should contain `\n` for newlines.
    ```json
    {
      "apiVersion": "googlesheets.ext.grafana.com/v1",
      "kind": "Datasource",
      "metadata": {
        "name": "12345"
      },
      "spec": {
        "authType": "jwt",
        "apiKey": "",
        "defaultProject": "",
        "jwt": "",
        "clientEmail": "",
        "tokenUri": "https://oauth2.googleapis.com/token",
        "authenticationType": "jwt",
        "privateKey": ""
      }
    }
    ```
6. Once applied, you should be able to do query and resource requests of following formats:
    ```json
    {
      "queries": [{
        "refId": "A",
        "datasource": {
          "type": "grafana-googlesheets-datasource",
          "uid": "b1808c48-9fc9-4045-82d7-081781f8a553"
        },
        "cacheDurationSeconds": 300,
        "spreadsheet": "spreadsheetID",
        "range": "",
        "datasourceId": 4,
        "intervalMs": 30000,
        "maxDataPoints": 794
      }],
      "from": "1692624667389",
      "to": "1692646267389"
    }

    ```

    ```shell
    curl -ik 'https://localhost:6443/apis/googlesheets.ext.grafana.com/v1/namespaces/default/datasources/12345/resource/spreadsheets'
   {"spreadsheets":{}}%    
    curl -X POST -ik 'https://localhost:6443/apis/googlesheets.ext.grafana.com/v1/namespaces/default/datasources/12345/query' -d @./request.json
    ```
