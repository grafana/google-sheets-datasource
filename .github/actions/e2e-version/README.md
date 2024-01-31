# E2E Grafana version resolver

This Action resolves what versions of Grafana to use when E2E testing a Grafana plugin in a Github Action.

## Features

The action supports two modes.

**plugin-grafana-dependency (default)**
The will return all the latest patch release of every minor version of Grafana since the version that was specified as grafanaDependency in the plugin.json. This requires the plugin.json file to be placed in the `<root>/src` directory.

### Example

At the time of writing, the most recent release of Grafana is 10.2.2. If the plugin has specified >=8.0.0 as `grafanaDependency` in the plugin.json file, the output would be:

```json
[
  "10.2.2",
  "10.1.5",
  "10.0.9",
  "9.5.14",
  "9.4.17",
  "9.3.16",
  "9.2.20",
  "9.1.8",
  "9.0.8",
  "8.5.27",
  "8.4.11",
  "8.3.11",
  "8.2.7",
  "8.1.8",
  "8.0.7"
]
```

Please note that the output changes as new versions of Grafana are being released.

**version-support-policy**
This will resolve versions according to Grafana's plugin compatibility support policy. Specifically, it retrieves the latest patch release for each minor version within the current major version of Grafana. Additionally, it includes the most recent release for the latest minor version of the previous major Grafana version.```

### Example

At the time of writing, the most recent release of Grafana is 10.2.2. The output for `version-support-policy` would be:

```json
["10.2.2", "10.1.5", "10.0.9", "9.5.14"]
```

### Output

The result of this action is a JSON array that lists the latest patch version for each Grafana minor version. These values can be employed to define a version matrix in a subsequent workflow job.

## Workflow example

### plugin-grafana-dependency

```yaml
name: E2E tests - Playwright
on:
  pull_request:

jobs:
  setup-matrix:
    runs-on: ubuntu-latest
    outputs:
      matrix: ${{ steps.resolve-versions.outputs.matrix }}
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      - name: Resolve Grafana E2E versions
        id: resolve-versions
        uses: grafana/plugin-actions/e2e-version
        with:
          # target all minor versions of Grafana that have been released since the version that was specified as grafanaDependency in the plugin
          version-resolver-type: plugin-grafana-dependency

  playwright-tests:
    needs: setup-matrix
    strategy:
      matrix:
        # use matrix from output in previous job
        GRAFANA_VERSION: ${{fromJson(needs.setup-matrix.outputs.matrix)}}
    runs-on: ubuntu-latest
    steps:
      ...
      - name: Start Grafana
        run: docker run --rm -d -p 3000:3000 --name=grafana grafana/grafana:${{ matrix.GRAFANA_VERSION }}
      ...
```

### version-support-policy

```yaml
name: E2E tests - Playwright
on:
  pull_request:

jobs:
  setup-matrix:
    runs-on: ubuntu-latest
    outputs:
      matrix: ${{ steps.resolve-versions.outputs.matrix }}
    steps:
      - name: Checkout
        uses: actions/checkout@v4
      - name: Resolve Grafana E2E versions
        id: resolve-versions
        uses: grafana/plugin-actions/e2e-version
        with:
          #target all minors for the current major version of Grafana and the last minor of the previous major version of Grafana
          version-resolver-type: version-support-policy

  playwright-tests:
    needs: setup-matrix
    strategy:
      matrix:
        # use matrix from output in previous job
        GRAFANA_VERSION: ${{fromJson(needs.setup-matrix.outputs.matrix)}}
    runs-on: ubuntu-latest
    steps:
      ...
      - name: Start Grafana
        run: docker run --rm -d -p 3000:3000 --name=grafana grafana/grafana:${{ matrix.GRAFANA_VERSION }}
      ...
```

## Development

```bash
cd e2e-versions
npm i

#before pushing to main
npm run bundle
```
