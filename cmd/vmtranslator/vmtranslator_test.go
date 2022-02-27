package main

import "testing"

func Test_normalizeSpaces(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name               string
		args               args
		wantNormalizedLine string
	}{
		{name: "Normal String",
			args:               args{line: "bundle o' joy"},
			wantNormalizedLine: "bundle o' joy"},

		{name: "lots of spaces String",
			args:               args{line: " bundle    o' joy   "},
			wantNormalizedLine: "bundle o' joy"},

		{name: "trailing spaces String",
			args:               args{line: "bundle o' joy      "},
			wantNormalizedLine: "bundle o' joy"},

		{name: "spaces everywhere String",
			args:               args{line: "    bundle     o'   joy  "},
			wantNormalizedLine: "bundle o' joy"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotNormalizedLine := normalizeSpaces(tt.args.line); gotNormalizedLine != tt.wantNormalizedLine {
				t.Errorf("normalizeSpaces() = %v, want %v", gotNormalizedLine, tt.wantNormalizedLine)
			}
		})
	}
}
