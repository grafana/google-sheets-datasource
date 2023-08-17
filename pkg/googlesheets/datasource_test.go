package googlesheets

import (
	"context"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
)

func TestGoogleSheetsDatasource_CheckHealth(t *testing.T) {
	ds := &GoogleSheetsDatasource{
		googlesheets: &GoogleSheets{},
	}

	t.Run("should return error when unable to load settings", func(t *testing.T) {
		req := &backend.CheckHealthRequest{
			PluginContext: backend.PluginContext{
				DataSourceInstanceSettings: &backend.DataSourceInstanceSettings{
					JSONData: nil,
				},
			},
		}

		res, err := ds.CheckHealth(context.Background(), req)

		assert.Equal(t, backend.HealthStatusError, res.Status)
		assert.Equal(t, "Unable to load settings", res.Message)
		assert.Nil(t, err)
	})

	t.Run("should return error when unable to create client", func(t *testing.T) {
		req := &backend.CheckHealthRequest{
			PluginContext: backend.PluginContext{
				DataSourceInstanceSettings: &backend.DataSourceInstanceSettings{
					JSONData: []byte(`{"client_id":"test","client_secret":"test","token":"test"}`),
				},
			},
		}

		res, err := ds.CheckHealth(context.Background(), req)

		assert.Equal(t, backend.HealthStatusError, res.Status)
		assert.Equal(t, "Unable to create client", res.Message)
		assert.Nil(t, err)
	})

	t.Run("should return error when permissions check failed", func(t *testing.T) {
		req := &backend.CheckHealthRequest{
			PluginContext: backend.PluginContext{
				DataSourceInstanceSettings: &backend.DataSourceInstanceSettings{
					JSONData:                []byte(`{"authenticationType":"key"}`),
					DecryptedSecureJSONData: map[string]string{"apiKey": "token"},
				},
			},
		}

		ds.googlesheets = &GoogleSheets{
			Cache: cache.New(300*time.Second, 5*time.Second),
		}

		res, err := ds.CheckHealth(context.Background(), req)

		assert.Equal(t, backend.HealthStatusError, res.Status)
		assert.Equal(t, "Permissions check failed", res.Message)
		assert.Nil(t, err)
	})
}

func TestGoogleSheetsDatasource_QueryData(t *testing.T) {
	ds := &GoogleSheetsDatasource{
		googlesheets: &GoogleSheets{},
	}

	t.Run("should return error when unable to load settings", func(t *testing.T) {
		req := &backend.QueryDataRequest{
			PluginContext: backend.PluginContext{
				DataSourceInstanceSettings: &backend.DataSourceInstanceSettings{
					JSONData: nil,
				},
			},
		}

		res, err := ds.QueryData(context.Background(), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})

	t.Run("should return error when failed to read query", func(t *testing.T) {
		req := &backend.QueryDataRequest{
			PluginContext: backend.PluginContext{
				DataSourceInstanceSettings: &backend.DataSourceInstanceSettings{
					JSONData: []byte(`{"client_id":"test","client_secret":"test","token":"test"}`),
				},
			},
			Queries: []backend.DataQuery{
				{
					RefID: "test",
				},
			},
		}

		res, err := ds.QueryData(context.Background(), req)

		assert.Nil(t, res)
		assert.NotNil(t, err)
	})
}
