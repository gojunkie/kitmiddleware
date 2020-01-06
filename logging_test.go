package kitmiddleware

import (
	"testing"
)

func TestDefaultRequestFormatter(t *testing.T) {
	suites := []struct {
		req  interface{}
		want []interface{}
	}{
		{
			struct {
				Name     string
				Age      int
				Email    string `val:"email"`
				Password string `val:"-"`
			}{
				Name:     "name2",
				Age:      22,
				Email:    "test@example.com",
				Password: "secret",
			},
			[]interface{}{"req", `[Name:name2 Age:22 email:test@example.com]`},
		},
	}

	valuer := NewDefaultRequestValuer()
	for _, suite := range suites {
		got := valuer.KeyValues(suite.req)
		if got[1] != suite.want[1] {
			t.Errorf("Format(%v): got %q, want %q", suite.req, got, suite.want)
		}
	}
}
