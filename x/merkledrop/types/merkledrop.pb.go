// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: bitsong/merkledrop/v1beta1/merkledrop.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
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

type Merkledrop struct {
	Id          uint64                                  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	MerkleRoot  []byte                                  `protobuf:"bytes,2,opt,name=merkle_root,json=merkleRoot,proto3" json:"merkle_root,omitempty" yaml:"merkle_root"`
	TotalAmount github_com_cosmos_cosmos_sdk_types.Coin `protobuf:"bytes,3,opt,name=total_amount,json=totalAmount,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Coin" json:"total_amount" yaml:"total_amount"`
	Owner       string                                  `protobuf:"bytes,4,opt,name=owner,proto3" json:"owner,omitempty"`
}

func (m *Merkledrop) Reset()      { *m = Merkledrop{} }
func (*Merkledrop) ProtoMessage() {}
func (*Merkledrop) Descriptor() ([]byte, []int) {
	return fileDescriptor_21aba39fc2313837, []int{0}
}
func (m *Merkledrop) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Merkledrop) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Merkledrop.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Merkledrop) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Merkledrop.Merge(m, src)
}
func (m *Merkledrop) XXX_Size() int {
	return m.Size()
}
func (m *Merkledrop) XXX_DiscardUnknown() {
	xxx_messageInfo_Merkledrop.DiscardUnknown(m)
}

var xxx_messageInfo_Merkledrop proto.InternalMessageInfo

type EventMerkledropCreate struct {
	Owner        string `protobuf:"bytes,1,opt,name=owner,proto3" json:"owner,omitempty"`
	MerkledropId uint64 `protobuf:"varint,2,opt,name=merkledrop_id,json=merkledropId,proto3" json:"merkledrop_id,omitempty"`
}

func (m *EventMerkledropCreate) Reset()         { *m = EventMerkledropCreate{} }
func (m *EventMerkledropCreate) String() string { return proto.CompactTextString(m) }
func (*EventMerkledropCreate) ProtoMessage()    {}
func (*EventMerkledropCreate) Descriptor() ([]byte, []int) {
	return fileDescriptor_21aba39fc2313837, []int{1}
}
func (m *EventMerkledropCreate) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *EventMerkledropCreate) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_EventMerkledropCreate.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *EventMerkledropCreate) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventMerkledropCreate.Merge(m, src)
}
func (m *EventMerkledropCreate) XXX_Size() int {
	return m.Size()
}
func (m *EventMerkledropCreate) XXX_DiscardUnknown() {
	xxx_messageInfo_EventMerkledropCreate.DiscardUnknown(m)
}

var xxx_messageInfo_EventMerkledropCreate proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Merkledrop)(nil), "bitsong.merkledrop.v1beta1.Merkledrop")
	proto.RegisterType((*EventMerkledropCreate)(nil), "bitsong.merkledrop.v1beta1.EventMerkledropCreate")
}

func init() {
	proto.RegisterFile("bitsong/merkledrop/v1beta1/merkledrop.proto", fileDescriptor_21aba39fc2313837)
}

var fileDescriptor_21aba39fc2313837 = []byte{
	// 388 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x52, 0xbf, 0x6f, 0xda, 0x40,
	0x14, 0xf6, 0x51, 0x5a, 0xb5, 0x07, 0xed, 0xe0, 0xd2, 0xca, 0x65, 0x38, 0x5b, 0xee, 0x50, 0x4b,
	0x15, 0xb6, 0x68, 0x87, 0x44, 0x6c, 0x01, 0x25, 0x52, 0x86, 0x2c, 0x1e, 0x32, 0x64, 0x41, 0xfe,
	0x71, 0x38, 0x27, 0x6c, 0x3f, 0x64, 0x1f, 0x24, 0xec, 0x19, 0x32, 0x66, 0xcc, 0xc8, 0x9f, 0xc3,
	0xc8, 0x18, 0x65, 0xb0, 0x12, 0xf8, 0x0f, 0x98, 0x33, 0x44, 0xf8, 0x0c, 0xf6, 0x74, 0xf7, 0xbd,
	0xf7, 0xbe, 0xf7, 0xbe, 0xf7, 0x03, 0xff, 0x75, 0x19, 0x4f, 0x21, 0x0e, 0xac, 0x88, 0x26, 0xe3,
	0x90, 0xfa, 0x09, 0x4c, 0xac, 0x59, 0xd7, 0xa5, 0xdc, 0xe9, 0x56, 0x4c, 0xe6, 0x24, 0x01, 0x0e,
	0x72, 0xbb, 0x08, 0x36, 0x2b, 0x9e, 0x22, 0xb8, 0xdd, 0x0a, 0x20, 0x80, 0x3c, 0xcc, 0xda, 0xfd,
	0x04, 0xa3, 0x4d, 0x3c, 0x48, 0x23, 0x48, 0x2d, 0xd7, 0x49, 0xe9, 0x21, 0xaf, 0x07, 0x2c, 0x16,
	0x7e, 0xfd, 0x0d, 0x61, 0x7c, 0x71, 0x48, 0x26, 0x7f, 0xc3, 0x35, 0xe6, 0x2b, 0x48, 0x43, 0x46,
	0xdd, 0xae, 0x31, 0x5f, 0x3e, 0xc2, 0x0d, 0x51, 0x6a, 0x98, 0x00, 0x70, 0xa5, 0xa6, 0x21, 0xa3,
	0xd9, 0xff, 0xb9, 0xcd, 0x54, 0x79, 0xee, 0x44, 0x61, 0x4f, 0xaf, 0x38, 0x75, 0x1b, 0x0b, 0x64,
	0x03, 0x70, 0xf9, 0x0e, 0xe1, 0x26, 0x07, 0xee, 0x84, 0x43, 0x27, 0x82, 0x69, 0xcc, 0x95, 0x0f,
	0x1a, 0x32, 0x1a, 0xff, 0x7e, 0x99, 0x42, 0x8f, 0xb9, 0xd3, 0xb3, 0x97, 0x6e, 0x0e, 0x80, 0xc5,
	0xfd, 0xb3, 0x65, 0xa6, 0x4a, 0xcf, 0x99, 0xfa, 0x27, 0x60, 0xfc, 0x7a, 0xea, 0x9a, 0x1e, 0x44,
	0x56, 0x21, 0x5e, 0x3c, 0x9d, 0xd4, 0x1f, 0x5b, 0x7c, 0x3e, 0xa1, 0x69, 0x4e, 0xd8, 0x66, 0xea,
	0x77, 0x21, 0xa2, 0x5a, 0x47, 0xb7, 0x1b, 0x39, 0x3c, 0xc9, 0x91, 0xdc, 0xc2, 0x1f, 0xe1, 0x26,
	0xa6, 0x89, 0x52, 0xd7, 0x90, 0xf1, 0xc5, 0x16, 0xa0, 0xf7, 0xf9, 0x7e, 0xa1, 0x4a, 0x8f, 0x0b,
	0x55, 0xd2, 0x6d, 0xfc, 0xe3, 0x74, 0x46, 0x63, 0x5e, 0x8e, 0x60, 0x90, 0x50, 0x87, 0xd3, 0x92,
	0x88, 0x2a, 0x44, 0xf9, 0x37, 0xfe, 0x5a, 0x4e, 0x7e, 0xc8, 0xfc, 0x7c, 0x20, 0x75, 0xbb, 0x59,
	0x1a, 0xcf, 0xfd, 0xfe, 0xe5, 0xf2, 0x95, 0x48, 0xcb, 0x35, 0x41, 0xab, 0x35, 0x41, 0x2f, 0x6b,
	0x82, 0x1e, 0x36, 0x44, 0x5a, 0x6d, 0x88, 0xf4, 0xb4, 0x21, 0xd2, 0xd5, 0x71, 0xa5, 0xbd, 0x62,
	0x9b, 0x30, 0x1a, 0x31, 0x8f, 0x39, 0xa1, 0x15, 0x40, 0x67, 0x7f, 0x0d, 0xb7, 0xd5, 0x7b, 0xc8,
	0x9b, 0x76, 0x3f, 0xe5, 0x1b, 0xfb, 0xff, 0x1e, 0x00, 0x00, 0xff, 0xff, 0xed, 0xb8, 0x11, 0x6a,
	0x32, 0x02, 0x00, 0x00,
}

func (m *Merkledrop) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Merkledrop) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Merkledrop) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintMerkledrop(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0x22
	}
	{
		size := m.TotalAmount.Size()
		i -= size
		if _, err := m.TotalAmount.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintMerkledrop(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	if len(m.MerkleRoot) > 0 {
		i -= len(m.MerkleRoot)
		copy(dAtA[i:], m.MerkleRoot)
		i = encodeVarintMerkledrop(dAtA, i, uint64(len(m.MerkleRoot)))
		i--
		dAtA[i] = 0x12
	}
	if m.Id != 0 {
		i = encodeVarintMerkledrop(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *EventMerkledropCreate) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *EventMerkledropCreate) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *EventMerkledropCreate) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.MerkledropId != 0 {
		i = encodeVarintMerkledrop(dAtA, i, uint64(m.MerkledropId))
		i--
		dAtA[i] = 0x10
	}
	if len(m.Owner) > 0 {
		i -= len(m.Owner)
		copy(dAtA[i:], m.Owner)
		i = encodeVarintMerkledrop(dAtA, i, uint64(len(m.Owner)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintMerkledrop(dAtA []byte, offset int, v uint64) int {
	offset -= sovMerkledrop(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Merkledrop) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Id != 0 {
		n += 1 + sovMerkledrop(uint64(m.Id))
	}
	l = len(m.MerkleRoot)
	if l > 0 {
		n += 1 + l + sovMerkledrop(uint64(l))
	}
	l = m.TotalAmount.Size()
	n += 1 + l + sovMerkledrop(uint64(l))
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovMerkledrop(uint64(l))
	}
	return n
}

func (m *EventMerkledropCreate) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Owner)
	if l > 0 {
		n += 1 + l + sovMerkledrop(uint64(l))
	}
	if m.MerkledropId != 0 {
		n += 1 + sovMerkledrop(uint64(m.MerkledropId))
	}
	return n
}

func sovMerkledrop(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozMerkledrop(x uint64) (n int) {
	return sovMerkledrop(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Merkledrop) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMerkledrop
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
			return fmt.Errorf("proto: Merkledrop: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Merkledrop: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMerkledrop
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MerkleRoot", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMerkledrop
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthMerkledrop
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthMerkledrop
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MerkleRoot = append(m.MerkleRoot[:0], dAtA[iNdEx:postIndex]...)
			if m.MerkleRoot == nil {
				m.MerkleRoot = []byte{}
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field TotalAmount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMerkledrop
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
				return ErrInvalidLengthMerkledrop
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthMerkledrop
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.TotalAmount.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMerkledrop
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
				return ErrInvalidLengthMerkledrop
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMerkledrop
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipMerkledrop(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMerkledrop
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
func (m *EventMerkledropCreate) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowMerkledrop
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
			return fmt.Errorf("proto: EventMerkledropCreate: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: EventMerkledropCreate: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Owner", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMerkledrop
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
				return ErrInvalidLengthMerkledrop
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthMerkledrop
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Owner = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MerkledropId", wireType)
			}
			m.MerkledropId = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowMerkledrop
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MerkledropId |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipMerkledrop(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthMerkledrop
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
func skipMerkledrop(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowMerkledrop
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
					return 0, ErrIntOverflowMerkledrop
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
					return 0, ErrIntOverflowMerkledrop
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
				return 0, ErrInvalidLengthMerkledrop
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupMerkledrop
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthMerkledrop
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthMerkledrop        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowMerkledrop          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupMerkledrop = fmt.Errorf("proto: unexpected end of group")
)