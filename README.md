# Google Sheets Datasource

Visualize your Google Spreadsheets with Grafana

![Visualize temperature date in Grafana Google Spreadsheets data source](./src/docs/img/dashboard.png =600x)

![Average temperatures in Google Sheets](./src/docs/img/spreadsheet.png =600x)

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
