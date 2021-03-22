// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: coinswap/query.proto

package coinswap

import (
	context "context"
	fmt "fmt"
	_ "github.com/gogo/protobuf/gogoproto"
	grpc1 "github.com/gogo/protobuf/grpc"
	proto "github.com/gogo/protobuf/proto"
	types "github.com/irisnet/irishub-sdk-go/types"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// QueryLiquidityRequest is request type for the Query/Liquidity RPC method
type QueryLiquidityRequest struct {
	Denom string `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty"`
}

func (m *QueryLiquidityRequest) Reset()         { *m = QueryLiquidityRequest{} }
func (m *QueryLiquidityRequest) String() string { return proto.CompactTextString(m) }
func (*QueryLiquidityRequest) ProtoMessage()    {}
func (*QueryLiquidityRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cabf8423404f12f, []int{0}
}
func (m *QueryLiquidityRequest) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryLiquidityRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryLiquidityRequest.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryLiquidityRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryLiquidityRequest.Merge(m, src)
}
func (m *QueryLiquidityRequest) XXX_Size() int {
	return m.Size()
}
func (m *QueryLiquidityRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryLiquidityRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryLiquidityRequest proto.InternalMessageInfo

func (m *QueryLiquidityRequest) GetDenom() string {
	if m != nil {
		return m.Denom
	}
	return ""
}

// QueryLiquidityResponse is response type for the Query/Liquidity RPC method
type QueryLiquidityResponse struct {
	Standard  types.Coin `protobuf:"bytes,1,opt,name=standard,proto3" json:"standard"`
	Token     types.Coin `protobuf:"bytes,2,opt,name=token,proto3" json:"token"`
	Liquidity types.Coin `protobuf:"bytes,3,opt,name=liquidity,proto3" json:"liquidity"`
	Fee       string     `protobuf:"bytes,4,opt,name=fee,proto3" json:"fee,omitempty"`
}

func (m *QueryLiquidityResponse) Reset()         { *m = QueryLiquidityResponse{} }
func (m *QueryLiquidityResponse) String() string { return proto.CompactTextString(m) }
func (*QueryLiquidityResponse) ProtoMessage()    {}
func (*QueryLiquidityResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_2cabf8423404f12f, []int{1}
}
func (m *QueryLiquidityResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *QueryLiquidityResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_QueryLiquidityResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *QueryLiquidityResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryLiquidityResponse.Merge(m, src)
}
func (m *QueryLiquidityResponse) XXX_Size() int {
	return m.Size()
}
func (m *QueryLiquidityResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryLiquidityResponse.DiscardUnknown(m)
}

var xxx_messageInfo_QueryLiquidityResponse proto.InternalMessageInfo

func (m *QueryLiquidityResponse) GetStandard() types.Coin {
	if m != nil {
		return m.Standard
	}
	return types.Coin{}
}

func (m *QueryLiquidityResponse) GetToken() types.Coin {
	if m != nil {
		return m.Token
	}
	return types.Coin{}
}

func (m *QueryLiquidityResponse) GetLiquidity() types.Coin {
	if m != nil {
		return m.Liquidity
	}
	return types.Coin{}
}

func (m *QueryLiquidityResponse) GetFee() string {
	if m != nil {
		return m.Fee
	}
	return ""
}

func init() {
	proto.RegisterType((*QueryLiquidityRequest)(nil), "irismod.coinswap.QueryLiquidityRequest")
	proto.RegisterType((*QueryLiquidityResponse)(nil), "irismod.coinswap.QueryLiquidityResponse")
}

func init() { proto.RegisterFile("coinswap/query.proto", fileDescriptor_2cabf8423404f12f) }

var fileDescriptor_2cabf8423404f12f = []byte{
	// 380 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x92, 0xbf, 0x6b, 0xdb, 0x40,
	0x14, 0xc7, 0x75, 0xfe, 0x51, 0xea, 0xeb, 0x62, 0x0e, 0xb7, 0xa8, 0xa6, 0xa8, 0xc5, 0x50, 0xec,
	0x45, 0x77, 0xd8, 0xa5, 0x53, 0xe9, 0xe2, 0xae, 0x5e, 0xea, 0xb1, 0xdb, 0xc9, 0xba, 0xca, 0x87,
	0xad, 0x7b, 0xb2, 0xee, 0x94, 0x60, 0x42, 0x96, 0xec, 0x81, 0x40, 0x86, 0xfc, 0x4b, 0x1e, 0x0d,
	0x59, 0xb2, 0x24, 0x04, 0x3b, 0x7f, 0x48, 0xd0, 0xc9, 0x56, 0xc0, 0x04, 0xe2, 0x49, 0x4f, 0xf7,
	0xde, 0xe7, 0xde, 0xf7, 0x7d, 0xdf, 0xe1, 0xd6, 0x04, 0xa4, 0xd2, 0xa7, 0x3c, 0x61, 0x8b, 0x4c,
	0xa4, 0x4b, 0x9a, 0xa4, 0x60, 0x80, 0x34, 0x65, 0x2a, 0x75, 0x0c, 0x21, 0xdd, 0x67, 0xdb, 0xde,
	0x04, 0x74, 0x0c, 0x9a, 0x05, 0x5c, 0x0b, 0x76, 0xd2, 0x0f, 0x84, 0xe1, 0x7d, 0x96, 0x67, 0x0b,
	0xa2, 0xdd, 0x8a, 0x20, 0x02, 0x1b, 0xb2, 0x3c, 0xda, 0x9d, 0x7e, 0x89, 0x00, 0xa2, 0xb9, 0x60,
	0x3c, 0x91, 0x8c, 0x2b, 0x05, 0x86, 0x1b, 0x09, 0x4a, 0x17, 0xd9, 0x8e, 0x8f, 0x3f, 0xfe, 0xcd,
	0x9b, 0x8e, 0xe4, 0x22, 0x93, 0xa1, 0x34, 0xcb, 0xb1, 0x58, 0x64, 0x42, 0x1b, 0xd2, 0xc2, 0xf5,
	0x50, 0x28, 0x88, 0x5d, 0xf4, 0x0d, 0xf5, 0x1a, 0xe3, 0xe2, 0xa7, 0x73, 0x8f, 0xf0, 0xa7, 0xc3,
	0x7a, 0x9d, 0x80, 0xd2, 0x82, 0xfc, 0xc2, 0xef, 0xb5, 0xe1, 0x2a, 0xe4, 0x69, 0x68, 0x99, 0x0f,
	0x83, 0xcf, 0xb4, 0x10, 0x4c, 0x73, 0xc1, 0x74, 0x27, 0x98, 0xfe, 0x01, 0xa9, 0x86, 0xb5, 0xd5,
	0xc3, 0x57, 0x67, 0x5c, 0x02, 0xe4, 0x27, 0xae, 0x1b, 0x98, 0x09, 0xe5, 0x56, 0x8e, 0x23, 0x8b,
	0x6a, 0xf2, 0x1b, 0x37, 0xe6, 0x7b, 0x21, 0x6e, 0xf5, 0x38, 0xf4, 0x85, 0x20, 0x4d, 0x5c, 0xfd,
	0x2f, 0x84, 0x5b, 0xb3, 0x13, 0xe6, 0xe1, 0xe0, 0x06, 0xe1, 0xba, 0x9d, 0x8f, 0x5c, 0x22, 0xdc,
	0x28, 0x87, 0x24, 0x5d, 0x7a, 0xb8, 0x0d, 0xfa, 0xaa, 0x6d, 0xed, 0xde, 0xdb, 0x85, 0x85, 0x5f,
	0x1d, 0xff, 0xe2, 0xf6, 0xe9, 0xba, 0xd2, 0x25, 0xdf, 0xd9, 0x8e, 0x60, 0xe5, 0x33, 0xd8, 0x2b,
	0x94, 0x42, 0xb3, 0x33, 0x6b, 0xfc, 0xf9, 0x70, 0xb4, 0xda, 0x78, 0x68, 0xbd, 0xf1, 0xd0, 0xe3,
	0xc6, 0x43, 0x57, 0x5b, 0xcf, 0x59, 0x6f, 0x3d, 0xe7, 0x6e, 0xeb, 0x39, 0xff, 0x06, 0x91, 0x34,
	0xd3, 0x2c, 0xa0, 0x13, 0x88, 0xed, 0x55, 0x4a, 0x18, 0xfb, 0x9d, 0x66, 0x81, 0xaf, 0xc3, 0x99,
	0x1f, 0x01, 0x8b, 0x21, 0xcc, 0xe6, 0x42, 0x97, 0x1d, 0x82, 0x77, 0x76, 0xfb, 0x3f, 0x9e, 0x03,
	0x00, 0x00, 0xff, 0xff, 0x4f, 0x9c, 0xdd, 0x1a, 0x7b, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// QueryClient is the client API for Query service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type QueryClient interface {
	// Liquidity returns the total liquidity available for the provided
	// denomination
	Liquidity(ctx context.Context, in *QueryLiquidityRequest, opts ...grpc.CallOption) (*QueryLiquidityResponse, error)
}

type queryClient struct {
	cc grpc1.ClientConn
}

func NewQueryClient(cc grpc1.ClientConn) QueryClient {
	return &queryClient{cc}
}

func (c *queryClient) Liquidity(ctx context.Context, in *QueryLiquidityRequest, opts ...grpc.CallOption) (*QueryLiquidityResponse, error) {
	out := new(QueryLiquidityResponse)
	err := c.cc.Invoke(ctx, "/irismod.coinswap.Query/Liquidity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QueryServer is the server API for Query service.
type QueryServer interface {
	// Liquidity returns the total liquidity available for the provided
	// denomination
	Liquidity(context.Context, *QueryLiquidityRequest) (*QueryLiquidityResponse, error)
}

// UnimplementedQueryServer can be embedded to have forward compatible implementations.
type UnimplementedQueryServer struct {
}

func (*UnimplementedQueryServer) Liquidity(ctx context.Context, req *QueryLiquidityRequest) (*QueryLiquidityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Liquidity not implemented")
}

func RegisterQueryServer(s grpc1.Server, srv QueryServer) {
	s.RegisterService(&_Query_serviceDesc, srv)
}

func _Query_Liquidity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryLiquidityRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QueryServer).Liquidity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/irismod.coinswap.Query/Liquidity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QueryServer).Liquidity(ctx, req.(*QueryLiquidityRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Query_serviceDesc = grpc.ServiceDesc{
	ServiceName: "irismod.coinswap.Query",
	HandlerType: (*QueryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Liquidity",
			Handler:    _Query_Liquidity_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "coinswap/query.proto",
}

func (m *QueryLiquidityRequest) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryLiquidityRequest) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryLiquidityRequest) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *QueryLiquidityResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *QueryLiquidityResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *QueryLiquidityResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Fee) > 0 {
		i -= len(m.Fee)
		copy(dAtA[i:], m.Fee)
		i = encodeVarintQuery(dAtA, i, uint64(len(m.Fee)))
		i--
		dAtA[i] = 0x22
	}
	{
		size, err := m.Liquidity.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size, err := m.Token.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size, err := m.Standard.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintQuery(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintQuery(dAtA []byte, offset int, v uint64) int {
	offset -= sovQuery(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *QueryLiquidityRequest) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func (m *QueryLiquidityResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Standard.Size()
	n += 1 + l + sovQuery(uint64(l))
	l = m.Token.Size()
	n += 1 + l + sovQuery(uint64(l))
	l = m.Liquidity.Size()
	n += 1 + l + sovQuery(uint64(l))
	l = len(m.Fee)
	if l > 0 {
		n += 1 + l + sovQuery(uint64(l))
	}
	return n
}

func sovQuery(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozQuery(x uint64) (n int) {
	return sovQuery(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *QueryLiquidityRequest) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryLiquidityRequest: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryLiquidityRequest: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *QueryLiquidityResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: QueryLiquidityResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: QueryLiquidityResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Standard", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Standard.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Token", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Token.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Liquidity", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Liquidity.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Fee", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthQuery
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthQuery
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Fee = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipQuery(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthQuery
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipQuery(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowQuery
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowQuery
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthQuery
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupQuery
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthQuery
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthQuery        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowQuery          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupQuery = fmt.Errorf("proto: unexpected end of group")
)
