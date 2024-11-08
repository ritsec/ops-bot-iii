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
	createCmd := exec.Command("/root/automated_new_member.sh", email)

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

func DebugCheckIfExists(s *discordgo.Session, user *discordgo.User, span ddtrace.Span, email string) (result bool, error error) {
	checkIfExistsCmd := exec.Command("/root/automated_check_if_exists.sh", email)

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
	output := stdout.String()

	if output == "1" {
		return true, nil
	} else {
		return false, nil
	}
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

func DebugSourceOpenRC(s *discordgo.Session, user *discordgo.User, span ddtrace.Span) error {
	sourceCmd := exec.Command("bash", "-c", "source /root/ops-openrc.sh")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	sourceCmd.Stdout = stdout
	sourceCmd.Stderr = stderr
	// Create a combined output buffer
	combinedOutput := &bytes.Buffer{}
	sourceCmd.Stdout = io.MultiWriter(stdout, combinedOutput)
	sourceCmd.Stderr = io.MultiWriter(stderr, combinedOutput)

	err := sourceCmd.Run()
	if err != nil {
		logging.Debug(s, combinedOutput.String(), user, span)
		return err
	}

	logging.Debug(s, combinedOutput.String(), user, span)
	return nil
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
