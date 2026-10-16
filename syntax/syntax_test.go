package syntax_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/Qu1ncyRy4n/Agents/syntax"
)

func parse(t *testing.T, source string) *syntax.File {
	t.Helper()
	f, err := syntax.Parse("example.mogent", []byte(source))
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func field(t *testing.T, entry *syntax.Entry, name string) *syntax.Value {
	t.Helper()
	for _, child := range entry.Body {
		if child.Kind == syntax.Field && child.Name == name {
			return child.Value
		}
	}
	t.Fatalf("field %q missing from %s", name, entry.Name)
	return nil
}

func text(t *testing.T, value *syntax.Value, bindings map[string]string) string {
	t.Helper()
	got, err := value.ResolveString(bindings)
	if err != nil {
		t.Fatal(err)
	}
	return got
}

func TestParseLibraryPreservesTypedHierarchyAndMetadata(t *testing.T) {
	source := `// Team library, not a filesystem operation.
library { format = 1; name = "Team" }
md-section "guidance" {
  curate = "foundation"
  md-section "workflow" {
    source = "guidance/workflow.md"
    vars = ["repo_name"]
  }
}
dir "skills" {
  source = "skills"
  curate = "choose"
  node "review" {
    role = "skill"
    content = "SKILL.md"
    default = true
    tags = ["workflow", "quality",]
  }
  node "testing" { default = true }
  node "deployment" { default = false }
}
md-file "reference" { source = "reference.md" }
file "settings" { source = "settings.json" }
preset "recommended" {
  include "self:guidance"
  include "self:skills" { accept_defaults = true }
}
group "persona" { mode = "at-most-one" }
`
	f := parse(t, source)
	if f.Kind != syntax.Library || f.Original() != source || len(f.Entries) != 7 {
		t.Fatalf("library shape: kind=%s entries=%d", f.Kind, len(f.Entries))
	}
	if len(f.Comments) != 1 || !strings.Contains(f.Comments[0].Text, "not a filesystem") {
		t.Fatalf("comments lost: %#v", f.Comments)
	}
	skills := f.Entries[2]
	if skills.Name != "dir" || text(t, skills.Label, nil) != "skills" {
		t.Fatalf("physical root=%#v", skills)
	}
	review, testing, deployment := skills.Body[2], skills.Body[3], skills.Body[4]
	if text(t, review.Label, nil) != "review" || text(t, testing.Label, nil) != "testing" || text(t, deployment.Label, nil) != "deployment" {
		t.Fatal("library declaration order changed")
	}
	if !field(t, review, "default").Bool || field(t, deployment, "default").Bool {
		t.Fatal("boolean metadata changed")
	}
	tags := field(t, review, "tags")
	if tags.Kind != syntax.List || len(tags.Items) != 2 || text(t, &tags.Items[1], nil) != "quality" {
		t.Fatalf("tags=%#v", tags)
	}
}

func TestSeparatorsAndOptionalFinalSemicolon(t *testing.T) {
	for name, source := range map[string]string{
		"one-field":           `plan { format = 1 }`,
		"separated-fields":    `build { format = 1; origin = "authored" }`,
		"optional-trailing":   `build { format = 1; origin = "authored"; }`,
		"single-line-actions": `plan { format = 1 }; include "team:a"; include "team:b"`,
		"block-then-bare":     "plan { format = 1 }\ninclude \"team:a\" { into = \"AGENTS.md\" }\ninclude \"team:b\"",
		"empty-selection":     "plan { format = 1 }\ninclude \"team:skills\" []",
		"multiline-array":     "plan { format = 1 }\ninclude \"team:skills\" [\n\"review\", // chosen\n\"testing\",\n]",
	} {
		t.Run(name, func(t *testing.T) { parse(t, source) })
	}
}

func TestParsePlanAndBuildRetainActionOrderAndLiteralOptions(t *testing.T) {
	f := parse(t, `plan { format = 1 }
let repo_name = "Example"
source "team" { local = "../library" }
include "team:skills" ["review", "testing"] {
  exclude = ["deployment"]
}
replace "local:skills/review"
include { text = "# ${repo_name}"; into = "AGENTS.md" }
include-build "team:builds/common.mogent.build"
`)
	if f.Kind != syntax.Plan || len(f.Entries) != 7 {
		t.Fatalf("kind=%s entries=%d", f.Kind, len(f.Entries))
	}
	if f.Entries[1].Kind != syntax.Binding || f.Entries[1].Name != "repo_name" {
		t.Fatal("binding not represented distinctly")
	}
	selection := f.Entries[3].Selection
	if len(selection.Items) != 2 || text(t, &selection.Items[0], nil) != "review" {
		t.Fatal("child-list order lost")
	}
	if f.Entries[4].Name != "replace" || f.Entries[4].HasBody {
		t.Fatal("bare action lost")
	}
	if got := text(t, field(t, f.Entries[5], "text"), map[string]string{"repo_name": "Example"}); got != "# Example" {
		t.Fatalf("text=%q", got)
	}
	b := parse(t, `build { format = 1; origin = "generated"; generated_from = "mogent.plan" }
resolved {
  input "local:skills/review" { sha256 = "example" }
  decision "team:skills/deployment" { action = "exclude"; reason = "consumer" }
  owner ".agents/skills" { kind = "directory" }
}
include "team:guidance/workflow" {
  values { repo_name = "Example" }
}
`)
	if b.Kind != syntax.Build || b.Entries[1].Name != "resolved" {
		t.Fatal("resolved metadata lost")
	}
}

func TestQuotedStringsAndInterpolationAreExplicit(t *testing.T) {
	f := parse(t, "plan {format=1}\ninclude {text=\"Hi ${name}; \\${PATH}; \\\"quoted\\\"; \\\\; \\n\\t\"}")
	v := field(t, f.Entries[1], "text")
	if got := text(t, v, map[string]string{"name": "${do_not_evaluate}"}); got != "Hi ${do_not_evaluate}; ${PATH}; \"quoted\"; \\; \n\t" {
		t.Fatalf("substitution/escaping=%q", got)
	}
	_, err := v.ResolveString(nil)
	var diagnostic *syntax.Diagnostic
	if !errors.As(err, &diagnostic) || !strings.Contains(err.Error(), `undefined binding "name"`) {
		t.Fatalf("missing binding diagnostic=%v", err)
	}
	if _, err := (syntax.Value{Kind: syntax.Integer}).ResolveString(nil); err == nil {
		t.Fatal("integer coerced to string")
	}
}

func TestMultilineTextDedentsAndLiteralTextDoesNotInterpolate(t *testing.T) {
	source := "plan { format = 1 }\ninclude {\n  text = \"\"\"\n    # ${repo}\n\n      indented\n    \"\"\"\n}\ninclude {\n  text = '''\n    ${SHELL} \\n \\path \"\"\"\n    '''\n}\n"
	for _, source := range []string{source, strings.ReplaceAll(source, "\n", "\r\n"), strings.ReplaceAll(source, "\n", "\r")} {
		f := parse(t, source)
		if got := text(t, field(t, f.Entries[1], "text"), map[string]string{"repo": "Example"}); got != "# Example\n\n  indented" {
			t.Fatalf("multiline=%q", got)
		}
		raw := field(t, f.Entries[2], "text")
		if !raw.Literal || text(t, raw, nil) != "${SHELL} \\n \\path \"\"\"" {
			t.Fatalf("literal=%#v", raw)
		}
	}
	for _, delimiter := range []string{`"""`, `'''`} {
		f := parse(t, "plan {format=1}\ninclude {text="+delimiter+"\n"+delimiter+"}")
		if got := text(t, field(t, f.Entries[1], "text"), nil); got != "" {
			t.Fatalf("empty triple=%q", got)
		}
	}
}

func TestDiagnosticsLocateOriginalBytesAndUnicodeColumns(t *testing.T) {
	source := "plan { format = 1 }\r\ninclude { text = \"é\"; into = \"${missing}\" }\r\n"
	f := parse(t, source)
	v := field(t, f.Entries[1], "into")
	_, err := v.ResolveString(nil)
	var diagnostic *syntax.Diagnostic
	if !errors.As(err, &diagnostic) {
		t.Fatalf("err=%v", err)
	}
	r := diagnostic.Range
	start := strings.Index(source, `"${missing}"`)
	lineStart := strings.Index(source, "\n") + 1
	if r.Filename != "example.mogent" || r.Start.Offset != start || r.Start.Line != 2 || r.Start.Column != utf8.RuneCountInString(source[lineStart:start])+1 {
		t.Fatalf("location=%#v", r)
	}
	if got := source[r.Start.Offset:r.End.Offset]; got != `"${missing}"` {
		t.Fatalf("source range=%q", got)
	}
	if got := source[f.Entries[1].Range.Start.Offset:f.Entries[1].Range.End.Offset]; !strings.HasPrefix(got, "include {") || !strings.HasSuffix(got, "}") {
		t.Fatalf("entry range=%q", got)
	}
}

func TestInvalidSyntaxIsRejectedWithoutExpressionEvaluation(t *testing.T) {
	cases := []struct{ name, source, want string }{
		{"empty", "// empty", "require a library"},
		{"missing-header", `include "team:a"`, "first declaration"},
		{"version", `plan {format=2}`, "format = 1"},
		{"missing-origin", `build {format=1}`, "requires origin"},
		{"bad-origin", `build {format=1; origin="mixed"}`, "origin must be"},
		{"field-comma", `build {format=1, origin="authored"}`, "commas only"},
		{"missing-semicolon", `build {format=1 origin="authored"}`, "between statements"},
		{"array-comma", "plan {format=1}\ninclude \"x\" [\"a\"\n\"b\"]", "between array items"},
		{"array-number", "plan {format=1}\ninclude \"x\" [1]", "quoted strings"},
		{"array-semicolon", "plan {format=1}\ninclude \"x\" [\"a\";\"b\"]", "between array items"},
		{"object", "plan {format=1}\ninclude {text={key=\"value\"}}", "object assignments"},
		{"function", "plan {format=1}\ninclude {text=env(\"HOME\")}", "expressions"},
		{"expression", "plan {format=1}\ninclude {text=\"${name + other}\"}", "not an expression"},
		{"unfinished-reference", "plan {format=1}\ninclude {text=\"${name\"}", "unterminated interpolation"},
		{"unknown-block", "plan {format=1}\nexec \"script\"", "unexpected"},
		{"unknown-option", "plan {format=1}\ninclude \"x\" {mode=\"copy\"}", "unknown field"},
		{"duplicate-option", "plan {format=1}\ninclude \"x\" {into=\"a\";into=\"b\"}", "duplicate field"},
		{"duplicate-binding", "plan {format=1}\nlet x=\"a\"\nlet x=\"b\"", "duplicate binding"},
		{"binding-type", "plan {format=1}\nlet x=1", "must be strings"},
		{"duplicate-source", "plan {format=1}\nsource \"x\" {local=\"a\"}\nsource \"x\" {local=\"b\"}", "duplicate"},
		{"duplicate-sibling", "library {format=1}\nmd-section \"g\" {node \"x\" {}\nnode \"x\" {}}", "duplicate sibling"},
		{"duplicate-public-name", "library {format=1}\ndir \"x\" {}\npreset \"x\" {}", "duplicate"},
		{"alias-conflict", "library {format=1}\ndir \"x\" {type=\"file\"}", "contradicts"},
		{"replace-false", "plan {format=1}\nreplace \"x\" {replace=false}", "cannot disable"},
		{"replace-inline", "plan {format=1}\nreplace \"x\" {text=\"x\"}", "cannot combine"},
		{"unnamed-include", "plan {format=1}\ninclude {}", "source reference or inline text"},
		{"duplicate-values", "plan {format=1}\ninclude \"x\" {values {}\nvalues {}}", "only one values"},
		{"values-type", "plan {format=1}\ninclude \"x\" {values {repo=1}}", "values must be strings"},
		{"recipe-options", "plan {format=1}\ninclude-build \"x\" {}", "does not take"},
		{"reserved-records", "plan {format=1}\nresolved {}", "belong in a build"},
		{"missing-brace", "plan {format=1}\ninclude \"x\" {", "closing brace"},
		{"extra-brace", "plan {format=1}\n}", "unexpected closing"},
		{"inline-triple", "plan {format=1}\ninclude {text=\"\"\"inline\"\"\"}", "begin with a newline"},
		{"closing-triple", "plan {format=1}\ninclude {text=\"\"\"\nx\"\"\"}", "own line"},
		{"indent", "plan {format=1}\ninclude {text=\"\"\"\nx\n  \"\"\"}", "less indentation"},
		{"unterminated", "plan {format=1}\ninclude {text=\"x", "unterminated string"},
		{"escape", "plan {format=1}\ninclude {text=\"\\q\"}", "unsupported string escape"},
		{"literal-single", "plan {format=1}\ninclude {text='x'}", "triple single"},
		{"overflow", "plan {format=9999999999999999999999999999999}", "64-bit range"},
		{"utf8", "plan {format=1}\n\xff", "UTF-8"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := syntax.Parse("bad.mogent", []byte(tc.source))
			var diagnostic *syntax.Diagnostic
			if !errors.As(err, &diagnostic) || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error=%v; want located %q", err, tc.want)
			}
		})
	}
}

func TestCanonicalFormattingExpandsAliasesWithoutMutatingInput(t *testing.T) {
	source := "// header\nlibrary {format=1;name=\"Team\"}\n// skills\ndir \"skills\" {source=\"skills\"}\npreset \"recommended\" {replace \"self:skills/review\"}\n"
	f := parse(t, source)
	out, err := syntax.Format(f, syntax.FormatOptions{ExpandAliases: true})
	if err != nil {
		t.Fatal(err)
	}
	want := "// header\nlibrary {\n  format = 1\n  name = \"Team\"\n}\n// skills\nnode \"skills\" {\n  type = \"dir\"\n  source = \"skills\"\n}\npreset \"recommended\" {\n  include \"self:skills/review\" {\n    replace = true\n  }\n}\n"
	if string(out) != want {
		t.Fatalf("canonical output:\n%s\nwant:\n%s", out, want)
	}
	if f.Original() != source || f.Entries[1].Name != "dir" || f.Entries[2].Body[0].Name != "replace" {
		t.Fatal("format changed authored syntax")
	}
	second, err := syntax.Format(parse(t, string(out)), syntax.FormatOptions{ExpandAliases: true})
	if err != nil || !bytes.Equal(out, second) {
		t.Fatalf("format not stable: %v\n%s", err, second)
	}
	aliases, err := syntax.Format(f, syntax.FormatOptions{})
	if err != nil || !strings.Contains(string(aliases), `dir "skills"`) || !strings.Contains(string(aliases), `replace "self:skills/review"`) {
		t.Fatalf("alias spelling lost: %v\n%s", err, aliases)
	}
}

func TestFormattingPreservesCommentsAndLiteralInterpolationSemantics(t *testing.T) {
	source := "plan {format=1} // after header\ninclude \"team:skills\" [\"review\", // inside list\n\"testing\"]\ninclude { // within block\ntext='''\n${literal} \\path \"\"\"\n'''\n}\n// end\n"
	f := parse(t, source)
	out, err := syntax.Format(f, syntax.FormatOptions{})
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range f.Comments {
		if strings.Count(string(out), c.Text) != 1 {
			t.Fatalf("comment lost/duplicated %q:\n%s", c.Text, out)
		}
	}
	reparsed := parse(t, string(out))
	if got := text(t, field(t, reparsed.Entries[2], "text"), nil); got != "${literal} \\path \"\"\"" {
		t.Fatalf("literal text re-evaluated: %q", got)
	}
	second, err := syntax.Format(reparsed, syntax.FormatOptions{})
	if err != nil || !bytes.Equal(out, second) {
		t.Fatalf("comment formatting not stable: %v\n%s", err, second)
	}
}

func TestParserBoundsDeepNesting(t *testing.T) {
	source := "library {format=1}\n" + strings.Repeat("node \"x\" {\n", 130) + strings.Repeat("}\n", 130)
	_, err := syntax.Parse("deep.mogent", []byte(source))
	if err == nil || !strings.Contains(err.Error(), "nesting exceeds") {
		t.Fatalf("depth error=%v", err)
	}
}

func TestFormatterRejectsUnsupportedProgrammaticValues(t *testing.T) {
	f := parse(t, "plan {format=1}\ninclude {text=\"hello\"}")
	field(t, f.Entries[1], "text").Kind = "expression"
	if _, err := syntax.Format(f, syntax.FormatOptions{}); err == nil || !strings.Contains(err.Error(), "unsupported value kind") {
		t.Fatalf("format error=%v", err)
	}
}
