package command

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultRunnerNoCommand(t *testing.T) {

	r := &DefaultRunner{}

	stdout, stderr, err := r.Run("")

	assert.Equal(t, err.Error(), "No command given")
	assert.Equal(t, stdout, "")
	assert.Equal(t, stderr, "")
}

func TestDefaultRunner(t *testing.T) {

	r := &DefaultRunner{}

	stdout, stderr, err := r.Run("ls command_test.go")

	assert.Equal(t, err, nil)
	assert.Equal(t, stderr, "")
	assert.Equal(t, stdout, "command_test.go")
}

func TestDefaultRunnerStdErr(t *testing.T) {

	r := &DefaultRunner{}

	stdout, stderr, err := r.Run("ls ls.go")

	assert.Equal(t, fmt.Sprintf("%T", err), "*exec.ExitError")
	assert.Regexp(t, "No such file or directory", stderr)
	assert.Equal(t, stdout, "")
}

func TestDefaultRunnerStdErrExitZero(t *testing.T) {

	r := &DefaultRunner{}

	stdout, stderr, err := r.Run("./stderr.sh")

	assert.Equal(t, err, nil)
	assert.Equal(t, stderr, "stderr")
	assert.Equal(t, stdout, "")
}
