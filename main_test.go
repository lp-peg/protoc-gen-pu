package main

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_render(t *testing.T) {
	type args struct {
		classes []*class
		refs    []*reference
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			args: args{
				classes: []*class{
					{
						Name: "Animal",
						Members: []member{
							{
								Name: "name",
								Type: "string",
							},
							{
								Name: "category",
								Type: "int32",
							},
						},
					},
					{
						Name: "Pet",
						Members: []member{
							{
								Name: "name",
								Type: "string",
							},
							{
								Name: "age",
								Type: "int64",
							},
						},
					},
				},
				refs: []*reference{
					{
						From: "Pet",
						To:   "Animal",
					},
				},
			},
			want: []byte(`
@startuml

skinparam linetype ortho

entity Animal {
  name (string)
  category (int32)
}

entity Pet {
  name (string)
  age (int64)
}

Animal <-- Pet

@enduml
`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := render(tt.args.classes, tt.args.refs)
			if (err != nil) != tt.wantErr {
				t.Errorf("render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			tt.want = bytes.TrimLeft(tt.want, "\n")
			tt.want = bytes.TrimRight(tt.want, "\n")
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Error(diff)
			}
		})
	}
}
