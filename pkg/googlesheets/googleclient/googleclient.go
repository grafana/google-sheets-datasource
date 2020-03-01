package googleclient

import (
	"context"
	"fmt"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// GoogleClient struct
type GoogleClient struct {
	sheetsService *sheets.Service
	driveService  *drive.Service
}

// Auth struct
type Auth struct {
	APIKey   string
	AuthType string
	JWT      string
}

// NewAuth creates a new auth struct
func NewAuth(apiKey string, authType string, jwt string) *Auth {
	return &Auth{
		APIKey:   apiKey,
		AuthType: authType,
		JWT:      jwt,
	}
}

// New creates a new client and initializes a sheet service and a drive service
func New(ctx context.Context, auth *Auth) (*GoogleClient, error) {
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
	}, nil
}

// GetSpreadsheet gets a google spreadsheet struct by id and range
func (gc *GoogleClient) GetSpreadsheet(spreadSheetID string, sheetRange string, includeGridData bool) (*sheets.Spreadsheet, error) {
	return gc.sheetsService.Spreadsheets.Get(spreadSheetID).Ranges(sheetRange).IncludeGridData(true).Do()
}

// GetSpreadsheetFiles lists all files with spreadsheet mimetype that the service account that was used to initialize the client has access to
func (gc *GoogleClient) GetSpreadsheetFiles() ([]*drive.File, error) {
	var fs []*drive.File
	pageToken := ""
	for {
		q := gc.driveService.Files.List().Q("mimeType='application/vnd.google-apps.spreadsheet'")
		if pageToken != "" {
			q = q.PageToken(pageToken)
		}
		r, err := q.Do()
		if err != nil {
			return fs, err
		}
		fs = append(fs, r.Files...)
		pageToken = r.NextPageToken
		if pageToken == "" {
			break
		}
	}
	return fs, nil
}

func createSheetsService(ctx context.Context, auth *Auth) (*sheets.Service, error) {
	if len(auth.AuthType) == 0 {
		return nil, fmt.Errorf("missing AuthType setting")
	}

	if auth.AuthType == "key" {
		if len(auth.APIKey) == 0 {
			return nil, fmt.Errorf("Missing API Key")
		}
		return sheets.NewService(ctx, option.WithAPIKey(auth.APIKey))
	}

	if auth.APIKey == "jwt" {
		jwtConfig, err := google.JWTConfigFromJSON([]byte(auth.JWT), "https://www.googleapis.com/auth/spreadsheets.readonly", "https://www.googleapis.com/auth/spreadsheets")
		if err != nil {
			return nil, fmt.Errorf("error parsing JWT file: %w", err)
		}

		client := jwtConfig.Client(ctx)
		return sheets.NewService(ctx, option.WithHTTPClient(client))
	}

	return nil, fmt.Errorf("invalid Auth Type: %s", auth.AuthType)
}

func createDriveService(ctx context.Context, auth *Auth) (*drive.Service, error) {
	if len(auth.AuthType) == 0 {
		return nil, fmt.Errorf("Missing AuthType setting")
	}

	if auth.AuthType == "key" {
		if len(auth.APIKey) == 0 {
			return nil, fmt.Errorf("Missing API Key")
		}
		return drive.NewService(ctx, option.WithAPIKey(auth.APIKey))
	}

	if auth.APIKey == "jwt" {
		jwtConfig, err := google.JWTConfigFromJSON([]byte(auth.JWT), drive.DriveMetadataReadonlyScope)
		if err != nil {
			return nil, fmt.Errorf("error parsing JWT file: %w", err)
		}

		client := jwtConfig.Client(ctx)
		return drive.NewService(ctx, option.WithHTTPClient(client))
	}
	return nil, fmt.Errorf("Invalid Auth Type: %s", auth.AuthType)
}
