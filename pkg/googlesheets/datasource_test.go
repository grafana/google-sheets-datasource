package googlesheets

import (
	"context"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/backend"
	"github.com/grafana/grafana-plugin-sdk-go/backend/datasource"
	"github.com/grafana/grafana-plugin-sdk-go/backend/instancemgmt"
	"github.com/grafana/grafana-plugin-sdk-go/experimental/datasourcetest"
	"github.com/patrickmn/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestGoogleSheetsDatasource_CheckHealth(t *testing.T) {
	ds := &Datasource{
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
	ds := &Datasource{
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

func TestGoogleSheetsMultiTenancy(t *testing.T) {
	const (
		tenantID1 = "abc123"
		tenantID2 = "def456"
		addr      = "127.0.0.1:8000"
	)

	var instances []instancemgmt.Instance
	factoryInvocations := 0
	factory := datasource.InstanceFactoryFunc(func(ctx context.Context, settings backend.DataSourceInstanceSettings) (instancemgmt.Instance, error) {
		factoryInvocations++
		i, err := NewDatasource(ctx, settings)
		if err == nil {
			instances = append(instances, i)
		}
		return i, err
	})

	tp, err := datasourcetest.Manage(factory, datasourcetest.ManageOpts{Address: addr})
	require.NoError(t, err)
	defer func() {
		err = tp.Shutdown()
		if err != nil {
			t.Log("plugin shutdown error", err)
		}
	}()

	pCtx := backend.PluginContext{DataSourceInstanceSettings: &backend.DataSourceInstanceSettings{
		ID: 1,
		DecryptedSecureJSONData: map[string]string{
			"privateKey": "randomPrivateKey",
		},
		JSONData: []byte(`{"authenticationType":"jwt","defaultProject": "raintank-dev", "processingLocation": "us-west1","tokenUri":"token","clientEmail":"test@grafana.com"}`),
	}}

	t.Run("Request without tenant information creates an instance", func(t *testing.T) {
		qdr := &backend.QueryDataRequest{PluginContext: pCtx}
		crr := &backend.CallResourceRequest{PluginContext: pCtx}
		chr := &backend.CheckHealthRequest{PluginContext: pCtx}
		responseSender := newTestCallResourceResponseSender()
		ctx := context.Background()

		resp, err := tp.Client.QueryData(ctx, qdr)
		require.NoError(t, err)
		require.NotNil(t, resp)
		require.Equal(t, 1, factoryInvocations)

		err = tp.Client.CallResource(ctx, crr, responseSender)
		require.NoError(t, err)
		require.Equal(t, 1, factoryInvocations)

		t.Run("Request from tenant #1 creates new instance", func(t *testing.T) {
			ctx = metadata.AppendToOutgoingContext(context.Background(), "tenantID", tenantID1)
			resp, err = tp.Client.QueryData(ctx, qdr)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, 2, factoryInvocations)

			// subsequent requests from tenantID1 with same settings will reuse instance
			resp, err = tp.Client.QueryData(ctx, qdr)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, 2, factoryInvocations)

			var chRes *backend.CheckHealthResult
			chRes, err = tp.Client.CheckHealth(ctx, chr)
			require.NoError(t, err)
			require.NotNil(t, chRes)
			require.Equal(t, 2, factoryInvocations)

			t.Run("Request from tenant #2 creates new instance", func(t *testing.T) {
				ctx = metadata.AppendToOutgoingContext(context.Background(), "tenantID", tenantID2)
				resp, err = tp.Client.QueryData(ctx, qdr)
				require.NoError(t, err)
				require.NotNil(t, resp)
				require.Equal(t, 3, factoryInvocations)

				// subsequent requests from tenantID2 with same settings will reuse instance
				err = tp.Client.CallResource(ctx, crr, responseSender)
				require.NoError(t, err)
				require.Equal(t, 3, factoryInvocations)
			})

			// subsequent requests from tenantID1 with same settings will reuse instance
			ctx = metadata.AppendToOutgoingContext(context.Background(), "tenantID", tenantID1)
			resp, err = tp.Client.QueryData(ctx, qdr)
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.Equal(t, 3, factoryInvocations)

			chRes, err = tp.Client.CheckHealth(ctx, chr)
			require.NoError(t, err)
			require.NotNil(t, chRes)
			require.Equal(t, 3, factoryInvocations)
		})
	})

	require.Len(t, instances, 3)
	require.NotEqual(t, instances[0], instances[1])
	require.NotEqual(t, instances[0], instances[2])
	require.NotEqual(t, instances[1], instances[2])
}

type testCallResourceResponseSender struct{}

func newTestCallResourceResponseSender() *testCallResourceResponseSender {
	return &testCallResourceResponseSender{}
}

func (s *testCallResourceResponseSender) Send(_ *backend.CallResourceResponse) error {
	return nil
}
