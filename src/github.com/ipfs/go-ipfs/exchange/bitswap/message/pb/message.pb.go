// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message.proto

/*
Package bitswap_message_pb is a generated protocol buffer package.

It is generated from these files:
	message.proto

It has these top-level messages:
	Message
*/
package bitswap_message_pb

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
	Wantlist         *Message_Wantlist `protobuf:"bytes,1,opt,name=wantlist" json:"wantlist,omitempty"`
	Blocks           [][]byte          `protobuf:"bytes,2,rep,name=blocks" json:"blocks,omitempty"`
	Payload          []*Message_Block  `protobuf:"bytes,3,rep,name=payload" json:"payload,omitempty"`
	XXX_unrecognized []byte            `json:"-"`
}

func (m *Message) Reset()                    { *m = Message{} }
func (m *Message) String() string            { return proto.CompactTextString(m) }
func (*Message) ProtoMessage()               {}
func (*Message) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *Message) GetWantlist() *Message_Wantlist {
	if m != nil {
		return m.Wantlist
	}
	return nil
}

func (m *Message) GetBlocks() [][]byte {
	if m != nil {
		return m.Blocks
	}
	return nil
}

func (m *Message) GetPayload() []*Message_Block {
	if m != nil {
		return m.Payload
	}
	return nil
}

type Message_Wantlist struct {
	Entries          []*Message_Wantlist_Entry `protobuf:"bytes,1,rep,name=entries" json:"entries,omitempty"`
	Full             *bool                     `protobuf:"varint,2,opt,name=full" json:"full,omitempty"`
	XXX_unrecognized []byte                    `json:"-"`
}

func (m *Message_Wantlist) Reset()                    { *m = Message_Wantlist{} }
func (m *Message_Wantlist) String() string            { return proto.CompactTextString(m) }
func (*Message_Wantlist) ProtoMessage()               {}
func (*Message_Wantlist) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

func (m *Message_Wantlist) GetEntries() []*Message_Wantlist_Entry {
	if m != nil {
		return m.Entries
	}
	return nil
}

func (m *Message_Wantlist) GetFull() bool {
	if m != nil && m.Full != nil {
		return *m.Full
	}
	return false
}

type Message_Wantlist_Entry struct {
	Block            *string `protobuf:"bytes,1,opt,name=block" json:"block,omitempty"`
	Priority         *int32  `protobuf:"varint,2,opt,name=priority" json:"priority,omitempty"`
	Cancel           *bool   `protobuf:"varint,3,opt,name=cancel" json:"cancel,omitempty"`
	Delegate         *bool   `protobuf:"varint,4,opt,name=delegate" json:"delegate,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *Message_Wantlist_Entry) Reset()                    { *m = Message_Wantlist_Entry{} }
func (m *Message_Wantlist_Entry) String() string            { return proto.CompactTextString(m) }
func (*Message_Wantlist_Entry) ProtoMessage()               {}
func (*Message_Wantlist_Entry) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0, 0} }

func (m *Message_Wantlist_Entry) GetBlock() string {
	if m != nil && m.Block != nil {
		return *m.Block
	}
	return ""
}

func (m *Message_Wantlist_Entry) GetPriority() int32 {
	if m != nil && m.Priority != nil {
		return *m.Priority
	}
	return 0
}

func (m *Message_Wantlist_Entry) GetCancel() bool {
	if m != nil && m.Cancel != nil {
		return *m.Cancel
	}
	return false
}

func (m *Message_Wantlist_Entry) GetDelegate() bool {
	if m != nil && m.Delegate != nil {
		return *m.Delegate
	}
	return false
}

type Message_Block struct {
	Prefix           []byte `protobuf:"bytes,1,opt,name=prefix" json:"prefix,omitempty"`
	Data             []byte `protobuf:"bytes,2,opt,name=data" json:"data,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *Message_Block) Reset()                    { *m = Message_Block{} }
func (m *Message_Block) String() string            { return proto.CompactTextString(m) }
func (*Message_Block) ProtoMessage()               {}
func (*Message_Block) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 1} }

func (m *Message_Block) GetPrefix() []byte {
	if m != nil {
		return m.Prefix
	}
	return nil
}

func (m *Message_Block) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*Message)(nil), "bitswap.message.pb.Message")
	proto.RegisterType((*Message_Wantlist)(nil), "bitswap.message.pb.Message.Wantlist")
	proto.RegisterType((*Message_Wantlist_Entry)(nil), "bitswap.message.pb.Message.Wantlist.Entry")
	proto.RegisterType((*Message_Block)(nil), "bitswap.message.pb.Message.Block")
}

func init() { proto.RegisterFile("message.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 248 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x8f, 0x31, 0x4b, 0xc4, 0x30,
	0x18, 0x86, 0x49, 0x7b, 0xb5, 0xf5, 0xbb, 0x9e, 0x48, 0xa6, 0x70, 0x53, 0x15, 0x85, 0xe2, 0x90,
	0xa1, 0x83, 0x8b, 0x9b, 0xe0, 0x22, 0x38, 0x3b, 0xa7, 0xed, 0x77, 0x47, 0x30, 0x36, 0x21, 0xf9,
	0xe4, 0xec, 0xcf, 0xf1, 0x1f, 0xf9, 0x93, 0x24, 0xb9, 0xa3, 0x8b, 0x20, 0x8e, 0xf9, 0x78, 0xde,
	0x27, 0xef, 0x0b, 0x9b, 0x77, 0x0c, 0x41, 0xed, 0x51, 0x3a, 0x6f, 0xc9, 0x72, 0xde, 0x6b, 0x0a,
	0x07, 0xe5, 0xe4, 0x72, 0xee, 0xaf, 0xbf, 0x33, 0x28, 0x5f, 0x8e, 0x4f, 0x7e, 0x0f, 0xd5, 0x41,
	0x4d, 0x64, 0x74, 0x20, 0xc1, 0x1a, 0xd6, 0xae, 0xbb, 0x1b, 0xf9, 0x3b, 0x22, 0x4f, 0xb8, 0x7c,
	0x3d, 0xb1, 0xfc, 0x02, 0xce, 0x7a, 0x63, 0x87, 0xb7, 0x20, 0xb2, 0x26, 0x6f, 0x6b, 0xde, 0x41,
	0xe9, 0xd4, 0x6c, 0xac, 0x1a, 0x45, 0xde, 0xe4, 0xed, 0xba, 0xbb, 0xfa, 0x4b, 0xf3, 0x18, 0xa3,
	0xdb, 0x2f, 0x06, 0xd5, 0x22, 0x7c, 0x80, 0x12, 0x27, 0xf2, 0x1a, 0x83, 0x60, 0x49, 0x70, 0xf7,
	0x9f, 0x1e, 0xf2, 0x69, 0x22, 0x3f, 0xf3, 0x1a, 0x56, 0xbb, 0x0f, 0x63, 0x44, 0xd6, 0xb0, 0xb6,
	0xda, 0x3e, 0x43, 0x71, 0x3c, 0x6f, 0xa0, 0x48, 0x25, 0xd3, 0xb2, 0x73, 0x7e, 0x09, 0x95, 0xf3,
	0xda, 0x7a, 0x4d, 0x73, 0x22, 0x8b, 0xb8, 0x62, 0x50, 0xd3, 0x80, 0x46, 0xe4, 0x31, 0x19, 0x89,
	0x11, 0x0d, 0xee, 0x15, 0xa1, 0x58, 0x25, 0xd7, 0x2d, 0x14, 0xa9, 0x6c, 0x44, 0x9d, 0xc7, 0x9d,
	0xfe, 0x4c, 0xb2, 0x3a, 0x7e, 0x39, 0x2a, 0x52, 0x49, 0x54, 0xff, 0x04, 0x00, 0x00, 0xff, 0xff,
	0x6e, 0xe5, 0xae, 0xe2, 0x76, 0x01, 0x00, 0x00,
}
