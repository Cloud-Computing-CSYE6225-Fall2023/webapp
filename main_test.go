package main

import (
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestIntegrations(t *testing.T) {
	// Start the server in a goroutine
	go func() {
		main()
	}()

	// Allow some time for the server to start
	time.Sleep(1 * time.Second)

	// Run integration tests
	cmd := "npx"
	args := []string{
		"newman",
		"run",
		"./tests/collection.json",
		"--reporters",
		"cli,junit",
		"--reporter-junit-export",
		"integration_report.xml",
		"--insecure",
	}

	//Create a command
	command := exec.Command(cmd, args...)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	err := command.Run()
	if err != nil {
		t.Errorf("Expected All Integration Tests to pass but got error")
	}
}
