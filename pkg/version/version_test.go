package version

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetDocOptVersionString(t *testing.T) {

	version = "1"
	goversion = "1.11"

	exp := "test 1 built with 1.11"

	res := GetDocoptVersionString("test")

	assert.Equal(t, res, exp)
}
