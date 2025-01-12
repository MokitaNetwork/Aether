// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: aeth/savings/v1beta1/store.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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

// Params defines the parameters for the savings module.
type Params struct {
	SupportedDenoms []string `protobuf:"bytes,1,rep,name=supported_denoms,json=supportedDenoms,proto3" json:"supported_denoms,omitempty"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_f7110366fa182786, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

// Deposit defines an amount of coins deposited into a savings module account.
type Deposit struct {
	Depositor github_com_cosmos_cosmos_sdk_types.AccAddress `protobuf:"bytes,1,opt,name=depositor,proto3,casttype=github.com/cosmos/cosmos-sdk/types.AccAddress" json:"depositor,omitempty"`
	Amount    github_com_cosmos_cosmos_sdk_types.Coins      `protobuf:"bytes,2,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
}

func (m *Deposit) Reset()         { *m = Deposit{} }
func (m *Deposit) String() string { return proto.CompactTextString(m) }
func (*Deposit) ProtoMessage()    {}
func (*Deposit) Descriptor() ([]byte, []int) {
	return fileDescriptor_f7110366fa182786, []int{1}
}
func (m *Deposit) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Deposit) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Deposit.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Deposit) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Deposit.Merge(m, src)
}
func (m *Deposit) XXX_Size() int {
	return m.Size()
}
func (m *Deposit) XXX_DiscardUnknown() {
	xxx_messageInfo_Deposit.DiscardUnknown(m)
}

var xxx_messageInfo_Deposit proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Params)(nil), "aeth.savings.v1beta1.Params")
	proto.RegisterType((*Deposit)(nil), "aeth.savings.v1beta1.Deposit")
}

func init() { proto.RegisterFile("aeth/savings/v1beta1/store.proto", fileDescriptor_f7110366fa182786) }

var fileDescriptor_f7110366fa182786 = []byte{
	// 335 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x51, 0xc1, 0x4e, 0xc2, 0x40,
	0x14, 0xec, 0x4a, 0x82, 0xa1, 0x1e, 0x34, 0x95, 0x03, 0x70, 0x58, 0x1a, 0x4e, 0xe5, 0xd0, 0x5d,
	0x91, 0x2f, 0xa0, 0x92, 0xe8, 0xd1, 0x70, 0xf4, 0x42, 0xb6, 0xed, 0x5a, 0x1b, 0x6c, 0x5f, 0xd3,
	0xb7, 0x10, 0xf9, 0x0b, 0xbf, 0xc3, 0xb3, 0x1f, 0xc1, 0x91, 0x78, 0x30, 0x9e, 0x50, 0xe1, 0x2f,
	0x3c, 0x99, 0xb6, 0x2b, 0x7a, 0xf4, 0xb4, 0x6f, 0xe7, 0xcd, 0x4c, 0x26, 0xf3, 0x4c, 0x7b, 0x26,
	0x16, 0x82, 0xa3, 0x58, 0xc4, 0x69, 0x84, 0x7c, 0x31, 0xf0, 0xa5, 0x12, 0x03, 0x8e, 0x0a, 0x72,
	0xc9, 0xb2, 0x1c, 0x14, 0x58, 0xcd, 0x82, 0xc1, 0x34, 0x83, 0x69, 0x46, 0x87, 0x06, 0x80, 0x09,
	0x20, 0xf7, 0x05, 0xca, 0xbd, 0x2c, 0x80, 0x38, 0xad, 0x54, 0x9d, 0x76, 0xb5, 0x9f, 0x96, 0x3f,
	0x5e, 0x7d, 0xf4, 0xaa, 0x19, 0x41, 0x04, 0x15, 0x5e, 0x4c, 0x15, 0xda, 0x1b, 0x9a, 0xf5, 0x6b,
	0x91, 0x8b, 0x04, 0xad, 0xbe, 0x79, 0x82, 0xf3, 0x2c, 0x83, 0x5c, 0xc9, 0x70, 0x1a, 0xca, 0x14,
	0x12, 0x6c, 0x11, 0xbb, 0xe6, 0x34, 0x26, 0xc7, 0x7b, 0x7c, 0x5c, 0xc2, 0xbd, 0x57, 0x62, 0x1e,
	0x8e, 0x65, 0x06, 0x18, 0x2b, 0xeb, 0xd6, 0x6c, 0x84, 0xd5, 0x08, 0x79, 0x8b, 0xd8, 0xc4, 0x69,
	0x78, 0x57, 0x5f, 0x9b, 0xae, 0x1b, 0xc5, 0xea, 0x6e, 0xee, 0xb3, 0x00, 0x12, 0x1d, 0x43, 0x3f,
	0x2e, 0x86, 0x33, 0xae, 0x96, 0x99, 0x44, 0x36, 0x0a, 0x82, 0x51, 0x18, 0xe6, 0x12, 0xf1, 0xe5,
	0xd9, 0x3d, 0xd5, 0x61, 0x35, 0xe2, 0x2d, 0x95, 0xc4, 0xc9, 0xaf, 0xb5, 0x15, 0x98, 0x75, 0x91,
	0xc0, 0x3c, 0x55, 0xad, 0x03, 0xbb, 0xe6, 0x1c, 0x9d, 0xb7, 0x99, 0x16, 0x14, 0x55, 0xfc, 0xf4,
	0xc3, 0x2e, 0x20, 0x4e, 0xbd, 0xb3, 0xd5, 0xa6, 0x6b, 0x3c, 0xbd, 0x77, 0x9d, 0x7f, 0x64, 0x28,
	0x04, 0x38, 0xd1, 0xd6, 0xde, 0xe5, 0xea, 0x93, 0x1a, 0xab, 0x2d, 0x25, 0xeb, 0x2d, 0x25, 0x1f,
	0x5b, 0x4a, 0x1e, 0x77, 0xd4, 0x58, 0xef, 0xa8, 0xf1, 0xb6, 0xa3, 0xc6, 0x4d, 0xff, 0x8f, 0x5f,
	0x71, 0x1d, 0xf7, 0x5e, 0xf8, 0x58, 0x4e, 0xfc, 0x61, 0x7f, 0xcb, 0xd2, 0xd6, 0xaf, 0x97, 0xed,
	0x0e, 0xbf, 0x03, 0x00, 0x00, 0xff, 0xff, 0xef, 0x42, 0x46, 0x9e, 0xe8, 0x01, 0x00, 0x00,
}

func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.SupportedDenoms) > 0 {
		for iNdEx := len(m.SupportedDenoms) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.SupportedDenoms[iNdEx])
			copy(dAtA[i:], m.SupportedDenoms[iNdEx])
			i = encodeVarintStore(dAtA, i, uint64(len(m.SupportedDenoms[iNdEx])))
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func (m *Deposit) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Deposit) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Deposit) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Amount) > 0 {
		for iNdEx := len(m.Amount) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Amount[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintStore(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	if len(m.Depositor) > 0 {
		i -= len(m.Depositor)
		copy(dAtA[i:], m.Depositor)
		i = encodeVarintStore(dAtA, i, uint64(len(m.Depositor)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintStore(dAtA []byte, offset int, v uint64) int {
	offset -= sovStore(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.SupportedDenoms) > 0 {
		for _, s := range m.SupportedDenoms {
			l = len(s)
			n += 1 + l + sovStore(uint64(l))
		}
	}
	return n
}

func (m *Deposit) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Depositor)
	if l > 0 {
		n += 1 + l + sovStore(uint64(l))
	}
	if len(m.Amount) > 0 {
		for _, e := range m.Amount {
			l = e.Size()
			n += 1 + l + sovStore(uint64(l))
		}
	}
	return n
}

func sovStore(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozStore(x uint64) (n int) {
	return sovStore(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStore
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
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SupportedDenoms", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStore
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
				return ErrInvalidLengthStore
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStore
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.SupportedDenoms = append(m.SupportedDenoms, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStore(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthStore
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
func (m *Deposit) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowStore
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
			return fmt.Errorf("proto: Deposit: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Deposit: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Depositor", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStore
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
				return ErrInvalidLengthStore
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthStore
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Depositor = github_com_cosmos_cosmos_sdk_types.AccAddress(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowStore
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
				return ErrInvalidLengthStore
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthStore
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Amount = append(m.Amount, types.Coin{})
			if err := m.Amount[len(m.Amount)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipStore(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthStore
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
func skipStore(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowStore
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
					return 0, ErrIntOverflowStore
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
					return 0, ErrIntOverflowStore
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
				return 0, ErrInvalidLengthStore
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupStore
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthStore
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthStore        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowStore          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupStore = fmt.Errorf("proto: unexpected end of group")
)
