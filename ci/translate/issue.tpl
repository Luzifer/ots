---
title: Missing Translations
---
> As a developer I want my application to have correct translations in all available languages and not to have them to experience mixed translations in their native language and English.

In order to achieve this we need to fix the following missing translations:

{{ range $lang, $translation := .Translations -}}
{{ if MissingTranslations $lang -}}
### Language: `{{ $lang }}`

Please add the following translations:
{{ range MissingTranslations $lang }}
- `{{ . }}`
  > {{ English . }}
{{ end }}
_{{ Ping $lang }}_

{{ end -}}
{{ end -}}
