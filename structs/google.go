package structs

// GoogleConfig is the config for google sheets
type GoogleConfig struct {
	// KeyFile is the path to the JSON key file
	KeyFile string

	// SheetName is the name of the sheet
	SheetName string

	// SheetID is the ID of the sheet
	SheetID string
}
