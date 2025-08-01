// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v5.29.3
// source: account/rpc_update_account_balance.proto

package account

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type UpdateAccountBalanceRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Balance       int64                  `protobuf:"varint,1,opt,name=balance,proto3" json:"balance,omitempty"`
	Currency      string                 `protobuf:"bytes,2,opt,name=currency,proto3" json:"currency,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateAccountBalanceRequest) Reset() {
	*x = UpdateAccountBalanceRequest{}
	mi := &file_account_rpc_update_account_balance_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateAccountBalanceRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateAccountBalanceRequest) ProtoMessage() {}

func (x *UpdateAccountBalanceRequest) ProtoReflect() protoreflect.Message {
	mi := &file_account_rpc_update_account_balance_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateAccountBalanceRequest.ProtoReflect.Descriptor instead.
func (*UpdateAccountBalanceRequest) Descriptor() ([]byte, []int) {
	return file_account_rpc_update_account_balance_proto_rawDescGZIP(), []int{0}
}

func (x *UpdateAccountBalanceRequest) GetBalance() int64 {
	if x != nil {
		return x.Balance
	}
	return 0
}

func (x *UpdateAccountBalanceRequest) GetCurrency() string {
	if x != nil {
		return x.Currency
	}
	return ""
}

type UpdateAccountBalanceResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Account       *Account               `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpdateAccountBalanceResponse) Reset() {
	*x = UpdateAccountBalanceResponse{}
	mi := &file_account_rpc_update_account_balance_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpdateAccountBalanceResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateAccountBalanceResponse) ProtoMessage() {}

func (x *UpdateAccountBalanceResponse) ProtoReflect() protoreflect.Message {
	mi := &file_account_rpc_update_account_balance_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateAccountBalanceResponse.ProtoReflect.Descriptor instead.
func (*UpdateAccountBalanceResponse) Descriptor() ([]byte, []int) {
	return file_account_rpc_update_account_balance_proto_rawDescGZIP(), []int{1}
}

func (x *UpdateAccountBalanceResponse) GetAccount() *Account {
	if x != nil {
		return x.Account
	}
	return nil
}

var File_account_rpc_update_account_balance_proto protoreflect.FileDescriptor

const file_account_rpc_update_account_balance_proto_rawDesc = "" +
	"\n" +
	"(account/rpc_update_account_balance.proto\x12\x02pb\x1a\x15account/account.proto\"S\n" +
	"\x1bUpdateAccountBalanceRequest\x12\x18\n" +
	"\abalance\x18\x01 \x01(\x03R\abalance\x12\x1a\n" +
	"\bcurrency\x18\x02 \x01(\tR\bcurrency\"E\n" +
	"\x1cUpdateAccountBalanceResponse\x12%\n" +
	"\aaccount\x18\x01 \x01(\v2\v.pb.AccountR\aaccountB-Z+github.com/auronvila/simple-bank/pb/accountb\x06proto3"

var (
	file_account_rpc_update_account_balance_proto_rawDescOnce sync.Once
	file_account_rpc_update_account_balance_proto_rawDescData []byte
)

func file_account_rpc_update_account_balance_proto_rawDescGZIP() []byte {
	file_account_rpc_update_account_balance_proto_rawDescOnce.Do(func() {
		file_account_rpc_update_account_balance_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_account_rpc_update_account_balance_proto_rawDesc), len(file_account_rpc_update_account_balance_proto_rawDesc)))
	})
	return file_account_rpc_update_account_balance_proto_rawDescData
}

var file_account_rpc_update_account_balance_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_account_rpc_update_account_balance_proto_goTypes = []any{
	(*UpdateAccountBalanceRequest)(nil),  // 0: pb.UpdateAccountBalanceRequest
	(*UpdateAccountBalanceResponse)(nil), // 1: pb.UpdateAccountBalanceResponse
	(*Account)(nil),                      // 2: pb.Account
}
var file_account_rpc_update_account_balance_proto_depIdxs = []int32{
	2, // 0: pb.UpdateAccountBalanceResponse.account:type_name -> pb.Account
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_account_rpc_update_account_balance_proto_init() }
func file_account_rpc_update_account_balance_proto_init() {
	if File_account_rpc_update_account_balance_proto != nil {
		return
	}
	file_account_account_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_account_rpc_update_account_balance_proto_rawDesc), len(file_account_rpc_update_account_balance_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_account_rpc_update_account_balance_proto_goTypes,
		DependencyIndexes: file_account_rpc_update_account_balance_proto_depIdxs,
		MessageInfos:      file_account_rpc_update_account_balance_proto_msgTypes,
	}.Build()
	File_account_rpc_update_account_balance_proto = out.File
	file_account_rpc_update_account_balance_proto_goTypes = nil
	file_account_rpc_update_account_balance_proto_depIdxs = nil
}
