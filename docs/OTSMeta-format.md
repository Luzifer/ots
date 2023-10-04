> OTS uses two different formats to store secrets under the hood. Both of them can be read and written by the frontend implementation as well as by the `ots-cli` application.

## Simple Format

The simple format is the format used by OTS since day one and is the, well, most simple one. It only consists of the secret encrypted using OpenSSL AES-256-CBC compatible encryption. This format is preferred for backwards compatibility when no other reasons require the use of the OTS-Meta format.

```console
# openssl enc -aes-256-cbc -pbkdf2 -md sha512 -iter 300000 -pass pass:12345678 -a -A <<<"I'm a secret"
U2FsdGVkX19G3GuIw3LGM0PVQmavPU/LnWvJhcLeYvs=

# curl -H 'Content-Type: application/json' -d '{"secret": "U2FsdGVkX19G3GuIw3LGM0PVQmavPU/LnWvJhcLeYvs="}' https://ots.fyi/api/create
{"success":true,"expires_at":"2023-10-11T19:45:01.315587714Z","secret_id":"bbd53ec5-8ee9-4df5-a630-9561313a348a"}

# ots-cli fetch "https://ots.fyi/#bbd53ec5-8ee9-4df5-a630-9561313a348a%7C12345678"
INFO[0000] fetching secret...
I'm a secret
```

## OTSMeta Format

The OTSMeta format was first introduced in `v1.9.0` of OTS together with the possibility to attach files to the secret. It contains structured data with a banner to differentiate between a simple JSON shared through OTS and the OTSMeta format. The OTSMeta structure itself is a simple JSON document containing a secret and a number of attachments having their contents base64 encoded:

```json
{
  "secret": "I'm a secret",
  "attachments": [
    {
      "name": "file.txt",
      "type": "text/plain",
      "data": "SSdtIGZpbGUgY29udGVudAo="
    }
  ]
}
```

This structure is prefixed with the Banner `OTSMeta` and then shared the same way as a simple secret would be:

```console
# ots-cli create -f file.txt <<<"I'm a secret"
INFO[0000] reading secret content...
INFO[0000] attaching file...                             file=file.txt
INFO[0000] creating the secret...
INFO[0000] secret created, see URL below                 expires-at="2023-10-11 19:52:30.816059504 +0000 UTC"
https://ots.fyi/#6a6be08c-97d7-4970-a202-5bb6964460d8%7CwNUURZ0LRrQAhaczdZfj

# curl -sS https://ots.fyi/api/get/6a6be08c-97d7-4970-a202-5bb6964460d8 | jq -r .secret >/tmp/secret.bin
# openssl enc -aes-256-cbc -pbkdf2 -md sha512 -iter 300000 -pass pass:wNUURZ0LRrQAhaczdZfj -a -A -d </tmp/secret.bin
OTSMeta{"secret":"I'm a secret\n","attachments":[{"name":"file.txt","type":"text/plain; charset=utf-8","data":"SSdtIGZpbGUgY29udGVudAo="}]}
```

Of course it's also possible to share a simple secret in OTSMeta format but the recommended way would be to omit the OTSMeta wrapping:
```
OTSMeta{"secret":"I'm a secret"}
```

When programmatically reading secrets you therefore need to check whether the secret starts with `OTSMeta` and decode the remaining as a JSON document and if it does not just use all the content as the secret.
