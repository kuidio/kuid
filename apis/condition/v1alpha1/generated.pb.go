/*
Copyright 2024 Nokia.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/ // Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/kuidio/kuid/apis/condition/v1alpha1/generated.proto

package v1alpha1

import (
	fmt "fmt"

	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"

	proto "github.com/gogo/protobuf/proto"
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

func (m *Condition) Reset()      { *m = Condition{} }
func (*Condition) ProtoMessage() {}
func (*Condition) Descriptor() ([]byte, []int) {
	return fileDescriptor_a4ff95c4762e0ac7, []int{0}
}
func (m *Condition) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Condition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *Condition) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Condition.Merge(m, src)
}
func (m *Condition) XXX_Size() int {
	return m.Size()
}
func (m *Condition) XXX_DiscardUnknown() {
	xxx_messageInfo_Condition.DiscardUnknown(m)
}

var xxx_messageInfo_Condition proto.InternalMessageInfo

func (m *ConditionedStatus) Reset()      { *m = ConditionedStatus{} }
func (*ConditionedStatus) ProtoMessage() {}
func (*ConditionedStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_a4ff95c4762e0ac7, []int{1}
}
func (m *ConditionedStatus) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ConditionedStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *ConditionedStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConditionedStatus.Merge(m, src)
}
func (m *ConditionedStatus) XXX_Size() int {
	return m.Size()
}
func (m *ConditionedStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_ConditionedStatus.DiscardUnknown(m)
}

var xxx_messageInfo_ConditionedStatus proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Condition)(nil), "github.com.kuidio.kuid.apis.condition.v1alpha1.Condition")
	proto.RegisterType((*ConditionedStatus)(nil), "github.com.kuidio.kuid.apis.condition.v1alpha1.ConditionedStatus")
}

func init() {
	proto.RegisterFile("github.com/kuidio/kuid/apis/condition/v1alpha1/generated.proto", fileDescriptor_a4ff95c4762e0ac7)
}

var fileDescriptor_a4ff95c4762e0ac7 = []byte{
	// 297 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0x31, 0x4e, 0xc3, 0x30,
	0x14, 0x86, 0x63, 0x31, 0xd5, 0x9d, 0x9a, 0xa9, 0xea, 0xe0, 0x56, 0x9d, 0xba, 0xf0, 0xac, 0x56,
	0x0c, 0xb0, 0x30, 0x94, 0x1b, 0x14, 0x26, 0x26, 0xdc, 0xc4, 0x38, 0x56, 0x48, 0x1c, 0x25, 0x4e,
	0x24, 0x36, 0xc4, 0x09, 0x38, 0x56, 0xc6, 0x8e, 0x9d, 0x2a, 0x62, 0x2e, 0x82, 0xe2, 0x50, 0x27,
	0x12, 0x53, 0xa7, 0xe7, 0x67, 0xfb, 0xff, 0xde, 0x27, 0x1b, 0xdf, 0x0b, 0xa9, 0xa3, 0x72, 0x0f,
	0x81, 0x4a, 0x68, 0x5c, 0xca, 0x50, 0x2a, 0x5b, 0x28, 0xcb, 0x64, 0x41, 0x03, 0x95, 0x86, 0x52,
	0x4b, 0x95, 0xd2, 0x6a, 0xcd, 0xde, 0xb2, 0x88, 0xad, 0xa9, 0xe0, 0x29, 0xcf, 0x99, 0xe6, 0x21,
	0x64, 0xb9, 0xd2, 0xca, 0x87, 0x3e, 0x0f, 0x5d, 0xde, 0x16, 0x68, 0xf3, 0xe0, 0xf2, 0x70, 0xce,
	0xcf, 0xae, 0x07, 0xf3, 0x84, 0x12, 0x8a, 0x5a, 0xcc, 0xbe, 0x7c, 0xb5, 0x9d, 0x6d, 0xec, 0xaa,
	0xc3, 0xcf, 0x6e, 0xe2, 0xdb, 0x02, 0xa4, 0x6a, 0x75, 0x12, 0x16, 0x44, 0x32, 0xe5, 0xf9, 0x3b,
	0xcd, 0x62, 0xd1, 0xf9, 0x25, 0x5c, 0x33, 0x5a, 0xfd, 0x93, 0x5a, 0x26, 0x78, 0xf4, 0x70, 0x1e,
	0xed, 0xbf, 0xe0, 0x91, 0xf3, 0x98, 0xa2, 0x05, 0x5a, 0x8d, 0x37, 0x14, 0x3a, 0x2c, 0x0c, 0xb1,
	0x90, 0xc5, 0xa2, 0xd3, 0x6e, 0xb1, 0x50, 0xad, 0xc1, 0x31, 0xb6, 0x93, 0xfa, 0x34, 0xf7, 0xcc,
	0x69, 0xde, 0x63, 0x77, 0x3d, 0x74, 0xf9, 0x89, 0xf0, 0xc4, 0x1d, 0xf0, 0xf0, 0x51, 0x33, 0x5d,
	0x16, 0x7e, 0x82, 0xb1, 0xbb, 0x52, 0x4c, 0xd1, 0xe2, 0x6a, 0x35, 0xde, 0xdc, 0x5d, 0xf8, 0x5c,
	0x03, 0x05, 0xff, 0x4f, 0x01, 0xbb, 0xad, 0x62, 0x37, 0x18, 0xb0, 0x7d, 0xaa, 0x1b, 0xe2, 0x1d,
	0x1a, 0xe2, 0x1d, 0x1b, 0xe2, 0x7d, 0x18, 0x82, 0x6a, 0x43, 0xd0, 0xc1, 0x10, 0x74, 0x34, 0x04,
	0x7d, 0x1b, 0x82, 0xbe, 0x7e, 0x88, 0xf7, 0x0c, 0x97, 0x7d, 0xf7, 0x6f, 0x00, 0x00, 0x00, 0xff,
	0xff, 0x91, 0x25, 0x36, 0x7b, 0x1f, 0x02, 0x00, 0x00,
}

func (m *Condition) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Condition) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Condition) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Condition.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenerated(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *ConditionedStatus) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ConditionedStatus) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ConditionedStatus) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Conditions) > 0 {
		for iNdEx := len(m.Conditions) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Conditions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenerated(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintGenerated(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenerated(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Condition) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.Condition.Size()
	n += 1 + l + sovGenerated(uint64(l))
	return n
}

func (m *ConditionedStatus) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Conditions) > 0 {
		for _, e := range m.Conditions {
			l = e.Size()
			n += 1 + l + sovGenerated(uint64(l))
		}
	}
	return n
}

func sovGenerated(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenerated(x uint64) (n int) {
	return sovGenerated(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Condition) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Condition{`,
		`Condition:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.Condition), "Condition", "v1.Condition", 1), `&`, ``, 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *ConditionedStatus) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForConditions := "[]Condition{"
	for _, f := range this.Conditions {
		repeatedStringForConditions += strings.Replace(strings.Replace(f.String(), "Condition", "Condition", 1), `&`, ``, 1) + ","
	}
	repeatedStringForConditions += "}"
	s := strings.Join([]string{`&ConditionedStatus{`,
		`Conditions:` + repeatedStringForConditions + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringGenerated(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Condition) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
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
			return fmt.Errorf("proto: Condition: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Condition: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Condition", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
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
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Condition.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenerated
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
func (m *ConditionedStatus) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
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
			return fmt.Errorf("proto: ConditionedStatus: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ConditionedStatus: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Conditions", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
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
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Conditions = append(m.Conditions, Condition{})
			if err := m.Conditions[len(m.Conditions)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenerated
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
func skipGenerated(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenerated
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
					return 0, ErrIntOverflowGenerated
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
					return 0, ErrIntOverflowGenerated
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
				return 0, ErrInvalidLengthGenerated
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenerated
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenerated
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenerated        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenerated          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenerated = fmt.Errorf("proto: unexpected end of group")
)
