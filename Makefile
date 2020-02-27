DSNAME=sheets-datasource
GO = GO111MODULE=on go
all: build-frontend build


dev: build watch-frontend

build:
	GOOS=linux $(GO) build -o ./dist/${DSNAME}_linux_amd64 -tags netgo -ldflags '-w' ./pkg

build-debug:
	GOOS=linux $(GO) build -o ./dist/${DSNAME}_linux_amd64 -gcflags=all="-N -l" ./pkg

build-darwin:
	$(GO) build -o ./dist/${DSNAME}_darwin_amd64 -tags netgo -ldflags '-w' ./pkg

build-debug-darwin:
	$(GO) build -o ./dist/${DSNAME}_darwin_amd64 -gcflags=all="-N -l" ./pkg

# Note frontend deletes backend file
build-frontend:
	npx grafana-toolkit plugin:build

watch-frontend: 
	npx grafana-toolkit plugin:dev --watch
