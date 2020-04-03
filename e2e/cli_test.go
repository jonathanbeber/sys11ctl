package e2e

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestCLI(t *testing.T) {
	// Create a temp dir
	dir, err := ioutil.TempDir("", "sys11ctl")
	if err != nil {
		fmt.Println("here")
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)
	binPath := dir + "sys11ctl"

	// build project
	build := exec.Command("go", "build", "-o", binPath, "..")
	cmdErr := build.Run()
	if cmdErr != nil {
		t.Fatalf("Failed to build project: %v", cmdErr)
	}

	t.Run("Test empty root command", func(t *testing.T) {
		root := exec.Command(binPath)
		var out bytes.Buffer
		root.Stdout = &out
		rootErr := root.Run()
		if rootErr != nil {
			t.Fatalf("Not expected error, received '%v'", rootErr)
		}
		if !strings.HasPrefix(out.String(), "sys11ctl is a SysEleven Metakube command line interface.") {
			t.Fatalf("Expected 'sys11ctl is a SysEleven Metakube command line interface.' on help message, received '%s'", out.String())
		}
	})

	t.Run("Test not existing command", func(t *testing.T) {
		nec := exec.Command(binPath, "notExistingCommand")
		var stdErr bytes.Buffer
		nec.Stderr = &stdErr
		necErr := nec.Run()
		if necErr == nil {
			t.Fatal("Expected error when a non existing command, none received")
		}

		if !strings.HasPrefix(stdErr.String(), "Error: unknown command") {
			t.Fatalf("Expected error 'Error: unknown command ...' on stdErr message, received '%s'", stdErr.String())
		}
	})

	t.Run("Test empty get command", func(t *testing.T) {
		get := exec.Command(binPath, "get")
		var getStdErr bytes.Buffer
		get.Stderr = &getStdErr
		getErr := get.Run()
		if getErr == nil {
			t.Fatal("Expected error when calling empty get command, none received")
		}

		if !strings.HasPrefix(getStdErr.String(), "Error: requires resources type argument") {
			t.Fatalf("Expected error 'Error: requires resources type argument' on stdErr message, received '%s'", getStdErr.String())
		}
	})

	t.Run("Test get command invalid resource type", func(t *testing.T) {
		get := exec.Command(binPath, "get", "invalidResource")
		var getStdErr bytes.Buffer
		get.Stderr = &getStdErr
		getErr := get.Run()
		if getErr == nil {
			t.Fatal("Expected error when calling get command with invalid resource type, none received")
		}

		if !strings.HasPrefix(getStdErr.String(), "Error: found not valid resource type 'invalidResource'") {
			t.Fatalf("Expected error 'Error: found not valid resource type 'invalidResource'' on stdErr message, received '%s'", getStdErr.String())
		}
	})

	// MOCK
	t.Run("Test get command valid resource type", func(t *testing.T) {
		get := exec.Command(binPath, "get", "projects")
		var getOut bytes.Buffer
		get.Stdout = &getOut
		getErr := get.Run()
		if getErr != nil {
			t.Fatalf("Foudn not expected error when calling  get command with valid resource type: %v", getErr)
		}

		if !strings.HasPrefix(getOut.String(), "Using config file") {
			t.Fatalf("Expected output 'get called', received '%s'", getOut.String())
		}
	})
}
