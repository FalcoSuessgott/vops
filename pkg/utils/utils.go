package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ghodss/yaml"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	yml "gopkg.in/yaml.v3"
)

// PrintHeader print outs the command header.
func PrintHeader(name, config string) string {
	c := cases.Title(language.English)

	return fmt.Sprintf(
		"[ %s ]\n"+
			"reading %s",
		c.String(name), config,
	)
}

// GetEnvs returns a map with all environment variables.
func GetEnvs() map[string]interface{} {
	m := make(map[string]interface{})

	for _, e := range os.Environ() {
		if i := strings.Index(e, "="); i >= 0 {
			m[e[:i]] = e[i+1:]
		}
	}

	return m
}

// ToJSON marshalls a given map to json.
func ToJSON(m interface{}) []byte {
	out, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		log.Fatalf("cannot marshal %v to JSON: %v", m, err)
	}

	return out
}

// FromJSON marshalls a given map to json.
func FromJSON(data []byte, i interface{}) error {
	if err := json.Unmarshal(data, i); err != nil {
		log.Fatalf("cannot unmarshal %v to %T: %v", data, i, err)
	}

	return nil
}

// FromYAML takes a yaml byte array and marshalls it into a map.
func FromYAML(b []byte, o interface{}) {
	if err := yaml.Unmarshal(b, &o); err != nil {
		log.Fatalf("cannot marshal %s to YAML: %v", string(b), err)
	}
}

// ToYAML marshalls a given map to yaml.
func ToYAML(m interface{}) []byte {
	var b bytes.Buffer

	yamlEncoder := yml.NewEncoder(&b)
	yamlEncoder.SetIndent(2)

	if err := yamlEncoder.Encode(m); err != nil {
		log.Fatalf("cannot marshal %T to YAML: %v", m, err)
	}

	return b.Bytes()
}

// GetCurrentDate returns the current date in YYYY-DD-MM-HH:mm:ss format.
func GetCurrentDate() string {
	return time.Now().Format("2006-01-02-15:04:05")
}
