package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

var (
	//go:embed templates/uml.pu.tmpl
	umlTmpl []byte
	//go:embed templates/entity.pu.tmpl
	entityTmpl []byte
	//go:embed templates/reference.pu.tmpl
	referenceTmpl []byte
	//go:embed templates/skinparam.pu.tmpl
	skinparamTmpl []byte
)

func str(s string) *string {
	return &s
}

type ss []string

func (s ss) tail() string {
	if len(s) == 0 {
		return ""
	}
	return (s)[len(s)-1]
}

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	req := new(pluginpb.CodeGeneratorRequest)
	if err := proto.Unmarshal(input, req); err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal(err)
	}
	res, err := run(req)
	if err != nil {
		log.Fatal(err)
	}
	out, err := proto.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := os.Stdout.Write(out); err != nil {
		log.Fatal(err)
	}
}

func render(classes []*class, refs []*reference, skinparams []*skinparam, showCircle bool) ([]byte, error) {
	tmpl := template.Must(template.New("uml").Parse(string(umlTmpl)))
	tmpl = template.Must(tmpl.AddParseTree("entity", template.Must(template.New("entity").Parse(string(entityTmpl))).Tree))
	tmpl = template.Must(tmpl.AddParseTree("references", template.Must(template.New("references").Parse(string(referenceTmpl))).Tree))
	tmpl = template.Must(tmpl.AddParseTree("skinparam", template.Must(template.New("skinparam").Parse(string(skinparamTmpl))).Tree))
	buf := new(bytes.Buffer)
	err := tmpl.ExecuteTemplate(buf, "uml", map[string]interface{}{
		"Classes":    classes,
		"References": refs,
		"Skinparams": skinparams,
		"ShowCircle": showCircle,
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func run(in *pluginpb.CodeGeneratorRequest) (*pluginpb.CodeGeneratorResponse, error) {
	genFiles := map[string]struct{}{}
	for _, file := range in.FileToGenerate {
		genFiles[file] = struct{}{}
	}
	classes := make([]*class, 0)
	refs := make([]*reference, 0)
	for _, file := range in.GetProtoFile() {
		if _, exists := genFiles[file.GetName()]; !exists {
			continue
		}
		for _, message := range file.GetMessageType() {
			c, r := parseMessage(message)
			classes = append(classes, c...)
			refs = append(refs, r...)
		}
		for _, enum := range file.GetEnumType() {
			classes = append(classes, parseEnum(enum)...)
		}
	}
	opt := option(*in)
	out, err := render(classes, refs, opt.getSkinparams(), opt.getCircleShow())
	if err != nil {
		return nil, err
	}
	return &pluginpb.CodeGeneratorResponse{
		File: []*pluginpb.CodeGeneratorResponse_File{
			{
				Name:    str((*option)(in).getOutputPath()),
				Content: str(string(out)),
			},
		},
	}, nil
}

func parseMessage(message *descriptorpb.DescriptorProto) ([]*class, []*reference) {
	c := class{
		Name: message.GetName(),
	}
	resClasses := []*class{&c}
	resRefs := make([]*reference, 0)
	for _, field := range message.GetField() {
		switch field.GetType() {
		case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, descriptorpb.FieldDescriptorProto_TYPE_ENUM:
			typename := ss(strings.Split(field.GetTypeName(), ".")).tail()
			c.Members = append(c.Members, member{
				Name: field.GetName(),
				Type: typename + "<FK>",
			})
			resRefs = append(resRefs, &reference{
				From: c.Name,
				To:   typename,
			})
		default:
			c.Members = append(c.Members, member{
				Name: field.GetName(),
				Type: strings.TrimPrefix(field.GetType().String(), "TYPE_"),
			})
		}
	}
	for _, nestedType := range message.GetNestedType() {
		nestedClasses, nestedRefs := parseMessage(nestedType)
		resClasses = append(resClasses, nestedClasses...)
		resRefs = append(resRefs, nestedRefs...)
	}
	for _, nestedEnum := range message.GetEnumType() {
		resClasses = append(resClasses, parseEnum(nestedEnum)...)
	}
	return resClasses, resRefs
}

func parseEnum(enum *descriptorpb.EnumDescriptorProto) []*class {
	c := class{
		Name: enum.GetName(),
	}
	for _, e := range enum.GetValue() {
		c.Members = append(c.Members, member{
			Name: e.GetName(),
			Type: fmt.Sprintf("ENUM: %v", e.GetNumber()),
		})
	}
	return []*class{&c}
}
