![](https://badges.fyi/github/license/Luzifer/ots)
![](https://badges.fyi/github/latest-release/Luzifer/ots)
![](https://badges.fyi/github/downloads/Luzifer/ots)
[![Go Report Card](https://goreportcard.com/badge/github.com/Luzifer/ots)](https://goreportcard.com/report/github.com/Luzifer/ots)

# Luzifer / OTS

`ots` is a one-time-secret sharing platform. The secret is encrypted with a symmetric 256bit AES encryption in the browser before being sent to the server. Afterwards an URL containing the ID of the secret and the password is generated. The password is never sent to the server so the server will never be able to decrypt the secrets it delivers with a reasonable effort. Also the secret is immediately deleted on the first read.

## Features

- AES 256bit encryption
- Server does never get the password
- Secret is deleted on first read

## Setup

- Download the [release](https://github.com/Luzifer/ots/releases)
- Start it and you can access the server on http://localhost:3000/
- Consult `./ots --help` for more options
- See [Wiki](https://github.com/Luzifer/ots/wiki) for a more detailed overview

For a better setup you can choose the backend which is used to store the secrets:

- `mem` - In memory storage (wiped on restart of the daemon)
- `redis` - Storing the secrets in a hash under one key
  - `REDIS_URL` - Redis connection string `redis://USR:PWD@HOST:PORT/DB`  
    (pre Redis v6 use `auth` as user, afterwards use a user available in your ACLs)
  - `REDIS_KEY` - Key prefix to store the keys under (Default `io.luzifer.ots`)
- Common options
  - `SECRET_EXPIRY` - Expiry of the keys in seconds (Default `0` = no expiry)

### Customization

In order to be adjustable to your needs there are some ways to customize your OTS setup. All of those require you to create a YAML file containing the definitions of your customizations and to load this file through the `--customize=path/to/customize.yaml`:

```yaml
# Override the app-icon, present a path to the image to use, if unset
# or empty the default FontAwesome icon will be displayed. Recommended
# is a height of 30px.
appIcon: ''

# Override the app-title, if unset or empty the default app-title
# "OTS - One Time Secrets" will be used
appTitle: ''

# Disable display of the app-title (for example if you included the
# title within the appIcon)
disableAppTitle: false

# Disable the dropdown and the API functionality to override the
# secret expiry
disableExpiryOverride: false

# Disable the footer linking back to the project. If you disable it
# please consider a donation to support the project.
disablePoweredBy: false

# Disable the button to display and the generation of the QR-Code
# for the secret URL
disableQRSupport: false

# Disable the switcher for dark / light theme in the top right corner
# for example if your custom theme does not support two themes.
disableThemeSwitcher: false

# Override the choices to be displayed in the expiry dropdown. Values
# are given in seconds and the order of the values controls the order
# in the dropdown.
expiryChoices: [300, ...]

# Custom path to override embedded resources. You can override any
# file present in the `frontend` directory (which is baked into the
# binary during compile-time). You also can add new files (for
# example the appIcon given above). Those files are available at the
# root of the application (i.e. an app.png would be served at
# https://ots.example.com/app.png).
overlayFSPath: /path/to/ots-customization

# Switch to formal translations for languages having those defined.
# Languages not having a formal version will still display the normal
# translations in the respective language.
useFormalLanguage: false

# Define which file types are selectable by the user when uploading
# files to attach. This fuels the `accept` attribute of the file
# select and requires the same format. Pay attention this is not
# suited as a security measure as this is purely a frontend
# implementation and can be circumvented.
# https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/file#accept
acceptedFileTypes: ''

# Disable the file attachment functionality alltogether
disableFileAttachment: false

# Define how big all attachments might be in bytes. Leave it set to
# zero to use the internal limit of 64 MiB (which is there to ensure
# the encrypted object does not cause the frontend to break).
maxAttachmentSizeTotal: 0
```

To override the styling of the application have a look at the [`src/style.scss`](./src/style.scss) file how the theme of the application is built and present the compiled `app.css` in the `overlayFSPath`.

After modifying files in the `overlayFSPath` make sure to restart the application as otherwise the file integrity hashes are no longer matching and your resources will be blocked by the browsers.

If you want to disable secret creation for users not logged into your company SSO you can apply an ACL on the `/api/create` and `/api/isWritable` endpoints to allow access to them only for logged in users. This will also disable the secret-creation interface for all not having access to the `/api/isWritable` endpoint.

## Creating secrets through CLI / scripts

As `ots` is designed to never let the server know the secret you are sharing you should not just send the plain secret to it though it is possible.

### Sharing an encrypted secret (strongly recommended!)

This is slightly more complex as you first need to encrypt your secret before sending it to the API but in this case you can be sure the server will in no case be able to access the secret. Especially if you are using ots.fyi (my public hosted instance) you should not trust me with your secret but use an encrypted secret:

```console
# echo "my password" | openssl aes-256-cbc -base64 -pass pass:mypass -iter 300000 -md sha512
U2FsdGVkX18wJtHr6YpTe8QrvMUUdaLZ+JMBNi1OvOQ=

# curl -X POST -H 'content-type: application/json' -i -s -d '{"secret": "U2FsdGVkX18wJtHr6YpTe8QrvMUUdaLZ+JMBNi1OvOQ="}' https://ots.fyi/api/create
HTTP/2 201
server: nginx
date: Wed, 29 Jan 2020 14:08:54 GMT
content-type: application/json
content-length: 68
cache-control: no-cache

{"secret_id":"5e0065ee-5734-4548-9fd3-bb0bcd4c899d","success":true}
```

You will now need to supply the web application with the password in addition to the ID of the secret: `https://ots.fyi/#5e0065ee-5734-4548-9fd3-bb0bcd4c899d|mypass`

In this case due to how browsers are handling hashes in URLs (the part after the `#`) the only URL the server gets to know is `https://ots.fyi/` which loads the frontend. Afterwards the Javascript executed in the browser fetches the encrypted secret at the given ID and decrypts it with the given password (in this case `mypass`). I will not be able to tell the content of your secret and just see the AES 256bit encrypted content.

You can find a script [`cli_create.sh`](cli_create.sh) in this repo demonstrating the creation of the secret with all steps.

### Sharing the plain secret

```console
# curl -X POST -H 'content-type: application/json' -i -s -d '{"secret": "my password"}' https://ots.fyi/api/create

HTTP/2 201
server: nginx
date: Wed, 29 Jan 2020 14:02:42 GMT
content-type: application/json
content-length: 68
cache-control: no-cache

{"secret_id":"1cb08e53-46b9-4f21-bbd9-f1eea1594ad9","success":true}
```

You can then use the URL `https://ots.fyi/#1cb08e53-46b9-4f21-bbd9-f1eea1594ad9` to access the secret.

## Localize to your own language

If you want to help translating the application to your own language please see the [`i18n.yaml`](https://github.com/Luzifer/ots/blob/master/i18n.yaml) file from this repository and translate the English strings inside. Afterwards please [open an issue](https://github.com/Luzifer/ots/issues/new) and attach your translation including the information which language you translated the strings into.

Of course you also could open a pull-request to add the new translations to the `i18n.yaml` file.

Same goes with when you're finding translation errors: Just open an issue and let me know!

The format for the `i18n.yaml` is as follows:
```yaml
reference:                 # Reference strings (English)
  deeplLanguage: en        # Source language for DeepL automated translations
  languageKey: en          # Browser language to use this translation for
  translations: {}         # Map of translation keys to their translations

translations:              # Translations into other languages
  de:                      # Identifier for the language, used as `languageKey`
    deeplLanguage: de      # Target language for DeepL automated translations
    translations: {}       # Informal / base translations for the language.
                           # Missing keys will be loaded from the `reference`
                           # and therefore get displayed in English. Missing
                           # keys can be generated through DeepL through the
                           # translation tool included in `ci/translate` but
                           # will have low quality as partial sentences or
                           # even only words lack the context for the
                           # translation
    formalTranslations: {} # Formal translations for the language (these will
                           # be merged over the `translations` for this language
                           # so you don't have to copy keys being equal in formal
                           # and informal translation.)
```
