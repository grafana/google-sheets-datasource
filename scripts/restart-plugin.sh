#!/bin/bash

make build
docker exec google-sheets-datasource_grafana_1 pkill -f "/var/lib/grafana/plugins/google_sheets_datasource/dist/sheets_datasource_linux_amd64"
