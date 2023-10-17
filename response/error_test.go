package response

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_slugify(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "'abc' returns 'abc'",
			args: args{"abc"},
			want: "abc",
		},
		{
			name: "'aBc' returns 'abc'",
			args: args{"aBc"},
			want: "abc",
		},
		{
			name: "'aBc  DeF' returns 'abc__def'",
			args: args{"aBc  DeF"},
			want: "abc__def",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, slugify(tt.args.s), "slugify(%v)", tt.args.s)
		})
	}
}
