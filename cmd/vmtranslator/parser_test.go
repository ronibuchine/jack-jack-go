package vmtranslator

import (
	"testing"
)

func TestParseCommand(t *testing.T) {
	/* var (
		str   string
		cmd   Command
		wants Command
		err   error
	) */

	CompileAllRegex()

	tests := []struct {
		name string
		arg  string
		want Command_t
	}{
		{name: "pop",
			arg:  "pop local 5",
			want: Command_t{_cmdType: C_POP, _arg1: "local", _arg2: 5},
		},
		{name: "add",
			arg:  "   add  \t",
			want: Command_t{_cmdType: C_ARITHMETIC, _arg1: "add", _arg2: -1},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd, err := parseCommand(tt.arg)
			arg1, err := cmd.arg1()
			if err != nil {
				arg1 = ""
			}
			arg2, err := cmd.arg2()
			if err != nil {
				arg2 = -1
			}

			if tt.want._cmdType != cmd.cmdType() || tt.want._arg1 != arg1 || tt.want._arg2 != arg2 {
				t.Errorf("parseCommand() = %v, want %v", cmd, tt.want)
			}
		})
	}

	/* str = "pop local 5"
	cmd, err = parseCommand(str)
	wants = Command_t{_cmdType: C_POP, _arg1: "local", _arg2: 5}
	if cmd != wants || err != nil {
		t.Fatalf("expected %+v, but got %+v", wants, cmd)
	}

	str = "  add   \t"
	cmd, err = parseCommand(str)
	wants = Command_t{_cmdType: C_ARITHMETIC, _arg1: "add", _arg2: -1}
	if cmd != wants || err != nil {
		t.Fatalf("expected %+v, but got %+v", wants, cmd)
	} */
}
