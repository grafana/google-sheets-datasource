---
title: Configure the Google Sheets data source plugin
menuTitle: Configure
description: How to configure the Google Sheets data source plugin
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

# Configure the Google Sheets data source plugin

To configure the Google Sheets data source plugin, you need to perform the following steps:

1. Navigate into Grafana and click on the menu option on the top left.
1. Browse to the **Connections** menu and then click on the **Data sources**.
1. Click on the **Add new data source** button.
1. Click on the Google Sheets data source plugin which you have installed.
1. Go to its settings tab and find the **Authentication** section.
1. It supports the following three ways of authentication:

   {{< tabs >}}

   {{< tab-content name="with an API key" >}}
   Before you begin, [create an API key](../authenticate/#authenticate-with-an-api-key).

   1. Select the **API Key** option in the **Authentication type**.
   1. Paste the API key.
   1. Click **Save & Test** button and you should see a confirmation dialog box that says "Data source is working".

   {{< /tab-content >}}

   {{< tab-content name="with a service account JWT" >}}
   Before you begin, [create a service account and download the JWT file](../authenticate/#authenticate-with-a-service-account-jwt).

   1. Select the **Google JWT File** option in the **Authentication type**.

   1. You can perform one of the following three options:

      1. Upload the Google JWT file by clicking the **Click to browse files** and select the JSON file you downloaded.
      1. Click the **Paste JWT Token** button and paste the complete JWT token manually
      1. Click the **Fill In JWT Token manually** button and provide the JWT details including Project ID, Client email, Token URI, and Private key.

   1. Click **Save & Test** button and you should see a confirmation dialog box that says "Data source is working".

   {{< /tab-content >}}

   {{< tab-content name="with the default GCE service account" >}}

   Before you begin, set up [authentication with the default GCE service account](../authenticate/#authenticate-with-the-default-gce-service-account)

   1. Select the **GCE Default Service Account** option in the **Authentication type**.
   1. Type the **Default project** name
   1. Click **Save & Test** button and you should see a confirmation dialog box that says "Data source is working".

   {{< admonition type="tip" >}}
   If you see errors, check the Grafana logs for troubleshooting.
   {{< /admonition >}}

   {{< /tab-content >}}
   {{< /tabs >}}
