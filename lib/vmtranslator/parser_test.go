package vmtranslator

import (
	"reflect"
	"testing"
)

func Test_parseCommand(t *testing.T) {
	CompileAllRegex()
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    *Command
		wantErr bool
	}{
		{name: "pop",
			args:    args{s: "pop local 5"},
			want:    &Command{cmdType: CPop, arg1: "local", arg2: "5"},
			wantErr: false,
		},
		{name: "add",
			args:    args{s: "   add  \t"},
			want:    &Command{cmdType: CArithmetic, arg1: "add", arg2: ""},
			wantErr: false,
		},
		{name: "poopy butt",
			args:    args{s: "poopy butt\t"},
			want:    nil,
			wantErr: true,
		},
		{name: "return // poopybutt",
			args:    args{s: "return // poopybutt"},
			want:    &Command{cmdType: CReturn, arg1: "", arg2: ""},
			wantErr: false,
		},
		{name: "// comment",
			args:    args{s: "// comment"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCommand(tt.args.s, "VMFILE")
			if (got == nil) != tt.wantErr {
				t.Errorf("parseCommand() error = %v, wantErr %v", got == nil, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}
