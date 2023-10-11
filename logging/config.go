package logging

import (
	"github.com/ritsec/ops-bot-iii/config"
	"github.com/ritsec/ops-bot-iii/structs"
)

const (
	// Logging Levels
	DebugLowLevel structs.LogLevel = iota
	DebugLevel
	InfoLevel
	WarningLevel
	ErrorLevel
	CriticalLevel
	AuditLevel
)

var (
	// LevelChannelMap maps a logging level to a channel
	LevelChannelMap map[structs.LogLevel]*string = make(map[structs.LogLevel]*string)

	// LevelNameMap maps a logging level to a name
	LevelNameMap map[structs.LogLevel]string = make(map[structs.LogLevel]string)

	// NameLevelMap maps a name to a logging level
	NameLevelMap map[string]structs.LogLevel = make(map[string]structs.LogLevel)
)

func init() {
	// Channel Map
	LevelChannelMap[DebugLowLevel] = &config.Logging.DebugLowChannel
	LevelChannelMap[DebugLevel] = &config.Logging.DebugChannel
	LevelChannelMap[InfoLevel] = &config.Logging.InfoChannel
	LevelChannelMap[WarningLevel] = &config.Logging.WarningChannel
	LevelChannelMap[ErrorLevel] = &config.Logging.ErrorChannel
	LevelChannelMap[CriticalLevel] = &config.Logging.CriticalChannel
	LevelChannelMap[AuditLevel] = &config.Logging.AuditChannel

	// Name Map
	LevelNameMap[DebugLowLevel] = "Debug Low"
	LevelNameMap[DebugLevel] = "Debug"
	LevelNameMap[InfoLevel] = "Info"
	LevelNameMap[WarningLevel] = "Warning"
	LevelNameMap[ErrorLevel] = "Error"
	LevelNameMap[CriticalLevel] = "Critical"
	LevelNameMap[AuditLevel] = "Audit"

	// Level Map
	NameLevelMap["DebugLow"] = DebugLowLevel
	NameLevelMap["Debug"] = DebugLevel
	NameLevelMap["Info"] = InfoLevel
	NameLevelMap["Warning"] = WarningLevel
	NameLevelMap["Error"] = ErrorLevel
	NameLevelMap["Critical"] = CriticalLevel
	NameLevelMap["Audit"] = AuditLevel
}

// LogLevel returns the logging level
func LogLevel() structs.LogLevel {
	return NameLevelMap[config.Logging.Level]
}
