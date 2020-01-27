#!/bin/bash

make build
cid=$(docker ps -aqf "name=google-sheets-datasource_grafana_1")
docker exec -u 0 -it $cid bash > pkill -f "/var/lib/grafana/plugins/google_sheets_datasource/dist/sheets_datasource_linux_amd64"

