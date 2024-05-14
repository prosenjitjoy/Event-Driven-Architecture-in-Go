// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        (unknown)
// source: api_baskets.proto

package basketspb

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

type Item struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StoreId      string  `protobuf:"bytes,1,opt,name=store_id,json=storeId,proto3" json:"store_id,omitempty"`
	ProductId    string  `protobuf:"bytes,2,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	StoreName    string  `protobuf:"bytes,3,opt,name=store_name,json=storeName,proto3" json:"store_name,omitempty"`
	ProductName  string  `protobuf:"bytes,4,opt,name=product_name,json=productName,proto3" json:"product_name,omitempty"`
	ProductPrice float64 `protobuf:"fixed64,5,opt,name=product_price,json=productPrice,proto3" json:"product_price,omitempty"`
	Quantity     int32   `protobuf:"varint,6,opt,name=quantity,proto3" json:"quantity,omitempty"`
}

func (x *Item) Reset() {
	*x = Item{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_baskets_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Item) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Item) ProtoMessage() {}

func (x *Item) ProtoReflect() protoreflect.Message {
	mi := &file_api_baskets_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Item.ProtoReflect.Descriptor instead.
func (*Item) Descriptor() ([]byte, []int) {
	return file_api_baskets_proto_rawDescGZIP(), []int{0}
}

func (x *Item) GetStoreId() string {
	if x != nil {
		return x.StoreId
	}
	return ""
}

func (x *Item) GetProductId() string {
	if x != nil {
		return x.ProductId
	}
	return ""
}

func (x *Item) GetStoreName() string {
	if x != nil {
		return x.StoreName
	}
	return ""
}

func (x *Item) GetProductName() string {
	if x != nil {
		return x.ProductName
	}
	return ""
}

func (x *Item) GetProductPrice() float64 {
	if x != nil {
		return x.ProductPrice
	}
	return 0
}

func (x *Item) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

type Basket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Items []*Item `protobuf:"bytes,2,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *Basket) Reset() {
	*x = Basket{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_baskets_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Basket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Basket) ProtoMessage() {}

func (x *Basket) ProtoReflect() protoreflect.Message {
	mi := &file_api_baskets_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Basket.ProtoReflect.Descriptor instead.
func (*Basket) Descriptor() ([]byte, []int) {
	return file_api_baskets_proto_rawDescGZIP(), []int{1}
}

func (x *Basket) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Basket) GetItems() []*Item {
	if x != nil {
		return x.Items
	}
	return nil
}

type StartBasketRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CustomerId string `protobuf:"bytes,1,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
}

func (x *StartBasketRequest) Reset() {
	*x = StartBasketRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_baskets_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartBasketRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartBasketRequest) ProtoMessage() {}

func (x *StartBasketRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_baskets_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartBasketRequest.ProtoReflect.Descriptor instead.
func (*StartBasketRequest) Descriptor() ([]byte, []int) {
	return file_api_baskets_proto_rawDescGZIP(), []int{2}
}

func (x *StartBasketRequest) GetCustomerId() string {
	if x != nil {
		return x.CustomerId
	}
	return ""
}

type StartBasketResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *StartBasketResponse) Reset() {
	*x = StartBasketResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_baskets_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartBasketResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartBasketResponse) ProtoMessage() {}

func (x *StartBasketResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_baskets_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartBasketResponse.ProtoReflect.Descriptor instead.
func (*StartBasketResponse) Descriptor() ([]byte, []int) {
	return file_api_baskets_proto_rawDescGZIP(), []int{3}
}

func (x *StartBasketResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetBasketRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetBasketRequest) Reset() {
	*x = GetBasketRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_baskets_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetBasketRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBasketRequest) ProtoMessage() {}

func (x *GetBasketRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_baskets_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBasketRequest.ProtoReflect.Descriptor instead.
func (*GetBasketRequest) Descriptor() ([]byte, []int) {
	return file_api_baskets_proto_rawDescGZIP(), []int{4}
}

func (x *GetBasketRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type GetBasketResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Basket *Basket `protobuf:"bytes,1,opt,name=basket,proto3" json:"basket,omitempty"`
}

func (x *GetBasketResponse) Reset() {
	*x = GetBasketResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_baskets_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetBasketResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetBasketResponse) ProtoMessage() {}

func (x *GetBasketResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_baskets_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetBasketResponse.ProtoReflect.Descriptor instead.
func (*GetBasketResponse) Descriptor() ([]byte, []int) {
	return file_api_baskets_proto_rawDescGZIP(), []int{5}
}

func (x *GetBasketResponse) GetBasket() *Basket {
	if x != nil {
		return x.Basket
	}
	return nil
}

type CancelBasketRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *CancelBasketRequest) Reset() {
	*x = CancelBasketRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_baskets_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CancelBasketRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CancelBasketRequest) ProtoMessage() {}

func (x *CancelBasketRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_baskets_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CancelBasketRequest.ProtoReflect.Descriptor instead.
func (*CancelBasketRequest) Descriptor() ([]byte, []int) {
	return file_api_baskets_proto_rawDescGZIP(), []int{6}
}

func (x *CancelBasketRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type CancelBasketResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CancelBasketResponse) Reset() {
	*x = CancelBasketResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_baskets_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CancelBasketResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CancelBasketResponse) ProtoMessage() {}

func (x *CancelBasketResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_baskets_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CancelBasketResponse.ProtoReflect.Descriptor instead.
func (*CancelBasketResponse) Descriptor() ([]byte, []int) {
	return file_api_baskets_proto_rawDescGZIP(), []int{7}
}

type CheckoutBasketRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	PaymentId string `protobuf:"bytes,2,opt,name=payment_id,json=paymentId,proto3" json:"payment_id,omitempty"`
}

func (x *CheckoutBasketRequest) Reset() {
	*x = CheckoutBasketRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_baskets_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckoutBasketRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckoutBasketRequest) ProtoMessage() {}

func (x *CheckoutBasketRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_baskets_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckoutBasketRequest.ProtoReflect.Descriptor instead.
func (*CheckoutBasketRequest) Descriptor() ([]byte, []int) {
	return file_api_baskets_proto_rawDescGZIP(), []int{8}
}

func (x *CheckoutBasketRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CheckoutBasketRequest) GetPaymentId() string {
	if x != nil {
		return x.PaymentId
	}
	return ""
}

type CheckoutBasketResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CheckoutBasketResponse) Reset() {
	*x = CheckoutBasketResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_baskets_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CheckoutBasketResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CheckoutBasketResponse) ProtoMessage() {}

func (x *CheckoutBasketResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_baskets_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CheckoutBasketResponse.ProtoReflect.Descriptor instead.
func (*CheckoutBasketResponse) Descriptor() ([]byte, []int) {
	return file_api_baskets_proto_rawDescGZIP(), []int{9}
}

type AddItemRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ProductId string `protobuf:"bytes,3,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	Quantity  int32  `protobuf:"varint,4,opt,name=quantity,proto3" json:"quantity,omitempty"`
}

func (x *AddItemRequest) Reset() {
	*x = AddItemRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_baskets_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddItemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddItemRequest) ProtoMessage() {}

func (x *AddItemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_baskets_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddItemRequest.ProtoReflect.Descriptor instead.
func (*AddItemRequest) Descriptor() ([]byte, []int) {
	return file_api_baskets_proto_rawDescGZIP(), []int{10}
}

func (x *AddItemRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *AddItemRequest) GetProductId() string {
	if x != nil {
		return x.ProductId
	}
	return ""
}

func (x *AddItemRequest) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

type AddItemResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AddItemResponse) Reset() {
	*x = AddItemResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_baskets_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddItemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddItemResponse) ProtoMessage() {}

func (x *AddItemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_baskets_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddItemResponse.ProtoReflect.Descriptor instead.
func (*AddItemResponse) Descriptor() ([]byte, []int) {
	return file_api_baskets_proto_rawDescGZIP(), []int{11}
}

type RemoveItemRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id        string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	ProductId string `protobuf:"bytes,3,opt,name=product_id,json=productId,proto3" json:"product_id,omitempty"`
	Quantity  int32  `protobuf:"varint,4,opt,name=quantity,proto3" json:"quantity,omitempty"`
}

func (x *RemoveItemRequest) Reset() {
	*x = RemoveItemRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_baskets_proto_msgTypes[12]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveItemRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveItemRequest) ProtoMessage() {}

func (x *RemoveItemRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_baskets_proto_msgTypes[12]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveItemRequest.ProtoReflect.Descriptor instead.
func (*RemoveItemRequest) Descriptor() ([]byte, []int) {
	return file_api_baskets_proto_rawDescGZIP(), []int{12}
}

func (x *RemoveItemRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *RemoveItemRequest) GetProductId() string {
	if x != nil {
		return x.ProductId
	}
	return ""
}

func (x *RemoveItemRequest) GetQuantity() int32 {
	if x != nil {
		return x.Quantity
	}
	return 0
}

type RemoveItemResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *RemoveItemResponse) Reset() {
	*x = RemoveItemResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_baskets_proto_msgTypes[13]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RemoveItemResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RemoveItemResponse) ProtoMessage() {}

func (x *RemoveItemResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_baskets_proto_msgTypes[13]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RemoveItemResponse.ProtoReflect.Descriptor instead.
func (*RemoveItemResponse) Descriptor() ([]byte, []int) {
	return file_api_baskets_proto_rawDescGZIP(), []int{13}
}

var File_api_baskets_proto protoreflect.FileDescriptor

var file_api_baskets_proto_rawDesc = []byte{
	0x0a, 0x11, 0x61, 0x70, 0x69, 0x5f, 0x62, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x09, 0x62, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x73, 0x70, 0x62, 0x22, 0xc3,
	0x01, 0x0a, 0x04, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x74, 0x6f, 0x72, 0x65,
	0x49, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x69, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49,
	0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x4e, 0x61, 0x6d, 0x65,
	0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x70,
	0x72, 0x69, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0c, 0x70, 0x72, 0x6f, 0x64,
	0x75, 0x63, 0x74, 0x50, 0x72, 0x69, 0x63, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75, 0x61, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x18, 0x06, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x71, 0x75, 0x61, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x22, 0x3f, 0x0a, 0x06, 0x42, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x25,
	0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e,
	0x62, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05,
	0x69, 0x74, 0x65, 0x6d, 0x73, 0x22, 0x35, 0x0a, 0x12, 0x53, 0x74, 0x61, 0x72, 0x74, 0x42, 0x61,
	0x73, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1f, 0x0a, 0x0b, 0x63,
	0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0a, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49, 0x64, 0x22, 0x25, 0x0a, 0x13,
	0x53, 0x74, 0x61, 0x72, 0x74, 0x42, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x22, 0x22, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x42, 0x61, 0x73, 0x6b, 0x65, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x3e, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x42, 0x61,
	0x73, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x29, 0x0a, 0x06,
	0x62, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x62,
	0x61, 0x73, 0x6b, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x42, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x52,
	0x06, 0x62, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x22, 0x25, 0x0a, 0x13, 0x43, 0x61, 0x6e, 0x63, 0x65,
	0x6c, 0x42, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x16,
	0x0a, 0x14, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x42, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x46, 0x0a, 0x15, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x6f,
	0x75, 0x74, 0x42, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x1d, 0x0a, 0x0a, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x22, 0x18,
	0x0a, 0x16, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x6f, 0x75, 0x74, 0x42, 0x61, 0x73, 0x6b, 0x65, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x5b, 0x0a, 0x0e, 0x41, 0x64, 0x64, 0x49,
	0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x72,
	0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x71, 0x75, 0x61,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x71, 0x75, 0x61,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x22, 0x11, 0x0a, 0x0f, 0x41, 0x64, 0x64, 0x49, 0x74, 0x65, 0x6d,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x5e, 0x0a, 0x11, 0x52, 0x65, 0x6d, 0x6f,
	0x76, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a,
	0x0a, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x74, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08,
	0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08,
	0x71, 0x75, 0x61, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x22, 0x14, 0x0a, 0x12, 0x52, 0x65, 0x6d, 0x6f,
	0x76, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xe6,
	0x03, 0x0a, 0x0d, 0x42, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x4e, 0x0a, 0x0b, 0x53, 0x74, 0x61, 0x72, 0x74, 0x42, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x12,
	0x1d, 0x2e, 0x62, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x53, 0x74, 0x61, 0x72,
	0x74, 0x42, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1e,
	0x2e, 0x62, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74,
	0x42, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00,
	0x12, 0x48, 0x0a, 0x09, 0x47, 0x65, 0x74, 0x42, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x12, 0x1b, 0x2e,
	0x62, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x61, 0x73,
	0x6b, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x62, 0x61, 0x73,
	0x6b, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x42, 0x61, 0x73, 0x6b, 0x65, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x51, 0x0a, 0x0c, 0x43, 0x61,
	0x6e, 0x63, 0x65, 0x6c, 0x42, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x12, 0x1e, 0x2e, 0x62, 0x61, 0x73,
	0x6b, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x42, 0x61, 0x73,
	0x6b, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x62, 0x61, 0x73,
	0x6b, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x43, 0x61, 0x6e, 0x63, 0x65, 0x6c, 0x42, 0x61, 0x73,
	0x6b, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x57, 0x0a,
	0x0e, 0x43, 0x68, 0x65, 0x63, 0x6b, 0x6f, 0x75, 0x74, 0x42, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x12,
	0x20, 0x2e, 0x62, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x43, 0x68, 0x65, 0x63,
	0x6b, 0x6f, 0x75, 0x74, 0x42, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x21, 0x2e, 0x62, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x43, 0x68,
	0x65, 0x63, 0x6b, 0x6f, 0x75, 0x74, 0x42, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x42, 0x0a, 0x07, 0x41, 0x64, 0x64, 0x49, 0x74, 0x65,
	0x6d, 0x12, 0x19, 0x2e, 0x62, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x41, 0x64,
	0x64, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x62,
	0x61, 0x73, 0x6b, 0x65, 0x74, 0x73, 0x70, 0x62, 0x2e, 0x41, 0x64, 0x64, 0x49, 0x74, 0x65, 0x6d,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x4b, 0x0a, 0x0a, 0x52, 0x65,
	0x6d, 0x6f, 0x76, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x1c, 0x2e, 0x62, 0x61, 0x73, 0x6b, 0x65,
	0x74, 0x73, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x62, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x73,
	0x70, 0x62, 0x2e, 0x52, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x7c, 0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x2e, 0x62,
	0x61, 0x73, 0x6b, 0x65, 0x74, 0x73, 0x70, 0x62, 0x42, 0x0f, 0x41, 0x70, 0x69, 0x42, 0x61, 0x73,
	0x6b, 0x65, 0x74, 0x73, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x16, 0x6d, 0x61, 0x6c,
	0x6c, 0x2f, 0x62, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x73, 0x2f, 0x62, 0x61, 0x73, 0x6b, 0x65, 0x74,
	0x73, 0x70, 0x62, 0xa2, 0x02, 0x03, 0x42, 0x58, 0x58, 0xaa, 0x02, 0x09, 0x42, 0x61, 0x73, 0x6b,
	0x65, 0x74, 0x73, 0x70, 0x62, 0xca, 0x02, 0x09, 0x42, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x73, 0x70,
	0x62, 0xe2, 0x02, 0x15, 0x42, 0x61, 0x73, 0x6b, 0x65, 0x74, 0x73, 0x70, 0x62, 0x5c, 0x47, 0x50,
	0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x09, 0x42, 0x61, 0x73, 0x6b,
	0x65, 0x74, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_baskets_proto_rawDescOnce sync.Once
	file_api_baskets_proto_rawDescData = file_api_baskets_proto_rawDesc
)

func file_api_baskets_proto_rawDescGZIP() []byte {
	file_api_baskets_proto_rawDescOnce.Do(func() {
		file_api_baskets_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_baskets_proto_rawDescData)
	})
	return file_api_baskets_proto_rawDescData
}

var file_api_baskets_proto_msgTypes = make([]protoimpl.MessageInfo, 14)
var file_api_baskets_proto_goTypes = []interface{}{
	(*Item)(nil),                   // 0: basketspb.Item
	(*Basket)(nil),                 // 1: basketspb.Basket
	(*StartBasketRequest)(nil),     // 2: basketspb.StartBasketRequest
	(*StartBasketResponse)(nil),    // 3: basketspb.StartBasketResponse
	(*GetBasketRequest)(nil),       // 4: basketspb.GetBasketRequest
	(*GetBasketResponse)(nil),      // 5: basketspb.GetBasketResponse
	(*CancelBasketRequest)(nil),    // 6: basketspb.CancelBasketRequest
	(*CancelBasketResponse)(nil),   // 7: basketspb.CancelBasketResponse
	(*CheckoutBasketRequest)(nil),  // 8: basketspb.CheckoutBasketRequest
	(*CheckoutBasketResponse)(nil), // 9: basketspb.CheckoutBasketResponse
	(*AddItemRequest)(nil),         // 10: basketspb.AddItemRequest
	(*AddItemResponse)(nil),        // 11: basketspb.AddItemResponse
	(*RemoveItemRequest)(nil),      // 12: basketspb.RemoveItemRequest
	(*RemoveItemResponse)(nil),     // 13: basketspb.RemoveItemResponse
}
var file_api_baskets_proto_depIdxs = []int32{
	0,  // 0: basketspb.Basket.items:type_name -> basketspb.Item
	1,  // 1: basketspb.GetBasketResponse.basket:type_name -> basketspb.Basket
	2,  // 2: basketspb.BasketService.StartBasket:input_type -> basketspb.StartBasketRequest
	4,  // 3: basketspb.BasketService.GetBasket:input_type -> basketspb.GetBasketRequest
	6,  // 4: basketspb.BasketService.CancelBasket:input_type -> basketspb.CancelBasketRequest
	8,  // 5: basketspb.BasketService.CheckoutBasket:input_type -> basketspb.CheckoutBasketRequest
	10, // 6: basketspb.BasketService.AddItem:input_type -> basketspb.AddItemRequest
	12, // 7: basketspb.BasketService.RemoveItem:input_type -> basketspb.RemoveItemRequest
	3,  // 8: basketspb.BasketService.StartBasket:output_type -> basketspb.StartBasketResponse
	5,  // 9: basketspb.BasketService.GetBasket:output_type -> basketspb.GetBasketResponse
	7,  // 10: basketspb.BasketService.CancelBasket:output_type -> basketspb.CancelBasketResponse
	9,  // 11: basketspb.BasketService.CheckoutBasket:output_type -> basketspb.CheckoutBasketResponse
	11, // 12: basketspb.BasketService.AddItem:output_type -> basketspb.AddItemResponse
	13, // 13: basketspb.BasketService.RemoveItem:output_type -> basketspb.RemoveItemResponse
	8,  // [8:14] is the sub-list for method output_type
	2,  // [2:8] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_api_baskets_proto_init() }
func file_api_baskets_proto_init() {
	if File_api_baskets_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_baskets_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Item); i {
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
		file_api_baskets_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Basket); i {
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
		file_api_baskets_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartBasketRequest); i {
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
		file_api_baskets_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StartBasketResponse); i {
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
		file_api_baskets_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetBasketRequest); i {
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
		file_api_baskets_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetBasketResponse); i {
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
		file_api_baskets_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CancelBasketRequest); i {
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
		file_api_baskets_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CancelBasketResponse); i {
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
		file_api_baskets_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckoutBasketRequest); i {
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
		file_api_baskets_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CheckoutBasketResponse); i {
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
		file_api_baskets_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddItemRequest); i {
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
		file_api_baskets_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddItemResponse); i {
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
		file_api_baskets_proto_msgTypes[12].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveItemRequest); i {
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
		file_api_baskets_proto_msgTypes[13].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RemoveItemResponse); i {
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
			RawDescriptor: file_api_baskets_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   14,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_baskets_proto_goTypes,
		DependencyIndexes: file_api_baskets_proto_depIdxs,
		MessageInfos:      file_api_baskets_proto_msgTypes,
	}.Build()
	File_api_baskets_proto = out.File
	file_api_baskets_proto_rawDesc = nil
	file_api_baskets_proto_goTypes = nil
	file_api_baskets_proto_depIdxs = nil
}
