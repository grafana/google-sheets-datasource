#!/bin/bash
set -eo pipefail

make build
docker exec google-sheets-datasource_grafana_1 pkill -f "/var/lib/grafana/plugins/google-sheets-datasource/dist/sheets-datasource_linux_amd64"
