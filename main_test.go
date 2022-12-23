package main

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_render(t *testing.T) {
	type args struct {
		classes    []*class
		refs       []*reference
		skinparam  []*skinparam
		showCircle bool
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
					{
						Name: "Owner",
						Members: []member{
							{
								Name: "name",
								Type: "string",
							},
						},
					},
				},
				refs: []*reference{
					{
						From: "Pet",
						To:   "Animal",
					},
					{
						From: "Pet",
						To:   "Owner",
					},
				},
				skinparam: []*skinparam{
					{
						Param: "linetype",
						Value: "ortho",
					},
					{
						Param: "classFontSize",
						Value: "10",
					},
				},
				showCircle: false,
			},
			want: []byte(`
@startuml

hide circle
skinparam linetype ortho
skinparam classFontSize 10

entity Animal {
  name (string)
  category (int32)
}

entity Pet {
  name (string)
  age (int64)
}

entity Owner {
  name (string)
}

Animal <-- Pet
Owner <-- Pet

@enduml
`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := render(tt.args.classes, tt.args.refs, tt.args.skinparam, tt.args.showCircle)
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

func Test_option_getParam(t *testing.T) {
	type args struct {
		param string
	}
	tests := []struct {
		name string
		o    option
		args args
		want []string
	}{
		{
			o: option{
				Parameter: str("out=foo.pu"),
			},
			args: args{param: "out"},
			want: []string{"foo.pu"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.o.getParam(tt.args.param)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("option.getParam() = %v, want %v", got, tt.want)
			}
		})
	}
}
