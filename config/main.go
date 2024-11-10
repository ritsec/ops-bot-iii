package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/ritsec/ops-bot-iii/structs"
	"github.com/spf13/viper"
)

var (
	// Token is the bot token
	Token string

	// AppID is the bot's application ID
	AppID string

	// GuildID is the bot's guild ID
	GuildID string

	// Logging is the logging configuration
	Logging structs.LoggingConfig

	// Google is the google configuration
	Google structs.GoogleConfig

	// Web is the web configuration
	Web structs.WebConfig

	// MailGun is the mailgun configuration
	MailGun structs.MailGunConfig

	// ServiceName is the name of the application
	ServiceName string

	// Environment is the environment the application is running in
	Environment string

	// Openstack is the boolean option for openstack configuration
	OpenstackEnabled bool
)

func init() {
	viper.SetConfigFile("config.yml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	UpdateConfigs()

	viper.OnConfigChange(
		func(e fsnotify.Event) {
			if e.Op == fsnotify.Write {
				UpdateConfigs()
			}
		},
	)

	viper.WatchConfig()

}

// UpdateConfigs updates the config variables
func UpdateConfigs() {
	Token = viper.GetString("token")
	AppID = viper.GetString("app_id")
	GuildID = viper.GetString("guild_id")
	ServiceName = viper.GetString("name")
	Environment = viper.GetString("env")
	OpenstackEnabled = viper.GetBool("openstack.enabled")
	Logging = logging()
	Google = google()
	Web = web()
	MailGun = mailgun()
}

// web returns the web configuration
func web() structs.WebConfig {
	return structs.WebConfig{
		Port:     viper.GetString("web.port"),
		Hostname: viper.GetString("web.hostname"),
	}
}

// logging returns the logging configuration
func logging() structs.LoggingConfig {
	return structs.LoggingConfig{
		Level:           viper.GetString("logging.level"),
		DebugLowChannel: viper.GetString("logging.debug_low_channel"),
		DebugChannel:    viper.GetString("logging.debug_channel"),
		InfoChannel:     viper.GetString("logging.info_channel"),
		WarningChannel:  viper.GetString("logging.warning_channel"),
		ErrorChannel:    viper.GetString("logging.error_channel"),
		CriticalChannel: viper.GetString("logging.critical_channel"),
		AuditChannel:    viper.GetString("logging.audit_channel"),
		LogFile:         viper.GetString("logging.log_file"),
	}
}

// google returns the google configuration
func google() structs.GoogleConfig {
	return structs.GoogleConfig{
		Enabled:   viper.GetBool("google.enabled"),
		KeyFile:   viper.GetString("google.key_file"),
		SheetName: viper.GetString("google.sheet_name"),
		SheetID:   viper.GetString("google.sheet_id"),
	}
}

// mailgun returns the mailgun configuration
func mailgun() structs.MailGunConfig {
	return structs.MailGunConfig{
		APIKey: viper.GetString("mailgun.api_key"),
		Domain: viper.GetString("mailgun.domain"),
	}
}

func openstack() {

}

// SetLoggingLevel sets the logging level
func SetLoggingLevel(level string) {
	viper.Set("logging.level", level)
	err := viper.WriteConfig()
	if err != nil {
		panic(err)
	}

	Logging.Level = level
}

// Get String Value from config
func GetString(key string) string {
	return viper.GetString(key)
}

// Get Int Value from config
func GetInt(key string) int {
	return viper.GetInt(key)
}

// Get Bool Value from config
func GetBool(key string) bool {
	return viper.GetBool(key)
}
