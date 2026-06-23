package lua

import (
	"fmt"
	"strconv"
	"strings"
)

// unquoteLuaString converts a Lua string-literal token (including its quotes or
// long brackets) to its raw byte value. ok is false for forms it cannot safely
// decode (backtick interpolation, unknown escapes, malformed input) — callers
// then leave the literal unencrypted.
func unquoteLuaString(text string) (string, bool) {
	if len(text) < 2 {
		return "", false
	}
	if text[0] == '[' {
		return unquoteLongString(text)
	}
	q := text[0]
	if q != '"' && q != '\'' || text[len(text)-1] != q {
		return "", false
	}
	body := text[1 : len(text)-1]
	var b strings.Builder
	for i := 0; i < len(body); {
		c := body[i]
		if c != '\\' {
			b.WriteByte(c)
			i++
			continue
		}
		i++
		if i >= len(body) {
			return "", false
		}
		switch e := body[i]; {
		case e == 'n':
			b.WriteByte('\n')
			i++
		case e == 't':
			b.WriteByte('\t')
			i++
		case e == 'r':
			b.WriteByte('\r')
			i++
		case e == 'a':
			b.WriteByte(7)
			i++
		case e == 'b':
			b.WriteByte(8)
			i++
		case e == 'f':
			b.WriteByte(12)
			i++
		case e == 'v':
			b.WriteByte(11)
			i++
		case e == '\\':
			b.WriteByte('\\')
			i++
		case e == '"':
			b.WriteByte('"')
			i++
		case e == '\'':
			b.WriteByte('\'')
			i++
		case e == '\n':
			b.WriteByte('\n')
			i++
		case e == 'x':
			if i+2 >= len(body) {
				return "", false
			}
			v, err := strconv.ParseUint(body[i+1:i+3], 16, 8)
			if err != nil {
				return "", false
			}
			b.WriteByte(byte(v))
			i += 3
		case e >= '0' && e <= '9':
			j := i
			for j < len(body) && j < i+3 && body[j] >= '0' && body[j] <= '9' {
				j++
			}
			v, err := strconv.ParseUint(body[i:j], 10, 32)
			if err != nil || v > 255 {
				return "", false
			}
			b.WriteByte(byte(v))
			i = j
		default:
			return "", false // unknown escape (\z, \u{...}, etc.) — skip encryption
		}
	}
	return b.String(), true
}

// unquoteLongString decodes a [[...]] / [=*[...]=*] literal to its raw bytes.
func unquoteLongString(text string) (string, bool) {
	i := 1
	for i < len(text) && text[i] == '=' {
		i++
	}
	if i >= len(text) || text[i] != '[' {
		return "", false
	}
	level := i - 1
	open := i + 1
	closer := "]" + strings.Repeat("=", level) + "]"
	if !strings.HasSuffix(text, closer) || len(text) < open+len(closer) {
		return "", false
	}
	body := text[open : len(text)-len(closer)]
	// Lua drops a single newline immediately after the opening bracket.
	if strings.HasPrefix(body, "\r\n") {
		body = body[2:]
	} else if strings.HasPrefix(body, "\n") {
		body = body[1:]
	}
	return body, true
}

// encodeString returns a double-quoted Lua literal whose bytes are value[i]^key,
// each written as a 3-digit decimal escape (valid in Lua 5.1 and Luau).
func encodeString(value string, key byte) string {
	var b strings.Builder
	b.WriteByte('"')
	for i := 0; i < len(value); i++ {
		fmt.Fprintf(&b, "\\%03d", value[i]^key)
	}
	b.WriteByte('"')
	return b.String()
}

// DecoderPrelude returns the `local _d=...` definition that decodes strings
// produced by encodeString with the same key. It uses arithmetic XOR (no
// bit32) so it runs under Lua 5.1 as well as Roblox Luau.
func DecoderPrelude(key byte) string {
	return fmt.Sprintf("local _d=(function()"+
		"local function x(a,b)local r,p=0,1 for _=1,8 do local m,n=a%%2,b%%2 "+
		"if m~=n then r=r+p end a,b,p=(a-m)/2,(b-n)/2,p*2 end return r end "+
		"return function(s)local t={}for i=1,#s do t[i]=string.char(x(s:byte(i),%d))end "+
		"return table.concat(t)end end)()", key)
}
