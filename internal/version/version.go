// Package version exposes the current themoviedb-cli release version and a
// small semver helper. The string is reported by `themoviedb-cli --version`
// and `themoviedb-cli version`.
package version

import (
	"runtime/debug"
	"strconv"
	"strings"
)

// Constant is the version baked into the binary at release time. It's the
// fallback used when runtime/debug doesn't know the module version (e.g.
// `go run .` or local builds without a tag). Bump on every release.
const Constant = "1.3.0"

// Current returns the active themoviedb-cli version, preferring the module
// version reported by runtime/debug (set by `go install github.com/yareeh/
// themoviedb-cli@vX.Y.Z`) and falling back to Constant.
func Current() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		v := strings.TrimPrefix(info.Main.Version, "v")
		if v != "" && v != "(devel)" {
			return v
		}
	}
	return Constant
}

// GTE reports whether a >= b under a tolerant semver comparison.
func GTE(a, b string) bool {
	return cmp(a, b) >= 0
}

func cmp(a, b string) int {
	aa := parts(a)
	bb := parts(b)
	n := len(aa)
	if len(bb) > n {
		n = len(bb)
	}
	for i := 0; i < n; i++ {
		var ai, bi int
		if i < len(aa) {
			ai = aa[i]
		}
		if i < len(bb) {
			bi = bb[i]
		}
		if ai != bi {
			if ai < bi {
				return -1
			}
			return 1
		}
	}
	return 0
}

func parts(v string) []int {
	v = strings.TrimPrefix(v, "v")
	if i := strings.IndexAny(v, "-+"); i >= 0 {
		v = v[:i]
	}
	raw := strings.Split(v, ".")
	out := make([]int, len(raw))
	for i, s := range raw {
		n, _ := strconv.Atoi(s)
		out[i] = n
	}
	return out
}
