// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.12.4
// source: core.proto

package panda

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type KeyExchange_Status int32

const (
	KeyExchange_INIT      KeyExchange_Status = 0
	KeyExchange_EXCHANGE1 KeyExchange_Status = 1
	KeyExchange_EXCHANGE2 KeyExchange_Status = 2
)

// Enum value maps for KeyExchange_Status.
var (
	KeyExchange_Status_name = map[int32]string{
		0: "INIT",
		1: "EXCHANGE1",
		2: "EXCHANGE2",
	}
	KeyExchange_Status_value = map[string]int32{
		"INIT":      0,
		"EXCHANGE1": 1,
		"EXCHANGE2": 2,
	}
)

func (x KeyExchange_Status) Enum() *KeyExchange_Status {
	p := new(KeyExchange_Status)
	*p = x
	return p
}

func (x KeyExchange_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (KeyExchange_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_core_proto_enumTypes[0].Descriptor()
}

func (KeyExchange_Status) Type() protoreflect.EnumType {
	return &file_core_proto_enumTypes[0]
}

func (x KeyExchange_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Do not use.
func (x *KeyExchange_Status) UnmarshalJSON(b []byte) error {
	num, err := protoimpl.X.UnmarshalJSONEnum(x.Descriptor(), b)
	if err != nil {
		return err
	}
	*x = KeyExchange_Status(num)
	return nil
}

// Deprecated: Use KeyExchange_Status.Descriptor instead.
func (KeyExchange_Status) EnumDescriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{0, 0}
}

type KeyExchange struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Status           *KeyExchange_Status `protobuf:"varint,1,req,name=status,enum=panda.KeyExchange_Status" json:"status,omitempty"`
	KeyExchangeBytes []byte              `protobuf:"bytes,2,req,name=key_exchange_bytes,json=keyExchangeBytes" json:"key_exchange_bytes,omitempty"`
	SharedSecret     []byte              `protobuf:"bytes,3,opt,name=shared_secret,json=sharedSecret" json:"shared_secret,omitempty"`
	DhPrivate        []byte              `protobuf:"bytes,4,opt,name=dh_private,json=dhPrivate" json:"dh_private,omitempty"`
	Key              []byte              `protobuf:"bytes,5,opt,name=key" json:"key,omitempty"`
	Meeting1         []byte              `protobuf:"bytes,6,opt,name=meeting1" json:"meeting1,omitempty"`
	Meeting2         []byte              `protobuf:"bytes,7,opt,name=meeting2" json:"meeting2,omitempty"`
	Message1         []byte              `protobuf:"bytes,8,opt,name=message1" json:"message1,omitempty"`
	Message2         []byte              `protobuf:"bytes,9,opt,name=message2" json:"message2,omitempty"`
	SharedKey        []byte              `protobuf:"bytes,10,opt,name=shared_key,json=sharedKey" json:"shared_key,omitempty"`
	SharedRandom     []byte              `protobuf:"bytes,11,opt,name=shared_random,json=sharedRandom" json:"shared_random,omitempty"`
}

func (x *KeyExchange) Reset() {
	*x = KeyExchange{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KeyExchange) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KeyExchange) ProtoMessage() {}

func (x *KeyExchange) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KeyExchange.ProtoReflect.Descriptor instead.
func (*KeyExchange) Descriptor() ([]byte, []int) {
	return file_core_proto_rawDescGZIP(), []int{0}
}

func (x *KeyExchange) GetStatus() KeyExchange_Status {
	if x != nil && x.Status != nil {
		return *x.Status
	}
	return KeyExchange_INIT
}

func (x *KeyExchange) GetKeyExchangeBytes() []byte {
	if x != nil {
		return x.KeyExchangeBytes
	}
	return nil
}

func (x *KeyExchange) GetSharedSecret() []byte {
	if x != nil {
		return x.SharedSecret
	}
	return nil
}

func (x *KeyExchange) GetDhPrivate() []byte {
	if x != nil {
		return x.DhPrivate
	}
	return nil
}

func (x *KeyExchange) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *KeyExchange) GetMeeting1() []byte {
	if x != nil {
		return x.Meeting1
	}
	return nil
}

func (x *KeyExchange) GetMeeting2() []byte {
	if x != nil {
		return x.Meeting2
	}
	return nil
}

func (x *KeyExchange) GetMessage1() []byte {
	if x != nil {
		return x.Message1
	}
	return nil
}

func (x *KeyExchange) GetMessage2() []byte {
	if x != nil {
		return x.Message2
	}
	return nil
}

func (x *KeyExchange) GetSharedKey() []byte {
	if x != nil {
		return x.SharedKey
	}
	return nil
}

func (x *KeyExchange) GetSharedRandom() []byte {
	if x != nil {
		return x.SharedRandom
	}
	return nil
}

var File_core_proto protoreflect.FileDescriptor

var file_core_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x61,
	0x6e, 0x64, 0x61, 0x22, 0xaa, 0x03, 0x0a, 0x0b, 0x4b, 0x65, 0x79, 0x45, 0x78, 0x63, 0x68, 0x61,
	0x6e, 0x67, 0x65, 0x12, 0x31, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x01, 0x20,
	0x02, 0x28, 0x0e, 0x32, 0x19, 0x2e, 0x70, 0x61, 0x6e, 0x64, 0x61, 0x2e, 0x4b, 0x65, 0x79, 0x45,
	0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06,
	0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x2c, 0x0a, 0x12, 0x6b, 0x65, 0x79, 0x5f, 0x65, 0x78,
	0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18, 0x02, 0x20, 0x02,
	0x28, 0x0c, 0x52, 0x10, 0x6b, 0x65, 0x79, 0x45, 0x78, 0x63, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x42,
	0x79, 0x74, 0x65, 0x73, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x5f, 0x73,
	0x65, 0x63, 0x72, 0x65, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0c, 0x73, 0x68, 0x61,
	0x72, 0x65, 0x64, 0x53, 0x65, 0x63, 0x72, 0x65, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x64, 0x68, 0x5f,
	0x70, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09, 0x64,
	0x68, 0x50, 0x72, 0x69, 0x76, 0x61, 0x74, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65,
	0x65, 0x74, 0x69, 0x6e, 0x67, 0x31, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x6d, 0x65,
	0x65, 0x74, 0x69, 0x6e, 0x67, 0x31, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x65, 0x74, 0x69, 0x6e,
	0x67, 0x32, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x6d, 0x65, 0x65, 0x74, 0x69, 0x6e,
	0x67, 0x32, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x31, 0x18, 0x08,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x31, 0x12, 0x1a,
	0x0a, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x08, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x68,
	0x61, 0x72, 0x65, 0x64, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x09,
	0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x4b, 0x65, 0x79, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x68, 0x61,
	0x72, 0x65, 0x64, 0x5f, 0x72, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x0c, 0x73, 0x68, 0x61, 0x72, 0x65, 0x64, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x22, 0x30,
	0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x08, 0x0a, 0x04, 0x49, 0x4e, 0x49, 0x54,
	0x10, 0x00, 0x12, 0x0d, 0x0a, 0x09, 0x45, 0x58, 0x43, 0x48, 0x41, 0x4e, 0x47, 0x45, 0x31, 0x10,
	0x01, 0x12, 0x0d, 0x0a, 0x09, 0x45, 0x58, 0x43, 0x48, 0x41, 0x4e, 0x47, 0x45, 0x32, 0x10, 0x02,
	0x42, 0x1d, 0x5a, 0x1b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b,
	0x61, 0x74, 0x7a, 0x65, 0x6e, 0x70, 0x6f, 0x73, 0x74, 0x2f, 0x70, 0x61, 0x6e, 0x64, 0x61,
}

var (
	file_core_proto_rawDescOnce sync.Once
	file_core_proto_rawDescData = file_core_proto_rawDesc
)

func file_core_proto_rawDescGZIP() []byte {
	file_core_proto_rawDescOnce.Do(func() {
		file_core_proto_rawDescData = protoimpl.X.CompressGZIP(file_core_proto_rawDescData)
	})
	return file_core_proto_rawDescData
}

var file_core_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_core_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_core_proto_goTypes = []interface{}{
	(KeyExchange_Status)(0), // 0: panda.KeyExchange.Status
	(*KeyExchange)(nil),     // 1: panda.KeyExchange
}
var file_core_proto_depIdxs = []int32{
	0, // 0: panda.KeyExchange.status:type_name -> panda.KeyExchange.Status
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_core_proto_init() }
func file_core_proto_init() {
	if File_core_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_core_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KeyExchange); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_core_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_core_proto_goTypes,
		DependencyIndexes: file_core_proto_depIdxs,
		EnumInfos:         file_core_proto_enumTypes,
		MessageInfos:      file_core_proto_msgTypes,
	}.Build()
	File_core_proto = out.File
	file_core_proto_rawDesc = nil
	file_core_proto_goTypes = nil
	file_core_proto_depIdxs = nil
}