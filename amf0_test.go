package amf

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

func TestEncodeAMF0(t *testing.T) {
	cases := []struct {
		in   interface{}
		want []byte
	}{
		{3.14, []byte{0x0, 0x40, 0x9, 0x1e, 0xb8, 0x51, 0xeb, 0x85, 0x1f}},
		{1, []byte{0x00, 0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{-1, []byte{0x00, 0xbf, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{true, []byte{0x01, 0x01}},
		{false, []byte{0x01, 0x00}},
		{"foo", []byte{0x02, 0x00, 0x03, 0x66, 0x6f, 0x6f}},
		{"", []byte{0x02, 0x00, 0x00}},
		{nil, []byte{0x05}},
		{map[string]interface{}{
			"1": 1,
			"2": 3.14,
			"3": "three",
			"4": nil,
			"5": true}, []byte{0x03,
			0x02, 0x00, 0x01, 0x31, 0x00, 0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x02, 0x00, 0x01, 0x32, 0x00, 0x40, 0x9, 0x1e, 0xb8, 0x51, 0xeb, 0x85, 0x1f,
			0x02, 0x00, 0x01, 0x33, 0x02, 0x00, 0x05, 0x74, 0x68, 0x72, 0x65, 0x65,
			0x02, 0x00, 0x01, 0x34, 0x05,
			0x02, 0x00, 0x01, 0x35, 0x01, 0x01,
			0x00, 0x00, 0x09}},
		{Amf0ECMAArray(map[string]interface{}{
			"1": 1,
			"2": 3.14,
			"3": "three",
			"4": nil,
			"5": true}), []byte{0x08, 0x00, 0x00, 0x00, 0x05,
			0x02, 0x00, 0x01, 0x31, 0x00, 0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x02, 0x00, 0x01, 0x32, 0x00, 0x40, 0x9, 0x1e, 0xb8, 0x51, 0xeb, 0x85, 0x1f,
			0x02, 0x00, 0x01, 0x33, 0x02, 0x00, 0x05, 0x74, 0x68, 0x72, 0x65, 0x65,
			0x02, 0x00, 0x01, 0x34, 0x05,
			0x02, 0x00, 0x01, 0x35, 0x01, 0x01}},
		{time.Unix(123456789, 123456789), []byte{0x0b, 0x00, 0x00, 0x00, 0x1c, 0xbe, 0x99, 0x1a, 0x83, 0x00, 0x00}},
		{[]interface{}{"one", "two", "three"}, []byte{0x0a, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00,
			0x00, 0x02, 0x00, 0x03, 0x6f, 0x6e, 0x65, 0x02, 0x00, 0x03, 0x74, 0x77, 0x6f, 0x02, 0x00,
			0x05, 0x74, 0x68, 0x72, 0x65, 0x65}},
		{[]interface{}{1, 2, 1}, []byte{0x0a, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		{[]interface{}{true, false, true}, []byte{0x0a, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00,
			0x01, 0x01, 0x01, 0x00, 0x01, 0x01}},
	}
	for _, c := range cases {
		got := EncodeAMF0(c.in)
		if bytes.Compare(got, c.want) != 0 {
			t.Errorf("EncodeAMF0(%#v) == %#v (%d), want %#v (%d)", c.in, got, len(got), c.want, len(c.want))
		}
	}
}

func TestDecodeAMF0(t *testing.T) {
	cases := []struct {
		in   []byte
		want interface{}
	}{
		{[]byte{0x00, 0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, 1.0},
		{[]byte{0x00, 0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, float64(1)},
		{[]byte{0x00, 0xbf, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, float64(-1)},
		{[]byte{0x01, 0x01}, true},
		{[]byte{0x01, 0x00}, false},
		{[]byte{0x02, 0x00, 0x03, 0x66, 0x6f, 0x6f}, "foo"},
		{[]byte{0x02, 0x00, 0x00}, ""},
		{[]byte{0x05}, nil},
		{[]byte{0x03,
			0x02, 0x00, 0x01, 0x31, 0x00, 0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x02, 0x00, 0x01, 0x32, 0x00, 0x40, 0x9, 0x1e, 0xb8, 0x51, 0xeb, 0x85, 0x1f,
			0x02, 0x00, 0x01, 0x33, 0x02, 0x00, 0x05, 0x74, 0x68, 0x72, 0x65, 0x65,
			0x02, 0x00, 0x01, 0x34, 0x05,
			0x02, 0x00, 0x01, 0x35, 0x01, 0x01,
			0x00, 0x00, 0x09}, map[string]interface{}{
			"1": 1.0, // AMF0 only has Number
			"2": 3.14,
			"3": "three",
			"4": nil,
			"5": true}},
		{[]byte{0x08, 0x00, 0x00, 0x00, 0x05,
			0x02, 0x00, 0x01, 0x31, 0x00, 0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x02, 0x00, 0x01, 0x32, 0x00, 0x40, 0x9, 0x1e, 0xb8, 0x51, 0xeb, 0x85, 0x1f,
			0x02, 0x00, 0x01, 0x33, 0x02, 0x00, 0x05, 0x74, 0x68, 0x72, 0x65, 0x65,
			0x02, 0x00, 0x01, 0x34, 0x05,
			0x02, 0x00, 0x01, 0x35, 0x01, 0x01}, Amf0ECMAArray(map[string]interface{}{
			"1": 1.0, // AMF0 only has Number
			"2": 3.14,
			"3": "three",
			"4": nil,
			"5": true})},
		{[]byte{0x0b, 0x00, 0x00, 0x00, 0x1c, 0xbe, 0x99, 0x1a, 0x83, 0x00, 0x00}, time.Unix(123456789, 123000000)},
		{[]byte{0x0a, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00,
			0x00, 0x02, 0x00, 0x03, 0x6f, 0x6e, 0x65, 0x02, 0x00, 0x03, 0x74, 0x77, 0x6f, 0x02, 0x00,
			0x05, 0x74, 0x68, 0x72, 0x65, 0x65}, []interface{}{"one", "two", "three"}},
		{[]byte{0x0a, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00,
			0x00, 0x00, 0x00, 0x3f, 0xf0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			[]interface{}{float64(1), float64(2), float64(1)}},
		{[]byte{0x0a, 0x00, 0x00, 0x00, 0x03, 0x00, 0x00, 0x00, 0x00,
			0x01, 0x01, 0x01, 0x00, 0x01, 0x01}, []interface{}{true, false, true}},
	}
	for _, c := range cases {
		got := DecodeAMF0(c.in)
		switch got.(type) {
		case []interface{}, map[string]interface{}, Amf0ECMAArray:
			if !reflect.DeepEqual(c.want, got) {
				t.Errorf("DecodeAMF0(%#v) == %#v, want %#v", c.in, got, c.want)
			}
			continue
		}
		if got != c.want {
			t.Errorf("DecodeAMF0(%#v) == %#v, want %#v", c.in, got, c.want)
		}
	}
}
