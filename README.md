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

For a better setup you can choose the backend which is used to store the secrets:

- `mem` - In memory storage (wiped on restart of the daemon)
- `redis` - Storing the secrets in a hash under one key
  - `REDIS_URL` - Redis connection string `tcp://auth:PWD@HOST:PORT/DB`
  - `REDIS_EXPIRY` - Expiry of the keys in seconds (Default `0` = no expiry)
  - `REDIS_KEY` - Key prefix to store the keys under (Default `io.luzifer.ots`)

**Hint:** Starting in `v0.7.0` the secrets in Redis are no longer stored in a hash but in own keys. This allows for individual expiry. At the first start of `v0.7.0` the old data will be migrated automatically and afterwards be subject of expiry if you set `REDIS_EXPIRY`. My hosted instance uses an expiry of 90d (= 7776000s).

## Localize to your own language

If you want to help translating the application to your own language please download the [`en.json`](https://github.com/Luzifer/ots/blob/master/src/langs/en.json) file from this repository and translate the strings inside. Afterwards please [open an issue](https://github.com/Luzifer/ots/issues/new) and attach your translation including the information which language you translated the strings into.

Of course you also could open a pull-request to add the new file to the `src/langs` folder. In this case please also edit the `langs.js` file and add your translation.

Same goes with when you're finding translation errors: Just open an issue and let me know!
