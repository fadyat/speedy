// Code generated by protoc-gen-go. DO NOT EDIT.
// source: cache_service.proto

package api

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	_ "github.com/golang/protobuf/ptypes/empty"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// get/put messages
type GetRequest struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetRequest) Reset()         { *m = GetRequest{} }
func (m *GetRequest) String() string { return proto.CompactTextString(m) }
func (*GetRequest) ProtoMessage()    {}
func (*GetRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd22456f3fa912cc, []int{0}
}

func (m *GetRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetRequest.Unmarshal(m, b)
}
func (m *GetRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetRequest.Marshal(b, m, deterministic)
}
func (m *GetRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetRequest.Merge(m, src)
}
func (m *GetRequest) XXX_Size() int {
	return xxx_messageInfo_GetRequest.Size(m)
}
func (m *GetRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetRequest proto.InternalMessageInfo

func (m *GetRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type GetResponse struct {
	Value                string   `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetResponse) Reset()         { *m = GetResponse{} }
func (m *GetResponse) String() string { return proto.CompactTextString(m) }
func (*GetResponse) ProtoMessage()    {}
func (*GetResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd22456f3fa912cc, []int{1}
}

func (m *GetResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetResponse.Unmarshal(m, b)
}
func (m *GetResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetResponse.Marshal(b, m, deterministic)
}
func (m *GetResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetResponse.Merge(m, src)
}
func (m *GetResponse) XXX_Size() int {
	return xxx_messageInfo_GetResponse.Size(m)
}
func (m *GetResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GetResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GetResponse proto.InternalMessageInfo

func (m *GetResponse) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type PutRequest struct {
	Key                  string   `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value                string   `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PutRequest) Reset()         { *m = PutRequest{} }
func (m *PutRequest) String() string { return proto.CompactTextString(m) }
func (*PutRequest) ProtoMessage()    {}
func (*PutRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd22456f3fa912cc, []int{2}
}

func (m *PutRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PutRequest.Unmarshal(m, b)
}
func (m *PutRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PutRequest.Marshal(b, m, deterministic)
}
func (m *PutRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PutRequest.Merge(m, src)
}
func (m *PutRequest) XXX_Size() int {
	return xxx_messageInfo_PutRequest.Size(m)
}
func (m *PutRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PutRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PutRequest proto.InternalMessageInfo

func (m *PutRequest) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *PutRequest) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

type LengthResponse struct {
	Length               uint32   `protobuf:"varint,1,opt,name=length,proto3" json:"length,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LengthResponse) Reset()         { *m = LengthResponse{} }
func (m *LengthResponse) String() string { return proto.CompactTextString(m) }
func (*LengthResponse) ProtoMessage()    {}
func (*LengthResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd22456f3fa912cc, []int{3}
}

func (m *LengthResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LengthResponse.Unmarshal(m, b)
}
func (m *LengthResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LengthResponse.Marshal(b, m, deterministic)
}
func (m *LengthResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LengthResponse.Merge(m, src)
}
func (m *LengthResponse) XXX_Size() int {
	return xxx_messageInfo_LengthResponse.Size(m)
}
func (m *LengthResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_LengthResponse.DiscardUnknown(m)
}

var xxx_messageInfo_LengthResponse proto.InternalMessageInfo

func (m *LengthResponse) GetLength() uint32 {
	if m != nil {
		return m.Length
	}
	return 0
}

// election messages
type ElectionRequest struct {
	CallerPid            int32    `protobuf:"varint,1,opt,name=caller_pid,json=callerPid,proto3" json:"caller_pid,omitempty"`
	CallerNodeId         string   `protobuf:"bytes,2,opt,name=caller_node_id,json=callerNodeId,proto3" json:"caller_node_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ElectionRequest) Reset()         { *m = ElectionRequest{} }
func (m *ElectionRequest) String() string { return proto.CompactTextString(m) }
func (*ElectionRequest) ProtoMessage()    {}
func (*ElectionRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd22456f3fa912cc, []int{4}
}

func (m *ElectionRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ElectionRequest.Unmarshal(m, b)
}
func (m *ElectionRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ElectionRequest.Marshal(b, m, deterministic)
}
func (m *ElectionRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ElectionRequest.Merge(m, src)
}
func (m *ElectionRequest) XXX_Size() int {
	return xxx_messageInfo_ElectionRequest.Size(m)
}
func (m *ElectionRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ElectionRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ElectionRequest proto.InternalMessageInfo

func (m *ElectionRequest) GetCallerPid() int32 {
	if m != nil {
		return m.CallerPid
	}
	return 0
}

func (m *ElectionRequest) GetCallerNodeId() string {
	if m != nil {
		return m.CallerNodeId
	}
	return ""
}

type PidRequest struct {
	CallerPid            int32    `protobuf:"varint,1,opt,name=caller_pid,json=callerPid,proto3" json:"caller_pid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PidRequest) Reset()         { *m = PidRequest{} }
func (m *PidRequest) String() string { return proto.CompactTextString(m) }
func (*PidRequest) ProtoMessage()    {}
func (*PidRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd22456f3fa912cc, []int{5}
}

func (m *PidRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PidRequest.Unmarshal(m, b)
}
func (m *PidRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PidRequest.Marshal(b, m, deterministic)
}
func (m *PidRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PidRequest.Merge(m, src)
}
func (m *PidRequest) XXX_Size() int {
	return xxx_messageInfo_PidRequest.Size(m)
}
func (m *PidRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PidRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PidRequest proto.InternalMessageInfo

func (m *PidRequest) GetCallerPid() int32 {
	if m != nil {
		return m.CallerPid
	}
	return 0
}

type PidResponse struct {
	Pid                  int32    `protobuf:"varint,1,opt,name=pid,proto3" json:"pid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PidResponse) Reset()         { *m = PidResponse{} }
func (m *PidResponse) String() string { return proto.CompactTextString(m) }
func (*PidResponse) ProtoMessage()    {}
func (*PidResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd22456f3fa912cc, []int{6}
}

func (m *PidResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PidResponse.Unmarshal(m, b)
}
func (m *PidResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PidResponse.Marshal(b, m, deterministic)
}
func (m *PidResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PidResponse.Merge(m, src)
}
func (m *PidResponse) XXX_Size() int {
	return xxx_messageInfo_PidResponse.Size(m)
}
func (m *PidResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_PidResponse.DiscardUnknown(m)
}

var xxx_messageInfo_PidResponse proto.InternalMessageInfo

func (m *PidResponse) GetPid() int32 {
	if m != nil {
		return m.Pid
	}
	return 0
}

type HeartbeatRequest struct {
	CallerNodeId         string   `protobuf:"bytes,1,opt,name=caller_node_id,json=callerNodeId,proto3" json:"caller_node_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HeartbeatRequest) Reset()         { *m = HeartbeatRequest{} }
func (m *HeartbeatRequest) String() string { return proto.CompactTextString(m) }
func (*HeartbeatRequest) ProtoMessage()    {}
func (*HeartbeatRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd22456f3fa912cc, []int{7}
}

func (m *HeartbeatRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HeartbeatRequest.Unmarshal(m, b)
}
func (m *HeartbeatRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HeartbeatRequest.Marshal(b, m, deterministic)
}
func (m *HeartbeatRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeartbeatRequest.Merge(m, src)
}
func (m *HeartbeatRequest) XXX_Size() int {
	return xxx_messageInfo_HeartbeatRequest.Size(m)
}
func (m *HeartbeatRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_HeartbeatRequest.DiscardUnknown(m)
}

var xxx_messageInfo_HeartbeatRequest proto.InternalMessageInfo

func (m *HeartbeatRequest) GetCallerNodeId() string {
	if m != nil {
		return m.CallerNodeId
	}
	return ""
}

type HeartbeatResponse struct {
	NodeId               string   `protobuf:"bytes,1,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	Status               string   `protobuf:"bytes,2,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *HeartbeatResponse) Reset()         { *m = HeartbeatResponse{} }
func (m *HeartbeatResponse) String() string { return proto.CompactTextString(m) }
func (*HeartbeatResponse) ProtoMessage()    {}
func (*HeartbeatResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd22456f3fa912cc, []int{8}
}

func (m *HeartbeatResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HeartbeatResponse.Unmarshal(m, b)
}
func (m *HeartbeatResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HeartbeatResponse.Marshal(b, m, deterministic)
}
func (m *HeartbeatResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HeartbeatResponse.Merge(m, src)
}
func (m *HeartbeatResponse) XXX_Size() int {
	return xxx_messageInfo_HeartbeatResponse.Size(m)
}
func (m *HeartbeatResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_HeartbeatResponse.DiscardUnknown(m)
}

var xxx_messageInfo_HeartbeatResponse proto.InternalMessageInfo

func (m *HeartbeatResponse) GetNodeId() string {
	if m != nil {
		return m.NodeId
	}
	return ""
}

func (m *HeartbeatResponse) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

type LeaderRequest struct {
	Caller               string   `protobuf:"bytes,1,opt,name=caller,proto3" json:"caller,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LeaderRequest) Reset()         { *m = LeaderRequest{} }
func (m *LeaderRequest) String() string { return proto.CompactTextString(m) }
func (*LeaderRequest) ProtoMessage()    {}
func (*LeaderRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd22456f3fa912cc, []int{9}
}

func (m *LeaderRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LeaderRequest.Unmarshal(m, b)
}
func (m *LeaderRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LeaderRequest.Marshal(b, m, deterministic)
}
func (m *LeaderRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LeaderRequest.Merge(m, src)
}
func (m *LeaderRequest) XXX_Size() int {
	return xxx_messageInfo_LeaderRequest.Size(m)
}
func (m *LeaderRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_LeaderRequest.DiscardUnknown(m)
}

var xxx_messageInfo_LeaderRequest proto.InternalMessageInfo

func (m *LeaderRequest) GetCaller() string {
	if m != nil {
		return m.Caller
	}
	return ""
}

type LeaderResponse struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LeaderResponse) Reset()         { *m = LeaderResponse{} }
func (m *LeaderResponse) String() string { return proto.CompactTextString(m) }
func (*LeaderResponse) ProtoMessage()    {}
func (*LeaderResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd22456f3fa912cc, []int{10}
}

func (m *LeaderResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LeaderResponse.Unmarshal(m, b)
}
func (m *LeaderResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LeaderResponse.Marshal(b, m, deterministic)
}
func (m *LeaderResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LeaderResponse.Merge(m, src)
}
func (m *LeaderResponse) XXX_Size() int {
	return xxx_messageInfo_LeaderResponse.Size(m)
}
func (m *LeaderResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_LeaderResponse.DiscardUnknown(m)
}

var xxx_messageInfo_LeaderResponse proto.InternalMessageInfo

func (m *LeaderResponse) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

type NewLeaderAnnouncement struct {
	LeaderId             string   `protobuf:"bytes,1,opt,name=leader_id,json=leaderId,proto3" json:"leader_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NewLeaderAnnouncement) Reset()         { *m = NewLeaderAnnouncement{} }
func (m *NewLeaderAnnouncement) String() string { return proto.CompactTextString(m) }
func (*NewLeaderAnnouncement) ProtoMessage()    {}
func (*NewLeaderAnnouncement) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd22456f3fa912cc, []int{11}
}

func (m *NewLeaderAnnouncement) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NewLeaderAnnouncement.Unmarshal(m, b)
}
func (m *NewLeaderAnnouncement) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NewLeaderAnnouncement.Marshal(b, m, deterministic)
}
func (m *NewLeaderAnnouncement) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NewLeaderAnnouncement.Merge(m, src)
}
func (m *NewLeaderAnnouncement) XXX_Size() int {
	return xxx_messageInfo_NewLeaderAnnouncement.Size(m)
}
func (m *NewLeaderAnnouncement) XXX_DiscardUnknown() {
	xxx_messageInfo_NewLeaderAnnouncement.DiscardUnknown(m)
}

var xxx_messageInfo_NewLeaderAnnouncement proto.InternalMessageInfo

func (m *NewLeaderAnnouncement) GetLeaderId() string {
	if m != nil {
		return m.LeaderId
	}
	return ""
}

// cluster config messages
type Node struct {
	Id                   string   `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Host                 string   `protobuf:"bytes,2,opt,name=host,proto3" json:"host,omitempty"`
	Port                 uint32   `protobuf:"varint,3,opt,name=port,proto3" json:"port,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Node) Reset()         { *m = Node{} }
func (m *Node) String() string { return proto.CompactTextString(m) }
func (*Node) ProtoMessage()    {}
func (*Node) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd22456f3fa912cc, []int{12}
}

func (m *Node) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Node.Unmarshal(m, b)
}
func (m *Node) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Node.Marshal(b, m, deterministic)
}
func (m *Node) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Node.Merge(m, src)
}
func (m *Node) XXX_Size() int {
	return xxx_messageInfo_Node.Size(m)
}
func (m *Node) XXX_DiscardUnknown() {
	xxx_messageInfo_Node.DiscardUnknown(m)
}

var xxx_messageInfo_Node proto.InternalMessageInfo

func (m *Node) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Node) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *Node) GetPort() uint32 {
	if m != nil {
		return m.Port
	}
	return 0
}

type ClusterConfigRequest struct {
	CallerNodeId         string   `protobuf:"bytes,1,opt,name=caller_node_id,json=callerNodeId,proto3" json:"caller_node_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ClusterConfigRequest) Reset()         { *m = ClusterConfigRequest{} }
func (m *ClusterConfigRequest) String() string { return proto.CompactTextString(m) }
func (*ClusterConfigRequest) ProtoMessage()    {}
func (*ClusterConfigRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd22456f3fa912cc, []int{13}
}

func (m *ClusterConfigRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClusterConfigRequest.Unmarshal(m, b)
}
func (m *ClusterConfigRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClusterConfigRequest.Marshal(b, m, deterministic)
}
func (m *ClusterConfigRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClusterConfigRequest.Merge(m, src)
}
func (m *ClusterConfigRequest) XXX_Size() int {
	return xxx_messageInfo_ClusterConfigRequest.Size(m)
}
func (m *ClusterConfigRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ClusterConfigRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ClusterConfigRequest proto.InternalMessageInfo

func (m *ClusterConfigRequest) GetCallerNodeId() string {
	if m != nil {
		return m.CallerNodeId
	}
	return ""
}

type ClusterConfig struct {
	Nodes                []*Node  `protobuf:"bytes,1,rep,name=nodes,proto3" json:"nodes,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ClusterConfig) Reset()         { *m = ClusterConfig{} }
func (m *ClusterConfig) String() string { return proto.CompactTextString(m) }
func (*ClusterConfig) ProtoMessage()    {}
func (*ClusterConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd22456f3fa912cc, []int{14}
}

func (m *ClusterConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ClusterConfig.Unmarshal(m, b)
}
func (m *ClusterConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ClusterConfig.Marshal(b, m, deterministic)
}
func (m *ClusterConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClusterConfig.Merge(m, src)
}
func (m *ClusterConfig) XXX_Size() int {
	return xxx_messageInfo_ClusterConfig.Size(m)
}
func (m *ClusterConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_ClusterConfig.DiscardUnknown(m)
}

var xxx_messageInfo_ClusterConfig proto.InternalMessageInfo

func (m *ClusterConfig) GetNodes() []*Node {
	if m != nil {
		return m.Nodes
	}
	return nil
}

type GenericResponse struct {
	Data                 string   `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GenericResponse) Reset()         { *m = GenericResponse{} }
func (m *GenericResponse) String() string { return proto.CompactTextString(m) }
func (*GenericResponse) ProtoMessage()    {}
func (*GenericResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_cd22456f3fa912cc, []int{15}
}

func (m *GenericResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GenericResponse.Unmarshal(m, b)
}
func (m *GenericResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GenericResponse.Marshal(b, m, deterministic)
}
func (m *GenericResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GenericResponse.Merge(m, src)
}
func (m *GenericResponse) XXX_Size() int {
	return xxx_messageInfo_GenericResponse.Size(m)
}
func (m *GenericResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_GenericResponse.DiscardUnknown(m)
}

var xxx_messageInfo_GenericResponse proto.InternalMessageInfo

func (m *GenericResponse) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func init() {
	proto.RegisterType((*GetRequest)(nil), "api.GetRequest")
	proto.RegisterType((*GetResponse)(nil), "api.GetResponse")
	proto.RegisterType((*PutRequest)(nil), "api.PutRequest")
	proto.RegisterType((*LengthResponse)(nil), "api.LengthResponse")
	proto.RegisterType((*ElectionRequest)(nil), "api.ElectionRequest")
	proto.RegisterType((*PidRequest)(nil), "api.PidRequest")
	proto.RegisterType((*PidResponse)(nil), "api.PidResponse")
	proto.RegisterType((*HeartbeatRequest)(nil), "api.HeartbeatRequest")
	proto.RegisterType((*HeartbeatResponse)(nil), "api.HeartbeatResponse")
	proto.RegisterType((*LeaderRequest)(nil), "api.LeaderRequest")
	proto.RegisterType((*LeaderResponse)(nil), "api.LeaderResponse")
	proto.RegisterType((*NewLeaderAnnouncement)(nil), "api.NewLeaderAnnouncement")
	proto.RegisterType((*Node)(nil), "api.Node")
	proto.RegisterType((*ClusterConfigRequest)(nil), "api.ClusterConfigRequest")
	proto.RegisterType((*ClusterConfig)(nil), "api.ClusterConfig")
	proto.RegisterType((*GenericResponse)(nil), "api.GenericResponse")
}

func init() { proto.RegisterFile("cache_service.proto", fileDescriptor_cd22456f3fa912cc) }

var fileDescriptor_cd22456f3fa912cc = []byte{
	// 644 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x54, 0xdf, 0x4f, 0x13, 0x41,
	0x10, 0xa6, 0x1c, 0x2d, 0x76, 0x28, 0xb4, 0x2e, 0x3f, 0x53, 0xa2, 0x90, 0x55, 0x23, 0x91, 0xe4,
	0x20, 0xd8, 0x07, 0x63, 0x84, 0x88, 0x48, 0x4e, 0x12, 0x42, 0x9a, 0x1a, 0x35, 0xf1, 0x85, 0x2c,
	0x77, 0x43, 0xbb, 0xf1, 0xb8, 0x3d, 0xef, 0xf6, 0x30, 0xbc, 0xfb, 0x87, 0x9b, 0xfd, 0xd1, 0xbb,
	0x5e, 0xed, 0x25, 0xfa, 0x36, 0x3b, 0xf7, 0xcd, 0xf7, 0xcd, 0xce, 0x7e, 0x73, 0xb0, 0xea, 0x33,
	0x7f, 0x84, 0xd7, 0x29, 0x26, 0xf7, 0xdc, 0x47, 0x37, 0x4e, 0x84, 0x14, 0xc4, 0x61, 0x31, 0xef,
	0x6e, 0x0f, 0x85, 0x18, 0x86, 0x78, 0xa0, 0x53, 0x37, 0xd9, 0xed, 0x01, 0xde, 0xc5, 0xf2, 0xc1,
	0x20, 0xe8, 0x53, 0x00, 0x0f, 0xe5, 0x00, 0x7f, 0x66, 0x98, 0x4a, 0xd2, 0x01, 0xe7, 0x07, 0x3e,
	0x6c, 0xd5, 0x76, 0x6b, 0x7b, 0xcd, 0x81, 0x0a, 0xe9, 0x33, 0x58, 0xd2, 0xdf, 0xd3, 0x58, 0x44,
	0x29, 0x92, 0x35, 0xa8, 0xdf, 0xb3, 0x30, 0x43, 0x0b, 0x31, 0x07, 0xda, 0x03, 0xe8, 0x67, 0xd5,
	0x24, 0x45, 0xd5, 0xfc, 0x64, 0xd5, 0x1e, 0xac, 0x5c, 0x62, 0x34, 0x94, 0xa3, 0x9c, 0x7d, 0x03,
	0x1a, 0xa1, 0xce, 0xe8, 0xe2, 0xe5, 0x81, 0x3d, 0xd1, 0xaf, 0xd0, 0x3e, 0x0f, 0xd1, 0x97, 0x5c,
	0x44, 0x63, 0x91, 0x27, 0x00, 0x3e, 0x0b, 0x43, 0x4c, 0xae, 0x63, 0x1e, 0x68, 0x78, 0x7d, 0xd0,
	0x34, 0x99, 0x3e, 0x0f, 0xc8, 0x73, 0x58, 0xb1, 0x9f, 0x23, 0x11, 0xe0, 0x35, 0x0f, 0xac, 0x74,
	0xcb, 0x64, 0xaf, 0x44, 0x80, 0x17, 0x01, 0xdd, 0x07, 0xe8, 0xf3, 0xe0, 0xdf, 0x28, 0xe9, 0x0e,
	0x2c, 0x69, 0xb0, 0xed, 0xb5, 0x03, 0x4e, 0x01, 0x53, 0x21, 0x7d, 0x03, 0x9d, 0x4f, 0xc8, 0x12,
	0x79, 0x83, 0x2c, 0x9f, 0xc5, 0xdf, 0x7d, 0xd4, 0x66, 0xf4, 0xf1, 0x11, 0x1e, 0x4f, 0x54, 0x5a,
	0x81, 0x4d, 0x58, 0x2c, 0xd7, 0x34, 0x22, 0x8d, 0x56, 0x53, 0x4a, 0x25, 0x93, 0x59, 0x6a, 0xef,
	0x64, 0x4f, 0xf4, 0x25, 0x2c, 0x5f, 0x22, 0x0b, 0x30, 0x19, 0x8b, 0x6f, 0x40, 0xc3, 0xc8, 0x8c,
	0x09, 0xcc, 0x89, 0xee, 0xaa, 0xc1, 0x1b, 0xa0, 0xd5, 0x5a, 0x81, 0xf9, 0x5c, 0x66, 0x9e, 0x07,
	0xb4, 0x07, 0xeb, 0x57, 0xf8, 0xcb, 0x80, 0x4e, 0xa3, 0x48, 0x64, 0x91, 0x8f, 0x77, 0x18, 0x49,
	0xb2, 0x0d, 0xcd, 0x50, 0x67, 0x8b, 0xb6, 0x1e, 0x99, 0xc4, 0x45, 0x40, 0x4f, 0x60, 0x41, 0x5d,
	0x68, 0x9a, 0x8d, 0x10, 0x58, 0x18, 0x89, 0x54, 0xda, 0x76, 0x75, 0xac, 0x72, 0xb1, 0x48, 0xe4,
	0x96, 0xa3, 0x1f, 0x5a, 0xc7, 0xf4, 0x1d, 0xac, 0x9d, 0x85, 0x59, 0x2a, 0x31, 0x39, 0x13, 0xd1,
	0x2d, 0x1f, 0xfe, 0xdf, 0x10, 0x0f, 0x61, 0xb9, 0x54, 0x4d, 0x76, 0xa0, 0xae, 0xf0, 0xe9, 0x56,
	0x6d, 0xd7, 0xd9, 0x5b, 0x3a, 0x6a, 0xba, 0x2c, 0xe6, 0xae, 0x02, 0x0f, 0x4c, 0x9e, 0xbe, 0x80,
	0xb6, 0x87, 0x11, 0x26, 0xdc, 0xcf, 0x07, 0x41, 0x60, 0x21, 0x60, 0x92, 0x59, 0x01, 0x1d, 0x1f,
	0xfd, 0xae, 0x43, 0xeb, 0x4c, 0x2d, 0xd7, 0x67, 0xb3, 0x5b, 0xe4, 0x15, 0x38, 0x1e, 0x4a, 0xd2,
	0xd6, 0x84, 0xc5, 0xf6, 0x74, 0x3b, 0x45, 0xc2, 0xd0, 0xd1, 0x39, 0x72, 0x08, 0x4e, 0x3f, 0x1b,
	0x63, 0x8b, 0x25, 0xe9, 0x6e, 0xb8, 0x66, 0x2b, 0xdd, 0xf1, 0x56, 0xba, 0xe7, 0x6a, 0x2b, 0xe9,
	0x1c, 0xe9, 0x81, 0x73, 0x89, 0x11, 0xa9, 0x00, 0x74, 0x57, 0x35, 0x53, 0x79, 0x71, 0xe8, 0x1c,
	0xd9, 0x87, 0x86, 0x87, 0x52, 0x59, 0xdf, 0x4a, 0xe5, 0xbe, 0xb6, 0x6d, 0x4d, 0x7a, 0xb7, 0x07,
	0x4d, 0x0f, 0xa5, 0x79, 0x5e, 0x42, 0x2c, 0xe1, 0x84, 0x73, 0x72, 0x91, 0x92, 0x49, 0x8e, 0xa1,
	0xe5, 0xa1, 0xcc, 0x8d, 0x4a, 0xd6, 0x35, 0x68, 0xda, 0xf2, 0x55, 0x37, 0x23, 0xef, 0xa1, 0xf5,
	0x25, 0x0e, 0x98, 0x44, 0xab, 0xdb, 0x35, 0xef, 0x31, 0xcb, 0x66, 0xdd, 0x35, 0x3b, 0xc9, 0xf2,
	0xe3, 0x1c, 0x43, 0xdb, 0x8a, 0x8c, 0xff, 0x06, 0xc4, 0x00, 0xa7, 0x7e, 0x0e, 0x15, 0xe5, 0x27,
	0xd0, 0xf1, 0x50, 0x96, 0x3d, 0x52, 0x35, 0x65, 0x33, 0x94, 0x32, 0xf6, 0x14, 0x56, 0xcd, 0x05,
	0xca, 0xe9, 0x19, 0xd0, 0xca, 0x19, 0xbc, 0x85, 0xcd, 0x01, 0x0e, 0xb9, 0x42, 0x2a, 0x23, 0x7e,
	0xe3, 0x72, 0x64, 0x0b, 0x49, 0x61, 0xcf, 0xd9, 0xed, 0x7f, 0x58, 0xfc, 0x5e, 0x77, 0x0f, 0x58,
	0xcc, 0x6f, 0x1a, 0x9a, 0xf4, 0xf5, 0x9f, 0x00, 0x00, 0x00, 0xff, 0xff, 0xe5, 0x81, 0xfc, 0x94,
	0xf2, 0x05, 0x00, 0x00,
}