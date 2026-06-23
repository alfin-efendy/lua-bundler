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

// EncryptStrings replaces eligible string literals in the chunk with _d("...")
// decoder calls (bytes XOR-ed with key). It never encrypts: backtick
// interpolation strings, the string argument of require(...), or the string
// argument of game:HttpGet(...) — those must survive as literals for module
// resolution. Literals it cannot safely decode are left unchanged.
func EncryptStrings(c *Chunk, key byte) {
	e := &encryptor{key: key}
	e.block(c.Body)
}

type encryptor struct{ key byte }

func (e *encryptor) block(stats []Stat) {
	for _, s := range stats {
		e.stat(s)
	}
}

func (e *encryptor) stat(s Stat) {
	switch v := s.(type) {
	case *LocalStat:
		e.exprs(v.Values)
	case *AssignStat:
		e.exprs(v.Targets)
		e.exprs(v.Values)
	case *CallStat:
		v.Call = e.expr(v.Call)
	case *DoStat:
		e.block(v.Body)
	case *WhileStat:
		v.Cond = e.expr(v.Cond)
		e.block(v.Body)
	case *RepeatStat:
		e.block(v.Body)
		v.Cond = e.expr(v.Cond)
	case *IfStat:
		for i := range v.Conds {
			v.Conds[i] = e.expr(v.Conds[i])
			e.block(v.Blocks[i])
		}
		if v.HasElse {
			e.block(v.Else)
		}
	case *NumericForStat:
		v.Start = e.expr(v.Start)
		v.Stop = e.expr(v.Stop)
		if v.Step != nil {
			v.Step = e.expr(v.Step)
		}
		e.block(v.Body)
	case *GenericForStat:
		e.exprs(v.Exprs)
		e.block(v.Body)
	case *FuncStat:
		e.block(v.Func.Body)
	case *LocalFuncStat:
		e.block(v.Func.Body)
	case *ReturnStat:
		e.exprs(v.Values)
	}
}

func (e *encryptor) exprs(xs []Expr) {
	for i := range xs {
		xs[i] = e.expr(xs[i])
	}
}

// expr returns x or its encrypted replacement, recursing into children.
func (e *encryptor) expr(x Expr) Expr {
	switch v := x.(type) {
	case *StringExpr:
		if repl, ok := e.encode(v); ok {
			return repl
		}
		return v
	case *ParenExpr:
		v.E = e.expr(v.E)
		return v
	case *BinExpr:
		v.L = e.expr(v.L)
		v.R = e.expr(v.R)
		return v
	case *UnExpr:
		v.E = e.expr(v.E)
		return v
	case *IndexExpr:
		v.Obj = e.expr(v.Obj)
		if v.Key != nil {
			v.Key = e.expr(v.Key)
		}
		return v
	case *CallExpr:
		v.Fn = e.expr(v.Fn)
		protected := isProtectedCall(v)
		for i, a := range v.Args {
			if protected {
				if _, isStr := a.(*StringExpr); isStr {
					continue // never encrypt require/HttpGet string args
				}
			}
			v.Args[i] = e.expr(a)
		}
		return v
	case *TableExpr:
		for i := range v.Fields {
			if v.Fields[i].Key != nil {
				v.Fields[i].Key = e.expr(v.Fields[i].Key)
			}
			v.Fields[i].Value = e.expr(v.Fields[i].Value)
		}
		return v
	case *FuncExpr:
		e.block(v.Body)
		return v
	case *IfExpr:
		v.Cond = e.expr(v.Cond)
		v.Then = e.expr(v.Then)
		for i := range v.ElifConds {
			v.ElifConds[i] = e.expr(v.ElifConds[i])
			v.ElifThen[i] = e.expr(v.ElifThen[i])
		}
		v.Else = e.expr(v.Else)
		return v
	default:
		return x // NameExpr, NumberExpr, BoolExpr, NilExpr, VarargExpr
	}
}

// encode turns a string literal into a _d("...") call, or reports ok=false to
// leave it unchanged (backtick interpolation or undecodable literal).
func (e *encryptor) encode(s *StringExpr) (Expr, bool) {
	if len(s.Text) > 0 && s.Text[0] == '`' {
		return nil, false
	}
	val, ok := unquoteLuaString(s.Text)
	if !ok {
		return nil, false
	}
	return &CallExpr{
		Fn:   &NameExpr{Name: "_d"},
		Args: []Expr{&StringExpr{Text: encodeString(val, e.key)}},
	}, true
}

// isProtectedCall reports whether v is a require(...) or X:HttpGet(...) call
// whose string argument must not be encrypted.
func isProtectedCall(v *CallExpr) bool {
	switch fn := v.Fn.(type) {
	case *NameExpr:
		return fn.Name == "require"
	case *IndexExpr:
		return fn.IsMethod && fn.Field == "HttpGet"
	}
	return false
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
