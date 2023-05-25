# Change Log

All notable changes to this project will be documented in this file.

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
