package lua

import "testing"

func TestResolve_LocalBound(t *testing.T) {
	c, err := Parse("local x = 1\nreturn x")
	if err != nil {
		t.Fatal(err)
	}
	resolve(c)
	decl := c.Body[0].(*LocalStat).Names[0]
	use := c.Body[1].(*ReturnStat).Values[0].(*NameExpr)
	if decl.Binding == nil || use.Binding == nil {
		t.Fatal("expected both to be bound")
	}
	if decl.Binding != use.Binding {
		t.Fatal("use must share the declaration's binding")
	}
}

func TestResolve_GlobalFree(t *testing.T) {
	c, _ := Parse("return print")
	resolve(c)
	use := c.Body[0].(*ReturnStat).Values[0].(*NameExpr)
	if use.Binding != nil {
		t.Fatal("global must be free (nil binding)")
	}
}

func TestResolve_ShadowOrder(t *testing.T) {
	// RHS x refers to the OUTER x, not the one being declared.
	c, _ := Parse("local x = 1\nlocal x = x\nreturn x")
	resolve(c)
	outer := c.Body[0].(*LocalStat).Names[0].Binding
	rhs := c.Body[1].(*LocalStat).Values[0].(*NameExpr).Binding
	inner := c.Body[1].(*LocalStat).Names[0].Binding
	last := c.Body[2].(*ReturnStat).Values[0].(*NameExpr).Binding
	if rhs != outer {
		t.Fatal("RHS x must bind to outer declaration")
	}
	if last != inner {
		t.Fatal("final x must bind to the inner (shadowing) declaration")
	}
}

func TestResolve_FieldNotResolved(t *testing.T) {
	c, _ := Parse("local t = {}\nreturn t.x")
	resolve(c)
	idx := c.Body[1].(*ReturnStat).Values[0].(*IndexExpr)
	if obj := idx.Obj.(*NameExpr); obj.Binding == nil {
		t.Fatal("t should be bound")
	}
	// idx.Field is a string, never a NameExpr — nothing to resolve, by design.
}

func TestResolve_Interp(t *testing.T) {
	c, err := Parse("local n = 1\nreturn `x {n}`")
	if err != nil {
		t.Fatal(err)
	}
	if !resolve(c) {
		t.Fatal("expected hasInterp=true for a chunk with a backtick string")
	}
	c2, err := Parse("local n = 1\nreturn n")
	if err != nil {
		t.Fatal(err)
	}
	if resolve(c2) {
		t.Fatal("expected hasInterp=false for a plain chunk")
	}
}

func TestResolve_MethodSelf(t *testing.T) {
	c, err := Parse("function a:m() return self end")
	if err != nil {
		t.Fatal(err)
	}
	resolve(c)
	fs := c.Body[0].(*FuncStat)
	ret := fs.Func.Body[0].(*ReturnStat)
	self := ret.Values[0].(*NameExpr)
	if self.Binding == nil {
		t.Fatal("self must be bound inside a method body")
	}
}

func TestResolve_RepeatUntilSharedScope(t *testing.T) {
	c, err := Parse("repeat local x = 1 until x")
	if err != nil {
		t.Fatal(err)
	}
	resolve(c)
	rs := c.Body[0].(*RepeatStat)
	decl := rs.Body[0].(*LocalStat).Names[0]
	use, ok := rs.Cond.(*NameExpr)
	if !ok || use.Binding == nil || use.Binding != decl.Binding {
		t.Fatalf("until condition must resolve x to the body's local: %#v", rs.Cond)
	}
}
