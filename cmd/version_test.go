package cmd

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	version := "v1.0.0"
	expected := "vops v1.0.0\n"

	c := versionCmd(version)
	b := bytes.NewBufferString("")
	c.SetOut(b)

	err := c.Execute()
	assert.NoError(t, err)

	out, _ := io.ReadAll(b)
	assert.Equal(t, expected, string(out))
}
