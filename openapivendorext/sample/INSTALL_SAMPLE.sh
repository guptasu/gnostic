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
./INSTALL.sh

pushd $GOPATH/src/github.com/googleapis/openapi-compiler/openapic
go install

# ensure tool to generate the vendor extension compiler is installed.
pushd $GOPATH/src/github.com/googleapis/openapi-compiler/openapivendorext/openapivendorextc
go install

####################################################################################33



# Now generate sample extension plugins and install them.
#
#
pushd $GOPATH/src/github.com/googleapis/openapi-compiler/openapivendorext

    EXTENSION_OUT_DIR="github.com/googleapis/openapi-compiler/openapivendorext/sample/generated"
    # For Google Extension Example
    #
    #
    GOOGLE_EXTENSION_SCHEMA="sample/x-google.json"
    GOOGLE_EXTENSION_PROTO_NAMESPACE="GoogleExtensions"

    openapivendorextc $GOOGLE_EXTENSION_SCHEMA --out_dir_relative_to_gopath_src=$EXTENSION_OUT_DIR \
    --proto_option_suffix=$GOOGLE_EXTENSION_PROTO_NAMESPACE \
    --extension_name_to_message=x-book:Book \
    --extension_name_to_message=x-shelve:Shelve \

    pushd $GOPATH/src/$EXTENSION_OUT_DIR/openapi_extensions_google/proto
        protoc --go_out=Mgoogle/protobuf/any.proto=github.com/golang/protobuf/ptypes/any:. *.proto
        go install
    popd

    pushd  $GOPATH/src/$EXTENSION_OUT_DIR/openapi_extensions_google
        go install
    popd

    # For IBM Extension Example
    #
    #
    IBM_EXTENSION_SCHEMA="sample/x-ibm.json"
    IBM_EXTENSION_PROTO_NAMESPACE="IbmExtensions"

    openapivendorextc $IBM_EXTENSION_SCHEMA --out_dir_relative_to_gopath_src=$EXTENSION_OUT_DIR \
    --proto_option_suffix=$IBM_EXTENSION_PROTO_NAMESPACE \
    --extension_name_to_message=x-ibm-book:IbmBook \
    --extension_name_to_message=x-ibm-shelve:IbmShelve \

    pushd $GOPATH/src/$EXTENSION_OUT_DIR/openapi_extensions_ibm/proto
        protoc --go_out=Mgoogle/protobuf/any.proto=github.com/golang/protobuf/ptypes/any:. *.proto
        go install
    popd

    pushd $GOPATH/src/$EXTENSION_OUT_DIR/openapi_extensions_ibm
        go install
    popd

popd
