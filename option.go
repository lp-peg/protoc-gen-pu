package main

import (
	"strings"

	"google.golang.org/protobuf/types/pluginpb"
)

type option pluginpb.CodeGeneratorRequest

func (o *option) getParam(p string) []string {
	if o.Parameter == nil {
		return nil
	}
	res := make([]string, 0)
	for _, param := range strings.Split(*o.Parameter, ",") {
		kv := strings.Split(param, "=")
		if len(kv) == 2 && kv[0] == p {
			res = append(res, kv[1])
		}
	}
	return res
}

func (o *option) getOutputPath() string {
	const (
		outputFilenameParamName = "out"
		defaultOutputFilename   = "out.pu"
	)
	if v := o.getParam(outputFilenameParamName); len(v) > 0 {
		return ss(v).tail()
	}
	return defaultOutputFilename
}

func (o option) getSkinparams() []*skinparam {
	const skinparamParamName = "skinparams"
	var defaultSkinParam = []*skinparam{
		{
			Param: "linetype",
			Value: "ortho",
		},
	}
	if v := o.getParam(skinparamParamName); len(v) > 0 {
		res := make([]*skinparam, 0)
		for _, kv := range v {
			param := strings.Split(kv, ":")
			if len(param) != 2 {
				continue
			}
			res = append(res, &skinparam{param[0], param[1]})
		}
		return res
	}
	return defaultSkinParam
}

func (o option) getCircleShow() bool {
	const circleShowParamName = "circle"
	p := o.getParam(circleShowParamName)
	return len(p) > 0 && p[len(p)-1] == "show"
}
