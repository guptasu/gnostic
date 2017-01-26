go get github.com/golang/protobuf/protoc-gen-go


# ONE TIME
#
#

############################# FOR E2E TESTING ########################################
pushd $GOPATH/src/github.com/googleapis/openapi-compiler/generator 
go build
cd ..
generator/generator

# ensure vendorextension proto contract is compiled.
pushd $GOPATH/src/github.com/googleapis/openapi-compiler/vendorextension
go install
popd

# ensure tool to generate the vendor extension compiler is installed.
pushd $GOPATH/src/github.com/googleapis/openapi-compiler/openapivendorextc
go install

####################################################################################33



# For every new extension
#
#

# run the tool to create your own vendor extension proto and compiler
pushd $GOPATH/src/github.com/googleapis/openapi-compiler/openapivendorextc
openapivendorextc sample/x-google-extensions.json --out_dir=$GOPATH/src/github.com/googleapis/openapi-compiler/vendorextension/generated_extension_plugins/sample_x_google_extension_plugin --extension_name_to_message=x-book:Book --extension_name_to_message=x-shelve:Shelve --proto_option_suffix=GoogleExtensions

# build and install your own vendor extension proto and compiler.
# This will allow openapic to discover your vendor extension handler
# and use it to create richer OpenAPI protobufs.
cd $GOPATH/src/github.com/googleapis/openapi-compiler/vendorextension/generated_extension_plugins/sample_x_google_extension_plugin
protoc --go_out=Mgoogle/protobuf/any.proto=github.com/golang/protobuf/ptypes/any:. *.proto
go install

popd
