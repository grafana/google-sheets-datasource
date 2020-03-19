# Google Sheets Datasource

Visualize your Google Spreadsheets with Grafana

![Visualize temperature date in Grafana Google Spreadsheets data source](./src/docs/img/dashboard.png)

![Average temperatures in Google Sheets](./src/docs/img/spreadsheet.png)

## Documentation

Check the [docs](https://github.com/grafana/google-sheets-datasource/blob/master/src/README.md) for information on how to use the data source.

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
