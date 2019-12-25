// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

/*
Package message is a generated protocol buffer package.

It is generated from these files:
	message.proto

It has these top-level messages:
	Message
*/
package message

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

type Message struct {
	Topic      string `protobuf:"bytes,1,opt,name=Topic" json:"Topic,omitempty"`
	Index      []byte `protobuf:"bytes,2,opt,name=Index,proto3" json:"Index,omitempty"`
	Payload    []byte `protobuf:"bytes,3,opt,name=Payload,proto3" json:"Payload,omitempty"`
	Qos        int32  `protobuf:"varint,4,opt,name=Qos" json:"Qos,omitempty"`
	TraceID    string `protobuf:"bytes,5,opt,name=TraceID" json:"TraceID,omitempty"`
	BizID      []byte `protobuf:"bytes,6,opt,name=BizID,proto3" json:"BizID,omitempty"`
	CreateTime int64  `protobuf:"varint,7,opt,name=CreateTime" json:"CreateTime,omitempty"`
	ExpireAt   int64  `protobuf:"varint,8,opt,name=ExpireAt" json:"ExpireAt,omitempty"`
	MessageID  int64  `protobuf:"varint,9,opt,name=MessageID" json:"MessageID,omitempty"`
}

func (m *Message) Reset()                    { *m = Message{} }
func (m *Message) String() string            { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()               {}
func (*Message) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Message) GetTopic() string {
	if m != nil {
		return m.Topic
	}
	return ""
}

func (m *Message) GetIndex() []byte {
	if m != nil {
		return m.Index
	}
	return nil
}

func (m *Message) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Message) GetQos() int32 {
	if m != nil {
		return m.Qos
	}
	return 0
}

func (m *Message) GetTraceID() string {
	if m != nil {
		return m.TraceID
	}
	return ""
}

func (m *Message) GetBizID() []byte {
	if m != nil {
		return m.BizID
	}
	return nil
}

func (m *Message) GetCreateTime() int64 {
	if m != nil {
		return m.CreateTime
	}
	return 0
}

func (m *Message) GetExpireAt() int64 {
	if m != nil {
		return m.ExpireAt
	}
	return 0
}

func (m *Message) GetMessageID() int64 {
	if m != nil {
		return m.MessageID
	}
	return 0
}

func init() {
	proto.RegisterType((*Message)(nil), "message.Message")
}

func init() { proto.RegisterFile("message.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 204 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x4c, 0x8f, 0x41, 0x4a, 0xc6, 0x30,
	0x10, 0x85, 0x89, 0xb5, 0x4d, 0x3b, 0x28, 0x48, 0x70, 0x31, 0x88, 0x48, 0x70, 0x95, 0x95, 0x1b,
	0x4f, 0xa0, 0xd6, 0x45, 0x16, 0x82, 0x86, 0x5e, 0x20, 0xb6, 0x83, 0x04, 0xac, 0x29, 0x69, 0x17,
	0xd5, 0x2b, 0x7b, 0x09, 0x49, 0xd2, 0xfa, 0xff, 0xbb, 0xf9, 0xbe, 0x37, 0x3c, 0x78, 0x70, 0x3e,
	0xd2, 0x3c, 0xdb, 0x0f, 0xba, 0x9b, 0x82, 0x5f, 0xbc, 0xe0, 0x1b, 0xde, 0xfe, 0x32, 0xe0, 0x2f,
	0xf9, 0x16, 0x97, 0x50, 0x76, 0x7e, 0x72, 0x3d, 0x32, 0xc9, 0x54, 0x63, 0x32, 0x44, 0xab, 0xbf,
	0x06, 0x5a, 0xf1, 0x44, 0x32, 0x75, 0x66, 0x32, 0x08, 0x04, 0xfe, 0x6a, 0xbf, 0x3f, 0xbd, 0x1d,
	0xb0, 0x48, 0x7e, 0x47, 0x71, 0x01, 0xc5, 0x9b, 0x9f, 0xf1, 0x54, 0x32, 0x55, 0x9a, 0x78, 0xc6,
	0xdf, 0x2e, 0xd8, 0x9e, 0x74, 0x8b, 0x65, 0x6a, 0xde, 0x31, 0x76, 0x3f, 0xba, 0x1f, 0xdd, 0x62,
	0x95, 0xbb, 0x13, 0x88, 0x1b, 0x80, 0xa7, 0x40, 0x76, 0xa1, 0xce, 0x8d, 0x84, 0x5c, 0x32, 0x55,
	0x98, 0x23, 0x23, 0xae, 0xa0, 0x7e, 0x5e, 0x27, 0x17, 0xe8, 0x61, 0xc1, 0x3a, 0xa5, 0xff, 0x2c,
	0xae, 0xa1, 0xd9, 0xe6, 0xe8, 0x16, 0x9b, 0x14, 0x1e, 0xc4, 0x7b, 0x95, 0xd6, 0xdf, 0xff, 0x05,
	0x00, 0x00, 0xff, 0xff, 0x47, 0x44, 0x95, 0xb9, 0x0e, 0x01, 0x00, 0x00,
}
