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

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"runtime"

	"path/filepath"

	"strings"

	"github.com/googleapis/openapi-compiler/generator/util"
	"github.com/googleapis/openapi-compiler/jsonschema"
)

const LICENSE = "" +
	"// Copyright 2016 Google Inc. All Rights Reserved.\n" +
	"//\n" +
	"// Licensed under the Apache License, Version 2.0 (the \"License\");\n" +
	"// you may not use this file except in compliance with the License.\n" +
	"// You may obtain a copy of the License at\n" +
	"//\n" +
	"//    http://www.apache.org/licenses/LICENSE-2.0\n" +
	"//\n" +
	"// Unless required by applicable law or agreed to in writing, software\n" +
	"// distributed under the License is distributed on an \"AS IS\" BASIS,\n" +
	"// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.\n" +
	"// See the License for the specific language governing permissions and\n" +
	"// limitations under the License.\n"

const additionalCompilerCodeWithMain = "" +
	"type documentHandler func(name string, version string, extensionName string, document string)\n" +
	"func forInputYamlFromOpenapic(handler documentHandler) {\n" +
	"	data, err := ioutil.ReadAll(os.Stdin)\n" +
	"\n" +
	"	if err != nil {\n" +
	"		fmt.Println(\"File error:\", err.Error())\n" +
	"		os.Exit(1)\n" +
	"	}\n" +
	"	request := &vendorextension.VendorExtensionHandlerRequest{}\n" +
	"	err = proto.Unmarshal(data, request)\n" +
	"	handler(request.Wrapper.Name, request.Wrapper.Version, request.Wrapper.ExtensionName, request.Wrapper.Yaml)\n" +
	"}\n" +
	"\n" +
	"func main() {\n" +
	"	response := &vendorextension.VendorExtensionHandlerResponse{}\n" +
	"	forInputYamlFromOpenapic(\n" +
	"		func(name string, version string, extensionName string, yamlInput string) {\n" +
	"		var info yaml.MapSlice\n" +
	"		var newObject proto.Message\n" +
	"       var err error\n" +
	"		err = yaml.Unmarshal([]byte(yamlInput), &info)\n" +
	"		if err != nil {\n" +
	"			response.Error = append(response.Error, err.Error())\n" +
	"			responseBytes, _ := proto.Marshal(response)\n" +
	"			os.Stdout.Write(responseBytes)\n" +
	"			os.Exit(0)\n" +
	"		}\n" +
	"      \n" +
	"\n" +
	"      switch extensionName {\n" +
	"      // All supported extensions\n" +
	"      %s\n" +
	"      default:\n" +
	"          responseBytes, _ := proto.Marshal(response)\n" +
	"          os.Stdout.Write(responseBytes)\n" +
	"          os.Exit(0)\n" +
	"       }" +
	"		// If we reach hear, then the extension is handled\n" +
	"		response.Handled = true\n" +
	"		if err != nil {\n" +
	"			response.Error = append(response.Error, err.Error())\n" +
	"			responseBytes, _ := proto.Marshal(response)\n" +
	"			os.Stdout.Write(responseBytes)\n" +
	"			os.Exit(0)\n" +
	"		}\n" +
	"		response.Value, err = ptypes.MarshalAny(newObject)\n" +
	"		if err != nil {\n" +
	"			response.Error = append(response.Error, err.Error())\n" +
	"			responseBytes, _ := proto.Marshal(response)\n" +
	"			os.Stdout.Write(responseBytes)\n" +
	"			os.Exit(0)\n" +
	"		}\n" +
	"		})\n" +
	"\n" +
	"	responseBytes, _ := proto.Marshal(response)\n" +
	"	os.Stdout.Write(responseBytes)\n" +
	"}\n"

const caseString = "\n" +
	"case \"%s\":\n" +
	"newObject, err = New%s(info, compiler.NewContextWithCustomAnyProtoGenerators(\"$root\", nil, nil))\n"

var PROTO_OPTIONS = []util.ProtoOption{
	util.ProtoOption{
		Name:  "java_multiple_files",
		Value: "true",
		Comment: "// This option lets the proto compiler generate Java code inside the package\n" +
			"// name (see below) instead of inside an outer class. It creates a simpler\n" +
			"// developer experience by reducing one-level of name nesting and be\n" +
			"// consistent with most programming languages that don't support outer classes.",
	},

	util.ProtoOption{
		Name:  "java_outer_classname",
		Value: "VendorExtensionProto",
		Comment: "// The Java outer classname should be the filename in UpperCamelCase. This\n" +
			"// class is only used to hold proto descriptor, so developers don't need to\n" +
			"// work with it directly.",
	},
}

func main() {
	// the OpenAPI schema file and API version are hard-coded for now

	usage := `
Usage: TODO
`
	outDir := ""
	schameFile := ""
	protoOptionSuffix := ""
	plugin_regex, _ := regexp.Compile("--(.+)=(.+)")
	plugin_extension_to_message_regex, _ := regexp.Compile("(.+):(.+)")

	extensionToMessage := make(map[string]string)

	for i, arg := range os.Args {
		if i == 0 {
			continue // skip the tool name
		}
		var m [][]byte
		if m = plugin_regex.FindSubmatch([]byte(arg)); m != nil {
			flagName := string(m[1])
			flagValue := string(m[2])
			switch flagName {
			case "out_dir":
				outDir = flagValue
			case "proto_option_suffix":
				protoOptionSuffix = flagValue
			case "extension_name_to_message":
				var t [][]byte
				if t = plugin_extension_to_message_regex.FindSubmatch([]byte(flagValue)); t != nil {
					extensionToMessage[string(t[1])] = string(t[2])
				} else {
					fmt.Printf("Unknown option: %s.\n%s\n", arg, usage)
					os.Exit(-1)
				}
			default:
				fmt.Printf("Unknown option: %s.\n%s\n", arg, usage)
				os.Exit(-1)
			}
		} else if arg[0] == '-' {
			fmt.Printf("Unknown option: %s.\n%s\n", arg, usage)
			os.Exit(-1)
		} else {
			schameFile = arg
		}
	}

	if len(extensionToMessage) == 0 {
		fmt.Printf("No extension_name_to_message specified.\n%s\n", usage)
		os.Exit(-1)
	}

	if schameFile == "" {
		fmt.Printf("No input json schema specified.\n%s\n", usage)
		os.Exit(-1)
	}
	if protoOptionSuffix == "" {
		fmt.Printf("No proto_option_suffix specified.\n%s\n", usage)
		os.Exit(-1)
	}
	if outDir == "" {
		fmt.Printf("Missing output directive.\n%s\n", usage)
		os.Exit(-1)
	}

	var cases string
	for extensionName, messagType := range extensionToMessage {
		cases += fmt.Sprintf(caseString, extensionName, messagType)
	}
	additionalCompilerCodeWithMainReplaced := fmt.Sprintf(additionalCompilerCodeWithMain, cases)
	outFileBaseName := filepath.Base(schameFile)
	outFileBaseName = outFileBaseName[0 : len(outFileBaseName)-len(filepath.Ext(outFileBaseName))]

	proto_packagename := "main"
	go_packagename := "main"

	additionalImports := []string{
		"io/ioutil",
		"os",
		"github.com/golang/protobuf/proto",
		"github.com/googleapis/openapi-compiler/vendorextension",
		"github.com/golang/protobuf/ptypes",
	}

	base_schema := jsonschema.NewSchemaFromFile("schema.json")
	base_schema.ResolveRefs()
	base_schema.ResolveAllOfs()

	openapi_schema := jsonschema.NewSchemaFromFile(schameFile)
	openapi_schema.ResolveRefs()
	openapi_schema.ResolveAllOfs()

	// build a simplified model of the types described by the schema
	cc := util.NewDomain(openapi_schema)
	cc.Build()

	var err error

	err = os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	// generate the protocol buffer description

	PROTO_OPTIONS = append(PROTO_OPTIONS,
		util.ProtoOption{Name: "java_package", Value: "org.openapi.extension." + strings.ToLower(protoOptionSuffix), Comment: "// The Java package name must be proto package name with proper prefix."},
		util.ProtoOption{Name: "objc_class_prefix", Value: strings.ToLower(protoOptionSuffix),
			Comment: "// A reasonable prefix for the Objective-C symbols generated from the package.\n" +
				"// It should at a minimum be 3 characters long, all uppercase, and convention\n" +
				"// is to use an abbreviation of the package name. Something short, but\n" +
				"// hopefully unique enough to not conflict with things that may come along in\n" +
				"// the future. 'GPB' is reserved for the protocol buffer implementation itself.",
		})

	proto := cc.GenerateProto(proto_packagename, LICENSE, PROTO_OPTIONS)
	proto_filename := outDir + "/" + outFileBaseName + ".proto"

	err = ioutil.WriteFile(proto_filename, []byte(proto), 0644)
	if err != nil {
		panic(err)
	}

	// generate the compiler
	compiler := cc.GenerateCompiler(go_packagename, LICENSE, additionalCompilerCodeWithMainReplaced, additionalImports)
	go_filename := outDir + "/" + outFileBaseName + ".go"
	err = ioutil.WriteFile(go_filename, []byte(compiler), 0644)
	if err != nil {
		panic(err)
	}

	// format the compiler
	err = exec.Command(runtime.GOROOT()+"/bin/gofmt", "-w", go_filename).Run()
}
