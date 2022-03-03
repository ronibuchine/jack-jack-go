package vmtranslator

import (
	"os"
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
			want:    &Command{cmdType: C_POP, arg1: "local", arg2: "5"},
			wantErr: false,
		},
		{name: "add",
			args:    args{s: "   add  \t"},
			want:    &Command{cmdType: C_ARITHMETIC, arg1: "add", arg2: ""},
			wantErr: false,
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
	type args struct {
		file *os.File
	}
	tests := []struct {
		name         string
		args         args
		wantCommands []Command
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCommands := ParseFile(tt.args.file); !reflect.DeepEqual(gotCommands, tt.wantCommands) {
				t.Errorf("ParseFile() = %v, want %v", gotCommands, tt.wantCommands)
			}
		})
	}
}
