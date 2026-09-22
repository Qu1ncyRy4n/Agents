package workspace

import (
	"reflect"
	"testing"
)

func TestManifestHeadingPathEscaping(t *testing.T) {
	tests := []struct {
		input string
		want  []string
	}{
		{input: "Workflow/Plan", want: []string{"Workflow", "Plan"}},
		{input: `Constraints \/ Safety`, want: []string{"Constraints / Safety"}},
		{input: `Parent\\Name/Child`, want: []string{"Parent\\Name", "Child"}},
		{input: "// Parent // Child /", want: []string{"Parent", "Child"}},
	}
	for _, test := range tests {
		got, err := ParseManifestHeadingPath(test.input)
		if err != nil {
			t.Fatalf("ParseManifestHeadingPath(%q): %v", test.input, err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("ParseManifestHeadingPath(%q) = %#v, want %#v", test.input, got, test.want)
		}
		if formatted := FormatManifestHeadingPath(got); formatted != FormatManifestHeadingPath(test.want) {
			t.Errorf("FormatManifestHeadingPath(%#v) = %q", got, formatted)
		}
	}
}

func TestManifestHeadingPathRejectsInvalidEscapes(t *testing.T) {
	for _, input := range []string{`Parent\x`, `Parent\`} {
		if _, err := ParseManifestHeadingPath(input); err == nil {
			t.Errorf("ParseManifestHeadingPath(%q) succeeded", input)
		}
	}
}
