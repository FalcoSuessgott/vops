package template

import (
	"os"
	"testing"
	"text/template"

	"github.com/FalcoSuessgott/vops/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	testCases := []struct {
		name string
		file string
		exp  []byte
		err  bool
	}{
		{
			name: "simple replacing",
			file: "testdata/config_1.yaml",
			exp: []byte(`Name: cluster-1
Addr: "http://127.0.0.1:8200"
KeyfilePath: "cluster-1.json"
Nodes:
  - "http://127.0.0.1:8200"
`),
			err: false,
		},
		{
			name: "recursive replacing",
			file: "testdata/config_2.yaml",
			exp: []byte(`Name: cluster-1
Addr: "http://127.0.0.1:8200"
TokenExecCmd: "jq -r root_token cluster-1.json"
KeyfilePath: "cluster-1.json"
Nodes:
  - "http://127.0.0.1:8200"
`),
			err: false,
		},
	}

	for _, tc := range testCases {
		out, err := os.ReadFile(tc.file)
		if err != nil {
			t.Fatalf("error reading file %s", tc.file)
		}

		var i interface{}

		utils.FromYAML(out, &i)

		res, err := Render(out, i)
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, string(tc.exp), res.String(), tc.name)
	}
}

func TestListTemplFields(t *testing.T) {
	str := "{{ .Name }} {{ .Field }} {{ .Value }}"
	exp := 3

	tpl, err := template.New("template").Option("missingkey=error").Parse(str)
	if err != nil {
		t.Fail()
	}

	res := ListTemplFields(tpl)

	assert.Equal(t, exp, len(res))
}
