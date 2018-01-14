package flag_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/terinjokes/bakelite/internal/flag"
)

func TestStringsValue(t *testing.T) {
	tests := []struct {
		v    string
		want []string
	}{
		{"", []string{}},
		{"a b c", []string{"a", "b", "c"}},
		{"foo bar baz", []string{"foo", "bar", "baz"}},
		{"foo  bar  baz", []string{"foo", "bar", "baz"}},
	}

	for i, tt := range tests {
		got := flag.StringsValue([]string{})
		got.Set(tt.v)

		if diff := cmp.Diff(tt.want, got.Get()); diff != "" {
			t.Errorf("#%d: manifest differs. (-got +want):\n%s", i, diff)
		}
	}
}
