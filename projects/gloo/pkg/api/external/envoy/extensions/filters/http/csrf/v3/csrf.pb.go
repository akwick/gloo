// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/solo-io/gloo/projects/gloo/api/external/envoy/extensions/filters/http/csrf/v3/csrf.proto

package v3

import (
	bytes "bytes"
	fmt "fmt"
	math "math"

	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	v3 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/config/core/v3"
	v31 "github.com/solo-io/gloo/projects/gloo/pkg/api/external/envoy/type/matcher/v3"
	_ "github.com/solo-io/gloo/projects/gloo/pkg/api/external/udpa/annotations"
	_ "github.com/solo-io/protoc-gen-ext/extproto"
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

// CSRF filter config.
type CsrfPolicy struct {
	// Specifies the % of requests for which the CSRF filter is enabled.
	//
	// If :ref:`runtime_key <envoy_api_field_config.core.v3.RuntimeFractionalPercent.runtime_key>` is specified,
	// Envoy will lookup the runtime key to get the percentage of requests to filter.
	//
	// .. note::
	//
	//   This field defaults to 100/:ref:`HUNDRED
	//   <envoy_api_enum_type.v3.FractionalPercent.DenominatorType>`.
	FilterEnabled *v3.RuntimeFractionalPercent `protobuf:"bytes,1,opt,name=filter_enabled,json=filterEnabled,proto3" json:"filter_enabled,omitempty"`
	// Specifies that CSRF policies will be evaluated and tracked, but not enforced.
	//
	// This is intended to be used when ``filter_enabled`` is off and will be ignored otherwise.
	//
	// If :ref:`runtime_key <envoy_api_field_config.core.v3.RuntimeFractionalPercent.runtime_key>` is specified,
	// Envoy will lookup the runtime key to get the percentage of requests for which it will evaluate
	// and track the request's *Origin* and *Destination* to determine if it's valid, but will not
	// enforce any policies.
	ShadowEnabled *v3.RuntimeFractionalPercent `protobuf:"bytes,2,opt,name=shadow_enabled,json=shadowEnabled,proto3" json:"shadow_enabled,omitempty"`
	// Specifies additional source origins that will be allowed in addition to
	// the destination origin.
	//
	// More information on how this can be configured via runtime can be found
	// :ref:`here <csrf-configuration>`.
	AdditionalOrigins    []*v31.StringMatcher `protobuf:"bytes,3,rep,name=additional_origins,json=additionalOrigins,proto3" json:"additional_origins,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *CsrfPolicy) Reset()         { *m = CsrfPolicy{} }
func (m *CsrfPolicy) String() string { return proto.CompactTextString(m) }
func (*CsrfPolicy) ProtoMessage()    {}
func (*CsrfPolicy) Descriptor() ([]byte, []int) {
	return fileDescriptor_067a8f826086f3e1, []int{0}
}
func (m *CsrfPolicy) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CsrfPolicy.Unmarshal(m, b)
}
func (m *CsrfPolicy) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CsrfPolicy.Marshal(b, m, deterministic)
}
func (m *CsrfPolicy) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CsrfPolicy.Merge(m, src)
}
func (m *CsrfPolicy) XXX_Size() int {
	return xxx_messageInfo_CsrfPolicy.Size(m)
}
func (m *CsrfPolicy) XXX_DiscardUnknown() {
	xxx_messageInfo_CsrfPolicy.DiscardUnknown(m)
}

var xxx_messageInfo_CsrfPolicy proto.InternalMessageInfo

func (m *CsrfPolicy) GetFilterEnabled() *v3.RuntimeFractionalPercent {
	if m != nil {
		return m.FilterEnabled
	}
	return nil
}

func (m *CsrfPolicy) GetShadowEnabled() *v3.RuntimeFractionalPercent {
	if m != nil {
		return m.ShadowEnabled
	}
	return nil
}

func (m *CsrfPolicy) GetAdditionalOrigins() []*v31.StringMatcher {
	if m != nil {
		return m.AdditionalOrigins
	}
	return nil
}

func init() {
	proto.RegisterType((*CsrfPolicy)(nil), "envoy.extensions.filters.http.csrf.v3.CsrfPolicy")
}

func init() {
	proto.RegisterFile("github.com/solo-io/gloo/projects/gloo/api/external/envoy/extensions/filters/http/csrf/v3/csrf.proto", fileDescriptor_067a8f826086f3e1)
}

var fileDescriptor_067a8f826086f3e1 = []byte{
	// 463 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x92, 0xc1, 0x6a, 0x14, 0x41,
	0x10, 0x86, 0x99, 0x89, 0x06, 0x9d, 0x90, 0x10, 0x07, 0xc1, 0x10, 0x30, 0xc6, 0xa0, 0x10, 0x90,
	0x74, 0xc3, 0xce, 0xcd, 0xe3, 0x8a, 0xde, 0xc4, 0x65, 0x83, 0x17, 0x11, 0x96, 0xde, 0x9e, 0xda,
	0xd9, 0xd6, 0xd9, 0xae, 0xa1, 0xbb, 0x76, 0xdc, 0xbd, 0x79, 0x14, 0x1f, 0xc1, 0x27, 0xf0, 0xe8,
	0x51, 0xbc, 0x0b, 0x5e, 0x3c, 0xf8, 0x0a, 0xbe, 0x83, 0x17, 0x4f, 0xd2, 0x5d, 0x13, 0x47, 0xc9,
	0x25, 0xe4, 0xd4, 0xd5, 0x55, 0xf5, 0x7f, 0x5d, 0xd5, 0x55, 0x99, 0xae, 0x0c, 0xcd, 0x97, 0x53,
	0xa1, 0x71, 0x21, 0x3d, 0xd6, 0x78, 0x62, 0x50, 0x56, 0x35, 0xa2, 0x6c, 0x1c, 0xbe, 0x02, 0x4d,
	0x9e, 0x6f, 0xaa, 0x31, 0x12, 0x56, 0x04, 0xce, 0xaa, 0x5a, 0x82, 0x6d, 0x71, 0x1d, 0xaf, 0xd6,
	0x1b, 0xb4, 0x5e, 0xce, 0x4c, 0x4d, 0xe0, 0xbc, 0x9c, 0x13, 0x35, 0x52, 0x7b, 0x37, 0x93, 0x6d,
	0x11, 0x4f, 0xd1, 0x38, 0x24, 0xcc, 0xef, 0x47, 0x85, 0xe8, 0x15, 0xa2, 0x53, 0x88, 0xa0, 0x10,
	0x31, 0xb3, 0x2d, 0xf6, 0xef, 0x30, 0x58, 0xa3, 0x9d, 0x99, 0x4a, 0x6a, 0x74, 0x10, 0x38, 0x53,
	0xe5, 0x81, 0x39, 0xfb, 0x47, 0x9c, 0x40, 0xeb, 0x06, 0xe4, 0x42, 0x91, 0x9e, 0x83, 0x0b, 0x19,
	0x9e, 0x9c, 0xb1, 0x55, 0x97, 0x73, 0x7b, 0x59, 0x36, 0x4a, 0x2a, 0x6b, 0x91, 0x14, 0xc5, 0xea,
	0x3c, 0x29, 0x5a, 0xfa, 0x2e, 0x7c, 0xf7, 0x5c, 0xb8, 0x05, 0x17, 0x6a, 0xea, 0x09, 0xb7, 0x5a,
	0x55, 0x9b, 0x52, 0x11, 0xc8, 0x33, 0xa3, 0x0b, 0xdc, 0xac, 0xb0, 0xc2, 0x68, 0xca, 0x60, 0x75,
	0xde, 0x1c, 0x56, 0xc4, 0x4e, 0x58, 0x11, 0xfb, 0x8e, 0xbe, 0xa7, 0x59, 0xf6, 0xc8, 0xbb, 0xd9,
	0x08, 0x6b, 0xa3, 0xd7, 0xf9, 0x24, 0xdb, 0xe1, 0x86, 0x27, 0x60, 0xd5, 0xb4, 0x86, 0x72, 0x2f,
	0x39, 0x4c, 0x8e, 0xb7, 0x06, 0x42, 0xf0, 0xc7, 0x70, 0xc7, 0x22, 0x74, 0x2c, 0xda, 0x42, 0x8c,
	0x97, 0x96, 0xcc, 0x02, 0x9e, 0x38, 0xa5, 0x43, 0x89, 0xaa, 0x1e, 0x81, 0xd3, 0x60, 0x69, 0x78,
	0xed, 0xf7, 0xf0, 0xea, 0xfb, 0x24, 0xdd, 0x4d, 0xc6, 0xdb, 0xcc, 0x7b, 0xcc, 0xb8, 0xfc, 0x79,
	0xb6, 0xe3, 0xe7, 0xaa, 0xc4, 0x37, 0x7f, 0x1f, 0x48, 0x2f, 0xf3, 0xc0, 0x78, 0x9b, 0x29, 0x67,
	0xd8, 0xd3, 0x2c, 0x57, 0x65, 0x69, 0x38, 0x67, 0x82, 0xce, 0x54, 0xc6, 0xfa, 0xbd, 0x8d, 0xc3,
	0x8d, 0xe3, 0xad, 0xc1, 0xbd, 0x0e, 0x1d, 0x86, 0x21, 0xba, 0x61, 0x04, 0xf6, 0x69, 0x1c, 0xc6,
	0x53, 0x76, 0x8c, 0x6f, 0xf4, 0xfa, 0x67, 0x2c, 0x7f, 0x38, 0xf8, 0xf0, 0xf5, 0xdd, 0xc1, 0x49,
	0xf6, 0xe0, 0xbf, 0xca, 0xb8, 0x9d, 0x7f, 0xd7, 0x61, 0x20, 0xfa, 0x0f, 0x1c, 0x7e, 0x4a, 0x3e,
	0xff, 0xba, 0x92, 0x7c, 0xfc, 0x79, 0x90, 0x7c, 0x79, 0xfb, 0xed, 0xc7, 0x66, 0xba, 0x9b, 0x66,
	0x85, 0x41, 0xae, 0xa0, 0x71, 0xb8, 0x5a, 0x8b, 0x0b, 0x6d, 0xd8, 0xf0, 0x7a, 0x64, 0x86, 0x11,
	0x8d, 0x92, 0x17, 0x2f, 0x2f, 0xb6, 0xfa, 0xcd, 0xeb, 0xea, 0x12, 0xeb, 0x3f, 0xdd, 0x8c, 0x9b,
	0x50, 0xfc, 0x09, 0x00, 0x00, 0xff, 0xff, 0x1a, 0x95, 0xd3, 0x6f, 0x61, 0x03, 0x00, 0x00,
}

func (this *CsrfPolicy) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*CsrfPolicy)
	if !ok {
		that2, ok := that.(CsrfPolicy)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if !this.FilterEnabled.Equal(that1.FilterEnabled) {
		return false
	}
	if !this.ShadowEnabled.Equal(that1.ShadowEnabled) {
		return false
	}
	if len(this.AdditionalOrigins) != len(that1.AdditionalOrigins) {
		return false
	}
	for i := range this.AdditionalOrigins {
		if !this.AdditionalOrigins[i].Equal(that1.AdditionalOrigins[i]) {
			return false
		}
	}
	if !bytes.Equal(this.XXX_unrecognized, that1.XXX_unrecognized) {
		return false
	}
	return true
}
