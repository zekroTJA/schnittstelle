package schnittstelle

import (
	"strings"
	"testing"
)

func TestExtract(t *testing.T) {
	expected := []string{
		"SimpleRefReceiver()",
		"SimpleValueReceiver()",
		"SimpleUnnamedRefReceiver()",
		"SimpleUnnamedValueReceiver()",
		"ParamsInline(a string, b int)",
		"ParamsMultiline(a string, b int, c interface{})",
		"ReturnsInline() (int, string)",
		"ReturnsInlineNamed() (a, b int, c string)",
		"ReturnsMultilineNamed() (a, b int, c string, d error)",
		"InlineEmpty()",
		"InlineFunc(a string) bool",
		"ReturnsMultilineNamedWithComment() (a, b int, d error)",
	}

	results, err := Extract("Example", "test", 5)
	if err != nil {
		t.Fatal(err)
	}

	compare(t, expected, results)
}

// See Issue #1
func TestExtractRootPathNotExists(t *testing.T) {
	_, err := Extract("Example", "not_exists", 5)
	if err == nil {
		t.Fatal("path error not returned")
	}

	_, err = Extract("Example", "go.mod", 5)
	if err != ErrPathIsFile {
		t.Fatal("wrong error returned:", err)
	}
}

// --- helpers ---

func compare(t *testing.T, expected, results []string) {
	t.Helper()

	contains := func(s []string, v string) bool {
		for _, c := range s {
			if c == v {
				return true
			}
		}
		return false
	}

	var extra, missing []string
	for _, e := range expected {
		if !contains(results, e) {
			missing = append(missing, e)
		}
	}
	for _, e := range results {
		if !contains(expected, e) {
			extra = append(extra, e)
		}
	}

	if len(extra) == 0 && len(missing) == 0 {
		return
	}

	t.Fatalf("Value missmatch:\n\nExtra:\n\t%s\n\nMissing:\n\t%s\n\n",
		strings.Join(extra, "\n\t"),
		strings.Join(missing, "\n\t"))
}
