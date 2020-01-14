package kitmiddleware

import (
	"testing"
)

func TestDefaultRequestFormatter(t *testing.T) {
	type user struct {
		Name     string
		Age      int
		Email    string `val:"email"`
		Password string `val:"-"`
	}

	suites := []struct {
		req  interface{}
		want []interface{}
	}{
		{
			user{
				Name:     "name2",
				Age:      22,
				Email:    "test@example.com",
				Password: "secret",
			},
			[]interface{}{"req", `[Name:name2 Age:22 email:test@example.com]`},
		},
		{
			&user{
				Name:  "name1",
				Email: "test1@example.com",
			},
			[]interface{}{"req", `[Name:name1 Age:0 email:test1@example.com]`},
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
