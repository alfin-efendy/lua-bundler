package lua

import (
	"crypto/rand"
	"math/big"
)

// Rename resolves scopes and assigns obfuscated names to all local bindings.
// It is a no-op if the chunk uses string interpolation, because identifiers
// inside `{...}` are not tracked in Phase 1 and renaming could break them.
func Rename(c *Chunk) {
	if resolve(c) { // returns hasInterp
		return
	}
	seen := map[*Binding]bool{}
	// Collect every binding by walking declarations; assign a fresh name once.
	walkBindings(c.Body, func(b *Binding) {
		if b == nil || seen[b] || b.OrigName == "self" {
			return
		}
		seen[b] = true
		b.NewName = generateName()
	})
}

// generateName returns an identifier like _0x1a2b3c.
func generateName() string {
	const chars = "0123456789abcdef"
	out := make([]byte, 6)
	for i := range out {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
		out[i] = chars[n.Int64()]
	}
	return "_0x" + string(out)
}

// walkBindings visits each declaration NameExpr's binding exactly where it is
// introduced. Because resolve() already shares one *Binding across a decl and
// its uses, assigning NewName on the declaration updates every use.
func walkBindings(stats []Stat, fn func(*Binding)) {
	for _, st := range stats {
		bindingsInStat(st, fn)
	}
}

func bindingsInStat(st Stat, fn func(*Binding)) {
	switch v := st.(type) {
	case *LocalStat:
		for _, n := range v.Names {
			fn(n.Binding)
		}
		for _, e := range v.Values {
			bindingsInExpr(e, fn)
		}
	case *LocalFuncStat:
		fn(v.Name.Binding)
		bindingsInFunc(v.Func, fn)
	case *FuncStat:
		bindingsInExpr(v.Target, fn)
		bindingsInFunc(v.Func, fn)
	case *AssignStat:
		for _, t := range v.Targets {
			bindingsInExpr(t, fn)
		}
		for _, e := range v.Values {
			bindingsInExpr(e, fn)
		}
	case *CallStat:
		bindingsInExpr(v.Call, fn)
	case *DoStat:
		walkBindings(v.Body, fn)
	case *WhileStat:
		bindingsInExpr(v.Cond, fn)
		walkBindings(v.Body, fn)
	case *RepeatStat:
		walkBindings(v.Body, fn)
		bindingsInExpr(v.Cond, fn)
	case *IfStat:
		for i := range v.Conds {
			bindingsInExpr(v.Conds[i], fn)
			walkBindings(v.Blocks[i], fn)
		}
		if v.HasElse {
			walkBindings(v.Else, fn)
		}
	case *NumericForStat:
		fn(v.Var.Binding)
		bindingsInExpr(v.Start, fn)
		bindingsInExpr(v.Stop, fn)
		if v.Step != nil {
			bindingsInExpr(v.Step, fn)
		}
		walkBindings(v.Body, fn)
	case *GenericForStat:
		for _, n := range v.Vars {
			fn(n.Binding)
		}
		for _, e := range v.Exprs {
			bindingsInExpr(e, fn)
		}
		walkBindings(v.Body, fn)
	case *ReturnStat:
		for _, e := range v.Values {
			bindingsInExpr(e, fn)
		}
	}
}

func bindingsInFunc(f *FuncExpr, fn func(*Binding)) {
	for _, p := range f.Params {
		fn(p.Binding)
	}
	walkBindings(f.Body, fn)
}

func bindingsInExpr(x Expr, fn func(*Binding)) {
	switch v := x.(type) {
	case *FuncExpr:
		bindingsInFunc(v, fn)
	case *IndexExpr:
		bindingsInExpr(v.Obj, fn)
		if v.Key != nil {
			bindingsInExpr(v.Key, fn)
		}
	case *CallExpr:
		bindingsInExpr(v.Fn, fn)
		for _, a := range v.Args {
			bindingsInExpr(a, fn)
		}
	case *TableExpr:
		for _, f := range v.Fields {
			if f.Key != nil {
				bindingsInExpr(f.Key, fn)
			}
			bindingsInExpr(f.Value, fn)
		}
	case *BinExpr:
		bindingsInExpr(v.L, fn)
		bindingsInExpr(v.R, fn)
	case *UnExpr:
		bindingsInExpr(v.E, fn)
	case *ParenExpr:
		bindingsInExpr(v.E, fn)
	case *IfExpr:
		bindingsInExpr(v.Cond, fn)
		bindingsInExpr(v.Then, fn)
		for i := range v.ElifConds {
			bindingsInExpr(v.ElifConds[i], fn)
			bindingsInExpr(v.ElifThen[i], fn)
		}
		bindingsInExpr(v.Else, fn)
	}
	// NameExpr here is a *use*; its binding (if any) already got a NewName at the
	// declaration site, so nothing to do. Leaves: nothing.
}
