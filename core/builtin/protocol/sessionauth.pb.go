// Code generated by protoc-gen-go.
// source: protocol/sessionauth.proto
// DO NOT EDIT!

package protocol

import proto "code.google.com/p/goprotobuf/proto"
import json "encoding/json"
import math "math"

// Reference proto, json, and math imports to suppress error if they are not otherwise used.
var _ = proto.Marshal
var _ = &json.SyntaxError{}
var _ = math.Inf

type SSPacketAuth struct {
	AuthKey          *string `protobuf:"bytes,1,req" json:"AuthKey,omitempty"`
	Timestamp        *int64  `protobuf:"varint,2,req" json:"Timestamp,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *SSPacketAuth) Reset()         { *m = SSPacketAuth{} }
func (m *SSPacketAuth) String() string { return proto.CompactTextString(m) }
func (*SSPacketAuth) ProtoMessage()    {}

func (m *SSPacketAuth) GetAuthKey() string {
	if m != nil && m.AuthKey != nil {
		return *m.AuthKey
	}
	return ""
}

func (m *SSPacketAuth) GetTimestamp() int64 {
	if m != nil && m.Timestamp != nil {
		return *m.Timestamp
	}
	return 0
}

func init() {
}
