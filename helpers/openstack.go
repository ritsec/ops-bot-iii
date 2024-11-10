package helpers

import (
	"bytes"
	"os"
	"os/exec"
	"strings"

	"github.com/ritsec/ops-bot-iii/config"
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

	err := resetCmd.Start()
	if err != nil {
		return "", "", err
	}
	err = resetCmd.Wait()
	if err != nil {
		return "", "", err
	}

	output := strings.Fields(stdout.String())
	username = output[0]
	password = output[1]

	return username, password, nil
}

func CheckIfExists(email string) (result bool, error error) {
	checkIfExistsCmd := exec.Command(check_if_exists, email)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	checkIfExistsCmd.Stdout = stdout
	checkIfExistsCmd.Stderr = stderr

	err := checkIfExistsCmd.Run()
	if err != nil {
		return false, err
	}

	output := strings.TrimSpace(stdout.String())
	if output == "0" {
		return false, nil
	} else {
		return true, nil
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
