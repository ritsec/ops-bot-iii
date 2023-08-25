package google

import (
	"context"
	"log"
	"net/http"
	"os"

	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/config"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var (
	// SheetsSrv is the sheets service
	sheetsSrv *sheets.Service

	// Client is the http client
	client *http.Client
)

func init() {
	// Read the JSON key file
	keyBytes, err := os.ReadFile(config.Google.KeyFile)
	if err != nil {
		log.Fatalf("Failed to read the JSON key file: %v", err)
	}

	// Create a JWT config from the JSON key file
	jwtConfig, err := google.JWTConfigFromJSON(keyBytes, sheets.SpreadsheetsScope, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Failed to create JWT config: %v", err)
	}

	// Create an OAuth2 client using the JWT config
	ctx := context.Background()
	_client := jwtConfig.Client(ctx)

	client = _client

	_sheetsSrv, err := sheets.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets client: %v", err)
	}
	sheetsSrv = _sheetsSrv
}
