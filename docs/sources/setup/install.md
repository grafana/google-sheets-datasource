---
title: Install the Google Sheets data source plugin for Grafana
menuTitle: Install
description: Learn how to install the Google Sheets data source plugin for Grafana
keywords:
  - data source
  - google sheets
  - spreadsheets
  - xls data
  - xlsx data
  - excel sheets
  - excel data
  - csv data
  - visualize spreadsheets
labels:
  products:
    - oss
    - enterprise
    - cloud
weight: 101
---

# Install the Google Sheets data source plugin for Grafana

You can any of the following sets of steps to install the Google Sheets data source plugin for Grafana.

## Install from plugin catalog

To install the plugin from the plugin catalog:

1. Sign in to Grafana as a server administrator.
1. Click **Administration** > **Plugins and data** > **Plugins** in the side navigation menu to view all plugins.
1. Type **Google Sheets** in the Search box.
1. Click the **All** in the **State** filter option.
1. Click the plugin logo.
1. Click **Install**.

## Install from the Grafana plugins page

To install the plugin from the Grafana plugins page, browse to the [Google Sheets data source plugin](https://grafana.com/grafana/plugins/grafana-googlesheets-datasource/?tab=installation) and follow the instructions provided there.

## Install from GitHub

To install the plugin from GitHub:

1. Browse to the [plugin GitHub releases page](https://github.com/grafana/google-sheets-datasource/releases).

1. Find the release you want to install.

1. Download the release by clicking the release asset called `grafana-googlesheets-datasource-<VERSION>.zip` where _`VERSION`_ is the version of the plugin you want to install.
   You may need to un-collapse the **Assets** section to see it.

1. Extract the plugin into the Grafana plugins directory:

   On Linux or macOS, run the following commands to extract the plugin:

   ```bash
   unzip grafana-googlesheets-datasource-<VERSION>.zip
   mv grafana-googlesheets-datasource /var/lib/grafana/plugins
   ```

   On Windows, run the following command to extract the plugin:

   ```powershell
   Expand-Archive -Path grafana-googlesheets-datasource-<VERSION>.zip -DestinationPath C:\grafana\data\plugins
   ```

1. Restart Grafana.

## Install using `grafana-cli`

To install the plugin using `grafana-cli`:

1. On Linux or macOS, open your terminal and run the following command:

   ```bash
   grafana-cli plugins install grafana-googlesheets-datasource
   ```

   On Windows, run the following command:

   ```shell
   grafana-cli.exe plugins install grafana-googlesheets-datasource
   ```

1. Restart Grafana.

### Install a custom version

If you need to install a custom version of the plugin using `grafana-cli`, use the `--pluginUrl` option.

```bash
grafana-cli --pluginUrl <ZIP_FILE_URL> plugins install grafana-googlesheets-datasource
```

For example, to install version `1.2.10` of the plugin on Linux or macOS:

```bash
grafana-cli --pluginUrl https://github.com/grafana/google-sheets-datasource/releases/download/v1.2.10/grafana-googlesheets-datasource-1.2.10.zip plugins install grafana-googlesheets-datasource
```

Or to install version `1.2.10` of the plugin on Windows:

```powershell
grafana-cli.exe --pluginUrl https://github.com/grafana/google-sheets-datasource/releases/download/v1.2.10/grafana-googlesheets-datasource-1.2.10.zip plugins install grafana-googlesheets-datasource
```
