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

// THIS FILE IS AUTOMATICALLY GENERATED.

package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/googleapis/openapi-compiler/compiler"
	"github.com/googleapis/openapi-compiler/vendorextension"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

func Version() string {
	return "main"
}

func NewAny(in interface{}, context *compiler.Context) (*Any, error) {
	errors := make([]error, 0)
	x := &Any{}
	bytes, _ := yaml.Marshal(in)
	x.Value = &any.Any{TypeUrl: string(bytes)}
	return x, compiler.NewErrorGroupOrNil(errors)
}

func NewMyTestDocument(in interface{}, context *compiler.Context) (*MyTestDocument, error) {
	errors := make([]error, 0)
	x := &MyTestDocument{}
	m, ok := compiler.UnpackMap(in)
	if !ok {
		message := fmt.Sprintf("has unexpected value: %+v", in)
		errors = append(errors, compiler.NewError(context, message))
	} else {
		requiredKeys := []string{"code", "message"}
		missingKeys := compiler.MissingKeysInMap(m, requiredKeys)
		if len(missingKeys) > 0 {
			message := fmt.Sprintf("is missing required %s: %+v", compiler.PluralProperties(len(missingKeys)), strings.Join(missingKeys, ", "))
			errors = append(errors, compiler.NewError(context, message))
		}
		allowedKeys := []string{"code", "message"}
		allowedPatterns := []string{}
		invalidKeys := compiler.InvalidKeysInMap(m, allowedKeys, allowedPatterns)
		if len(invalidKeys) > 0 {
			message := fmt.Sprintf("has invalid %s: %+v", compiler.PluralProperties(len(invalidKeys)), strings.Join(invalidKeys, ", "))
			errors = append(errors, compiler.NewError(context, message))
		}
		// int64 code = 1;
		v1 := compiler.MapValueForKey(m, "code")
		if v1 != nil {
			t, ok := v1.(int)
			if ok {
				x.Code = int64(t)
			} else {
				message := fmt.Sprintf("has unexpected value for code: %+v", v1)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
		// int64 message = 2;
		v2 := compiler.MapValueForKey(m, "message")
		if v2 != nil {
			t, ok := v2.(int)
			if ok {
				x.Message = int64(t)
			} else {
				message := fmt.Sprintf("has unexpected value for message: %+v", v2)
				errors = append(errors, compiler.NewError(context, message))
			}
		}
	}
	return x, compiler.NewErrorGroupOrNil(errors)
}

func NewStringArray(in interface{}, context *compiler.Context) (*StringArray, error) {
	errors := make([]error, 0)
	x := &StringArray{}
	a, ok := in.([]interface{})
	if !ok {
		message := fmt.Sprintf("has unexpected value for StringArray: %+v", in)
		errors = append(errors, compiler.NewError(context, message))
	} else {
		x.Value = make([]string, 0)
		for _, s := range a {
			x.Value = append(x.Value, s.(string))
		}
	}
	return x, compiler.NewErrorGroupOrNil(errors)
}

func (m *Any) ResolveReferences(root string) (interface{}, error) {
	errors := make([]error, 0)
	return nil, compiler.NewErrorGroupOrNil(errors)
}

func (m *MyTestDocument) ResolveReferences(root string) (interface{}, error) {
	errors := make([]error, 0)
	return nil, compiler.NewErrorGroupOrNil(errors)
}

func (m *StringArray) ResolveReferences(root string) (interface{}, error) {
	errors := make([]error, 0)
	return nil, compiler.NewErrorGroupOrNil(errors)
}

type documentHandler func(name string, version string, document string)

func forInputYamlFromOpenapic(handler documentHandler) {
	data, err := ioutil.ReadAll(os.Stdin)

	if err != nil {
		fmt.Println("File error:", err.Error())
		os.Exit(1)
	}
	request := &vendorextension.VendorExtensionHandlerRequest{}
	err = proto.Unmarshal(data, request)
	handler(request.Wrapper.Name, request.Wrapper.Version, request.Wrapper.Yaml)
}

func main() {
	response := &vendorextension.VendorExtensionHandlerResponse{}
	forInputYamlFromOpenapic(
		func(name string, version string, yamlInput string) {
			var info yaml.MapSlice
			err := yaml.Unmarshal([]byte(yamlInput), &info)
			if err != nil {
				response.Error = append(response.Error, err.Error())
				responseBytes, _ := proto.Marshal(response)
				os.Stdout.Write(responseBytes)
				os.Exit(0)
			}

			newObject, err := NewMyTestDocument(info, compiler.NewContextWithPatternFieldProtoGenerators("$root", nil, nil))
			if err != nil {
				response.Error = append(response.Error, err.Error())
				responseBytes, _ := proto.Marshal(response)
				os.Stdout.Write(responseBytes)
				os.Exit(0)
			}
			response.Value, err = ptypes.MarshalAny(newObject)
			if err != nil {
				response.Error = append(response.Error, err.Error())
				responseBytes, _ := proto.Marshal(response)
				os.Stdout.Write(responseBytes)
				os.Exit(0)
			}
		})

	responseBytes, _ := proto.Marshal(response)
	os.Stdout.Write(responseBytes)
}
