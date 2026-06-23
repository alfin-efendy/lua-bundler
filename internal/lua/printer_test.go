package lua

import "testing"

func roundTrip(t *testing.T, src string) string {
	t.Helper()
	c, err := Parse(src)
	if err != nil {
		t.Fatalf("parse %q: %v", src, err)
	}
	return c.Print()
}

func TestPrint_ExprRoundTrip(t *testing.T) {
	cases := map[string]string{
		"return 1 + 2 * 3":       "return 1+2*3",
		"return a.b.c(1)":        "return a.b.c(1)",
		"return obj:M(1, 2)":     "return obj:M(1,2)",
		`return {x = 1, [y]= 2}`: "return{x=1,[y]=2}",
		"return 1 .. 2":          "return 1 ..2",
		`return "a  b"`:          `return"a  b"`,
		"return not a and b":     "return not a and b",
	}
	for in, want := range cases {
		if got := roundTrip(t, in); got != want {
			t.Errorf("Print(%q) = %q, want %q", in, got, want)
		}
	}
}
