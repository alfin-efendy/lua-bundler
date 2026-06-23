package lua

import "testing"

func TestMinify_PreservesStringSpaces(t *testing.T) {
	got := Minify(`local s = "a    b"`)
	if got != `local s="a    b"` {
		t.Fatalf("string spaces corrupted: %q", got)
	}
}

func TestMinify_KeepsNumberConcatSpace(t *testing.T) {
	got := Minify(`return 1 .. 2`)
	if got != `return 1 ..2` && got != `return 1 .. 2` {
		t.Fatalf("number/concat merged into malformed number: %q", got)
	}
}

func TestMinify_DropsComments(t *testing.T) {
	got := Minify("local x = 1 -- hi\nreturn x")
	if got != "local x=1 return x" {
		t.Fatalf("comment not dropped or spacing wrong: %q", got)
	}
}
