DSNAME=test-datasource
GO = GO111MODULE=on go
all: build-frontend build

# TODO: This should build for the current arch, not linux
build:
	$(GO) build -o ./dist/${DSNAME}_linux_amd64 -tags netgo -ldflags '-w' ./datasource

build-darwin:
	$(GO) build -o ./dist/${DSNAME}_darwin_amd64 -tags netgo -ldflags '-w' ./datasource

# Note frontend deletes backend file
build-frontend:
	npx grafana-toolkit plugin:build