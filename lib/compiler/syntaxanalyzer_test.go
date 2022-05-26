package compiler

import (
	"reflect"
	"testing"
)

func TestReadStream(t *testing.T) {
	type args struct {
		tokenStream string
	}
	tests := []struct {
		name string
		args args
		want [][]string
	}{
		{
			name: "first",
			args: args{"sample.xml"},
			want: [][]string{{"if", "keyword"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ReadStream(tt.args.tokenStream)
			if got := NormalizedTokenStream; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadStream() = %v, want %v", got, tt.want)
			}
		})
	}
}
