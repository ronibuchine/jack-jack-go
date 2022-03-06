package vmtranslator

import (
	"reflect"
	"strings"
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
			want:    &Command{cmdType: C_POP, arg1: "local", arg2: "5"},
			wantErr: false,
		},
		{name: "add",
			args:    args{s: "   add  \t"},
			want:    &Command{cmdType: C_ARITHMETIC, arg1: "add", arg2: ""},
			wantErr: false,
		},
		{name: "poopy butt",
			args:    args{s: "poopy butt\t"},
			want:    nil,
			wantErr: true,
		},
		{name: "return // poopybutt",
			args:    args{s: "return // poopybutt"},
			want:    &Command{cmdType: C_RETURN, arg1: "", arg2: ""},
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
			got, err := parseCommand(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseCommand() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCommand() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseFile(t *testing.T) {
	tests := []struct {
		name         string
		arg          string
		wantCommands []*Command
	}{
		{name: "push and pop",
			arg:          "push local 5\npop local 9",
			wantCommands: []*Command{{cmdType: C_PUSH, arg1: "local", arg2: "5"}, {cmdType: C_POP, arg1: "local", arg2: "9"}},
		},
		{name: "push, push, add",
			arg:          "push constant 5//push 5\npop constant 9//push 9\nadd // now add them together\n",
            wantCommands: []*Command{{cmdType: C_PUSH, arg1: "constant", arg2: "5"}, {cmdType: C_POP, arg1: "constant", arg2: "9"}, {cmdType: C_ARITHMETIC, arg1: "add", arg2: ""}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCommands := ParseFile(strings.NewReader(tt.arg)); !reflect.DeepEqual(gotCommands, tt.wantCommands) {
				t.Errorf("ParseFile() = %v, want %v", gotCommands, tt.wantCommands)
			}
		})
	}
}
