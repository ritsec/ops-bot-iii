package helpers

import (
	"bytes"
	"io"
	"os/exec"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/ritsec/ops-bot-iii/logging"
	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace"
)

func DebugCreate(s *discordgo.Session, user *discordgo.User, span ddtrace.Span, email string) (username string, password string, error error) {
	createCmd := exec.Command("/root/automated_new_member.sh", email)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	createCmd.Stdout = stdout
	createCmd.Stderr = stderr

	combinedOutput := &bytes.Buffer{}
	createCmd.Stdout = io.MultiWriter(stdout, combinedOutput)
	createCmd.Stderr = io.MultiWriter(stderr, combinedOutput)

	err := createCmd.Run()
	if err != nil {
		logging.Error(s, combinedOutput.String(), user, span)
		return "", "", err
	}

	logging.Debug(s, combinedOutput.String(), user, span)
	output := strings.Fields(stdout.String())

	username = output[0]
	password = output[1]

	return username, password, nil
}

func Create(email string) (username string, password string, error error) {
	createCmd := exec.Command("/root/automated_new_member.sh", email)

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	createCmd.Stdout = stdout
	createCmd.Stderr = stderr

	err := createCmd.Run()
	if err != nil {
		return "", "", err
	}
	output := strings.Fields(stdout.String())

	username = output[0]
	password = output[1]

	return username, password, nil
}

func Reset(email string) (username string, password string, error error) {
	resetCmd := exec.Command("/root/automated_reset_password.sh", email)

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

func CheckIfExists(email string) (result bool, error error) {
	resetCmd := exec.Command("/root/automated_check_if_exists.sh", email)

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

func SourceOpenRC() error {
	sourceCmd := exec.Command("bash", "-c", "source /root/ops-openrc.sh")

	stderr := &bytes.Buffer{}
	sourceCmd.Stderr = stderr

	err := sourceCmd.Run()
	if err != nil {
		return err
	}
	return nil
}
