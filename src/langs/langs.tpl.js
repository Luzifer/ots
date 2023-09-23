// Auto-Generated, do not edit!

const switchFormal = (formal, informal) => window.useFormalLanguage ? formal : informal

export default {
{{- range $lang, $translation := .Translations -}}
{{- if .FormalTranslations }}
  '{{ $lang }}': switchFormal(
    JSON.parse('{{ mustMergeOverwrite (dict) .Translations .FormalTranslations | mustToJson | replace "'" "\\'" }}'),
    JSON.parse('{{ .Translations | mustToJson | replace "'" "\\'" }}'),
  ),
{{ else }}
  '{{ $lang }}': JSON.parse('{{ .Translations | mustToJson | replace "'" "\\'" }}'),
{{ end -}}
{{- end }}
}
