go get github.com/golang/protobuf/protoc-gen-go


# ONE TIME
#
#

############################# FOR E2E TESTING ########################################
# ensure vendorextension proto contract is compiled.
pushd $GOPATH/src/github.com/googleapis/openapi-compiler/openapivendorext/plugin
go install
popd

pushd $GOPATH/src/github.com/googleapis/openapi-compiler/generator 
go build
cd ..
generator/generator

# ensure tool to generate the vendor extension compiler is installed.
pushd $GOPATH/src/github.com/googleapis/openapi-compiler/openapivendorext/openapivendorextc
go install

####################################################################################33



# For every new extension
#
#

# run the tool to create your own vendor extension proto and compiler
pushd $GOPATH/src/github.com/googleapis/openapi-compiler/openapivendorext
openapivendorextc sample/x-google-extensions.json --out_dir_relative_to_gopath_src=github.com/googleapis/openapi-compiler/openapivendorext/sample/generated/sample_x_google_extension_plugin --extension_name_to_message=x-book:Book --extension_name_to_message=x-shelve:Shelve --proto_option_suffix=GoogleExtensions

# build and install your own vendor extension proto and compiler.
# This will allow openapic to discover your vendor extension handler
# and use it to create richer OpenAPI protobufs.
cd $GOPATH/src/github.com/googleapis/openapi-compiler/openapivendorext/sample/generated/sample_x_google_extension_plugin/proto
protoc --go_out=Mgoogle/protobuf/any.proto=github.com/golang/protobuf/ptypes/any:. *.proto
go install

cd $GOPATH/src/github.com/googleapis/openapi-compiler/openapivendorext/sample/generated/sample_x_google_extension_plugin
go install


popd
