---
title: Configure
menuTitle: Configure 
description: Configure the Google Sheets data source plugin
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
weight: 103
---

# Configure the Google Sheets Plugin

To configure the Google Sheets plugin, you need to perform the following steps:

1. Navigate into Grafana and click on the menu option on the top left.
1. Browse to the **Connections** menu and then click on the **Data sources**.
1. Click on the Google Sheets data source plugin which you have installed.
1. Go to its settings tab and at the bottom, you will find the **Authentication** section.
1. You will be provided with three options to authenticated it:
   1. API Key
   1. Google JWT File
   1. GCE Default Service Account
1. You can select any one from above that you have created from the [Authentication](./setup/authentication.md) page.
1. Click **Save & Test** button and you should see a confirmation dialog box that says “Data source is working”.

{{< admonition type="tip" >}}
If you see errors, check the Grafana logs for troubleshooting.
{{< /admonition >}}
