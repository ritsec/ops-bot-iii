package google

import (
	"fmt"

	"gitlab.ritsec.cloud/1nv8rZim/ops-bot-iii/config"
	"google.golang.org/api/sheets/v4"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

// SheetsAppendSignin appends a signin to the signin sheet
func SheetsAppendSignin(userID string, username string, signinType string, ctx ddtrace.SpanContext) error {
	span := tracer.StartSpan(
		"google.sheets:SheetsAppendSignin",
		tracer.ResourceName("Google.SheetsAppendSignin"),
		tracer.ChildOf(ctx),
	)
	defer span.Finish()

	insertRequest := &sheets.Request{
		InsertDimension: &sheets.InsertDimensionRequest{
			Range: &sheets.DimensionRange{
				SheetId:    0,
				Dimension:  "ROWS",
				StartIndex: 1,
				EndIndex:   2,
			},
			InheritFromBefore: false,
		},
	}

	batchUpdateReq := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: []*sheets.Request{insertRequest},
	}

	_, err := sheetsSrv.Spreadsheets.BatchUpdate(config.Google.SheetID, batchUpdateReq).Do()
	if err != nil {
		return err
	}

	valueRange := &sheets.ValueRange{
		Values: [][]interface{}{{userID, username, signinType}},
	}

	rangeToUpdate := fmt.Sprintf("%s!A2:C2", config.Google.SheetName)

	appendCall := sheetsSrv.Spreadsheets.Values.Update(config.Google.SheetID, rangeToUpdate, valueRange)
	appendCall.ValueInputOption("RAW")
	_, err = appendCall.Do()

	return err
}
