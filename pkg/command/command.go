package command

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

// Runner interface, mostly exists so we can mock things for unit
// testing
type Runner interface {
	Run(string) (string, string, error)
}

// DefaultRunner implements a simple command runner
type DefaultRunner struct {
}

// Run executes a command on the DefaultRunner. Returns stdout, stderr and error
func (d DefaultRunner) Run(cmd string) (string, string, error) {
	cmdArray := strings.Split(cmd, " ")
	if len(cmdArray) == 0 || cmdArray[0] == "" {
		return "", "", fmt.Errorf("No command given")
	}
	var execCmd *exec.Cmd
	if len(cmdArray) == 1 {
		execCmd = exec.Command(cmdArray[0])
	} else {
		execCmd = exec.Command(cmdArray[0], cmdArray[1:len(cmdArray)]...)
	}
	stderr, err := execCmd.StderrPipe()
	if err != nil {
		return "", "", err
	}

	stdout, err := execCmd.StdoutPipe()
	if err != nil {
		return "", "", err
	}

	if err := execCmd.Start(); err != nil {
		return "", "", err
	}

	stderrBytes, err := ioutil.ReadAll(stderr)
	if err != nil {
		return "", "", err
	}
	stdoutBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", "", err
	}

	err = execCmd.Wait()
	stdoutString := strings.TrimSpace(string(stdoutBytes))
	stderrString := strings.TrimSpace(string(stderrBytes))
	return stdoutString, stderrString, err
}

// EchoingRunner echoes commands before running them and prints stdout and
// stderr
type EchoingRunner struct {
	runner DefaultRunner
}

// NewEchoingRunner returns a new EchoingRunner
func NewEchoingRunner() EchoingRunner {
	return EchoingRunner{
		runner: DefaultRunner{},
	}
}

// Run executes a command on the EchoingRunner. Returns stdout, stderr and error
func (e EchoingRunner) Run(cmd string) (string, string, error) {
	fmt.Printf("Running command: '%s'\n", cmd)
	stdout, stderr, err := e.runner.Run(cmd)
	if stdout != "" {
		fmt.Printf("stdout: %s\n", stdout)
	}
	if stderr != "" {
		fmt.Printf("stderr: %s\n", stderr)
	}
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}

	return stdout, stderr, err
}
