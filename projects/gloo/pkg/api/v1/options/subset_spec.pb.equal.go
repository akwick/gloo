// Code generated by protoc-gen-ext. DO NOT EDIT.
// source: github.com/solo-io/gloo/projects/gloo/api/v1/options/subset_spec.proto

package options

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"
	equality "github.com/solo-io/protoc-gen-ext/pkg/equality"
)

// ensure the imports are used
var (
	_ = errors.New("")
	_ = fmt.Print
	_ = binary.LittleEndian
	_ = bytes.Compare
	_ = strings.Compare
	_ = equality.Equalizer(nil)
	_ = proto.Message(nil)
)

// Equal function
func (m *SubsetSpec) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*SubsetSpec)
	if !ok {
		that2, ok := that.(SubsetSpec)
		if ok {
			target = &that2
		} else {
			return false
		}
	}
	if target == nil {
		return m == nil
	} else if m == nil {
		return false
	}

	if len(m.GetSelectors()) != len(target.GetSelectors()) {
		return false
	}
	for idx, v := range m.GetSelectors() {

		if h, ok := interface{}(v).(equality.Equalizer); ok {
			if !h.Equal(target.GetSelectors()[idx]) {
				return false
			}
		} else {
			if !proto.Equal(v, target.GetSelectors()[idx]) {
				return false
			}
		}

	}

	return true
}

// Equal function
func (m *Selector) Equal(that interface{}) bool {
	if that == nil {
		return m == nil
	}

	target, ok := that.(*Selector)
	if !ok {
		that2, ok := that.(Selector)
		if ok {
			target = &that2
		} else {
			return false
		}
	}
	if target == nil {
		return m == nil
	} else if m == nil {
		return false
	}

	if len(m.GetKeys()) != len(target.GetKeys()) {
		return false
	}
	for idx, v := range m.GetKeys() {

		if strings.Compare(v, target.GetKeys()[idx]) != 0 {
			return false
		}

	}

	return true
}
