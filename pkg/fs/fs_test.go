package fs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadFile(t *testing.T) {
	path := "testdata/file_1.txt"
	content := []byte("Hello World")

	assert.Equal(t, content, ReadFile(path))
}
