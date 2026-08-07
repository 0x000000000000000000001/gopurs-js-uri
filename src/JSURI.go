package JSURI

import (
	"strings"
	"unicode/utf8"

	"gopurs/output/Effect.Exception"
	"gopurs/output/gopurs_runtime"
)

func isURIComponentSafe(c byte) bool {
	if c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' {
		return true
	}
	switch c {
	case '-', '_', '.', '!', '~', '*', '\'', '(', ')':
		return true
	}
	return false
}

func encodeStr(str string, form bool, safe func(byte) bool) (string, bool) {
	if !utf8.ValidString(str) {
		return "", false
	}
	var builder strings.Builder
	for i := 0; i < len(str); i++ {
		c := str[i]
		if form && c == ' ' {
			builder.WriteByte('+')
		} else if safe(c) {
			builder.WriteByte(c)
		} else {
			hex := "0123456789ABCDEF"
			builder.WriteByte('%')
			builder.WriteByte(hex[c>>4])
			builder.WriteByte(hex[c&15])
		}
	}
	return builder.String(), true
}

func decodeStr(str string, form bool) (string, bool) {
	// Simple URI decoding
	var builder strings.Builder
	for i := 0; i < len(str); i++ {
		c := str[i]
		if form && c == '+' {
			builder.WriteByte(' ')
		} else if c == '%' {
			if i+2 < len(str) {
				var b byte
				v1 := fromHexChar(str[i+1])
				v2 := fromHexChar(str[i+2])
				if v1 == 255 || v2 == 255 {
					return "", false
				}
				b = (v1 << 4) | v2
				builder.WriteByte(b)
				i += 2
			} else {
				return "", false
			}
		} else {
			builder.WriteByte(c)
		}
	}
	res := builder.String()
	if !utf8.ValidString(res) {
		return "", false
	}
	return res, true
}

func fromHexChar(c byte) byte {
	if c >= '0' && c <= '9' {
		return c - '0'
	}
	if c >= 'a' && c <= 'f' {
		return c - 'a' + 10
	}
	if c >= 'A' && c <= 'F' {
		return c - 'A' + 10
	}
	return 255
}

var _EncodeURIComponent = gopurs_runtime.Func3(func(fail, succeed, input gopurs_runtime.Value) gopurs_runtime.Value {
	str := input.StrVal()
	res, ok := encodeStr(str, false, isURIComponentSafe)
	if !ok {
		return gopurs_runtime.Apply(fail, gopurs_runtime.Str("URI malformed"))
	}
	return gopurs_runtime.Apply(succeed, gopurs_runtime.Str(res))
})

var _EncodeFormURLComponent = gopurs_runtime.Func3(func(fail, succeed, input gopurs_runtime.Value) gopurs_runtime.Value {
	str := input.StrVal()
	res, ok := encodeStr(str, true, isURIComponentSafe)
	if !ok {
		return gopurs_runtime.Apply(fail, gopurs_runtime.Str("URI malformed"))
	}
	return gopurs_runtime.Apply(succeed, gopurs_runtime.Str(res))
})

var _DecodeURIComponent = gopurs_runtime.Func3(func(fail, succeed, input gopurs_runtime.Value) gopurs_runtime.Value {
	str := input.StrVal()
	res, ok := decodeStr(str, false)
	if !ok {
		errVal := gopurs_runtime.Apply(Effect_Exception.Get_error(), gopurs_runtime.Str("URI malformed"))
		return gopurs_runtime.Apply(fail, errVal)
	}
	return gopurs_runtime.Apply(succeed, gopurs_runtime.Str(res))
})

var _DecodeFormURLComponent = gopurs_runtime.Func3(func(fail, succeed, input gopurs_runtime.Value) gopurs_runtime.Value {
	str := input.StrVal()
	res, ok := decodeStr(str, true)
	if !ok {
		errVal := gopurs_runtime.Apply(Effect_Exception.Get_error(), gopurs_runtime.Str("URI malformed"))
		return gopurs_runtime.Apply(fail, errVal)
	}
	return gopurs_runtime.Apply(succeed, gopurs_runtime.Str(res))
})

func isURISafe(c byte) bool {
	if c >= 'A' && c <= 'Z' || c >= 'a' && c <= 'z' || c >= '0' && c <= '9' {
		return true
	}
	switch c {
	case ';', ',', '/', '?', ':', '@', '&', '=', '+', '$', '#', '-', '_', '.', '!', '~', '*', '\'', '(', ')':
		return true
	}
	return false
}

var _EncodeURI = gopurs_runtime.Func3(func(fail, succeed, input gopurs_runtime.Value) gopurs_runtime.Value {
	str := input.StrVal()
	res, ok := encodeStr(str, false, isURISafe)
	if !ok {
		return gopurs_runtime.Apply(fail, gopurs_runtime.Str("URI malformed"))
	}
	return gopurs_runtime.Apply(succeed, gopurs_runtime.Str(res))
})

var _DecodeURI = gopurs_runtime.Func3(func(fail, succeed, input gopurs_runtime.Value) gopurs_runtime.Value {
	str := input.StrVal()
	res, ok := decodeStr(str, false)
	if !ok {
		errVal := gopurs_runtime.Apply(Effect_Exception.Get_error(), gopurs_runtime.Str("URI malformed"))
		return gopurs_runtime.Apply(fail, errVal)
	}
	return gopurs_runtime.Apply(succeed, gopurs_runtime.Str(res))
})
