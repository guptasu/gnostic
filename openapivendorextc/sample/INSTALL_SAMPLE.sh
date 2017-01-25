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
openapivendorextc sample/x-mytest.json --out_dir=$GOPATH/src/github.com/googleapis/openapi-compiler/vendorextension/generated/x_mytest --extension_name=x-mytest --proto_message_name_prefix=MyTest

# build and install your own vendor extension proto and compiler.
# This will allow openapic to discover your vendor extension handler
# and use it to create richer OpenAPI protobufs.
cd $GOPATH/src/github.com/googleapis/openapi-compiler/vendorextension/generated/x_mytest
protoc --go_out=Mgoogle/protobuf/any.proto=github.com/golang/protobuf/ptypes/any:. *.proto
go install

popd
