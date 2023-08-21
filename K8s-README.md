# K8s-README

Temporarily, while development is still ongoing on the POC for this plugin as an aggregated API server, this readme
contains information to get started for developers. Once this process is completed, any documnetation changes will
be absorbed into the main README at root.

## Background

1. Until entity-storage is exposed to the plugin in form of an API, `etcd` and standard K8s codegen are prerequisites for the POC.
2. `mage` target has been redefined in `.bra.toml` to provide better debugging capability using appropriate gcflags and ldflags. For now, the current platform binary is generated at root of the repo (not in `dist/`).
3. For now, we have found the ability to run in standalone mode as useful: this means, no Grafana is used and plugin binary is launched directly.
4. Authz is disabled for now. You can `kubectl` or `curl` against the API without worrying about tokens.

## Setup

1. `brew install etcd` and start `etcd` on standard port `2379`.
2. `mage watch` in one terminal tab.
3. `./gpx_sheets_darwin_arm64` launches the aggregated API server. 
4. Apply a datasource (google sheet connection) manually. `defaultProject`, `clientEmail` and `privateKey` must be filled from your google service account credentials. Private key should contain `\n` for newlines.
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
5. Once applied, you should be able to:
    ```shell
    curl -ik 'https://localhost:6443/apis/googlesheets.ext.grafana.com/v1/namespaces/default/datasources/12345/resource/spreadsheets
    []
    ```
