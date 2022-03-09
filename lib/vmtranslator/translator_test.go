package vmtranslator

import "testing"

func Test_vmArgumentAddressToAD(t *testing.T) {
	type args struct {
		command *Command
	}
	tests := []struct {
		name     string
		args     args
		wantHack string
	}{
		{
			name: "test local",
			args: args{&Command{cmdType: CPop,
				arg1: "local",
				arg2: "5"}},
			wantHack: "@LCL\nAD=M\n@5\nAD=D+A\n",
		},
		{
			name: "test argument",
			args: args{&Command{cmdType: CPop,
				arg1: "argument",
				arg2: "5"}},
			wantHack: "@ARG\nAD=M\n@5\nAD=D+A\n",
		},
		{
			name: "test squish",
			args: args{&Command{cmdType: CPop,
				arg1: "temp",
				arg2: "0"}},
			wantHack: "@5\nD=A\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotHack := vmArgumentAddressToAD(tt.args.command); gotHack != tt.wantHack {
				t.Errorf("vmArgumentAddressToAD() = %v, want %v", gotHack, tt.wantHack)
			}
		})
	}
}
