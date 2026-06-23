package lua

import (
	"strings"
	"testing"
)

func parseReturnExpr(t *testing.T, src string) Expr {
	t.Helper()
	c, err := Parse("return " + src)
	if err != nil {
		t.Fatalf("parse %q: %v", src, err)
	}
	if len(c.Body) != 1 {
		t.Fatalf("want 1 stat, got %d", len(c.Body))
	}
	r, ok := c.Body[0].(*ReturnStat)
	if !ok || len(r.Values) != 1 {
		t.Fatalf("want a return with one value, got %T", c.Body[0])
	}
	return r.Values[0]
}

func TestParseExpr_Precedence(t *testing.T) {
	e := parseReturnExpr(t, "1 + 2 * 3")
	b, ok := e.(*BinExpr)
	if !ok || b.Op != "+" {
		t.Fatalf("want top-level +, got %#v", e)
	}
	if _, ok := b.R.(*BinExpr); !ok {
		t.Fatalf("want * grouped on the right, got %#v", b.R)
	}
}

func TestParseExpr_CallAndIndex(t *testing.T) {
	e := parseReturnExpr(t, "a.b.c(1)")
	call, ok := e.(*CallExpr)
	if !ok {
		t.Fatalf("want CallExpr, got %T", e)
	}
	if _, ok := call.Fn.(*IndexExpr); !ok {
		t.Fatalf("want IndexExpr callee, got %T", call.Fn)
	}
}

func TestParseExpr_MethodCall(t *testing.T) {
	e := parseReturnExpr(t, "obj:Method(1, 2)")
	call := e.(*CallExpr)
	idx, ok := call.Fn.(*IndexExpr)
	if !ok || !idx.IsMethod || idx.Field != "Method" {
		t.Fatalf("want method index, got %#v", call.Fn)
	}
}

func TestParseExpr_Table(t *testing.T) {
	e := parseReturnExpr(t, `{x = 1, [y] = 2, 3}`)
	tbl, ok := e.(*TableExpr)
	if !ok || len(tbl.Fields) != 3 {
		t.Fatalf("want 3 fields, got %#v", e)
	}
	if tbl.Fields[0].KeyName != "x" {
		t.Fatalf("want name key x, got %q", tbl.Fields[0].KeyName)
	}
	if tbl.Fields[1].Key == nil {
		t.Fatalf("want computed key for field 2")
	}
}

func TestParseStat_RoundTrip(t *testing.T) {
	cases := []string{
		"local x = 1",
		"local a, b = 1, 2",
		"x = 1",
		"a.b = c",
		"a, b = b, a",
		"foo()",
		"obj:m(1)",
		"do local x = 1 end",
		"while a do b() end",
		"repeat b() until a",
		"if a then b() elseif c then d() else e() end",
		"for i = 1, 10 do f(i) end",
		"for i = 1, 10, 2 do f(i) end",
		"for k, v in pairs(t) do f(k, v) end",
		"function a.b.c() return 1 end",
		"function a:m() return self end",
		"local function f(x) return x end",
		"return",
		"return 1, 2",
		"break",
		"local t = {1, 2, 3}",
	}
	for _, in := range cases {
		c, err := Parse(in)
		if err != nil {
			t.Errorf("parse %q: %v", in, err)
			continue
		}
		if len(c.Body) == 0 {
			t.Errorf("parse %q: empty body", in)
		}
		_ = c.Print() // must not panic
	}
}

func TestParseLuau_Accepts(t *testing.T) {
	cases := []string{
		"local x: number = 1",
		"local t: {string} = {}",
		"function f(a: number, b: string): boolean return true end",
		"type Point = { x: number, y: number }",
		"export type Id = string",
		"local n = if a then 1 else 2",
		"local s = `hello {name}!`",
		"for i: number = 1, 10 do end",
		"continue",
	}
	for _, in := range cases {
		if _, err := Parse(in); err != nil {
			t.Errorf("Parse(%q) should accept Luau: %v", in, err)
		}
	}
}

func TestParseLuau_TypeAnnotationDropped(t *testing.T) {
	c, err := Parse("local x: number = 1\nreturn x")
	if err != nil {
		t.Fatal(err)
	}
	if len(c.Body) != 2 {
		t.Fatalf("want 2 stats (local, return), got %d: %#v", len(c.Body), c.Body)
	}
	ls, ok := c.Body[0].(*LocalStat)
	if !ok || len(ls.Names) != 1 || ls.Names[0].Name != "x" || len(ls.Values) != 1 {
		t.Fatalf("type annotation mis-parsed: %#v", c.Body[0])
	}
	if got := c.Print(); got != "local x=1 return x" {
		t.Fatalf("type not dropped on print: %q", got)
	}
}

func TestParseLuau_FuncTypesDropped(t *testing.T) {
	c, err := Parse("local function f(a: number, b: string): boolean return a end\nreturn f")
	if err != nil {
		t.Fatal(err)
	}
	got := c.Print()
	if strings.Contains(got, "number") || strings.Contains(got, "boolean") {
		t.Fatalf("param/return types not dropped: %q", got)
	}
}

func TestParseLuau_TypeAliasDropped(t *testing.T) {
	c, err := Parse("type Point = { x: number, y: number }\nlocal p = 1\nreturn p")
	if err != nil {
		t.Fatal(err)
	}
	got := c.Print()
	if strings.Contains(got, "Point") || strings.Contains(got, "number") {
		t.Fatalf("type alias not dropped: %q", got)
	}
}
