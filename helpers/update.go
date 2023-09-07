package helpers

import (
	"bytes"
	"fmt"
	"os/exec"
)

// UpdateMainBranch switches to the main branch, fetches from origin, pulls from origin, and returns true if an update was pulled
func UpdateMainBranch() (bool, error) {
	switchCmd := exec.Command("git", "switch", "main")

	stderr := &bytes.Buffer{}
	switchCmd.Stderr = stderr

	err := switchCmd.Run()
	if err != nil {
		return false, fmt.Errorf("error switching to main branch: %s", stderr.String())
	}

	fetchCmd := exec.Command("git", "fetch", "origin", "main")

	stderr = &bytes.Buffer{}
	fetchCmd.Stderr = stderr

	err = fetchCmd.Run()
	if err != nil {
		return false, fmt.Errorf("error fetching from origin: %s", stderr.String())
	}

	commitCount := exec.Command("git", "rev-list", "--count", "HEAD...origin/main")

	stdout := &bytes.Buffer{}
	stderr = &bytes.Buffer{}

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

// BuildOBIII builds the OBIII binary
func BuildOBIII() error {
	buildCmd := exec.Command("/usr/local/go/bin/go", "build", "-o", "OBIII", "main.go")

	stderr := &bytes.Buffer{}
	buildCmd.Stderr = stderr

	err := buildCmd.Run()
	if err != nil {
		return fmt.Errorf("error building obiii:\n%s\n%s", err.Error(), stderr.String())
	}

	return nil
}

// Exit restarts the OBIII service
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
