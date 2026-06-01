package version

import "testing"

func TestConstantIsValidSemver(t *testing.T) {
	if Constant == "" {
		t.Fatal("Constant is empty — set it to the next release version")
	}
	if v := parts(Constant); len(v) != 3 {
		t.Errorf("Constant = %q, want MAJOR.MINOR.PATCH, got %d segments", Constant, len(v))
	}
}

func TestCurrentFallsBackToConstant(t *testing.T) {
	got := Current()
	if got != Constant {
		t.Errorf("Current() = %q, want %q (Constant) — build info should not be set during go test", got, Constant)
	}
}

func TestGTE(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"1.3.0", "1.2.0", true},
		{"1.2.0", "1.3.0", false},
		{"1.3.0", "1.3.0", true},
		{"v1.3.0", "1.3.0", true},
		{"1.3.0-rc1", "1.3.0", true},
		{"", "0.0.0", true},
	}
	for _, c := range cases {
		if got := GTE(c.a, c.b); got != c.want {
			t.Errorf("GTE(%q, %q) = %v, want %v", c.a, c.b, got, c.want)
		}
	}
}
