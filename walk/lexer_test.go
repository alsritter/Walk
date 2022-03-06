package walk

import (
	"reflect"
	"strings"
	"testing"
)

var txt = "while i < 10 { \n" + // note: "//n" is a Token
	"  sum = sum + i \n" +
	"  i = i + 1 \n" +
	"}\n" +
	"sum"

func Test_lexer_Read(t *testing.T) {
	lexer := NewLexer(strings.NewReader(txt))

	tests := []struct {
		name string
		want Token
	}{
		{
			name: "id Token 01",
			want: NewIdToken(1, "while"),
		},
		{
			name: "id Token 02",
			want: NewIdToken(1, "i"),
		},
		{
			name: "id Token 03",
			want: NewIdToken(1, "<"),
		},
		{
			name: "number Token 01",
			want: NewNumToken(1, 10),
		},
		{
			name: "id Token 04",
			want: NewIdToken(1, "{"),
		},
		{
			name: "id Token 05",
			want: NewIdToken(1, "\\n"),
		},
		{
			name: "id Token 06",
			want: NewIdToken(2, "sum"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := lexer.Read(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("lexer.Read() = %v, want %v", got, tt.want)
			}
		})
	}
}
