package tplfunc

import (
	"errors"
	"strings"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type stringer string

func (s stringer) String() string { return string(s) }

func TestDict(t *testing.T) {
	t.Parallel()

	dict := dict(
		"string", "value",
		[]byte("bytes"), 42,
		errors.New("err"), true,
		stringer("stringer"), []string{"a", "b"},
		"missing",
	)

	assert.Equal(t, "value", dict["string"])
	assert.Equal(t, 42, dict["bytes"])
	assert.Equal(t, true, dict["err"])
	assert.Equal(t, []string{"a", "b"}, dict["stringer"])
	assert.Empty(t, dict["missing"])
}

func TestGetFuncs(t *testing.T) {
	t.Parallel()

	funcs := FuncMap()

	for _, name := range []string{"dict", "list", "mustMergeOverwrite", "mustToJson", "replace"} {
		assert.Contains(t, funcs, name)
	}
}

func TestList(t *testing.T) {
	t.Parallel()

	assert.Equal(t, []any{"a", 1, true}, list("a", 1, true))
}

func TestMustMergeOverwrite(t *testing.T) {
	t.Parallel()

	left := map[string]any{
		"keep":     "left",
		"override": "left",
	}
	right := map[string]any{
		"add":      "right",
		"override": "right",
	}

	out := mustMergeOverwrite(left, nil, right)

	assert.Equal(t, map[string]any{
		"add":      "right",
		"keep":     "left",
		"override": "right",
	}, out)
	assert.Equal(t, "left", left["override"])
}

func TestMustToJSON(t *testing.T) {
	t.Parallel()

	out, err := mustToJSON(map[string]any{"key": "value"})
	require.NoError(t, err)

	assert.JSONEq(t, `{"key":"value"}`, out)
}

func TestMustToJSONError(t *testing.T) {
	t.Parallel()

	_, err := mustToJSON(make(chan string))
	require.Error(t, err)
	assert.ErrorContains(t, err, "marshalling data")
}

func TestReplace(t *testing.T) {
	t.Parallel()

	assert.Equal(t, `can\'t stop`, replace("'", `\'`, "can't stop"))
}

func TestTemplateFuncs(t *testing.T) {
	t.Parallel()

	tpl, err := template.New("test").Funcs(FuncMap()).Parse(strings.TrimSpace(`
{{- range list "one" "two" }}{{ . }} {{ end -}}
{{- $merged := mustMergeOverwrite (dict "a" "left" "b" "left") (dict "b" "right") -}}
{{- $merged | mustToJson | replace "\"" "'" -}}
`))
	require.NoError(t, err)

	out := new(strings.Builder)
	require.NoError(t, tpl.Execute(out, nil))

	assert.Equal(t, "one two {'a':'left','b':'right'}", out.String())
}
