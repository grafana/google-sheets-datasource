package googlesheets

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/grafana/google-sheets-datasource/pkg/models"
	"github.com/grafana/grafana-google-sdk-go/pkg/tokenprovider"
	"github.com/grafana/grafana-plugin-sdk-go/backend/httpclient"
	"github.com/pkg/errors"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

const (
	sheetsRoute = "sheets"
	driveRoute  = "drive"
)

type routeInfo struct {
	method string
	scopes []string
}

var routes = map[string]routeInfo{
	sheetsRoute: {
		method: "GET",
		scopes: []string{sheets.SpreadsheetsReadonlyScope},
	},
	driveRoute: {
		method: "GET",
		scopes: []string{drive.DriveReadonlyScope},
	},
}

// GoogleClient struct
type GoogleClient struct {
	sheetsService *sheets.Service
	driveService  *drive.Service
	auth          string
}

type client interface {
	GetSpreadsheet(spreadSheetID string, sheetRange string, includeGridData bool) (*sheets.Spreadsheet, error)
}

// NewGoogleClient creates a new client and initializes a sheet service and a drive service
func NewGoogleClient(ctx context.Context, settings models.DatasourceSettings) (*GoogleClient, error) {
	sheetsService, err := createSheetsService(ctx, settings)
	if err != nil {
		return nil, err
	}

	driveService, err := createDriveService(ctx, settings)
	if err != nil {
		return nil, err
	}

	return &GoogleClient{
		sheetsService: sheetsService,
		driveService:  driveService,
		auth:          settings.AuthenticationType,
	}, nil
}

// TestClient checks that the client can connect to required services
func (gc *GoogleClient) TestClient() error {
	// When using JWT, check the drive API
	if gc.auth == "jwt" {
		_, err := gc.driveService.Files.List().PageSize(1).Do()
		if err != nil {
			return err
		}
	}

	// Test spreadsheet from google
	spreadsheetId := "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"
	readRange := "Class Data!A2:E"
	_, err := gc.sheetsService.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		return err
	}
	return nil
}

// GetSpreadsheet gets a google spreadsheet struct by id and range
func (gc *GoogleClient) GetSpreadsheet(spreadSheetID string, sheetRange string, includeGridData bool) (*sheets.Spreadsheet, error) {
	req := gc.sheetsService.Spreadsheets.Get(spreadSheetID)
	if len(sheetRange) > 0 {
		req = req.Ranges(sheetRange)
	}
	return req.IncludeGridData(true).Do()
}

// GetSpreadsheetFiles lists all files with spreadsheet mimetype that the client has access to.
func (gc *GoogleClient) GetSpreadsheetFiles() ([]*drive.File, error) {
	fs := []*drive.File{}
	pageToken := ""
	for {
		q := gc.driveService.Files.List().Q("mimeType='application/vnd.google-apps.spreadsheet'")
		if pageToken != "" {
			q = q.PageToken(pageToken)
		}
		r, err := q.Do()
		if err != nil {
			return nil, fmt.Errorf("failed to list spreadsheet files, page token %q: %w", pageToken, err)
		}

		fs = append(fs, r.Files...)
		pageToken = r.NextPageToken
		if pageToken == "" {
			break
		}
	}

	return fs, nil
}

func createSheetsService(ctx context.Context, settings models.DatasourceSettings) (*sheets.Service, error) {
	if len(settings.AuthenticationType) == 0 {
		return nil, fmt.Errorf("missing AuthenticationType setting")
	}

	if settings.AuthenticationType == "api" {
		if len(settings.APIKey) == 0 {
			return nil, fmt.Errorf("missing API Key")
		}
		return sheets.NewService(ctx, option.WithAPIKey(settings.APIKey))
	}

	client, err := newHTTPClient(settings, httpclient.Options{}, sheetsRoute)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to create http client")
	}

	srv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}

	return srv, nil
}

func createDriveService(ctx context.Context, settings models.DatasourceSettings) (*drive.Service, error) {
	if len(settings.AuthenticationType) == 0 {
		return nil, fmt.Errorf("missing AuthenticationType setting")
	}

	if settings.AuthenticationType == "api" {
		if len(settings.APIKey) == 0 {
			return nil, fmt.Errorf("missing API Key")
		}
		return drive.NewService(ctx, option.WithAPIKey(settings.APIKey))
	}

	client, err := newHTTPClient(settings, httpclient.Options{}, driveRoute)
	if err != nil {
		return nil, errors.WithMessage(err, "Failed to create http client")
	}

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	return srv, nil
}

func getMiddleware(settings models.DatasourceSettings, routePath string) (httpclient.Middleware, error) {
	providerConfig := tokenprovider.Config{
		RoutePath:         routePath,
		RouteMethod:       routes[routePath].method,
		DataSourceID:      settings.InstanceSettings.ID,
		DataSourceUpdated: settings.InstanceSettings.Updated,
		Scopes:            routes[routePath].scopes,
	}

	var provider tokenprovider.TokenProvider
	switch settings.AuthenticationType {
	case "gce":
		provider = tokenprovider.NewGceAccessTokenProvider(providerConfig)
	case "jwt":
		err := validateDataSourceSettings(settings)

		if err != nil {
			return nil, err
		}

		providerConfig.JwtTokenConfig = &tokenprovider.JwtTokenConfig{
			Email:      settings.ClientEmail,
			URI:        settings.TokenUri,
			PrivateKey: []byte(settings.PrivateKey),
		}
		provider = tokenprovider.NewJwtAccessTokenProvider(providerConfig)
	}

	return tokenprovider.AuthMiddleware(provider), nil
}

func newHTTPClient(settings models.DatasourceSettings, opts httpclient.Options, route string) (*http.Client, error) {
	m, err := getMiddleware(settings, route)
	if err != nil {
		return nil, err
	}

	opts.Middlewares = append(opts.Middlewares, m)
	return httpclient.New(opts)
}

func validateDataSourceSettings(settings models.DatasourceSettings) error {
	if settings.DefaultProject == "" || settings.ClientEmail == "" || settings.PrivateKey == "" || settings.TokenUri == "" {
		return fmt.Errorf("datasource is missing authentication details")
	}

	return nil
}
