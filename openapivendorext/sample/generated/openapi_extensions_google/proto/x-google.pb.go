// Code generated by protoc-gen-go.
// source: x-google.proto
// DO NOT EDIT!

/*
Package google is a generated protocol buffer package.

It is generated from these files:
	x-google.proto

It has these top-level messages:
	Book
	Shelve
*/
package google

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Book struct {
	Code    int64 `protobuf:"varint,1,opt,name=code" json:"code,omitempty"`
	Message int64 `protobuf:"varint,2,opt,name=message" json:"message,omitempty"`
}

func (m *Book) Reset()                    { *m = Book{} }
func (m *Book) String() string            { return proto.CompactTextString(m) }
func (*Book) ProtoMessage()               {}
func (*Book) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Book) GetCode() int64 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *Book) GetMessage() int64 {
	if m != nil {
		return m.Message
	}
	return 0
}

type Shelve struct {
	Foo1 int64 `protobuf:"varint,1,opt,name=foo1" json:"foo1,omitempty"`
	Bar  int64 `protobuf:"varint,2,opt,name=bar" json:"bar,omitempty"`
}

func (m *Shelve) Reset()                    { *m = Shelve{} }
func (m *Shelve) String() string            { return proto.CompactTextString(m) }
func (*Shelve) ProtoMessage()               {}
func (*Shelve) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *Shelve) GetFoo1() int64 {
	if m != nil {
		return m.Foo1
	}
	return 0
}

func (m *Shelve) GetBar() int64 {
	if m != nil {
		return m.Bar
	}
	return 0
}

func init() {
	proto.RegisterType((*Book)(nil), "google.Book")
	proto.RegisterType((*Shelve)(nil), "google.Shelve")
}

func init() { proto.RegisterFile("x-google.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 161 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xab, 0xd0, 0x4d, 0xcf,
	0xcf, 0x4f, 0xcf, 0x49, 0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x83, 0xf0, 0x94, 0x54,
	0xb9, 0x58, 0x9c, 0xf2, 0xf3, 0xb3, 0x85, 0x78, 0xb8, 0x58, 0x92, 0xf3, 0x53, 0x52, 0x25, 0x18,
	0x15, 0x18, 0x35, 0x98, 0x85, 0xf8, 0xb9, 0xd8, 0x73, 0x53, 0x8b, 0x8b, 0x13, 0xd3, 0x53, 0x25,
	0x98, 0x40, 0x02, 0x4a, 0xca, 0x5c, 0x6c, 0xc1, 0x19, 0xa9, 0x39, 0x65, 0xa9, 0x20, 0x85, 0x69,
	0xf9, 0xf9, 0x86, 0x50, 0x85, 0xdc, 0x5c, 0xcc, 0x49, 0x89, 0x45, 0x10, 0x45, 0x4e, 0xf6, 0x5c,
	0x32, 0xf9, 0x45, 0xe9, 0x7a, 0xf9, 0x05, 0xa9, 0x79, 0x89, 0x05, 0x99, 0x7a, 0xa9, 0x15, 0x25,
	0xa9, 0x79, 0xc5, 0x99, 0xf9, 0x79, 0x7a, 0x10, 0xbb, 0x9c, 0x44, 0xc2, 0x52, 0xf3, 0x52, 0xf2,
	0x8b, 0x5c, 0x61, 0xe2, 0x01, 0x20, 0x97, 0x04, 0x30, 0x2e, 0x62, 0x82, 0x3a, 0x26, 0x89, 0x0d,
	0xec, 0x36, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb5, 0x2d, 0x5a, 0x40, 0xad, 0x00, 0x00,
	0x00,
}
