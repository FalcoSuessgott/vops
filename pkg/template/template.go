package template

import (
	"bytes"
	"text/template"
	"text/template/parse"

	"github.com/FalcoSuessgott/vops/pkg/utils"
)

// String renders byte array input with the given data.
func Render(content []byte, input interface{}) (bytes.Buffer, error) {
	var buf bytes.Buffer

	// fmt.Println(string(content))
	// fmt.Printf("%#v\n", input)

	tpl, err := template.New("template").Option("missingkey=error").Parse(string(content))
	if err != nil {
		return buf, err
	}

	if err := tpl.Execute(&buf, &input); err != nil {
		return buf, err
	}

	if len(ListTemplFields(tpl)) > 0 {
		var i interface{}

		utils.FromYAML(buf.Bytes(), &i)

		return Render(buf.Bytes(), i)
	}

	return buf, nil
}

func ListTemplFields(t *template.Template) []string {
	return listNodeFields(t.Tree.Root, nil)
}

func listNodeFields(node parse.Node, res []string) []string {
	if node.Type() == parse.NodeAction {
		res = append(res, node.String())
	}

	if ln, ok := node.(*parse.ListNode); ok {
		for _, n := range ln.Nodes {
			res = listNodeFields(n, res)
		}
	}

	return res
}
