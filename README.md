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
  - `SECRET_EXPIRY` - Expiry of the keys in seconds (Default `0` = no expiry)
- `redis` - Storing the secrets in a hash under one key
  - `REDIS_URL` - Redis connection string `tcp://auth:PWD@HOST:PORT/DB`
  - `REDIS_KEY` - Key prefix to store the keys under (Default `io.luzifer.ots`)
  - `SECRET_EXPIRY` - Expiry of the keys in seconds (Default `0` = no expiry)

## Creating secrets through CLI / scripts

As `ots` is designed to never let the server know the secret you are sharing you should not just send the plain secret to it though it is possible.

### Sharing an encrypted secret (strongly recommended!)

This is slightly more complex as you first need to encrypt your secret before sending it to the API but in this case you can be sure the server will in no case be able to access the secret. Especially if you are using ots.fyi (my public hosted instance) you should not trust me with your secret but use an encrypted secret:

```console
# echo "my password" | openssl aes-256-cbc -base64 -pass pass:mypass -md md5
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

--------

### Detailed usage - OTS

Download the latest stable release depending upon the operating system [here](https://github.com/Luzifer/ots/releases) ([https://github.com/Luzifer/ots/releases](https://github.com/Luzifer/ots/releases))

#### For Windows User -

Download ots binary available at URL  -[https://github.com/Luzifer/ots/releases/download/v0.21.0/ots_windows_amd64.zip](https://github.com/Luzifer/ots/releases/download/v0.21.0/ots_windows_amd64.zip)

![ots2.PNG](images/ots2.png)

Double click to unzip/extract the `ots_windows_amd64.exe` file.

Run/execute the `ots_windows_amd64.exe` file. Notice a command prompt as shown below - 

![ots3.PNG](images/ots3.png)

Open your favourite browser and navigate to `http://localhost:3000`

![ots4.PNG](images/ots4.png)

Type your secret in the input field as shown below `secret123`

![ots5.PNG](images/ots5.png)

Click on **Create the secret!**

![ots6.PNG](images/ots6.png)

- You can share the above URL `http://localhost:3000/#2f6e1abc......` with the person you wanted to share the password with.
    
    **Note** - This URL is only valid for one use. Once you retrieve the password the URL will become invalid 
    

Ask the person to navigate to the URL `http://localhost:3000/#2f6e1abc......`

Click on **Show me the secret!** 

![ots7.PNG](images/ots7.png)

#### For Linux users -

Download file available at URL  - [https://github.com/Luzifer/ots/releases/download/v0.21.0/ots_linux_amd64.tar.gz](https://github.com/Luzifer/ots/releases/download/v0.21.0/ots_linux_amd64.tar.gz)

- Next Run below command

```jsx
tar -zxvf ots_linux_amd64.tar.gz
```

- Next,

```jsx
./ots_linux_amd64
```

- Browse URL http://localhost:3000

Use the OTS application to share secrets

To add an extra layer of security, encrypt the secret using OpenSSL before creating the secret in OTS

#### Encrypting secrets using a passphrase using OpenSSL

```bash
**echo "Secret124" | openssl enc -aes-256-cbc -a -pbkdf2**
```

Output - U2FsdGVkX19N6IuZJ9POuXkgQe++Kdr/+nsuGc+HqbA=

Now, pass this encrypted token into OTS application and share the OTS URL with the relevant person.

![Untitled](images/Untitled.png)

Once the receiver opens the URL and clicks on `show me the secret!` button. They will notice the encrypted token.

#### Decrypting token using OpenSSL

```bash
**echo "U2FsdGVkX19N6IuZJ9POuXkgQe++Kdr/+nsuGc+HqbA=" | openssl enc -aes-256-cbc -d -a -pbkdf2**
```

![Untitled](images/Untitled%201.png)


-------

## Localize to your own language

If you want to help translating the application to your own language please download the [`en.json`](https://github.com/Luzifer/ots/blob/master/src/langs/en.json) file from this repository and translate the strings inside. Afterwards please [open an issue](https://github.com/Luzifer/ots/issues/new) and attach your translation including the information which language you translated the strings into.

Of course you also could open a pull-request to add the new file to the `src/langs` folder. In this case please also edit the `langs.js` file and add your translation.

Same goes with when you're finding translation errors: Just open an issue and let me know!
