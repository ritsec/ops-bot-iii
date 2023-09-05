package helpers

import (
	"bytes"
	"fmt"
	"os/exec"
)

func UpdateMainBranch() (bool, error) {
	err := exec.Command("git", "switch", "main").Run()
	if err != nil {
		return false, err
	}

	err = exec.Command("git", "fetch", "origin", "main").Run()
	if err != nil {
		return false, err
	}

	commitCount := exec.Command("git", "rev-list", "--count", "HEAD...origin/main")

	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}

	commitCount.Stdout = stdout
	commitCount.Stderr = stderr

	err = commitCount.Run()
	if err != nil {
		return false, fmt.Errorf("error getting commit count: %s", stderr.String())
	}

	if stdout.String() == "0\n" {
		return false, nil
	}

	pullCmd := exec.Command("git", "pull", "origin", "main")

	stderr = &bytes.Buffer{}
	pullCmd.Stderr = stderr

	err = pullCmd.Run()
	if err != nil {
		return false, fmt.Errorf("error pulling from origin: %s", stderr.String())
	}

	return true, nil
}

func BuildOBIII() error {
	buildCmd := exec.Command("go", "build", "-o", "OBIII", "main.go")

	stderr := &bytes.Buffer{}
	buildCmd.Stderr = stderr

	err := buildCmd.Run()
	if err != nil {
		return fmt.Errorf("error building obiii: %s", stderr.String())
	}

	return nil
}

func Exit() error {
	exitCmd := exec.Command("systemctl", "restart", "OBIII")

	stderr := &bytes.Buffer{}
	exitCmd.Stderr = stderr

	err := exitCmd.Run()
	if err != nil {
		return fmt.Errorf("error restarting obiii: %s", stderr.String())
	}

	return nil
}
