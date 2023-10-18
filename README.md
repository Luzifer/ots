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

To shorten the README this documentation has been moved to the Wiki:
https://github.com/Luzifer/ots/wiki/Customization

## Creating secrets through CLI / scripts

As `ots` is designed to never let the server know the secret you are sharing you should not just send the plain secret to it though it is possible.

### OTS-CLI

Download OTS-CLI from the [Releases](https://github.com/Luzifer/ots/releases) section of the repo or build it yourself having a Go toolchain available from the `./cmd/ots-cli` directory.

Afterwards you can just create and fetch secrets:

```console
# echo "my password" | ots-cli create
INFO[0000] reading secret content...                    
INFO[0000] creating the secret...                       
INFO[0000] secret created, see URL below                 expires-at="2023-10-16 16:33:27.422174121 +0000 UTC"
https://ots.fyi/#37a75a7f-0c2d-4ae6-bcca-4208b6d596ab%7CHGShVWm5umv4lmswfM73

# ots-cli fetch 'https://ots.fyi/#37a75a7f-0c2d-4ae6-bcca-4208b6d596ab%7CHGShVWm5umv4lmswfM73'
INFO[0000] fetching secret...                           
my password
```

To set the instance to send the secret to or to attach files see `ots-cli create --help` and to define where downloaded files are stored see `ots-cli fetch --help`.

Both commands can be used in scripts:
- `create` reads from `STDIN` or the specified file and yields the URL to `STDOUT`
- `fetch` prints the secret to `STDOUT` and stores files to the given directory
- both sends logs to `STDERR` which you can disable (`--log-level=fatal`) or ignore in your script

In case your instance needs credentials to use the `/api/create` endpoint you can pass them to OTS-CLI like you would do with curl:
- `ots-cli create --instance ... -u myuser:mypass` for basic-auth
- `ots-cli create --instance ... -H 'Authorization: Token abcde'` for token-auth (you can set any header you need, just repeat `-H ...`)

### Bash: Sharing an encrypted secret (strongly recommended!)

This is slightly more complex as you first need to encrypt your secret before sending it to the API but in this case you can be sure the server will in no case be able to access the secret. Especially if you are using ots.fyi (my public hosted instance) you should not trust me with your secret but use an encrypted secret:

```console
# echo "my password" | openssl aes-256-cbc -base64 -pass pass:mypass -pbkdf2 -iter 300000 -md sha512
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
    translators: []        # Array of Github usernames who "own" the translation
                           # and are pinged in the translation issue when there
                           # are translations missing (as of new features being
                           # added or features being improved). Add your username
                           # to this array to get pinged by the bot when stuff
                           # needs to be translated.
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
