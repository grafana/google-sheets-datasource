# Google Sheets data source

Visualize your Google Spreadsheets with Grafana

![Visualize temperature date in Grafana Google Spreadsheets data source](https://github.com/user-attachments/assets/7857ac8e-835d-401e-b51d-3daf3d4aa89a)

![Average temperatures in Google Sheets](https://github.com/user-attachments/assets/218e3346-f968-495b-85ae-d29688516bba)

## Documentation

For the plugin documentation, visit [plugin documentation website](https://grafana.com/docs/plugins/grafana-googlesheets-datasource/).

## Video Tutorial

Watch this video to learn more about setting up the Grafana Google Sheets data source plugin:

[![Google Sheets data source plugin | Visualize Spreadsheets using Grafana | Tutorial](https://img.youtube.com/vi/hqeqeQFrtSA/hq720.jpg)](https://youtu.be/hqeqeQFrtSA "Grafana Google Sheets data source plugin")

> [!TIP]
> 
> ## Give it a try using Grafana Play
> 
> With Grafana Play, you can explore and see how it works, learning from practical examples to accelerate your development. This feature can be seen on [Google Sheets data source plugin demo](https://play.grafana.org/d/ddkar8yanj56oa/visualizing-google-sheets-data).

## Development guide

This is a basic guide on how to set up your local environment, make the desired changes and see the result with a fresh Grafana installation.

## Requirements

You need to install the following first:

- [Mage](https://magefile.org/)
- [Yarn](https://yarnpkg.com/)
- [Docker Compose](https://docs.docker.com/compose/)

## Running the development version

### Compiling the Backend

If you have made any changes to any `go` files, you can use [mage](https://github.com/magefile/mage) to recompile the plugin.

```sh
mage buildAll
```

### Compiling the Frontend

After you made the desired changes, you can build and test the new version of the plugin using `yarn`:

```sh
yarn run dev # builds and puts the output at ./dist
```

The built plug-in will be in the `dist/` directory.

### Docker Compose

To test the plug-in running inside Grafana, we recommend using our Docker Compose setup:

```sh
docker-compose up
```

### Test spreadsheet

Need at publicly available spreadsheet to play around with? Feel free to use [this](https://docs.google.com/spreadsheets/d/1TZlZX67Y0s4CvRro_3pCYqRCKuXer81oFp_xcsjPpe8/edit?usp=sharing) demo spreadsheet that is suitable for visualization in graphs and in tables.
