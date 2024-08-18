---
title: Provisioning the Google Sheets data source in Grafana
menuTitle: Provisioning
description: Provisioning the Google Sheets source plugin
keywords:
  - data source
  - google sheets
  - spreadsheets
  - xls data
  - xlsx data
  - excel sheets
  - excel data
  - csv data
labels:
  products:
    - oss
    - enterprise
    - cloud
weight: 104
---

# Provisioning the Google Sheets data source in Grafana

You can define and configure the Google Sheets data source in YAML files with Grafana provisioning. For more information about provisioning a data source, and for available configuration options, refer to [Provision Grafana](https://grafana.com/docs/grafana/latest/administration/provisioning/#data-sources).

**Example**

The following YAML snippet can be used to provision the Google Sheets data source for Grafana if you are using [Prometheus Operator](https://github.com/prometheus-operator/prometheus-operator):

```yaml
promop:
  grafana:
    additionalDataSources:
      - name: GitHub Repo Insights
        type: grafana-github-datasource
        jsonData:
          owner: ’<REPOSITORY_OWNER>’
          repository: ’<REPOSITORY_NAME>’
        secureJsonData:
          accessToken: ’<PERSONAL_ACCESS_TOKEN>’
```
