package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintHeaders(t *testing.T) {
	exp := "[ Example ]\nreading file.txt"

	assert.Equal(t, exp, PrintHeader("example", "file.txt"))
}

func TestToJson(t *testing.T) {
	input := map[string]interface{}{
		"key_1": "value",
		"key_2": false,
	}

	exp := []byte(`{
  "key_1": "value",
  "key_2": false
}`)

	out := ToJSON(input)
	assert.Equal(t, exp, out)
}

func TestFromYAML(t *testing.T) {
	input := []byte(`key_1: value
key_2: false
`)
	expected := map[string]interface{}{
		"key_1": "value",
		"key_2": false,
	}

	m := make(map[string]interface{})
	FromYAML(input, &m)

	assert.Equal(t, expected, m)
}
