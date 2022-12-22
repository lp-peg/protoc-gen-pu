package main

import (
	"bytes"
	_ "embed"
	"text/template"
)

var (
	//go:embed templates/uml.pu.tmpl
	umlTmpl []byte
	//go:embed templates/entity.pu.tmpl
	entityTmpl []byte
	//go:embed templates/reference.pu.tmpl
	referenceTmpl []byte
)

func main() {}

func render(classes []class, refs []reference) ([]byte, error) {
	tmpl := template.Must(template.New("uml").Parse(string(umlTmpl)))
	tmpl = template.Must(tmpl.AddParseTree("entity", template.Must(template.New("entity").Parse(string(entityTmpl))).Tree))
	tmpl = template.Must(tmpl.AddParseTree("references", template.Must(template.New("references").Parse(string(referenceTmpl))).Tree))
	buf := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(buf, "uml", map[string]interface{}{
		"Classes":    classes,
		"References": refs,
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
