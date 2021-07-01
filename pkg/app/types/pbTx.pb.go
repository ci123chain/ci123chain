// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: pbTx.proto

package types

import (
	fmt "fmt"
	codectypes "github.com/ci123chain/ci123chain/pkg/abci/codec/types"
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

type PbTx struct {
	From      []byte          `protobuf:"bytes,1,opt,name=from,proto3" json:"from,omitempty"`
	Nonce     uint64          `protobuf:"varint,2,opt,name=nonce,proto3" json:"nonce,omitempty"`
	Gas       uint64          `protobuf:"varint,3,opt,name=gas,proto3" json:"gas,omitempty"`
	Msgs      []*codectypes.Any    `protobuf:"bytes,4,rep,name=msgs,proto3" json:"msgs,omitempty"`
	PubKey    []byte          `protobuf:"bytes,5,opt,name=pub_key,json=pubKey,proto3" json:"pub_key,omitempty"`
	Signature []byte          `protobuf:"bytes,6,opt,name=signature,proto3" json:"signature,omitempty"`
}

func (m *PbTx) Reset()         { *m = PbTx{} }
func (m *PbTx) String() string { return proto.CompactTextString(m) }
func (*PbTx) ProtoMessage()    {}
func (*PbTx) Descriptor() ([]byte, []int) {
	return fileDescriptor_0fd2153dc07d3b5c, []int{0}
}
func (m *PbTx) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *PbTx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_PbTx.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *PbTx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PbTx.Merge(m, src)
}
func (m *PbTx) XXX_Size() int {
	return m.Size()
}
func (m *PbTx) XXX_DiscardUnknown() {
	xxx_messageInfo_PbTx.DiscardUnknown(m)
}

var xxx_messageInfo_PbTx proto.InternalMessageInfo

func (m *PbTx) GetFrom() []byte {
	if m != nil {
		return m.From
	}
	return nil
}

func (m *PbTx) GetNonce() uint64 {
	if m != nil {
		return m.Nonce
	}
	return 0
}

func (m *PbTx) GetGas() uint64 {
	if m != nil {
		return m.Gas
	}
	return 0
}

//func (m *PbTx) GetMsgs() []*codectypes.Any {
//	if m != nil {
//		return m.Msgs
//	}
//	return nil
//}

func (m *PbTx) GetPubKey() []byte {
	if m != nil {
		return m.PubKey
	}
	return nil
}

func (m *PbTx) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}

func init() {
	proto.RegisterType((*PbTx)(nil), "weelink.tx.PbTx")
}

func init() { proto.RegisterFile("tx.proto", fileDescriptor_0fd2153dc07d3b5c) }

var fileDescriptor_0fd2153dc07d3b5c = []byte{
	// 259 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x4c, 0x8e, 0xb1, 0x4a, 0xc4, 0x30,
	0x1c, 0xc6, 0x1b, 0x9b, 0xab, 0x1a, 0x1d, 0x24, 0x14, 0x8c, 0x87, 0x84, 0xe2, 0xd4, 0x29, 0xc5,
	0xbb, 0xcd, 0x4d, 0x57, 0x17, 0x29, 0x4e, 0x2e, 0xd2, 0x94, 0x5c, 0x2e, 0xdc, 0x35, 0x29, 0x6d,
	0x82, 0xed, 0x5b, 0xf8, 0x0e, 0xbe, 0x8c, 0xe3, 0x8d, 0x8e, 0xd2, 0xbe, 0x88, 0x34, 0x45, 0x74,
	0xfb, 0xfd, 0xbf, 0xef, 0x83, 0xdf, 0x1f, 0x9d, 0xd8, 0x8e, 0xd5, 0x8d, 0xb1, 0x06, 0xa3, 0x37,
	0x21, 0xf6, 0x4a, 0xef, 0x98, 0xed, 0x96, 0xb1, 0x34, 0xd2, 0xf8, 0x38, 0x9b, 0x68, 0x5e, 0x2c,
	0xaf, 0xa4, 0x31, 0x72, 0x2f, 0x32, 0x7f, 0x71, 0xb7, 0xc9, 0x0a, 0xdd, 0xcf, 0xd5, 0xcd, 0x07,
	0x40, 0xf0, 0x89, 0x3f, 0x77, 0x18, 0x23, 0xb8, 0x69, 0x4c, 0x45, 0x40, 0x02, 0xd2, 0xf3, 0xdc,
	0x33, 0x8e, 0xd1, 0x42, 0x1b, 0x5d, 0x0a, 0x72, 0x94, 0x80, 0x14, 0xe6, 0xf3, 0x81, 0x2f, 0x50,
	0x28, 0x8b, 0x96, 0x84, 0x3e, 0x9b, 0x10, 0xa7, 0x08, 0x56, 0xad, 0x6c, 0x09, 0x4c, 0xc2, 0xf4,
	0x6c, 0x15, 0xb3, 0x59, 0xc7, 0x7e, 0x75, 0xec, 0x5e, 0xf7, 0xb9, 0x5f, 0xe0, 0x4b, 0x74, 0x5c,
	0x3b, 0xfe, 0xba, 0x13, 0x3d, 0x59, 0x78, 0x51, 0x54, 0x3b, 0xfe, 0x28, 0x7a, 0x7c, 0x8d, 0x4e,
	0x5b, 0x25, 0x75, 0x61, 0x5d, 0x23, 0x48, 0xe4, 0xab, 0xbf, 0xe0, 0xe1, 0xee, 0x73, 0xa0, 0xe0,
	0x30, 0x50, 0xf0, 0x3d, 0x50, 0xf0, 0x3e, 0xd2, 0xe0, 0x30, 0xd2, 0xe0, 0x6b, 0xa4, 0xc1, 0x4b,
	0x22, 0x95, 0xdd, 0x3a, 0xce, 0x4a, 0x53, 0x65, 0xa5, 0xba, 0x5d, 0xad, 0xcb, 0xad, 0x2a, 0xf4,
	0x3f, 0xe4, 0x91, 0x7f, 0x63, 0xfd, 0x13, 0x00, 0x00, 0xff, 0xff, 0x22, 0xc2, 0xf2, 0x23, 0x31,
	0x01, 0x00, 0x00,
}

func (m *PbTx) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *PbTx) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *PbTx) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Signature) > 0 {
		i -= len(m.Signature)
		copy(dAtA[i:], m.Signature)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Signature)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.PubKey) > 0 {
		i -= len(m.PubKey)
		copy(dAtA[i:], m.PubKey)
		i = encodeVarintTx(dAtA, i, uint64(len(m.PubKey)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Msgs) > 0 {
		for iNdEx := len(m.Msgs) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Msgs[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTx(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x22
		}
	}
	if m.Gas != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.Gas))
		i--
		dAtA[i] = 0x18
	}
	if m.Nonce != 0 {
		i = encodeVarintTx(dAtA, i, uint64(m.Nonce))
		i--
		dAtA[i] = 0x10
	}
	if len(m.From) > 0 {
		i -= len(m.From)
		copy(dAtA[i:], m.From)
		i = encodeVarintTx(dAtA, i, uint64(len(m.From)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *PbTx) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.From)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if m.Nonce != 0 {
		n += 1 + sovTx(uint64(m.Nonce))
	}
	if m.Gas != 0 {
		n += 1 + sovTx(uint64(m.Gas))
	}
	if len(m.Msgs) > 0 {
		for _, e := range m.Msgs {
			l = e.Size()
			n += 1 + l + sovTx(uint64(l))
		}
	}
	l = len(m.PubKey)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Signature)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *PbTx) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: PbTx: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: PbTx: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field From", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.From = append(m.From[:0], dAtA[iNdEx:postIndex]...)
			if m.From == nil {
				m.From = []byte{}
			}
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Nonce", wireType)
			}
			m.Nonce = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Nonce |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Gas", wireType)
			}
			m.Gas = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Gas |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Msgs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Msgs = append(m.Msgs, &codectypes.Any{})
			if err := m.Msgs[len(m.Msgs)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PubKey", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PubKey = append(m.PubKey[:0], dAtA[iNdEx:postIndex]...)
			if m.PubKey == nil {
				m.PubKey = []byte{}
			}
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Signature", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Signature = append(m.Signature[:0], dAtA[iNdEx:postIndex]...)
			if m.Signature == nil {
				m.Signature = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)