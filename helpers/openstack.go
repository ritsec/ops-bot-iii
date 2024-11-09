package helpers

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/config"
	"github.com/ritsec/ops-bot-iii/logging"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
)

var (
	// Environment variables required for openstack CLI
	OS_AUTH_URL             string = config.GetString("openstack.ENV.OS_AUTH_URL")
	OS_PROJECT_ID           string = config.GetString("openstack.ENV.OS_PROJECT_ID")
	OS_PROJECT_NAME         string = config.GetString("openstack.ENV.OS_PROJECT_NAME")
	OS_USER_DOMAIN_NAME     string = config.GetString("openstack.ENV.OS_USER_DOMAIN_NAME")
	OS_PROJECT_DOMAIN_ID    string = config.GetString("openstack.ENV.OS_PROJECT_DOMAIN_ID")
	OS_USERNAME             string = config.GetString("openstack.ENV.OS_USERNAME")
	OS_PASSWORD             string = config.GetString("openstack.ENV.OS_PASSWORD")
	OS_REGION_NAME          string = config.GetString("openstack.ENV.OS_REGION_NAME")
	OS_INTERFACE            string = config.GetString("openstack.ENV.OS_INTERFACE")
	OS_IDENTITY_API_VERSION string = config.GetString("openstack.ENV.OS_IDENTITY_API_VERSION")

	// Paths for scripts to automate openstack user management
	new_member      string = config.GetString("openstack.SCRIPTS.new_member")
	reset_password  string = config.GetString("openstack.SCRIPTS.reset_password")
	check_if_exists string = config.GetString("openstack.SCRIPTS.check_if_exists")
)

func DebugCreate(s *discordgo.Session, user *discordgo.User, span ddtrace.Span, email string) (username string, password string, error error) {
	createCmd := exec.Command(new_member, email)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	createCmd.Stdout = stdout
	createCmd.Stderr = stderr

	// Create a combined output buffer
	combinedOutput := &bytes.Buffer{}
	createCmd.Stdout = io.MultiWriter(stdout, combinedOutput)
	createCmd.Stderr = io.MultiWriter(stderr, combinedOutput)

	err := createCmd.Run()
	if err != nil {
		logging.Debug(s, combinedOutput.String(), user, span)
		return "", "", err
	}

	logging.Debug(s, combinedOutput.String(), user, span)

	output := strings.Fields(stdout.String())

	username = output[0]
	password = output[1]

	return username, password, nil
}

func Create(email string) (username string, password string, error error) {
	createCmd := exec.Command(new_member, email)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	createCmd.Stdout = stdout
	createCmd.Stderr = stderr

	err := createCmd.Start()
	if err != nil {
		return "", "", err
	}
	err = createCmd.Wait()
	if err != nil {
		return "", "", err
	}
	output := strings.Fields(stdout.String())

	username = output[0]
	password = output[1]

	return username, password, nil
}

func Reset(email string) (username string, password string, error error) {
	resetCmd := exec.Command(reset_password, email)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	resetCmd.Stdout = stdout
	resetCmd.Stderr = stderr

	err := resetCmd.Run()
	if err != nil {
		return "", "", err
	}
	output := strings.Fields(stdout.String())

	username = output[0]
	password = output[1]

	return username, password, nil
}

func DebugCheckIfExists(s *discordgo.Session, user *discordgo.User, span ddtrace.Span, email string) (result bool, error error) {
	checkIfExistsCmd := exec.Command(check_if_exists, email)

	logging.Debug(s, checkIfExistsCmd.String(), user, span)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	checkIfExistsCmd.Stdout = stdout
	checkIfExistsCmd.Stderr = stderr
	// Create a combined output buffer
	combinedOutput := &bytes.Buffer{}
	checkIfExistsCmd.Stdout = io.MultiWriter(stdout, combinedOutput)
	checkIfExistsCmd.Stderr = io.MultiWriter(stderr, combinedOutput)

	err := checkIfExistsCmd.Run()
	if err != nil {
		logging.Debug(s, combinedOutput.String(), user, span)
		return false, err
	}

	logging.Debug(s, combinedOutput.String(), user, span)
	logging.Debug(s, fmt.Sprintf("combinedOutput type: %T", combinedOutput.String()), user, span)
	output := stdout.String()
	logging.Debug(s, fmt.Sprintf("Output type: %T", output), user, span)

	if string(output) == "0" {
		return false, nil
	} else {
		return true, nil
	}
}

func CheckIfExists(email string) (result bool, error error) {
	resetCmd := exec.Command(check_if_exists, email)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	resetCmd.Stdout = stdout
	resetCmd.Stderr = stderr

	err := resetCmd.Run()
	if err != nil {
		return false, err
	}
	output := stdout.String()

	if output == "1" {
		return true, nil
	} else {
		return false, nil
	}
}

func SetOpenstackRC() {
	os.Setenv("OS_AUTH_URL", OS_AUTH_URL)
	os.Setenv("OS_PROJECT_ID", OS_PROJECT_ID)
	os.Setenv("OS_PROJECT_NAME", OS_PROJECT_NAME)
	os.Setenv("OS_USER_DOMAIN_NAME", OS_USER_DOMAIN_NAME)
	os.Setenv("OS_PROJECT_DOMAIN_ID", OS_PROJECT_DOMAIN_ID)
	os.Setenv("OS_USERNAME", OS_USERNAME)
	os.Setenv("OS_PASSWORD", OS_PASSWORD)
	os.Setenv("OS_REGION_NAME", OS_REGION_NAME)
	os.Setenv("OS_INTERFACE", OS_INTERFACE)
	os.Setenv("OS_IDENTITY_API_VERSION", OS_IDENTITY_API_VERSION)
}
