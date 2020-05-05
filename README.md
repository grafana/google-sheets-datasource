# Google Sheets data source

Visualize your Google Spreadsheets with Grafana

![Visualize temperature date in Grafana Google Spreadsheets data source](https://raw.githubusercontent.com/grafana/google-sheets-datasource/master/src/docs/img/dashboard.png)

![Average temperatures in Google Sheets](https://raw.githubusercontent.com/grafana/google-sheets-datasource/master/src/docs/img/spreadsheet.png)

## Documentation

Check the [docs](https://github.com/grafana/google-sheets-datasource/blob/master/src/README.md) for information on how to use the data source.

## Development

You need to install the following first:

- [Mage](https://magefile.org/)
- [Yarn](https://yarnpkg.com/)
- [Docker Compose](https://docs.docker.com/compose/)


### Building the Plug-In
In order to build the plug-in, both front-end and back-end parts, do the following:

```
yarn install
yarn build
```

The built plug-in will be in the dist/ directory.

### Testing the Plug-In w/ Docker Compose
To test the plug-in running inside Grafana, we recommend using our Docker Compose setup:

```BASH
mage buildAll
```

In another terminal

```BASH
docker-compose up
```

To restart the plug-in after backend changes:
`./scripts/restart-plugin.sh`

### Test spreadsheet

Need at publicly available spreadsheet to play around with? Feel free to use [this](https://docs.google.com/spreadsheets/d/1TZlZX67Y0s4CvRro_3pCYqRCKuXer81oFp_xcsjPpe8/edit?usp=sharing) demo spreadsheet that is suitable for visualization in graphs and in tables.
