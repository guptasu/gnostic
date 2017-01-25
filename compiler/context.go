// Copyright 2016 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package compiler

import (
	"bytes"
	"fmt"
	"os/exec"

	"strings"

	"errors"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/googleapis/openapi-compiler/vendorextension"
	yaml "gopkg.in/yaml.v2"
)

type PatternFieldProtoGenerator struct {
	PatternFieldName   string
	ProtoGeneratorName string
}

func (pluginCall *PatternFieldProtoGenerator) Perform(in interface{}) (*any.Any, error) {
	if pluginCall.ProtoGeneratorName != "" {
		binary, _ := yaml.Marshal(in)

		request := &vendorextension.VendorExtensionHandlerRequest{}
		request.Parameter = ""

		version := &vendorextension.Version{}
		version.Major = 0
		version.Minor = 1
		version.Patch = 0
		request.CompilerVersion = version

		request.Wrapper = &vendorextension.Wrapper{}
		request.Wrapper.Name = "TESTETEST"
		request.Wrapper.Version = "v2"

		request.Wrapper.Yaml = string(binary)
		requestBytes, _ := proto.Marshal(request)

		cmd := exec.Command(pluginCall.ProtoGeneratorName)
		cmd.Stdin = bytes.NewReader(requestBytes)
		output, err := cmd.Output()
		if err != nil {
			fmt.Printf("Error: %+v\n", err)
			return nil, err
		}
		response := &vendorextension.VendorExtensionHandlerResponse{}
		err = proto.Unmarshal(output, response)
		if err != nil {
			fmt.Printf("Error: %+v\n", err)
			fmt.Printf("%s\n", string(output))
			return nil, err
		}
		if len(response.Error) != 0 {
			message := fmt.Sprintf("Errors when parsing: %+v for field %s by vendor extension handler %s. Details %+v", in, pluginCall.PatternFieldName, pluginCall.ProtoGeneratorName, strings.Join(response.Error, ","))
			return nil, errors.New(message)
		}
		return response.Value, nil
	}
	return nil, nil
}

type Context struct {
	Parent *Context
	Name   string

	// TODO: Figure out a better way to pass the patternFieldProtoGenerators to the generated compiler.
	PatternFieldProtoGenerators *[]PatternFieldProtoGenerator
}

func NewContextWithPatternFieldProtoGenerators(name string, parent *Context, patternFieldProtoGenerators *[]PatternFieldProtoGenerator) *Context {
	return &Context{Name: name, Parent: parent, PatternFieldProtoGenerators: patternFieldProtoGenerators}
}

func NewContext(name string, parent *Context) *Context {
	return &Context{Name: name, Parent: parent, PatternFieldProtoGenerators: parent.PatternFieldProtoGenerators}
}

func (context *Context) Description() string {
	if context.Parent != nil {
		return context.Parent.Description() + "." + context.Name
	} else {
		return context.Name
	}
}
