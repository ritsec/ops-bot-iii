package structs

// LogLevel is the type for log levels
type LogLevel uint8

type LoggingConfig struct {
	// Level is the log level
	Level string

	// DebugLowChannel is the channel to log debug messages to
	DebugLowChannel string

	// DebugHighChannel is the channel to log debug messages to
	DebugChannel string

	// InfoChannel is the channel to log info messages to
	InfoChannel string

	// WarningChannel is the channel to log warning messages to
	WarningChannel string

	// ErrorChannel is the channel to log error messages to
	ErrorChannel string

	// CriticalChannel is the channel to log critical messages to
	CriticalChannel string

	// AuditChannel is the channel to log audit messages to
	AuditChannel string

	// LogFile is the file to log to
	LogFile string
}
