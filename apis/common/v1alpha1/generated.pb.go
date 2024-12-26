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
*/
// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/kuidio/kuid/apis/common/v1alpha1/generated.proto

package v1alpha1

import (
	fmt "fmt"

	io "io"
	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"

	proto "github.com/gogo/protobuf/proto"
	github_com_gogo_protobuf_sortkeys "github.com/gogo/protobuf/sortkeys"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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

func (m *ClaimLabels) Reset()      { *m = ClaimLabels{} }
func (*ClaimLabels) ProtoMessage() {}
func (*ClaimLabels) Descriptor() ([]byte, []int) {
	return fileDescriptor_ff64f0e0dd0fa4cf, []int{0}
}
func (m *ClaimLabels) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ClaimLabels) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *ClaimLabels) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ClaimLabels.Merge(m, src)
}
func (m *ClaimLabels) XXX_Size() int {
	return m.Size()
}
func (m *ClaimLabels) XXX_DiscardUnknown() {
	xxx_messageInfo_ClaimLabels.DiscardUnknown(m)
}

var xxx_messageInfo_ClaimLabels proto.InternalMessageInfo

func (m *UserDefinedLabels) Reset()      { *m = UserDefinedLabels{} }
func (*UserDefinedLabels) ProtoMessage() {}
func (*UserDefinedLabels) Descriptor() ([]byte, []int) {
	return fileDescriptor_ff64f0e0dd0fa4cf, []int{1}
}
func (m *UserDefinedLabels) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UserDefinedLabels) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *UserDefinedLabels) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UserDefinedLabels.Merge(m, src)
}
func (m *UserDefinedLabels) XXX_Size() int {
	return m.Size()
}
func (m *UserDefinedLabels) XXX_DiscardUnknown() {
	xxx_messageInfo_UserDefinedLabels.DiscardUnknown(m)
}

var xxx_messageInfo_UserDefinedLabels proto.InternalMessageInfo

func init() {
	proto.RegisterType((*ClaimLabels)(nil), "github.com.kuidio.kuid.apis.common.v1alpha1.ClaimLabels")
	proto.RegisterType((*UserDefinedLabels)(nil), "github.com.kuidio.kuid.apis.common.v1alpha1.UserDefinedLabels")
	proto.RegisterMapType((map[string]string)(nil), "github.com.kuidio.kuid.apis.common.v1alpha1.UserDefinedLabels.LabelsEntry")
}

func init() {
	proto.RegisterFile("github.com/kuidio/kuid/apis/common/v1alpha1/generated.proto", fileDescriptor_ff64f0e0dd0fa4cf)
}

var fileDescriptor_ff64f0e0dd0fa4cf = []byte{
	// 388 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x52, 0x3d, 0x6f, 0xda, 0x40,
	0x18, 0xf6, 0x81, 0x8a, 0xe0, 0x5c, 0x55, 0xc5, 0xea, 0x40, 0x19, 0x8e, 0x8a, 0x09, 0x09, 0xf5,
	0x4e, 0xd0, 0x0e, 0xb4, 0x95, 0x3a, 0xb8, 0xcd, 0x12, 0x65, 0x89, 0xa3, 0x2c, 0x91, 0x32, 0x1c,
	0xe6, 0x30, 0x27, 0x7f, 0x9c, 0xe5, 0x2f, 0x89, 0x2d, 0x5b, 0xd6, 0xfc, 0xa6, 0x4c, 0x8c, 0x8c,
	0x4c, 0x28, 0x38, 0xff, 0x21, 0x73, 0xe4, 0x3b, 0x93, 0xa0, 0x78, 0x42, 0x99, 0xde, 0xf7, 0xb1,
	0xfd, 0x7c, 0xbd, 0x32, 0xfc, 0xe3, 0xf0, 0x64, 0x91, 0x4e, 0xb1, 0x2d, 0x7c, 0xe2, 0xa6, 0x7c,
	0xc6, 0x85, 0x1c, 0x84, 0x86, 0x3c, 0x26, 0xb6, 0xf0, 0x7d, 0x11, 0x90, 0x6c, 0x44, 0xbd, 0x70,
	0x41, 0x47, 0xc4, 0x61, 0x01, 0x8b, 0x68, 0xc2, 0x66, 0x38, 0x8c, 0x44, 0x22, 0x8c, 0xe1, 0x2b,
	0x19, 0x2b, 0xb2, 0x1c, 0xb8, 0x20, 0x63, 0x45, 0xc6, 0x7b, 0x72, 0xf7, 0xfb, 0x81, 0x93, 0x23,
	0x1c, 0x41, 0xa4, 0xc6, 0x34, 0x9d, 0x4b, 0x24, 0x81, 0xdc, 0x94, 0x76, 0xf7, 0xa7, 0x3b, 0x89,
	0x31, 0x17, 0x45, 0x10, 0x9f, 0xda, 0x0b, 0x1e, 0xb0, 0x68, 0x49, 0x42, 0xd7, 0x51, 0xc9, 0x7c,
	0x96, 0x50, 0x92, 0x55, 0x12, 0xf5, 0x9f, 0x00, 0xd4, 0xff, 0x79, 0x94, 0xfb, 0x67, 0x74, 0xca,
	0xbc, 0xd8, 0xb8, 0x05, 0xb0, 0x9d, 0xc6, 0x2c, 0xfa, 0xcf, 0xe6, 0x3c, 0x60, 0x33, 0xf5, 0xb4,
	0x03, 0xbe, 0x81, 0x81, 0x3e, 0xfe, 0x8b, 0x8f, 0x88, 0x8f, 0x2f, 0xdf, 0xaa, 0x98, 0x5f, 0x57,
	0xdb, 0x9e, 0x96, 0x6f, 0x7b, 0xed, 0xca, 0x2b, 0xab, 0xea, 0x69, 0x5c, 0xc3, 0x66, 0xcc, 0x3c,
	0x66, 0x27, 0x22, 0xea, 0xd4, 0xa4, 0xff, 0x0f, 0xac, 0x2a, 0xe2, 0xc3, 0x8a, 0x38, 0x74, 0x1d,
	0x15, 0xa0, 0xa8, 0x88, 0xb3, 0x11, 0x96, 0xfc, 0x8b, 0x92, 0x6a, 0x7e, 0xcc, 0xb7, 0xbd, 0xe6,
	0x1e, 0x59, 0x2f, 0x92, 0xfd, 0x7b, 0x00, 0xab, 0x39, 0x8c, 0x08, 0x36, 0xbc, 0x7d, 0xe5, 0xfa,
	0x40, 0x1f, 0x9f, 0xbe, 0xaf, 0xb2, 0xca, 0x12, 0x9f, 0x04, 0x49, 0xb4, 0x34, 0x3f, 0x95, 0xf5,
	0x1b, 0x65, 0xe7, 0xd2, 0xa9, 0xfb, 0x0b, 0xea, 0x07, 0x9f, 0x19, 0x9f, 0x61, 0xdd, 0x65, 0x4b,
	0x79, 0xf2, 0x96, 0x55, 0xac, 0xc6, 0x17, 0xf8, 0x21, 0xa3, 0x5e, 0xca, 0xe4, 0x19, 0x5a, 0x96,
	0x02, 0xbf, 0x6b, 0x13, 0x60, 0x9e, 0xaf, 0x76, 0x48, 0x5b, 0xef, 0x90, 0xb6, 0xd9, 0x21, 0xed,
	0x26, 0x47, 0x60, 0x95, 0x23, 0xb0, 0xce, 0x11, 0xd8, 0xe4, 0x08, 0x3c, 0xe4, 0x08, 0xdc, 0x3d,
	0x22, 0xed, 0x6a, 0x78, 0xc4, 0x2f, 0xfb, 0x1c, 0x00, 0x00, 0xff, 0xff, 0x24, 0x9f, 0x6b, 0xba,
	0xe0, 0x02, 0x00, 0x00,
}

func (m *ClaimLabels) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ClaimLabels) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ClaimLabels) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Selector != nil {
		{
			size, err := m.Selector.MarshalToSizedBuffer(dAtA[:i])
			if err != nil {
				return 0, err
			}
			i -= size
			i = encodeVarintGenerated(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	{
		size, err := m.UserDefinedLabels.MarshalToSizedBuffer(dAtA[:i])
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

func (m *UserDefinedLabels) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UserDefinedLabels) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UserDefinedLabels) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Labels) > 0 {
		keysForLabels := make([]string, 0, len(m.Labels))
		for k := range m.Labels {
			keysForLabels = append(keysForLabels, string(k))
		}
		github_com_gogo_protobuf_sortkeys.Strings(keysForLabels)
		for iNdEx := len(keysForLabels) - 1; iNdEx >= 0; iNdEx-- {
			v := m.Labels[string(keysForLabels[iNdEx])]
			baseI := i
			i -= len(v)
			copy(dAtA[i:], v)
			i = encodeVarintGenerated(dAtA, i, uint64(len(v)))
			i--
			dAtA[i] = 0x12
			i -= len(keysForLabels[iNdEx])
			copy(dAtA[i:], keysForLabels[iNdEx])
			i = encodeVarintGenerated(dAtA, i, uint64(len(keysForLabels[iNdEx])))
			i--
			dAtA[i] = 0xa
			i = encodeVarintGenerated(dAtA, i, uint64(baseI-i))
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
func (m *ClaimLabels) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.UserDefinedLabels.Size()
	n += 1 + l + sovGenerated(uint64(l))
	if m.Selector != nil {
		l = m.Selector.Size()
		n += 1 + l + sovGenerated(uint64(l))
	}
	return n
}

func (m *UserDefinedLabels) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Labels) > 0 {
		for k, v := range m.Labels {
			_ = k
			_ = v
			mapEntrySize := 1 + len(k) + sovGenerated(uint64(len(k))) + 1 + len(v) + sovGenerated(uint64(len(v)))
			n += mapEntrySize + 1 + sovGenerated(uint64(mapEntrySize))
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
func (this *ClaimLabels) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&ClaimLabels{`,
		`UserDefinedLabels:` + strings.Replace(strings.Replace(this.UserDefinedLabels.String(), "UserDefinedLabels", "UserDefinedLabels", 1), `&`, ``, 1) + `,`,
		`Selector:` + strings.Replace(fmt.Sprintf("%v", this.Selector), "LabelSelector", "v1.LabelSelector", 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *UserDefinedLabels) String() string {
	if this == nil {
		return "nil"
	}
	keysForLabels := make([]string, 0, len(this.Labels))
	for k := range this.Labels {
		keysForLabels = append(keysForLabels, k)
	}
	github_com_gogo_protobuf_sortkeys.Strings(keysForLabels)
	mapStringForLabels := "map[string]string{"
	for _, k := range keysForLabels {
		mapStringForLabels += fmt.Sprintf("%v: %v,", k, this.Labels[k])
	}
	mapStringForLabels += "}"
	s := strings.Join([]string{`&UserDefinedLabels{`,
		`Labels:` + mapStringForLabels + `,`,
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
func (m *ClaimLabels) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: ClaimLabels: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ClaimLabels: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field UserDefinedLabels", wireType)
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
			if err := m.UserDefinedLabels.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Selector", wireType)
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
			if m.Selector == nil {
				m.Selector = &v1.LabelSelector{}
			}
			if err := m.Selector.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *UserDefinedLabels) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: UserDefinedLabels: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UserDefinedLabels: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Labels", wireType)
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
			if m.Labels == nil {
				m.Labels = make(map[string]string)
			}
			var mapkey string
			var mapvalue string
			for iNdEx < postIndex {
				entryPreIndex := iNdEx
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
				if fieldNum == 1 {
					var stringLenmapkey uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGenerated
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapkey |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapkey := int(stringLenmapkey)
					if intStringLenmapkey < 0 {
						return ErrInvalidLengthGenerated
					}
					postStringIndexmapkey := iNdEx + intStringLenmapkey
					if postStringIndexmapkey < 0 {
						return ErrInvalidLengthGenerated
					}
					if postStringIndexmapkey > l {
						return io.ErrUnexpectedEOF
					}
					mapkey = string(dAtA[iNdEx:postStringIndexmapkey])
					iNdEx = postStringIndexmapkey
				} else if fieldNum == 2 {
					var stringLenmapvalue uint64
					for shift := uint(0); ; shift += 7 {
						if shift >= 64 {
							return ErrIntOverflowGenerated
						}
						if iNdEx >= l {
							return io.ErrUnexpectedEOF
						}
						b := dAtA[iNdEx]
						iNdEx++
						stringLenmapvalue |= uint64(b&0x7F) << shift
						if b < 0x80 {
							break
						}
					}
					intStringLenmapvalue := int(stringLenmapvalue)
					if intStringLenmapvalue < 0 {
						return ErrInvalidLengthGenerated
					}
					postStringIndexmapvalue := iNdEx + intStringLenmapvalue
					if postStringIndexmapvalue < 0 {
						return ErrInvalidLengthGenerated
					}
					if postStringIndexmapvalue > l {
						return io.ErrUnexpectedEOF
					}
					mapvalue = string(dAtA[iNdEx:postStringIndexmapvalue])
					iNdEx = postStringIndexmapvalue
				} else {
					iNdEx = entryPreIndex
					skippy, err := skipGenerated(dAtA[iNdEx:])
					if err != nil {
						return err
					}
					if (skippy < 0) || (iNdEx+skippy) < 0 {
						return ErrInvalidLengthGenerated
					}
					if (iNdEx + skippy) > postIndex {
						return io.ErrUnexpectedEOF
					}
					iNdEx += skippy
				}
			}
			m.Labels[mapkey] = mapvalue
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