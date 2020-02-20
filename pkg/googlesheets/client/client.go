package client

import (
	"context"
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

// Client struct
type Client struct {
	sheetsService *sheets.Service
	driveService  *drive.Service
}

// Auth struct
type Auth struct {
	APIKey   string
	AuthType string
	JWT      json.RawMessage
}

// NewAuth creates a new auth struct
func NewAuth(apiKey string, authType string, jwt json.RawMessage) *Auth {
	return &Auth{
		APIKey:   apiKey,
		AuthType: authType,
		JWT:      jwt,
	}
}

// New creates a new client and initializes a sheet service and a drive service
func New(ctx context.Context, auth *Auth) (*Client, error) {
	sheetsService, err := createSheetsService(ctx, auth)
	if err != nil {
		return nil, err
	}

	driveService, err := createDriveService(ctx, auth)
	if err != nil {
		return nil, err
	}

	return &Client{
		sheetsService: sheetsService,
		driveService:  driveService,
	}, nil
}

// GetSpreadsheet gets a google spreadsheet struct by id and range
func (c *Client) GetSpreadsheet(spreadSheetID string, sheetRange string, includeGridData bool) (*sheets.Spreadsheet, error) {
	return c.sheetsService.Spreadsheets.Get(spreadSheetID).Ranges(sheetRange).IncludeGridData(true).Do()
}

// GetSpreadsheetFiles lists all files with spreadsheet mimetype that the service account that was used to initialize the client has access to
func (c *Client) GetSpreadsheetFiles() ([]*drive.File, error) {
	var fs []*drive.File
	pageToken := ""
	for {
		q := c.driveService.Files.List().Q("mimeType='application/vnd.google-apps.spreadsheet'")
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
	if auth.AuthType == "none" {
		if len(auth.APIKey) == 0 {
			return nil, fmt.Errorf("Invalid API Key")
		}

		return sheets.NewService(ctx, option.WithAPIKey(auth.APIKey))
	}

	jwtConfig, err := google.JWTConfigFromJSON([]byte(auth.JWT), "https://www.googleapis.com/auth/spreadsheets.readonly", "https://www.googleapis.com/auth/spreadsheets")
	if err != nil {
		return nil, fmt.Errorf("Error parsin JWT file: %v", err)
	}

	client := jwtConfig.Client(ctx)

	return sheets.New(client)
}

func createDriveService(ctx context.Context, auth *Auth) (*drive.Service, error) {
	if auth.AuthType == "none" {
		return drive.NewService(ctx, option.WithAPIKey(auth.APIKey))
	}

	jwtConfig, err := google.JWTConfigFromJSON([]byte(auth.JWT), drive.DriveMetadataReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("Error parsin JWT file: %v", err)
	}

	client := jwtConfig.Client(ctx)

	return drive.New(client)
}
