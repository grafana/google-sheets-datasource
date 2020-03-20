# Google Sheets Datasource

Visualize your Google Spreadsheets with Grafana

![Visualize temperature date in Grafana Google Spreadsheets data source](./src/docs/img/dashboard.png)

![Average temperatures in Google Sheets](./src/docs/img/spreadsheet.png)

## Development

You need to install the following first:

- [Mage](https://magefile.org/)
- [Yarn](https://yarnpkg.com/)
- [Docker Compose](https://docs.docker.com/compose/)

```BASH
mage watch
```

In another terminal

```BASH
docker-compose up
```

To restart after backend changes:
`./scripts/restart-plugin.sh`
