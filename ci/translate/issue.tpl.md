---
title: Missing Translations
---
> As a developer I want my application to have correct translations in all available languages and not to have them to experience mixed translations in their native language and English.

**In order to achieve this we need to fix the following missing translations.**

To help translating please either **create a pull-request** updating the `i18n.yaml` in the root of the repository and add the missing translations to the corresponding language or **just leave a comment** below and ping @Luzifer in your comment. He then will integrate the new translation strings and mark your comment hidden after this issue has been automatically updated (kind of a to-do list for translations until we have something better in place).

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
