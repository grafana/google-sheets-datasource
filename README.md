# Google Sheets Datasource

Visualize your Google Spreadsheets with Grafana

![Visualize temperature date in Grafana Google Spreadsheets data source](./docs/img/dashboard.png)

![Average temperatures in Google Sheets](./docs/img/spreadsheet.png)

## Development

You need to install the following first:

- [Mage](https://magefile.org/)
- [Yarn](https://yarnpkg.com/)
- [Docker Compose](https://docs.docker.com/compose/)

```
mage watch
```

In another terminal

```
docker-compose up
```

To restart after backend changes:
`./scripts/restart-plugin.sh`
