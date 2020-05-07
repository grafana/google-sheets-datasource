package googlesheets

import (
	"context"
	"fmt"

	"github.com/grafana/google-sheets-datasource/pkg/models"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// GoogleClient struct
type GoogleClient struct {
	sheetsService *sheets.Service
	driveService  *drive.Service
	auth          *models.DatasourceSettings
}

type client interface {
	GetSpreadsheet(spreadSheetID string, sheetRange string, includeGridData bool) (*sheets.Spreadsheet, error)
}

// NewGoogleClient creates a new client and initializes a sheet service and a drive service
func NewGoogleClient(ctx context.Context, auth *models.DatasourceSettings) (*GoogleClient, error) {
	sheetsService, err := createSheetsService(ctx, auth)
	if err != nil {
		return nil, err
	}

	driveService, err := createDriveService(ctx, auth)
	if err != nil {
		return nil, err
	}

	return &GoogleClient{
		sheetsService: sheetsService,
		driveService:  driveService,
		auth:          auth,
	}, nil
}

// TestClient checks that the client can connect to required services
func (gc *GoogleClient) TestClient() error {
	// When using JWT, check the drive API
	if gc.auth.AuthType == "jwt" {
		q := gc.driveService.Files.List().Q("mimeType='application/vnd.google-apps.spreadsheet'")
		_, err := q.Do()
		if err != nil {
			return err
		}
	}

	// TODO: check the sheets API
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

func createSheetsService(ctx context.Context, auth *models.DatasourceSettings) (*sheets.Service, error) {
	if len(auth.AuthType) == 0 {
		return nil, fmt.Errorf("missing AuthType setting")
	}

	if auth.AuthType == "key" {
		if len(auth.APIKey) == 0 {
			return nil, fmt.Errorf("missing API Key")
		}
		return sheets.NewService(ctx, option.WithAPIKey(auth.APIKey))
	}

	if auth.AuthType == "jwt" {
		jwtConfig, err := google.JWTConfigFromJSON([]byte(auth.JWT),
			// Only need readonly access to spreadsheets (and drive?)
			"https://www.googleapis.com/auth/spreadsheets.readonly")
		if err != nil {
			return nil, fmt.Errorf("error parsing JWT file: %w", err)
		}

		client := jwtConfig.Client(ctx)
		return sheets.NewService(ctx, option.WithHTTPClient(client))
	}

	return nil, fmt.Errorf("invalid Auth Type: %s", auth.AuthType)
}

func createDriveService(ctx context.Context, auth *models.DatasourceSettings) (*drive.Service, error) {
	if len(auth.AuthType) == 0 {
		return nil, fmt.Errorf("missing AuthType setting")
	}

	if auth.AuthType == "key" {
		if len(auth.APIKey) == 0 {
			return nil, fmt.Errorf("missing API Key")
		}
		return drive.NewService(ctx, option.WithAPIKey(auth.APIKey))
	}

	if auth.AuthType == "jwt" {
		jwtConfig, err := google.JWTConfigFromJSON([]byte(auth.JWT), drive.DriveMetadataReadonlyScope)
		if err != nil {
			return nil, fmt.Errorf("error parsing JWT file: %w", err)
		}

		client := jwtConfig.Client(ctx)
		return drive.NewService(ctx, option.WithHTTPClient(client))
	}
	return nil, fmt.Errorf("invalid Auth Type: %s", auth.AuthType)
}
