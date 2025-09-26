# Change Log

## 2.2.1

🐛 Bump `github.com/grafana/grafana-google-sdk-go` to 0.4.2

## 2.2.0

- Add support for filtering in template variables

## 2.1.0

- Add template variables functionality

## 2.0.3

- Bump `github.com/grafana/grafana-plugin-sdk-go` to 0.277.1
- Docs and configuration page updated with inline docs

## 2.0.2

- Update frontend dependencies
- Built with go 1.24

## 2.0.1

- Update documentation
- Bump `github.com/grafana/grafana-plugin-sdk-go` to 0.263.0

## 2.0.0

- Plugin now requires Grafana 10.4.8 or newer

## 1.2.18

- Backend dependencies update
- Readme changes
- Replace @grafana/experimental with @grafana/plugin-ui (#294)

## 1.2.17

- Improve data source documentation
- Fix error source for various http errors
- Bump `cross-spawn` to 7.0.6
- Bump `github.com/grafana/grafana-plugin-sdk-go` to 0.259.2

## 1.2.16

- Bump `uplot` to 1.6.31
- Bump `github.com/grafana/grafana-plugin-sdk-go` from 0.251.0 to 0.258.0

## 1.2.15

- New documentation on grafana.com #255

## 1.2.14

- Bump `path-to-regexp` from 1.8.0 to 1.9.0 #268
- Bump `github.com/grafana/grafana-plugin-sdk-go` from 0.250.0 to 0.251.0. #275

## 1.2.13

- Update `github.com/grafana/grafana-plugin-sdk-go` to `v0.248.0`
- Add logging for response size
- Fix error source for invalid JWT information

## 1.2.12

- Update `github.com/grafana/grafana-plugin-sdk-go` to `v0.245.0`
- Bump micromatch from 4.0.7 to 4.0.8
- Bump webpack from 5.92.0 to 5.94.0
- Improve efficiency when processing response
- Fix timeout error source
- Add the make docs procedure

## 1.2.11

- Update `github.com/grafana/grafana-plugin-sdk-go` to `v0.241.0`
- Fix context canceled errors to be marked as downstream

## 1.2.10

- Improve handling of unknown non-api errors
- Add error source to error responses

## 1.2.9

- Fix showing of correct percentages values
- Upgrade dependencies

## 1.2.8

- Upgrade dependencies

## 1.2.7

- Upgrade dependencies

## 1.2.6

- Build with go 1.22
- Configuration help: Add additional instruction to enable Google Sheets API

## 1.2.5

- Upgrade grafana-plugin-sdk-go to latest
- Added lint github workflow
- Remove legacy form styles

## 1.2.4

- Added feature tracking
- Upgrade dependencies

## 1.2.3

- Make sure we don't mutate the options object in the config page. This prevents crashes that occurred intermittently.

## 1.2.2

- Handle error messages more gracefully

## v1.2.0

- Refactored authentication to use grafana google sdks. With this change you can now use GCE authentication with google sheets.

There was a change in the plugin configuration. Please take a look at the provisioning example in the [documentation](src/docs/provisioning.md).
The change is backward compatible so you can still use the old configuration.

## v1.1.8

- **Chore**: Backend binaries are now compiled with golang 1.20.4

## v1.1.7

- **Chore**: Update to Golang 1.20

## v1.1.6

- Fix: Don't panic when the user selects a range of empty cells.

## v1.1.5

- **Chore**: Update to Golang 1.19 #160

## v1.1.4

- Fix: deprecated link variant button for v9

## v1.1.3

- Bump grafana dependencies to 8.3.4

## v1.1.2

- Change release pipeline

## v1.1.1

- Targeting Grafana 8.1+
- Documentation and link fixes
- Add explicit explanation to the auth constraint

## v1.1.0

- Targeting Grafana 7.2+
- Adding support for annotations
- Include arm builds

## v1.0.0

- Works with Grafana 7+
- Avoid crashing on unknown timezone (#69)
- improved support for formula values
- supports template variables

## v0.9.0

- First official release (grafana 6.7)

## v0.1.0

- Initial Release (preview)
