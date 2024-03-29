// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.19.4
// source: user/user.proto

package model

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

type Reply_Gender int32

const (
	Reply_MALE   Reply_Gender = 0
	Reply_BOY    Reply_Gender = 0
	Reply_FEMALE Reply_Gender = 1
	Reply_GIRL   Reply_Gender = 1
	Reply_OTHER  Reply_Gender = 2
)

// Enum value maps for Reply_Gender.
var (
	Reply_Gender_name = map[int32]string{
		0: "MALE",
		// Duplicate value: 0: "BOY",
		1: "FEMALE",
		// Duplicate value: 1: "GIRL",
		2: "OTHER",
	}
	Reply_Gender_value = map[string]int32{
		"MALE":   0,
		"BOY":    0,
		"FEMALE": 1,
		"GIRL":   1,
		"OTHER":  2,
	}
)

func (x Reply_Gender) Enum() *Reply_Gender {
	p := new(Reply_Gender)
	*p = x
	return p
}

func (x Reply_Gender) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Reply_Gender) Descriptor() protoreflect.EnumDescriptor {
	return file_user_user_proto_enumTypes[0].Descriptor()
}

func (Reply_Gender) Type() protoreflect.EnumType {
	return &file_user_user_proto_enumTypes[0]
}

func (x Reply_Gender) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Reply_Gender.Descriptor instead.
func (Reply_Gender) EnumDescriptor() ([]byte, []int) {
	return file_user_user_proto_rawDescGZIP(), []int{2, 0}
}

// 参数1
type Args struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    int64  `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Param string `protobuf:"bytes,2,opt,name=Param,proto3" json:"Param,omitempty"`
}

func (x *Args) Reset() {
	*x = Args{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Args) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Args) ProtoMessage() {}

func (x *Args) ProtoReflect() protoreflect.Message {
	mi := &file_user_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Args.ProtoReflect.Descriptor instead.
func (*Args) Descriptor() ([]byte, []int) {
	return file_user_user_proto_rawDescGZIP(), []int{0}
}

func (x *Args) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Args) GetParam() string {
	if x != nil {
		return x.Param
	}
	return ""
}

// 参数2
type Args2 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    int64  `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Param string `protobuf:"bytes,2,opt,name=Param,proto3" json:"Param,omitempty"`
}

func (x *Args2) Reset() {
	*x = Args2{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Args2) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Args2) ProtoMessage() {}

func (x *Args2) ProtoReflect() protoreflect.Message {
	mi := &file_user_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Args2.ProtoReflect.Descriptor instead.
func (*Args2) Descriptor() ([]byte, []int) {
	return file_user_user_proto_rawDescGZIP(), []int{1}
}

func (x *Args2) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Args2) GetParam() string {
	if x != nil {
		return x.Param
	}
	return ""
}

// 返回值
type Reply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    int64             `protobuf:"varint,1,opt,name=Code,proto3" json:"Code,omitempty"`
	Message string            `protobuf:"bytes,2,opt,name=Message,proto3" json:"Message,omitempty"`
	Gender  Reply_Gender      `protobuf:"varint,3,opt,name=gender,proto3,enum=model.Reply_Gender" json:"gender,omitempty"`
	List    []*Reply_Data     `protobuf:"bytes,4,rep,name=list,proto3" json:"list,omitempty"`
	Detail  map[string]string `protobuf:"bytes,5,rep,name=detail,proto3" json:"detail,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Reply) Reset() {
	*x = Reply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_user_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Reply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Reply) ProtoMessage() {}

func (x *Reply) ProtoReflect() protoreflect.Message {
	mi := &file_user_user_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Reply.ProtoReflect.Descriptor instead.
func (*Reply) Descriptor() ([]byte, []int) {
	return file_user_user_proto_rawDescGZIP(), []int{2}
}

func (x *Reply) GetCode() int64 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Reply) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *Reply) GetGender() Reply_Gender {
	if x != nil {
		return x.Gender
	}
	return Reply_MALE
}

func (x *Reply) GetList() []*Reply_Data {
	if x != nil {
		return x.List
	}
	return nil
}

func (x *Reply) GetDetail() map[string]string {
	if x != nil {
		return x.Detail
	}
	return nil
}

type Reply_Data struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Model string `protobuf:"bytes,1,opt,name=model,proto3" json:"model,omitempty"`
	Brand string `protobuf:"bytes,2,opt,name=brand,proto3" json:"brand,omitempty"`
}

func (x *Reply_Data) Reset() {
	*x = Reply_Data{}
	if protoimpl.UnsafeEnabled {
		mi := &file_user_user_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Reply_Data) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Reply_Data) ProtoMessage() {}

func (x *Reply_Data) ProtoReflect() protoreflect.Message {
	mi := &file_user_user_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Reply_Data.ProtoReflect.Descriptor instead.
func (*Reply_Data) Descriptor() ([]byte, []int) {
	return file_user_user_proto_rawDescGZIP(), []int{2, 0}
}

func (x *Reply_Data) GetModel() string {
	if x != nil {
		return x.Model
	}
	return ""
}

func (x *Reply_Data) GetBrand() string {
	if x != nil {
		return x.Brand
	}
	return ""
}

var File_user_user_proto protoreflect.FileDescriptor

var file_user_user_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x05, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x22, 0x2c, 0x0a, 0x04, 0x41, 0x72, 0x67, 0x73,
	0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x64,
	0x12, 0x14, 0x0a, 0x05, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x22, 0x2d, 0x0a, 0x05, 0x41, 0x72, 0x67, 0x73, 0x32, 0x12,
	0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x49, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x50, 0x61, 0x72, 0x61, 0x6d, 0x22, 0xec, 0x02, 0x0a, 0x05, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12,
	0x12, 0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x43,
	0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2b, 0x0a,
	0x06, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e,
	0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x2e, 0x47, 0x65, 0x6e, 0x64,
	0x65, 0x72, 0x52, 0x06, 0x67, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x25, 0x0a, 0x04, 0x6c, 0x69,
	0x73, 0x74, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x2e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x52, 0x04, 0x6c, 0x69, 0x73,
	0x74, 0x12, 0x30, 0x0a, 0x06, 0x64, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x18, 0x05, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x18, 0x2e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x2e,
	0x44, 0x65, 0x74, 0x61, 0x69, 0x6c, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x64, 0x65, 0x74,
	0x61, 0x69, 0x6c, 0x1a, 0x32, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x12, 0x14, 0x0a, 0x05, 0x6d,
	0x6f, 0x64, 0x65, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6d, 0x6f, 0x64, 0x65,
	0x6c, 0x12, 0x14, 0x0a, 0x05, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x05, 0x62, 0x72, 0x61, 0x6e, 0x64, 0x1a, 0x39, 0x0a, 0x0b, 0x44, 0x65, 0x74, 0x61, 0x69,
	0x6c, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02,
	0x38, 0x01, 0x22, 0x40, 0x0a, 0x06, 0x47, 0x65, 0x6e, 0x64, 0x65, 0x72, 0x12, 0x08, 0x0a, 0x04,
	0x4d, 0x41, 0x4c, 0x45, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03, 0x42, 0x4f, 0x59, 0x10, 0x00, 0x12,
	0x0a, 0x0a, 0x06, 0x46, 0x45, 0x4d, 0x41, 0x4c, 0x45, 0x10, 0x01, 0x12, 0x08, 0x0a, 0x04, 0x47,
	0x49, 0x52, 0x4c, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x4f, 0x54, 0x48, 0x45, 0x52, 0x10, 0x02,
	0x1a, 0x02, 0x10, 0x01, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2f, 0x3b, 0x6d, 0x6f, 0x64, 0x65, 0x6c,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_user_user_proto_rawDescOnce sync.Once
	file_user_user_proto_rawDescData = file_user_user_proto_rawDesc
)

func file_user_user_proto_rawDescGZIP() []byte {
	file_user_user_proto_rawDescOnce.Do(func() {
		file_user_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_user_user_proto_rawDescData)
	})
	return file_user_user_proto_rawDescData
}

var file_user_user_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_user_user_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_user_user_proto_goTypes = []interface{}{
	(Reply_Gender)(0),  // 0: model.Reply.Gender
	(*Args)(nil),       // 1: model.Args
	(*Args2)(nil),      // 2: model.Args2
	(*Reply)(nil),      // 3: model.Reply
	(*Reply_Data)(nil), // 4: model.Reply.Data
	nil,                // 5: model.Reply.DetailEntry
}
var file_user_user_proto_depIdxs = []int32{
	0, // 0: model.Reply.gender:type_name -> model.Reply.Gender
	4, // 1: model.Reply.list:type_name -> model.Reply.Data
	5, // 2: model.Reply.detail:type_name -> model.Reply.DetailEntry
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_user_user_proto_init() }
func file_user_user_proto_init() {
	if File_user_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_user_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Args); i {
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
		file_user_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Args2); i {
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
		file_user_user_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Reply); i {
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
		file_user_user_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Reply_Data); i {
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
			RawDescriptor: file_user_user_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_user_user_proto_goTypes,
		DependencyIndexes: file_user_user_proto_depIdxs,
		EnumInfos:         file_user_user_proto_enumTypes,
		MessageInfos:      file_user_user_proto_msgTypes,
	}.Build()
	File_user_user_proto = out.File
	file_user_user_proto_rawDesc = nil
	file_user_user_proto_goTypes = nil
	file_user_user_proto_depIdxs = nil
}
