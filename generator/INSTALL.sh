go get github.com/golang/protobuf/protoc-gen-go

pushd $GOPATH/src/github.com/googleapis/openapi-compiler/generator 
go build
cd ..
generator/generator

pushd $GOPATH/src/github.com/googleapis/openapi-compiler/OpenAPIv2
protoc \
--go_out=Mgoogle/protobuf/any.proto=github.com/golang/protobuf/ptypes/any:. OpenAPIv2.proto 
go build
go install
popd