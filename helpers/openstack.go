package helpers

import (
	"bytes"
	"os/exec"
	"strings"
)

func Create(email string) (error error, username string, password string) {
	createCmd := exec.Command("root/automated_new_member.sh", email)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	createCmd.Stdout = stdout
	createCmd.Stderr = stderr

	err := createCmd.Run()
	if err != nil {
		return err, "", ""
	}
	output := strings.Fields(stdout.String())

	username = output[0]
	password = output[1]

	return nil, username, password
}

func Reset(email string) (error error, username string, password string) {
	resetCmd := exec.Command("root/automated_reset_password.sh", email)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	resetCmd.Stdout = stdout
	resetCmd.Stderr = stderr

	err := resetCmd.Run()
	if err != nil {
		return err, "", ""
	}
	output := strings.Fields(stdout.String())

	username = output[0]
	password = output[1]

	return nil, username, password
}

func CheckIfExists(email string) (error error, result bool) {
	resetCmd := exec.Command("root/automated_reset_password.sh", email)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	resetCmd.Stdout = stdout
	resetCmd.Stderr = stderr

	err := resetCmd.Run()
	if err != nil {
		return err, false
	}
	output := stdout.String()

	if output == "1" {
		return nil, true
	} else {
		return nil, false
	}
}
