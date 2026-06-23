package lua

// scope is a lexical scope mapping names to bindings.
type scope struct {
	parent *scope
	names  map[string]*Binding
}

func newScope(parent *scope) *scope { return &scope{parent: parent, names: map[string]*Binding{}} }

func (s *scope) lookup(name string) *Binding {
	for c := s; c != nil; c = c.parent {
		if b, ok := c.names[name]; ok {
			return b
		}
	}
	return nil
}

// declare binds name in this scope (shadowing any outer binding) and returns it.
func (s *scope) declare(name string) *Binding {
	b := &Binding{OrigName: name}
	s.names[name] = b
	return b
}

type resolver struct{ hasInterp bool }

// resolve attaches bindings to NameExprs. Returns whether the chunk uses string
// interpolation (renamer bails if so — see Task 8).
func resolve(c *Chunk) bool {
	r := &resolver{}
	top := newScope(nil)
	r.block(c.Body, top)
	return r.hasInterp
}

func (r *resolver) block(stats []Stat, s *scope) {
	for _, st := range stats {
		r.stat(st, s)
	}
}

func (r *resolver) stat(st Stat, s *scope) {
	switch v := st.(type) {
	case *LocalStat:
		// RHS resolves in the CURRENT scope (before the new names are visible).
		for _, val := range v.Values {
			r.expr(val, s)
		}
		for _, n := range v.Names {
			n.Binding = s.declare(n.Name)
		}
	case *LocalFuncStat:
		// Name is visible inside its own body (recursion): declare first.
		v.Name.Binding = s.declare(v.Name.Name)
		r.funcExpr(v.Func, s, false)
	case *FuncStat:
		// Target's leading name is a use; trailing fields are not names.
		r.expr(v.Target, s)
		r.funcExpr(v.Func, s, v.IsMethod)
	case *AssignStat:
		for _, t := range v.Targets {
			r.expr(t, s)
		}
		for _, val := range v.Values {
			r.expr(val, s)
		}
	case *CallStat:
		r.expr(v.Call, s)
	case *DoStat:
		r.block(v.Body, newScope(s))
	case *WhileStat:
		r.expr(v.Cond, s)
		r.block(v.Body, newScope(s))
	case *RepeatStat:
		// until-condition sees the body's locals: one shared scope.
		inner := newScope(s)
		r.block(v.Body, inner)
		r.expr(v.Cond, inner)
	case *IfStat:
		for i, cond := range v.Conds {
			r.expr(cond, s)
			r.block(v.Blocks[i], newScope(s))
		}
		if v.HasElse {
			r.block(v.Else, newScope(s))
		}
	case *NumericForStat:
		r.expr(v.Start, s)
		r.expr(v.Stop, s)
		if v.Step != nil {
			r.expr(v.Step, s)
		}
		inner := newScope(s)
		v.Var.Binding = inner.declare(v.Var.Name)
		r.block(v.Body, inner)
	case *GenericForStat:
		for _, e := range v.Exprs {
			r.expr(e, s)
		}
		inner := newScope(s)
		for _, n := range v.Vars {
			n.Binding = inner.declare(n.Name)
		}
		r.block(v.Body, inner)
	case *ReturnStat:
		for _, e := range v.Values {
			r.expr(e, s)
		}
	case *BreakStat, *ContinueStat, *GotoStat, *LabelStat, *TypeAliasStat:
		// nothing to resolve
	}
}

func (r *resolver) funcExpr(f *FuncExpr, parent *scope, isMethod bool) {
	inner := newScope(parent)
	if isMethod {
		inner.declare("self") // implicit self; not renamed (see Task 8 guard)
	}
	for _, pm := range f.Params {
		pm.Binding = inner.declare(pm.Name)
	}
	r.block(f.Body, inner)
}

func (r *resolver) expr(x Expr, s *scope) {
	switch v := x.(type) {
	case *NameExpr:
		v.Binding = s.lookup(v.Name) // nil => free/global
	case *IndexExpr:
		r.expr(v.Obj, s)
		if v.Key != nil {
			r.expr(v.Key, s)
		}
		// v.Field is a plain string: never resolved (field access).
	case *CallExpr:
		r.expr(v.Fn, s)
		for _, a := range v.Args {
			r.expr(a, s)
		}
	case *FuncExpr:
		r.funcExpr(v, s, false)
	case *TableExpr:
		for _, f := range v.Fields {
			if f.Key != nil {
				r.expr(f.Key, s)
			}
			// f.KeyName is a plain string: never resolved (table key).
			r.expr(f.Value, s)
		}
	case *BinExpr:
		r.expr(v.L, s)
		r.expr(v.R, s)
	case *UnExpr:
		r.expr(v.E, s)
	case *ParenExpr:
		r.expr(v.E, s)
	case *IfExpr:
		r.expr(v.Cond, s)
		r.expr(v.Then, s)
		for i := range v.ElifConds {
			r.expr(v.ElifConds[i], s)
			r.expr(v.ElifThen[i], s)
		}
		r.expr(v.Else, s)
	case *StringExpr:
		if len(v.Text) > 0 && v.Text[0] == '`' {
			r.hasInterp = true
		}
	case *NumberExpr, *BoolExpr, *NilExpr, *VarargExpr:
		// leaves
	}
}
